package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var cred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}
	cred.Email = strings.ToLower(strings.TrimSpace(cred.Email))

	if cred.Email == "" || cred.Password == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Email and password are required", "Missing fields")
		return
	}
	user, err := models.GetUserByEmail(cred.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Invalid email or password", "Unauthorized")
		return
	} else if err != nil {
		log.Error().Err(err).Str("email", cred.Email).Msg("Database error during login")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}

	passwordMatch := user.MatchPassword(cred.Password)
	if !passwordMatch {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Invalid email or password", "Unauthorized")
		return
	}
	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT token during login")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to generate token", "Internal Server Error")
		return
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get system settings during login")
		settings = &models.SystemSettings{SecureCookies: false}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "mist_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   settings.SecureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 30,
	})

	models.LogUserAudit(user.ID, "login", "user", &user.ID, map[string]interface{}{
		"email": user.Email,
	})

	handlers.SendResponse(w, http.StatusOK, true, user, "Login successful", "")
}
