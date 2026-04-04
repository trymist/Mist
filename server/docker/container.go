package docker

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

func StopContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error making moby client: %s", err.Error())
	}
	_, err = cli.ContainerStop(ctx, containerName, client.ContainerStopOptions{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker stop timed out after 2 minutes for container %s", containerName)
		}
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

func StartContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}
	_, err = cli.ContainerStart(ctx, containerName, client.ContainerStartOptions{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker start timed out after 1 minute for container %s", containerName)
		}
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil

}

func RestartContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	_, err = cli.ContainerRestart(ctx, containerName, client.ContainerRestartOptions{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("dockre restart timed out after 3 minuts for container %s", containerName)
		}
		return fmt.Errorf("failed to restart container: %w", err)
	}

	pollCtx, pollCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer pollCancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pollCtx.Done():
			return fmt.Errorf("container %s restart command completed but container did not reach running state within timeout", containerName)
		case <-ticker.C:
			inspectResult, err := cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
			if err != nil {
				continue
			}
			if inspectResult.Container.State != nil && inspectResult.Container.State.Running {
				return nil
			}
		}
	}
}

func GetContainerLogs(containerName string, tail int) (string, error) {
	if !ContainerExists(containerName) {
		return "", fmt.Errorf("container %s does not exist", containerName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating moby client: %s", err.Error())
	}

	tailStr := fmt.Sprintf("%d", tail)
	options := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tailStr,
	}

	logReader, err := cli.ContainerLogs(ctx, containerName, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logReader.Close()

	var logs strings.Builder
	buf := make([]byte, 8192)
	for {
		n, err := logReader.Read(buf)
		if n > 0 {
			logs.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	return logs.String(), nil

}

func GetContainerName(appName string, appId int64) string {
	return fmt.Sprintf("app-%d", appId)
}

func StopAndRemoveContainer(containerName string, logfile *os.File) error {
	ifExists := ContainerExists(containerName)
	if !ifExists {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}
	_, err = cli.ContainerStop(ctx, containerName, client.ContainerStopOptions{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker stop timed out after 2 minutes for container %s", containerName)
		}
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}

	_, err = cli.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	}

	return nil

}
func ContainerExists(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		log.Error().Err(err).Msg("failed to create docker client")
		return false
	}
	_, err = cli.ContainerInspect(ctx, name, client.ContainerInspectOptions{})
	if err != nil {
		log.Error().Msg("container not found")
		return false
	}
	return true

}

func CreateAndStartContainer(ctx context.Context, app *models.App, imageTag, containerName string, domains []string, Port int, runtimeEnvVars map[string]string, logfile *os.File) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}
	var restartPolicy container.RestartPolicyMode
	switch string(app.RestartPolicy) {

	case string(models.RestartPolicyNo):
		restartPolicy = container.RestartPolicyDisabled
	case string(models.RestartPolicyAlways):
		restartPolicy = container.RestartPolicyAlways
	case string(models.RestartPolicyOnFailure):
		restartPolicy = container.RestartPolicyOnFailure
	case string(models.RestartPolicyUnlessStopped):
		restartPolicy = container.RestartPolicyUnlessStopped
	default:
		restartPolicy = container.RestartPolicyUnlessStopped
	}

	var volumeBinds []string
	volumes, err := models.GetVolumesByAppID(app.ID)
	if err == nil {
		for _, vol := range volumes {
			volumeBindArg := fmt.Sprintf("%s:%s", vol.HostPath, vol.ContainerPath)
			if vol.ReadOnly {
				volumeBindArg += ":ro"
			}
			volumeBinds = append(volumeBinds, volumeBindArg)
		}
	}

	var envList []string
	for key, value := range runtimeEnvVars {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	labels := make(map[string]string)

	networkMode := ""
	exposedPorts := make(network.PortSet)
	portBindings := make(network.PortMap)

	switch app.AppType {
	case models.AppTypeWeb:
		networkMode = "traefik-net"

		if len(domains) > 0 {
			labels["traefik.enable"] = "true"

			var hostRules []string
			for _, domain := range domains {
				hostRules = append(hostRules, fmt.Sprintf("Host(`%s`)", domain))
			}
			hostRule := strings.Join(hostRules, " || ")

			labels[fmt.Sprintf("traefik.http.routers.%s.rule", containerName)] = hostRule
			labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName)] = "websecure"
			labels[fmt.Sprintf("traefik.http.routers.%s.tls", containerName)] = "true"
			labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", containerName)] = "le"
			labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", containerName)] = fmt.Sprintf("%d", Port)

			labels[fmt.Sprintf("traefik.http.routers.%s-http.rule", containerName)] = hostRule
			labels[fmt.Sprintf("traefik.http.routers.%s-http.entrypoints", containerName)] = "web"
			labels[fmt.Sprintf("traefik.http.routers.%s-http.middlewares", containerName)] = fmt.Sprintf("%s-https-redirect", containerName)
			labels[fmt.Sprintf("traefik.http.middlewares.%s-https-redirect.redirectscheme.scheme", containerName)] = "https"
		}

		shouldExpose := app.ShouldExpose != nil && *app.ShouldExpose

		if shouldExpose {
			exposePort := Port
			if app.ExposePort != nil && *app.ExposePort > 0 {
				exposePort = int(*app.ExposePort)
			}

			port, err := network.ParsePort(fmt.Sprintf("%d/tcp", Port))
			if err != nil {
				return fmt.Errorf("failed to parse port: %w", err)
			}
			hostIP, err := netip.ParseAddr("0.0.0.0")
			if err != nil {
				return fmt.Errorf("failed to parse host IP: %w", err)
			}
			exposedPorts[port] = struct{}{}
			portBindings[port] = []network.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: fmt.Sprintf("%d", exposePort),
				},
			}
		}

	case models.AppTypeService, models.AppTypeDatabase:
		networkMode = "traefik-net"

		shouldExpose := app.ShouldExpose != nil && *app.ShouldExpose

		if shouldExpose {
			exposePort := Port
			if app.ExposePort != nil && *app.ExposePort > 0 {
				exposePort = int(*app.ExposePort)
			}

			port, err := network.ParsePort(fmt.Sprintf("%d/tcp", Port))
			if err != nil {
				return fmt.Errorf("failed to parse port: %w", err)
			}
			hostIP, err := netip.ParseAddr("0.0.0.0")
			if err != nil {
				return fmt.Errorf("failed to parse host IP: %w", err)
			}
			exposedPorts[port] = struct{}{}
			portBindings[port] = []network.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: fmt.Sprintf("%d", exposePort),
				},
			}
		}

	default:
		networkMode = "traefik-net"

		shouldExpose := app.ShouldExpose != nil && *app.ShouldExpose

		if shouldExpose && app.AppType != models.AppTypeCompose {
			exposePort := Port
			if app.ExposePort != nil && *app.ExposePort > 0 {
				exposePort = int(*app.ExposePort)
			}

			port, err := network.ParsePort(fmt.Sprintf("%d/tcp", Port))
			if err != nil {
				return fmt.Errorf("failed to parse port: %w", err)
			}
			hostIP, err := netip.ParseAddr("0.0.0.0")
			if err != nil {
				return fmt.Errorf("failed to parse host IP: %w", err)
			}
			exposedPorts[port] = struct{}{}
			portBindings[port] = []network.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: fmt.Sprintf("%d", exposePort),
				},
			}
		}
	}

	hostConfig := container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: restartPolicy,
		},
		Resources:    container.Resources{},
		Binds:        volumeBinds,
		NetworkMode:  container.NetworkMode(networkMode),
		PortBindings: portBindings,
	}

	config := container.Config{
		Image:        imageTag,
		Env:          envList,
		Labels:       labels,
		ExposedPorts: exposedPorts,
	}

	if app.CPULimit != nil && *app.CPULimit > 0 {
		hostConfig.Resources.NanoCPUs = int64(*app.CPULimit * 1e9)
	}
	if app.MemoryLimit != nil && *app.MemoryLimit > 0 {
		hostConfig.Resources.Memory = int64(*app.MemoryLimit) * 1024 * 1024
	}

	resp, err := cli.ContainerCreate(timeoutCtx, client.ContainerCreateOptions{
		Name:       containerName,
		Config:     &config,
		HostConfig: &hostConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	_, err = cli.ContainerStart(timeoutCtx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil

}

type ContainerStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Uptime  string `json:"uptime"`
	Healthy bool   `json:"healthy"`
}

func GetContainerStatus(containerName string) (*ContainerStatus, error) {
	if !ContainerExists(containerName) {
		return &ContainerStatus{
			Name:    containerName,
			Status:  "not_found",
			State:   "stopped",
			Uptime:  "N/A",
			Healthy: false,
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating moby client: %s", err.Error())
	}

	inspectResult, err := cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	inspectData := inspectResult.Container

	uptime := "N/A"
	if inspectData.State != nil && inspectData.State.StartedAt != "" {
		uptime = inspectData.State.StartedAt
	}

	state := "stopped"
	status := ""
	if inspectData.State != nil {
		status = string(inspectData.State.Status)
		if inspectData.State.Running {
			state = "running"
		} else if inspectData.State.Status == "exited" {
			state = "stopped"
		} else {
			state = string(inspectData.State.Status)
		}
	}

	healthy := true
	if inspectData.State != nil && inspectData.State.Health != nil {
		healthy = inspectData.State.Health.Status == "healthy"
	}

	return &ContainerStatus{
		Name:    strings.TrimPrefix(inspectData.Name, "/"),
		Status:  status,
		State:   state,
		Uptime:  uptime,
		Healthy: healthy,
	}, nil

}

func RecreateContainer(app *models.App) error {
	containerName := GetContainerName(app.Name, app.ID)

	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	inspectResult, err := cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get container image: %w", err)
	}
	imageTag := inspectResult.Container.Image

	port, domains, envSet, err := FetchDeploymentConfigurationForApp(app)
	if err != nil {
		return fmt.Errorf("failed to fetch deployment configuration: %w", err)
	}

	if err := StopAndRemoveContainer(containerName, nil); err != nil {
		return fmt.Errorf("failed to stop and remove container: %w", err)
	}

	if err := CreateAndStartContainer(ctx, app, imageTag, containerName, domains, port, envSet.Runtime, nil); err != nil {
		return fmt.Errorf("failed to create and start container: %w", err)
	}

	return nil

}
