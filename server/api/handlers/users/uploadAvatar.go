package users

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"github.com/rs/zerolog/log"
)

func init() {
	avatarDir := constants.Constants["AvatarDirPath"].(string)
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		log.Warn().Err(err).Str("path", avatarDir).Msg("Failed to create avatar directory")
	}
}

func getAvatarDir() string {
	return constants.Constants["AvatarDirPath"].(string)
}

func getMaxAvatarSize() int64 {
	return int64(constants.Constants["MaxAvatarSize"].(int))
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	contextUser, ok := middleware.GetUser(r)
	if !ok || contextUser == nil {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized", "")
		return
	}
	userID := contextUser.ID

	maxSize := getMaxAvatarSize()
	if err := r.ParseMultipartForm(maxSize); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "File too large or invalid", err.Error())
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Failed to get file from request", err.Error())
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid file type. Only JPG, PNG, GIF, and WebP are allowed", "")
		return
	}

	if header.Size > maxSize {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "File size exceeds 5MB limit", "")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = getExtensionFromContentType(contentType)
	}

	filename := fmt.Sprintf("avatar_%d_%s%s", userID, utils.GenerateRandomString(8), ext)
	avatarDir := getAvatarDir()
	filePath := filepath.Join(avatarDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to save file", err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to save file", err.Error())
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		os.Remove(filePath)
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", err.Error())
		return
	}

	if user.AvatarURL != nil && *user.AvatarURL != "" {
		oldFilename := filepath.Base(*user.AvatarURL)
		oldPath := filepath.Join(avatarDir, oldFilename)
		if _, err := os.Stat(oldPath); err == nil {
			os.Remove(oldPath)
		}
	}

	avatarURL := fmt.Sprintf("/uploads/avatar/%s", filename)
	user.AvatarURL = &avatarURL

	if err := models.UpdateUser(user); err != nil {
		os.Remove(filePath)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update user avatar", err.Error())
		return
	}

	models.LogUserAudit(userID, "update", "user", &userID, map[string]interface{}{
		"action":    "avatar_upload",
		"avatarUrl": avatarURL,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"avatarUrl": avatarURL,
		"user":      user,
	}, "Avatar uploaded successfully", "")
}

func DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r)
	if !ok || user == nil {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized", "")
		return
	}

	user, err := models.GetUserByID(user.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", err.Error())
		return
	}

	if user.AvatarURL != nil && *user.AvatarURL != "" {
		filename := filepath.Base(*user.AvatarURL)
		avatarDir := getAvatarDir()
		filePath := filepath.Join(avatarDir, filename)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				log.Warn().Err(err).Str("file", filePath).Msg("Failed to delete avatar file")
			}
		}
	}

	user.AvatarURL = nil

	if err := models.UpdateUser(user); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update user", err.Error())
		return
	}

	models.LogUserAudit(user.ID, "delete", "user", &user.ID, map[string]interface{}{
		"action": "avatar_delete",
	})

	handlers.SendResponse(w, http.StatusOK, true, user, "Avatar deleted successfully", "")
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
