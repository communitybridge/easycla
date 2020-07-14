// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	v1GithubOrg "github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	v1Repositories "github.com/communitybridge/easycla/cla-backend-go/repositories"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	"github.com/jinzhu/copier"
)

func v2GithubOrgnizationsModel(in *v1Models.GithubOrganizations) (*models.GithubOrganizations, error) {
	var response models.GithubOrganizations
	err := copier.Copy(&response, in)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func v2GithubOrgnizationModel(in *v1Models.GithubOrganization) (*models.GithubOrganization, error) {
	var response models.GithubOrganization
	err := copier.Copy(&response, in)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Service contains functions of GithubOrganizations service
type Service interface {
	GetGithubOrganizations(projectSFID string) (*models.GithubOrganizations, error)
	AddGithubOrganization(projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(projectSFID string, githubOrgName string) error
}

type service struct {
	repo         v1GithubOrg.Repository
	ghRepository v1Repositories.Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo v1GithubOrg.Repository, ghRepository v1Repositories.Repository) Service {
	return service{
		repo:         repo,
		ghRepository: ghRepository,
	}
}

func (s service) GetGithubOrganizations(projectSFID string) (*models.GithubOrganizations, error) {
	resp, err := s.repo.GetGithubOrganizations("", projectSFID)
	if err != nil {
		return nil, err
	}
	return v2GithubOrgnizationsModel(resp)
}

func (s service) AddGithubOrganization(projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error) {
	var in v1Models.CreateGithubOrganization
	err := copier.Copy(&in, input)
	if err != nil {
		return nil, err
	}
	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(projectSFID)
	if err != nil {
		return nil, err
	}
	var externalProjectID string
	if project.Parent == "" {
		externalProjectID = projectSFID
	} else {
		externalProjectID = project.Parent
	}
	resp, err := s.repo.AddGithubOrganization(externalProjectID, projectSFID, &in)
	if err != nil {
		return nil, err
	}
	return v2GithubOrgnizationModel(resp)
}

func (s service) DeleteGithubOrganization(projectSFID string, githubOrgName string) error {
	err := s.ghRepository.DeleteRepositoriesOfGithubOrganization("", projectSFID, githubOrgName)
	if err != nil {
		return err
	}
	return s.repo.DeleteGithubOrganization("", projectSFID, githubOrgName)
}
