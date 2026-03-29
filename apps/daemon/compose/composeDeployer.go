// exec.Command is still used for docker compose, because moby doesn't support compose yet

package compose

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

func DeployComposeApp(ctx context.Context, dep *models.Deployment, app *models.App, appContextPath string, db *gorm.DB, logfile *os.File, logger *utils.DeploymentLogger) error {
	logger.Info("Starting compose deployment process")

	logger.Info("Fetching deployment configuration with environment variables")
	port, domains, envSet, err := docker.FetchFullDeploymentConfiguration(dep.ID, app, db)
	if err != nil {
		logger.Error(err, "Failed to fetch deployment configuration")
		dep.Status = models.DeploymentStatusFailed
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to fetch deployment config: %v", err)
		dep.ErrorMessage = &errMsg
		docker.UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("fetch deployment config failed: %w", err)
	}

	mergedEnvVars := make(map[string]string)
	for k, v := range envSet.BuildTime {
		mergedEnvVars[k] = v
	}
	for k, v := range envSet.Runtime {
		mergedEnvVars[k] = v
	}

	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"domains":          domains,
		"port":             port,
		"totalEnvVars":     envSet.GetTotalEnvVarCount(),
		"buildTimeEnvVars": envSet.GetBuildTimeCount(),
		"runtimeEnvVars":   envSet.GetRuntimeCount(),
		"appType":          app.AppType,
	})

	dep.Status = models.DeploymentStatusDeploying
	dep.Stage = "deploying"
	dep.Progress = 50
	docker.UpdateDeploymentRecord(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 50, nil)

	logger.Info("Running docker compose up")

	err = ComposeUp(appContextPath, mergedEnvVars, logfile)
	if err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Compose deployment canceled")
			return ctx.Err()
		}
		logger.Error(err, "Docker compose up failed")
		dep.Status = models.DeploymentStatusFailed
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Compose up failed: %v", err)
		dep.ErrorMessage = &errMsg
		docker.UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		docker.UpdateApplicationStatus(app.ID, "error", db)
		return fmt.Errorf("compose up failed: %w", err)
	}

	logger.Info("Docker compose up completed successfully")

	dep.Status = models.DeploymentStatusSuccess
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	docker.UpdateDeploymentRecord(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.Info("Updating application status to running")
	err = docker.UpdateApplicationStatus(app.ID, "running", db)
	if err != nil {
		logger.Error(err, "Failed to update application status (non-fatal)")
	}

	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
		"deployment_id": dep.ID,
		"app_status":    "running",
	})

	return nil
}
