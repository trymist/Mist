package compose

import (
	"os"
	"os/exec"
)

func ComposeRestart(appContextPath string) {
	cmd := exec.Command("docker", "compose", "restart")
	cmd.Dir = appContextPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
