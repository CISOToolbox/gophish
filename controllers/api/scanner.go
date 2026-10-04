package api

import (
	"encoding/json"
	"net/http"

	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
)

// Scanner handles GET/PUT for the global scanner-detection settings
// (/api/scanner/). It is admin-only (wired with RequirePermission).
func (as *Server) Scanner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		JSONResponse(w, models.GetScannerSettings(), http.StatusOK)
	case http.MethodPut:
		s := models.ScannerSettings{}
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}
		if err := models.PutScannerSettings(&s); err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, s, http.StatusOK)
	}
}
