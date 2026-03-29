package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/rs/zerolog/log"
)

type WebhookPayload struct {
	Event        string              `json:"-"` // X-GitHub-Event
	Raw          json.RawMessage     `json:"-"` // full body
	Repository   *github.RepoFull    `json:"repository,omitempty"`
	Installation *github.InstallMini `json:"installation,omitempty"`
	Sender       *github.User        `json:"sender,omitempty"`
}

func verifyGitHubSignature(payload []byte, signature string, secret string) bool {
	if secret == "" {
		log.Warn().Msg("GitHub webhook secret not set - webhook signature verification disabled (security risk)")
		return true
	}

	if signature == "" {
		return false
	}

	if len(signature) < 7 || signature[:7] != "sha256=" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	receivedMAC := signature[7:]

	return hmac.Equal([]byte(expectedMAC), []byte(receivedMAC))
}

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Received GitHub webhook")

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		http.Error(w, "Missing X-GitHub-Event header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	// Get webhook secret from database
	settings, err := models.GetSystemSettings()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get system settings for webhook verification")
		http.Error(w, "Configuration error", http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifyGitHubSignature(body, signature, settings.GithubWebhookSecret) {
		log.Warn().Str("event", eventType).Msg("Invalid webhook signature")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if eventType == "push" {
		var evt github.PushEvent
		if err := json.Unmarshal(body, &evt); err != nil {
			http.Error(w, "Invalid push event payload", http.StatusBadRequest)
			return
		}

		log.Info().Str("repo", evt.Repository.FullName).Msg("Processing push event")
		depId, err := github.CreateDeploymentFromGithubPushEvent(evt)
		if err != nil {
			log.Error().Err(err).Str("repo", evt.Repository.FullName).Msg("Failed to create deployment from push event")
			http.Error(w, "Failed to handle push event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if depId != 0 {
			queue := queue.GetQueue()
			queue.AddJob(depId)
			log.Info().Int64("deployment_id", depId).Msg("Deployment queued")
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}
