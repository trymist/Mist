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

	//legacy exec method
	//
	//
	// cmd := exec.CommandContext(ctx, "docker", "stop", containerName)
	// if err := cmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker stop timed out after 2 minutes for container %s", containerName)
	// 	}
	// 	return fmt.Errorf("failed to stop container: %w", err)
	// }
	//
	// return nil
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

	// legacy exec method
	//
	//
	// cmd := exec.CommandContext(ctx, "docker", "start", containerName)
	// if err := cmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker start timed out after 1 minute for container %s", containerName)
	// 	}
	// 	return fmt.Errorf("failed to start container: %w", err)
	// }
	//
	// return nil
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

	// Wait for container to be fully running
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

	// legacy exec method
	//
	//
	// cmd := exec.CommandContext(ctx, "docker", "restart", containerName)
	// if err := cmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker restart timed out after 3 minutes for container %s", containerName)
	// 	}
	// 	return fmt.Errorf("failed to restart container: %w", err)
	// }
	//
	// return nil
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

	// legacy exec method
	//
	//
	// tailStr := fmt.Sprintf("%d", tail)
	// cmd := exec.Command("docker", "logs", "--tail", tailStr, containerName)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get container logs: %w", err)
	// }
	//
	// return string(output), nil
}

func GetContainerName(appName string, appId int64) string {
	return fmt.Sprintf("app-%d", appId)
}

func StopRemoveContainer(containerName string, logfile *os.File) error {
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

	// legacy exec method
	//
	//
	// stopCmd := exec.CommandContext(ctx, "docker", "stop", containerName)
	// stopCmd.Stdout = logfile
	// stopCmd.Stderr = logfile
	// if err := stopCmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker stop timed out after 2 minutes for container %s", containerName)
	// 	}
	// 	return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	// }
	//
	// removeCmd := exec.CommandContext(ctx, "docker", "rm", containerName)
	// removeCmd.Stdout = logfile
	// removeCmd.Stderr = logfile
	// if err := removeCmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	// }
	//
	// return nil
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

	// legacy exec method
	//
	//
	// cmd := exec.Command("docker", "inspect", name)
	// output, err := cmd.CombinedOutput()
	//
	// if err != nil {
	// 	if strings.Contains(string(output), "No such object") {
	// 		return false
	// 	}
	// 	return false
	// }
	//
	// return true
}

