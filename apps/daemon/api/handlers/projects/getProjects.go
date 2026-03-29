package projects

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetProjects(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	userId := userData.ID

	projects, err := models.GetProjectsUserIsPartOf(userId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve projects", err.Error())
		return
	}
	if projects == nil {
		projects = []models.Project{}
	}
	responseProjects := make([]interface{}, len(projects))
	for i, project := range projects {
		responseProjects[i] = project.ToJSON()
	}
	handlers.SendResponse(w, http.StatusOK, true, responseProjects, "Projects retrieved successfully", "")
}
