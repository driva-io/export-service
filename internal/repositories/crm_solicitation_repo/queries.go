package crm_solicitation_repo

const getQuery = `
	select * from crm.solicitation_v2 where list_id = $1
	`

const updateExportedCompanies = `
	UPDATE crm.solicitation_v2
	SET exported_companies = jsonb_set(
		COALESCE(exported_companies, '{}'),
		ARRAY[$1],
		$2::jsonb
	)
	WHERE list_id = $3
	RETURNING *;
	`

const createSolicitation = `
insert into crm.solicitation_v2 (list_id, user_email, status, exported_companies, owner_id, stage_id, pipeline_id, overwrite_data, create_deal, current, total, created_at, updated_at) values ($1, $2, 'In Progress', null, $3, $4, $5, $6, $7, $8, $9, now(), now()) returning *
`
