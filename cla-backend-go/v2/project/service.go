package project

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/project"
)

// Service interface defines the v2 project service methods
type Service interface {
	GetCCLAProjectsByExternalID(params *project.GetCCLAProjectsByExternalIDParams) (*models.Projects, error)
}

// service
type service struct {
	repo Repository
}

// NewService returns an instance of v2 porject service
func NewService(projectRepo Repository) Service {
	return service{
		repo: projectRepo,
	}
}

// GetCCLAProjectsByExternalID returns a list of CCLA enabled projects based on externalID
func (s service) GetCCLAProjectsByExternalID(params *project.GetCCLAProjectsByExternalIDParams) (*models.Projects, error) {
	projects, err := s.repo.GetCCLAProjectsByExternalID(params)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
