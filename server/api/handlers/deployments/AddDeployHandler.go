package deployments

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/git"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/rs/zerolog/log"
)

func AddDeployHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AppId int `json:"appId"`
	}
	queue := queue.GetQueue()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid request body", err.Error())
		return
	}

	user, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "unauthorized", "")
		return
	}

	app, err := models.GetApplicationByID(int64(req.AppId))
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to get app details", err.Error())
		return
	}

	var commitHash string
	var commitMessage string

	if app.AppType != models.AppTypeDatabase {
		userId := int64(user.ID)
		commit, err := git.GetLatestCommit(int64(req.AppId), userId)
		if err != nil {
			log.Error().Err(err).Int("app_id", req.AppId).Msg("Error getting latest commit")
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to get latest commit", err.Error())
			return
		}
		commitHash = commit.SHA
		commitMessage = commit.Message
	} else {
		if app.TemplateName != nil {
			template, err := models.GetServiceTemplateByName(*app.TemplateName)
			if err == nil && template != nil && template.DockerImageVersion != nil {
				commitHash = *template.DockerImageVersion
			} else {
				commitHash = "latest"
			}
		} else {
			commitHash = "latest"
		}
		commitMessage = "Deploy database service"
	}

	deployment := models.Deployment{
		AppID:         int64(req.AppId),
		CommitHash:    commitHash,
		CommitMessage: &commitMessage,
		Status:        models.DeploymentStatusPending,
	}
	err = deployment.CreateDeployment()

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to insert deployment", err.Error())
		return
	}

	// create github deployment
	if app.GitRepository != nil {
		depId, err := github.CreateDeployment(*app.GitRepository, app.GitBranch, int(user.ID))
		if err != nil {
			log.Err(err).Msg("failed to create github deployment")
		}

		deployment.GithubDepId = &depId
		err = deployment.UpdateDeployment()
		if err != nil {
			log.Err(err).Msg("failed to update deployment with GH dep id")
		}
	}
	if err := queue.AddJob(int64(deployment.ID)); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to add job to queue", err.Error())
		return
	}

	log.Info().Int64("deployment_id", deployment.ID).Int("app_id", req.AppId).Msg("Deployment added to queue")

	models.LogUserAudit(user.ID, "create", "deployment", &deployment.ID, map[string]interface{}{
		"app_id":         deployment.AppID,
		"commit_hash":    deployment.CommitHash,
		"commit_message": deployment.CommitMessage,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)

}
