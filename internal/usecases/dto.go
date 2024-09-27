package usecases

type ExportRequest struct {
	UserName        string `json:"user_name"`
	UserEmail       string `json:"user_email"`
	UserCompany     string `json:"user_company"`
	DataSource      string `json:"source"`
	DataDownloadURL string `json:"download_url"`
	ListID          string `json:"list_id"`
}
