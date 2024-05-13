package main

import (
	"embed"
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/NodeboxHQ/node-dashboard/services/xally"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"github.com/NodeboxHQ/node-dashboard/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/jet/v2"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed views/*
var viewsFS embed.FS

//go:embed public/*
var publicFS embed.FS

func main() {
	config.ShowAsciiArt()
	logger.InitLogger()
	logger.Info("Initializing NodeBox Dashboard")

	cfg, err := config.LoadConfig()

	logger.Info("Loaded config: ", cfg)

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	engine := jet.NewFileSystem(http.FS(viewsFS), ".jet")
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(publicFS),
		PathPrefix: "public",
		Browse:     true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index.jet", fiber.Map{
			"Title":    fmt.Sprintf("Dashboard - %s", cfg.Node),
			"NodeIP":   cfg.IPv4,
			"NodeType": cfg.Node,
		})
	})

	data := app.Group("/data")
	data.Get("/logo", services.GetLogo(cfg))
	data.Get("/cpu", services.GetCPUUsage(cfg.IPv4))
	data.Get("/ram", services.GetRAMUsage(cfg.IPv4))
	data.Get("/uptime", services.GetSystemUptime(cfg.IPv4))
	data.Get("/disk", services.GetDiskUsage(cfg.IPv4))
	data.Get("/activity", services.GetActivity(cfg))

	actions := app.Group("/actions")
	actions.Get("/restart-node", services.RestartNode(cfg))
	actions.Get("/restart-server", services.RestartServer)
	//actions.Get("/shutdown-server", services.ShutdownServer)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		checkPortOpen := func(port int) bool {
			timeout := time.Second
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), timeout)
			if err != nil {
				return false
			}
			if conn != nil {
				conn.Close()
				return true
			}
			return false
		}

		for !checkPortOpen(cfg.Port) {
			time.Sleep(time.Second)
		}

		utils.InstallService()
	}()

	updateTicker := time.NewTicker(10 * time.Minute)

	go func() {
		for {
			select {
			case <-updateTicker.C:
				utils.SelfUpdate(cfg.NodeboxDashboardVersion)
			}
		}
	}()

	if cfg.Node == "Xally" {
		webhookTicker := time.NewTicker(6 * time.Minute)

		go func() {
			for {
				select {
				case <-webhookTicker.C:
					xally.CheckRunning(cfg)
				}
			}
		}()
	}

	<-stopChan
}
