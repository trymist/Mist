package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/store"
	"gorm.io/gorm"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	if userData.Role != "owner" && userData.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}
	userIDParam := r.URL.Query().Get("id")
	if userIDParam == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "Missing 'id' parameter")
		return
	}

	if userIDParam == strconv.FormatInt(userData.ID, 10) {
		err := models.DeleteUserByID(userData.ID)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete user", err.Error())
			return
		}
		http.Redirect(w, r, "/api/auth/logout", http.StatusSeeOther)
		store.InitSetupRequired()
		return
	}
	id, err := strconv.Atoi(userIDParam)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid user ID", "User ID must be an integer")
		return
	}
	userToDeleteRole, err := models.GetUserRole(int64(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", "No such user")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve user role", err.Error())
		return
	}

	if userData.Role == "admin" && userToDeleteRole == "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized to delete owner", "Forbidden")
		return
	}
	if userData.Role == "admin" && userToDeleteRole == "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized to delete admin", "Forbidden")
		return
	}

	userToDelete, _ := models.GetUserByID(int64(id))

	err = models.DeleteUserByID(int64(id))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete user", err.Error())
		return
	}

	deletedUserID := int64(id)
	auditDetails := map[string]interface{}{
		"deleted_by": userData.Username,
		"role":       userToDeleteRole,
	}
	if userToDelete != nil {
		auditDetails["username"] = userToDelete.Username
		auditDetails["email"] = userToDelete.Email
	}
	models.LogUserAudit(userData.ID, "delete", "user", &deletedUserID, auditDetails)

	store.InitSetupRequired()
	handlers.SendResponse(w, http.StatusOK, true, nil, "User deleted successfully", "")

}
