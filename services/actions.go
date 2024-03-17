package services

import (
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/gofiber/fiber/v2"
	"os/exec"
)

func RestartNode(cfg *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if cfg.Node == "linea" {
			err := exec.Command("systemctl", "stop", "linea").Run()
			if err != nil {
				return c.SendString("Error stopping linea service")
			} else {
				return c.SendString("Linea service stopped")
			}
		} else if cfg.Node == "dusk" {
			err := exec.Command("service", "rusk", "stop").Run()
			if err != nil {
				return c.SendString("Error stopping rusk service")
			} else {
				return c.SendString("Rusk service stopped")
			}
		}
		return c.SendString("Unknown node")
	}
}

func RestartServer(c *fiber.Ctx) error {
	err := exec.Command("reboot").Run()
	if err != nil {
		return c.SendString("Error rebooting server")
	} else {
		return c.SendString("Server rebooting")
	}
}

func ShutdownServer(c *fiber.Ctx) error {
	err := exec.Command("shutdown", "now").Run()
	if err != nil {
		return c.SendString("Error shutting down server")
	} else {
		return c.SendString("Server shutting down")
	}
}
