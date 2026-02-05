package compose

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/moby/moby/client"
)

func ComposeRestart(appContextPath string) error {
	cmd := exec.Command("docker", "compose", "restart")
	cmd.Dir = appContextPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose restart failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	projectName := getProjectNameFromPath(appContextPath)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("compose restart command completed but containers did not reach running state within timeout")
		case <-ticker.C:
			containers, err := cli.ContainerList(ctx, client.ContainerListOptions{
				All: true,
			})
			if err != nil {
				continue
			}

			allRunning := true
			foundProjectContainers := false
			for _, container := range containers.Items {
				if projectName != "" {
					if composeProject, ok := container.Labels["com.docker.compose.project"]; ok && composeProject == projectName {
						foundProjectContainers = true
						if container.State != "running" {
							allRunning = false
							break
						}
					}
				}
			}

			if foundProjectContainers && allRunning {
				return nil
			}
		}
	}
}

func getProjectNameFromPath(path string) string {
	if len(path) == 0 {
		return ""
	}
	parts := make([]string, 0)
	for path != "" {
		dir := path
		if idx := len(path) - 1; idx >= 0 && path[idx] == '/' {
			path = path[:idx]
			continue
		}
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				dir = path[i+1:]
				path = path[:i]
				break
			}
			if i == 0 {
				dir = path
				path = ""
			}
		}
		if dir != "" {
			parts = append([]string{dir}, parts...)
		}
	}
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
