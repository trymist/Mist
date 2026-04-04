package users

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
)

type UpdatePasswordRequest struct {
	UserID          int64  `json:"userId"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.UserID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Current and new password are required", "")
		return
	}

	user, err := models.GetUserByID(req.UserID)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", err.Error())
		return
	}

	if !user.MatchPassword(req.CurrentPassword) {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Current password is incorrect", "")
		return
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to hash password", err.Error())
		return
	}

	if err := user.UpdatePassword(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update password", err.Error())
		return
	}

	models.LogUserAudit(req.UserID, "update", "user", &req.UserID, map[string]interface{}{
		"action": "password_change",
	})

	handlers.SendResponse(w, http.StatusOK, true, nil, "Password updated successfully", "")
}
