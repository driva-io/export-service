package crm_company_repo

const getQuery = `
	select * from crm.company where crm = $1 and name = $2
	`

const addHubspotQuery = `
	insert into crm.company (crm, name, user_who_installed_id, workspace_id, refresh_token, access_token, created_at, updated_at) values ('hubspot', $1, $2, $3, $4, $5, now(), now()) returning *;
`
