package websockets

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"
)

type DockerLogEntry struct {
	Stream         string                 `json:"stream"`
	Aux            map[string]interface{} `json:"aux,omitempty"`
	Error          string                 `json:"error,omitempty"`
	ErrorDetail    map[string]interface{} `json:"errorDetail,omitempty"`
	Status         string                 `json:"status,omitempty"`
	ID             string                 `json:"id,omitempty"`
	Progress       string                 `json:"progress,omitempty"`
	ProgressDetail map[string]interface{} `json:"progressDetail,omitempty"`
}

func WatcherLogs(ctx context.Context, filePath string, send chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			} else if err != nil {
				return err
			}

			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}

			if len(line) > 0 {
				processedLine := parseDockerLog(line)
				if processedLine != "" {
					send <- processedLine
				}
			}
		}
	}
}

func parseDockerLog(line string) string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "{") {
		return line
	}

	var dockerLog DockerLogEntry
	if err := json.Unmarshal([]byte(trimmed), &dockerLog); err != nil {
		return line
	}

	if dockerLog.Error != "" {
		return dockerLog.Error
	}

	if dockerLog.Stream != "" {
		return dockerLog.Stream
	}

	if dockerLog.Status != "" {
		switch dockerLog.Status {
		case "Downloading", "Extracting", "Waiting", "Verifying Checksum":
			return ""
		case "Pull complete", "Download complete", "Already exists":
			if dockerLog.ID != "" {
				return dockerLog.Status + ": " + dockerLog.ID
			}
			return dockerLog.Status
		default:
			if dockerLog.ID != "" {
				return dockerLog.Status + ": " + dockerLog.ID
			}
			return dockerLog.Status
		}
	}

	if dockerLog.Aux != nil {
		return ""
	}

	return line
}
