package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/fs"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

// to prevent concurrent deployments of same app
// two deployments of same app shouldn't be happening at the same time to prevent race conditions
// currently its impossible, bcz we only support one deployment at a time,
// but future plans include configuratble concurrent deployments
// then this will be helpful
var deploymentLocks sync.Map

func (q *Queue) HandleWork(id int64, db *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("panic during deployment: %v", r)
			models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Register(id, cancel)
	defer Unregister(id)

	appId, err := models.GetAppIDByDeploymentID(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get app ID: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	if _, loaded := deploymentLocks.LoadOrStore(appId, true); loaded {
		errMsg := fmt.Sprintf("Deployment already in progress for app %d", appId)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}
	// make sure to release the lock, once deployment is complete, else we'll not be able to deploy that app again
	defer deploymentLocks.Delete(appId)

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get app details: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	dep, err := docker.LoadDeployment(id, db)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load deployment: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logger := utils.NewDeploymentLogger(id, appId, dep.CommitHash)
	logger.Info("Starting deployment processing")

	if err := models.MarkDeploymentStarted(id); err != nil {
		logger.Error(err, "Failed to mark deployment as started")
		errMsg := fmt.Sprintf("Failed to update deployment start time: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logFile, _, err := fs.CreateDockerBuildLogFile(id)
	if err != nil {
		logger.Error(err, "Failed to create log file")
		errMsg := fmt.Sprintf("Failed to create log file: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}
	defer logFile.Close()

	if app.AppType != models.AppTypeDatabase {
		logger.Info("Cloning repository")
		models.UpdateDeploymentStatus(id, "cloning", "cloning", 20, nil)

		err = github.CloneRepo(ctx, appId, logFile)
		if err != nil {
			if ctx.Err() == context.Canceled {
				logger.Info("Deployment cancelled by user")
				errMsg := "deployment stopped by user"
				models.UpdateDeploymentStatus(id, "stopped", "stopped", dep.Progress, &errMsg)
				return
			}
			logger.Error(err, "Failed to clone repository")
			errMsg := fmt.Sprintf("Failed to clone repository: %v", err)
			models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
			return
		}

		logger.Info("Repository cloned successfully")
	} else {
		logger.Info("Skipping git clone for database app")
	}

	_, err = docker.DeployerMain(ctx, id, db, logFile, logger)
	if err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Deployment cancelled by user")
			errMsg := "deployment stopped by user"
			models.UpdateDeploymentStatus(id, "stopped", "stopped", dep.Progress, &errMsg)
			return
		}
		logger.Error(err, "Deployment failed")
		errMsg := fmt.Sprintf("Deployment failed: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logger.Info("Deployment completed successfully")
}
