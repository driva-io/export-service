package middlewares

import (
	"export-service/internal/gateways"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(g gateways.AuthServiceGateway) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := make(map[string]any)
		authHeader := c.Get(fiber.HeaderAuthorization)

		if authHeader != "" {
			headers["Authorization"] = strings.Replace(authHeader, "Bearer ", "", 1)
		} else {
			authHeader = c.Cookies("session")
			headers["Cookie"] = "session=" + c.Cookies("session")
		}

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header or session cookie",
			})
		}

		user, err := g.GetUserByToken(headers)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func TokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-API-KEY")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-API-KEY header",
			})
		}
		if token != os.Getenv("API_KEY") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid X-API-KEY header",
			})
		}
		return c.Next()
	}
}
