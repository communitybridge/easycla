// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
)

// Service interface defines the project service methods/functions
type Service interface {
	CreateProject(project *models.Project) (*models.Project, error)
	GetProjects(params *project.GetProjectsParams) (*models.Projects, error)
	GetProjectByID(projectID string) (*models.Project, error)
	DeleteProject(projectID string) error
	UpdateProject(projectModel *models.Project) (*models.Project, error)
}

// service
type service struct {
	repo Repository
}

// NewService returns an instance of the project service
func NewService(projectRepo Repository) Service {
	return service{
		repo: projectRepo,
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
	return s.repo.GetProject(projectID)
}

// DeleteProject service method
func (s service) DeleteProject(projectID string) error {
	return s.repo.DeleteProject(projectID)
}

// UpdateProject service method
func (s service) UpdateProject(projectModel *models.Project) (*models.Project, error) {
	return s.repo.UpdateProject(projectModel)
}
