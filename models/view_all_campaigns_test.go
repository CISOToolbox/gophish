package models

import (
	check "gopkg.in/check.v1"
)

// TestViewAllCampaignsPermissions verifies the role/permission wiring seeded by
// the migration: auditor and admin can view all campaigns, a plain user cannot,
// and auditor remains read-only (no modify permission).
func (s *ModelsSuite) TestViewAllCampaignsPermissions(c *check.C) {
	auditor, err := GetRoleBySlug(RoleAuditor)
	c.Assert(err, check.Equals, nil)
	admin, err := GetRoleBySlug(RoleAdmin)
	c.Assert(err, check.Equals, nil)
	user, err := GetRoleBySlug(RoleUser)
	c.Assert(err, check.Equals, nil)

	auditorUser := User{RoleID: auditor.ID}
	adminUser := User{RoleID: admin.ID}
	plainUser := User{RoleID: user.ID}

	has := func(u User, perm string) bool {
		ok, err := u.HasPermission(perm)
		c.Assert(err, check.Equals, nil)
		return ok
	}

	// view_all_campaigns: auditor + admin yes, plain user no.
	c.Assert(has(auditorUser, PermissionViewAllCampaigns), check.Equals, true)
	c.Assert(has(adminUser, PermissionViewAllCampaigns), check.Equals, true)
	c.Assert(has(plainUser, PermissionViewAllCampaigns), check.Equals, false)

	// auditor is read-only: it can view objects but not modify them.
	c.Assert(has(auditorUser, PermissionViewObjects), check.Equals, true)
	c.Assert(has(auditorUser, PermissionModifyObjects), check.Equals, false)
	c.Assert(has(auditorUser, PermissionModifySystem), check.Equals, false)
}

// TestGetCampaignsAllCrossUser verifies that the all-view queries return a
// campaign owned by another user, while the owner-scoped queries do not.
func (s *ModelsSuite) TestGetCampaignsAllCrossUser(c *check.C) {
	campaign := s.createCampaign(c) // owned by user id 1
	const otherUID = int64(2)

	// The owner-scoped listing for a different user sees nothing...
	scoped, err := GetCampaigns(otherUID)
	c.Assert(err, check.Equals, nil)
	c.Assert(len(scoped), check.Equals, 0)

	// ...but the all-view listing sees the campaign.
	all, err := GetCampaignsAll()
	c.Assert(err, check.Equals, nil)
	found := false
	for _, cg := range all {
		if cg.Id == campaign.Id {
			found = true
		}
	}
	c.Assert(found, check.Equals, true)

	// GetCampaignByID resolves it without an owner filter.
	got, err := GetCampaignByID(campaign.Id)
	c.Assert(err, check.Equals, nil)
	c.Assert(got.Id, check.Equals, campaign.Id)

	// The all-view summary and results also resolve it.
	summary, err := GetCampaignSummaryAll(campaign.Id)
	c.Assert(err, check.Equals, nil)
	c.Assert(summary.Id, check.Equals, campaign.Id)

	results, err := GetCampaignResultsAll(campaign.Id)
	c.Assert(err, check.Equals, nil)
	c.Assert(results.Id, check.Equals, campaign.Id)
}
