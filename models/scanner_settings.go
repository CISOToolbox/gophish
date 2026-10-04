package models

import (
	"strings"
	"sync"

	log "github.com/gophish/gophish/logger"
	"gorm.io/gorm"
)

// scannerSettingsID is the fixed primary key of the single settings row.
const scannerSettingsID = 1

// DefaultScannerUserAgents is the built-in seed list of lowercased user-agent
// substrings that identify mail scanners, sandboxes and link/image proxies.
// It is written to the database on first run and is fully editable from the
// admin UI afterwards.
var DefaultScannerUserAgents = []string{
	"headlesschrome", "phantomjs", "slimerjs", "electron",
	"python-requests", "go-http-client", "libwww-perl", "curl/", "wget/",
	"java/", "okhttp", "apache-httpclient", "httpclient",
	"bingpreview", "skypeuripreview", "office365", "microsoft",
	"proofpoint", "mimecast", "barracuda", "forcepoint", "fireeye",
	"symantec", "trendmicro", "mcafee", "kaspersky", "avast", "avira",
	"sophos", "paloalto", "ironport", "cisco", "fortinet", "zscaler",
	"messagelabs", "cloudmark", "spamtitan", "vadesecure", "gatefy",
	"googleimageproxy", "google-safety", "googlebot", "yahoo! slurp",
	"bitdefender", "urlscan", "virustotal", "safebrowsing",
}

// Default scanner-detection tuning used to seed the settings row. A minimum
// version of 0 disables the outdated-browser check for that engine.
const (
	defaultScannerWindowSeconds = 120
	defaultScannerMinChrome     = 125
	defaultScannerMinFirefox    = 120
	defaultScannerMinSafari     = 0
)

// ScannerSettings holds the (single-row) global configuration for detecting
// security-scanner / sandbox interactions. It is editable from the admin UI.
type ScannerSettings struct {
	Id            int64 `json:"-" gorm:"column:id;primaryKey"`
	WindowSeconds int   `json:"window_seconds" gorm:"column:window_seconds"`
	MinChrome     int   `json:"min_chrome" gorm:"column:min_chrome"`
	MinFirefox    int   `json:"min_firefox" gorm:"column:min_firefox"`
	MinSafari     int   `json:"min_safari" gorm:"column:min_safari"`
	// UserAgents is a newline-separated list of lowercased substrings.
	UserAgents string `json:"user_agents" gorm:"column:user_agents"`
}

// TableName overrides the GORM table name.
func (ScannerSettings) TableName() string { return "scanner_settings" }

// UserAgentList returns the parsed, lowercased, non-empty scanner user-agent
// substrings.
func (s ScannerSettings) UserAgentList() []string {
	out := []string{}
	for _, line := range strings.Split(s.UserAgents, "\n") {
		if v := strings.ToLower(strings.TrimSpace(line)); v != "" {
			out = append(out, v)
		}
	}
	return out
}

var (
	scannerSettingsMu    sync.RWMutex
	scannerSettingsCache *ScannerSettings
)

// ensureScannerSettings creates the settings row with defaults if it does not
// yet exist. Called once at startup from Setup.
func ensureScannerSettings() error {
	s := ScannerSettings{}
	err := db.First(&s, scannerSettingsID).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	seed := ScannerSettings{
		Id:            scannerSettingsID,
		WindowSeconds: defaultScannerWindowSeconds,
		MinChrome:     defaultScannerMinChrome,
		MinFirefox:    defaultScannerMinFirefox,
		MinSafari:     defaultScannerMinSafari,
		UserAgents:    strings.Join(DefaultScannerUserAgents, "\n"),
	}
	return db.Create(&seed).Error
}

// GetScannerSettings returns the global scanner-detection settings, cached in
// memory. It falls back to built-in defaults if the row cannot be read.
func GetScannerSettings() ScannerSettings {
	scannerSettingsMu.RLock()
	cached := scannerSettingsCache
	scannerSettingsMu.RUnlock()
	if cached != nil {
		return *cached
	}
	s := ScannerSettings{}
	if err := db.First(&s, scannerSettingsID).Error; err != nil {
		log.Error(err)
		return ScannerSettings{
			WindowSeconds: defaultScannerWindowSeconds,
			MinChrome:     defaultScannerMinChrome,
			MinFirefox:    defaultScannerMinFirefox,
			MinSafari:     defaultScannerMinSafari,
			UserAgents:    strings.Join(DefaultScannerUserAgents, "\n"),
		}
	}
	scannerSettingsMu.Lock()
	scannerSettingsCache = &s
	scannerSettingsMu.Unlock()
	return s
}

// PutScannerSettings persists the settings (always the single row) and updates
// the in-memory cache.
func PutScannerSettings(s *ScannerSettings) error {
	s.Id = scannerSettingsID
	if err := db.Save(s).Error; err != nil {
		log.Error(err)
		return err
	}
	scannerSettingsMu.Lock()
	cp := *s
	scannerSettingsCache = &cp
	scannerSettingsMu.Unlock()
	return nil
}
