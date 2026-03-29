package updates

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

// GetCurrentVersion returns the current version from the database
func GetCurrentVersion(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userInfo.ID).Msg("Failed to verify user role for version check")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		log.Warn().Int64("user_id", userInfo.ID).Str("role", role).Msg("Non-owner attempted to view version")
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can view version information", "Forbidden")
		return
	}

	currentVersion, err := models.GetSystemSetting("version")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve current version from database")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve current version", err.Error())
		return
	}

	if currentVersion == "" {
		currentVersion = "1.0.0"
		if err := models.SetSystemSetting("version", currentVersion); err != nil {
			log.Error().Err(err).Msg("Failed to set default version")
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to set default version", err.Error())
			return
		}
		log.Info().Msg("Set default version to 1.0.0")
	}

	log.Info().Str("version", currentVersion).Int64("user_id", userInfo.ID).Msg("Version retrieved")
	handlers.SendResponse(w, http.StatusOK, true, map[string]string{
		"version": currentVersion,
	}, "Current version retrieved successfully", "")
}

// CheckForUpdates checks GitHub releases for newer versions
func CheckForUpdates(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userInfo.ID).Msg("Failed to verify user role for update check")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		log.Warn().Int64("user_id", userInfo.ID).Str("role", role).Msg("Non-owner attempted to check for updates")
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can check for updates", "Forbidden")
		return
	}

	currentVersion, err := models.GetSystemSetting("version")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve current version for update check")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve current version", err.Error())
		return
	}

	if currentVersion == "" {
		currentVersion = "1.0.0"
	}

	log.Info().Str("current_version", currentVersion).Int64("user_id", userInfo.ID).Msg("Checking for updates")

	// Fetch latest release from GitHub
	resp, err := http.Get("https://api.github.com/repos/corecollectives/mist/releases/latest")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch latest release from GitHub")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to check for updates", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("GitHub API returned non-OK status")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch latest release", "GitHub API error")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read GitHub API response")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to read response", err.Error())
		return
	}

	var release GithubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		log.Error().Err(err).Msg("Failed to parse GitHub release data")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to parse release data", err.Error())
		return
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	updateAvailable := compareVersions(currentVersion, latestVersion)

	log.Info().
		Str("current_version", currentVersion).
		Str("latest_version", latestVersion).
		Bool("update_available", updateAvailable).
		Msg("Update check completed")

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"currentVersion":  currentVersion,
		"latestVersion":   latestVersion,
		"updateAvailable": updateAvailable,
		"releaseNotes":    release.Body,
		"releaseName":     release.Name,
	}, "Update check completed", "")
}

func compareVersions(current, latest string) bool {
	return latest > current
}
