package settings

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func GetSystemSettings(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can view system settings", "Forbidden")
		return
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve system settings", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, settings, "System settings retrieved successfully", "")
}

func UpdateSystemSettings(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can update system settings", "Forbidden")
		return
	}

	var req struct {
		WildcardDomain        *string `json:"wildcardDomain"`
		MistAppName           *string `json:"mistAppName"`
		AllowedOrigins        *string `json:"allowedOrigins"`
		ProductionMode        *bool   `json:"productionMode"`
		SecureCookies         *bool   `json:"secureCookies"`
		AutoCleanupContainers *bool   `json:"autoCleanupContainers"`
		AutoCleanupImages     *bool   `json:"autoCleanupImages"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve current settings", err.Error())
		return
	}

	if req.WildcardDomain != nil || req.MistAppName != nil {
		wildcardDomain := settings.WildcardDomain
		if req.WildcardDomain != nil {
			if *req.WildcardDomain == "" {
				wildcardDomain = nil
			} else {
				wildcardDomain = req.WildcardDomain
			}
		}

		mistAppName := settings.MistAppName
		if req.MistAppName != nil {
			if *req.MistAppName == "" {
				handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Mist app name cannot be empty", "Invalid value")
				return
			}
			mistAppName = *req.MistAppName
		}

		settings, err = models.UpdateSystemSettings(wildcardDomain, mistAppName)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update system settings", err.Error())
			return
		}
	}

	if req.AllowedOrigins != nil || req.ProductionMode != nil || req.SecureCookies != nil {
		currentSettings, err := models.GetSystemSettings()
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve current settings", err.Error())
			return
		}

		allowedOrigins := currentSettings.AllowedOrigins
		if req.AllowedOrigins != nil {
			allowedOrigins = *req.AllowedOrigins
		}

		productionMode := currentSettings.ProductionMode
		if req.ProductionMode != nil {
			productionMode = *req.ProductionMode
		}

		secureCookies := currentSettings.SecureCookies
		if req.SecureCookies != nil {
			secureCookies = *req.SecureCookies
		}

		if err := models.UpdateSecuritySettings(allowedOrigins, productionMode, secureCookies); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update security settings", err.Error())
			return
		}

		settings, err = models.GetSystemSettings()
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve updated settings", err.Error())
			return
		}
	}

	if req.AutoCleanupContainers != nil || req.AutoCleanupImages != nil {
		currentSettings, err := models.GetSystemSettings()
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve current settings", err.Error())
			return
		}

		autoCleanupContainers := currentSettings.AutoCleanupContainers
		if req.AutoCleanupContainers != nil {
			autoCleanupContainers = *req.AutoCleanupContainers
		}

		autoCleanupImages := currentSettings.AutoCleanupImages
		if req.AutoCleanupImages != nil {
			autoCleanupImages = *req.AutoCleanupImages
		}

		if err := models.UpdateDockerSettings(autoCleanupContainers, autoCleanupImages); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update Docker settings", err.Error())
			return
		}

		settings, err = models.GetSystemSettings()
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve updated settings", err.Error())
			return
		}
	}

	if err := utils.GenerateDynamicConfig(settings.WildcardDomain, settings.MistAppName); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to generate Traefik configuration", err.Error())
		return
	}

	dummyID := int64(1)
	auditData := map[string]any{}
	if req.WildcardDomain != nil {
		auditData["wildcardDomain"] = *req.WildcardDomain
	}
	if req.MistAppName != nil {
		auditData["mistAppName"] = *req.MistAppName
	}
	if req.AllowedOrigins != nil {
		auditData["allowedOrigins"] = *req.AllowedOrigins
	}
	if req.ProductionMode != nil {
		auditData["productionMode"] = *req.ProductionMode
	}
	if req.SecureCookies != nil {
		auditData["secureCookies"] = *req.SecureCookies
	}
	if req.AutoCleanupContainers != nil {
		auditData["autoCleanupContainers"] = *req.AutoCleanupContainers
	}
	if req.AutoCleanupImages != nil {
		auditData["autoCleanupImages"] = *req.AutoCleanupImages
	}
	models.LogUserAudit(userInfo.ID, "update", "system_settings", &dummyID, auditData)

	handlers.SendResponse(w, http.StatusOK, true, settings, "System settings updated successfully", "")
}
