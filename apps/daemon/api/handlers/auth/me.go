package auth

import (
	"errors"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/store"
	"gorm.io/gorm"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("mist_token")
	setupRequired := store.IsSetupRequired()
	if err != nil {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "No auth cookie", "")
		return
	}

	tokenStr := cookie.Value

	claims, err := middleware.VerifyJWT(tokenStr)
	if err != nil {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "Invalid token", "")
		return
	}

	userId := claims.UserID
	user, err := models.GetUserByID(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "User not found", "")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}
	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": user}, "User fetched successfully", "")
}
