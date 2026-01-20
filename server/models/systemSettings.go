package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SystemSettings struct {
	WildcardDomain        *string `json:"wildcardDomain"`
	MistAppName           string  `json:"mistAppName"`
	JwtSecret             string  `json:"-"`
	GithubWebhookSecret   string  `json:"-"`
	AllowedOrigins        string  `json:"allowedOrigins"`
	ProductionMode        bool    `json:"productionMode"`
	SecureCookies         bool    `json:"secureCookies"`
	AutoCleanupContainers bool    `json:"autoCleanupContainers"`
	AutoCleanupImages     bool    `json:"autoCleanupImages"`
}

type SystemSettingEntry struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SystemSettingEntry) TableName() string {
	return "system_settings"
}

func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GetSystemSetting(key string) (string, error) {
	var entry SystemSettingEntry
	err := db.Where("key = ?", key).First(&entry).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return entry.Value, nil
}

func SetSystemSetting(key, value string) error {
	entry := SystemSettingEntry{
		Key:   key,
		Value: value,
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&entry).Error
}

func GetSystemSettings() (*SystemSettings, error) {
	var settings SystemSettings

	wildcardVal, err := GetSystemSetting("wildcard_domain")
	if err != nil {
		return nil, err
	}
	if wildcardVal != "" {
		settings.WildcardDomain = &wildcardVal
	}

	mistAppName, err := GetSystemSetting("mist_app_name")
	if err != nil {
		return nil, err
	}
	if mistAppName == "" {
		mistAppName = "mist"
	}
	settings.MistAppName = mistAppName

	jwtSecret, err := GetSystemSetting("jwt_secret")
	if err != nil {
		return nil, err
	}
	if jwtSecret == "" {
		jwtSecret, err = generateRandomSecret(64)
		if err != nil {
			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
		}
		if err := SetSystemSetting("jwt_secret", jwtSecret); err != nil {
			return nil, fmt.Errorf("failed to save JWT secret: %w", err)
		}
		log.Info().Msg("Auto-generated JWT secret and saved to database")
	}
	settings.JwtSecret = jwtSecret

	githubSecret, err := GetSystemSetting("github_webhook_secret")
	if err != nil {
		return nil, err
	}
	settings.GithubWebhookSecret = githubSecret

	allowedOrigins, err := GetSystemSetting("allowed_origins")
	if err != nil {
		return nil, err
	}
	settings.AllowedOrigins = allowedOrigins

	prodMode, err := GetSystemSetting("production_mode")
	if err != nil {
		return nil, err
	}
	settings.ProductionMode = prodMode == "true"

	secureCookies, err := GetSystemSetting("secure_cookies")
	if err != nil {
		return nil, err
	}
	if secureCookies == "" {
		settings.SecureCookies = settings.ProductionMode
	} else {
		settings.SecureCookies = secureCookies == "true"
	}

	autoCleanupContainers, err := GetSystemSetting("auto_cleanup_containers")
	if err != nil {
		return nil, err
	}
	settings.AutoCleanupContainers = autoCleanupContainers == "true"

	autoCleanupImages, err := GetSystemSetting("auto_cleanup_images")
	if err != nil {
		return nil, err
	}
	settings.AutoCleanupImages = autoCleanupImages == "true"

	return &settings, nil
}

func UpdateSystemSettings(wildcardDomain *string, mistAppName string) (*SystemSettings, error) {
	wildcardValue := ""
	if wildcardDomain != nil {
		wildcardValue = *wildcardDomain
	}

	if err := SetSystemSetting("wildcard_domain", wildcardValue); err != nil {
		return nil, err
	}

	if err := SetSystemSetting("mist_app_name", mistAppName); err != nil {
		return nil, err
	}

	return GetSystemSettings()
}

func UpdateSecuritySettings(allowedOrigins string, productionMode, secureCookies bool) error {
	if err := SetSystemSetting("allowed_origins", allowedOrigins); err != nil {
		return err
	}

	prodModeStr := "false"
	if productionMode {
		prodModeStr = "true"
	}
	if err := SetSystemSetting("production_mode", prodModeStr); err != nil {
		return err
	}

	secureCookiesStr := "false"
	if secureCookies {
		secureCookiesStr = "true"
	}
	if err := SetSystemSetting("secure_cookies", secureCookiesStr); err != nil {
		return err
	}

	return nil
}

func UpdateDockerSettings(autoCleanupContainers, autoCleanupImages bool) error {
	cleanupContainersStr := "false"
	if autoCleanupContainers {
		cleanupContainersStr = "true"
	}
	if err := SetSystemSetting("auto_cleanup_containers", cleanupContainersStr); err != nil {
		return err
	}

	cleanupImagesStr := "false"
	if autoCleanupImages {
		cleanupImagesStr = "true"
	}
	if err := SetSystemSetting("auto_cleanup_images", cleanupImagesStr); err != nil {
		return err
	}

	return nil
}

