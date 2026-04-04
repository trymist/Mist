package applications

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/git"
	"github.com/corecollectives/mist/models"
)

func GetLatestCommit(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID     int64 `json:"appId"`
		ProjectID int64 `json:"projectId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}
	if req.ProjectID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project ID is required", "Missing fields")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, req.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify project access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this project", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get app details", err.Error())
		return
	}

	if app.AppType == models.AppTypeDatabase {
		handlers.SendResponse(w, http.StatusOK, true, nil, "Database apps do not have commits", "")
		return
	}

	commit, err := git.GetLatestCommit(req.AppID, userInfo.ID)
	if err != nil {
		// Check if it's a "no repository configured" error
		errMsg := err.Error()
		if errMsg == "no git repository configured for this application" {
			handlers.SendResponse(w, http.StatusBadRequest, false, nil, "No git repository configured", "Please configure a git repository for this application")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get latest commit", errMsg)
		return
	}
	handlers.SendResponse(w, http.StatusOK, true, commit, "Latest commit retrieved successfully", "")

}
