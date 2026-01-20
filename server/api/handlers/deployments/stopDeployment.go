package deployments

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/rs/zerolog/log"
)

type stopDeployment struct {
	DeploymentID int64 `json:"deploymentId"`
}

func StopDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.SendResponse(w, http.StatusMethodNotAllowed, false, nil, "Method not allowed", "Only POST method is allowed")
		return
	}
	var req stopDeployment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	deployment, err := models.GetDeploymentByID(req.DeploymentID)
	if err != nil {
		log.Error().Err(err).Int64("deployment_id", req.DeploymentID).Msg("Failed to get deployment details")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get deployment details", err.Error())
		return
	}

	logPath := filepath.Join(constants.Constants["LogPath"].(string), deployment.CommitHash+strconv.FormatInt(req.DeploymentID, 10)+"_build_logs")

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Int64("deployment_id", req.DeploymentID).Str("log_path", logPath).Msg("Failed to open log file")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to open log file", err.Error())
		return
	}
	defer file.Close()

	logLine := "deployment was stopped by user\n"
	_, err = file.WriteString(logLine)
	if err != nil {
		log.Error().Err(err).Int64("deployment_id", req.DeploymentID).Str("log_path", logPath).Msg("Failed to write to log file")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to write to log file", err.Error())
		return
	}

	wasRunning := queue.Cancel(req.DeploymentID)
	message := "Deployment marked as stopped"
	if wasRunning {
		message = "Deployment aborted immediately"
	}

	errorMsg := "deployment stopped by user"
	err = models.UpdateDeploymentStatus(req.DeploymentID, "stopped", "stopped", deployment.Progress, &errorMsg)
	if err != nil {
		log.Error().Err(err).Int64("deployment_id", req.DeploymentID).Msg("Failed to update deployment status")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to stop deployment", err.Error())
		return
	}

	log.Info().Int64("deployment_id", req.DeploymentID).Bool("was_running", wasRunning).Msg("Deployment stopped successfully")
	handlers.SendResponse(w, http.StatusOK, true, nil, message, "")
}
