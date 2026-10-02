package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gophish/gophish/dialer"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gophish/gophish/util"
	"github.com/jordan-wright/email"
)

// maxStylesheetSize caps how much we download for a single inlined stylesheet.
const maxStylesheetSize = 5 << 20 // 5 MiB

// maxEmbedSize caps a single resource embedded as a data: URI, and embedBudget
// caps the total embedded across a single import so the page stays reasonable.
const (
	maxEmbedSize = 2 << 20  // 2 MiB per resource
	embedBudget  = 12 << 20 // 12 MiB total
)

var (
	cssURLRe    = regexp.MustCompile(`(?i)url\(\s*['"]?([^'")]+?)['"]?\s*\)`)
	cssImportRe = regexp.MustCompile(`(?i)@import\s+['"]([^'"]+)['"]`)
	cssEndTagRe = regexp.MustCompile(`(?i)</style`)
)

// cloneHTTPClient builds the HTTP client used to fetch a site and its
// resources, going through the SSRF-restricted dialer.
func cloneHTTPClient() *http.Client {
	restrictedDialer := dialer.Dialer()
	return &http.Client{
		Transport: &http.Transport{
			DialContext: restrictedDialer.DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

type cloneRequest struct {
	URL              string `json:"url"`
	IncludeResources bool   `json:"include_resources"`
	// EmbedImages additionally downloads images (and, when the CSS is inlined,
	// the resources it references) and embeds them as data: URIs so the landing
	// page is fully self-contained.
	EmbedImages bool `json:"embed_images"`
}

func (cr *cloneRequest) validate() error {
	if cr.URL == "" {
		return errors.New("No URL Specified")
	}
	return nil
}

type cloneResponse struct {
	HTML string `json:"html"`
}

type emailResponse struct {
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Subject string `json:"subject"`
}

// ImportGroup imports a CSV of group members
func (as *Server) ImportGroup(w http.ResponseWriter, r *http.Request) {
	ts, err := util.ParseCSV(r)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Error parsing CSV"}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, ts, http.StatusOK)
}

// ImportEmail allows for the importing of email.
// Returns a Message object
func (as *Server) ImportEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusBadRequest)
		return
	}
	ir := struct {
		Content      string `json:"content"`
		ConvertLinks bool   `json:"convert_links"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&ir)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Error decoding JSON Request"}, http.StatusBadRequest)
		return
	}
	e, err := email.NewEmailFromReader(strings.NewReader(ir.Content))
	if err != nil {
		log.Error(err)
	}
	// If the user wants to convert links to point to
	// the landing page, let's make it happen by changing up
	// e.HTML
	if ir.ConvertLinks {
		d, err := goquery.NewDocumentFromReader(bytes.NewReader(e.HTML))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		d.Find("a").Each(func(i int, a *goquery.Selection) {
			a.SetAttr("href", "{{.URL}}")
		})
		h, err := d.Html()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		e.HTML = []byte(h)
	}
	er := emailResponse{
		Subject: e.Subject,
		Text:    string(e.Text),
		HTML:    string(e.HTML),
	}
	JSONResponse(w, er, http.StatusOK)
}

// ImportSite allows for the importing of HTML from a website
// Without "include_resources" set, it will merely place a "base" tag
// so that all resources can be loaded relative to the given URL.
func (as *Server) ImportSite(w http.ResponseWriter, r *http.Request) {
	cr := cloneRequest{}
	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Error decoding JSON Request"}, http.StatusBadRequest)
		return
	}
	if err = cr.validate(); err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	client := cloneHTTPClient()
	resp, err := client.Get(cr.URL)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	baseURL, err := url.Parse(cr.URL)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	// Optionally embed images (and CSS-referenced resources) as data: URIs so
	// the landing page is fully self-contained. The embedder carries a shared
	// size budget across stylesheets and images.
	var embedder *resourceEmbedder
	if cr.EmbedImages {
		embedder = newResourceEmbedder(client)
	}
	// When requested, download each linked stylesheet and inline it so the
	// landing page keeps its styling locally, without depending on the origin.
	if cr.IncludeResources {
		inlineStylesheets(d, baseURL, client, embedder)
	}
	if cr.EmbedImages {
		embedImages(d, baseURL, embedder)
	}
	// Make sure a base href points at the original site so the remaining
	// relative resources (images, scripts, etc.) resolve against it rather
	// than the host serving the landing page. If the page already ships a base
	// (often root-relative, e.g. <base href="/"> on SPAs), rewrite it to an
	// absolute origin URL instead of leaving it to resolve against our host.
	baseSel := d.Find("head base").First()
	if baseSel.Length() == 0 {
		d.Find("head").PrependHtml(fmt.Sprintf("<base href=\"%s\">", cr.URL))
	} else if href, ok := baseSel.Attr("href"); ok {
		if ref, perr := url.Parse(strings.TrimSpace(href)); perr == nil {
			if abs := baseURL.ResolveReference(ref); abs.Scheme == "http" || abs.Scheme == "https" {
				baseSel.SetAttr("href", abs.String())
			}
		}
	}
	forms := d.Find("form")
	forms.Each(func(i int, f *goquery.Selection) {
		// We'll want to store where we got the form from
		// (the current URL)
		url := f.AttrOr("action", cr.URL)
		if !strings.HasPrefix(url, "http") {
			url = fmt.Sprintf("%s%s", cr.URL, url)
		}
		f.PrependHtml(fmt.Sprintf("<input type=\"hidden\" name=\"__original_url\" value=\"%s\"/>", url))
	})
	h, err := d.Html()
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	cs := cloneResponse{HTML: h}
	JSONResponse(w, cs, http.StatusOK)
}

// inlineStylesheets downloads every linked stylesheet in the document and
// replaces the <link> with an inline <style> element, rewriting relative
// url()/@import references to absolute URLs so fonts and background images
// still resolve. On a fetch error the <link> is kept but its href is made
// absolute as a fallback.
func inlineStylesheets(d *goquery.Document, base *url.URL, client *http.Client, embedder *resourceEmbedder) {
	d.Find(`link[rel="stylesheet"]`).Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || strings.TrimSpace(href) == "" {
			return
		}
		ref, err := url.Parse(strings.TrimSpace(href))
		if err != nil {
			return
		}
		abs := base.ResolveReference(ref)
		if abs.Scheme != "http" && abs.Scheme != "https" {
			return
		}
		css, err := fetchResource(client, abs.String())
		if err != nil {
			log.Warnf("import: could not fetch stylesheet %s: %v", abs.String(), err)
			s.SetAttr("href", abs.String())
			return
		}
		css = rewriteCSSURLs(css, abs, embedder)
		// Prevent the stylesheet contents from closing the <style> element.
		css = cssEndTagRe.ReplaceAllString(css, `<\/style`)
		var b strings.Builder
		b.WriteString("<style")
		if media := strings.TrimSpace(s.AttrOr("media", "")); media != "" {
			b.WriteString(fmt.Sprintf(" media=%q", media))
		}
		b.WriteString(">\n")
		b.WriteString(css)
		b.WriteString("\n</style>")
		s.ReplaceWithHtml(b.String())
	})
}

// fetchResource downloads a resource (bounded in size) and returns its body.
func fetchResource(client *http.Client, u string) (string, error) {
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxStylesheetSize))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// rewriteCSSURLs resolves relative url(...) and @import "..." references in a
// stylesheet against the stylesheet's own URL so they keep working once the
// CSS has been moved (inlined) into the landing page.
func rewriteCSSURLs(css string, base *url.URL, embedder *resourceEmbedder) string {
	css = cssURLRe.ReplaceAllStringFunc(css, func(m string) string {
		abs := absoluteResourceURL(cssURLRe.FindStringSubmatch(m)[1], base)
		if abs == "" {
			return m
		}
		// Embed the referenced resource (font/background image) when asked and
		// it isn't already a data: URI; otherwise keep the absolute URL.
		if embedder != nil && !strings.HasPrefix(abs, "data:") {
			if uri, ok := embedder.embed(abs); ok {
				return "url(" + uri + ")"
			}
		}
		return "url(" + abs + ")"
	})
	css = cssImportRe.ReplaceAllStringFunc(css, func(m string) string {
		sub := cssImportRe.FindStringSubmatch(m)
		if abs := absoluteResourceURL(sub[1], base); abs != "" {
			return `@import "` + abs + `"`
		}
		return m
	})
	return css
}

// absoluteResourceURL resolves a raw reference against base, leaving data: URIs
// and already-absolute http(s) URLs untouched. Returns "" when it can't parse.
func absoluteResourceURL(raw string, base *url.URL) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "data:") {
		return raw
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return raw
	}
	ref, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}

// resourceEmbedder downloads resources and returns them as data: URIs, bounded
// by a per-resource size cap and a shared total budget so an import can't
// produce an unbounded landing page.
type resourceEmbedder struct {
	client *http.Client
	budget int
}

func newResourceEmbedder(client *http.Client) *resourceEmbedder {
	return &resourceEmbedder{client: client, budget: embedBudget}
}

// embed fetches u and returns a data: URI for it. ok is false (and the caller
// should keep the absolute URL) on any error, when the resource exceeds the
// per-resource cap, or when the total budget is exhausted.
func (e *resourceEmbedder) embed(u string) (string, bool) {
	if e == nil || e.budget <= 0 {
		return "", false
	}
	resp, err := e.client.Get(u)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	limit := maxEmbedSize
	if e.budget < limit {
		limit = e.budget
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(limit)+1))
	if err != nil || len(body) == 0 || len(body) > limit {
		return "", false
	}
	e.budget -= len(body)
	ct := resp.Header.Get("Content-Type")
	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = ct[:i]
	}
	ct = strings.TrimSpace(ct)
	if ct == "" {
		ct = mime.TypeByExtension(path.Ext(u))
	}
	if ct == "" {
		ct = "application/octet-stream"
	}
	return fmt.Sprintf("data:%s;base64,%s", ct, base64.StdEncoding.EncodeToString(body)), true
}

// embedImages downloads <img> sources and replaces them with data: URIs. On
// failure it falls back to an absolute URL so the image still loads remotely.
func embedImages(d *goquery.Document, base *url.URL, embedder *resourceEmbedder) {
	d.Find("img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" || strings.HasPrefix(src, "data:") {
			return
		}
		ref, err := url.Parse(src)
		if err != nil {
			return
		}
		abs := base.ResolveReference(ref)
		if abs.Scheme != "http" && abs.Scheme != "https" {
			return
		}
		if uri, ok := embedder.embed(abs.String()); ok {
			s.SetAttr("src", uri)
			// A remote srcset would otherwise be preferred over the embedded src.
			s.RemoveAttr("srcset")
		} else {
			s.SetAttr("src", abs.String())
		}
	})
}
