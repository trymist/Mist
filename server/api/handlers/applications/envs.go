package applications

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func CreateEnvVariable(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64  `json:"appId"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 || req.Key == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID and key are required", "Missing fields")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	env, err := models.CreateEnvVariable(req.AppID, strings.TrimSpace(req.Key), req.Value)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create environment variable", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "create", "env_variable", &env.ID, map[string]interface{}{
		"app_id": req.AppID,
		"key":    req.Key,
	})

	response := map[string]interface{}{
		"envVariable":    env,
		"actionRequired": "redeploy",
		"actionMessage":  "Environment variable changes require a full redeployment to take effect. Would you like to redeploy now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Environment variable created successfully", "")
}

func GetEnvVariables(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to access this application", "Forbidden")
		return
	}

	envs, err := models.GetEnvVariablesByAppID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get environment variables", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, envs, "Environment variables retrieved successfully", "")
}

func UpdateEnvVariable(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID    int64  `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 || req.Key == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID and key are required", "Missing fields")
		return
	}

	env, err := models.GetEnvVariableByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get environment variable", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, env.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	err = models.UpdateEnvVariable(req.ID, strings.TrimSpace(req.Key), req.Value)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update environment variable", err.Error())
		return
	}

	updatedEnv, err := models.GetEnvVariableByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch updated environment variable", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "update", "env_variable", &req.ID, map[string]interface{}{
		"app_id": env.AppID,
		"key":    req.Key,
	})

	response := map[string]interface{}{
		"envVariable":    updatedEnv,
		"actionRequired": "redeploy",
		"actionMessage":  "Environment variable changes require a full redeployment to take effect. Would you like to redeploy now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Environment variable updated successfully", "")
}

func DeleteEnvVariable(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID is required", "Missing fields")
		return
	}

	env, err := models.GetEnvVariableByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get environment variable", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, env.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	models.LogUserAudit(userInfo.ID, "delete", "env_variable", &req.ID, map[string]interface{}{
		"app_id": env.AppID,
		"key":    env.Key,
	})

	err = models.DeleteEnvVariable(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete environment variable", err.Error())
		return
	}

	response := map[string]interface{}{
		"actionRequired": "redeploy",
		"actionMessage":  "Environment variable changes require a full redeployment to take effect. Would you like to redeploy now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Environment variable deleted successfully", "")
}
