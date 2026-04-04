package websockets

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/gorilla/websocket"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

type ContainerLogsEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var containerLogsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOriginWithSettings,
}

func ContainerLogsHandler(w http.ResponseWriter, r *http.Request) {
	appIDStr := r.URL.Query().Get("appId")
	if appIDStr == "" {
		http.Error(w, "appId is required", http.StatusBadRequest)
		return
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid appId", http.StatusBadRequest)
		return
	}

	app, err := models.GetApplicationByID(appID)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	conn, err := containerLogsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection for container logs")
		return
	}
	defer conn.Close()

	log.Info().Int64("app_id", appID).Str("app_name", app.Name).Msg("Container logs client connected")

	containerName := docker.GetContainerName(app.Name, appID)

	if app.AppType != models.AppTypeCompose {
		if !docker.ContainerExists(containerName) {
			conn.WriteJSON(ContainerLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": "Container not found",
				},
			})
			return
		}

		status, err := docker.GetContainerStatus(containerName)
		if err != nil {
			conn.WriteJSON(ContainerLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": fmt.Sprintf("Failed to get container status: %v", err),
				},
			})
			return
		}

		conn.WriteJSON(ContainerLogsEvent{
			Type:      "status",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"container": containerName,
				"state":     status.State,
				"status":    status.Status,
			},
		})

		if status.State != "running" {
			conn.WriteJSON(ContainerLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": fmt.Sprintf("Container is not running (state: %s)", status.State),
				},
			})
			return
		}
	} else {
		conn.WriteJSON(ContainerLogsEvent{
			Type:      "status",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"container": app.Name,
				"state":     "running",
				"status":    "Compose Stack",
			},
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type logMessage struct {
		line       string
		streamType string
	}

	logChan := make(chan logMessage, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(logChan)

		if app.AppType == models.AppTypeCompose {
			// Streaming for Compose Apps
			path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
			cmd := exec.CommandContext(ctx, "docker", "compose", "logs", "--follow", "--tail", "100")
			cmd.Dir = path

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				errChan <- fmt.Errorf("failed to get stdout pipe: %w", err)
				return
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				errChan <- fmt.Errorf("failed to get stderr pipe: %w", err)
				return
			}

			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to start compose logs: %w", err)
				return
			}

			readStream := func(r io.Reader, streamType string) {
				scanner := bufio.NewScanner(r)
				for scanner.Scan() {
					select {
					case <-ctx.Done():
						return
					case logChan <- logMessage{line: scanner.Text(), streamType: streamType}:
					}
				}
				if err := scanner.Err(); err != nil && err != io.EOF {
				}
			}

			go readStream(stdout, "stdout")
			go readStream(stderr, "stderr")

			if err := cmd.Wait(); err != nil {
				if ctx.Err() == nil {
				}
			}
			return
		}

		// Standard Docker Container Streaming
		cli, err := client.New(client.FromEnv)
		if err != nil {
			errChan <- fmt.Errorf("failed to create docker client: %w", err)
			return
		}
		defer cli.Close()

		options := client.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Tail:       "100",
			Timestamps: false,
		}

		logReader, err := cli.ContainerLogs(ctx, containerName, options)
		if err != nil {
			errChan <- fmt.Errorf("failed to get container logs: %w", err)
			return
		}
		defer logReader.Close()

		buf := make([]byte, 8)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, err := io.ReadFull(logReader, buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("failed to read log header: %w", err)
				return
			}

			streamType := "stdout"
			if buf[0] == 2 {
				streamType = "stderr"
			}

			payloadSize := binary.BigEndian.Uint32(buf[4:8])

			payload := make([]byte, payloadSize)
			_, err = io.ReadFull(logReader, payload)
			if err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("failed to read log payload: %w", err)
				return
			}

			lines := string(payload)
			currentLine := ""
			for i := 0; i < len(lines); i++ {
				if lines[i] == '\n' {
					if currentLine != "" {
						select {
						case <-ctx.Done():
							return
						case logChan <- logMessage{line: currentLine, streamType: streamType}:
						}
					}
					currentLine = ""
				} else {
					currentLine += string(lines[i])
				}
			}
			if currentLine != "" {
				select {
				case <-ctx.Done():
					return
				case logChan <- logMessage{line: currentLine, streamType: streamType}:
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					cancel()
					return
				}
			}
		}
	}()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Info().Int64("app_id", appID).Msg("Container logs client disconnected")
				}
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int64("app_id", appID).Msg("Container logs context cancelled")
			return

		case err := <-errChan:
			log.Error().Err(err).Int64("app_id", appID).Msg("Container logs stream error")
			conn.WriteJSON(ContainerLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
			return

		case msg, ok := <-logChan:
			if !ok {
				conn.WriteJSON(ContainerLogsEvent{
					Type:      "end",
					Timestamp: time.Now().Format(time.RFC3339),
					Data: map[string]interface{}{
						"message": "Log stream ended",
					},
				})
				return
			}

			event := ContainerLogsEvent{
				Type:      "log",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"line":   msg.line,
					"stream": msg.streamType,
				},
			}

			if err := conn.WriteJSON(event); err != nil {
				log.Warn().Err(err).Int64("app_id", appID).Msg("Failed to send container log message to client")
				return
			}
		}
	}
}
