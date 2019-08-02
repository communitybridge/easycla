// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"context"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

var (
	userID = "<redacted>"
)

// Service interface defines the project service methods/functions
type Service interface {
	GetProjects(ctx context.Context) ([]models.Project, error)
	GetProjectByID(ctx context.Context, projectID string) (models.Project, error)
}

// service
type service struct {
	projectRepo Repository
}

// NewService returns an instance of the project service
func NewService(projectRepo Repository) service {
	return service{
		projectRepo: projectRepo,
	}
}

// GetProjects returns a list of projects
func (s service) GetProjects(ctx context.Context) ([]models.Project, error) {
	//TODO: how to get the user ID
	// projectIDs
	_, err := s.projectRepo.GetProjectIDsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// TODO: Get Projects from SFDC
	projects := make([]models.Project, 2)
	projects[0].Name = "CCLA & ICLA Project"
	projects[0].Description = "This is a test project with both a CCLA and ICLA"
	projects[0].LogoURL = "https://s3.amazonaws.com/cla-project-logo-staging/a092M00001F1Yv4QAF.png"
	projects[0].SfdcID = "456789"
	projects[1].Name = "ICLA Project"
	projects[1].Description = "This is a test project with a ICLA"
	projects[1].LogoURL = "https://s3.amazonaws.com/cla-project-logo-staging/a092M00001F1YvEQAV.png"
	projects[1].SfdcID = "123sfdc"

	return projects, nil
}

// GetProjectByID returns the project based on the specified project id value
func (s service) GetProjectByID(ctx context.Context, projectID string) (models.Project, error) {
	m := make(map[string]models.Project)
	m["456789"] = models.Project{Name: "CCLA & ICLA Project", Description: "This is a test project with both a CCLA and ICLA", LogoURL: "https://s3.amazonaws.com/cla-project-logo-staging/a092M00001F1Yv4QAF.png", SfdcID: "456789"}
	m["123sfdc"] = models.Project{Name: "ICLA Project", Description: "This is a test project with a ICLA", LogoURL: "https://s3.amazonaws.com/cla-project-logo-staging/a092M00001F1YvEQAV.png", SfdcID: "123sfdc"}

	return m[projectID], nil
}
