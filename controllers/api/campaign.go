package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gophish/gophish/audit"
	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// auditActor returns the username of the authenticated caller for audit
// entries, or "-" when it cannot be resolved.
func auditActor(r *http.Request) string {
	if u, ok := ctx.Get(r, "user").(models.User); ok {
		return u.Username
	}
	return "-"
}

// canViewAllCampaigns reports whether the requesting user has been granted
// PermissionViewAllCampaigns, i.e. read-only visibility over every user's
// campaigns and results.
func canViewAllCampaigns(r *http.Request) bool {
	u, ok := ctx.Get(r, "user").(models.User)
	if !ok {
		return false
	}
	ok, err := u.HasPermission(models.PermissionViewAllCampaigns)
	if err != nil {
		log.Error(err)
		return false
	}
	return ok
}

// Campaigns returns a list of campaigns if requested via GET.
// If requested via POST, APICampaigns creates a new campaign and returns a reference to it.
func (as *Server) Campaigns(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		var cs []models.Campaign
		var err error
		if canViewAllCampaigns(r) {
			cs, err = models.GetCampaignsAll()
		} else {
			cs, err = models.GetCampaigns(ctx.Get(r, "user_id").(int64))
		}
		if err != nil {
			log.Error(err)
		}
		JSONResponse(w, cs, http.StatusOK)
	//POST: Create a new campaign and return it as JSON
	case r.Method == "POST":
		c := models.Campaign{}
		// Put the request into a campaign
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		err = models.PostCampaign(&c, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		audit.LogEvent(r, auditActor(r), audit.EventCampaignCreate, fmt.Sprintf("campaign '%s' (id %d)", c.Name, c.Id))
		// If the campaign is scheduled to launch immediately, send it to the worker.
		// Otherwise, the worker will pick it up at the scheduled time
		if c.Status == models.CampaignInProgress {
			audit.LogEvent(r, auditActor(r), audit.EventCampaignLaunch, fmt.Sprintf("campaign '%s' (id %d) launched immediately", c.Name, c.Id))
			go as.worker.LaunchCampaign(c)
		}
		JSONResponse(w, c, http.StatusCreated)
	}
}

// CampaignsSummary returns the summary for the current user's campaigns
func (as *Server) CampaignsSummary(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		var cs models.CampaignSummaries
		var err error
		if canViewAllCampaigns(r) {
			cs, err = models.GetCampaignSummariesAll()
		} else {
			cs, err = models.GetCampaignSummaries(ctx.Get(r, "user_id").(int64))
		}
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, cs, http.StatusOK)
	}
}

// Campaign returns details about the requested campaign. If the campaign is not
// valid, APICampaign returns null.
func (as *Server) Campaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	uid := ctx.Get(r, "user_id").(int64)
	var c models.Campaign
	var err error
	if canViewAllCampaigns(r) {
		c, err = models.GetCampaignByID(id)
	} else {
		c, err = models.GetCampaign(id, uid)
	}
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, c, http.StatusOK)
	case r.Method == "DELETE":
		// Deletion stays owner-scoped even for users who can view all
		// campaigns: PermissionViewAllCampaigns is read-only.
		if c.UserId != uid {
			JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
			return
		}
		err = models.DeleteCampaign(id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting campaign"}, http.StatusInternalServerError)
			return
		}
		audit.LogEvent(r, auditActor(r), audit.EventCampaignDelete, fmt.Sprintf("campaign '%s' (id %d)", c.Name, id))
		JSONResponse(w, models.Response{Success: true, Message: "Campaign deleted successfully!"}, http.StatusOK)
	}
}

// CampaignResults returns just the results for a given campaign to
// significantly reduce the information returned.
func (as *Server) CampaignResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	var cr models.CampaignResults
	var err error
	if canViewAllCampaigns(r) {
		cr, err = models.GetCampaignResultsAll(id)
	} else {
		cr, err = models.GetCampaignResults(id, ctx.Get(r, "user_id").(int64))
	}
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		JSONResponse(w, cr, http.StatusOK)
		return
	}
}

// CampaignSummary returns the summary for a given campaign.
func (as *Server) CampaignSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	switch {
	case r.Method == "GET":
		var cs models.CampaignSummary
		var err error
		if canViewAllCampaigns(r) {
			cs, err = models.GetCampaignSummaryAll(id)
		} else {
			cs, err = models.GetCampaignSummary(id, ctx.Get(r, "user_id").(int64))
		}
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
			} else {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			}
			log.Error(err)
			return
		}
		JSONResponse(w, cs, http.StatusOK)
	}
}

// CampaignResultReport lets an admin manually declare that a recipient
// reported the simulated phishing email, specifying the channel it was
// reported through and the time it happened.
func (as *Server) CampaignResultReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	rid := vars["rid"]
	uid := ctx.Get(r, "user_id").(int64)
	// Declaring a report mutates the campaign, so it stays owner-scoped.
	if _, err := models.GetCampaign(id, uid); err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Campaign not found"}, http.StatusNotFound)
		return
	}
	result, err := models.GetResult(rid)
	if err != nil || result.CampaignId != id || result.UserId != uid {
		JSONResponse(w, models.Response{Success: false, Message: "Result not found"}, http.StatusNotFound)
		return
	}
	req := struct {
		Channel    string    `json:"channel"`
		ReportDate time.Time `json:"report_date"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}
	if err := result.HandleEmailReportManual(req.ReportDate, req.Channel); err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, models.Response{Success: true, Message: "Report recorded"}, http.StatusOK)
}

// CampaignComplete effectively "ends" a campaign.
// Future phishing emails clicked will return a simple "404" page.
func (as *Server) CampaignComplete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	switch {
	case r.Method == "GET":
		err := models.CompleteCampaign(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error completing campaign"}, http.StatusInternalServerError)
			return
		}
		audit.LogEvent(r, auditActor(r), audit.EventCampaignComplete, fmt.Sprintf("campaign id %d", id))
		JSONResponse(w, models.Response{Success: true, Message: "Campaign completed successfully!"}, http.StatusOK)
	}
}
