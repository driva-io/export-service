package routes

import (
	"export-service/api/middlewares"
	"export-service/internal/gateways"
	"export-service/internal/handlers"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/server"

	"github.com/gofiber/fiber/v2"
)

func RegisterPresentationSpecRoutes(s *server.FiberServer, p *presentation_spec_repo.PgPresentationSpecRepository, a gateways.AuthServiceGateway) {
	presentationSpecRoutes := s.App.Group("/presentation-spec")

	presentationSpecRoutes.Use(middlewares.AuthMiddleware(a))

	presentationSpecRoutes.Get("/", func(c *fiber.Ctx) error {
		return handlers.GetPresentationSpecHandler(c, p)
	})
	presentationSpecRoutes.Post("/", func(c *fiber.Ctx) error {
		return handlers.AddPresentationSpecHandler(c, p)
	})
	presentationSpecRoutes.Patch("/:id", func(c *fiber.Ctx) error {
		return handlers.PatchPresentationSpecHandler(c, p)
	})
	presentationSpecRoutes.Patch("/:id/:key", func(c *fiber.Ctx) error {
		return handlers.PatchPresentationSpecKeyHandler(c, p)
	})
	presentationSpecRoutes.Delete("/:id", func(c *fiber.Ctx) error {
		return handlers.DeletePresentationSpecHandler(c, p)
	})
}
