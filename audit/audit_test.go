package audit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// withTempLog points the audit log at a throwaway file for the duration of a
// test and restores the previous path afterwards.
func withTempLog(t *testing.T) string {
	t.Helper()
	prev := logPath
	p := filepath.Join(t.TempDir(), "audit.log")
	logPath = p
	t.Cleanup(func() { logPath = prev })
	return p
}

func TestLogEventWritesEntry(t *testing.T) {
	p := withTempLog(t)
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "203.0.113.7:54321"

	LogEvent(req, "alice", EventLogin, "SUCCESS")

	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("reading audit log: %v", err)
	}
	line := string(data)
	for _, want := range []string{EventLogin, "alice", "203.0.113.7", "SUCCESS"} {
		if !strings.Contains(line, want) {
			t.Fatalf("audit line missing %q: %s", want, line)
		}
	}
}

func TestLogEventHonorsXForwardedFor(t *testing.T) {
	p := withTempLog(t)
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "10.0.0.1:1234" // the proxy
	req.Header.Set("X-Forwarded-For", "198.51.100.9, 10.0.0.1")

	LogEvent(req, "bob", EventLogin, "FAILED - Invalid Username/Password")

	data, _ := os.ReadFile(p)
	line := string(data)
	if !strings.Contains(line, "198.51.100.9") {
		t.Fatalf("expected client IP from X-Forwarded-For, got: %s", line)
	}
}

func TestLogEventNilRequestAndEmptyUser(t *testing.T) {
	p := withTempLog(t)
	LogEvent(nil, "", EventCampaignComplete, "campaign id 5")

	data, _ := os.ReadFile(p)
	line := string(data)
	// Nil request → "-" IP; empty username → "-".
	if !strings.Contains(line, "IP: -") || !strings.Contains(line, "| -") {
		t.Fatalf("expected placeholders for nil request/empty user: %s", line)
	}
	if !strings.Contains(line, EventCampaignComplete) {
		t.Fatalf("missing event: %s", line)
	}
}

func TestLogEventAppends(t *testing.T) {
	p := withTempLog(t)
	LogEvent(nil, "a", EventUserCreate, "one")
	LogEvent(nil, "b", EventUserDelete, "two")
	data, _ := os.ReadFile(p)
	if n := strings.Count(string(data), "\n"); n != 2 {
		t.Fatalf("expected 2 appended lines, got %d: %s", n, string(data))
	}
}
