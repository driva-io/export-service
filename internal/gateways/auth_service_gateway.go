package gateways

type AuthServiceGateway interface {
	Login(email string, password string, company string) (response map[string]interface{}, err error)
	HasCredits(companyId string, amount int) (hasCredits bool, err error)
	TakeCredits()
	RefundCredits()
	GetBlacklist()
	GetAdminToken() (adminToken string)
	GetUserByToken(headers map[string]any) (AuthUser, error)
}

type AuthUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	WorkspaceID string `json:"workspaceId"`
	OwnerId     string `json:"ownerId"`
}
