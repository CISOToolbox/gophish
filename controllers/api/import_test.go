package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gophish/gophish/dialer"
	"github.com/gophish/gophish/models"
)

func makeImportRequest(ctx *testContext, allowedHosts []string, url string) *httptest.ResponseRecorder {
	orig := dialer.DefaultDialer.AllowedHosts()
	dialer.SetAllowedHosts(allowedHosts)
	req := httptest.NewRequest(http.MethodPost, "/api/import/site",
		bytes.NewBuffer([]byte(fmt.Sprintf(`
			{
				"url" : "%s"
			}
		`, url))))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ctx.apiServer.ImportSite(response, req)
	dialer.SetAllowedHosts(orig)
	return response
}

func TestDefaultDeniedImport(t *testing.T) {
	ctx := setupTest(t)
	metadataURL := "http://169.254.169.254/latest/meta-data/"
	response := makeImportRequest(ctx, []string{}, metadataURL)
	expectedCode := http.StatusBadRequest
	if response.Code != expectedCode {
		t.Fatalf("incorrect status code received. expected %d got %d", expectedCode, response.Code)
	}
	got := &models.Response{}
	err := json.NewDecoder(response.Body).Decode(got)
	if err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	if !strings.Contains(got.Message, "upstream connection denied") {
		t.Fatalf("incorrect response error provided: %s", got.Message)
	}
}

func TestDefaultAllowedImport(t *testing.T) {
	ctx := setupTest(t)
	h := "<html><head></head><body><img src=\"/test.png\"/></body></html>"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, h)
	}))
	defer ts.Close()
	response := makeImportRequest(ctx, []string{}, ts.URL)
	expectedCode := http.StatusOK
	if response.Code != expectedCode {
		t.Fatalf("incorrect status code received. expected %d got %d", expectedCode, response.Code)
	}
}

func TestCustomDeniedImport(t *testing.T) {
	ctx := setupTest(t)
	h := "<html><head></head><body><img src=\"/test.png\"/></body></html>"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, h)
	}))
	defer ts.Close()
	response := makeImportRequest(ctx, []string{"192.168.1.1"}, ts.URL)
	expectedCode := http.StatusBadRequest
	if response.Code != expectedCode {
		t.Fatalf("incorrect status code received. expected %d got %d", expectedCode, response.Code)
	}
	got := &models.Response{}
	err := json.NewDecoder(response.Body).Decode(got)
	if err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	if !strings.Contains(got.Message, "upstream connection denied") {
		t.Fatalf("incorrect response error provided: %s", got.Message)
	}
}

func makeImportRequestResources(ctx *testContext, allowedHosts []string, target string, includeResources bool) *httptest.ResponseRecorder {
	orig := dialer.DefaultDialer.AllowedHosts()
	dialer.SetAllowedHosts(allowedHosts)
	body := fmt.Sprintf(`{"url":"%s","include_resources":%t}`, target, includeResources)
	req := httptest.NewRequest(http.MethodPost, "/api/import/site", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ctx.apiServer.ImportSite(response, req)
	dialer.SetAllowedHosts(orig)
	return response
}

func TestImportSiteInlineCSS(t *testing.T) {
	ctx := setupTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/style.css" {
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, ".x{color:red}body{background:url(bg.png)}")
			return
		}
		fmt.Fprint(w, `<html><head><link rel="stylesheet" href="/style.css"></head><body><p class="x">hi</p></body></html>`)
	}))
	defer ts.Close()

	response := makeImportRequestResources(ctx, []string{}, ts.URL, true)
	if response.Code != http.StatusOK {
		t.Fatalf("incorrect status code: expected 200 got %d", response.Code)
	}
	got := &cloneResponse{}
	if err := json.NewDecoder(response.Body).Decode(got); err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	if !strings.Contains(got.HTML, "<style") || !strings.Contains(got.HTML, "color:red") {
		t.Fatalf("stylesheet was not inlined: %s", got.HTML)
	}
	if strings.Contains(got.HTML, `rel="stylesheet"`) {
		t.Fatalf("the <link> stylesheet should have been replaced: %s", got.HTML)
	}
	wantBG := fmt.Sprintf("url(%s/bg.png)", ts.URL)
	if !strings.Contains(got.HTML, wantBG) {
		t.Fatalf("relative url() not rewritten to absolute (want %q): %s", wantBG, got.HTML)
	}
}

func TestImportSiteWithoutResourcesKeepsLink(t *testing.T) {
	ctx := setupTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><head><link rel="stylesheet" href="/style.css"></head><body>hi</body></html>`)
	}))
	defer ts.Close()

	response := makeImportRequestResources(ctx, []string{}, ts.URL, false)
	if response.Code != http.StatusOK {
		t.Fatalf("incorrect status code: expected 200 got %d", response.Code)
	}
	got := &cloneResponse{}
	if err := json.NewDecoder(response.Body).Decode(got); err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	if !strings.Contains(got.HTML, `rel="stylesheet"`) {
		t.Fatalf("without include_resources the <link> should be kept: %s", got.HTML)
	}
}
