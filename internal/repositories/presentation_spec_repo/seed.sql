create schema presentation_spec;

create table presentation_spec.basic_info (
    id UUID primary key,
    version int not null default 2,
    base text not null,
    created_at timestamp
    with
        time zone not null DEFAULT now(),
        updated_at timestamp
    with
        time zone not null DEFAULT now(),
        user_email text not null,
        user_company text not null,
        service text not null,
        is_default bool default false,
        unique (
            user_email,
            user_company,
            service,
            base
        ) -- TODO: adicionar esse unique no banco de produção
);

create unique index unique_default_spec_idx on presentation_spec.basic_info (service, base)
where (is_default);

create table presentation_spec.sheet_options (presentation_spec_id UUID references presentation_spec.basic_info (id) on delete cascade, key text, active_columns text[], position integer, should_explode bool);

create table presentation_spec.specs (
    presentation_spec_id UUID references presentation_spec.basic_info (id) on delete cascade,
    key text,
    value JSONB
);

insert into
    presentation_spec.basic_info (
        id,
        version,
        base,
        created_at,
        updated_at,
        user_email,
        user_company,
        service,
        is_default
    )
values (
        '123e4567-e89b-12d3-a456-426655440000',
        2,
        'empresas',
        '2022-01-01 00:00:00 -03:00',
        '2022-01-01 00:00:00 -03:00',
        'victor@driva.com.br',
        'Driva',
        'enrichment_test',
        false
    ),
    (
        '123e4567-e89b-12d3-a456-426655440001',
        2,
        'empresas',
        '2022-01-01 00:00:00 -03:00',
        '2022-01-01 00:00:00 -03:00',
        '',
        '',
        'enrichment_test',
        true
    );

insert into
    presentation_spec.specs (
        presentation_spec_id,
        key,
        value
    )
values (
        '123e4567-e89b-12d3-a456-426655440000',
        'RFB',
        '{"CNPJ": "cnpj"}'
    ),
    (
        '123e4567-e89b-12d3-a456-426655440001',
        'RFB',
        '{"CNPJ": "cnpj", "Nome": "razao_social"}'
    );

insert into
    presentation_spec.sheet_options (
        presentation_spec_id,
        key,
        active_columns,
        position,
        should_explode
    )
values (
        '123e4567-e89b-12d3-a456-426655440000',
        'RFB',
        '{"CNPJ"}',
        0,
        false
    ),
    (
        '123e4567-e89b-12d3-a456-426655440001',
        'RFB',
        '{"CNPJ", "Nome"}',
        0,
        false
    );