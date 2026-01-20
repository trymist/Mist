package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

func DeployerMain(ctx context.Context, Id int64, db *gorm.DB, logFile *os.File, logger *utils.DeploymentLogger) (string, error) {
	dep, err := LoadDeployment(Id, db)
	if err != nil {
		logger.Error(err, "Failed to load deployment")
		return "", fmt.Errorf("failed to load deployment: %w", err)
	}

	var appId int64
	// err = db.QueryRow("SELECT app_id FROM deployments WHERE id = ?", Id).Scan(&appId)
	err = db.Table("deployments").Select("app_id").Where("id = ?", Id).Take(&appId).Error

	if err != nil {
		logger.Error(err, "Failed to get app_id")
		return "", fmt.Errorf("failed to get app_id: %w", err)
	}

	appPtr, err := models.GetApplicationByID(appId)
	if err != nil {
		logger.Error(err, "Failed to get app details")
		return "", fmt.Errorf("failed to get app details: %w", err)
	}
	app := *appPtr

	logger.InfoWithFields("App details loaded", map[string]interface{}{
		"app_name":   app.Name,
		"project_id": app.ProjectID,
		"app_type":   app.AppType,
	})

	appContextPath := filepath.Join(constants.Constants["RootPath"].(string), fmt.Sprintf("projects/%d/apps/%s", app.ProjectID, app.Name))
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", app.ID)

	err = DeployApp(ctx, dep, &app, appContextPath, imageTag, containerName, db, logFile, logger)
	if err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Deployment cancelled by user")
			return "", context.Canceled
		}
		logger.Error(err, "DeployApp failed")
		dep.Status = "failed"
		dep.Stage = "failed"
		errMsg := err.Error()
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		return "", err
	}

	logger.Info("Deployment completed successfully")

	settings, err := models.GetSystemSettings()
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to get system settings for cleanup: %v", err))
	} else {
		if settings.AutoCleanupContainers {
			logger.Info("Running automatic container cleanup")
			if err := CleanupStoppedContainers(); err != nil {
				logger.Warn(fmt.Sprintf("Container cleanup failed: %v", err))
			} else {
				logger.Info("Container cleanup completed")
			}
		}

		if settings.AutoCleanupImages {
			logger.Info("Running automatic image cleanup")
			if err := CleanupDanglingImages(); err != nil {
				logger.Warn(fmt.Sprintf("Image cleanup failed: %v", err))
			} else {
				logger.Info("Image cleanup completed")
			}
		}
	}

	return "Deployment started", nil
}
