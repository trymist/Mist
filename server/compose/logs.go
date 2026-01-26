package compose

import (
	"fmt"
	"os/exec"
)

func GetComposeLogs(appContextPath string, tail int) (string, error) {
	tailStr := fmt.Sprintf("%d", tail)
	cmd := exec.Command("docker", "compose", "logs", "--tail", tailStr)
	cmd.Dir = appContextPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get compose logs: %w", err)
	}

	return string(output), nil
}
