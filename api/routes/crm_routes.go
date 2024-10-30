package routes

import (
	"export-service/api/middlewares"
	"export-service/internal/gateways"
	"export-service/internal/handlers"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/server"
	"export-service/internal/services/crm_exporter"

	"github.com/gofiber/fiber/v2"

	"github.com/belong-inc/go-hubspot"
)

func RegisterCrmRoutes(s *server.FiberServer, a gateways.AuthServiceGateway, co *crm_company_repo.PgCrmCompanyRepository) {
	crmRoutes := s.App.Group("/crm/v1")

	crmRoutes.Use(middlewares.AuthMiddleware(a))
	crmRoutes.Use("/:crm/*", middlewares.ValidateCrmMiddleware(co))
	crmRoutes.Use("/:crm/*", middlewares.AuthenticateCrmMiddleware())

	crmRoutes.Get("/:crm/pipelines", func(c *fiber.Ctx) error {
		return handlers.GetPipelinesHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient").(*hubspot.Client))
	})
	crmRoutes.Get("/:crm/fields", func(c *fiber.Ctx) error {
		return handlers.GetFieldsHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient").(*hubspot.Client))
	})
	crmRoutes.Get("/:crm/owners", func(c *fiber.Ctx) error {
		return handlers.GetOwnersHandler(c, c.Locals("crm").(crm_exporter.Crm), c.Locals("crmClient").(*hubspot.Client))
	})
}
