// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"errors"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// errors
var (
	ErrCclaWhitelistRequestAlreadyExists = errors.New("CCLA whiltelist request already exist")
)

type service struct {
	repo        Repository
	userRepo    users.UserRepository
	companyRepo company.CompanyRepository
	projectRepo project.ProjectRepository
	httpClient  *http.Client
}

// NewService creates a new whitelist service
func NewService(repo Repository, userRepo users.UserRepository, companyRepo company.CompanyRepository, projectRepo project.ProjectRepository, httpClient *http.Client) service {
	return service{
		repo:        repo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
		projectRepo: projectRepo,
		httpClient:  httpClient,
	}
}

func (s service) AddCclaWhitelistRequest(companyID string, projectID string, args models.CclaWhitelistRequestInput) (string, error) {
	list, err := s.repo.ListCclaWhitelistRequest(companyID, &projectID, &args.UserID)
	if err != nil {
		return "", err
	}
	if len(list.List) > 0 {
		return "", ErrCclaWhitelistRequestAlreadyExists
	}
	companyModel, err := s.companyRepo.GetCompany(companyID)
	if err != nil {
		return "", err
	}
	projectModel, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return "", err
	}
	userModel, err := s.userRepo.GetUser(args.UserID)
	if err != nil {
		return "", err
	}
	if userModel == nil {
		return "", errors.New("invalid user")
	}
	return s.repo.AddCclaWhitelistRequest(companyModel, projectModel, userModel)
}

// DeleteCclaWhitelistRequest is the handler for the Delete CLA Whitelist request
func (s service) DeleteCclaWhitelistRequest(requestID string) error {
	return s.repo.DeleteCclaWhitelistRequest(requestID)
}

// ListCclaWhitelistRequest is the handler for the list CLA Whitelist request
func (s service) ListCclaWhitelistRequest(companyID string, projectID *string) (*models.CclaWhitelistRequestList, error) {
	return s.repo.ListCclaWhitelistRequest(companyID, projectID, nil)
}
