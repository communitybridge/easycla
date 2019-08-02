// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package contractgroup

import (
	"context"

	"github.com/aymerick/raymond"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

var (
	userID = "<redacted>"
)

// Service interface functions
type Service interface {
	CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (models.ContractGroup, error)
	GetContractGroups(ctx context.Context, projectID string) ([]models.ContractGroup, error)

	CreateContractTemplate(ctx context.Context, contractTemplate models.ContractTemplate, contractID string) (models.ContractTemplate, error)

	CreateGitHubOrganization(ctx context.Context, contractID string, githubOrg models.Github) (models.Github, error)

	CreateGerritInstance(ctx context.Context, projectSfdcID, contractID, userID string, gerritInstance models.Gerrit) (models.Gerrit, error)
	DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error

	GetContractGroupSignatures(ctx context.Context, projectSfdcID string, contractID string) (models.ContractGroupSignatures, error)
}

type service struct {
	contractGroupRepo Repository
}

// NewService API call constructs a new service instance
func NewService(contractGroupRepo Repository) service {
	return service{
		contractGroupRepo: contractGroupRepo,
	}
}

// InjectProjectInformationIntoTemplate
func (s service) InjectProjectInformationIntoTemplate(projectName, shortProjectName, documentType, majorVersion, minorVersion, contactEmail string) string {
	// DocRaptor API likes HTML in single line
	templateBefore := `<html><body><p style=\"text-align: center\">{{projectName}}<br />{{documentType}} Contributor License Agreement (\"Agreement\")v{{majorVersion}}.{{minorVersion}}</p><p>Thank you for your interest in {{projectName}} project (“{{shortProjectName}}”) of The Linux Foundation (the “Foundation”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Foundation must have a Contributor License Agreement (“CLA”) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of {{shortProjectName}}, the Foundation and its users; it does not change your rights to use your own Contributions for any other purpose.</p><p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Foundation or its third-party service providers, or email a PDF of the signed agreement to {{contactEmail}}. Please read this document carefully before signing and keep a copy for your records.</p></body></html>`
	fieldsMap := map[string]string{
		"projectName":      projectName,
		"shortProjectName": shortProjectName,
		"documentType":     documentType,
		"majorVersion":     majorVersion,
		"minorVersion":     minorVersion,
		"contactEmail":     contactEmail,
	}

	templateAfter, err := raymond.Render(templateBefore, fieldsMap)
	if err != nil {
		log.Warnf("Failed to enter fields into HTML, error: %v", err)
	}

	return templateAfter
}

// CreateContractGroup creates a new contract group based on the specified parameters
func (s service) CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (models.ContractGroup, error) {
	contractGroupID, err := s.contractGroupRepo.CreateContractGroup(ctx, projectSfdcID, contractGroup)
	if err != nil {
		log.Warnf("Failed to create contract group using project ID: %s and contract group: %v, error: %v", projectSfdcID, contractGroup, err)
		return models.ContractGroup{}, err
	}

	contractGroup.ContractGroupID = contractGroupID
	contractGroup.ProjectSfdcID = projectSfdcID

	return contractGroup, nil
}

// GetContractGroups returns a list of contract group models based on the specified input parameters
func (s service) GetContractGroups(ctx context.Context, projectSfdcID string) ([]models.ContractGroup, error) {
	contractGroups, err := s.contractGroupRepo.GetContractGroups(ctx, projectSfdcID)
	if err != nil {
		log.Warnf("Error retrieving contract group using SFID: %s, error: %v", projectSfdcID, err)
		return nil, err
	}

	for i, contractGroup := range contractGroups {
		contractGroups[i].CclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "CCLA")
		if err != nil {
			log.Warnf("Error retrieving latest CCLA contract template contract using group ID: %s, error: %v", contractGroup.ContractGroupID, err)
			return nil, err
		}

		contractGroups[i].IclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "ICLA")
		if err != nil {
			log.Warnf("Error retrieving latest ICLA contract template using contract group ID: %s, error: %v", contractGroup.ContractGroupID, err)
			return nil, err
		}

		contractGroups[i].GithubOrganizations, err = s.contractGroupRepo.GetGithubOrganizations(ctx, contractGroup.ContractGroupID)
		if err != nil {
			log.Warnf("Error retrieving GH organizations using contract group ID: %s, error: %v", contractGroup.ContractGroupID, err)
			return nil, err
		}

		contractGroups[i].GerritInstances, err = s.contractGroupRepo.GetGerritInstances(ctx, contractGroup.ContractGroupID)
		if err != nil {
			log.Warnf("Error retrieving Gerrit instances using contract group ID: %s, error: %v", contractGroup.ContractGroupID, err)
			return nil, err
		}
	}

	return contractGroups, nil
}

