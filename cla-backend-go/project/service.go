// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// Service interface defines the project service methods/functions
type Service interface {
	GetProjects() ([]models.Project, error)
	GetProjectByID(projectID string) (*models.Project, error)
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

func (s service) GetProjects() ([]models.Project, error) {
	return s.repo.GetProjects()
}

func (s service) GetProjectByID(projectID string) (*models.Project, error) {
	return s.repo.GetProject(projectID)
}
