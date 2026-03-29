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

func GetProjectFromId(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	userID := userData.ID

	projectIDStr := r.URL.Query().Get("id")
	if projectIDStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "no id provided")
		return
	}

	projectId, err := strconv.ParseInt(projectIDStr, 10, 64)
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

	hasAccess, err := models.HasUserAccessToProject(userID, projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}
	if !hasAccess {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied to this project", "forbidden")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, project.ToJSON(), "Project retrieved successfully", "")

}
