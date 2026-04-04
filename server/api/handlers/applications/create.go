package applications

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		Name         string            `json:"name"`
		Description  string            `json:"description"`
		ProjectID    int64             `json:"projectId"`
		AppType      string            `json:"appType"`      // "web", "service", "database"
		TemplateName *string           `json:"templateName"` // For database type
		Port         *int              `json:"port"`         // For web type
		ShouldExpose *bool             `json:"shouldExpose"` // For web type
		ExposePort   *int              `json:"exposePort"`   // For web type
		EnvVars      map[string]string `json:"envVars"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.Name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Application name is required", "Missing fields")
		return
	}

	if req.ProjectID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project ID is required", "Missing fields")
		return
	}

	if req.AppType == "" {
		req.AppType = "web"
	}

	if req.AppType != "web" && req.AppType != "service" && req.AppType != "database" && req.AppType != "compose" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid app type", "Must be 'web', 'service', 'database', or 'compose'")
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

	app := models.App{
		Name:        req.Name,
		Description: &req.Description,
		ProjectID:   req.ProjectID,
		CreatedBy:   userInfo.ID,
		AppType:     models.AppType(req.AppType),
	}

	switch req.AppType {
	case "web":
		if req.Port != nil {
			port := int64(*req.Port)
			app.Port = &port
		} else {
			defaultPort := int64(3000)
			app.Port = &defaultPort
		}

		if req.ShouldExpose != nil {
			app.ShouldExpose = req.ShouldExpose
		} else {
			defaultShouldExpose := true
			app.ShouldExpose = &defaultShouldExpose
		}

		if req.ExposePort != nil {
			exposePort := int64(*req.ExposePort)
			app.ExposePort = &exposePort
		}

	case "service":
		// this port doesn't matter becuase yeh externally available nhi hai its something that will work without being exposed to external http
		internalPort := int64(3000)
		app.Port = &internalPort

		if req.ShouldExpose != nil {
			app.ShouldExpose = req.ShouldExpose
		} else {
			defaultShouldExpose := false
			app.ShouldExpose = &defaultShouldExpose
		}

		if req.ExposePort != nil {
			exposePort := int64(*req.ExposePort)
			app.ExposePort = &exposePort
		}

	case "database":
		if req.TemplateName == nil || *req.TemplateName == "" {
			handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Template name is required for database apps", "Missing template")
			return
		}

		app.TemplateName = req.TemplateName

		template, err := models.GetServiceTemplateByName(*req.TemplateName)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch template", err.Error())
			return
		}
		if template == nil {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "Service template not found", "")
			return
		}

		port := int64(template.DefaultPort)
		app.Port = &port

		if template.RecommendedCPU != nil {
			app.CPULimit = template.RecommendedCPU
		}
		if template.RecommendedMemory != nil {
			app.MemoryLimit = template.RecommendedMemory
		}

		if req.ShouldExpose != nil {
			app.ShouldExpose = req.ShouldExpose
		} else {
			defaultShouldExpose := false
			app.ShouldExpose = &defaultShouldExpose
		}

		if req.ExposePort != nil {
			exposePort := int64(*req.ExposePort)
			app.ExposePort = &exposePort
		}
	}

	if err := app.InsertInDB(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create application", err.Error())
		return
	}

	if req.AppType == "web" {
		project, err := models.GetProjectByID(req.ProjectID)
		if err == nil {
			autoDomain, err := models.GenerateAutoDomain(project.Name, app.Name)
			if err == nil && autoDomain != "" {
				_, err = models.CreateDomain(app.ID, autoDomain)
				if err != nil {
				}
			}
		}
	}

	if len(req.EnvVars) > 0 {
		for key, value := range req.EnvVars {
			_, err := models.CreateEnvVariable(app.ID, key, value)
			if err != nil {
				continue
			}
		}
	}
	models.LogUserAudit(userInfo.ID, "create", "application", &app.ID, map[string]interface{}{
		"name":        app.Name,
		"description": app.Description,
		"project_id":  app.ProjectID,
	})

	handlers.SendResponse(w, http.StatusOK, true, app.ToJson(), "Application created successfully", "")

}
