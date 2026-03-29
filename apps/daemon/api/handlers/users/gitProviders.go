package users

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetUserGitProviders(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	providers, err := models.GetGitProvidersByUser(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch git providers", err.Error())
		return
	}

	providersList := make([]map[string]interface{}, 0, len(providers))
	for _, p := range providers {
		providersList = append(providersList, map[string]interface{}{
			"id":       p.ID,
			"provider": p.Provider,
			"username": p.Username,
			"email":    p.Email,
		})
	}

	handlers.SendResponse(w, http.StatusOK, true, providersList, "Git providers fetched successfully", "")
}
