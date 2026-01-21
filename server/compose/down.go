package compose

import (
	"os"
	"os/exec"
)

func ComposeDown(appContextPath string) {
	cmd := exec.Command("docker", "compose", "down")
	cmd.Dir = appContextPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

