// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// Service contains functions of Github Repository service
type Service interface {
	AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error)
	DeleteGithubRepository(externalProjectID string, repositoryID string) error
	ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error)
	GetGithubRepository(repositoryID string) (*models.GithubRepository, error)
}

type service struct {
	repo Repository
	//projectRepository project.Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
		//projectRepository: projectRepository,
	}
}

func (s *service) AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	return s.repo.AddGithubRepository(externalProjectID, input)
}
func (s *service) DeleteGithubRepository(externalProjectID string, repositoryID string) error {
	return s.repo.DeleteGithubRepository(externalProjectID, repositoryID)
}
func (s *service) ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error) {
	return s.repo.ListProjectRepositories(externalProjectID)
}
func (s *service) GetGithubRepository(repositoryID string) (*models.GithubRepository, error) {
	return s.repo.GetGithubRepository(repositoryID)
}
