package models

import (
	"fmt"
	"time"

	"github.com/corecollectives/mist/utils"
)

type DeploymentStrategy string
type AppStatus string
type AppType string
type RestartPolicy string

const (
	DeploymentAuto   DeploymentStrategy = "auto"
	DeploymentManual DeploymentStrategy = "manual"

	StatusStopped   AppStatus = "stopped"
	StatusRunning   AppStatus = "running"
	StatusError     AppStatus = "error"
	StatusBuilding  AppStatus = "building"
	StatusDeploying AppStatus = "deploying"

	AppTypeWeb      AppType = "web"
	AppTypeService  AppType = "service"
	AppTypeDatabase AppType = "database"

	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)

type App struct {
	ID                  int64              `gorm:"primaryKey;autoIncrement:false" json:"id"`
	ProjectID           int64              `gorm:"uniqueIndex:idx_project_app_name;index;not null" json:"project_id"`
	Name                string             `gorm:"uniqueIndex:idx_project_app_name;not null" json:"name"`
	CreatedBy           int64              `gorm:"index" json:"created_by"`
	Description         *string            `json:"description,omitempty"`
	AppType             AppType            `gorm:"default:'web';index" json:"app_type"`
	TemplateName        *string            `json:"template_name,omitempty"`
	GitProviderID       *int64             `json:"git_provider_id,omitempty"`
	GitRepository       *string            `json:"git_repository,omitempty"`
	GitBranch           string             `gorm:"default:'main'" json:"git_branch,omitempty"`
	GitCloneURL         *string            `json:"git_clone_url,omitempty"`
	DeploymentStrategy  DeploymentStrategy `gorm:"default:'auto'" json:"deployment_strategy"`
	Port                *int64             `json:"port,omitempty"`
	RootDirectory       string             `gorm:"default:'.'" json:"root_directory,omitempty"`
	BuildCommand        *string            `json:"build_command,omitempty"`
	StartCommand        *string            `json:"start_command,omitempty"`
	DockerfilePath      *string            `gorm:"default:'DOCKERFILE'" json:"dockerfile_path,omitempty"`
	CPULimit            *float64           `json:"cpu_limit,omitempty"`
	MemoryLimit         *int               `json:"memory_limit,omitempty"`
	RestartPolicy       RestartPolicy      `gorm:"default:'unless-stopped'" json:"restart_policy"`
	HealthcheckPath     *string            `json:"healthcheck_path,omitempty"`
	HealthcheckInterval int                `gorm:"default:30" json:"healthcheck_interval"`
	HealthcheckTimeout  int                `gorm:"default:10" json:"healthcheck_timeout"`
	HealthcheckRetries  int                `gorm:"default:3" json:"healthcheck_retries"`
	Status              AppStatus          `gorm:"default:'stopped';index" json:"status"`
	CreatedAt           time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

func (a *App) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":                  a.ID,
		"projectId":           a.ProjectID,
		"createdBy":           a.CreatedBy,
		"name":                a.Name,
		"description":         a.Description,
		"appType":             a.AppType,
		"templateName":        a.TemplateName,
		"gitProviderId":       a.GitProviderID,
		"gitRepository":       a.GitRepository,
		"gitBranch":           a.GitBranch,
		"gitCloneUrl":         a.GitCloneURL,
		"deploymentStrategy":  a.DeploymentStrategy,
		"port":                a.Port,
		"rootDirectory":       a.RootDirectory,
		"buildCommand":        a.BuildCommand,
		"startCommand":        a.StartCommand,
		"dockerfilePath":      a.DockerfilePath,
		"cpuLimit":            a.CPULimit,
		"memoryLimit":         a.MemoryLimit,
		"restartPolicy":       a.RestartPolicy,
		"healthcheckPath":     a.HealthcheckPath,
		"healthcheckInterval": a.HealthcheckInterval,
		"healthcheckTimeout":  a.HealthcheckTimeout,
		"healthcheckRetries":  a.HealthcheckRetries,
		"status":              a.Status,
		"createdAt":           a.CreatedAt,
		"updatedAt":           a.UpdatedAt,
	}
}

func (a *App) InsertInDB() error {
	a.ID = utils.GenerateRandomId()
	if a.AppType == "" {
		a.AppType = AppTypeWeb
	}
	if a.RestartPolicy == "" {
		a.RestartPolicy = RestartPolicyUnlessStopped
	}
	if a.DeploymentStrategy == "" {
		a.DeploymentStrategy = DeploymentAuto
	}
	if a.Status == "" {
		a.Status = StatusStopped
	}

	return db.Create(a).Error
}

func GetApplicationByProjectID(projectId int64) ([]App, error) {
	var apps []App
	result := db.Where("project_id=?", projectId).Find(&apps)
	return apps, result.Error
}

