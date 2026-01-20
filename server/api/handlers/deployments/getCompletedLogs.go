package deployments

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

type GetDeploymentLogsResponse struct {
	Deployment *models.Deployment `json:"deployment"`
	Logs       string             `json:"logs"`
}

func GetCompletedDeploymentLogsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Authentication required", "")
		return
	}

	depIdstr := r.URL.Query().Get("id")
	depId, err := strconv.ParseInt(depIdstr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid deployment id", err.Error())
		return
	}

	dep, err := models.GetDeploymentByID(depId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "deployment not found", err.Error())
		return
	}

	app, err := models.GetApplicationByID(dep.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}

	project, err := models.GetProjectByID(app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get project", err.Error())
		return
	}

	hasAccess := false
	if project.OwnerID == currentUser.ID {
		hasAccess = true
	} else {
		for _, member := range project.ProjectMembers {
			if member.ID == currentUser.ID {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "You don't have access to this deployment")
		return
	}

	if dep.Status != "success" && dep.Status != "failed" && dep.Status != "stopped" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "deployment is still in progress, use WebSocket endpoint", "")
		return
	}

	logPath := docker.GetLogsPath(dep.CommitHash, depId)
	logContent := ""

	if _, err := os.Stat(logPath); err == nil {
		file, err := os.Open(logPath)
		if err != nil {
			log.Error().Err(err).Int64("deployment_id", depId).Msg("Failed to open log file")
		} else {
			defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				log.Error().Err(err).Int64("deployment_id", depId).Msg("Failed to read log file")
			} else {
				logContent = string(content)
			}
		}
	} else {
		log.Warn().Int64("deployment_id", depId).Str("log_path", logPath).Msg("Deployment log file not found")
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Deployment log file not found", "")
		return
	}

	response := GetDeploymentLogsResponse{
		Deployment: dep,
		Logs:       logContent,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "Deployment logs retrieved successfully",
	})
}
