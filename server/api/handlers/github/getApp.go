package github

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetApp(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok || userInfo == nil {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized", "user not authenticated")
		return
	}

	app, isInstalled, err := models.GetApp(int(userInfo.ID))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	if app.ID == 0 {
		handlers.SendResponse(w, http.StatusNotFound, true, map[string]interface{}{
			"app":         nil,
			"isInstalled": false,
		}, "GitHub App not found", "")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"app":         app,
		"isInstalled": isInstalled,
	}, "GitHub App retrieved successfully", "")
}
