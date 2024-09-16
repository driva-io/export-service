package usecases

type ExportRequest struct {
	UserEmail       string `json:"user_email"`
	UserCompany     string `json:"user_company"`
	DataSource      string `json:"source"`
	DataDownloadURL string `json:"download_url"`
}
