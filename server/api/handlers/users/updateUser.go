package users

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

type UpdateUserRequest struct {
	ID       int64   `json:"id"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Role     *string `json:"role,omitempty"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Authentication required", "")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "")
		return
	}

	if currentUser.ID != req.ID {
		if currentUser.Role != "owner" && currentUser.Role != "admin" {
			handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "You can only update your own profile")
			return
		}
	}

	if req.Role != nil && *req.Role != "" {
		if currentUser.Role != "owner" {
			handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "Only owner can change user roles")
			return
		}
	}

	user, err := models.GetUserByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", err.Error())
		return
	}

	oldUsername := user.Username
	oldEmail := user.Email
	oldRole := user.Role

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := models.UpdateUser(user); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update user", err.Error())
		return
	}

	models.LogUserAudit(currentUser.ID, "update", "user", &req.ID, map[string]interface{}{
		"before": map[string]interface{}{
			"username": oldUsername,
			"email":    oldEmail,
			"role":     oldRole,
		},
		"after": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})

	handlers.SendResponse(w, http.StatusOK, true, user, "User updated successfully", "")
}
