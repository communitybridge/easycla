// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service interface defines methods of event service
type Service interface {
	LogEvent(args *LogEventArgs)
	SearchEvents(params *eventOps.SearchEventsParams) (*models.EventList, error)
	GetRecentEvents(paramPageSize *int64) (*models.EventList, error)
	GetRecentEventsForCompanyProject(companyID, projectSFID string, pageSize *int64) (*models.EventList, error)
	GetFoundationSFDCEvents(foundationSFDC string, paramPageSize *int64) (*models.EventList, error)
	GetProjectSFDCEvents(projectSFDC string, paramPageSize *int64) (*models.EventList, error)
}

// CombinedRepo contains the various methods of other repositories
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
func (s *service) GetRecentEventsForCompanyProject(companyID, projectSFID string, paramPageSize *int64) (*models.EventList, error) {
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if paramPageSize != nil {
		pageSize = *paramPageSize
	}
	return s.repo.GetRecentEventsForCompanyProject(companyID, projectSFID, pageSize)
}

// GetFoundationSFDCEvents returns the list of foundation events
func (s *service) GetFoundationSFDCEvents(foundationSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return s.repo.GetFoundationSFDCEvents(foundationSFDC, paramPageSize)
}

// GetProjectSFDCEvents returns the list of project events
func (s *service) GetProjectSFDCEvents(projectSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return s.repo.GetProjectSFDCEvents(projectSFDC, paramPageSize)
}

// LogEventArgs is argument to LogEvent function
// EventType, EventData are compulsory.
// One of LfUsername, UserID must be present
type LogEventArgs struct {
	EventType         string
	ProjectID         string
	ProjectModel      *models.Project
	CompanyID         string
	CompanyModel      *models.Company
	LfUsername        string
	UserID            string
	UserModel         *models.User
	ExternalProjectID string
	EventData         EventData
	userName          string
	projectName       string
	companyName       string
}

func (s *service) loadCompany(args *LogEventArgs) error {
	if args.CompanyModel != nil {
		args.companyName = args.CompanyModel.CompanyName
		args.CompanyID = args.CompanyModel.CompanyID
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
		args.ProjectID = args.ProjectModel.ProjectID
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
		args.userName = args.UserModel.Username
		args.UserID = args.UserModel.UserID
		args.LfUsername = args.UserModel.LfUsername
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
	args.userName = userModel.Username
	args.UserID = userModel.UserID
	args.LfUsername = userModel.LfUsername
	return nil
}

func (s *service) loadDetails(args *LogEventArgs) error {
	err := s.loadCompany(args)
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

// LogEvent logs the event in database
func (s *service) LogEvent(args *LogEventArgs) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic occured in CreateEvent", fmt.Errorf("%v", r))
		}
	}()
	if args == nil || args.EventType == "" || args.EventData == nil || (args.UserID == "" && args.LfUsername == "") {
		log.Warnf("invalid arguments to LogEvent: args %#v", args)
		return
	}
	err := s.loadDetails(args)
	if err != nil {
		log.Error("unable to load details for event", err)
		return
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
		UserName:               args.userName,
		LfUsername:             args.LfUsername,
	}
	err = s.repo.CreateEvent(&event)
	if err != nil {
		log.Error(fmt.Sprintf("unable to create event for args %#v", args), err)
	}
}
