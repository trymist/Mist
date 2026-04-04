package projects

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	projectIdStr := r.URL.Query().Get("id")
	if projectIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "project id is required")
		return
	}
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid project ID", err.Error())
		return
	}

	project, err := models.GetProjectByID(projectId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no project with that ID")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	if project.OwnerID != userData.ID {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only the project owner can delete the project", "forbidden")
		return
	}

	models.LogUserAudit(userData.ID, "delete", "project", &projectId, map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
	})

	err = models.DeleteProjectByID(projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete project", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, nil, "Project deleted successfully", "")
}
