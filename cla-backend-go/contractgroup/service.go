package contractgroup

import (
	"context"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
)

var (
	userID = "***REMOVED***"
)

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

func NewService(contractGroupRepo Repository) service {
	return service{
		contractGroupRepo: contractGroupRepo,
	}
}

func (s service) CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (models.ContractGroup, error) {
	contractGroupID, err := s.contractGroupRepo.CreateContractGroup(ctx, projectSfdcID, contractGroup)
	if err != nil {
		return models.ContractGroup{}, err
	}

	contractGroup.ContractGroupID = contractGroupID
	contractGroup.ProjectSfdcID = projectSfdcID

	return contractGroup, nil
}

func (s service) GetContractGroups(ctx context.Context, projectSfdcID string) ([]models.ContractGroup, error) {
	contractGroups, err := s.contractGroupRepo.GetContractGroups(ctx, projectSfdcID)
	if err != nil {
		return nil, err
	}

	for i, contractGroup := range contractGroups {
		contractGroups[i].CclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "CCLA")
		if err != nil {
			return nil, err
		}

		contractGroups[i].IclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "ICLA")
		if err != nil {
			return nil, err
		}

		contractGroups[i].GithubOrganizations, err = s.contractGroupRepo.GetGithubOrganizatons(ctx, contractGroup.ContractGroupID)
		if err != nil {
			return nil, err
		}

		contractGroups[i].GerritInstances, err = s.contractGroupRepo.GetGerritInstances(ctx, contractGroup.ContractGroupID)
		if err != nil {
			return nil, err
		}
	}

	return contractGroups, nil
}

func (s service) CreateContractTemplate(ctx context.Context, contractTemplate models.ContractTemplate, contractID string) (models.ContractTemplate, error) {
	contractTemplateID, err := s.contractGroupRepo.CreateContractTemplate(ctx, contractID, contractTemplate)
	if err != nil {
		return models.ContractTemplate{}, err
	}

	contractTemplate.ContractTemplateID = contractTemplateID

	return contractTemplate, nil
}

func (s service) CreateGitHubOrganization(ctx context.Context, contractID string, githubOrg models.Github) (models.Github, error) {
	githubOrgID, err := s.contractGroupRepo.CreateGitHubOrganization(ctx, contractID, userID, githubOrg)
	if err != nil {
		return models.Github{}, err
	}

	githubOrg.GithubOrganizationID = githubOrgID

	return githubOrg, nil
}

func (s service) CreateGerritInstance(ctx context.Context, projectSFDCID, contractID, userID string, gerritInstance models.Gerrit) (models.Gerrit, error) {
	gerritInstanceID, err := s.contractGroupRepo.CreateGerritInstance(ctx, projectSFDCID, contractID, userID, gerritInstance)
	if err != nil {
		return models.Gerrit{}, err
	}

	gerritInstance.GerritInstanceID = gerritInstanceID

	return gerritInstance, nil

}

func (s service) DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error {
	err := s.contractGroupRepo.DeleteGerritInstance(ctx, projectSfdcID, contractID, gerritInstanceID)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetContractGroupSignatures(ctx context.Context, projectSFDCID string, contractID string) (models.ContractGroupSignatures, error) {
	contractGoupSignatures := models.ContractGroupSignatures{ContractGroupID: contractID}
	var err error

	contractGoupSignatures.CclaSignatures, err = s.contractGroupRepo.GetContractGroupCCLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		return models.ContractGroupSignatures{}, err
	}

	contractGoupSignatures.IclaSignatures, err = s.contractGroupRepo.GetContractGroupICLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		return models.ContractGroupSignatures{}, err
	}

	return contractGoupSignatures, nil
}
