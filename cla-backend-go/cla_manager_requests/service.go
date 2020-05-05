// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager_requests

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// IService interface defining the functions for the company service
type IService interface {
	CreateRequest(reqModel *CLAManagerRequest) (*models.ClaManagerRequest, error)
	GetRequests(companyID, projectID string) (*models.ClaManagerRequestList, error)
	GetRequest(requestID string) (*models.ClaManagerRequest, error)

	ApproveRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error)
	DenyRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error)
}

type service struct {
	repo IRepository
}

// NewService creates a new service object
func NewService(repo IRepository) IService {
	return service{
		repo: repo,
	}
}

// CreateRequest creates a request based on the specified parameters
func (s service) CreateRequest(reqModel *CLAManagerRequest) (*models.ClaManagerRequest, error) {
	request, err := s.repo.CreateRequest(reqModel)
	if err != nil {
		log.Warnf("problem with approving request for company ID: %s, project ID: %s, user ID: %s, user name: %s, error :%+v",
			reqModel.CompanyID, reqModel.ProjectID, reqModel.UserID, reqModel.UserName, err)
		return nil, err
	}

	return dbModelToServiceModel(request), err
}

// GetRequests returns a requests object based on the specified parameters
func (s service) GetRequests(companyID, projectID string) (*models.ClaManagerRequestList, error) {
	requests, err := s.repo.GetRequests(companyID, projectID)
	if err != nil {
		log.Warnf("problem with fetching request for company ID: %s, project ID: %s, error :%+v",
			companyID, projectID, err)
		return nil, err
	}

	// Convert to a service response model
	responseModel := models.ClaManagerRequestList{}
	for _, request := range requests.Requests {
		resp := dbModelToServiceModel(&request)
		responseModel.Requests = append(responseModel.Requests, *resp)
	}

	return &responseModel, nil
}

// GetRequest returns the request object based on the specified parameters
func (s service) GetRequest(requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.GetRequest(requestID)
	if err != nil {
		log.Warnf("problem with fetching request for request ID: %s, error :%+v",
			requestID, err)
		return nil, err
	}

	if request == nil {
		log.Debugf("request not found for Request ID: %s", requestID)
		return nil, nil
	}

	return dbModelToServiceModel(request), err
}

// ApproveRequest approves the request based on the specified parameters
func (s service) ApproveRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.ApproveRequest(companyID, projectID, requestID)
	if err != nil {
		log.Warnf("problem with approving request for company ID: %s, project ID: %s, request ID: %s, error :%+v",
			companyID, projectID, requestID, err)
		return nil, err
	}

	// Send email

	return dbModelToServiceModel(request), err
}

// DenyRequest denies the request based on the specified parameters
func (s service) DenyRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.DenyRequest(companyID, projectID, requestID)
	if err != nil {
		log.Warnf("problem with denying request for company ID: %s, project ID: %s, request ID: %s, error :%+v",
			companyID, projectID, requestID, err)
		return nil, err
	}

	// Send email

	return dbModelToServiceModel(request), err
}

// dbModelToServiceModel converts a database model to a service model
func dbModelToServiceModel(dbModel *CLAManagerRequest) *models.ClaManagerRequest {
	return &models.ClaManagerRequest{
		RequestID:         dbModel.RequestID,
		CompanyID:         dbModel.CompanyID,
		CompanyExternalID: dbModel.CompanyExternalID,
		CompanyName:       dbModel.CompanyName,
		ProjectID:         dbModel.ProjectID,
		ProjectExternalID: dbModel.ProjectExternalID,
		ProjectName:       dbModel.ProjectName,
		UserID:            dbModel.UserID,
		UserExternalID:    dbModel.UserExternalID,
		UserName:          dbModel.UserName,
		UserEmail:         dbModel.UserEmail,
		Status:            dbModel.Status,
		Created:           dbModel.Created,
		Updated:           dbModel.Updated,
	}
}
