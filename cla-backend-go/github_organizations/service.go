// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"context"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
)

// Service contains functions of GithubOrganizations service
type Service interface {
	GetGithubOrganizations(ctx context.Context, externalProjectID string) (*models.GithubOrganizations, error)
	GetGithubOrganizationByName(ctx context.Context, githubOrgName string) (*models.GithubOrganization, error)
	AddGithubOrganization(ctx context.Context, externalProjectID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(ctx context.Context, externalProjectID string, githubOrgName string) error
	UpdateGithubOrganization(ctx context.Context, projectSFID string, organizationName string, autoEnabled bool, branchProtectionEnabled bool) error
}

type service struct {
	repo         Repository
	ghRepository repositories.Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository, ghRepository repositories.Repository) Service {
	return service{
		repo:         repo,
		ghRepository: ghRepository,
	}
}

func (s service) GetGithubOrganizations(ctx context.Context, projectSFID string) (*models.GithubOrganizations, error) {
	externalProjectID, err := projectHelper(projectSFID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetGithubOrganizations(ctx, externalProjectID, projectSFID)
}

func (s service) GetGithubOrganizationByName(ctx context.Context, githubOrgName string) (*models.GithubOrganization, error) {
	gitHubOrgs, err := s.repo.GetGithubOrganizationByName(ctx, githubOrgName)
	if err != nil {
		return nil, err
	}
	if len(gitHubOrgs.List) == 0 {
		log.Debugf("no matching github organization matches organization name: %s", githubOrgName)
		return nil, nil
	}

	if len(gitHubOrgs.List) > 1 {
		log.Warnf("More than 1 github organization matches organization name: %s - using first one", githubOrgName)
	}

	return gitHubOrgs.List[0], err
}

func (s service) AddGithubOrganization(ctx context.Context, projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error) {
	externalProjectID, err := projectHelper(projectSFID)
	if err != nil {
		return nil, err
	}
	return s.repo.AddGithubOrganization(ctx, externalProjectID, projectSFID, input)
}

func (s service) DeleteGithubOrganization(ctx context.Context, projectSFID string, githubOrgName string) error {
	externalProjectID, err := projectHelper(projectSFID)
	if err != nil {
		return err
	}

	err = s.ghRepository.DisableRepositoriesOfGithubOrganization(ctx, externalProjectID, githubOrgName)
	if err != nil {
		return err
	}
	return s.repo.DeleteGithubOrganization(ctx, externalProjectID, projectSFID, githubOrgName)
}

func (s service) UpdateGithubOrganization(ctx context.Context, projectSFID string, organizationName string, autoEnabled bool, branchProtectionEnabled bool) error {
	return s.repo.UpdateGithubOrganization(ctx, projectSFID, organizationName, autoEnabled, branchProtectionEnabled)
}

func projectHelper(projectSFID string) (string, error) {
	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(projectSFID)
	if err != nil {
		return "", err
	}

	var externalProjectID string
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		externalProjectID = projectSFID
	} else {
		externalProjectID = project.Parent
	}

	return externalProjectID, nil
}
