package presentation_spec_repo

const getQuery = `
	with spec as (
		select presentation_spec_id as id, jsonb_object_agg(key, value) as spec from presentation_spec.specs
		group by presentation_spec_id
	),
	options as (
		select presentation_spec_id as id, jsonb_agg(jsonb_build_object('key', key,'active_columns', active_columns , 'position', position, 'should_explode', should_explode) order by position) as sheet_options from presentation_spec.sheet_options
		group by presentation_spec_id
	),
	basic_info_with_default as (
		select * from presentation_spec.basic_info
		where (($1 = '' OR user_email = $1) AND user_company = $2 AND service = $3 AND base = $4) or (service = $3 AND base = $4 AND is_default)
		order by is_default -- first not default
		limit 1
	)
	SELECT * FROM basic_info_with_default inner join spec using (id) inner join options using (id)
	`

const getByIdQuery = `
with spec as (
	select presentation_spec_id as id, jsonb_object_agg(key, value) as spec from presentation_spec.specs
	group by presentation_spec_id
),
options as (
	select presentation_spec_id as id, jsonb_agg(jsonb_build_object('key', key,'active_columns', active_columns , 'position', position, 'should_explode', should_explode) order by position) as sheet_options from presentation_spec.sheet_options
	group by presentation_spec_id
),
basic_info_with_default as (
	select * from presentation_spec.basic_info
	where id = $1
	order by is_default -- first not default
	limit 1
)
SELECT * FROM basic_info_with_default inner join spec using (id) inner join options using (id)
`

const addBasicInfoQuery = `
	insert into presentation_spec.basic_info (id, base, user_email, user_company, service, is_default) values ($1, $2, $3, $4, $5, false);
`

const addOptionsQuery = `
	insert into presentation_spec.sheet_options (presentation_spec_id, key, active_columns, position, should_explode) values ($1, $2, $3, $4, $5);
`

const addSpecQuery = `
	insert into presentation_spec.specs (presentation_spec_id, key, value) values ($1, $2, $3);
`

const deleteQuery = `
	delete from presentation_spec.basic_info where id = $1;
`

const deleteSpecsQuery = `delete from presentation_spec.specs where presentation_spec_id = $1`

const deleteSheetOptionsQuery = `delete from presentation_spec.sheet_options where presentation_spec_id = $1`

const patchBasicInfo = `
	update presentation_spec.basic_info set updated_at = now() where id = $1;
`

const patchKeyOptions = `
	update presentation_spec.sheet_options set key = $1, active_columns = $2, position = $3, should_explode = $4 where presentation_spec_id = $5 and key = $6;
`

const patchKeyValueQuery = `
	update presentation_spec.specs set value = $1 where key = $2 and presentation_spec_id = $3;
`

const GetKeyValueQuery = `
	select value from presentation_spec.specs where key = $1 and presentation_spec_id = $2;
`

const patchKeySpec = `
	update presentation_spec.specs set value = $1, key = $2 where presentation_spec_id = $3 and key = $4
`