func GetApplicationByID(appId int64) (*App, error) {
	var app App
	result := db.First(&app, "id=?", appId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

func (a *App) UpdateApplication() error {
	return db.Model(a).Select("Name", "Description", "AppType", "TemplateName",
		"GitProviderID", "GitRepository", "GitBranch", "GitCloneURL",
		"DeploymentStrategy", "Port", "RootDirectory",
		"BuildCommand", "StartCommand", "DockerfilePath",
		"CPULimit", "MemoryLimit", "RestartPolicy",
		"HealthcheckPath", "HealthcheckInterval", "HealthcheckTimeout", "HealthcheckRetries",
		"Status", "UpdatedAt").Updates(a).Error
}

func IsUserApplicationOwner(userId int64, appId int64) (bool, error) {
	var count int64
	err := db.Model(&App{}).
		Where("id = ? AND created_by = ?", appId, userId).
		Count(&count).Error

	return count > 0, err
}

func FindApplicationIDByGitRepoAndBranch(gitRepo string, gitBranch string) (int64, error) {
	var app App
	err := db.Select("id").
		Where("git_repository = ? AND git_branch = ?", gitRepo, gitBranch).
		First(&app).Error

	if err != nil {
		return 0, err
	}
	return app.ID, nil
}

func GetUserIDByAppID(appID int64) (*int64, error) {
	var app App
	err := db.Select("created_by").First(&app, appID).Error
	if err != nil {
		return nil, err
	}
	return &app.CreatedBy, nil
}

func GetAppIDByDeploymentID(depId int64) (int64, error) {
	var result struct {
		AppID int64
	}
	err := db.Table("deployments").Select("app_id").Where("id=?", depId).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.AppID, nil
}

func GetAppRepoInfo(appId int64) (string, string, int64, string, error) {
	var app App
	err := db.Select("git_repository, git_branch, project_id, name").
		First(&app, appId).Error

	if err != nil {
		return "", "", 0, "", err
	}

	repo := ""
	if app.GitRepository != nil {
		repo = *app.GitRepository
	}

	return repo, app.GitBranch, app.ProjectID, app.Name, nil
}

func GetAppRepoAndBranch(appId int64) (string, string, error) {
	var app App
	err := db.Select("git_repository, git_branch").First(&app, appId).Error
	if err != nil {
		return "", "", err
	}
	if app.GitRepository == nil || *app.GitRepository == "" {
		return "", "", fmt.Errorf("app has no git repository configured")
	}

	branch := app.GitBranch
	if branch == "" {
		branch = "main"
	}

	return *app.GitRepository, branch, nil
}

func DeleteApplication(appID int64) error {
	return db.Delete(&App{}, appID).Error
}

// if git_clone_url is set, it uses that
// if not, it falls back to git_repository
// assuming github, becuase in v1.0.0, we only had github as a git provider and we were only storing git_repository and appending the https://github.com/ manually, bad design decision but whatever
func GetAppCloneURL(appID int64, userID int64) (string, string, bool, error) {
	gitProviderID, gitRepository, _, gitCloneURL, _, _, err := GetAppGitInfo(appID)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get app git info: %w", err)
	}

	// case 1 - git_clone_url is set, use it
	if gitCloneURL != nil && *gitCloneURL != "" {
		var accessToken string
		if gitProviderID != nil {
			token, _, needsRefresh, err := GetGitProviderAccessToken(*gitProviderID)
			if err != nil {
				// try to get GitHub access token as fallback
				_, err := GetInstallationIDByUserID(userID)
				if err == nil {
					installation, err := GetInstallationByUserID(int(userID))
					if err == nil {
						accessToken = installation.AccessToken
					}
				}
			} else {
				// Check if token needs refresh
				if needsRefresh {
					// Try to refresh the token
					newToken, err := RefreshGitProviderToken(*gitProviderID)
					if err == nil {
						token = newToken
					}
					// If refresh fails, continue with expired token - deployment will fail but at least we tried
				}
				accessToken = token
			}
		}
		return *gitCloneURL, accessToken, false, nil
	}

	// Case 2 - git_clone_url missing but git_repository present (legacy github apps)
	if gitRepository != nil && *gitRepository != "" {
		// assume GitHub and construct the clone URL
		cloneURL := fmt.Sprintf("https://github.com/%s.git", *gitRepository)

		// ary to get access token from GitHub installation
		var accessToken string
		_, err := GetInstallationIDByUserID(userID)
		if err == nil {
			installation, err := GetInstallationByUserID(int(userID))
			if err == nil {
				accessToken = installation.AccessToken
			}
		}

		// mark that we should migrate this app's data
		return cloneURL, accessToken, true, nil
	}

	return "", "", false, fmt.Errorf("app has no git repository or clone URL configured")
}
