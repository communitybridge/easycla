// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"errors"

	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
)

// Service interface defines methods of event service
type Service interface {
	LogEvent(args *LogEventArgs) error
	SearchEvents(params *eventOps.SearchEventsParams) (*models.EventList, error)
	GetRecentEvents(paramPageSize *int64) (*models.EventList, error)
}

type CombinedRepo interface {
	GetProjectByID(projectID string) (*models.Project, error)
	GetCompany(companyID string) (*models.Company, error)
	GetUserByUserName(userName string, fullMatch bool) (*models.User, error)
	GetUser(userID string) (*models.User, error)
}

type service struct {
	repo         Repository
	combinedRepo CombinedRepo
}

// NewService creates new instance of event service
func NewService(repo Repository, combinedRepo CombinedRepo) Service {
	return &service{
		repo:         repo,
		combinedRepo: combinedRepo,
	}
}

func (s *service) CreateEvent(event models.Event) error {
	return s.repo.CreateEvent(&event)
}

// CreateAuditEventWithUserID creates an audit event record in the database
func (s *service) CreateAuditEventWithUserID(eventType, userID, projectID, companyID, data string, containsPII bool) {
	/*
		var projectName = "not defined"
		var projectExternalID = "not defined"
		if projectID != "" {
			projectModel, projectErr := s.combinedRepo.GetProjectByID(projectID)
			if projectErr != nil || projectModel == nil {
				log.Warnf("error looking up the project by id: %s, error: %+v", projectID, projectErr)
			}
			if projectModel != nil {
				projectName = projectModel.ProjectName
				projectExternalID = projectModel.ProjectExternalID
			}
		}
		var companyName = "not defined"
		if companyID != "" {
			companyModel, companyErr := s.repo.GetCompany(companyID)
			if companyErr != nil || companyModel == nil {
				log.Warnf("error looking up the company by id: %s, error: %+v", companyID, companyErr)
			}
			if companyModel != nil {
				companyName = companyModel.CompanyName
			}
		}
		var userName = "not defined"
		if userID != "" {
			userModel, userErr := s.repo.GetUserByUserName(userID, true)
			if userErr != nil || userModel == nil {
				log.Warnf("error looking up the user by id: %s, error: %+v", userID, userErr)
			}
			if userModel != nil {
				userName = userModel.Username
			}
		}
		// Create and log the event
		eventErr := s.CreateEvent(models.Event{
			UserID:                 userID,
			UserName:               userName,
			EventProjectID:         projectID,
			EventProjectName:       projectName,
			EventCompanyID:         companyID,
			EventCompanyName:       companyName,
			EventType:              eventType,
			EventTime:              time.Now().UTC().Format(time.RFC3339),
			EventTimeEpoch:         time.Now().Unix(),
			EventData:              data,
			EventProjectExternalID: projectExternalID,
			ContainsPII:            containsPII,
		})
		if eventErr != nil {
			log.Warnf("error adding event type: %s by user %s (%s) for project: %s (%s) and company: %s (%s) to the event log, error: %v",
				eventType, userName, userID, projectName, projectID, companyName, companyID, eventErr)
		}
	*/
}

