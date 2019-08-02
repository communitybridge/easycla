// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package contractgroup

import (
	"context"
	"database/sql"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/jmoiron/sqlx"
)

// Repository interface contains a list of service methods
type Repository interface {
	CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (string, error)
	GetContractGroups(ctx context.Context, projectSfdcID string) ([]models.ContractGroup, error)

	CreateContractTemplate(ctx context.Context, contractID string, contractTemplate models.ContractTemplate) (string, error)
	GetLatestContractTemplate(ctx context.Context, contractGroupID string, contractType string) (models.ContractTemplate, error)

	CreateGitHubOrganization(ctx context.Context, contractID, userID string, githubOrg models.Github) (string, error)
	GetGithubOrganizations(ctx context.Context, contractGroupID string) ([]models.Github, error)

	CreateGerritInstance(ctx context.Context, projectSFDCID, contractID, userID string, gerritInstance models.Gerrit) (string, error)
	GetGerritInstances(ctx context.Context, contractGroupID string) ([]models.Gerrit, error)
	DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error

	GetContractGroupCCLASignatures(ctx context.Context, projectSFDCID string, contractID string) ([]models.CclaSignatureDetails, error)
	GetContractGroupICLASignatures(ctx context.Context, projectSFDCID string, contractID string) ([]models.IclaSignatureDetails, error)
}

type repository struct {
	db *sqlx.DB
	//session *session.Session
}

// NewRepository creates a new instance of the contract group service
func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

// CreateContractGroup create contract group creates a new contract group
func (repo repository) CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (string, error) {
	sql := `
		INSERT INTO cla.contract_group (
			project_sfdc_id, 
			name,
			individual_cla_enabled,
			corporate_cla_enabled,
			corporate_cla_requires_individual_cla)
		VALUES (
			$1,
			$2,
			$3, 
			$4,
			$5
		)
		RETURNING 
			contract_group_id;`

	var contractGroupID string
	err := repo.db.QueryRowx(
		sql,
		projectSfdcID,
		contractGroup.Name,
		contractGroup.IndividualClaEnabled,
		contractGroup.CorporateClaEnabled,
		contractGroup.CorporateClaRequiresIndividualCla,
	).Scan(&contractGroupID)
	if err != nil {
		log.Warnf("Error scanning for contract group, error: %v", err)
		return "", err
	}

	return contractGroupID, nil
}

// GetContractGroups returns a list of contract groups
func (repo repository) GetContractGroups(ctx context.Context, projectSfdcID string) ([]models.ContractGroup, error) {
	getContractGroupsSQL := `
		SELECT
			contract_group_id,
			project_sfdc_id,
			"name",
			corporate_cla_requires_individual_cla,
			individual_cla_enabled,
			corporate_cla_enabled
		FROM
			cla.contract_group
		WHERE
			project_sfdc_id = $1;`

	rows, err := repo.db.Queryx(getContractGroupsSQL, projectSfdcID)
	if err != nil && err != sql.ErrNoRows {
		log.Warnf("Error querying for contract groups, error: %v", err)
		return nil, err
	}
	if err == sql.ErrNoRows {
		log.Info("No rows returned from contract group query")
		return []models.ContractGroup{}, nil
	}

	contractGroups := []models.ContractGroup{}
	for rows.Next() {
		contractGroup := models.ContractGroup{}
		err := rows.StructScan(&contractGroup)
		if err != nil {
			rows.Close()
			return nil, err
		}

		contractGroups = append(contractGroups, contractGroup)
	}

	return contractGroups, nil
}

// CreateContractTemplate creates a new contract template
func (repo repository) CreateContractTemplate(ctx context.Context, contractID string, contractTemplate models.ContractTemplate) (string, error) {
	sql := `
		INSERT INTO cla.contract_template (
			contract_group_id, 
			"type", "document", 
			major_version, 
			minor_version)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
			
		)
		RETURNING 
			contract_template_id`

	var contractTemplateID string
	err := repo.db.QueryRowx(sql,
		contractID,
		contractTemplate.Type,
		contractTemplate.Document,
		contractTemplate.MajorVersion,
		contractTemplate.MinorVersion,
	).Scan(&contractTemplateID)
	if err != nil {
		log.Warnf("Error querying for contract templates, error: %v", err)
		return "", err
	}

	return contractTemplateID, nil
}

