package main

import (
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ShowAsciiArt()
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Static("/", "./public")

	metrics := app.Group("/metrics")
	metrics.Get("/cpu", services.GetCPUUsage)
	metrics.Get("/ram", services.GetRAMUsage)
	metrics.Get("/uptime", services.GetSystemUptime)
	metrics.Get("/disk", services.GetDiskUsage)
	metrics.Get("/activity", services.GetActivity(cfg))

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
