package controllers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gophish/gophish/models"
)

var (
	reEdge    = regexp.MustCompile(`edg(?:a|ios)?/(\d+)`)
	reChrome  = regexp.MustCompile(`chrome/(\d+)`)
	reFirefox = regexp.MustCompile(`firefox/(\d+)`)
	reSafari  = regexp.MustCompile(`version/(\d+)[.\d]* safari`)
)

// isOutdatedBrowser reports whether a user-agent advertises a browser whose
// major version is below the configured minimum for its engine. A minimum of 0
// disables the check for that engine. Unknown user-agents are not judged here
// (the scanner user-agent list handles those).
func isOutdatedBrowser(ua string, s models.ScannerSettings) bool {
	parse := func(re *regexp.Regexp) (int, bool) {
		if m := re.FindStringSubmatch(ua); m != nil {
			if v, err := strconv.Atoi(m[1]); err == nil {
				return v, true
			}
		}
		return 0, false
	}
	below := func(v, min int) bool { return min > 0 && v < min }
	// Edge reports both "edg/" and "chrome/"; check it first so it isn't
	// mistaken for Chrome. Edge is part of the Chromium family, so it shares
	// the Chrome minimum.
	if v, ok := parse(reEdge); ok {
		return below(v, s.MinChrome)
	}
	if v, ok := parse(reChrome); ok {
		return below(v, s.MinChrome)
	}
	if v, ok := parse(reFirefox); ok {
		return below(v, s.MinFirefox)
	}
	if v, ok := parse(reSafari); ok {
		return below(v, s.MinSafari)
	}
	return false
}

// matchesKnownScanner reports whether the user-agent contains any of the
// configured scanner substrings (all lowercased).
func matchesKnownScanner(ua string, markers []string) bool {
	for _, marker := range markers {
		if marker != "" && strings.Contains(ua, marker) {
			return true
		}
	}
	return false
}

// isScannerInteraction reports whether an interaction with the given result
// should be attributed to a security scanner rather than a real recipient,
// using the admin-configured ScannerSettings.
//
// Decision (AND on timing, OR on the user-agent signals):
//  1. the interaction happens within the configured window of the email being
//     sent (sandboxes detonate near-instantly), AND
//  2. its user-agent looks like a scanner — either a configured scanner string,
//     or an outdated browser version (sandboxes such as Safe Links pin an old
//     Chrome build while real recipients run current browsers).
func isScannerInteraction(rs models.Result, d models.EventDetails) bool {
	s := models.GetScannerSettings()
	window := time.Duration(s.WindowSeconds) * time.Second
	if window <= 0 || rs.SendDate.IsZero() || time.Since(rs.SendDate) > window {
		return false
	}
	ua := strings.ToLower(d.Browser["user-agent"])
	if ua == "" {
		return false
	}
	return matchesKnownScanner(ua, s.UserAgentList()) || isOutdatedBrowser(ua, s)
}