// CreateAuditEvent creates an audit event record in the database
func (s *service) CreateAuditEvent(eventType string, claUser *user.CLAUser, projectID, companyID, data string, containsPII bool) {
	/*

		var projectName = "not defined"
		var projectExternalID = "not defined"
		if projectID != "" {
			projectModel, projectErr := s.repo.GetProject(projectID)
			if projectErr != nil || projectModel == nil {
				log.Warnf("error looking up the project by id: %s, error: %+v", projectID, projectErr)
			}
			if projectModel != nil {
				projectName = projectModel.ProjectName
				projectExternalID = projectModel.ProjectExternalID
			}
		}

		var companyName = "not defined"
		if companyID != "" {
			companyModel, companyErr := s.repo.GetCompany(companyID)
			if companyErr != nil || companyModel == nil {
				log.Warnf("error looking up the company by id: %s, error: %+v", companyID, companyErr)
			}
			if companyModel != nil {
				companyName = companyModel.CompanyName
			}
		}

		// Create and log the event
		eventErr := s.CreateEvent(models.Event{
			UserID:                 claUser.UserID,
			UserName:               claUser.Name,
			EventProjectID:         projectID,
			EventProjectName:       projectName,
			EventCompanyID:         companyID,
			EventCompanyName:       companyName,
			EventType:              eventType,
			EventTime:              time.Now().UTC().Format(time.RFC3339),
			EventTimeEpoch:         time.Now().Unix(),
			EventData:              data,
			EventProjectExternalID: projectExternalID,
			ContainsPII:            containsPII,
		})

		if eventErr != nil {
			log.Warnf("error adding event type: %s by user %s (%s) for project: %s (%s) and company: %s (%s) to the event log, error: %v",
				eventType, claUser.Name, claUser.UserID, projectName, projectID, companyName, companyID, eventErr)
		}

	*/
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

// GetRecentEvents returns event list of recent events
func (s *service) GetRecentEvents(paramPageSize *int64) (*models.EventList, error) {
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if paramPageSize != nil {
		pageSize = *paramPageSize
	}
	return s.repo.GetRecentEvents(pageSize)
}

type LogEventArgs struct {
	EventType         string
	ProjectID         string
	ProjectModel      *models.Project
	CompanyID         string
	CompanyModel      *models.Company
	LfUsername        string
	UserID            string
	UserName          string
	UserModel         *models.User
	ExternalProjectID string
	EventData         EventData
	projectName       string
	companyName       string
}

func (s *service) loadCompany(args *LogEventArgs) error {
	if args.CompanyModel != nil {
		args.companyName = args.CompanyModel.CompanyName
		return nil
	}
	if args.CompanyID != "" {
		companyModel, err := s.combinedRepo.GetCompany(args.CompanyID)
		if err != nil {
			return err
		}
		args.CompanyModel = companyModel
		args.companyName = companyModel.CompanyName
	}
	return nil
}

func (s *service) loadProject(args *LogEventArgs) error {
	if args.ProjectModel != nil {
		args.projectName = args.ProjectModel.ProjectName
		args.ExternalProjectID = args.ProjectModel.ProjectExternalID
		return nil
	}
	if args.ProjectID != "" {
		projectModel, err := s.combinedRepo.GetProjectByID(args.ProjectID)
		if err != nil {
			return err
		}
		args.ProjectModel = projectModel
		args.projectName = projectModel.ProjectName
		args.ExternalProjectID = projectModel.ProjectExternalID
	}
	return nil
}

func (s *service) loadUser(args *LogEventArgs) error {
	if args.UserModel != nil {
		args.UserName = args.UserModel.Username
		args.UserID = args.UserModel.UserID
		return nil
	}
	if args.UserID == "" && args.LfUsername == "" {
		return errors.New("require userID or LfUsername")
	}
	var userModel *models.User
	var err error
	if args.LfUsername != "" {
		userModel, err = s.combinedRepo.GetUserByUserName(args.LfUsername, true)
		if err != nil {
			return err
		}
	}
	if args.UserID != "" {
		userModel, err = s.combinedRepo.GetUser(args.UserID)
		if err != nil {
			return err
		}
	}
	args.UserModel = userModel
	args.UserName = userModel.Username
	args.UserID = userModel.UserID
	return nil
}

func (s *service) loadDetails(args *LogEventArgs) error {
	var err error
	err = s.loadCompany(args)
	if err != nil {
		return err
	}
	err = s.loadProject(args)
	if err != nil {
		return err
	}
	err = s.loadUser(args)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) LogEvent(args *LogEventArgs) error {
	if args == nil {
		errors.New("invalid arguments to LogEvent")
	}
	err := s.loadDetails(args)
	if err != nil {
		return err
	}
	eventData, containsPII := args.EventData.GetEventString(args)
	event := models.Event{
		ContainsPII:            containsPII,
		EventCompanyID:         args.CompanyID,
		EventCompanyName:       args.companyName,
		EventData:              eventData,
		EventProjectExternalID: args.ExternalProjectID,
		EventProjectID:         args.ProjectID,
		EventProjectName:       args.projectName,
		EventType:              args.EventType,
		UserID:                 args.UserID,
		UserName:               args.UserName,
	}
	return s.repo.CreateEvent(&event)
}
