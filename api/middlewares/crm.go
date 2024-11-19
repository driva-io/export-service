package middlewares

import (
	"errors"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/services/crm_exporter"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ValidateCrmMiddleware(co *crm_company_repo.PgCrmCompanyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		crm := c.Params("crm")

		crmService, exists := crm_exporter.GetCrm(crm, co)
		if !exists {
			return errors.New("Crm " + crm + "not implemented.")
		}

		c.Locals("crm", crmService)

		return c.Next()
	}
}

func AuthenticateCrmMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		companyName := c.Query("company")
		crm := c.Locals("crm").(crm_exporter.Crm)

		log.Printf("Authenticating CRM for company: %v", companyName)

		crmClient, err := crm.Authorize(ctx, companyName)
		if err != nil {
			return err
		}

		c.Locals("crmClient", crmClient)

		return c.Next()
	}
}
