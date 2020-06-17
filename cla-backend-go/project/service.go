// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/repositories"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gerrits"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service interface defines the project service methods/functions
type Service interface {
	CreateProject(project *models.Project) (*models.Project, error)
	GetProjects(params *project.GetProjectsParams) (*models.Projects, error)
	GetProjectByID(projectID string) (*models.Project, error)
	GetProjectsByExternalSFID(projectSFID string) (*models.Projects, error)
	GetProjectsByExternalID(params *project.GetProjectsByExternalIDParams) (*models.Projects, error)
	GetProjectByName(projectName string) (*models.Project, error)
	DeleteProject(projectID string) error
	UpdateProject(projectModel *models.Project) (*models.Project, error)
	GetClaGroupsByFoundationSFID(foundationSFID string, loadRepoDetails bool) (*models.Projects, error)
}

// service
type service struct {
	repo             ProjectRepository
	repositoriesRepo repositories.Repository
	gerritRepo       gerrits.Repository
}

// NewService returns an instance of the project service
func NewService(projectRepo ProjectRepository, repositoriesRepo repositories.Repository, gerritRepo gerrits.Repository) Service {
	return service{
		repo:             projectRepo,
		repositoriesRepo: repositoriesRepo,
		gerritRepo:       gerritRepo,
	}
}

// CreateProject service method
func (s service) CreateProject(projectModel *models.Project) (*models.Project, error) {
	return s.repo.CreateProject(projectModel)
}

// GetProjects service method
func (s service) GetProjects(params *project.GetProjectsParams) (*models.Projects, error) {
	return s.repo.GetProjects(params)
}

// GetProjectByID service method
func (s service) GetProjectByID(projectID string) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	s.fillRepoInfo(project)
	return project, nil
}

// GetProjectsByExternalSFID returns a list of projects based on the external SFID parameter
func (s service) GetProjectsByExternalSFID(projectSFID string) (*models.Projects, error) {
	return s.GetProjectsByExternalID(&project.GetProjectsByExternalIDParams{
		HTTPRequest: nil,
		NextKey:     nil,
		PageSize:    nil,
		ProjectSFID: projectSFID,
	})
}

// GetProjectsByExternalID returns a list of projects based on the external ID parameters
func (s service) GetProjectsByExternalID(params *project.GetProjectsByExternalIDParams) (*models.Projects, error) {
	log.Debugf("Project Service Handler - GetProjectsByExternalID")
	projects, err := s.repo.GetProjectsByExternalID(params, LoadRepoDetails)
	if err != nil {
		return nil, err
	}
	numberOfProjects := len(projects.Projects)
	if numberOfProjects == 0 {
		return projects, nil
	}
	var wg sync.WaitGroup
	wg.Add(numberOfProjects)
	for i := range projects.Projects {
		go func(project *models.Project) {
			defer wg.Done()
			s.fillRepoInfo(project)
		}(&projects.Projects[i])
	}
	wg.Wait()
	return projects, nil
}

func (s service) fillRepoInfo(project *models.Project) {
	var wg sync.WaitGroup
	wg.Add(2)
	var ghrepos []*models.GithubRepositoriesGroupByOrgs
	var gerrits []*models.Gerrit
	go func() {
		defer wg.Done()
		var err error
		ghrepos, err = s.repositoriesRepo.GetProjectRepositoriesGroupByOrgs(project.ProjectID)
		if err != nil {
			log.Error("unable to get github repositories for project.", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		gerrits, err = s.gerritRepo.GetProjectGerrits(project.ProjectID)
		if err != nil {
			log.Error("unable to get gerrit instances:w for project.", err)
			return
		}

	}()
	wg.Wait()
	project.GithubRepositories = ghrepos
	project.Gerrits = gerrits
}

// GetProjectByName service method
func (s service) GetProjectByName(projectName string) (*models.Project, error) {
	return s.repo.GetProjectByName(projectName)
}

// DeleteProject service method
func (s service) DeleteProject(projectID string) error {
	return s.repo.DeleteProject(projectID)
}

// UpdateProject service method
func (s service) UpdateProject(projectModel *models.Project) (*models.Project, error) {
	return s.repo.UpdateProject(projectModel)
}

// GetClaGroupsByFoundationSFID service method
func (s service) GetClaGroupsByFoundationSFID(foundationSFID string, loadRepoDetails bool) (*models.Projects, error) {
	return s.repo.GetClaGroupsByFoundationSFID(foundationSFID, loadRepoDetails)
}
