// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"fmt"
	"strconv"

	"github.com/labstack/gommon/log"

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
	GetGithubOrganizations(projectSFID string) (*models.ProjectGithubOrganizations, error)
	AddGithubOrganization(projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(projectSFID string, githubOrgName string) error
	UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error
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

const (
	Connected         = "connected"
	PartialConnection = "partial_connection"
	ConnectionFailure = "connection_failure"
	NoConnection      = "no_connection"
)

func (s service) GetGithubOrganizations(projectSFID string) (*models.ProjectGithubOrganizations, error) {
	orgs, err := s.repo.GetGithubOrganizations("", projectSFID)
	if err != nil {
		return nil, err
	}
	out := &models.ProjectGithubOrganizations{
		List: make([]*models.ProjectGithubOrganization, 0),
	}
	type githubRepoInfo struct {
		orgName  string
		repoInfo *v1Models.GithubRepositoryInfo
	}
	// connectedRepo contains list of repositories for which github app have permission
	connectedRepo := make(map[string]*githubRepoInfo)
	orgmap := make(map[string]*models.ProjectGithubOrganization)
	for _, org := range orgs.List {
		for _, repoInfo := range org.Repositories.List {
			key := fmt.Sprintf("%s#%v", org.OrganizationName, repoInfo.RepositoryGithubID)
			connectedRepo[key] = &githubRepoInfo{
				orgName:  org.OrganizationName,
				repoInfo: repoInfo,
			}
		}
		rorg := &models.ProjectGithubOrganization{
			AutoEnabled:            org.AutoEnabled,
			ConnectionStatus:       "",
			GithubOrganizationName: org.OrganizationName,
			Repositories:           make([]*models.ProjectGithubRepository, 0),
		}
		orgmap[org.OrganizationName] = rorg
		out.List = append(out.List, rorg)
		if org.OrganizationInstallationID == 0 {
			rorg.ConnectionStatus = NoConnection
		} else {
			if org.Repositories.Error != "" {
				rorg.ConnectionStatus = ConnectionFailure
			} else {
				rorg.ConnectionStatus = Connected
			}
		}
	}
	repos, err := s.ghRepository.ListProjectRepositories("", projectSFID, true)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos.List {
		rorg, ok := orgmap[repo.RepositoryOrganizationName]
		if !ok {
			log.Warnf("repositories table contain stale data for organization %s", repo.RepositoryOrganizationName)
			continue
		}
		key := fmt.Sprintf("%s#%v", repo.RepositoryOrganizationName, repo.RepositoryExternalID)
		if _, ok := connectedRepo[key]; ok {
			repoGithubID, err := strconv.ParseInt(repo.RepositoryExternalID, 10, 64)
			if err != nil {
				log.Warnf("repository github id is not integer. error = %s", err)
			}
			rorg.Repositories = append(rorg.Repositories, &models.ProjectGithubRepository{
				ConnectionStatus:   Connected,
				Enabled:            true,
				RepositoryID:       repo.RepositoryID,
				RepositoryName:     repo.RepositoryName,
				RepositoryGithubID: repoGithubID,
			})
			// delete it from connectedRepo array since we have processed it
			// connectedArray after this loop will contain repo for which github app have permission but
			// they are enabled in cla
			delete(connectedRepo, key)
		} else {
			rorg.Repositories = append(rorg.Repositories, &models.ProjectGithubRepository{
				ConnectionStatus: ConnectionFailure,
				Enabled:          true,
				RepositoryID:     repo.RepositoryID,
				RepositoryName:   repo.RepositoryName,
			})
			if rorg.ConnectionStatus == Connected {
				rorg.ConnectionStatus = PartialConnection
			}
		}
	}
	for _, notEnabledRepo := range connectedRepo {
		rorg, ok := orgmap[notEnabledRepo.orgName]
		if !ok {
			log.Warnf("failed to get org %s", notEnabledRepo.orgName)
			continue
		}
		rorg.Repositories = append(rorg.Repositories, &models.ProjectGithubRepository{
			ConnectionStatus:   Connected,
			Enabled:            false,
			RepositoryID:       "",
			RepositoryName:     notEnabledRepo.repoInfo.RepositoryName,
			RepositoryGithubID: notEnabledRepo.repoInfo.RepositoryGithubID,
		})
	}
	return out, nil
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
	err := s.ghRepository.DisableRepositoriesOfGithubOrganization(projectSFID, githubOrgName)
	if err != nil {
		return err
	}
	return s.repo.DeleteGithubOrganization("", projectSFID, githubOrgName)
}

func (s service) UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error {
	return s.repo.UpdateGithubOrganization(projectSFID, organizationName, autoEnabled)
}
