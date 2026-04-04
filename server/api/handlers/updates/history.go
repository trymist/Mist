package updates

import (
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

func GetUpdateHistory(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userInfo.ID).Msg("Failed to verify user role for update history")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		log.Warn().Int64("user_id", userInfo.ID).Str("role", role).Msg("Non-owner attempted to view update history")
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can view update history", "Forbidden")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := models.GetUpdateLogs(limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve update logs")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve update history", err.Error())
		return
	}

	log.Info().Int64("user_id", userInfo.ID).Int("count", len(logs)).Msg("Update history retrieved")
	handlers.SendResponse(w, http.StatusOK, true, logs, "Update history retrieved successfully", "")
}

func GetUpdateLogByID(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userInfo.ID).Msg("Failed to verify user role for update log")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		log.Warn().Int64("user_id", userInfo.ID).Str("role", role).Msg("Non-owner attempted to view update log")
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can view update logs", "Forbidden")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing id parameter", "")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid id parameter", err.Error())
		return
	}

	updateLog, err := models.GetUpdateLogByID(id)
	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to retrieve update log")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve update log", err.Error())
		return
	}

	if updateLog == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Update log not found", "")
		return
	}

	log.Info().Int64("user_id", userInfo.ID).Int64("update_log_id", id).Msg("Update log retrieved")
	handlers.SendResponse(w, http.StatusOK, true, updateLog, "Update log retrieved successfully", "")
}
