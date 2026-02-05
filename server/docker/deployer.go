package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

func ExecuteContainerDeployment(ctx context.Context, dep *models.Deployment, app *models.App, appContextPath, imageTag, containerName string, db *gorm.DB, logfile *os.File, logger *utils.DeploymentLogger) error {

	logger.Info("Starting container deployment process")

	logger.Info("Fetching deployment configuration with environment variables")
	port, domains, envSet, err := FetchFullDeploymentConfiguration(dep.ID, app, db)
	if err != nil {
		logger.Error(err, "Failed to fetch deployment configuration")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to fetch deployment config: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("fetch deployment config failed: %w", err)
	}

	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"domains":          domains,
		"port":             port,
		"totalEnvVars":     envSet.GetTotalEnvVarCount(),
		"buildTimeEnvVars": envSet.GetBuildTimeCount(),
		"runtimeEnvVars":   envSet.GetRuntimeCount(),
		"appType":          app.AppType,
	})

	if app.AppType == models.AppTypeDatabase {
		logger.Info("Database app detected - pulling Docker image instead of building")

		if app.TemplateName == nil || *app.TemplateName == "" {
			logger.Error(nil, "Database app missing template name")
			dep.Status = "failed"
			dep.Stage = "failed"
			dep.Progress = 0
			errMsg := "Database app requires a template name"
			dep.ErrorMessage = &errMsg
			UpdateDeploymentRecord(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			return fmt.Errorf("database app missing template")
		}

		template, err := models.GetServiceTemplateByName(*app.TemplateName)
		if err != nil {
			logger.Error(err, "Failed to get service template")
			dep.Status = "failed"
			dep.Stage = "failed"
			dep.Progress = 0
			errMsg := fmt.Sprintf("Failed to get template: %v", err)
			dep.ErrorMessage = &errMsg
			UpdateDeploymentRecord(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			return fmt.Errorf("get template failed: %w", err)
		}

		if template == nil {
			logger.Error(nil, "Template not found")
			dep.Status = "failed"
			dep.Stage = "failed"
			dep.Progress = 0
			errMsg := fmt.Sprintf("Template not found: %s", *app.TemplateName)
			dep.ErrorMessage = &errMsg
			UpdateDeploymentRecord(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			return fmt.Errorf("template not found")
		}

		dep.Status = "building"
		dep.Stage = "pulling"
		dep.Progress = 50
		UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "building", "pulling", 50, nil)

		imageName := template.DockerImage
		if template.DockerImageVersion != nil && *template.DockerImageVersion != "" {
			imageName = imageName + ":" + *template.DockerImageVersion
		}

		logger.InfoWithFields("Pulling prebuilt Docker image", map[string]interface{}{
			"image": imageName,
		})

		if err := PullPrebuiltDockerImage(ctx, imageName, logfile); err != nil {
			if ctx.Err() == context.Canceled {
				logger.Info("Docker image pull canceled")
				return ctx.Err()
			}
			logger.Error(err, "Docker image pull failed")
			dep.Status = "failed"
			dep.Stage = "failed"
			dep.Progress = 0
			errMsg := fmt.Sprintf("Pull failed: %v", err)
			dep.ErrorMessage = &errMsg
			UpdateDeploymentRecord(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			UpdateApplicationStatus(app.ID, "error", db)
			return fmt.Errorf("pull image failed: %w", err)
		}

		logger.Info("Docker image pulled successfully")
		imageTag = imageName

	} else {
		dep.Status = "building"
		dep.Stage = "building"
		dep.Progress = 50
		UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "building", "building", 50, nil)

		logger.InfoWithFields("Building Docker image with build-time arguments", map[string]interface{}{
			"buildArgsCount": envSet.GetBuildTimeCount(),
		})
		if err := BuildDockerImageWithBuildArgs(ctx, imageTag, appContextPath, envSet.BuildTime, logfile); err != nil {
			if ctx.Err() == context.Canceled {
				logger.Info("Docker image build canceled")
				return ctx.Err()
			}
			logger.Error(err, "Docker image build failed")
			dep.Status = "failed"
			dep.Stage = "failed"
			dep.Progress = 0
			errMsg := fmt.Sprintf("Build failed: %v", err)
			dep.ErrorMessage = &errMsg
			UpdateDeploymentRecord(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			UpdateApplicationStatus(app.ID, "error", db)
			return fmt.Errorf("build image failed: %w", err)
		}

		logger.Info("Docker image built successfully")
	}

	dep.Status = "deploying"
	dep.Stage = "deploying"
	dep.Progress = 80
	UpdateDeploymentRecord(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 80, nil)

	logger.Info("Stopping existing container if exists")
	err = StopAndRemoveContainer(containerName, logfile)
	if err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Container stop/remove canceled")
			return ctx.Err()
		}
		logger.Error(err, "Failed to stop/remove existing container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to stop/remove container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateApplicationStatus(app.ID, "error", db)
		return fmt.Errorf("stop/remove container failed: %w", err)
	}

	logger.InfoWithFields("Creating and starting container", map[string]interface{}{
		"domains":        domains,
		"port":           port,
		"runtimeEnvVars": envSet.GetRuntimeCount(),
		"appType":        app.AppType,
	})

	if err := CreateAndStartContainer(ctx, app, imageTag, containerName, domains, port, envSet.Runtime, logfile); err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Container creation canceled")
			return ctx.Err()
		}
		logger.Error(err, "Failed to create and start container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to create and start container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeploymentRecord(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateApplicationStatus(app.ID, "error", db)
		return fmt.Errorf("create and start container failed: %w", err)
	}

	dep.Status = "success"
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	UpdateDeploymentRecord(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.Info("Updating application status to running")
	err = UpdateApplicationStatus(app.ID, "running", db)
	if err != nil {
		logger.Error(err, "Failed to update application status (non-fatal)")
	}

	logger.Info("Cleaning up old Docker images")
	if err := CleanupOldImages(app.ID, 5); err != nil {
		logger.Error(err, "Failed to cleanup old images (non-fatal)")
	}

	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
		"deployment_id": dep.ID,
		"container":     containerName,
		"app_status":    "running",
	})

	return nil
}

