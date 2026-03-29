package projects

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userData.Role != "owner" && userData.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
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

	project := models.Project{
		Name:    input.Name,
		OwnerID: userData.ID,
	}
	if input.Description != "" {
		desc := input.Description
		project.Description = &desc
	} else {
		project.Description = nil
	}
	if input.Tags != nil {
		project.Tags = input.Tags
	}
	err := project.InsertInDB()

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create project", err.Error())
		return
	}

	models.LogUserAudit(userData.ID, "create", "project", &project.ID, map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
	})

	handlers.SendResponse(w, http.StatusCreated, true, project.ToJSON(), "Project created successfully", "")
}
