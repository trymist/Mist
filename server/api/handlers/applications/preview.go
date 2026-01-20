package applications

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetPreviewURL(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to access this application", "Forbidden")
		return
	}

	domain, err := models.GetPrimaryDomainByAppID(req.AppID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app, appErr := models.GetApplicationByID(req.AppID)
			if appErr != nil {
				handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", appErr.Error())
				return
			}

			serverIP := getOutboundIP()

			port := 80
			if app.Port != nil {
				port = int(*app.Port)
			}

			previewURL := fmt.Sprintf("http://%s:%d", serverIP, port)
			handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
				"url":  previewURL,
				"type": "ip",
				"ip":   serverIP,
				"port": port,
			}, "Preview URL retrieved (using IP:port)", "")
			return
		}

		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domain", err.Error())
		return
	}

	previewURL := "http://" + domain.Domain
	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"url":    previewURL,
		"domain": domain.Domain,
		"type":   "domain",
	}, "Preview URL retrieved successfully", "")
}
