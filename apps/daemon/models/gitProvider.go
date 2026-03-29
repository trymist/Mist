package models

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type GitProviderType string

const (
	GitProviderGitHub    GitProviderType = "github"
	GitProviderGitLab    GitProviderType = "gitlab"
	GitProviderBitbucket GitProviderType = "bitbucket"
	GitProviderGitea     GitProviderType = "gitea"
)

type GitProvider struct {
	ID           int64           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID       int64           `gorm:"uniqueIndex:idx_user_provider;index;not null;constraint:OnDelete:CASCADE" json:"user_id"`
	Provider     GitProviderType `gorm:"uniqueIndex:idx_user_provider;not null" json:"provider"`
	AccessToken  string          `gorm:"not null" json:"access_token"`
	RefreshToken *string         `json:"refresh_token,omitempty"`
	ExpiresAt    *time.Time      `json:"expires_at,omitempty"`
	Username     *string         `json:"username,omitempty"`
	Email        *string         `json:"email,omitempty"`
}

func (gp *GitProvider) InsertInDB() error {
	return db.Create(gp).Error
}

func GetGitProviderByID(id int64) (*GitProvider, error) {
	var gp GitProvider
	err := db.First(&gp, id).Error
	if err != nil {
		return nil, err
	}
	return &gp, nil
}

func GetGitProviderByUserAndProvider(userID int64, provider GitProviderType) (*GitProvider, error) {
	var gp GitProvider
	err := db.Where("user_id = ? AND provider = ?", userID, provider).First(&gp).Error
	if err != nil {
		return nil, err
	}
	return &gp, nil
}

func GetGitProvidersByUser(userID int64) ([]GitProvider, error) {
	var providers []GitProvider
	err := db.Where("user_id = ?", userID).Find(&providers).Error
	return providers, err
}

func (gp *GitProvider) UpdateToken(accessToken string, refreshToken *string, expiresAt *time.Time) error {
	return db.Model(gp).Updates(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    expiresAt,
	}).Error
}

func DeleteGitProvider(id int64) error {
	return db.Delete(&GitProvider{}, id).Error
}

func GetGitProviderAccessToken(providerID int64) (string, GitProviderType, bool, error) {
	var gp GitProvider
	err := db.Select("access_token, provider, expires_at").First(&gp, providerID).Error
	if err != nil {
		return "", "", false, err
	}

	needsRefresh := false
	if gp.ExpiresAt != nil && time.Now().After(*gp.ExpiresAt) {
		needsRefresh = true
	}

	return gp.AccessToken, gp.Provider, needsRefresh, nil
}

func GetAppGitInfo(appID int64) (*int64, *string, string, *string, int64, string, error) {
	var app App
	err := db.Select("git_provider_id, git_repository, git_branch, git_clone_url, project_id, name").
		First(&app, appID).Error

	if err != nil {
		return nil, nil, "", nil, 0, "", err
	}

	gitBranch := app.GitBranch
	if gitBranch == "" {
		gitBranch = "main"
	}

	return app.GitProviderID, app.GitRepository, gitBranch, app.GitCloneURL, app.ProjectID, app.Name, nil
}

func UpdateAppGitCloneURL(appID int64, gitCloneURL string, gitProviderID *int64) error {
	return db.Model(&App{ID: appID}).Updates(map[string]interface{}{
		"git_clone_url":   gitCloneURL,
		"git_provider_id": gitProviderID,
	}).Error
}

func RefreshGitProviderToken(providerID int64) (string, error) {
	provider, err := GetGitProviderByID(providerID)
	if err != nil {
		return "", err
	}

	switch provider.Provider {
	case GitProviderGitHub:
		return refreshGitHubToken(provider.UserID)
	case GitProviderGitLab, GitProviderBitbucket, GitProviderGitea:
		return "", fmt.Errorf("token refresh not implemented for %s", provider.Provider)
	default:
		return "", fmt.Errorf("unknown provider type: %s", provider.Provider)
	}
}

func refreshGitHubToken(userID int64) (string, error) {
	installation, err := GetInstallationByUserID(int(userID))
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub installation: %w", err)
	}

	appID, privateKey, err := GetGithubAppIDAndPrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub App credentials: %w", err)
	}

	jwtToken, err := generateGithubJwt(appID, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate GitHub JWT: %w", err)
	}

	newToken, newExpiry, err := regenerateGithubInstallationToken(jwtToken, installation.InstallationID)
	if err != nil {
		return "", fmt.Errorf("failed to regenerate installation token: %w", err)
	}

	err = UpdateInstallationToken(installation.InstallationID, newToken, newExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to update installation token: %w", err)
	}

	var gp GitProvider
	if err := db.Where("user_id = ? AND provider = ?", userID, GitProviderGitHub).First(&gp).Error; err == nil {
		_ = gp.UpdateToken(newToken, nil, &newExpiry)
	}

	return newToken, nil
}

