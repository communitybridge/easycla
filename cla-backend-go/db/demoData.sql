
--USERS
	
	INSERT INTO cla."user" (user_id, "name")
	VALUES ('11ebaa98-3471-4fcf-99e8-729549e4f326','Test user');
	
	INSERT INTO cla."user" (user_id, "name")
	VALUES ('d76bf2b0-0593-407b-a9fe-d6532f5ace38','Test user 2');
	
	INSERT INTO cla.user_auth_provider(user_id, provider, provider_user_id)
	VALUES ('11ebaa98-3471-4fcf-99e8-729549e4f326', 'lfid', 'lfid_username');

	INSERT INTO cla.user_auth_provider(user_id, provider, provider_user_id)
	VALUES ('11ebaa98-3471-4fcf-99e8-729549e4f326', 'github', 'github_username');

	INSERT INTO cla.user_auth_provider(user_id, provider, provider_user_id)
	VALUES ('d76bf2b0-0593-407b-a9fe-d6532f5ace38', 'lfid', 'foobarski');

	INSERT INTO cla.user_auth_provider(user_id, provider, provider_user_id)
	VALUES ('d76bf2b0-0593-407b-a9fe-d6532f5ace38', 'github', 'user two');

	INSERT INTO cla.project_manager(user_id, project_sfdc_id)
	VALUES ('11ebaa98-3471-4fcf-99e8-729549e4f326', 'sfdc_id_one');

	INSERT INTO cla.project_manager(user_id, project_sfdc_id)
	VALUES ('11ebaa98-3471-4fcf-99e8-729549e4f326', 'sfdc_id_two');

	INSERT INTO cla.project_manager(user_id, project_sfdc_id)
	VALUES ('d76bf2b0-0593-407b-a9fe-d6532f5ace38', 'sfdc_id_one');

	INSERT INTO cla.project_manager(user_id, project_sfdc_id)
	VALUES ('d76bf2b0-0593-407b-a9fe-d6532f5ace38', 'sfdc_id_two');

	INSERT INTO cla."user" (user_id, "name")
	VALUES ('fd1abddd-a370-4de8-a95d-0bec5b21e485','Test user 3');
	
	-- COMPANY
	INSERT INTO cla.company (company_id, "name")
	VALUES ('445a532e-e938-431f-92cc-62a67e26cd1e','Test Comany 1');
	
	INSERT INTO cla.company (company_id, "name")
	VALUES ('5d6120e1-95fb-4975-90c1-54fbf063fc90','Test Comany 2');
	
	-- CONTRACT GOUPS
	INSERT INTO cla.contract_group 
	(
	project_sfdc_id, 
	"name",
	corporate_cla_requires_individual_cla,
	corporate_cla_enabled)
	VALUES (
	'123456',
	'demo CCLA Only', 
	false,
	true);
	
	INSERT INTO cla.contract_group 
	(contract_group_id,
	project_sfdc_id, 
	"name",
	corporate_cla_requires_individual_cla,
	corporate_cla_enabled)
	VALUES ('0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50' ,
	'456789',
	'demo Contract Group 2', 
	false,
	true);
	
		INSERT INTO cla.contract_group 
	(contract_group_id,
	project_sfdc_id, 
	"name",
	corporate_cla_requires_individual_cla,
	corporate_cla_enabled)
	VALUES ('0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50' ,
	'456789',
	'demo Contract Group 2', 
	false,
	true);
	
	--Project Manager 
	INSERT INTO cla.project_manager
	(user_id,
	project_sfdc_id) 
	VALUES 
	('11ebaa98-3471-4fcf-99e8-729549e4f326', 
	'123sfdc');
	
	-- Contract Templates
	INSERT INTO cla.contract_template 
	(contract_template_id,
	contract_group_id, 
	"type",
	"document", 
	major_version, 
	minor_version,
	"name")
	VALUES ('b65da042-3d6b-408a-aaed-6155c8fdf577',
	'0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50',
	'CCLA', 
	'{"name": "Paint house", "tags": ["Improvements", "Office"], "finished": true}',
	1,
	0,
	'test template 1');
	
	INSERT INTO "cla"."contract_template"
	("contract_template_id",
	"contract_group_id",
	"type",
	"document", 
	"major_version",
	"minor_version",
	"name") 
	VALUES ('35ef4864-5174-4394-b07e-408fa1247cb6',
	'0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50', 
	'ICLA', '{"name": "Paint house", "tags": ["Improvements", "Office"], "finished": true}',
	1, 
	0,
	'test template 2');
	
	INSERT INTO "cla"."contract_template"
	("contract_template_id",
	"contract_group_id",
	"type",
	"document", 
	"major_version",
	"minor_version",
	"name") 
	VALUES ('4694efcf-2d5f-46bf-a924-80abcdcd837c',
	'ea3bac44-08c0-4947-8c81-8c02c3435a25', 
	'ICLA', '{"name": "Paint house"}',
	1, 
	0,
	'test template 3');
	
	INSERT INTO "cla"."contract_template"
	("contract_template_id",
	"contract_group_id",
	"type",
	"document", 
	"major_version",
	"minor_version",
	"name") 
	VALUES ('e7ccdbb3-64a7-4943-a1a4-21260af52a3a',
	'0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50', 
	'CCLA', 
	'{"name": "Paint house", "tags": ["Improvements", "Office"], "finished": true}',
	2, 
	0,
	'test template 4');
	
	-- Corporate cla group
	
	INSERT INTO cla.corporate_cla_group
	(corporate_cla_group_id, 
	email_whitelist, 
	company_id, 
	contract_group_id)
	VALUES ('e630255b-7974-47f1-969f-2b9fb3d271b4',
	'{"email":"test@test.com"}',
	'5d6120e1-95fb-4975-90c1-54fbf063fc90',
	'0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50');
	
	-- CLA Manager 
	INSERT INTO cla.cla_manager
	(user_id, corporate_cla_group_id) 
	VALUES 
	('d76bf2b0-0593-407b-a9fe-d6532f5ace38', 
	'e630255b-7974-47f1-969f-2b9fb3d271b4'); 
	
	--Docusign document 
	INSERT INTO cla.docusign (docusign_id, envelope_id)
	VALUES ('97be5d8e-0f3d-49d6-b3bd-f287f4b4929c',
	'333');
	
	--Docusign document
	INSERT INTO cla.docusign (docusign_id, envelope_id)
	VALUES ('33473f71-b696-4547-bb13-abb0e6aec910',
	'444');

	
	-- Corporate Cla
	INSERT INTO cla.corporate_cla
	(corporate_cla_id, 
	corporate_cla_group_id,
	contract_template_id,
	docusign_id,
	signatory_email,
	signed_by,
	signed)
	VALUES
	('5f31f687-2f03-43ff-a8e6-b0081ba22cab',
	'e630255b-7974-47f1-969f-2b9fb3d271b4',
	'e7ccdbb3-64a7-4943-a1a4-21260af52a3a',
	'33473f71-b696-4547-bb13-abb0e6aec910',
	'signatory@email.com',
	'daea7b2e-9fad-4628-8aa9-4f6d158350db',
	true);

	--USER 1 ICLA signed = true
	INSERT INTO cla.individual_cla 
	(individual_cla_id, 
	contract_template_id, 
	user_id, 
	docusign_id, 
	signed)
	VALUES ('97be5d8e-0f3d-49d6-b3bd-f287f4b4929c',
	'35ef4864-5174-4394-b07e-408fa1247cb6',
	'11ebaa98-3471-4fcf-99e8-729549e4f326',
	'97be5d8e-0f3d-49d6-b3bd-f287f4b4929c',
	true);
	
	-- User 2 ICLA signed = false 	
	INSERT INTO cla.individual_cla 
	(individual_cla_id, 
	contract_template_id, 
	user_id, 
	docusign_id, 
	signed)
	VALUES ('741fb220-7b79-41ed-aec4-76e49dc48fa3',
	'35ef4864-5174-4394-b07e-408fa1247cb6',
	'd76bf2b0-0593-407b-a9fe-d6532f5ace38',
	'33473f71-b696-4547-bb13-abb0e6aec910',
	false);
	
	-- github org
	INSERT INTO cla.github_organization
	(github_organization_id, 
	contract_group_id,
	"name",
	installation_id,
	authorizing_user_name,
	authorizing_github_id,
	created_by,
	updated_by)
	VALUES 
	('7f415c29-a2f7-465d-8251-541fe48c1f5e',
	'ea3bac44-08c0-4947-8c81-8c02c3435a25',
	'Demo ICLA Org',
	'1111',
	'Autorizing Username',
	'authorizing GH ID',
	'11ebaa98-3471-4fcf-99e8-729549e4f326',
	'11ebaa98-3471-4fcf-99e8-729549e4f326');
	
	-- Gerrit Instace 
	INSERT INTO cla.gerrit_instance
	(gerrit_instance_id,
	contract_group_id,
	ldap_group_id,
	ldap_group_name,
	url,
	created_by,
	updated_by)
	VALUES
	('a61924cc-ab10-4b45-b23c-142ef609b85d',
	'0e8eaca6-667e-4cc6-a354-b6ea1cfa8a50',
	'1234',
	'LDAP group name',
	'ldap url',
	'd76bf2b0-0593-407b-a9fe-d6532f5ace38',
	'd76bf2b0-0593-407b-a9fe-d6532f5ace38');
	
	
	
