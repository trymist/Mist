package templates

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
)

func ListServiceTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := models.GetAllServiceTemplates()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch service templates", err.Error())
		return
	}

	var templatesJson []map[string]interface{}
	for _, template := range templates {
		templatesJson = append(templatesJson, template.ToJson())
	}

	handlers.SendResponse(w, http.StatusOK, true, templatesJson, "Service templates fetched successfully", "")
}

func GetServiceTemplateByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Template name is required", "Missing name parameter")
		return
	}

	template, err := models.GetServiceTemplateByName(name)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch service template", err.Error())
		return
	}

	if template == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Service template not found", "")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, template.ToJson(), "Service template fetched successfully", "")
}