func generateGithubJwt(appID int64, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": fmt.Sprintf("%d", appID),
	})

	return token.SignedString(privateKey)
}

func regenerateGithubInstallationToken(appJWT string, installationID int64) (string, time.Time, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", time.Time{}, fmt.Errorf("failed to create token, status %d", resp.StatusCode)
	}

	var tokenResp struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, err
	}

	return tokenResp.Token, tokenResp.ExpiresAt, nil
}

//########################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"encoding/pem"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// )

// type GitProviderType string

// const (
// 	GitProviderGitHub    GitProviderType = "github"
// 	GitProviderGitLab    GitProviderType = "gitlab"
// 	GitProviderBitbucket GitProviderType = "bitbucket"
// 	GitProviderGitea     GitProviderType = "gitea"
// )

// type GitProvider struct {
// 	ID           int64           `db:"id" json:"id"`
// 	UserID       int64           `db:"user_id" json:"user_id"`
// 	Provider     GitProviderType `db:"provider" json:"provider"`
// 	AccessToken  string          `db:"access_token" json:"access_token"`
// 	RefreshToken *string         `db:"refresh_token" json:"refresh_token,omitempty"`
// 	ExpiresAt    *time.Time      `db:"expires_at" json:"expires_at,omitempty"`
// 	Username     *string         `db:"username" json:"username,omitempty"`
// 	Email        *string         `db:"email" json:"email,omitempty"`
// }

// func (gp *GitProvider) InsertInDB() error {
// 	query := `
// 		INSERT INTO git_providers (user_id, provider, access_token, refresh_token, expires_at, username, email)
// 		VALUES (?, ?, ?, ?, ?, ?, ?)
// 		RETURNING id
// 	`
// 	err := db.QueryRow(query, gp.UserID, gp.Provider, gp.AccessToken, gp.RefreshToken, gp.ExpiresAt, gp.Username, gp.Email).Scan(&gp.ID)
// 	return err
// }

// func GetGitProviderByID(id int64) (*GitProvider, error) {
// 	var gp GitProvider
// 	query := `
// 		SELECT id, user_id, provider, access_token, refresh_token, expires_at, username, email
// 		FROM git_providers
// 		WHERE id = ?
// 	`
// 	err := db.QueryRow(query, id).Scan(
// 		&gp.ID, &gp.UserID, &gp.Provider, &gp.AccessToken, &gp.RefreshToken, &gp.ExpiresAt, &gp.Username, &gp.Email,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &gp, nil
// }

// func GetGitProviderByUserAndProvider(userID int64, provider GitProviderType) (*GitProvider, error) {
// 	var gp GitProvider
// 	query := `
// 		SELECT id, user_id, provider, access_token, refresh_token, expires_at, username, email
// 		FROM git_providers
// 		WHERE user_id = ? AND provider = ?
// 	`
// 	err := db.QueryRow(query, userID, provider).Scan(
// 		&gp.ID, &gp.UserID, &gp.Provider, &gp.AccessToken, &gp.RefreshToken, &gp.ExpiresAt, &gp.Username, &gp.Email,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &gp, nil
// }

// func GetGitProvidersByUser(userID int64) ([]GitProvider, error) {
// 	var providers []GitProvider
// 	query := `
// 		SELECT id, user_id, provider, access_token, refresh_token, expires_at, username, email
// 		FROM git_providers
// 		WHERE user_id = ?
// 	`
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var gp GitProvider
// 		err := rows.Scan(
// 			&gp.ID, &gp.UserID, &gp.Provider, &gp.AccessToken, &gp.RefreshToken, &gp.ExpiresAt, &gp.Username, &gp.Email,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		providers = append(providers, gp)
// 	}
// 	return providers, rows.Err()
// }

// func (gp *GitProvider) UpdateToken(accessToken string, refreshToken *string, expiresAt *time.Time) error {
// 	query := `
// 		UPDATE git_providers
// 		SET access_token = ?, refresh_token = ?, expires_at = ?
// 		WHERE id = ?
// 	`
// 	_, err := db.Exec(query, accessToken, refreshToken, expiresAt, gp.ID)
// 	return err
// }

// func DeleteGitProvider(id int64) error {
// 	query := `DELETE FROM git_providers WHERE id = ?`
// 	_, err := db.Exec(query, id)
// 	return err
// }

// func GetGitProviderAccessToken(providerID int64) (string, GitProviderType, bool, error) {
// 	var token string
// 	var provider GitProviderType
// 	var expiresAt sql.NullTime
// 	query := `SELECT access_token, provider, expires_at FROM git_providers WHERE id = ?`
// 	err := db.QueryRow(query, providerID).Scan(&token, &provider, &expiresAt)
// 	if err != nil {
// 		return "", "", false, err
// 	}

// 	// check if token is expired
// 	needsRefresh := false
// 	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
// 		needsRefresh = true
// 	}