// CreateContractTemplate creates a enw contract template based on the specified parameters
func (s service) CreateContractTemplate(ctx context.Context, contractTemplate models.ContractTemplate, contractID string) (models.ContractTemplate, error) {
	contractTemplateID, err := s.contractGroupRepo.CreateContractTemplate(ctx, contractID, contractTemplate)
	if err != nil {
		log.Warnf("Error creating contract template using contract ID: %s with template: %v, error: %v", contractID, contractTemplate, err)
		return models.ContractTemplate{}, err
	}

	contractTemplate.ContractTemplateID = contractTemplateID

	return contractTemplate, nil
}

// CreateGitHubOrganization creates a GitHub organization contract group
func (s service) CreateGitHubOrganization(ctx context.Context, contractID string, githubOrg models.Github) (models.Github, error) {
	githubOrgID, err := s.contractGroupRepo.CreateGitHubOrganization(ctx, contractID, userID, githubOrg)
	if err != nil {
		log.Warnf("Error creating GH organization using contractID: %s, userID: %s, GH org: %s, error: %v", contractID, userID, githubOrg, err)
		return models.Github{}, err
	}

	githubOrg.GithubOrganizationID = githubOrgID

	return githubOrg, nil
}

// CreateGerritInstance creates a Gerrit instance contract group
func (s service) CreateGerritInstance(ctx context.Context, projectSFDCID, contractID, userID string, gerritInstance models.Gerrit) (models.Gerrit, error) {
	gerritInstanceID, err := s.contractGroupRepo.CreateGerritInstance(ctx, projectSFDCID, contractID, userID, gerritInstance)
	if err != nil {
		log.Warnf("Error creating Gerrit instances using contractID: %s, userID: %s, gerrit instance: %v, error: %v", contractID, userID, gerritInstance, err)
		return models.Gerrit{}, err
	}

	gerritInstance.GerritInstanceID = gerritInstanceID

	return gerritInstance, nil

}

// DeleteGerritInstances deletes the specified Gerrit instance
func (s service) DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error {
	err := s.contractGroupRepo.DeleteGerritInstance(ctx, projectSfdcID, contractID, gerritInstanceID)
	if err != nil {
		log.Warnf("Error deleting Gerrit instances using contractID: %s, userID: %s, gerrit instance: %s, error: %v", contractID, userID, gerritInstanceID, err)
		return err
	}
	return nil
}

// GetContractGroupSignatures returns a data model containing the contract group signatures
func (s service) GetContractGroupSignatures(ctx context.Context, projectSFDCID string, contractID string) (models.ContractGroupSignatures, error) {
	contractGoupSignatures := models.ContractGroupSignatures{ContractGroupID: contractID}
	var err error

	contractGoupSignatures.CclaSignatures, err = s.contractGroupRepo.GetContractGroupCCLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		log.Warnf("Error fetching CCLA contract group signatures using SFID: %s, contract ID: %s, error: %v", projectSFDCID, contractID, err)
		return models.ContractGroupSignatures{}, err
	}

	contractGoupSignatures.IclaSignatures, err = s.contractGroupRepo.GetContractGroupICLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		log.Warnf("Error fetching ICLA contract group signatures using SFID: %s, contract ID: %s, error: %v", projectSFDCID, contractID, err)
		return models.ContractGroupSignatures{}, err
	}

	return contractGoupSignatures, nil
}