// GetLatestContractTemplate returns the latest contract template
func (repo repository) GetLatestContractTemplate(ctx context.Context, contractGroupID string, contractType string) (models.ContractTemplate, error) {
	getContractTempleteSQL :=
		`SELECT
			contract_template_id,
			contract_group_id,
			name,
			type,
			document,
			major_version,
			minor_version,
			created_at
		FROM
			cla.contract_template
		WHERE
			contract_group_id = $1
		AND
			type = $2
		ORDER BY
			created_at DESC
		LIMIT 1;`

	template := models.ContractTemplate{}
	err := repo.db.QueryRowx(
		getContractTempleteSQL,
		contractGroupID,
		contractType,
	).StructScan(&template)
	if err != nil && err != sql.ErrNoRows {
		log.Warnf("Error querying for contract template, error: %v", err)
		return models.ContractTemplate{}, err
	}
	if err == sql.ErrNoRows {
		log.Infof("No contract template found using group ID: %s, type: %s", contractGroupID, contractType)
		return models.ContractTemplate{}, nil
	}

	return template, nil
}

// GetGithubOrganizations returns a list of GH organizations
func (repo repository) GetGithubOrganizations(ctx context.Context, contractGroupID string) ([]models.Github, error) {
	getGithubOrganizationsSQL := `
		SELECT
			github_organization_id,
			contract_group_id,
			name,
			COALESCE( installation_id, '') AS installation_id,
			COALESCE( authorizing_github_id, '') AS authorizing_github_id,
			created_by,
			updated_by
		FROM
			cla.github_organization
		WHERE
			contract_group_id = $1;`

	rows, err := repo.db.Queryx(getGithubOrganizationsSQL, contractGroupID)
	if err != nil && err != sql.ErrNoRows {
		log.Warnf("Error querying for GH organizations, error: %v", err)
		return nil, err
	}
	if err == sql.ErrNoRows {
		log.Infof("No GH organizations found using group id: %s", contractGroupID)
		return []models.Github{}, nil
	}

	githubOrgs := []models.Github{}
	for rows.Next() {
		githubOrg := models.Github{}
		err := rows.StructScan(&githubOrg)
		if err != nil {
			rows.Close()
			return nil, err
		}

		githubOrgs = append(githubOrgs, githubOrg)
	}

	return githubOrgs, nil
}

// CreateGitHubOrganization creates a new GH organization
func (repo repository) CreateGitHubOrganization(ctx context.Context, contractID, userID string, githubOrg models.Github) (string, error) {
	sql := `
		INSERT INTO cla.github_organization (
			contract_group_id,name,
			created_by,updated_by
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$4
		)
		RETURNING 
			contract_group_id;`

	var githubOrgID string
	err := repo.db.QueryRowx(
		sql,
		contractID,
		githubOrg.Name,
		githubOrg.UpdatedBy,
	).Scan(&githubOrgID)
	if err != nil {
		return "", err
	}

	return githubOrgID, nil
}

// CreateGerritInstance creates a new gerrit instance
func (repo repository) CreateGerritInstance(ctx context.Context, projectSFDCID, contractID, userID string, gerritInstance models.Gerrit) (string, error) {
	// We have to verify that the provided Contract Group belongs to the specified Salesforce.com Project, so
	// a malicious user doesn't manipulate ownership through URL parameters. We verify using the WHERE clause
	// below. If the Contract Group relates to the SFDC Project, the Gerrit Instance inserts. Otherwise, the
	// insert fails, and a SQL No Rows error returns.
	sql := `
		INSERT INTO cla.gerrit_instance (
			contract_group_id,
			ldap_group_id,
			ldap_group_name,
			url,
			created_by,
			updated_by
		)
		SELECT
			$1,
			$2,
			$3,
			$4,
			$5,
			$5
		FROM
			cla.contract_group cg
		WHERE
			cg.project_sfdc_id = $6
		AND
			cg.contract_group_id = $1
		RETURNING
			gerrit_instance_id;`

	var gerritInstanceID string
	err := repo.db.QueryRowx(
		sql,
		contractID,
		gerritInstance.LdapGroupID,
		gerritInstance.LdapGroupName,
		gerritInstance.URL,
		userID,
		projectSFDCID,
	).Scan(&gerritInstanceID)
	if err != nil {
		return "", err
	}

	return gerritInstanceID, nil
}

