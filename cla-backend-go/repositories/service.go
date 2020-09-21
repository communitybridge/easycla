// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
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

// GithubOrgRepo provide method to get github organization by name
type GithubOrgRepo interface {
	GetGithubOrganizationByName(githubOrganizationName string) (*models.GithubOrganizations, error)
	GetGithubOrganization(githubOrganizationName string) (*models.GithubOrganization, error)
}

type service struct {
	repo                  Repository
	ghOrgRepo             GithubOrgRepo
	projectsClaGroupsRepo projects_cla_groups.Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository, ghOrgRepo GithubOrgRepo, pcgRepo projects_cla_groups.Repository) Service {
	return &service{
		repo:                  repo,
		ghOrgRepo:             ghOrgRepo,
		projectsClaGroupsRepo: pcgRepo,
	}
}

func (s *service) AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	if input.RepositoryName != nil && *input.RepositoryName == "" {
		return nil, errors.New("github repository name required")
	}
	projectSFID := externalProjectID

	allMappings, err := s.projectsClaGroupsRepo.GetProjectsIdsForClaGroup(aws.StringValue(input.RepositoryProjectID))
	if err != nil {
		return nil, err
	}
	var valid bool
	for _, cgm := range allMappings {
		if cgm.ProjectSFID == projectSFID || cgm.FoundationSFID == projectSFID {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("provided cla group id %s is not linked to project sfid %s", utils.StringValue(input.RepositoryProjectID), projectSFID)
	}

	org, err := s.ghOrgRepo.GetGithubOrganizationByName(utils.StringValue(input.RepositoryOrganizationName))
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
