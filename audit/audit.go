// Package audit provides a lightweight, append-only audit log for
// security-relevant events (authentication and sensitive admin actions). It
// writes human-readable lines to a configurable file, independent of the
// application log, so the audit trail can be shipped and retained separately.
package audit

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Event names used across the codebase. Keeping them as constants avoids typos
// and makes the set of audited actions easy to see at a glance.
const (
	EventLogin            = "LOGIN"
	EventLogout           = "LOGOUT"
	EventCampaignCreate   = "CAMPAIGN_CREATE"
	EventCampaignLaunch   = "CAMPAIGN_LAUNCH"
	EventCampaignDelete   = "CAMPAIGN_DELETE"
	EventCampaignComplete = "CAMPAIGN_COMPLETE"
	EventUserCreate       = "USER_CREATE"
	EventUserModify       = "USER_MODIFY"
	EventUserDelete       = "USER_DELETE"
	EventAPIKeyReset      = "APIKEY_RESET"
)

var (
	mu      sync.Mutex
	logPath = "gophish_audit.log"
)

// SetLogPath sets the audit log file path from configuration. An empty path is
// ignored so the default is kept.
func SetLogPath(path string) {
	if path != "" {
		logPath = path
	}
}

// clientIP returns a best-effort client IP for the request. It honors the
// left-most X-Forwarded-For entry when present (meaningful only behind a
// trusted reverse proxy that sets it), otherwise the direct remote address.
func clientIP(r *http.Request) string {
	if r == nil {
		return "-"
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if first := strings.TrimSpace(strings.Split(xff, ",")[0]); first != "" {
			return first
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// LogEvent appends a single audit entry. r may be nil for events that are not
// tied to an HTTP request. Failures to write are intentionally swallowed: the
// audit log must never break the request path.
func LogEvent(r *http.Request, username, event, status string) {
	mu.Lock()
	defer mu.Unlock()

	if dir := filepath.Dir(logPath); dir != "." {
		os.MkdirAll(dir, 0750)
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return
	}
	defer f.Close()

	if username == "" {
		username = "-"
	}
	line := fmt.Sprintf("%s | %-16s | %-24s | IP: %-20s | %s\n",
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		event,
		username,
		clientIP(r),
		status,
	)
	f.WriteString(line)
}
