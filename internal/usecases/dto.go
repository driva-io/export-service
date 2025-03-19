package usecases

type ExportRequest struct {
	UserName        string `json:"user_name"`
	UserEmail       string `json:"user_email"`
	UserCompany     string `json:"user_company"`
	DataSource      string `json:"source"`
	DataDownloadURL string `json:"download_url"`
	ListID          string `json:"list_id"`
	ListName        string `json:"list_name"`
}

type CrmExportRequest struct {
	UserName        string `json:"user_name"`
	UserEmail       string `json:"user_email"`
	UserCompany     string `json:"user_company"`
	DataSource      string `json:"source"`
	DataDownloadURL string `json:"download_url"`
	ListID          string `json:"list_id"`
}

type CrmExportHeaders struct {
	Crm           string `json:"crm"`
	OwnerId       string `json:"owner_id"`
	PipelineId    string `json:"pipeline_id"`
	StageId       string `json:"stage_id"`
	OverwriteData bool   `json:"overwrite_data"`
	CreateDeal    bool   `json:"create_deal"`
}
