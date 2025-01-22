package middlewares

import (
	"errors"
	"export-service/internal/core/ports"
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

		var company crm_company_repo.Company
		var err error
		if workspaceId != "" {
			log.Printf("Authenticating CRM for workspace: %v", workspaceId)
			company, err = co.GetCompanyByWorkspaceId(ctx, ports.CrmCompanyQueryParams{Crm: crm, WorkspaceId: workspaceId})
			if err != nil {
				return err
			}
		} else if companyName != "" {
			log.Printf("Authenticating CRM for company: %v", companyName)
			company, err = co.GetByCompanyName(ctx, ports.CrmGetByCompanyNameQueryParams{Crm: crm, CompanyName: companyName})
			if err != nil {
				return err
			}
		} else {
			return errors.New("either workspace_id or company query parameter is required")
		}

		crmClient, err := crmService.Authorize(ctx, company.WorkspaceId.String)
		if err != nil {
			return err
		}

		c.Locals("crmClient", crmClient)

		return c.Next()
	}
}
