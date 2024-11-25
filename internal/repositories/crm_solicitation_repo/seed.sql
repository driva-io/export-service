create schema crm;

CREATE TABLE crm.company (
    id UUID PRIMARY KEY,
    crm VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    refresh_token VARCHAR(255),
    access_token VARCHAR(255),
    expires_in VARCHAR(50),
    refreshed_at VARCHAR(255),
    environment VARCHAR(255),
    token VARCHAR(255),
    webhook VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    instance_url VARCHAR(255),
    merge VARCHAR(255),
    mapping TEXT,
    mapping_linkedin TEXT,
    company_id VARCHAR(255),
    user_who_installed_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    workspace_id VARCHAR(255)
);

INSERT INTO crm.company (
    id, crm, name, refresh_token, access_token, expires_in, refreshed_at, 
    environment, token, webhook, email, password, instance_url, merge, 
    mapping, mapping_linkedin, company_id, user_who_installed_id, 
    created_at, updated_at, workspace_id
) VALUES (
    '123e4567-e89b-12d3-a456-426655440003', 'hubspot', 'Driva Teste F', 'refresh_token_1', 'access_token_1', '3600', 
    '2022-01-01 00:00:00 -03:00', 'production', 'token_1', 'https://webhook.url/1', 
    'francisco.becheli@driva.com.br', null, 'https://instance.url/1', 'merge_A', 
    null, null, 
    'company_id_1', 'user_id_1', '2022-01-01 00:00:00 -03:00', '2022-01-01 00:00:00 -03:00', 'workspace_1'
);

INSERT INTO crm.company (
    id, crm, name, refresh_token, access_token, expires_in, refreshed_at, 
    environment, token, webhook, email, password, instance_url, merge, 
    mapping, mapping_linkedin, company_id, user_who_installed_id, 
    created_at, updated_at, workspace_id
) VALUES (
    '123e4567-e89b-12d3-a456-426655440004', 'salesforce', 'Driva Teste F', 'refresh_token_1', 'access_token_1', '3600', 
    '2022-01-01 00:00:00 -03:00', 'production', 'token_1', 'https://webhook.url/1', 
    'francisco.becheli@driva.com.br', null, 'https://instance.url/1', 'merge_A', 
    null, '{"linkedin_mapping": "value"}', 
    'company_id_1', 'user_id_1', '2022-01-01 00:00:00 -03:00', '2022-01-01 00:00:00 -03:00', 'workspace_1'
);
