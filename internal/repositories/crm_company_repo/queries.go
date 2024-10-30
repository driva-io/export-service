package crm_company_repo

const getQuery = `
	select * from crm.company where crm = $1 and name = $2
	`
