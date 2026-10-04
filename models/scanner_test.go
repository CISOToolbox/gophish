package models

import (
	check "gopkg.in/check.v1"
)

// TestCampaignIgnoreScannersRoundTrip verifies the ignore_scanners column is
// persisted and reloaded (exercises the migration + model).
func (s *ModelsSuite) TestCampaignIgnoreScannersRoundTrip(c *check.C) {
	campaign := s.createCampaignDependencies(c)
	campaign.IgnoreScanners = true
	c.Assert(PostCampaign(&campaign, campaign.UserId), check.Equals, nil)

	got, err := GetCampaign(campaign.Id, campaign.UserId)
	c.Assert(err, check.Equals, nil)
	c.Assert(got.IgnoreScanners, check.Equals, true)
}

// TestHandleScannerDoesNotAdvanceStatus verifies that recording a scanner
// interaction logs an event but leaves the result status unchanged (so it is
// not counted in the funnel).
func (s *ModelsSuite) TestHandleScannerDoesNotAdvanceStatus(c *check.C) {
	campaign := s.createCampaign(c)
	cr, err := GetCampaignResults(campaign.Id, campaign.UserId)
	c.Assert(err, check.Equals, nil)
	c.Assert(len(cr.Results) > 0, check.Equals, true)

	r := cr.Results[0]
	before := r.Status

	c.Assert(r.HandleScannerClicked(EventDetails{Browser: map[string]string{"user-agent": "HeadlessChrome"}}), check.Equals, nil)
	c.Assert(r.HandleScannerOpened(EventDetails{Browser: map[string]string{"user-agent": "HeadlessChrome"}}), check.Equals, nil)

	// Status must be unchanged (the scanner handlers never advance it).
	c.Assert(r.Status, check.Equals, before)

	reloaded, err := GetResult(r.RId)
	c.Assert(err, check.Equals, nil)
	c.Assert(reloaded.Status, check.Equals, before)

	// The scanner events are present in the timeline for transparency.
	after, err := GetCampaignResults(campaign.Id, campaign.UserId)
	c.Assert(err, check.Equals, nil)
	seenClick, seenOpen := false, false
	for _, e := range after.Events {
		if e.Message == EventScannerClicked {
			seenClick = true
		}
		if e.Message == EventScannerOpened {
			seenOpen = true
		}
	}
	c.Assert(seenClick, check.Equals, true)
	c.Assert(seenOpen, check.Equals, true)
}
