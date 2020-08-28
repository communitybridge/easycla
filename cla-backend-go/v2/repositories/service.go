// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/github"

	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	v1Repositories "github.com/communitybridge/easycla/cla-backend-go/repositories"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
)

// Service contains functions of Github Repository service
type Service interface {
	AddGithubRepository(projectSFID string, input *models.GithubRepositoryInput) (*v1Models.GithubRepository, error)
	EnableRepository(repositoryID string) error
	DisableRepository(repositoryID string) error
	ListProjectRepositories(projectSFID string) (*v1Models.ListGithubRepositories, error)
	GetRepository(repositoryID string) (*v1Models.GithubRepository, error)
	DisableCLAGroupRepositories(claGroupID string) error
}

// GithubOrgRepo provide method to get github organization by name
type GithubOrgRepo interface {
	GetGithubOrganization(githubOrganizationName string) (*v1Models.GithubOrganization, error)
}

type service struct {
	repo                  v1Repositories.Repository
	projectsClaGroupsRepo projects_cla_groups.Repository
	ghOrgRepo             GithubOrgRepo
}

// NewService creates a new githubOrganizations service
func NewService(repo v1Repositories.Repository, pcgRepo projects_cla_groups.Repository, ghOrgRepo GithubOrgRepo) Service {
	return &service{
		repo:                  repo,
		projectsClaGroupsRepo: pcgRepo,
		ghOrgRepo:             ghOrgRepo,
	}
}

func (s *service) AddGithubRepository(projectSFID string, input *models.GithubRepositoryInput) (*v1Models.GithubRepository, error) {
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
	allMappings, err := s.projectsClaGroupsRepo.GetProjectsIdsForClaGroup(aws.StringValue(input.ClaGroupID))
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
		return nil, fmt.Errorf("provided cla group id %s is not linked to project sfid %s", utils.StringValue(input.ClaGroupID), projectSFID)
	}
	org, err := s.ghOrgRepo.GetGithubOrganization(utils.StringValue(input.GithubOrganizationName))
	if err != nil {
		return nil, err
	}
	if org.OrganizationInstallationID == 0 {
		return nil, errors.New("github app not installed on github organization")
	}
	repoGithubID, err := strconv.ParseInt(utils.StringValue(input.RepositoryGithubID), 10, 64)
	if err != nil {
		return nil, err
	}
	ghRepo, err := github.GetRepositoryByExternalID(org.OrganizationInstallationID, repoGithubID)
	if err != nil {
		return nil, err
	}
	in := &v1Models.GithubRepositoryInput{
		RepositoryExternalID:       input.RepositoryGithubID,
		RepositoryName:             ghRepo.FullName,
		RepositoryOrganizationName: input.GithubOrganizationName,
		RepositoryProjectID:        input.ClaGroupID,
		RepositoryType:             aws.String("github"),
		RepositoryURL:              ghRepo.HTMLURL,
	}
	return s.repo.AddGithubRepository(externalProjectID, projectSFID, in)
}

func (s *service) EnableRepository(repositoryID string) error {
	return s.repo.EnableRepository(repositoryID)
}

func (s *service) DisableRepository(repositoryID string) error {
	return s.repo.DisableRepository(repositoryID)
}

func (s *service) ListProjectRepositories(projectSFID string) (*v1Models.ListGithubRepositories, error) {
	return s.repo.ListProjectRepositories("", projectSFID, true)
}

func (s *service) GetRepository(repositoryID string) (*v1Models.GithubRepository, error) {
	return s.repo.GetRepository(repositoryID)
}

func (s *service) DisableCLAGroupRepositories(claGroupID string) error {
	var deleteErr error
	ghOrgs, err := s.repo.GetCLAGroupRepositoriesGroupByOrgs(claGroupID, true)
	if err != nil {
		return err
	}
	if len(ghOrgs) > 0 {
		log.Debugf("Deleting repositories for cla-group :%s", claGroupID)
		for _, ghOrg := range ghOrgs {
			for _, item := range ghOrg.List {
				deleteErr = s.repo.DisableRepository(item.RepositoryID)
				if deleteErr != nil {
					log.Warnf("Unable to remove repository: %s for project :%s error :%v", item.RepositoryID, claGroupID, deleteErr)
				}
			}
		}
	}
	return nil
}