func (s *SystemSettings) UpdateSystemSettings() error {
	wildcardValue := ""
	if s.WildcardDomain != nil {
		wildcardValue = *s.WildcardDomain
	}
	if err := SetSystemSetting("wildcard_domain", wildcardValue); err != nil {
		return err
	}

	if err := SetSystemSetting("mist_app_name", s.MistAppName); err != nil {
		return err
	}

	prodModeStr := "false"
	if s.ProductionMode {
		prodModeStr = "true"
	}
	if err := SetSystemSetting("production_mode", prodModeStr); err != nil {
		return err
	}

	secureCookiesStr := "false"
	if s.SecureCookies {
		secureCookiesStr = "true"
	}
	if err := SetSystemSetting("secure_cookies", secureCookiesStr); err != nil {
		return err
	}

	cleanupContainersStr := "false"
	if s.AutoCleanupContainers {
		cleanupContainersStr = "true"
	}
	if err := SetSystemSetting("auto_cleanup_containers", cleanupContainersStr); err != nil {
		return err
	}

	cleanupImagesStr := "false"
	if s.AutoCleanupImages {
		cleanupImagesStr = "true"
	}
	if err := SetSystemSetting("auto_cleanup_images", cleanupImagesStr); err != nil {
		return err
	}

	return nil
}

func GenerateAutoDomain(projectName, appName string) (string, error) {
	settings, err := GetSystemSettings()
	if err != nil {
		return "", err
	}

	if settings.WildcardDomain == nil || *settings.WildcardDomain == "" {
		return "", nil
	}

	wildcardDomain := *settings.WildcardDomain
	if len(wildcardDomain) > 0 && wildcardDomain[0] == '*' {
		wildcardDomain = wildcardDomain[1:]
	}
	if len(wildcardDomain) > 0 && wildcardDomain[0] == '.' {
		wildcardDomain = wildcardDomain[1:]
	}

	return projectName + "-" + appName + "." + wildcardDomain, nil
}

//############################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"crypto/rand"
// 	"database/sql"
// 	"encoding/base64"
// 	"fmt"

// 	"github.com/rs/zerolog/log"
// )

// type SystemSettings struct {
// 	WildcardDomain        *string `json:"wildcardDomain"`
// 	MistAppName           string  `json:"mistAppName"`
// 	JwtSecret             string  `json:"-"`
// 	GithubWebhookSecret   string  `json:"-"`
// 	AllowedOrigins        string  `json:"allowedOrigins"`
// 	ProductionMode        bool    `json:"productionMode"`
// 	SecureCookies         bool    `json:"secureCookies"`
// 	AutoCleanupContainers bool    `json:"autoCleanupContainers"`
// 	AutoCleanupImages     bool    `json:"autoCleanupImages"`
// }

// func generateRandomSecret(length int) (string, error) {
// 	bytes := make([]byte, length)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(bytes), nil
// }

// func GetSystemSetting(key string) (string, error) {
// 	var value string
// 	err := db.QueryRow(`SELECT value FROM system_settings WHERE key = ?`, key).Scan(&value)
// 	if err == sql.ErrNoRows {
// 		return "", nil
// 	}
// 	if err != nil {
// 		return "", err
// 	}
// 	return value, nil
// }

// func SetSystemSetting(key, value string) error {
// 	_, err := db.Exec(`
// 		INSERT INTO system_settings (key, value, updated_at)
// 		VALUES (?, ?, CURRENT_TIMESTAMP)
// 		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
// 	`, key, value, value)
// 	return err
// }

// func GetSystemSettings() (*SystemSettings, error) {
// 	var settings SystemSettings

// 	var wildcardDomain sql.NullString
// 	err := db.QueryRow(`SELECT value FROM system_settings WHERE key = ?`, "wildcard_domain").Scan(&wildcardDomain)
// 	if err != nil && err != sql.ErrNoRows {
// 		return nil, err
// 	}
// 	if wildcardDomain.Valid && wildcardDomain.String != "" {
// 		settings.WildcardDomain = &wildcardDomain.String
// 	}

// 	var mistAppName string
// 	err = db.QueryRow(`SELECT value FROM system_settings WHERE key = ?`, "mist_app_name").Scan(&mistAppName)
// 	if err != nil && err != sql.ErrNoRows {
// 		return nil, err
// 	}
// 	if mistAppName == "" {
// 		mistAppName = "mist"
// 	}
// 	settings.MistAppName = mistAppName

// 	jwtSecret, err := GetSystemSetting("jwt_secret")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if jwtSecret == "" {
// 		jwtSecret, err = generateRandomSecret(64)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
// 		}
// 		if err := SetSystemSetting("jwt_secret", jwtSecret); err != nil {
// 			return nil, fmt.Errorf("failed to save JWT secret: %w", err)
// 		}
// 		log.Info().Msg("Auto-generated JWT secret and saved to database")
// 	}
// 	settings.JwtSecret = jwtSecret

// 	githubSecret, err := GetSystemSetting("github_webhook_secret")
// 	if err != nil {
// 		return nil, err
// 	}
// 	settings.GithubWebhookSecret = githubSecret

// 	allowedOrigins, err := GetSystemSetting("allowed_origins")
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Default to empty string - same-origin requests are always allowed
// 	// Users only need to configure this for cross-origin requests
// 	settings.AllowedOrigins = allowedOrigins

