package crm_solicitation_repo

import (
	"export-service/internal/services/crm_exporter"
	"time"
)

type UpdateExportedCompaniesParms struct {
	Cnpj               any                      `json:"cnpj"` //Any to support changes to cnpj type
	NewExportedCompany crm_exporter.CreatedLead `json:"new_exported_company"`
}

type SolicitationStatus string

// Define constants representing the enum values
const (
	Interrupted SolicitationStatus = "Interrupted"
	InProgress  SolicitationStatus = "In Progress"
	Completed   SolicitationStatus = "Completed"
)

type Solicitation struct {
	ListId            string
	UserEmail         string
	Status            SolicitationStatus
	ExportedCompanies map[string]any

	OwnerId       string
	PipelineId    string
	StageId       string
	OverwriteData bool
	CreateDeal    bool
	Current       int
	Total         int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateSolicitation struct {
	ListId    string
	UserEmail string
	Current   int
	Total     int

	//Potentialy make these fields optional for future CRMs
	OwnerId       string
	PipelineId    string
	StageId       string
	OverwriteData bool
	CreateDeal    bool
}
