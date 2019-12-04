// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

type service struct {
	repo OnboardRepository
}

// Service interface defining the public functions
type Service interface { // nolint
	CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail string) (*models.OnboardClaManagerRequest, error)
	GetCLAManagerRequestsByLFID(lfid string) (*models.OnboardClaManagerRequests, error)
	DeleteCLAManagerRequestsByRequestID(requestID string) error
}

// NewService creates a new company service object
func NewService(repo OnboardRepository) Service {
	return service{
		repo: repo,
	}
}

// CreateCLAManagerRequest creates a new CLA manager request
func (s service) CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail string) (*models.OnboardClaManagerRequest, error) {
	return s.repo.CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail)
}

// GetCLAManagerRequestsByLFID get the CLA Manager requests by LFID
func (s service) GetCLAManagerRequestsByLFID(lfid string) (*models.OnboardClaManagerRequests, error) {
	return s.repo.GetCLAManagerRequestsByLFID(lfid)
}

// DeleteCLAManagerRequestsByRequestID invokes the backend to delete the request by ID
func (s service) DeleteCLAManagerRequestsByRequestID(requestID string) error {
	return s.repo.DeleteCLAManagerRequestsByRequestID(requestID)
}
