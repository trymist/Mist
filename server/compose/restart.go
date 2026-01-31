package compose

import (
	"os"
	"os/exec"
)

func ComposeRestart(appContextPath string) error {
	cmd := exec.Command("docker", "compose", "restart")
	cmd.Dir = appContextPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
