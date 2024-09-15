package domain

import "time"

type PresentationSpecSheetOptions struct {
	Key           string   `json:"key"`
	ActiveColumns []string `json:"active_columns"`
	Position      int      `json:"position"`
	ShouldExplode bool     `json:"should_explode"`
}
type PresentationSpecSpec struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}

type PresentationSpec struct {
	ID           string                         `json:"id" binding:"required"`
	Version      int                            `json:"version"  binding:"required"`
	Base         string                         `json:"base"  binding:"required"`
	UserEmail    string                         `json:"user_email"`
	UserCompany  string                         `json:"user_company"`
	Service      string                         `json:"service"  binding:"required"`
	SheetOptions []PresentationSpecSheetOptions `json:"sheet_options"`
	Spec         []PresentationSpecSpec         `json:"spec"  binding:"required"`
	CreatedAt    time.Time                      `json:"created_at"  binding:"required"`
	UpdatedAt    time.Time                      `json:"updated_at"  binding:"required"`
	IsDefault    bool                           `json:"is_default"  binding:"required"`
}
