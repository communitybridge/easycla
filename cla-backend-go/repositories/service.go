// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

// Service contains functions of Github Repository service
type Service interface {
	AddGithubRepository(ctx context.Context, externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error)
	EnableRepository(ctx context.Context, repositoryID string) error
	DisableRepository(ctx context.Context, repositoryID string) error
	ListProjectRepositories(ctx context.Context, externalProjectID string) (*models.ListGithubRepositories, error)
	GetRepository(ctx context.Context, repositoryID string) (*models.GithubRepository, error)
	DisableRepositoriesByProjectID(ctx context.Context, projectID string) (int, error)
	GetRepositoriesByCLAGroup(ctx context.Context, claGroupID string) ([]*models.GithubRepository, error)
	GetRepositoriesByOrganizationName(ctx context.Context, gitHubOrgName string) ([]*models.GithubRepository, error)
}

// GithubOrgRepo provide method to get github organization by name
type GithubOrgRepo interface {
	GetGithubOrganizationByName(ctx context.Context, githubOrganizationName string) (*models.GithubOrganizations, error)
	GetGithubOrganization(ctx context.Context, githubOrganizationName string) (*models.GithubOrganization, error)
}

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetCLAGroupByID(ctx context.Context, projectID string, loadRepoDetails bool) (*models.ClaGroup, error)
}

type service struct {
	repo        Repository
	ghOrgRepo   GithubOrgRepo
	projectRepo ProjectRepo
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository, ghOrgRepo GithubOrgRepo, projectRepo ProjectRepo) Service {
	return &service{
		repo:        repo,
		ghOrgRepo:   ghOrgRepo,
		projectRepo: projectRepo,
	}
}

func (s *service) AddGithubRepository(ctx context.Context, externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	if input.RepositoryName != nil && *input.RepositoryName == "" {
		return nil, errors.New("github repository name required")
	}
	projectSFID := externalProjectID

	_, projErr := s.projectRepo.GetCLAGroupByID(ctx, externalProjectID, false)
	if projErr != nil {
		return nil, errors.New("provided repository projectID is invalid")
	}

	org, err := s.ghOrgRepo.GetGithubOrganizationByName(ctx, utils.StringValue(input.RepositoryOrganizationName))
	if err != nil {
		return nil, err
	}
	if len(org.List) == 0 {
		return nil, errors.New("github app not installed on github organization")
	}
	repoGithubID, err := strconv.ParseInt(utils.StringValue(input.RepositoryExternalID), 10, 64)
	if err != nil {
		return nil, err
	}
	ghRepo, err := github.GetRepositoryByExternalID(org.List[0].OrganizationInstallationID, repoGithubID)
	if err != nil {
		return nil, err
	}
	log.Debugf("ghRepo.HTMLURL %s, input.RepositoryURL  %s", *ghRepo.HTMLURL, *input.RepositoryURL)
	if !strings.EqualFold(*ghRepo.HTMLURL, *input.RepositoryURL) {
		return nil, errors.New("github repository not found")
	}
	return s.repo.AddGithubRepository(ctx, externalProjectID, projectSFID, input)
}

func (s *service) EnableRepository(ctx context.Context, repositoryID string) error {
	return s.repo.EnableRepository(ctx, repositoryID)
}

func (s *service) DisableRepository(ctx context.Context, repositoryID string) error {
	return s.repo.DisableRepository(ctx, repositoryID)
}

func (s *service) ListProjectRepositories(ctx context.Context, externalProjectID string) (*models.ListGithubRepositories, error) {
	return s.repo.ListProjectRepositories(ctx, externalProjectID, "", true)
}
func (s *service) GetRepository(ctx context.Context, repositoryID string) (*models.GithubRepository, error) {
	return s.repo.GetRepository(ctx, repositoryID)
}

// DisableRepositoriesByProjectID disables the repositories by project ID
func (s *service) DisableRepositoriesByProjectID(ctx context.Context, projectID string) (int, error) {
	var deleteErr error
	// Return the list of GitHub repositories by CLA Group for those that are currently enabled
	ghOrgs, err := s.repo.GetCLAGroupRepositoriesGroupByOrgs(ctx, projectID, true)
	if err != nil {
		return 0, err
	}
	if len(ghOrgs) > 0 {
		log.Debugf("Deleting repositories for project :%s", projectID)
		for _, ghOrg := range ghOrgs {
			for _, item := range ghOrg.List {
				deleteErr = s.repo.DisableRepository(ctx, item.RepositoryID)
				if deleteErr != nil {
					log.Warnf("Unable to remove repository: %s for project :%s error :%v", item.RepositoryID, projectID, deleteErr)
				}
			}
		}
	}

	return len(ghOrgs), nil
}

// GetRepositoriesByCLAGroup returns the list of repositories for the specified CLA Group
func (s *service) GetRepositoriesByCLAGroup(ctx context.Context, claGroupID string) ([]*models.GithubRepository, error) {
	// Return the list of github repositories that are enabled
	return s.repo.GetRepositoriesByCLAGroup(ctx, claGroupID, true)
}

// GetRepositoriesByOrganizationName get repositories by organization name
func (s *service) GetRepositoriesByOrganizationName(ctx context.Context, gitHubOrgName string) ([]*models.GithubRepository, error) {
	return s.repo.GetRepositoriesByOrganizationName(ctx, gitHubOrgName)
}
