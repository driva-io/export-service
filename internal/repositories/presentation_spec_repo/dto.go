package presentation_spec_repo

import "time"

type PresentationSpecSheetOptions struct {
	PresentationSpecID string   `json:"presentation_spec_id"`
	Key                string   `json:"key"`
	ActiveColumns      []string `json:"active_columns"`
	Position           int      `json:"position"`
	ShouldExplode      bool     `json:"should_explode"`
}
type PresentationSpec struct {
	PresentationSpecID string         `json:"presentation_spec_id"`
	Key                string         `json:"key"`
	Value              map[string]any `json:"value"`
}

type PresentationSpecBasicInfo struct {
	ID          string    `json:"id"`
	Version     int       `json:"version"`
	Base        string    `json:"base"`
	UserEmail   string    `json:"user_email"`
	UserCompany string    `json:"user_company"`
	Service     string    `json:"service"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