func RunContainer(ctx context.Context, app *models.App, imageTag, containerName string, domains []string, Port int, envVars map[string]string, logfile *os.File) error {

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
	for key, value := range envVars {
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

	// legacey exec method
	//
	//
	// runArgs := []string{
	// 	"run", "-d",
	// 	"--name", containerName,
	// }

	// restartPolicy := string(app.RestartPolicy)
	// if restartPolicy == "" {
	// 	restartPolicy = "unless-stopped"
	// }
	// runArgs = append(runArgs, "--restart", restartPolicy)
	//
	// if app.CPULimit != nil && *app.CPULimit > 0 {
	// 	runArgs = append(runArgs, "--cpus", fmt.Sprintf("%.2f", *app.CPULimit))
	// }
	//
	// if app.MemoryLimit != nil && *app.MemoryLimit > 0 {
	// 	runArgs = append(runArgs, "-m", fmt.Sprintf("%dm", *app.MemoryLimit))
	// }
	//
	// // Add volumes from the volumes table (user-configurable)
	// volumes, err := models.GetVolumesByAppID(app.ID)
	// if err == nil {
	// 	for _, vol := range volumes {
	// 		volumeArg := fmt.Sprintf("%s:%s", vol.HostPath, vol.ContainerPath)
	// 		if vol.ReadOnly {
	// 			volumeArg += ":ro"
	// 		}
	// 		runArgs = append(runArgs, "-v", volumeArg)
	// 	}
	// }
	//
	// for key, value := range envVars {
	// 	runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	// }
	//
	// switch app.AppType {
	// case models.AppTypeWeb:
	// 	if len(domains) > 0 {
	// 		runArgs = append(runArgs,
	// 			"--network", "traefik-net",
	// 			"-l", "traefik.enable=true",
	// 		)
	//
	// 		var hostRules []string
	// 		for _, domain := range domains {
	// 			hostRules = append(hostRules, fmt.Sprintf("Host(`%s`)", domain))
	// 		}
	// 		hostRule := strings.Join(hostRules, " || ")
	//
	// 		runArgs = append(runArgs,
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s.rule=%s", containerName, hostRule),
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s.entrypoints=websecure", containerName),
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s.tls=true", containerName),
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=le", containerName),
	// 			"-l", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", containerName, Port),
	// 		)
	//
	// 		runArgs = append(runArgs,
	//
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s-http.rule=%s", containerName, hostRule),
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s-http.entrypoints=web", containerName),
	// 			"-l", fmt.Sprintf("traefik.http.routers.%s-http.middlewares=%s-https-redirect", containerName, containerName),
	//
	// 			"-l", fmt.Sprintf("traefik.http.middlewares.%s-https-redirect.redirectscheme.scheme=https", containerName),
	// 		)
	// 	} else {
	// 		runArgs = append(runArgs,
	// 			"-p", fmt.Sprintf("%d:%d", Port, Port),
	// 		)
	// 	}
	//
	// case models.AppTypeService:
	// 	runArgs = append(runArgs, "--network", "traefik-net")
	//
	// case models.AppTypeDatabase:
	// 	runArgs = append(runArgs, "--network", "traefik-net")
	//
	// default:
	// 	runArgs = append(runArgs,
	// 		"-p", fmt.Sprintf("%d:%d", Port, Port),
	// 	)
	// }
	//
	// runArgs = append(runArgs, imageTag)
	//
	// // ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// // defer cancel()
	//
	// cmd := exec.CommandContext(ctx, "docker", runArgs...)
	// cmd.Stdout = logfile
	// cmd.Stderr = logfile
	//
	// if err := cmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker run timed out after 5 minutes")
	// 	}
	// 	exitCode := -1
	// 	if exitErr, ok := err.(*exec.ExitError); ok {
	// 		exitCode = exitErr.ExitCode()
	// 	}
	// 	return fmt.Errorf("docker run failed with exit code %d: %w", exitCode, err)
	// }
	//
	// return nil
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

	// legacy exec method
	//
	//
	// cmd := exec.Command("docker", "inspect", containerName, "--format", "{{json .}}")
	// output, err := cmd.Output()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to inspect container: %w", err)
	// }
	//
	// var inspectData struct {
	// 	State struct {
	// 		Status  string `json:"Status"`
	// 		Running bool   `json:"Running"`
	// 		Paused  bool   `json:"Paused"`
	// 		Health  *struct {
	// 			Status string `json:"Status"`
	// 		} `json:"Health"`
	// 	} `json:"State"`
	// 	Name string `json:"Name"`
	// }
	//
	// if err := json.Unmarshal(output, &inspectData); err != nil {
	// 	return nil, fmt.Errorf("failed to parse inspect output: %w", err)
	// }
	//
	// uptimeCmd := exec.Command("docker", "inspect", containerName, "--format", "{{.State.StartedAt}}")
	// uptimeOutput, err := uptimeCmd.Output()
	// uptime := "N/A"
	// if err == nil {
	// 	uptime = strings.TrimSpace(string(uptimeOutput))
	// }
	//
	// state := "stopped"
	// if inspectData.State.Running {
	// 	state = "running"
	// } else if inspectData.State.Status == "exited" {
	// 	state = "stopped"
	// } else {
	// 	state = inspectData.State.Status
	// }
	//
	// healthy := true
	// if inspectData.State.Health != nil {
	// 	healthy = inspectData.State.Health.Status == "healthy"
	// }
	//
	// return &ContainerStatus{
	// 	Name:    strings.TrimPrefix(inspectData.Name, "/"),
	// 	Status:  inspectData.State.Status,
	// 	State:   state,
	// 	Uptime:  uptime,
	// 	Healthy: healthy,
	// }, nil
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

	port, domains, envVars, err := GetDeploymentConfigForApp(app)
	if err != nil {
		return fmt.Errorf("failed to get deployment configuration: %w", err)
	}

	if err := StopRemoveContainer(containerName, nil); err != nil {
		return fmt.Errorf("failed to stop/remove container: %w", err)
	}

	if err := RunContainer(ctx, app, imageTag, containerName, domains, port, envVars, nil); err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	return nil

	// legacy exec method
	//
	//
	// cmd := exec.Command("docker", "inspect", containerName, "--format", "{{.Config.Image}}")
	// output, err := cmd.Output()
	// if err != nil {
	// 	return fmt.Errorf("failed to get container image: %w", err)
	// }
	// imageTag := strings.TrimSpace(string(output))
	//
	// port, domains, envVars, err := GetDeploymentConfigForApp(app)
	// if err != nil {
	// 	return fmt.Errorf("failed to get deployment configuration: %w", err)
	// }
	//
	// if err := StopRemoveContainer(containerName, nil); err != nil {
	// 	return fmt.Errorf("failed to stop/remove container: %w", err)
	// }
	//
	// if err := RunContainer(app, imageTag, containerName, domains, port, envVars, nil); err != nil {
	// 	return fmt.Errorf("failed to run container: %w", err)
	// }
	//
	// return nil
}
