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

func UpdateProject(w http.ResponseWriter, r *http.Request) {
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
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}
	if input.Name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project name is required", "name field is empty")
		return
	}

	existingProject, err := models.GetProjectByID(projectId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no such project")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	if existingProject.OwnerID != userData.ID {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}

	// tags := make([]string, len(input.Tags))
	// for i, tag := range input.Tags {
	// 	tags[i] = tag
	// }
	var desc *string
	if input.Description != "" {
		desc = &input.Description

	} else {
		desc = nil
	}

	project := &models.Project{
		ID:          projectId,
		Name:        input.Name,
		Description: desc,
		Tags:        input.Tags,
	}

	err = models.UpdateProject(project)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update project", err.Error())
		return
	}

	updatedProject, err := models.GetProjectByID(projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch updated project", err.Error())
		return
	}

	models.LogUserAudit(userData.ID, "update", "project", &project.ID, map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
	})

	handlers.SendResponse(w, http.StatusOK, true, updatedProject.ToJSON(), "Project updated successfully", "")
}
