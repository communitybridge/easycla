-- migrate:up
CREATE SCHEMA cla;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE cla."user"(
	user_id UUID NOT NULL DEFAULT gen_random_uuid(),
	name TEXT,

	PRIMARY KEY(user_id)
);

CREATE TABLE cla.user_auth_provider(
	user_id UUID NOT NULL REFERENCES cla."user"(user_id),
	provider_user_id TEXT NOT NULL,
	provider TEXT NOT NULL,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	UNIQUE (provider_user_id, user_id, provider) 
);

CREATE TABLE cla.company (
	company_id UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" TEXT,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(company_id)
);

CREATE TABLE cla.user_email(
	user_id UUID REFERENCES cla."user"(user_id),
	email TEXT NOT NULL,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now())
);

CREATE TABLE cla.contract_group (
	contract_group_id UUID NOT NULL DEFAULT gen_random_uuid(),
	project_sfdc_id TEXT NOT NULL,
	name TEXT NOT NULL,
	individual_cla_enabled BOOLEAN DEFAULT FALSE,
	corporate_cla_enabled BOOLEAN DEFAULT FALSE,
	corporate_cla_requires_individual_cla BOOLEAN DEFAULT FALSE,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(contract_group_id)
);

CREATE TABLE cla.corporate_cla_group (
	corporate_cla_group_id UUID NOT NULL DEFAULT gen_random_uuid(),
	email_whitelist jsonb,
	company_id UUID NOT NULL REFERENCES cla.company(company_id),
	contract_group_id UUID NOT NULL REFERENCES cla.contract_group(contract_group_id),
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	UNIQUE (company_id, contract_group_id),

	PRIMARY KEY(corporate_cla_group_id)
);

CREATE TYPE cla.contract_template_type AS ENUM (
    'CCLA',
    'ICLA'
);

CREATE TABLE cla.contract_template (
	contract_template_id UUID NOT NULL DEFAULT gen_random_uuid(),
	contract_group_id UUID NOT NULL REFERENCES cla.contract_group(contract_group_id),
	"type" cla.contract_template_type NOT NULL,
	"name" TEXT,
	document jsonb,
	major_version int NOT NULL,
	minor_version int NOT NULL,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	UNIQUE(contract_group_id, "type", major_version, minor_version),

	PRIMARY KEY(contract_template_id)
);

CREATE TABLE cla.github_organization (
	github_organization_id UUID NOT NULL DEFAULT gen_random_uuid(),
	contract_group_id UUID NOT NULL REFERENCES cla.contract_group(contract_group_id),
	name TEXT,
	installation_id TEXT,
	authorizing_user_name TEXT,
	authorizing_github_id TEXT,
	created_by UUID NOT NULL REFERENCES cla.user(user_id),
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_by UUID NOT NULL REFERENCES cla.user(user_id),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(github_organization_id)
);

CREATE TABLE cla.gerrit_instance (
	gerrit_instance_id UUID NOT NULL DEFAULT gen_random_uuid(),
	contract_group_id UUID NOT NULL REFERENCES cla.contract_group(contract_group_id),
	ldap_group_id TEXT,
	ldap_group_name TEXT,
	url TEXT,
	created_by UUID NOT NULL REFERENCES cla.user(user_id),
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_by UUID NOT NULL REFERENCES cla.user(user_id),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(gerrit_instance_id)
);

CREATE TABLE cla.corporate_cla_group_confirmed_users (
	user_id UUID NOT NULL REFERENCES cla.user(user_id),
	corporate_cla_group_id UUID NOT NULL REFERENCES cla.corporate_cla_group(corporate_cla_group_id)
);

CREATE TABLE cla.docusign (
	docusign_id UUID NOT NULL DEFAULT gen_random_uuid(),
	envelope_id INT,
	callback_url TEXT,
	sign_url TEXT,
	return_url TEXT,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(docusign_id)
);

CREATE TABLE cla.individual_cla (
	individual_cla_id UUID NOT NULL DEFAULT gen_random_uuid(),
	contract_template_id UUID NOT NULL REFERENCES cla.contract_template(contract_template_id),
	user_id UUID NOT NULL REFERENCES cla.user(user_id),
	docusign_id UUID NOT NULL REFERENCES cla.docusign(docusign_id),
	signed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY(individual_cla_id)
);

CREATE TABLE cla.corporate_cla (
	corporate_cla_id UUID NOT NULL DEFAULT gen_random_uuid(),
	corporate_cla_group_id UUID NOT NULL REFERENCES cla.corporate_cla_group(corporate_cla_group_id),
	contract_template_id UUID NOT NULL REFERENCES cla.contract_template(contract_template_id),
	docusign_id UUID NOT NULL REFERENCES cla.docusign(docusign_id),
	signatory_email TEXT,
	signed_by UUID,
	signed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY (corporate_cla_id)
);

CREATE TABLE cla.project_manager (
	user_id UUID NOT NULL REFERENCES cla.user(user_id),
	project_sfdc_id TEXT NOT NULL,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY (user_id, project_sfdc_id)
);

CREATE TABLE cla.cla_manager (
	user_id UUID NOT NULL REFERENCES cla.user(user_id),
	corporate_cla_group_id UUID NOT NULL REFERENCES cla.corporate_cla_group(corporate_cla_group_id),
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),

	PRIMARY KEY (user_id, corporate_cla_group_id)
);

-- migrate:down
DROP TABLE cla.cla_manager;
DROP TABLE cla.project_manager;
DROP TABLE cla.corporate_cla;
DROP TABLE cla.individual_cla;
DROP TABLE cla.docusign;
DROP TABLE cla.corporate_cla_group_confirmed_users;
DROP TABLE cla.gerrit_instance;
DROP TABLE cla.github_organization;
DROP TABLE cla.contract_template;
DROP TYPE cla.contract_template_type;
DROP TABLE cla.corporate_cla_group;
DROP TABLE cla.contract_group;
DROP TABLE cla.user_email;
DROP TABLE cla.company;
DROP TABLE cla.user_auth_provider;
DROP TABLE cla.user;
DROP SCHEMA cla;
