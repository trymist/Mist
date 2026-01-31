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

func DeployApp(ctx context.Context, dep *models.Deployment, app *models.App, appContextPath, imageTag, containerName string, db *gorm.DB, logfile *os.File, logger *utils.DeploymentLogger) error {

	logger.Info("Starting deployment process")

	logger.Info("Getting port, domains, and environment variables")
	port, domains, envVars, err := GetDeploymentConfig(dep.ID, app, db)
	if err != nil {
		logger.Error(err, "Failed to get deployment configuration")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to get deployment config: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("get deployment config failed: %w", err)
	}

	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"domains": domains,
		"port":    port,
		"envVars": len(envVars),
		"appType": app.AppType,
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
			UpdateDeployment(dep, db)
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
			UpdateDeployment(dep, db)
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
			UpdateDeployment(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			return fmt.Errorf("template not found")
		}

		dep.Status = "building"
		dep.Stage = "pulling"
		dep.Progress = 50
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "building", "pulling", 50, nil)

		imageName := template.DockerImage
		if template.DockerImageVersion != nil && *template.DockerImageVersion != "" {
			imageName = imageName + ":" + *template.DockerImageVersion
		}

		logger.InfoWithFields("Pulling Docker image", map[string]interface{}{
			"image": imageName,
		})

		if err := PullDockerImage(ctx, imageName, logfile); err != nil {
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
			UpdateDeployment(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			UpdateAppStatus(app.ID, "error", db)
			return fmt.Errorf("pull image failed: %w", err)
		}

		logger.Info("Docker image pulled successfully")
		imageTag = imageName

	} else {
		dep.Status = "building"
		dep.Stage = "building"
		dep.Progress = 50
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "building", "building", 50, nil)

		logger.Info("Building Docker image with environment variables")
		if err := BuildImage(ctx, imageTag, appContextPath, envVars, logfile); err != nil {
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
			UpdateDeployment(dep, db)
			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
			UpdateAppStatus(app.ID, "error", db)
			return fmt.Errorf("build image failed: %w", err)
		}

		logger.Info("Docker image built successfully")
	}

	dep.Status = "deploying"
	dep.Stage = "deploying"
	dep.Progress = 80
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 80, nil)

	logger.Info("Stopping existing container if exists")
	err = StopRemoveContainer(containerName, logfile)
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
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateAppStatus(app.ID, "error", db)
		return fmt.Errorf("stop/remove container failed: %w", err)
	}

	logger.InfoWithFields("Running container", map[string]interface{}{
		"domains": domains,
		"port":    port,
		"envVars": len(envVars),
		"appType": app.AppType,
	})

	if err := RunContainer(ctx, app, imageTag, containerName, domains, port, envVars, logfile); err != nil {
		if ctx.Err() == context.Canceled {
			logger.Info("Container run canceled")
			return ctx.Err()
		}
		logger.Error(err, "Failed to run container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to run container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateAppStatus(app.ID, "error", db)
		return fmt.Errorf("run container failed: %w", err)
	}

	dep.Status = "success"
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.Info("Updating app status to running")
	err = UpdateAppStatus(app.ID, "running", db)
	if err != nil {
		logger.Error(err, "Failed to update app status (non-fatal)")
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

func UpdateDeployment(dep *models.Deployment, db *gorm.DB) error {
	return db.Model(dep).Updates(map[string]interface{}{
		"status":        dep.Status,
		"stage":         dep.Stage,
		"progress":      dep.Progress,
		"logs":          dep.Logs,
		"error_message": dep.ErrorMessage,
		"finished_at":   dep.FinishedAt,
	}).Error
}

func GetLogsPath(commitHash string, depId int64) string {
	return filepath.Join(constants.Constants["LogPath"].(string), commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
}

func UpdateAppStatus(appID int64, status string, db *gorm.DB) error {
	return db.Model(&models.App{ID: appID}).Update("status", status).Error
}

func GetDeploymentConfig(deploymentID int64, app *models.App, db *gorm.DB) (int, []string, map[string]string, error) {
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

	envMap := make(map[string]string)

	if app.AppType == models.AppTypeDatabase && app.TemplateName != nil {
		template, err := models.GetServiceTemplateByName(*app.TemplateName)
		if err == nil && template != nil && template.DefaultEnvVars != nil {
			var defaultEnvs map[string]string
			if err := json.Unmarshal([]byte(*template.DefaultEnvVars), &defaultEnvs); err == nil {
				for k, v := range defaultEnvs {
					envMap[k] = v
				}
			}
		}
	}

	for _, env := range envs {
		envMap[env.Key] = env.Value
	}

	return finalPort, domainStrings, envMap, nil
}

// package docker

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"time"

// 	"github.com/corecollectives/mist/constants"
// 	"github.com/corecollectives/mist/models"
// 	"github.com/corecollectives/mist/utils"
// 	"gorm.io/gorm"
// )

// func DeployApp(dep *models.Deployment, app *models.App, appContextPath, imageTag, containerName string, db *gorm.DB, logfile *os.File, logger *utils.DeploymentLogger) error {

// 	logger.Info("Starting deployment process")

// 	logger.Info("Getting port, domains, and environment variables")
// 	port, domains, envVars, err := GetDeploymentConfig(dep.ID, app, db)
// 	if err != nil {
// 		logger.Error(err, "Failed to get deployment configuration")
// 		dep.Status = "failed"
// 		dep.Stage = "failed"
// 		dep.Progress = 0
// 		errMsg := fmt.Sprintf("Failed to get deployment config: %v", err)
// 		dep.ErrorMessage = &errMsg
// 		UpdateDeployment(dep, db)
// 		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 		return fmt.Errorf("get deployment config failed: %w", err)
// 	}

// 	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
// 		"domains": domains,
// 		"port":    port,
// 		"envVars": len(envVars),
// 		"appType": app.AppType,
// 	})

// 	if app.AppType == models.AppTypeDatabase {
// 		logger.Info("Database app detected - pulling Docker image instead of building")

// 		if app.TemplateName == nil || *app.TemplateName == "" {
// 			logger.Error(nil, "Database app missing template name")
// 			dep.Status = "failed"
// 			dep.Stage = "failed"
// 			dep.Progress = 0
// 			errMsg := "Database app requires a template name"
// 			dep.ErrorMessage = &errMsg
// 			UpdateDeployment(dep, db)
// 			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 			return fmt.Errorf("database app missing template")
// 		}

// 		template, err := models.GetServiceTemplateByName(*app.TemplateName)
// 		if err != nil {
// 			logger.Error(err, "Failed to get service template")
// 			dep.Status = "failed"
// 			dep.Stage = "failed"
// 			dep.Progress = 0
// 			errMsg := fmt.Sprintf("Failed to get template: %v", err)
// 			dep.ErrorMessage = &errMsg
// 			UpdateDeployment(dep, db)
// 			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 			return fmt.Errorf("get template failed: %w", err)
// 		}

// 		if template == nil {
// 			logger.Error(nil, "Template not found")
// 			dep.Status = "failed"
// 			dep.Stage = "failed"
// 			dep.Progress = 0
// 			errMsg := fmt.Sprintf("Template not found: %s", *app.TemplateName)
// 			dep.ErrorMessage = &errMsg
// 			UpdateDeployment(dep, db)
// 			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 			return fmt.Errorf("template not found")
// 		}

// 		dep.Status = "building"
// 		dep.Stage = "pulling"
// 		dep.Progress = 50
// 		UpdateDeployment(dep, db)
// 		models.UpdateDeploymentStatus(dep.ID, "building", "pulling", 50, nil)

// 		imageName := template.DockerImage
// 		if template.DockerImageVersion != nil && *template.DockerImageVersion != "" {
// 			imageName = imageName + ":" + *template.DockerImageVersion
// 		}

// 		logger.InfoWithFields("Pulling Docker image", map[string]interface{}{
// 			"image": imageName,
// 		})

// 		if err := PullDockerImage(imageName, logfile); err != nil {
// 			logger.Error(err, "Docker image pull failed")
// 			dep.Status = "failed"
// 			dep.Stage = "failed"
// 			dep.Progress = 0
// 			errMsg := fmt.Sprintf("Pull failed: %v", err)
// 			dep.ErrorMessage = &errMsg
// 			UpdateDeployment(dep, db)
// 			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 			UpdateAppStatus(app.ID, "error", db)
// 			return fmt.Errorf("pull image failed: %w", err)
// 		}

// 		logger.Info("Docker image pulled successfully")
// 		imageTag = imageName

// 	} else {
// 		dep.Status = "building"
// 		dep.Stage = "building"
// 		dep.Progress = 50
// 		UpdateDeployment(dep, db)
// 		models.UpdateDeploymentStatus(dep.ID, "building", "building", 50, nil)

// 		logger.Info("Building Docker image with environment variables")
// 		if err := BuildImage(imageTag, appContextPath, envVars, logfile); err != nil {
// 			logger.Error(err, "Docker image build failed")
// 			dep.Status = "failed"
// 			dep.Stage = "failed"
// 			dep.Progress = 0
// 			errMsg := fmt.Sprintf("Build failed: %v", err)
// 			dep.ErrorMessage = &errMsg
// 			UpdateDeployment(dep, db)
// 			models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 			UpdateAppStatus(app.ID, "error", db)
// 			return fmt.Errorf("build image failed: %w", err)
// 		}

// 		logger.Info("Docker image built successfully")
// 	}

// 	dep.Status = "deploying"
// 	dep.Stage = "deploying"
// 	dep.Progress = 80
// 	UpdateDeployment(dep, db)
// 	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 80, nil)

// 	logger.Info("Stopping existing container if exists")
// 	err = StopRemoveContainer(containerName, logfile)
// 	if err != nil {
// 		logger.Error(err, "Failed to stop/remove existing container")
// 		dep.Status = "failed"
// 		dep.Stage = "failed"
// 		dep.Progress = 0
// 		errMsg := fmt.Sprintf("Failed to stop/remove container: %v", err)
// 		dep.ErrorMessage = &errMsg
// 		UpdateDeployment(dep, db)
// 		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 		UpdateAppStatus(app.ID, "error", db)
// 		return fmt.Errorf("stop/remove container failed: %w", err)
// 	}

// 	logger.InfoWithFields("Running container", map[string]interface{}{
// 		"domains": domains,
// 		"port":    port,
// 		"envVars": len(envVars),
// 		"appType": app.AppType,
// 	})

// 	if err := RunContainer(app, imageTag, containerName, domains, port, envVars, logfile); err != nil {
// 		logger.Error(err, "Failed to run container")
// 		dep.Status = "failed"
// 		dep.Stage = "failed"
// 		dep.Progress = 0
// 		errMsg := fmt.Sprintf("Failed to run container: %v", err)
// 		dep.ErrorMessage = &errMsg
// 		UpdateDeployment(dep, db)
// 		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
// 		UpdateAppStatus(app.ID, "error", db)
// 		return fmt.Errorf("run container failed: %w", err)
// 	}

// 	dep.Status = "success"
// 	dep.Stage = "success"
// 	dep.Progress = 100
// 	now := time.Now()
// 	dep.FinishedAt = &now
// 	UpdateDeployment(dep, db)
// 	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

// 	logger.Info("Updating app status to running")
// 	err = UpdateAppStatus(app.ID, "running", db)
// 	if err != nil {
// 		logger.Error(err, "Failed to update app status (non-fatal)")
// 	}

// 	logger.Info("Cleaning up old Docker images")
// 	if err := CleanupOldImages(app.ID, 5); err != nil {
// 		logger.Error(err, "Failed to cleanup old images (non-fatal)")
// 	}

// 	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
// 		"deployment_id": dep.ID,
// 		"container":     containerName,
// 		"app_status":    "running",
// 	})

// 	return nil
// }

// func UpdateDeployment(dep *models.Deployment, db *gorm.DB) error {
// 	stmt, err := db.Prepare("UPDATE deployments SET status=?, stage=?, progress=?, logs=?, error_message=?, finished_at=? WHERE id=?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(dep.Status, dep.Stage, dep.Progress, dep.Logs, dep.ErrorMessage, dep.FinishedAt, dep.ID)
// 	return err
// }

// func GetLogsPath(commitHash string, depId int64) string {
// 	return filepath.Join(constants.Constants["LogPath"].(string), commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
// }

// func UpdateAppStatus(appID int64, status string, db *sql.DB) error {
// 	query := `UPDATE apps SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
// 	_, err := db.Exec(query, status, appID)
// 	return err
// }

// func GetDeploymentConfig(deploymentID int64, app *models.App, db *sql.DB) (int, []string, map[string]string, error) {
// 	appID, err := models.GetAppIDByDeploymentID(deploymentID)
// 	if err != nil {
// 		return 0, nil, nil, fmt.Errorf("get app ID failed: %w", err)
// 	}

// 	var port *int
// 	err = db.QueryRow("SELECT port FROM apps WHERE id = ?", appID).Scan(&port)
// 	if err != nil {
// 		return 0, nil, nil, fmt.Errorf("get port failed: %w", err)
// 	}
// 	if port == nil {
// 		defaultPort := 3000
// 		port = &defaultPort
// 	}

// 	domains, err := models.GetDomainsByAppID(appID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return 0, nil, nil, fmt.Errorf("get domains failed: %w", err)
// 	}

// 	var domainStrings []string
// 	for _, d := range domains {
// 		domainStrings = append(domainStrings, d.Domain)
// 	}

// 	envs, err := models.GetEnvVariablesByAppID(appID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return 0, nil, nil, fmt.Errorf("get env variables failed: %w", err)
// 	}

// 	envMap := make(map[string]string)

// 	if app.AppType == models.AppTypeDatabase && app.TemplateName != nil {
// 		template, err := models.GetServiceTemplateByName(*app.TemplateName)
// 		if err == nil && template != nil && template.DefaultEnvVars != nil {
// 			var defaultEnvs map[string]string
// 			if err := json.Unmarshal([]byte(*template.DefaultEnvVars), &defaultEnvs); err == nil {
// 				for k, v := range defaultEnvs {
// 					envMap[k] = v
// 				}
// 			}
// 		}
// 	}

// 	for _, env := range envs {
// 		envMap[env.Key] = env.Value
// 	}

// 	return *port, domainStrings, envMap, nil
// }
