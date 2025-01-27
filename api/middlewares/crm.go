package middlewares

import (
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/services/crm_exporter"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateCrmMiddleware(co *crm_company_repo.PgCrmCompanyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		crm := c.Params("crm")

		crmService, exists := crm_exporter.GetCrm(crm, co)
		if !exists {
			return errors.New("Crm " + crm + "not implemented.")
		}

		c.Locals("crm_service", crmService)

		return c.Next()
	}
}

func AuthenticateCrmMiddleware(co *crm_company_repo.PgCrmCompanyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		companyName := c.Query("company")
		workspaceId := c.Query("workspace_id")
		crmService := c.Locals("crm_service").(crm_exporter.Crm)
		crm := c.Params("crm")

		if workspaceId == "" && companyName == "" {
			return errors.New("either workspace_id or company query parameter is required")
		}

		if workspaceId == "" {
			company, err := co.GetByCompanyName(ctx, ports.CrmGetByCompanyNameQueryParams{
				Crm:         crm,
				CompanyName: companyName,
			})

			if company.WorkspaceId.String == "" {
				if strings.HasSuffix(c.Path(), "/validate") {
					c.Locals("crmClient", nil)
					return c.JSON(fiber.Map{"valid": false})
				}
				return errors.New("workspace_id not found for the given company")
			}

			if err != nil {
				return err
			}

			workspaceId = company.WorkspaceId.String
		}

		log.Printf("Authenticating CRM for workspace: %v", workspaceId)
		crmClient, err := crmService.Authorize(ctx, workspaceId)
		if err != nil {
			if strings.HasSuffix(c.Path(), "/validate") {
				c.Locals("crmClient", nil)
				return c.JSON(fiber.Map{"valid": false})
			}
			return err
		}

		c.Locals("crmClient", crmClient)

		return c.Next()
	}
}
