package models

import (
	check "gopkg.in/check.v1"
)

func (s *ModelsSuite) TestPostEducationalPage(c *check.C) {
	p := EducationalPage{
		Name:   "Test Educational Page",
		HTML:   "<html><body>You were part of a phishing simulation.</body></html>",
		UserId: 1,
	}
	err := PostEducationalPage(&p)
	c.Assert(err, check.Equals, nil)

	got, err := GetEducationalPage(p.Id, 1)
	c.Assert(err, check.Equals, nil)
	c.Assert(got.Name, check.Equals, p.Name)
	c.Assert(got.HTML, check.Equals, p.HTML)
}

func (s *ModelsSuite) TestEducationalPageValidation(c *check.C) {
	p := EducationalPage{HTML: "<html></html>", UserId: 1}
	c.Assert(PostEducationalPage(&p), check.Equals, ErrEducationalPageNameNotSpecified)
}

func (s *ModelsSuite) TestPutEducationalPage(c *check.C) {
	p := EducationalPage{Name: "Edu", HTML: "<html>a</html>", UserId: 1}
	c.Assert(PostEducationalPage(&p), check.Equals, nil)
	p.HTML = "<html>b</html>"
	c.Assert(PutEducationalPage(&p), check.Equals, nil)
	got, err := GetEducationalPage(p.Id, 1)
	c.Assert(err, check.Equals, nil)
	c.Assert(got.HTML, check.Equals, "<html>b</html>")
}

func (s *ModelsSuite) TestDeleteEducationalPage(c *check.C) {
	p := EducationalPage{Name: "ToDelete", HTML: "<html></html>", UserId: 1}
	c.Assert(PostEducationalPage(&p), check.Equals, nil)
	c.Assert(DeleteEducationalPage(p.Id, 1), check.Equals, nil)
	_, err := GetEducationalPage(p.Id, 1)
	c.Assert(err, check.NotNil)
}

// TestCampaignWithEducationalPage verifies that a campaign can reference an
// educational page by name, that the association is persisted via
// EducationalPageId, and that it is reloaded by GetCampaign.
func (s *ModelsSuite) TestCampaignWithEducationalPage(c *check.C) {
	edu := EducationalPage{
		Name:   "Awareness",
		HTML:   "<html><body>Teaching moment</body></html>",
		UserId: 1,
	}
	c.Assert(PostEducationalPage(&edu), check.Equals, nil)

	campaign := s.createCampaignDependencies(c)
	campaign.EducationalPage = EducationalPage{Name: "Awareness"}
	c.Assert(PostCampaign(&campaign, campaign.UserId), check.Equals, nil)
	c.Assert(campaign.EducationalPageId, check.Equals, edu.Id)

	reloaded, err := GetCampaign(campaign.Id, campaign.UserId)
	c.Assert(err, check.Equals, nil)
	c.Assert(reloaded.EducationalPageId, check.Equals, edu.Id)
	c.Assert(reloaded.EducationalPage.HTML, check.Equals, edu.HTML)
}

// TestCampaignWithoutEducationalPage ensures the educational page stays
// optional: a campaign created without one has a zero EducationalPageId.
func (s *ModelsSuite) TestCampaignWithoutEducationalPage(c *check.C) {
	campaign := s.createCampaignDependencies(c)
	c.Assert(PostCampaign(&campaign, campaign.UserId), check.Equals, nil)

	reloaded, err := GetCampaign(campaign.Id, campaign.UserId)
	c.Assert(err, check.Equals, nil)
	c.Assert(reloaded.EducationalPageId, check.Equals, int64(0))
}

// TestCampaignMissingEducationalPage ensures a named-but-nonexistent
// educational page is rejected.
func (s *ModelsSuite) TestCampaignMissingEducationalPage(c *check.C) {
	campaign := s.createCampaignDependencies(c)
	campaign.EducationalPage = EducationalPage{Name: "Does Not Exist"}
	c.Assert(PostCampaign(&campaign, campaign.UserId), check.Equals, ErrEducationalPageNotFound)
}
