package settings

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
)

func DockerCleanup(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can perform Docker cleanup", "Forbidden")
		return
	}

	var req struct {
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	var result string
	var cleanupErr error

	switch req.Type {
	case "containers":
		cleanupErr = docker.CleanupStoppedContainers()
		if cleanupErr == nil {
			result = "Successfully cleaned up stopped containers"
		}
	case "images":
		cleanupErr = docker.CleanupDanglingImages()
		if cleanupErr == nil {
			result = "Successfully cleaned up dangling images"
		}
	case "system":
		result, cleanupErr = docker.SystemPrune()
	case "system-all":
		result, cleanupErr = docker.SystemPruneAll()
	default:
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid cleanup type", "Type must be: containers, images, system, or system-all")
		return
	}

	if cleanupErr != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Docker cleanup failed", cleanupErr.Error())
		return
	}

	dummyID := int64(1)
	models.LogUserAudit(userInfo.ID, "cleanup", "docker", &dummyID, map[string]any{
		"type":   req.Type,
		"result": result,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]string{
		"message": result,
		"type":    req.Type,
	}, result, "")
}
