package auditlogs

import (
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetAuditLogsByResource(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Access denied", "Only admins can view audit logs")
		return
	}

	resourceType := r.URL.Query().Get("resourceType")
	resourceIDStr := r.URL.Query().Get("resourceId")

	if resourceType == "" || resourceIDStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Resource type and ID are required", "Missing parameters")
		return
	}

	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid resource ID", err.Error())
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

	logs, err := models.GetAuditLogsByResource(resourceType, resourceID, limit, offset)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve audit logs", err.Error())
		return
	}

	response := map[string]interface{}{
		"logs":   logs,
		"limit":  limit,
		"offset": offset,
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "Audit logs retrieved successfully", "")
}
