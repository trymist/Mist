package auditlogs

import (
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetAllAuditLogs(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "Only admins can view audit logs")
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	resourceType := r.URL.Query().Get("resourceType")

	var logs []models.AuditLog
	var err error

	if resourceType != "" {
		logs, err = models.GetAuditLogsByResourceType(resourceType, limit, offset)
	} else {
		logs, err = models.GetAllAuditLogs(limit, offset)
	}

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve audit logs", err.Error())
		return
	}

	total, err := models.GetAuditLogsCount()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get audit logs count", err.Error())
		return
	}

	response := map[string]interface{}{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "Audit logs retrieved successfully", "")
}
