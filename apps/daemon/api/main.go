package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	// "github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/websockets"
	"github.com/rs/zerolog/log"
)

func InitApiServer() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	staticDir := "static"
	fs := http.FileServer(http.Dir(staticDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	go websockets.BroadcastMetrics()
	handler := mux
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Info().Msg("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
