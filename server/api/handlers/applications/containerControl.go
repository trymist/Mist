package applications

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/compose"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
)

func StopContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	// containerName declaration moved after check

	if app.AppType == models.AppTypeCompose {
		path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
		err = compose.ComposeDown(path)
	} else {
		containerName := docker.GetContainerName(app.Name, appId)
		err = docker.StopContainer(containerName)
	}

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to stop application", err.Error())
		return
	}

	app.Status = models.StatusStopped
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	containerName := ""
	if app.AppType != models.AppTypeCompose {
		containerName = docker.GetContainerName(app.Name, appId)
	}

	models.LogUserAudit(userInfo.ID, "stop", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
		"app_type":       app.AppType,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Application stopped successfully",
	}, "Application stopped successfully", "")
}

func StartContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	var startErr error
	if app.AppType == models.AppTypeCompose {
		path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
		envVars, _ := models.GetEnvVariablesByAppID(appId)
		envMap := make(map[string]string)
		for _, e := range envVars {
			envMap[e.Key] = e.Value
		}
		startErr = compose.ComposeUp(path, envMap, nil)
	} else {
		containerName := docker.GetContainerName(app.Name, appId)
		startErr = docker.StartContainer(containerName)
	}

	if startErr != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to start application", startErr.Error())
		return
	}

	app.Status = models.StatusRunning
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	containerName := ""
	if app.AppType != models.AppTypeCompose {
		containerName = docker.GetContainerName(app.Name, appId)
	}

	models.LogUserAudit(userInfo.ID, "start", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
		"app_type":       app.AppType,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Application started successfully",
	}, "Application started successfully", "")
}

func RestartContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	var restartErr error
	if app.AppType == models.AppTypeCompose {
		path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
		restartErr = compose.ComposeRestart(path)
	} else {
		containerName := docker.GetContainerName(app.Name, appId)
		restartErr = docker.RestartContainer(containerName)
	}

	if restartErr != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to restart application", restartErr.Error())
		return
	}

	app.Status = models.StatusRunning
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	containerName := ""
	if app.AppType != models.AppTypeCompose {
		containerName = docker.GetContainerName(app.Name, appId)
	}

	models.LogUserAudit(userInfo.ID, "restart", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
		"app_type":       app.AppType,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Application restarted successfully",
	}, "Application restarted successfully", "")
}

func GetContainerStatusHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to view this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	// containerName declaration moved

	if app.AppType == models.AppTypeCompose {
		path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
		status, err := compose.GetComposeStatus(path, app.Name)
		if err != nil {
			// If error, maybe it's just not running or folder doesn't exist
			handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
				"name":   app.Name,
				"status": "Unknown",
				"state":  "stopped",
				"error":  err.Error(),
			}, "Compose status retrieval failed (app might be stopped)", "")
			return
		}
		handlers.SendResponse(w, http.StatusOK, true, status, "Compose status retrieved successfully", "")
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)
	status, err := docker.GetContainerStatus(containerName)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get container status", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, status, "Container status retrieved successfully", "")
}

func GetContainerLogsHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to view this application", "Forbidden")
		return
	}

	tailStr := r.URL.Query().Get("tail")
	tail := 100
	if tailStr != "" {
		if parsedTail, err := strconv.Atoi(tailStr); err == nil && parsedTail > 0 {
			tail = parsedTail
		}
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	if app.AppType == models.AppTypeCompose {
		path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
		logs, err := compose.GetComposeLogs(path, tail)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get compose logs", err.Error())
			return
		}
		handlers.SendResponse(w, http.StatusOK, true, map[string]any{
			"logs": logs,
		}, "Compose logs retrieved successfully", "")
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	logs, err := docker.GetContainerLogs(containerName, tail)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get container logs", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"logs": logs,
	}, "Container logs retrieved successfully", "")
}
