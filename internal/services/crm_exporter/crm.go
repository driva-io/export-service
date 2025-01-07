package crm_exporter

import (
	"context"
	"export-service/internal/repositories/crm_company_repo"

	"github.com/gofiber/fiber/v2"
)

type Status string

const (
	Updated Status = "updated"
	Created Status = "created"
	Skipped Status = "skipped"
	Failed  Status = "failed"
)

type ObjectStatus struct {
	CrmId          any     `json:"crm_id"`
	Status         Status  `json:"status"`
	Message        string  `json:"message,omitempty"`
	DrivaContactId *string `json:"driva_contact_id,omitempty"`
}

type CreatedLead struct {
	Company  *ObjectStatus   `json:"company,omitempty"`
	Deal     *ObjectStatus   `json:"deal,omitempty"`
	Contacts *[]ObjectStatus `json:"contacts,omitempty"`
	Other    *[]ObjectStatus `json:"other,omitempty"`
}

type Crm interface {
	Authorize(ctx context.Context, companyName string) (any, error)
	Validate(c *fiber.Ctx, client any) bool
	Install(installData any) (any, error)
	OAuthCallback(c *fiber.Ctx, params ...any) (any, error)

	SendLead(client any, mappedStorageData map[string]any, correspondingRawData map[string]any, configs map[string]any, existingLead map[string]any) (CreatedLead, error)
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
