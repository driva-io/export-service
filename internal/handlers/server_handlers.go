package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]any{
		"data": map[string]any{
			"message": "Server running",
		},
	})
}
