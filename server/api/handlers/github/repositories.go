package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

type RepoListResponse struct {
	TotalCount   int   `json:"total_count"`
	Repositories []any `json:"repositories"`
}

func GetRepositories(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	installationID, err := models.GetInstallationID(int(userData.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "no installation found for user", "No installation found")
		return
	}
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "database error", err.Error())
		return
	}

	token, tokenExpires, appID, err := models.GetInstallationToken(installationID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to fetch installation info", err.Error())
		return
	}

	expiry, _ := time.Parse(time.RFC3339, tokenExpires)
	if time.Now().After(expiry) {
		appJWT, err := github.GenerateGithubJwt(appID)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to generate app jwt", err.Error())
			return
		}

		newToken, newExpiry, err := regenerateInstallationToken(appJWT, installationID)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to refresh token", err.Error())
			return
		}

		_ = models.UpdateInstallationToken(int64(installationID), newToken, newExpiry)

		token = newToken
	}

	allRepos := []any{}
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/installation/repositories?per_page=100&page=%d", page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "request error", err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			handlers.SendResponse(w, resp.StatusCode, false, nil, fmt.Sprintf("GitHub API returned %d", resp.StatusCode), "")
			return
		}

		var repoList RepoListResponse
		if err := json.NewDecoder(resp.Body).Decode(&repoList); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to parse GitHub response", err.Error())
			return
		}

		allRepos = append(allRepos, repoList.Repositories...)

		if len(repoList.Repositories) < 100 {
			break // no more pages
		}
		page++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allRepos)
}

func regenerateInstallationToken(appJWT string, installationID int64) (string, time.Time, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	req, _ := http.NewRequest("POST", url, nil)
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
