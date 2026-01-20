package github

import (
	"errors"
	"strings"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateDeploymentFromGithubPushEvent(evt PushEvent) (int64, error) {
	repoName := evt.Repository.FullName
	branch := evt.Ref
	commit := evt.After

	branch = strings.TrimPrefix(branch, "refs/heads/")

	log.Info().
		Str("repo", repoName).
		Str("branch", branch).
		Str("commit", commit).
		Msg("Push event received")

	appID, err := models.FindApplicationIDByGitRepoAndBranch(repoName, branch)

	if err != nil {
		log.Error().Err(err).
			Str("repo", repoName).
			Str("branch", branch).
			Msg("Error finding application")
		return 0, err
	}

	if appID == 0 {
		log.Warn().
			Str("repo", repoName).
			Str("branch", branch).
			Msg("No application found for this repository and branch")
		return 0, errors.New("no application found for this repository and branch")
	}

	app, err := models.GetApplicationByID(appID)
	if err != nil {
		log.Error().Err(err).
			Int64("app_id", appID).
			Msg("Error getting application")
		return 0, err
	}

	if app.DeploymentStrategy == models.DeploymentManual {
		log.Info().
			Int64("app_id", appID).
			Str("repo", repoName).
			Str("branch", branch).
			Msg("Skipping automatic deployment - deployment strategy is set to manual")
		return 0, nil
	}

	dep, err := models.GetDeploymentByAppIDAndCommitHash(appID, commit)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).
			Int64("app_id", appID).
			Str("commit", commit).
			Msg("Error checking for existing deployment")
		return 0, err
	}
	if dep != nil && dep.ID != 0 {
		log.Warn().
			Int64("dep_id", dep.ID).
			Int64("app_id", appID).
			Str("commit", commit).
			Msg("Deployment already exists for this commit, skipping duplicate")
		return 0, nil
	}

	commitMsg := evt.HeadCommit.Message
	deployment := models.Deployment{
		AppID:         appID,
		CommitHash:    commit,
		CommitMessage: &commitMsg,
	}

	if err := deployment.CreateDeployment(); err != nil {
		log.Error().Err(err).
			Int64("app_id", appID).
			Msg("Error creating deployment")
		return 0, err
	}

	log.Info().
		Int64("deployment_id", deployment.ID).
		Int64("app_id", appID).
		Msg("Deployment created from GitHub webhook")

	models.LogWebhookAudit("create", "deployment", &deployment.ID, map[string]interface{}{
		"app_id":         appID,
		"commit_hash":    commit,
		"commit_message": evt.HeadCommit.Message,
		"repository":     repoName,
		"branch":         branch,
		"pusher":         evt.Pusher.Name,
	})

	return deployment.ID, nil

}
