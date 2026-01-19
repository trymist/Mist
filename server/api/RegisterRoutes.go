package api

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/handlers/applications"
	"github.com/corecollectives/mist/api/handlers/auditlogs"
	"github.com/corecollectives/mist/api/handlers/auth"
	"github.com/corecollectives/mist/api/handlers/deployments"
	"github.com/corecollectives/mist/api/handlers/github"
	"github.com/corecollectives/mist/api/handlers/projects"
	"github.com/corecollectives/mist/api/handlers/settings"
	"github.com/corecollectives/mist/api/handlers/templates"
	"github.com/corecollectives/mist/api/handlers/updates"
	"github.com/corecollectives/mist/api/handlers/users"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/websockets"
)

func RegisterRoutes(mux *http.ServeMux) {

	avatarDir := constants.Constants["AvatarDirPath"].(string)
	mux.Handle("/uploads/avatar/", http.StripPrefix("/uploads/avatar/", http.FileServer(http.Dir(avatarDir))))

	mux.Handle("/api/ws/stats", middleware.AuthMiddleware()(http.HandlerFunc(websockets.StatWsHandler)))
	mux.HandleFunc("/api/ws/container/logs", websockets.ContainerLogsHandler)
	mux.HandleFunc("/api/ws/container/stats", websockets.ContainerStatsHandler)
	mux.Handle("/api/ws/system/logs", middleware.AuthMiddleware()(http.HandlerFunc(websockets.SystemLogsHandler)))
	mux.HandleFunc("GET /api/health", handlers.HealthCheckHandler)

	mux.HandleFunc("POST /api/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("POST /api/auth/login", auth.LoginHandler)
	mux.HandleFunc("GET /api/auth/me", auth.MeHandler)
	mux.HandleFunc("POST /api/auth/logout", auth.LogoutHandler)

	mux.Handle("POST /api/users/create", middleware.AuthMiddleware()(http.HandlerFunc(users.CreateUser)))
	mux.Handle("GET /api/users/getAll", middleware.AuthMiddleware()(http.HandlerFunc(users.GetUsers)))
	mux.Handle("GET /api/users/getFromId", middleware.AuthMiddleware()(http.HandlerFunc(users.GetUserById)))
	mux.Handle("PUT /api/users/update", middleware.AuthMiddleware()(http.HandlerFunc(users.UpdateUser)))
	mux.Handle("PUT /api/users/password", middleware.AuthMiddleware()(http.HandlerFunc(users.UpdatePassword)))
	mux.Handle("POST /api/users/avatar", middleware.AuthMiddleware()(http.HandlerFunc(users.UploadAvatar)))
	mux.Handle("DELETE /api/users/avatar", middleware.AuthMiddleware()(http.HandlerFunc(users.DeleteAvatar)))
	mux.Handle("DELETE /api/users/delete", middleware.AuthMiddleware()(http.HandlerFunc(users.DeleteUser)))
	mux.Handle("GET /api/users/git-providers", middleware.AuthMiddleware()(http.HandlerFunc(users.GetUserGitProviders)))

	mux.Handle("POST /api/projects/create", middleware.AuthMiddleware()(http.HandlerFunc(projects.CreateProject)))
	mux.Handle("GET /api/projects/getAll", middleware.AuthMiddleware()(http.HandlerFunc(projects.GetProjects)))
	mux.Handle("GET /api/projects/getFromId", middleware.AuthMiddleware()(http.HandlerFunc(projects.GetProjectFromId)))
	mux.Handle("PUT /api/projects/update", middleware.AuthMiddleware()(http.HandlerFunc(projects.UpdateProject)))
	mux.Handle("DELETE /api/projects/delete", middleware.AuthMiddleware()(http.HandlerFunc(projects.DeleteProject)))
	mux.Handle("PUT /api/projects/updateMembers", middleware.AuthMiddleware()(http.HandlerFunc(projects.UpdateMembers)))

	mux.Handle("POST /api/apps/create", middleware.AuthMiddleware()(http.HandlerFunc(applications.CreateApplication)))
	mux.Handle("POST /api/apps/getByProjectId", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetApplicationByProjectID)))
	mux.Handle("POST /api/apps/getById", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetApplicationById)))
	mux.Handle("PUT /api/apps/update", middleware.AuthMiddleware()(http.HandlerFunc(applications.UpdateApplication)))
	mux.Handle("DELETE /api/apps/delete", middleware.AuthMiddleware()(http.HandlerFunc(applications.DeleteApplication)))
	mux.Handle("POST /api/apps/getLatestCommit", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetLatestCommit)))
	mux.Handle("POST /api/apps/getPreviewUrl", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetPreviewURL)))

	mux.Handle("POST /api/apps/envs/create", middleware.AuthMiddleware()(http.HandlerFunc(applications.CreateEnvVariable)))
	mux.Handle("POST /api/apps/envs/get", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetEnvVariables)))
	mux.Handle("PUT /api/apps/envs/update", middleware.AuthMiddleware()(http.HandlerFunc(applications.UpdateEnvVariable)))
	mux.Handle("DELETE /api/apps/envs/delete", middleware.AuthMiddleware()(http.HandlerFunc(applications.DeleteEnvVariable)))

	mux.Handle("POST /api/apps/domains/create", middleware.AuthMiddleware()(http.HandlerFunc(applications.CreateDomain)))
	mux.Handle("POST /api/apps/domains/get", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetDomains)))
	mux.Handle("PUT /api/apps/domains/update", middleware.AuthMiddleware()(http.HandlerFunc(applications.UpdateDomain)))
	mux.Handle("DELETE /api/apps/domains/delete", middleware.AuthMiddleware()(http.HandlerFunc(applications.DeleteDomain)))
	mux.Handle("POST /api/apps/domains/verify", middleware.AuthMiddleware()(http.HandlerFunc(applications.VerifyDomainDNS)))
	mux.Handle("POST /api/apps/domains/instructions", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetDNSInstructions)))

	mux.Handle("POST /api/apps/volumes/get", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetVolumes)))
	mux.Handle("POST /api/apps/volumes/create", middleware.AuthMiddleware()(http.HandlerFunc(applications.CreateVolume)))
	mux.Handle("PUT /api/apps/volumes/update", middleware.AuthMiddleware()(http.HandlerFunc(applications.UpdateVolume)))
	mux.Handle("DELETE /api/apps/volumes/delete", middleware.AuthMiddleware()(http.HandlerFunc(applications.DeleteVolume)))

	mux.Handle("POST /api/apps/container/stop", middleware.AuthMiddleware()(http.HandlerFunc(applications.StopContainerHandler)))
	mux.Handle("POST /api/apps/container/start", middleware.AuthMiddleware()(http.HandlerFunc(applications.StartContainerHandler)))
	mux.Handle("POST /api/apps/container/restart", middleware.AuthMiddleware()(http.HandlerFunc(applications.RestartContainerHandler)))
	mux.Handle("GET /api/apps/container/status", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetContainerStatusHandler)))
	mux.Handle("GET /api/apps/container/logs", middleware.AuthMiddleware()(http.HandlerFunc(applications.GetContainerLogsHandler)))

	mux.Handle("GET /api/github/app", middleware.AuthMiddleware()(http.HandlerFunc(github.GetApp)))
	mux.Handle("GET /api/github/app/create", middleware.AuthMiddleware()(http.HandlerFunc(github.CreateGithubApp)))
	mux.Handle("GET /api/github/callback", http.HandlerFunc(github.CallBackHandler))
	mux.Handle("GET /api/github/installation/callback", http.HandlerFunc(github.HandleInstallationEvent))
	mux.Handle("GET /api/github/repositories", middleware.AuthMiddleware()(http.HandlerFunc(github.GetRepositories)))
	mux.Handle("POST /api/github/branches", middleware.AuthMiddleware()(http.HandlerFunc(github.GetBranches)))
	mux.HandleFunc("POST /api/github/webhook", github.GithubWebhook)

	mux.HandleFunc("/api/deployments/logs/stream", deployments.LogsHandler)
	mux.Handle("POST /api/deployments", middleware.AuthMiddleware()(http.HandlerFunc(deployments.AddDeployHandler)))
	mux.Handle("POST /api/deployments/create", middleware.AuthMiddleware()(http.HandlerFunc(deployments.AddDeployHandler)))
	mux.Handle("POST /api/deployments/getByAppId", middleware.AuthMiddleware()(http.HandlerFunc(deployments.GetByApplicationID)))
	mux.Handle("GET /api/deployments/logs", middleware.AuthMiddleware()(http.HandlerFunc(deployments.GetCompletedDeploymentLogsHandler)))
	mux.Handle("POST /api/deployments/stopDep", middleware.AuthMiddleware()(http.HandlerFunc(deployments.StopDeployment)))

	mux.Handle("GET /api/templates/list", middleware.AuthMiddleware()(http.HandlerFunc(templates.ListServiceTemplates)))
	mux.Handle("GET /api/templates/get", middleware.AuthMiddleware()(http.HandlerFunc(templates.GetServiceTemplateByName)))
	mux.Handle("GET /api/audit-logs", middleware.AuthMiddleware()(http.HandlerFunc(auditlogs.GetAllAuditLogs)))
	mux.Handle("GET /api/audit-logs/resource", middleware.AuthMiddleware()(http.HandlerFunc(auditlogs.GetAuditLogsByResource)))

	mux.Handle("GET /api/settings/system", middleware.AuthMiddleware()(http.HandlerFunc(settings.GetSystemSettings)))
	mux.Handle("PUT /api/settings/system", middleware.AuthMiddleware()(http.HandlerFunc(settings.UpdateSystemSettings)))
	mux.Handle("POST /api/settings/docker/cleanup", middleware.AuthMiddleware()(http.HandlerFunc(settings.DockerCleanup)))

	mux.Handle("GET /api/updates/version", middleware.AuthMiddleware()(http.HandlerFunc(updates.GetCurrentVersion)))
	mux.Handle("GET /api/updates/check", middleware.AuthMiddleware()(http.HandlerFunc(updates.CheckForUpdates)))
	mux.Handle("POST /api/updates/trigger", middleware.AuthMiddleware()(http.HandlerFunc(updates.TriggerUpdate)))
	mux.Handle("GET /api/updates/history", middleware.AuthMiddleware()(http.HandlerFunc(updates.GetUpdateHistory)))
	mux.Handle("GET /api/updates/log", middleware.AuthMiddleware()(http.HandlerFunc(updates.GetUpdateLogByID)))
	mux.Handle("POST /api/updates/clear", middleware.AuthMiddleware()(http.HandlerFunc(updates.ClearStuckUpdate)))

}
