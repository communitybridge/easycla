// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service contains functions of Github Repository service
type Service interface {
	AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error)
	EnableRepository(repositoryID string) error
	DisableRepository(repositoryID string) error
	ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error)
	GetRepository(repositoryID string) (*models.GithubRepository, error)
	DisableRepositoriesByProjectID(projectID string) (int, error)
	GetRepositoriesByCLAGroup(claGroupID string) ([]*models.GithubRepository, error)
}

type service struct {
	repo Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	projectSFID := externalProjectID
	return s.repo.AddGithubRepository(externalProjectID, projectSFID, input)
}

func (s *service) EnableRepository(repositoryID string) error {
	return s.repo.EnableRepository(repositoryID)
}

func (s *service) DisableRepository(repositoryID string) error {
	return s.repo.DisableRepository(repositoryID)
}

func (s *service) ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error) {
	return s.repo.ListProjectRepositories(externalProjectID, "", true)
}
func (s *service) GetRepository(repositoryID string) (*models.GithubRepository, error) {
	return s.repo.GetRepository(repositoryID)
}

// DisableRepositoriesByProjectID disables the repositories by project ID
func (s *service) DisableRepositoriesByProjectID(projectID string) (int, error) {
	var deleteErr error
	// Return the list of GitHub repositories by CLA Group for those that are currently enabled
	ghOrgs, err := s.repo.GetCLAGroupRepositoriesGroupByOrgs(projectID, true)
	if err != nil {
		return 0, err
	}
	if len(ghOrgs) > 0 {
		log.Debugf("Deleting repositories for project :%s", projectID)
		for _, ghOrg := range ghOrgs {
			for _, item := range ghOrg.List {
				deleteErr = s.repo.DisableRepository(item.RepositoryID)
				if deleteErr != nil {
					log.Warnf("Unable to remove repository: %s for project :%s error :%v", item.RepositoryID, projectID, deleteErr)
				}
			}
		}
	}

	return len(ghOrgs), nil
}

// GetRepositoriesByCLAGroup returns the list of repositories for the specified CLA Group
func (s *service) GetRepositoriesByCLAGroup(claGroupID string) ([]*models.GithubRepository, error) {
	// Return the list of github repositories that are enabled
	return s.repo.GetRepositoriesByCLAGroup(claGroupID, true)
}
