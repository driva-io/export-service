package crm_exporter

import (
	"context"
	"export-service/internal/repositories/crm_company_repo"

	"github.com/gofiber/fiber/v2"
)

type Status string

// Define constants representing the enum values
const (
	Updated Status = "updated"
	Created Status = "created"
	Skipped Status = "skipped"
	Failed  Status = "failed"
)

type ObjectStatus struct {
	Id      any
	Status  Status
	Message string
}

type CreatedLead struct {
	company  *ObjectStatus
	deal     *ObjectStatus
	contacts *[]ObjectStatus
	other    *[]ObjectStatus
}

type Crm interface {
	Authorize(ctx context.Context, companyName string) (any, error)
	Validate(c *fiber.Ctx, client any) (bool)
	Install(installData any) (any, error)
	OAuthCallback(c *fiber.Ctx, params ...any) (any, error)

	SendLead(client any, mappedStorageData map[string]any) (CreatedLead, error)
	GetPipelines(client any) ([]Pipeline, error)
	GetFields(client any) (CrmFields, error)
	GetOwners(client any) ([]Owner, error)
}

func GetCrm(crm string, co *crm_company_repo.PgCrmCompanyRepository) (Crm, bool) {
	crms := map[string]Crm{
		"hubspot": NewHubspotService(co),
	}

	crmService, exists := crms[crm]
	return crmService, exists
}
