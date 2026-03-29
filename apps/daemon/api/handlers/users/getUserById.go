package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

func GetUserById(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	userIDParam := r.URL.Query().Get("id")
	if userIDParam == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "Missing 'id' parameter")
		return
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid user ID", "User ID must be an integer")
		return
	}
	user, err := models.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", "No user exists with the given ID")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, user, "User retrieved successfully", "")
}