func UpdateDeploymentRecord(dep *models.Deployment, db *gorm.DB) error {
	return db.Model(dep).Updates(map[string]interface{}{
		"status":        dep.Status,
		"stage":         dep.Stage,
		"progress":      dep.Progress,
		"logs":          dep.Logs,
		"error_message": dep.ErrorMessage,
		"finished_at":   dep.FinishedAt,
	}).Error
}

func GetBuildLogsPath(commitHash string, depId int64) string {
	return filepath.Join(constants.Constants["LogPath"].(string), commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
}

func UpdateApplicationStatus(appID int64, status string, db *gorm.DB) error {
	return db.Model(&models.App{ID: appID}).Update("status", status).Error
}

func FetchFullDeploymentConfiguration(deploymentID int64, app *models.App, db *gorm.DB) (int, []string, *EnvironmentVariableSet, error) {
	appID, err := models.GetAppIDByDeploymentID(deploymentID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get app ID failed: %w", err)
	}

	var port *int
	err = db.Model(&models.App{}).Select("port").Where("id = ?", appID).Scan(&port).Error
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get port failed: %w", err)
	}

	finalPort := 3000
	if port != nil {
		finalPort = *port
	}

	domains, err := models.GetDomainsByAppID(appID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get domains failed: %w", err)
	}

	var domainStrings []string
	for _, d := range domains {
		domainStrings = append(domainStrings, d.Domain)
	}

	envs, err := models.GetEnvVariablesByAppID(appID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get env variables failed: %w", err)
	}

	envSet := CategorizeEnvironmentVariables(envs)

	if app.AppType == models.AppTypeDatabase && app.TemplateName != nil {
		template, err := models.GetServiceTemplateByName(*app.TemplateName)
		if err == nil && template != nil && template.DefaultEnvVars != nil {
			var defaultEnvs map[string]string
			if err := json.Unmarshal([]byte(*template.DefaultEnvVars), &defaultEnvs); err == nil {
				for k, v := range defaultEnvs {
					envSet.Runtime[k] = v
				}
			}
		}
	}

	return finalPort, domainStrings, envSet, nil
}

func GetDeploymentConfig(deploymentID int64, app *models.App, db *gorm.DB) (int, []string, map[string]string, error) {
	port, domains, envSet, err := FetchFullDeploymentConfiguration(deploymentID, app, db)
	if err != nil {
		return 0, nil, nil, err
	}

	merged := make(map[string]string)
	for k, v := range envSet.BuildTime {
		merged[k] = v
	}
	for k, v := range envSet.Runtime {
		merged[k] = v
	}

	return port, domains, merged, nil
}
