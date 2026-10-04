package controllers

import (
	"strings"
	"time"

	"github.com/gophish/gophish/models"
)

// scannerWindow is how long after an email was sent an interaction may still be
// attributed to an automated security scanner / sandbox. Mail sandboxes (e.g.
// Microsoft Defender Safe Links) detonate links within seconds of delivery, so
// an interaction arriving within this window AND carrying a known scanner
// user-agent is treated as a scan rather than a real recipient action.
const scannerWindow = 60 * time.Second

// scannerUserAgents holds lowercased substrings that identify automated mail
// scanners, sandboxes and link/image proxies. Matching is a case-insensitive
// substring test, so broad markers are fine: detection also requires the
// interaction to fall inside scannerWindow (logical AND), which keeps false
// positives low. "headlesschrome" is included because detonation chambers
// (including Safe Links) commonly drive a headless Chromium.
var scannerUserAgents = []string{
	// Headless / automation engines used by detonation sandboxes
	"headlesschrome", "phantomjs", "slimerjs", "electron",
	"python-requests", "go-http-client", "libwww-perl", "curl/", "wget/",
	"java/", "okhttp", "apache-httpclient", "httpclient",
	// Microsoft (Defender for Office 365 / Safe Links, Bing, Skype preview)
	"bingpreview", "skypeuripreview", "office365", "microsoft",
	// Mail security vendors
	"proofpoint", "mimecast", "barracuda", "forcepoint", "fireeye",
	"symantec", "trendmicro", "mcafee", "kaspersky", "avast", "avira",
	"sophos", "paloalto", "ironport", "cisco", "fortinet", "zscaler",
	"messagelabs", "cloudmark", "spamtitan", "vadesecure", "gatefy",
	// Link / image proxies and generic crawlers
	"googleimageproxy", "google-safety", "googlebot", "yahoo! slurp",
	"bitdefender", "urlscan", "virustotal", "safebrowsing",
}

// isScannerInteraction reports whether an interaction with the given result
// should be attributed to a security scanner rather than a real recipient.
//
// Both conditions must hold (AND):
//  1. the interaction happens within scannerWindow of the email being sent, and
//  2. the request's user-agent matches a known scanner.
//
// This is intentionally conservative: it only reclassifies interactions that
// look like a scanner on both timing and user-agent.
func isScannerInteraction(rs models.Result, d models.EventDetails) bool {
	// Timing: rs.SendDate holds the actual send time once the email was sent
	// (set by HandleEmailSent). If it is zero or the window has elapsed, this
	// is not a near-instant detonation.
	if rs.SendDate.IsZero() || time.Since(rs.SendDate) > scannerWindow {
		return false
	}
	ua := strings.ToLower(d.Browser["user-agent"])
	if ua == "" {
		return false
	}
	for _, marker := range scannerUserAgents {
		if strings.Contains(ua, marker) {
			return true
		}
	}
	return false
}
