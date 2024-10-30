package crm_exporter

type Stage struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Pipeline struct {
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Stages []Stage `json:"stages,omitempty"`
}

type FieldOptions struct {
	Id    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type CrmField struct {
	Id       interface{}     `json:"id"` // Can be string or number
	Label    string          `json:"label"`
	Type     string          `json:"type"`
	Options  *[]FieldOptions `json:"options,omitempty"`
	Required *bool           `json:"required,omitempty"`
}

type CrmFields struct {
	Deals     *[]CrmField `json:"deals,omitempty"`
	Companies *[]CrmField `json:"companies,omitempty"`
	Contacts  *[]CrmField `json:"contacts,omitempty"`
}

type Owner struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
