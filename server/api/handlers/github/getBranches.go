package github

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/github"
)

func GetBranches(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		Repo string `json:"repo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.Repo == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Repository is required", "Missing fields")
		return
	}

	token, _, err := github.GetGitHubAccessToken(int(userInfo.ID))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get GitHub access token", err.Error())
		return
	}
	branches, err := github.GetGitHubBranches(token, req.Repo)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get branches", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, branches, "Branches retrieved successfully", "")

}
