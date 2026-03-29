package websockets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/gorilla/websocket"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

type ContainerStatsEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type ContainerStatsData struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsed    float64 `json:"memory_used_mb"`
	MemoryLimit   float64 `json:"memory_limit_mb"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     float64 `json:"network_rx_mb"`
	NetworkTx     float64 `json:"network_tx_mb"`
	BlockRead     float64 `json:"block_read_mb"`
	BlockWrite    float64 `json:"block_write_mb"`
	PIDs          uint64  `json:"pids"`
}

var containerStatsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOriginWithSettings,
}

func ContainerStatsHandler(w http.ResponseWriter, r *http.Request) {
	appIDStr := r.URL.Query().Get("appId")
	if appIDStr == "" {
		http.Error(w, "appId is required", http.StatusBadRequest)
		return
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid appID", http.StatusBadRequest)
		return
	}

	app, err := models.GetApplicationByID(appID)
	if err != nil {
		http.Error(w, "application not found", http.StatusNotFound)
		return
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		http.Error(w, "unable to create docker client", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	conn, err := containerStatsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection for container stats")
		return
	}
	defer conn.Close()

	log.Info().Int64("app_id", appID).Str("app_name", app.Name).Msg("Container stats client connected")

	containerName := docker.GetContainerName(app.Name, appID)

	if !docker.ContainerExists(containerName) {
		conn.WriteJSON(ContainerStatsEvent{
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
		conn.WriteJSON(ContainerStatsEvent{
			Type:      "error",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"message": fmt.Sprintf("Failed to get container status: %v", err),
			},
		})
		return
	}

	conn.WriteJSON(ContainerStatsEvent{
		Type:      "status",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"container": containerName,
			"state":     status.State,
			"status":    status.Status,
		},
	})

	if status.State != "running" {
		conn.WriteJSON(ContainerStatsEvent{
			Type:      "error",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]interface{}{
				"message": fmt.Sprintf("Container is not running (state: %s)", status.State),
			},
		})
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	statsChan := make(chan *ContainerStatsData, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(statsChan)

		statsResult, err := cli.ContainerStats(ctx, containerName, client.ContainerStatsOptions{
			Stream: true,
		})
		if err != nil {
			errChan <- fmt.Errorf("failed to get container stats: %w", err)
			return
		}
		defer statsResult.Body.Close()

		decoder := json.NewDecoder(statsResult.Body)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var stats container.StatsResponse
				if err := decoder.Decode(&stats); err != nil {
					errChan <- fmt.Errorf("failed to decode stats: %w", err)
					return
				}

				if stats.PreCPUStats.SystemUsage == 0 {
					continue
				}

				cpu := calculateCPUPercent(&stats)
				memUsed, memLimit, memPercent := calculateMemory(&stats)
				netRx, netTx := calculateNetwork(&stats)
				blockRead, blockWrite := calculateBlockIO(&stats)

				statsData := &ContainerStatsData{
					CPUPercent:    cpu,
					MemoryUsed:    bytesToMB(memUsed),
					MemoryLimit:   bytesToMB(memLimit),
					MemoryPercent: memPercent,
					NetworkRx:     bytesToMB(netRx),
					NetworkTx:     bytesToMB(netTx),
					BlockRead:     bytesToMB(blockRead),
					BlockWrite:    bytesToMB(blockWrite),
					PIDs:          stats.PidsStats.Current,
				}

				select {
				case <-ctx.Done():
					return
				case statsChan <- statsData:
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
					log.Info().Int64("app_id", appID).Msg("Container stats client disconnected")
				}
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int64("app_id", appID).Msg("Container stats context cancelled")
			return

		case err := <-errChan:
			log.Error().Err(err).Int64("app_id", appID).Msg("Container stats stream error")
			conn.WriteJSON(ContainerStatsEvent{
				Type:      "error",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
			return

		case statsData, ok := <-statsChan:
			if !ok {
				conn.WriteJSON(ContainerStatsEvent{
					Type:      "end",
					Timestamp: time.Now().Format(time.RFC3339),
					Data: map[string]interface{}{
						"message": "Stats stream ended",
					},
				})
				return
			}

			event := ContainerStatsEvent{
				Type:      "stats",
				Timestamp: time.Now().Format(time.RFC3339),
				Data: map[string]interface{}{
					"cpu_percent":     statsData.CPUPercent,
					"memory_used_mb":  statsData.MemoryUsed,
					"memory_limit_mb": statsData.MemoryLimit,
					"memory_percent":  statsData.MemoryPercent,
					"network_rx_mb":   statsData.NetworkRx,
					"network_tx_mb":   statsData.NetworkTx,
					"block_read_mb":   statsData.BlockRead,
					"block_write_mb":  statsData.BlockWrite,
					"pids":            statsData.PIDs,
				},
			}

			if err := conn.WriteJSON(event); err != nil {
				log.Warn().Err(err).Int64("app_id", appID).Msg("Failed to send container stats message to client")
				return
			}
		}
	}
}

func calculateCPUPercent(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta <= 0 || cpuDelta <= 0 {
		return 0.0
	}

	onlineCPUs := float64(stats.CPUStats.OnlineCPUs)
	if onlineCPUs == 0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	return (cpuDelta / systemDelta) * onlineCPUs * 100.0
}

func calculateMemory(stats *container.StatsResponse) (used, limit uint64, percent float64) {
	used = stats.MemoryStats.Usage

	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		used -= cache
	}

	limit = stats.MemoryStats.Limit
	if limit > 0 {
		percent = (float64(used) / float64(limit)) * 100.0
	}

	return
}

func calculateNetwork(stats *container.StatsResponse) (rx, tx uint64) {
	for _, netStats := range stats.Networks {
		rx += netStats.RxBytes
		tx += netStats.TxBytes
	}
	return
}

func calculateBlockIO(stats *container.StatsResponse) (read, write uint64) {
	for _, bioEntry := range stats.BlkioStats.IoServiceBytesRecursive {
		switch bioEntry.Op {
		case "read", "Read":
			read += bioEntry.Value
		case "write", "Write":
			write += bioEntry.Value
		}
	}
	return
}

func bytesToMB(b uint64) float64 {
	return float64(b) / 1024 / 1024
}
