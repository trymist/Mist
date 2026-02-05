package applications

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

type CreateVolumeRequest struct {
	AppID         int64  `json:"appId"`
	Name          string `json:"name"`
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	ReadOnly      bool   `json:"readOnly"`
}

type UpdateVolumeRequest struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	ReadOnly      bool   `json:"readOnly"`
}

type DeleteVolumeRequest struct {
	ID int64 `json:"id"`
}

func GetVolumes(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "")
		return
	}

	volumes, err := models.GetVolumesByAppID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get volumes", err.Error())
		return
	}

	var volumesJSON []map[string]interface{}
	for _, vol := range volumes {
		volumesJSON = append(volumesJSON, vol.ToJson())
	}

	handlers.SendResponse(w, http.StatusOK, true, volumesJSON, "Volumes retrieved successfully", "")
}

func CreateVolume(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req CreateVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.AppID == 0 || req.Name == "" || req.HostPath == "" || req.ContainerPath == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "All fields are required", "")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "")
		return
	}

	volume, err := models.CreateVolume(req.AppID, req.Name, req.HostPath, req.ContainerPath, req.ReadOnly)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create volume", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "create", "volume", &volume.ID, map[string]interface{}{
		"app_id":         req.AppID,
		"name":           req.Name,
		"container_path": req.ContainerPath,
	})

	response := map[string]interface{}{
		"volume":         volume.ToJson(),
		"actionRequired": "restart",
		"actionMessage":  "Volume changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusCreated, true, response, "Volume created successfully", "")
}

func UpdateVolume(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req UpdateVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.ID == 0 || req.Name == "" || req.HostPath == "" || req.ContainerPath == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "All fields are required", "")
		return
	}

	volume, err := models.GetVolumeByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get volume", err.Error())
		return
	}
	if volume == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Volume not found", "")
		return
	}

	app, err := models.GetApplicationByID(volume.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this volume", "")
		return
	}

	err = models.UpdateVolume(req.ID, req.Name, req.HostPath, req.ContainerPath, req.ReadOnly)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update volume", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "update", "volume", &req.ID, map[string]interface{}{
		"name":           req.Name,
		"container_path": req.ContainerPath,
	})

	response := map[string]interface{}{
		"actionRequired": "restart",
		"actionMessage":  "Volume changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Volume updated successfully", "")
}

func DeleteVolume(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req DeleteVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Volume ID is required", "")
		return
	}

	volume, err := models.GetVolumeByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get volume", err.Error())
		return
	}
	if volume == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Volume not found", "")
		return
	}

	app, err := models.GetApplicationByID(volume.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this volume", "")
		return
	}

	err = models.DeleteVolume(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete volume", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "delete", "volume", &req.ID, map[string]interface{}{
		"app_id": volume.AppID,
		"name":   volume.Name,
	})

	response := map[string]interface{}{
		"actionRequired": "restart",
		"actionMessage":  "Volume changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Volume deleted successfully", "")
}
