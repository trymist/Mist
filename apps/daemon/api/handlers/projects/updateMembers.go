package projects

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

func UpdateMembers(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	projectIDStr := r.URL.Query().Get("id")
	if projectIDStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "project id is required")
		return
	}
	projectId, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid project ID", err.Error())
		return
	}

	var input struct {
		UserIDs []int64 `json:"userIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	project, err := models.GetProjectByID(projectId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no such project")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	if project.OwnerID != userData.ID {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Only the project owner can update members")
		return
	}

	err = models.UpdateProjectMembers(projectId, input.UserIDs)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update project members", err.Error())
		return
	}

	updatedProject, err := models.GetProjectByID(projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch updated project", err.Error())
		return
	}

	models.LogUserAudit(userData.ID, "update", "project", &projectId, map[string]interface{}{
		"action":  "members_update",
		"userIds": input.UserIDs,
	})

	handlers.SendResponse(w, http.StatusOK, true, updatedProject.ToJSON(), "Project members updated successfully", "")
}
