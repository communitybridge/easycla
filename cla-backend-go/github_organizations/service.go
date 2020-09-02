// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
)

// Service contains functions of GithubOrganizations service
type Service interface {
	GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error)
	AddGithubOrganization(externalProjectID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(externalProjectID string, githubOrgName string) error
	UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error
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

func (s service) GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error) {
	return s.repo.GetGithubOrganizations(externalProjectID, "")
}

func (s service) AddGithubOrganization(externalProjectID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error) {
	projectSFID := externalProjectID
	return s.repo.AddGithubOrganization(externalProjectID, projectSFID, input)
}

func (s service) DeleteGithubOrganization(externalProjectID string, githubOrgName string) error {
	err := s.ghRepository.DisableRepositoriesOfGithubOrganization(externalProjectID, githubOrgName)
	if err != nil {
		return err
	}
	return s.repo.DeleteGithubOrganization(externalProjectID, "", githubOrgName)
}

func (s service) UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error {
	return s.repo.UpdateGithubOrganization(projectSFID, organizationName, autoEnabled)
}
