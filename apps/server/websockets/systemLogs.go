package websockets

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type SystemLogsEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var systemLogsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOriginWithSettings,
}

func SystemLogsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := systemLogsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection for system logs")
		return
	}
	defer conn.Close()

	log.Info().Msg("System logs client connected")

	conn.WriteJSON(SystemLogsEvent{
		Type:      "connected",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"message": "Connected to system logs",
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logChan := make(chan string, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(logChan)
		cmd := exec.CommandContext(ctx, "journalctl", "-u", "mist", "-f", "-n", "100", "--no-pager")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errChan <- fmt.Errorf("failed to create stdout pipe: %w", err)
			return
		}

		if err := cmd.Start(); err != nil {
			errChan <- fmt.Errorf("failed to start journalctl: %w", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				return
			case logChan <- line:
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("scanner error: %w", err)
		}

		cmd.Wait()
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
					log.Info().Msg("System logs client disconnected")
				}
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("System logs context cancelled")
			return

		case err := <-errChan:
			log.Error().Err(err).Msg("System logs stream error")
			conn.WriteJSON(SystemLogsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
			return

		case line, ok := <-logChan:
			if !ok {
				conn.WriteJSON(SystemLogsEvent{
					Type:      "end",
					Timestamp: time.Now().Format(time.RFC3339),
					Data: map[string]interface{}{
						"message": "Log stream ended",
					},
				})
				return
			}

			event := SystemLogsEvent{
				Type:      "log",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"line": line,
				},
			}

			if err := conn.WriteJSON(event); err != nil {
				log.Warn().Err(err).Msg("Failed to send system log message to client")
				return
			}
		}
	}
}
