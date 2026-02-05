package applications

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func CreateDomain(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID  int64  `json:"appId"`
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 || req.Domain == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID and domain are required", "Missing fields")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	domain, err := models.CreateDomain(req.AppID, strings.TrimSpace(req.Domain))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create domain", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "create", "domain", &domain.ID, map[string]interface{}{
		"appId":  req.AppID,
		"domain": domain.Domain,
	})

	response := map[string]interface{}{
		"domain":         domain,
		"actionRequired": "restart",
		"actionMessage":  "Domain changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Domain created successfully", "")
}

func GetDomains(w http.ResponseWriter, r *http.Request) {
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

	domains, err := models.GetDomainsByAppID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domains", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, domains, "Domains retrieved successfully", "")
}

func UpdateDomain(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID     int64  `json:"id"`
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 || req.Domain == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID and domain are required", "Missing fields")
		return
	}

	domain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domain", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, domain.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	oldDomain := domain.Domain

	err = models.UpdateDomain(req.ID, strings.TrimSpace(req.Domain))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update domain", err.Error())
		return
	}

	updatedDomain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve updated domain", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "update", "domain", &req.ID, map[string]interface{}{
		"appId": domain.AppID,
		"before": map[string]interface{}{
			"domain": oldDomain,
		},
		"after": map[string]interface{}{
			"domain": strings.TrimSpace(req.Domain),
		},
	})

	response := map[string]interface{}{
		"domain":         updatedDomain,
		"actionRequired": "restart",
		"actionMessage":  "Domain changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Domain updated successfully", "")
}

func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID is required", "Missing fields")
		return
	}

	domain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domain", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, domain.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to modify this application", "Forbidden")
		return
	}

	err = models.DeleteDomain(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete domain", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "delete", "domain", &req.ID, map[string]interface{}{
		"appId":  domain.AppID,
		"domain": domain.Domain,
	})

	response := map[string]interface{}{
		"actionRequired": "restart",
		"actionMessage":  "Domain changes require restarting the container to take effect. Would you like to restart now?",
	}
	handlers.SendResponse(w, http.StatusOK, true, response, "Domain deleted successfully", "")
}

func VerifyDomainDNS(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID is required", "Missing fields")
		return
	}

	domain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domain", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, domain.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to access this application", "Forbidden")
		return
	}

	valid, validationErr := utils.ValidateDNSWithTimeout(domain.Domain, 5*time.Second)

	var errorMsg *string
	if validationErr != nil {
		errStr := validationErr.Error()
		errorMsg = &errStr
	}

	if err := models.UpdateDomainDnsStatus(req.ID, valid, errorMsg); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update DNS status", err.Error())
		return
	}

	updatedDomain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve domain", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "verify", "domain", &req.ID, map[string]interface{}{
		"appId":         domain.AppID,
		"domain":        domain.Domain,
		"dnsConfigured": valid,
	})

	serverIP, _ := utils.GetServerIP()

	response := map[string]interface{}{
		"domain":   updatedDomain,
		"valid":    valid,
		"serverIP": serverIP,
	}

	if !valid {
		response["error"] = errorMsg
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "DNS verification completed", "")
}

func GetDNSInstructions(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.ID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "ID is required", "Missing fields")
		return
	}

	domain, err := models.GetDomainByID(req.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get domain", err.Error())
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, domain.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to access this application", "Forbidden")
		return
	}

	serverIP, err := utils.GetServerIP()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get server IP", err.Error())
		return
	}

	instructions := map[string]interface{}{
		"domain":   domain.Domain,
		"serverIP": serverIP,
		"records": []map[string]string{
			{
				"type":  "A",
				"name":  "@",
				"value": serverIP,
			},
			{
				"type":  "A",
				"name":  "www",
				"value": serverIP,
			},
		},
	}

	handlers.SendResponse(w, http.StatusOK, true, instructions, "DNS instructions retrieved", "")
}
