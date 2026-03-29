package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

type GithubAppConversion struct {
	ID            int               `json:"id"`
	Slug          string            `json:"slug"`
	NodeID        string            `json:"node_id"`
	Name          string            `json:"name"`
	ClientID      string            `json:"client_id"`
	ClientSecret  string            `json:"client_secret"`
	WebhookSecret string            `json:"webhook_secret"`
	PEM           string            `json:"pem"`
	HTMLURL       string            `json:"html_url"`
	ExternalURL   string            `json:"external_url"`
	Permissions   map[string]string `json:"permissions"`
	Events        []string          `json:"events"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Owner         struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"owner"`
}

func CallBackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	baseFrontendURL := GetFrontendBaseUrl()

	if code == "" || state == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=missing_code_or_state", baseFrontendURL), http.StatusSeeOther)
		return
	}

	decodedState, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_state_encoding", baseFrontendURL), http.StatusSeeOther)
		return
	}

	var stateData StateData
	if err := json.Unmarshal(decodedState, &stateData); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_state_data", baseFrontendURL), http.StatusSeeOther)
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), nil)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=request_creation_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Mist-App")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=github_request_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=github_api_error&details=%s", baseFrontendURL, base64.URLEncoding.EncodeToString(body)), http.StatusSeeOther)
		return
	}

	var app GithubAppConversion
	if err := json.Unmarshal(body, &app); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=invalid_github_response", baseFrontendURL), http.StatusSeeOther)
		return
	}

	githubApp := models.GithubApp{
		AppID:         int64(app.ID),
		ClientID:      app.ClientID,
		ClientSecret:  app.ClientSecret,
		WebhookSecret: app.WebhookSecret,
		PrivateKey:    app.PEM,
		Name:          &app.Name,
		Slug:          app.Slug,
	}

	err = githubApp.InsertInDB()
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/callback?error=db_insert_failed", baseFrontendURL), http.StatusSeeOther)
		return
	}

	if err := models.SetSystemSetting("github_webhook_secret", app.WebhookSecret); err != nil {
		log.Warn().Err(err).Msg("Failed to save webhook secret to system_settings")
	}

	userID := int64(stateData.UserId)
	models.LogUserAudit(userID, "create", "github_app", &githubApp.AppID, map[string]interface{}{
		"appName": app.Name,
		"appSlug": app.Slug,
		"appId":   app.ID,
		"owner":   app.Owner.Login,
	})

	newState := GenerateState(app.ID, stateData.UserId)
	redirectURL := fmt.Sprintf("https://github.com/apps/%s/installations/new?state=%s", app.Slug, newState)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
