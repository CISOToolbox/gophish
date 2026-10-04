package controllers

import (
	"testing"

	"github.com/gophish/gophish/models"
)

// defaultSettings mirrors the seeded defaults for pure-function tests
// (Chrome/Edge floor 125, Firefox floor 120, Safari check disabled).
func defaultSettings() models.ScannerSettings {
	return models.ScannerSettings{
		WindowSeconds: 120,
		MinChrome:     125,
		MinFirefox:    120,
		MinSafari:     0,
		UserAgents:    "headlesschrome\nproofpoint\ncurl/",
	}
}

func TestIsOutdatedBrowser(t *testing.T) {
	s := defaultSettings()
	cases := []struct {
		ua   string
		want bool
	}{
		{"mozilla/5.0 (windows nt 10.0; win64; x64) applewebkit/537.36 (khtml, like gecko) chrome/109.0.0.0 safari/537.36", true},        // safe links pin
		{"mozilla/5.0 (macintosh; intel mac os x 10_15_7) applewebkit/537.36 (khtml, like gecko) chrome/143.0.0.0 safari/537.36", false}, // real, near-current
		{"mozilla/5.0 (windows nt 10.0) applewebkit/537.36 (khtml, like gecko) chrome/124.0.0.0 safari/537.36", true},                    // 124 < 125
		{"mozilla/5.0 (windows nt 10.0) applewebkit/537.36 (khtml, like gecko) chrome/130.0.0.0 safari/537.36", false},                   // 130 >= 125
		{"mozilla/5.0 applewebkit/537.36 (khtml, like gecko) chrome/143.0.0.0 safari/537.36 edg/109.0.0.0", true},                        // edge checked first
		{"mozilla/5.0 (x11; linux x86_64; rv:102.0) gecko/20100101 firefox/102.0", true},
		{"mozilla/5.0 (x11; linux x86_64; rv:128.0) gecko/20100101 firefox/128.0", false},
		{"mozilla/5.0 (macintosh) applewebkit/605 (khtml, like gecko) version/14.1 safari/605", false}, // safari check disabled (min 0)
		{"curl/8.0.1", false}, // no browser version
	}
	for _, tc := range cases {
		if got := isOutdatedBrowser(tc.ua, s); got != tc.want {
			t.Errorf("isOutdatedBrowser(%q) = %v, want %v", tc.ua, got, tc.want)
		}
	}
}

func TestMatchesKnownScanner(t *testing.T) {
	markers := defaultSettings().UserAgentList()
	if !matchesKnownScanner("mozilla/5.0 headlesschrome/120 safari/537", markers) {
		t.Error("expected headlesschrome to match")
	}
	if !matchesKnownScanner("proofpoint-scanner/1.0", markers) {
		t.Error("expected proofpoint to match")
	}
	if matchesKnownScanner("mozilla/5.0 (macintosh) chrome/143.0.0.0 safari/537.36", markers) {
		t.Error("did not expect a clean chrome UA to match")
	}
}
