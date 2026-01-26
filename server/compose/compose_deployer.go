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

	logger.Info("Getting port, domains, and environment variables")
	port, domains, envVars, err := docker.GetDeploymentConfig(dep.ID, app, db)
	if err != nil {
		logger.Error(err, "Failed to get deployment configuration")
		dep.Status = models.DeploymentStatusFailed
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to get deployment config: %v", err)
		dep.ErrorMessage = &errMsg
		docker.UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("get deployment config failed: %w", err)
	}

	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"domains": domains,
		"port":    port,
		"envVars": len(envVars),
		"appType": app.AppType,
	})

	dep.Status = models.DeploymentStatusDeploying
	dep.Stage = "deploying"
	dep.Progress = 50
	docker.UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 50, nil)

	logger.Info("Running docker compose up")

	err = ComposeUp(appContextPath, envVars, logfile)
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
		docker.UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		docker.UpdateAppStatus(app.ID, "error", db)
		return fmt.Errorf("compose up failed: %w", err)
	}

	logger.Info("Docker compose up completed successfully")

	dep.Status = models.DeploymentStatusSuccess
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	docker.UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.Info("Updating app status to running")
	err = docker.UpdateAppStatus(app.ID, "running", db)
	if err != nil {
		logger.Error(err, "Failed to update app status (non-fatal)")
	}

	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
		"deployment_id": dep.ID,
		"app_status":    "running",
	})

	return nil
}
