package presentation_spec_repo

const getQuery = `
	with spec as (
		select presentation_spec_id as id, jsonb_object_agg(key, value) as spec from presentation_spec.specs
		group by presentation_spec_id
	),
	options as (
		select presentation_spec_id as id, jsonb_agg(jsonb_build_object('key', key,'active_columns', active_columns , 'position', position, 'should_explode', should_explode)) as sheet_options from presentation_spec.sheet_options
		group by presentation_spec_id
	),
	basic_info_with_default as (
		select * from presentation_spec.basic_info
		where (user_email = $1 AND user_company = $2 AND service = $3 AND base = $4) or (service = $3 AND base = $4 AND is_default)
		order by is_default -- first not default
		limit 1
	)
	SELECT * FROM basic_info_with_default inner join spec using (id) inner join options using (id)
	`
