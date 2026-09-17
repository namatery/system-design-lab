package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/analytics"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/link"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/redirect"
	"go.uber.org/zap"
)

type Application struct {
	db     *gocql.Session
	logger *zap.Logger
	addr   string
}

func (app *Application) mount() {
	router := gin.Default()

	ctx, cancel := context.WithCancel(context.Background())
	service := analytics.NewAnalyticsService(app.db)
	collector := analytics.NewAnalyticsCollector(1024, service, app.logger)
	collector.Start(ctx, 4)
	defer func() { cancel(); collector.Wait() }()

	// Links
	{
		service := link.NewLinkService(app.db)
		handler := link.NewLinkController(service)

		v1 := router.Group("/v1/link")
		v1.POST("", handler.CreateShortenLink)
	}

	// Redirect
	{
		service := redirect.NewRedirectService(app.db)
		handler := redirect.NewRedirectControllerWithAnalytics(service, collector)

		v1 := router.Group("/")
		v1.GET(":alias", handler.Redirect)
	}

	server := &http.Server{
		Addr:              app.addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()

	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("HTTP server failed", zap.Error(err))
		}
	case <-stopCtx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			app.logger.Warn("HTTP shutdown failed", zap.Error(err))
			server.Close()
		}
	}
}
