package crm_company_repo

import (
	"database/sql"
	"time"
)

type Company struct {
	Id        string
	Crm       string
	CreatedAt time.Time
	UpdatedAt time.Time

	CrmId              sql.NullString
	RefreshedAt        sql.NullString
	Name               sql.NullString
	RefreshToken       sql.NullString
	AccessToken        sql.NullString
	ExpiresIn          sql.NullString
	Environment        sql.NullString
	Token              sql.NullString
	Webhook            sql.NullString
	Email              sql.NullString
	Password           sql.NullString
	InstanceUrl        sql.NullString
	Merge              sql.NullString
	Mapping            sql.NullString
	MappingLinkedin    sql.NullString
	CompanyId          sql.NullString
	UserWhoInstalledId sql.NullString
	WorkspaceId        sql.NullString
}
