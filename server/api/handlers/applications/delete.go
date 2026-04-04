package applications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	appIDStr := r.URL.Query().Get("id")
	var appID int64

	if appIDStr != "" {
		var err error
		appID, err = strconv.ParseInt(appIDStr, 10, 64)
		if err != nil {
			handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid app ID", "App ID must be a valid number")
			return
		}
	} else {
		var req struct {
			AppID int64 `json:"appId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request", "App ID must be provided as query parameter or in request body")
			return
		}
		appID = req.AppID
	}

	if appID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	app, err := models.GetApplicationByID(appID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", fmt.Sprintf("Error fetching application: %v", err))
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "No application with the given ID exists")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "Forbidden")
		return
	}

	containerName := docker.GetContainerName(app.Name, app.ID)

	if docker.ContainerExists(containerName) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		cli, err := client.New(client.FromEnv)
		if err != nil {
			log.Warn().Err(err).Str("container", containerName).Msg("Failed to create Docker client during app deletion")
		} else {
			_, err = cli.ContainerStop(ctx, containerName, client.ContainerStopOptions{})
			if err != nil {
				log.Warn().Err(err).Str("container", containerName).Msg("Failed to stop container during app deletion")
			}

			_, err = cli.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{})
			if err != nil {
				log.Warn().Err(err).Str("container", containerName).Msg("Failed to remove container during app deletion")
			}
		}

		// legacy exec method
		//
		//
		// stopCmd := exec.CommandContext(ctx, "docker", "stop", containerName)
		// if err := stopCmd.Run(); err != nil {
		// 	log.Warn().Err(err).Str("container", containerName).Msg("Failed to stop container during app deletion")
		// }
		//
		// removeCmd := exec.CommandContext(ctx, "docker", "rm", containerName)
		// if err := removeCmd.Run(); err != nil {
		// 	log.Warn().Err(err).Str("container", containerName).Msg("Failed to remove container during app deletion")
		// }
	}

	imagePattern := fmt.Sprintf("mist-app-%d-", app.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		log.Warn().Err(err).Int64("app_id", app.ID).Msg("Failed to create Docker client for image cleanup during app deletion")
	} else {
		filterArgs := make(client.Filters)
		filterArgs.Add("reference", fmt.Sprintf("%s*", imagePattern))

		imageListResult, err := cli.ImageList(ctx, client.ImageListOptions{
			Filters: filterArgs,
		})
		if err == nil && len(imageListResult.Items) > 0 {
			for _, img := range imageListResult.Items {
				_, err := cli.ImageRemove(ctx, img.ID, client.ImageRemoveOptions{
					Force: true,
				})
				if err != nil {
					log.Warn().Err(err).Str("image_id", img.ID).Int64("app_id", app.ID).Msg("Failed to remove Docker image during app deletion")
				}
			}
		}
	}

	// legacy exec method
	//
	//
	// listImagesCmd := exec.CommandContext(ctx, "docker", "images", "-q", "--filter", fmt.Sprintf("reference=%s*", imagePattern))
	// output, err := listImagesCmd.Output()
	// if err == nil && len(output) > 0 {
	// 	rmiCmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("docker images -q --filter 'reference=%s*' | xargs -r docker rmi -f", imagePattern))
	// 	if err := rmiCmd.Run(); err != nil {
	// 		log.Warn().Err(err).Int64("app_id", app.ID).Msg("Failed to remove Docker images during app deletion")
	// 	}
	// }

	appPath := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
	if _, err := os.Stat(appPath); err == nil {
		if err := os.RemoveAll(appPath); err != nil {
			log.Warn().Err(err).Str("path", appPath).Msg("Failed to remove app directory during deletion")
		}
	}

	logPath := constants.Constants["LogPath"].(string)
	logPattern := filepath.Join(logPath, fmt.Sprintf("*%d_build_logs", app.ID))
	matches, err := filepath.Glob(logPattern)
	if err == nil {
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				log.Warn().Err(err).Str("log_file", match).Msg("Failed to remove log file during app deletion")
			}
		}
	}

	err = models.DeleteApplication(appID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete application from database", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "delete", "application", &appID, map[string]interface{}{
		"app_name":   app.Name,
		"project_id": app.ProjectID,
	})

	handlers.SendResponse(w, http.StatusOK, true, nil, "Application deleted successfully", "")
}
