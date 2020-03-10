// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"time"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
)

// Service interface defines methods of event service
type Service interface {
	CreateEvent(event models.Event) error
	CreateAuditEvent(eventType string, claUser *user.CLAUser, projectID, companyID, data string)
	CreateAuditEventWithUserID(eventType string, userID, projectID, companyID, data string)
	SearchEvents(params *eventOps.SearchEventsParams) (*models.EventList, error)
	GetProject(projectID string) (*models.Project, error)
	GetCompany(companyID string) (*models.Company, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of event service
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateEvent(event models.Event) error {
	return s.repo.CreateEvent(&event)
}

// CreateAuditEventWithUserID creates an audit event record in the database
func (s *service) CreateAuditEventWithUserID(eventType, userID, projectID, companyID, data string) {
	var projectName = "not defined"
	if projectID != "" {
		projectModel, projectErr := s.GetProject(projectID)
		if projectErr != nil || projectModel == nil {
			log.Warnf("error looking up the project by id: %s, error: %+v", projectID, projectErr)
		}
		if projectModel != nil {
			projectName = projectModel.ProjectName
		}
	}
	var companyName = "not defined"
	if companyID != "" {
		companyModel, companyErr := s.GetCompany(companyID)
		if companyErr != nil || companyModel == nil {
			log.Warnf("error looking up the company by id: %s, error: %+v", companyID, companyErr)
		}
		if companyModel != nil {
			companyName = companyModel.CompanyName
		}
	}
	var userName = "not defined"
	if userID != "" {
		userModel, userErr := s.GetUser(userID)
		if userErr != nil || userModel == nil {
			log.Warnf("error looking up the user by id: %s, error: %+v", userID, userErr)
		}
		if userModel != nil {
			userName = userModel.Username
		}
	}
	// Create and log the event
	eventErr := s.CreateEvent(models.Event{
		UserID:           userID,
		UserName:         userName,
		EventProjectID:   projectID,
		EventProjectName: projectName,
		EventCompanyID:   companyID,
		EventCompanyName: companyName,
		EventType:        eventType,
		EventTime:        time.Now().UTC().Format(time.RFC3339),
		EventTimeEpoch:   time.Now().Unix(),
		EventData:        data,
	})
	if eventErr != nil {
		log.Warnf("error adding event type: %s by user %s (%s) for project: %s (%s) and company: %s (%s) to the event log, error: %v",
			eventType, userName, userID, projectName, projectID, companyName, companyID, eventErr)
	}
}

// CreateAuditEvent creates an audit event record in the database
func (s *service) CreateAuditEvent(eventType string, claUser *user.CLAUser, projectID, companyID, data string) {

	var projectName = "not defined"
	if projectID != "" {
		projectModel, projectErr := s.GetProject(projectID)
		if projectErr != nil || projectModel == nil {
			log.Warnf("error looking up the project by id: %s, error: %+v", projectID, projectErr)
		}
		if projectModel != nil {
			projectName = projectModel.ProjectName
		}
	}

	var companyName = "not defined"
	if companyID != "" {
		companyModel, companyErr := s.GetCompany(companyID)
		if companyErr != nil || companyModel == nil {
			log.Warnf("error looking up the company by id: %s, error: %+v", companyID, companyErr)
		}
		if companyModel != nil {
			companyName = companyModel.CompanyName
		}
	}

	// Create and log the event
	eventErr := s.CreateEvent(models.Event{
		UserID:           claUser.UserID,
		UserName:         claUser.Name,
		EventProjectID:   projectID,
		EventProjectName: projectName,
		EventCompanyID:   companyID,
		EventCompanyName: companyName,
		EventType:        eventType,
		EventTime:        time.Now().UTC().Format(time.RFC3339),
		EventTimeEpoch:   time.Now().Unix(),
		EventData:        data,
	})

	if eventErr != nil {
		log.Warnf("error adding event type: %s by user %s (%s) for project: %s (%s) and company: %s (%s) to the event log, error: %v",
			eventType, claUser.Name, claUser.UserID, projectName, projectID, companyName, companyID, eventErr)
	}
}

// SearchEvents service definition
func (s *service) SearchEvents(params *eventOps.SearchEventsParams) (*models.EventList, error) {
	const defaultPageSize int64 = 50
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	return s.repo.SearchEvents(params, pageSize)
}

// GetProject returns the project object based on the project id
func (s *service) GetProject(projectID string) (*models.Project, error) {
	return s.repo.GetProject(projectID)
}

// GetCompany returns the company object based on the company id
func (s *service) GetCompany(companyID string) (*models.Company, error) {
	return s.repo.GetCompany(companyID)
}

// GetUser returns the user object based on the user id
func (s *service) GetUser(userID string) (*models.User, error) {
	return s.repo.GetUserByUserName(userID, true)
}
