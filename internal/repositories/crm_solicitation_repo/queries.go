package crm_solicitation_repo

const getQuery = `
	select * from crm.solicitation_v2 where list_id = $1 and crm = $2
	`

const updateStatusQuery = `
	update crm.solicitation_v2 set status = $1 where list_id = $2 and crm = $3 returning *
`

const incrementCurrentQuery = `
	update crm.solicitation_v2 set current = current + 1 where list_id = $1 and crm = $2 returning *
`

const updateExportedCompanies = `
	UPDATE crm.solicitation_v2
	SET exported_companies = jsonb_set(
		COALESCE(exported_companies, '{}'),
		ARRAY[$1],
		$2::jsonb
	)
	WHERE list_id = $3 and crm = $4
	RETURNING *;
	`

const createSolicitationQuery = `
insert into crm.solicitation_v2 (list_id, user_email, status, exported_companies, owner_id, stage_id, pipeline_id, overwrite_data, create_deal, current, total, created_at, updated_at, crm) values ($1, $2, 'In Progress', null, $3, $4, $5, $6, $7, $8, $9, now(), now(), $10) returning *
`
