package main

import (
	"github.com/namatery/system-design-lab/00-url-shortener/internal/shared/config"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/shared/db"
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Load config failed", zap.Error(err))
	}
	logger.Info("Configuration loaded successfully")

	// Connect to ScyllaDB
	session, err := db.New([]string{cfg.ScyllaAddr}, cfg.ScyllaKeyspace)
	if err != nil {
		logger.Fatal("Failed to connect to ScyllaDB", zap.Error(err))
	}
	defer session.Close()
	logger.Info("Connected to ScyllaDB successfully")

	app := &Application{db: session, addr: cfg.Addr}
	app.mount()
}