// 	prodMode, err := GetSystemSetting("production_mode")
// 	if err != nil {
// 		return nil, err
// 	}
// 	settings.ProductionMode = prodMode == "true"

// 	secureCookies, err := GetSystemSetting("secure_cookies")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if secureCookies == "" {
// 		settings.SecureCookies = settings.ProductionMode
// 	} else {
// 		settings.SecureCookies = secureCookies == "true"
// 	}

// 	autoCleanupContainers, err := GetSystemSetting("auto_cleanup_containers")
// 	if err != nil {
// 		return nil, err
// 	}
// 	settings.AutoCleanupContainers = autoCleanupContainers == "true"

// 	autoCleanupImages, err := GetSystemSetting("auto_cleanup_images")
// 	if err != nil {
// 		return nil, err
// 	}
// 	settings.AutoCleanupImages = autoCleanupImages == "true"

// 	return &settings, nil
// }

// func UpdateSystemSettings(wildcardDomain *string, mistAppName string) (*SystemSettings, error) {
// 	wildcardValue := ""
// 	if wildcardDomain != nil {
// 		wildcardValue = *wildcardDomain
// 	}
// 	_, err := db.Exec(`
// 		INSERT INTO system_settings (key, value, updated_at)
// 		VALUES (?, ?, CURRENT_TIMESTAMP)
// 		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
// 	`, "wildcard_domain", wildcardValue, wildcardValue)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = db.Exec(`
// 		INSERT INTO system_settings (key, value, updated_at)
// 		VALUES (?, ?, CURRENT_TIMESTAMP)
// 		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
// 	`, "mist_app_name", mistAppName, mistAppName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return GetSystemSettings()
// }

// func UpdateSecuritySettings(allowedOrigins string, productionMode, secureCookies bool) error {
// 	if err := SetSystemSetting("allowed_origins", allowedOrigins); err != nil {
// 		return err
// 	}

// 	prodModeStr := "false"
// 	if productionMode {
// 		prodModeStr = "true"
// 	}
// 	if err := SetSystemSetting("production_mode", prodModeStr); err != nil {
// 		return err
// 	}

// 	secureCookiesStr := "false"
// 	if secureCookies {
// 		secureCookiesStr = "true"
// 	}
// 	if err := SetSystemSetting("secure_cookies", secureCookiesStr); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func UpdateDockerSettings(autoCleanupContainers, autoCleanupImages bool) error {
// 	cleanupContainersStr := "false"
// 	if autoCleanupContainers {
// 		cleanupContainersStr = "true"
// 	}
// 	if err := SetSystemSetting("auto_cleanup_containers", cleanupContainersStr); err != nil {
// 		return err
// 	}

// 	cleanupImagesStr := "false"
// 	if autoCleanupImages {
// 		cleanupImagesStr = "true"
// 	}
// 	if err := SetSystemSetting("auto_cleanup_images", cleanupImagesStr); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // UpdateSystemSettings updates all settings from the SystemSettings struct
// func (s *SystemSettings) UpdateSystemSettings() error {
// 	wildcardValue := ""
// 	if s.WildcardDomain != nil {
// 		wildcardValue = *s.WildcardDomain
// 	}
// 	if err := SetSystemSetting("wildcard_domain", wildcardValue); err != nil {
// 		return err
// 	}

// 	if err := SetSystemSetting("mist_app_name", s.MistAppName); err != nil {
// 		return err
// 	}

// 	prodModeStr := "false"
// 	if s.ProductionMode {
// 		prodModeStr = "true"
// 	}
// 	if err := SetSystemSetting("production_mode", prodModeStr); err != nil {
// 		return err
// 	}

// 	secureCookiesStr := "false"
// 	if s.SecureCookies {
// 		secureCookiesStr = "true"
// 	}
// 	if err := SetSystemSetting("secure_cookies", secureCookiesStr); err != nil {
// 		return err
// 	}

// 	cleanupContainersStr := "false"
// 	if s.AutoCleanupContainers {
// 		cleanupContainersStr = "true"
// 	}
// 	if err := SetSystemSetting("auto_cleanup_containers", cleanupContainersStr); err != nil {
// 		return err
// 	}

// 	cleanupImagesStr := "false"
// 	if s.AutoCleanupImages {
// 		cleanupImagesStr = "true"
// 	}
// 	if err := SetSystemSetting("auto_cleanup_images", cleanupImagesStr); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func GenerateAutoDomain(projectName, appName string) (string, error) {
// 	settings, err := GetSystemSettings()
// 	if err != nil {
// 		return "", err
// 	}

// 	if settings.WildcardDomain == nil || *settings.WildcardDomain == "" {
// 		return "", nil
// 	}

// 	wildcardDomain := *settings.WildcardDomain
// 	if len(wildcardDomain) > 0 && wildcardDomain[0] == '*' {
// 		wildcardDomain = wildcardDomain[1:]
// 	}
// 	if len(wildcardDomain) > 0 && wildcardDomain[0] == '.' {
// 		wildcardDomain = wildcardDomain[1:]
// 	}

// 	return projectName + "-" + appName + "." + wildcardDomain, nil
// }
