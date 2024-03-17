package main

import (
	"embed"
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/jet/v2"
	"net/http"
)

//go:embed views/*
var viewsFS embed.FS

//go:embed public/*
var publicFS embed.FS

func main() {
	config.ShowAsciiArt()
	cfg, err := config.LoadConfig()

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
			"Title": fmt.Sprintf("Dashboard - %s", cfg.Node),
		})
	})

	data := app.Group("/data")
	data.Get("/logo", services.GetLogo(cfg))
	data.Get("/cpu", services.GetCPUUsage)
	data.Get("/ram", services.GetRAMUsage)
	data.Get("/uptime", services.GetSystemUptime)
	data.Get("/disk", services.GetDiskUsage)
	data.Get("/activity", services.GetActivity(cfg))

	actions := app.Group("/actions")
	actions.Get("/restart-node", services.RestartNode(cfg))
	actions.Get("/restart-server", services.RestartServer)
	actions.Get("/shutdown-server", services.ShutdownServer)

	err = app.Listen(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		panic(err)
	}
}
