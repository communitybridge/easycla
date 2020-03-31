// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

// Service contains functions of GithubOrganizations service
type Service interface {
	GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error)
}

type service struct {
	repo Repository
}

// NewService creates a new githubOrganizations service
func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error) {
	return s.repo.GetGithubOrganizations(externalProjectID)
}
