package routes

import (
	"export-service/api/middlewares"
	"export-service/internal/gateways"
	"export-service/internal/handlers"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/server"
	"export-service/internal/services/crm_exporter"

	"github.com/gofiber/fiber/v2"
)

func RegisterCrmRoutes(s *server.FiberServer, a gateways.AuthServiceGateway, co *crm_company_repo.PgCrmCompanyRepository, p *presentation_spec_repo.PgPresentationSpecRepository) {
	noAuthRoutes := s.App.Group("/crm/v1")
	noAuthRoutes.Use("/:crm/*", middlewares.ValidateCrmMiddleware(co))
	noAuthRoutes.Get("/:crm/oauth_callback", func(c *fiber.Ctx) error {
		return handlers.OAuthCallBackHandler(c, c.Locals("crm").(crm_exporter.Crm))
	})

	noCrmAuthRoutes := s.App.Group("/crm/v1")
	noCrmAuthRoutes.Use(middlewares.AuthMiddleware(a))
	noCrmAuthRoutes.Use("/:crm/*", middlewares.ValidateCrmMiddleware(co))
	noCrmAuthRoutes.Post("/:crm/install", func(c *fiber.Ctx) error {
		return handlers.InstallHandler(c, c.Locals("crm").(crm_exporter.Crm))
	})

	crmRoutes := s.App.Group("/crm/v1")
	crmRoutes.Use(middlewares.AuthMiddleware(a))
	crmRoutes.Use("/:crm/*", middlewares.ValidateCrmMiddleware(co))
	crmRoutes.Use("/:crm/*", middlewares.AuthenticateCrmMiddleware(co))
	crmRoutes.Get("/:crm/pipelines", func(c *fiber.Ctx) error {
		return handlers.GetPipelinesHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient"))
	})
	crmRoutes.Get("/:crm/fields", func(c *fiber.Ctx) error {
		return handlers.GetFieldsHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient"))
	})
	crmRoutes.Get("/:crm/owners", func(c *fiber.Ctx) error {
		return handlers.GetOwnersHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient"))
	})
	crmRoutes.Get("/:crm/validate", func(c *fiber.Ctx) error {
		return handlers.ValidateHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient"))
	})
	crmRoutes.Post("/:crm/test-lead", func(c *fiber.Ctx) error {
		return handlers.TestLeadHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient"), p)
	})
}
