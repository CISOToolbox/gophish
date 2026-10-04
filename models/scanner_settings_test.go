package models

import (
	check "gopkg.in/check.v1"
)

func (s *ModelsSuite) TestScannerSettingsSeededAndRoundTrip(c *check.C) {
	// Setup() seeds the row; it should be present with defaults.
	got := GetScannerSettings()
	c.Assert(got.WindowSeconds, check.Equals, defaultScannerWindowSeconds)
	c.Assert(got.MinChrome, check.Equals, defaultScannerMinChrome)
	c.Assert(len(got.UserAgentList()) > 0, check.Equals, true)

	// Update and read back (also exercises the cache refresh).
	got.WindowSeconds = 90
	got.MinChrome = 130
	got.UserAgents = "headlesschrome\nProofPoint\n\n  curl/  "
	c.Assert(PutScannerSettings(&got), check.Equals, nil)

	reloaded := GetScannerSettings()
	c.Assert(reloaded.WindowSeconds, check.Equals, 90)
	c.Assert(reloaded.MinChrome, check.Equals, 130)
	// UserAgentList trims/lowercases/drops blanks.
	c.Assert(reloaded.UserAgentList(), check.DeepEquals, []string{"headlesschrome", "proofpoint", "curl/"})
}
