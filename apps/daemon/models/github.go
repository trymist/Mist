package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GithubApp struct {
	ID int64 `gorm:"primaryKey;autoIncrement:true" json:"id"`

	AppID         int64   `gorm:"not null" json:"appId"`
	ClientID      string  `gorm:"not null" json:"clientId"`
	ClientSecret  string  `gorm:"not null" json:"clientSecret"`
	WebhookSecret string  `gorm:"not null" json:"webhookSecret"`
	PrivateKey    string  `gorm:"not null" json:"privateKey"`
	Name          *string `json:"name"`
	Slug          string  `gorm:"not null" json:"slug"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (GithubApp) TableName() string {
	return "github_app"
}

type GithubInstallation struct {
	InstallationID int64 `gorm:"primaryKey;autoIncrement:false" json:"installation_id"`

	AccountLogin   string    `json:"account_login"`
	AccountType    string    `json:"account_type"`
	AccessToken    string    `json:"access_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`

	UserID int `gorm:"index" json:"user_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (app *GithubApp) InsertInDB() error {
	return db.Create(app).Error
}

func CheckIfAppExists() (bool, error) {
	var count int64
	err := db.Model(&GithubApp{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetApp(userID int) (GithubApp, bool, error) {
	var app GithubApp
	err := db.First(&app, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return GithubApp{}, false, nil
		}
		return GithubApp{}, false, err
	}

	var count int64
	err = db.Model(&GithubInstallation{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	isInstalled := count > 0
	return app, isInstalled, err
}

func (i *GithubInstallation) InsertOrReplace() error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "installation_id"}},
		UpdateAll: true,
	}).Create(i).Error
}

func GetInstallationID(userID int) (int64, error) {
	var inst GithubInstallation
	err := db.Select("installation_id").
		Where("user_id = ?", userID).
		First(&inst).Error
	if err != nil {
		return 0, err
	}
	return inst.InstallationID, nil
}

func GetInstallationByUserID(userID int) (GithubInstallation, error) {
	var inst GithubInstallation
	err := db.Where("user_id = ?", userID).First(&inst).Error
	if err != nil {
		return GithubInstallation{}, err
	}
	return inst, nil
}

func GetInstallationToken(installationID int64) (string, string, int, error) {
	type Result struct {
		AccessToken    string
		TokenExpiresAt time.Time
		AppID          int
	}

	var res Result
	err := db.Table("github_installations as i").
		Select("i.access_token, i.token_expires_at, a.app_id").
		Joins("JOIN github_app a ON a.id = 1").
		Where("i.installation_id = ?", installationID).
		Scan(&res).Error

	if err != nil {
		return "", "", 0, err
	}

	return res.AccessToken, res.TokenExpiresAt.Format(time.RFC3339), res.AppID, nil
}

func UpdateInstallationToken(installationID int64, token string, newExpiry time.Time) error {
	return db.Model(&GithubInstallation{InstallationID: installationID}).
		Updates(map[string]interface{}{
			"access_token":     token,
			"token_expires_at": newExpiry,
			"updated_at":       time.Now(),
		}).Error
}

func GetGithubAppCredentials(appID int) (string, string, error) {
	var app GithubApp
	err := db.Select("app_id, private_key").First(&app, 1).Error
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%d", app.AppID), app.PrivateKey, nil
}

func GetInstallationIDByUserID(userID int64) (int64, error) {
	var inst GithubInstallation
	err := db.Select("installation_id").Where("user_id = ?", userID).First(&inst).Error
	if err != nil {
		return 0, err
	}
	return inst.InstallationID, nil
}

func GetGithubAppIDAndPrivateKey() (int64, string, error) {
	var app GithubApp
	err := db.Select("app_id, private_key").First(&app).Error
	if err != nil {
		return 0, "", err
	}
	return app.AppID, app.PrivateKey, nil
}

//########################################################################################################################
//ARCHIVED CODE BELOW-------------------------------->

// package models

// import (
// 	"database/sql"
// 	"time"
// )

// type GithubApp struct {
// 	ID            int64     `json:"id"`
// 	AppID         int64     `json:"appId"`
// 	Name          *string   `json:"name"`
// 	Slug          string    `json:"slug"`
// 	ClientID      string    `json:"clientId"`
// 	ClientSecret  string    `json:"clientSecret"`
// 	WebhookSecret string    `json:"webhookSecret"`
// 	PrivateKey    string    `json:"privateKey"`
// 	CreatedAt     time.Time `json:"createdAt"`
// }

// func (app *GithubApp) InsertInDB() error {
// 	_, err := db.Exec(`
// 		INSERT INTO github_app (app_id, client_id, client_secret, webhook_secret, private_key, name, slug)
// 		VALUES (?, ?, ?, ?, ?, ?, ?)
// 	`, app.AppID, app.ClientID, app.ClientSecret, app.WebhookSecret, app.PrivateKey, app.Name, app.Slug)
// 	return err
// }

// func CheckIfAppExists() (bool, error) {
// 	var count int
// 	err := db.QueryRow(`SELECT COUNT(*) FROM github_app`).Scan(&count)
// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

// func GetApp(userID int) (GithubApp, bool, error) {
// 	query := `
// 		SELECT
// 			a.id,
// 			a.name,
// 			a.app_id,
// 			a.client_id,
// 			a.slug,
// 			a.created_at,
// 			CASE WHEN i.installation_id IS NOT NULL THEN 1 ELSE 0 END AS is_installed
// 		FROM github_app a
// 		LEFT JOIN github_installations i ON i.user_id = ?
// 		WHERE a.id = 1
// 	`

// 	row := db.QueryRow(query, userID)

// 	var app GithubApp
// 	var isInstalled bool

// 	err := row.Scan(
// 		&app.ID,
// 		&app.Name,
// 		&app.AppID,
// 		&app.ClientID,
// 		&app.Slug,
// 		&app.CreatedAt,
// 		&isInstalled,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return GithubApp{}, false, nil
// 		}
// 		return GithubApp{}, false, err
// 	}

// 	return app, isInstalled, nil
// }

// type GithubInstallation struct {
// 	InstallationID int64     `json:"installation_id"`
// 	AccountLogin   string    `json:"account_login"`
// 	AccountType    string    `json:"account_type"`
// 	AccessToken    string    `json:"access_token"`
// 	TokenExpiresAt time.Time `json:"token_expires_at"`
// 	UserID         int       `json:"user_id"`
// 	CreatedAt      time.Time `json:"created_at"`
// 	UpdatedAt      time.Time `json:"updated_at"`
// }

// func (i *GithubInstallation) InsertOrReplace() error {
// 	_, err := db.Exec(`
// 		INSERT OR REPLACE INTO github_installations
// 		(installation_id, account_login, account_type, access_token, token_expires_at, user_id, updated_at)
// 		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
// 	`, i.InstallationID, i.AccountLogin, i.AccountType, i.AccessToken, i.TokenExpiresAt, i.UserID)
// 	return err
// }

// func GetInstallationID(userID int) (int64, error) {
// 	var installationID int64
// 	err := db.QueryRow(`
// 		SELECT installation_id FROM github_installations
// 		WHERE user_id = ?
// 	`, userID).Scan(&installationID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return installationID, nil
// }

// func GetInstallationByUserID(userID int) (GithubInstallation, error) {
// 	var inst GithubInstallation
// 	err := db.QueryRow(`
// 		SELECT installation_id, account_login, account_type, access_token, token_expires_at, user_id, created_at, updated_at
// 		FROM github_installations
// 		WHERE user_id = ?
// 	`, userID).Scan(
// 		&inst.InstallationID,
// 		&inst.AccountLogin,
// 		&inst.AccountType,
// 		&inst.AccessToken,
// 		&inst.TokenExpiresAt,
// 		&inst.UserID,
// 		&inst.CreatedAt,
// 		&inst.UpdatedAt,
// 	)
// 	if err != nil {
// 		return GithubInstallation{}, err
// 	}
// 	return inst, nil
// }

// func GetInstallationToken(installationID int64) (string, string, int, error) {
// 	var (
// 		token        string
// 		tokenExpires string
// 		appID        int
// 	)
// 	err := db.QueryRow(`
// 		SELECT i.access_token, i.token_expires_at, a.app_id
// 		FROM github_installations i
// 		JOIN github_app a ON a.id = 1
// 		WHERE i.installation_id = ?
// 	`, installationID).Scan(&token, &tokenExpires, &appID)
// 	if err != nil {
// 		return "", "", 0, err
// 	}
// 	return token, tokenExpires, appID, nil
// }

// func UpdateInstallationToken(installationID int64, token string, newExpiry time.Time) error {
// 	_, err := db.Exec(`
// 		UPDATE github_installations
// 		SET access_token = ?, token_expires_at = ?, updated_at = CURRENT_TIMESTAMP
// 		WHERE installation_id = ?
// 	`, token, newExpiry.Format(time.RFC3339), installationID)
// 	return err
// }

// func GetGithubAppCredentials(appID int) (string, string, error) {
// 	var appNumericId string
// 	var appPrivateKeyPEM string

// 	err := db.QueryRow(`
// 		SELECT app_id, private_key FROM github_app WHERE id = 1
// 	`).Scan(&appNumericId, &appPrivateKeyPEM)

// 	if err != nil {
// 		return "", "", err
// 	}

// 	return appNumericId, appPrivateKeyPEM, nil
// }

// func GetInstallationIDByUserID(userID int64) (int64, error) {
// 	var installationID int64
// 	err := db.QueryRow(`SELECT installation_id FROM github_installations WHERE user_id = ?`, userID).
// 		Scan(&installationID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return installationID, nil
// }

// func GetGithubAppIDAndPrivateKey() (int64, string, error) {
// 	var appAppID int64
// 	var privateKey string
// 	err := db.QueryRow(`SELECT app_id, private_key FROM github_app LIMIT 1`).Scan(&appAppID, &privateKey)
// 	if err != nil {
// 		return 0, "", err
// 	}
// 	return appAppID, privateKey, nil
// }
