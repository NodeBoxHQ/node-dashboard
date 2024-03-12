package main

import (
	"github.com/NodeboxHQ/node-dashboard/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	metrics := app.Group("/metrics")
	metrics.Get("/cpu", services.GetCPUUsage)
	metrics.Get("/ram", services.GetRAMUsage)
	metrics.Get("/uptime", services.GetSystemUptime)
	metrics.Get("/disk", services.GetDiskUsage)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
