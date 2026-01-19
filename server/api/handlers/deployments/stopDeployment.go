package deployments

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
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
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get deployment details", err.Error())
		return
	}
	err = models.UpdateDeploymentStatus(req.DeploymentID, "stopped", "pending", deployment.Progress, deployment.ErrorMessage)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to stop deployment", err.Error())
		return
	}
	wasRunning := queue.Cancel(req.DeploymentID)
	message := "Deployment marked as stopped"
	if wasRunning {
		message = "Deployment aborted immediately"
	}

	handlers.SendResponse(w, http.StatusOK, true, nil, message, "")
}
