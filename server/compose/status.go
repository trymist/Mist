package compose

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ComposeStatus struct {
	Name     string           `json:"name"`
	Status   string           `json:"status"`
	State    string           `json:"state"`
	Uptime   string           `json:"uptime"`
	Services []ComposeService `json:"services"`
}

type ComposeService struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Status string `json:"status"`
}

type DockerComposePSOutput struct {
	Name    string `json:"Name"`
	Service string `json:"Service"`
	State   string `json:"State"`
	Status  string `json:"Status"`
	Created int64  `json:"Created"`
}

func GetComposeStatus(appContextPath string, appName string) (*ComposeStatus, error) {
	cmd := exec.Command("docker", "compose", "ps", "--format", "json")
	cmd.Dir = appContextPath
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run docker compose ps: %w", err)
	}

	var services []DockerComposePSOutput

	if len(output) == 0 {
		return &ComposeStatus{
			Name:     appName,
			Status:   "Not Created",
			State:    "stopped",
			Uptime:   "N/A",
			Services: []ComposeService{},
		}, nil
	}

	if err := json.Unmarshal(output, &services); err != nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var svc DockerComposePSOutput
			if err := json.Unmarshal([]byte(line), &svc); err == nil {
				services = append(services, svc)
			}
		}
	}

	if len(services) == 0 {
		return &ComposeStatus{
			Name:     appName,
			Status:   "Stopped",
			State:    "stopped",
			Uptime:   "N/A",
			Services: []ComposeService{},
		}, nil
	}

	runningCount := 0
	formattedServices := []ComposeService{}

	for _, s := range services {
		formattedServices = append(formattedServices, ComposeService{
			Name:   s.Service, // Or s.Name
			State:  s.State,
			Status: s.Status,
		})
		if s.State == "running" {
			runningCount++
		}
	}

	overallState := "stopped"
	if runningCount == len(services) {
		overallState = "running"
	} else if runningCount > 0 {
		overallState = "partial"
	}

	return &ComposeStatus{
		Name:     appName,
		Status:   fmt.Sprintf("%d/%d Running", runningCount, len(services)),
		State:    overallState,
		Uptime:   time.Now().Format(time.RFC3339),
		Services: formattedServices,
	}, nil
}
