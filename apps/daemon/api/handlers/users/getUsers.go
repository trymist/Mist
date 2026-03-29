package users

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Authentication required", "")
		return
	}

	if currentUser.Role != "owner" && currentUser.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "Only administrators can view all users")
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve users", err.Error())
		return
	}
	if users == nil {
		users = []models.User{}
	}
	handlers.SendResponse(w, http.StatusOK, true, users, "Users retrieved successfully", "")
}
