package compose

import (
	"fmt"
	"os"
	"os/exec"
)

func ComposeUp(appContextPath string, env map[string]string, logFile *os.File) {
	cmd := exec.Command("docker", "compose", "up", "-d")
	var envArray []string
	for k, v := range env {
		envArray = append(envArray, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = envArray
	cmd.Dir = appContextPath
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Run()
}