// 	return token, provider, needsRefresh, nil
// }

// func GetAppGitInfo(appID int64) (*int64, *string, string, *string, int64, string, error) {
// 	var gitProviderID sql.NullInt64
// 	var gitRepository sql.NullString
// 	var gitBranch string
// 	var gitCloneURL sql.NullString
// 	var projectID int64
// 	var appName string

// 	query := `
// 		SELECT git_provider_id, git_repository, COALESCE(git_branch, 'main'), git_clone_url, project_id, name
// 		FROM apps
// 		WHERE id = ?
// 	`
// 	err := db.QueryRow(query, appID).Scan(&gitProviderID, &gitRepository, &gitBranch, &gitCloneURL, &projectID, &appName)
// 	if err != nil {
// 		return nil, nil, "", nil, 0, "", err
// 	}

// 	var gitProviderIDPtr *int64
// 	if gitProviderID.Valid {
// 		gitProviderIDPtr = &gitProviderID.Int64
// 	}

// 	var gitRepositoryPtr *string
// 	if gitRepository.Valid {
// 		gitRepositoryPtr = &gitRepository.String
// 	}

// 	var gitCloneURLPtr *string
// 	if gitCloneURL.Valid {
// 		gitCloneURLPtr = &gitCloneURL.String
// 	}

// 	return gitProviderIDPtr, gitRepositoryPtr, gitBranch, gitCloneURLPtr, projectID, appName, nil
// }

// // this is used for migrating old apps that only have git_repository set
// func UpdateAppGitCloneURL(appID int64, gitCloneURL string, gitProviderID *int64) error {
// 	query := `
// 		UPDATE apps
// 		SET git_clone_url = ?, git_provider_id = ?
// 		WHERE id = ?
// 	`
// 	_, err := db.Exec(query, gitCloneURL, gitProviderID, appID)
// 	return err
// }

// // currently only supports GitHub via GitHub App installations
// func RefreshGitProviderToken(providerID int64) (string, error) {
// 	provider, err := GetGitProviderByID(providerID)
// 	if err != nil {
// 		return "", err
// 	}

// 	switch provider.Provider {
// 	case GitProviderGitHub:
// 		return refreshGitHubToken(provider.UserID)
// 	// will be implemented as we add more git-providers
// 	case GitProviderGitLab, GitProviderBitbucket, GitProviderGitea:
// 		return "", fmt.Errorf("token refresh not implemented for %s", provider.Provider)
// 	default:
// 		return "", fmt.Errorf("unknown provider type: %s", provider.Provider)
// 	}
// }

// func refreshGitHubToken(userID int64) (string, error) {
// 	installation, err := GetInstallationByUserID(int(userID))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get GitHub installation: %w", err)
// 	}

// 	appID, privateKey, err := GetGithubAppIDAndPrivateKey()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get GitHub App credentials: %w", err)
// 	}

// 	jwt, err := generateGithubJwt(appID, privateKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate GitHub JWT: %w", err)
// 	}

// 	newToken, newExpiry, err := regenerateGithubInstallationToken(jwt, installation.InstallationID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to regenerate installation token: %w", err)
// 	}

// 	err = UpdateInstallationToken(installation.InstallationID, newToken, newExpiry)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to update installation token: %w", err)
// 	}

// 	// update the git_provider token
// 	err = (&GitProvider{ID: installation.InstallationID}).UpdateToken(newToken, nil, &newExpiry)
// 	if err != nil {
// 		return newToken, nil
// 	}

// 	return newToken, nil
// }

// func generateGithubJwt(appID int64, privateKeyPEM string) (string, error) {
// 	block, _ := pem.Decode([]byte(privateKeyPEM))
// 	if block == nil {
// 		return "", fmt.Errorf("failed to decode PEM block")
// 	}

// 	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
// 	if err != nil {
// 		return "", err
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
// 		"iat": time.Now().Unix(),
// 		"exp": time.Now().Add(10 * time.Minute).Unix(),
// 		"iss": fmt.Sprintf("%d", appID),
// 	})

// 	signedToken, err := token.SignedString(privateKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return signedToken, nil
// }

// func regenerateGithubInstallationToken(appJWT string, installationID int64) (string, time.Time, error) {
// 	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

// 	req, err := http.NewRequest("POST", url, nil)
// 	if err != nil {
// 		return "", time.Time{}, err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+appJWT)
// 	req.Header.Set("Accept", "application/vnd.github+json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", time.Time{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		return "", time.Time{}, fmt.Errorf("failed to create token, status %d", resp.StatusCode)
// 	}

// 	var tokenResp struct {
// 		Token     string    `json:"token"`
// 		ExpiresAt time.Time `json:"expires_at"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
// 		return "", time.Time{}, err
// 	}

// 	return tokenResp.Token, tokenResp.ExpiresAt, nil
// }
