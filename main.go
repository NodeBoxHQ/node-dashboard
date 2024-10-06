package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodeboxhq/nodebox-dashboard/internal/cmd"
	"github.com/nodeboxhq/nodebox-dashboard/internal/config"
	"github.com/nodeboxhq/nodebox-dashboard/internal/db"
	"github.com/nodeboxhq/nodebox-dashboard/internal/handlers"
	"github.com/nodeboxhq/nodebox-dashboard/internal/handlers/middleware"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services"
)

//go:embed web/build/*
var webFS embed.FS

func main() {
	handlers.EmbeddedWebFS = webFS

	cmd.AsciiArt()
	cfg := config.ParseConfig(cmd.ParseFlags())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceRegistry := services.NewServiceRegistry(db.SetupDatabase(cfg))
	sS := serviceRegistry.StatsService
	hS := serviceRegistry.HostService
	nS := serviceRegistry.NodeService

	hS.InstallService()

	go sS.StartStatsCollection(ctx)
	go hS.StartHostInfoCollection(ctx)
	go hS.StartUpdateChecker(ctx)
	go nS.GetNodeInfo()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()
	r.Use(middleware.CORSConfig())
	handlers.RegisterRoutes(r, cfg.Environment, sS, hS, nS)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.IP, cfg.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal().Any("error", err).Msg("Failed to start server")
		}
	}()

	logger.L.Info().Msgf("Server started on %s:%d", cfg.IP, cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L.Info().Msg("Shutting down gracefully...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.L.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.L.Info().Msg("Server exited properly")
}