// GetGerritInstances returns a list of gerrit instances
func (repo repository) GetGerritInstances(ctx context.Context, contractGroupID string) ([]models.Gerrit, error) {
	getGerritInstanceSQL := `
		SELECT
			gerrit_instance_id,
			contract_group_id,
			ldap_group_id,
			ldap_group_name,
			url,
			created_by,
			updated_by
		FROM
			cla.gerrit_instance
		WHERE
			contract_group_id = $1`

	rows, err := repo.db.Queryx(getGerritInstanceSQL, contractGroupID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return []models.Gerrit{}, nil
	}

	gerritInstances := []models.Gerrit{}
	for rows.Next() {
		gerritInstance := models.Gerrit{}
		err := rows.StructScan(&gerritInstance)
		if err != nil {
			rows.Close()
			return nil, err
		}

		gerritInstances = append(gerritInstances, gerritInstance)
	}

	return gerritInstances, nil
}

// DeleteGerritInstance deletes a gerrit instance
func (repo repository) DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error {

	deleteGerritInstanceSQL := `
	DELETE FROM 
		cla.gerrit_instance gi
	USING 
		cla.contract_group cg
	WHERE 
		cg.project_sfdc_id = $1
	AND 
		cg.contract_group_id = $2
	AND 
		gi.gerrit_instance_id = $3
	RETURNING 
		gi.gerrit_instance_id`

	var deletedGerritInstnaceID string
	err := repo.db.QueryRowx(deleteGerritInstanceSQL,
		projectSfdcID,
		contractID,
		gerritInstanceID,
	).Scan(&deletedGerritInstnaceID)

	if err != nil {
		return err
	}

	return nil
}

// GetContractGroupCCLASignatures gets a list of CCLA signatures
func (repo repository) GetContractGroupCCLASignatures(ctx context.Context, projectSFDCID string, contractID string) ([]models.CclaSignatureDetails, error) {
	getCCLASignaturesSQL := `
		SELECT
			company.name AS company_name,
			u.name AS user_name,
			ct.minor_version,
			ct.major_version,
			signed,
			ccla.updated_at AS signed_on
		FROM
			cla.contract_group cg
		JOIN
			cla.contract_template ct
		ON
			cg.contract_group_id = ct.contract_group_id
		JOIN
			cla.corporate_cla ccla
		ON
			ct.contract_template_id = ccla.contract_template_id
		JOIN
			cla."user" u
		ON
			ccla.signed_by = u.user_id
		JOIN
			cla.corporate_cla_group ccg
		ON
			ccla.corporate_cla_group_id = ccg.corporate_cla_group_id
		JOIN
			cla.company company
		ON
			ccg.company_id = company.company_id
		WHERE
			cg.contract_group_id = $1
		AND
			ct.type = 'CCLA'
		AND
			cg.project_sfdc_id = $2;`

	rows, err := repo.db.Queryx(getCCLASignaturesSQL, contractID, projectSFDCID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return []models.CclaSignatureDetails{}, nil
	}

	cclaSignatures := []models.CclaSignatureDetails{}
	for rows.Next() {
		cclaSignature := models.CclaSignatureDetails{}
		err = rows.StructScan(&cclaSignature)
		if err != nil {
			rows.Close()
			return nil, err
		}

		cclaSignatures = append(cclaSignatures, cclaSignature)
	}

	return cclaSignatures, nil
}

// GetContractGroupICLASignatures returns list of ICLA signatures
func (repo repository) GetContractGroupICLASignatures(ctx context.Context, projectSFDCID string, contractID string) ([]models.IclaSignatureDetails, error) {
	getICLASignaturesSQL := `
		SELECT
			u.name AS user_name,
			ct.minor_version,
			ct.major_version,
			signed,
			icla.updated_at AS signed_on
		FROM
			cla.contract_group cg
		JOIN
			cla.contract_template ct
		ON
			cg.contract_group_id = ct.contract_group_id
		JOIN
			cla.individual_cla icla
		ON
			ct.contract_template_id = icla.contract_template_id
		JOIN
			cla."user" u
		ON
			icla.user_id = u.user_id
		WHERE
			cg.contract_group_id = $1
		AND
			ct.type = 'ICLA'
		AND
			cg.project_sfdc_id = $2;`

	rows, err := repo.db.Queryx(getICLASignaturesSQL, contractID, projectSFDCID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return []models.IclaSignatureDetails{}, nil
	}

	iclaSignatures := []models.IclaSignatureDetails{}
	for rows.Next() {
		iclaSignature := models.IclaSignatureDetails{}
		err = rows.StructScan(&iclaSignature)
		if err != nil {
			rows.Close()
			return nil, err
		}

		iclaSignatures = append(iclaSignatures, iclaSignature)
	}

	return iclaSignatures, nil
}
