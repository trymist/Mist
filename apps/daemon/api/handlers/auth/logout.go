package auth

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if ok && user != nil {
		models.LogUserAudit(user.ID, "logout", "user", &user.ID, map[string]interface{}{
			"username": user.Username,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "mist_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    nil,
		"message": "Logged out successfully",
		"error":   "",
	})
}
