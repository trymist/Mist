package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/corecollectives/mist/models"
)

func GetGitHubAccessToken(userID int) (string, time.Time, error) {
	app, _, err := models.GetApp(userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to fetch github app credentials: %w", err)
	}

	inst, err := models.GetInstallationByUserID(userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to fetch github installation: %w", err)
	}

	if time.Until(inst.TokenExpiresAt) > 5*time.Minute {
		return inst.AccessToken, inst.TokenExpiresAt, nil
	}

	jwt, err := GenerateGithubJwt(int(app.AppID))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create JWT: %w", err)
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", inst.InstallationID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to request installation token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var body bytes.Buffer
		body.ReadFrom(resp.Body)
		return "", time.Time{}, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, body.String())
	}

	var result struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	err = models.UpdateInstallationToken(inst.InstallationID, result.Token, result.ExpiresAt)

	return result.Token, result.ExpiresAt, nil
}
