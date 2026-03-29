package applications

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetApplicationById(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", fmt.Sprintf("Error fetching application: %v", err))
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "No application with the given ID exists")
		return
	}
	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "Forbidden")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, app.ToJson(),
		"Application retrieved successfully", "")
}
