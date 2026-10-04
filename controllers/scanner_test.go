package controllers

import (
	"testing"
	"time"

	"github.com/gophish/gophish/models"
)

func details(ua string) models.EventDetails {
	return models.EventDetails{Browser: map[string]string{"user-agent": ua}}
}

func resultSentAgo(d time.Duration) models.Result {
	return models.Result{SendDate: time.Now().Add(-d)}
}

func TestIsScannerInteraction(t *testing.T) {
	chrome := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
	headless := "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/120.0 Safari/537.36"

	cases := []struct {
		name string
		rs   models.Result
		d    models.EventDetails
		want bool
	}{
		{"scanner UA within window", resultSentAgo(5 * time.Second), details(headless), true},
		{"vendor UA within window", resultSentAgo(10 * time.Second), details("ProofPoint-Scanner/1.0"), true},
		{"scanner UA but window elapsed", resultSentAgo(5 * time.Minute), details(headless), false},
		{"real browser within window", resultSentAgo(5 * time.Second), details(chrome), false},
		{"empty user-agent within window", resultSentAgo(5 * time.Second), details(""), false},
		{"zero send date", models.Result{}, details(headless), false},
	}
	for _, tc := range cases {
		if got := isScannerInteraction(tc.rs, tc.d); got != tc.want {
			t.Errorf("%s: isScannerInteraction = %v, want %v", tc.name, got, tc.want)
		}
	}
}
