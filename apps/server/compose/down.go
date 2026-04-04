package compose

import (
	"os"
	"os/exec"
)

func ComposeDown(appContextPath string) error {
	cmd := exec.Command("docker", "compose", "down")
	cmd.Dir = appContextPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
