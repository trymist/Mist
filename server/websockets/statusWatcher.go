package websockets

import (
	"context"
	"strings"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

type DeploymentEvent struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type StatusUpdate struct {
	DeploymentID int64  `json:"deployment_id"`
	Status       string `json:"status"`
	Stage        string `json:"stage"`
	Progress     int    `json:"progress"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type LogUpdate struct {
	Line      string    `json:"line"`
	Stream    string    `json:"stream,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func DetectStreamType(line string) string {
	lineLower := strings.ToLower(line)

	stderrPatterns := []string{
		"error:",
		"err:",
		"fatal:",
		"panic:",
		"warning:",
		"warn:",
		"failed",
		"failure",
		"exception:",
		"traceback",
		"stack trace",
		" err ",
		"[error]",
		"[err]",
		"[fatal]",
		"[panic]",
		"[warning]",
		"[warn]",
	}

	for _, pattern := range stderrPatterns {
		if strings.Contains(lineLower, pattern) {
			return "stderr"
		}
	}

	return "stdout"
}

func WatchDeploymentStatus(ctx context.Context, depID int64, events chan<- DeploymentEvent) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	defer close(events)

	var lastStatus models.DeploymentStatus
	var lastStage string
	var lastProgress int

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := models.GetDeploymentByID(depID)
			if err != nil {
				continue
			}

			if dep.Status != lastStatus || dep.Stage != lastStage || dep.Progress != lastProgress {
				lastStatus = dep.Status
				lastStage = dep.Stage
				lastProgress = dep.Progress

				errMsg := ""
				if dep.ErrorMessage != nil {
					errMsg = *dep.ErrorMessage
				}

				select {
				case <-ctx.Done():
					return
				case events <- DeploymentEvent{
					Type:      "status",
					Timestamp: time.Now(),
					Data: StatusUpdate{
						DeploymentID: depID,
						Status:       string(dep.Status),
						Stage:        dep.Stage,
						Progress:     dep.Progress,
						Message:      utils.GetStageMessage(dep.Stage),
						ErrorMessage: errMsg,
					},
				}:
				}
			}

			if dep.Status == "success" || dep.Status == "failed" {
				time.Sleep(1 * time.Second)
				return
			}
		}
	}
}
