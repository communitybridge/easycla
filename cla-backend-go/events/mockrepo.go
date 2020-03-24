// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	"github.com/go-openapi/strfmt"
)

// MockRepository interface defines methods of event mock repository
type MockRepository interface {
	CreateEvent(event *models.Event) error
	SearchEvents(params *eventOps.SearchEventsParams, pageSize int64) (*models.EventList, error)
	GetProject(projectID string) (*models.Project, error)
	GetCompany(companyID string) (*models.Company, error)
	GetUserByUserName(userName string, fullMatch bool) (*models.User, error)
	GetRecentEvents(pageSize int64) (*models.EventList, error)
}

// mockRepository data model
type mockRepository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

var events []*models.Event

// NewMockRepository creates a new instance of the mock event repository
func NewMockRepository(awsSession *session.Session, stage string) MockRepository {
	return &mockRepository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func (repo *mockRepository) CreateEvent(event *models.Event) error {
	// Add to our in-memory list
	events = append(events, event)
	return nil
}

func (repo *mockRepository) SearchEvents(params *eventOps.SearchEventsParams, pageSize int64) (*models.EventList, error) {
	// Empty response model
	eventList := &models.EventList{
		Events:  []*models.Event{},
		NextKey: "",
	}

	if params.ProjectID != nil {
		for _, event := range events {
			if event.EventProjectID == *params.ProjectID && !eventInList(eventList.Events, event) {
				eventList.Events = append(eventList.Events, event)
			}
		}
	}

	if params.CompanyID != nil {
		for _, event := range events {
			if event.EventCompanyID == *params.CompanyID && !eventInList(eventList.Events, event) {
				eventList.Events = append(eventList.Events, event)
			}
		}
	}

	return eventList, nil
}

func eventInList(eventList []*models.Event, event *models.Event) bool {
	var retVal = false
	for _, theEvent := range eventList {
		if theEvent == event {
			retVal = true
			break
		}
	}

	return retVal
}

func (repo *mockRepository) GetProject(projectID string) (*models.Project, error) {
	return &models.Project{
		DateCreated:             "",
		DateModified:            "",
		ProjectACL:              nil,
		ProjectCCLAEnabled:      false,
		ProjectCCLARequiresICLA: false,
		ProjectExternalID:       "",
		ProjectICLAEnabled:      false,
		ProjectID:               "",
		ProjectName:             "",
		Version:                 "",
	}, nil

}

func (repo *mockRepository) GetCompany(companyID string) (*models.Company, error) {
	return &models.Company{
		CompanyACL:  []string{"foo", "bar"},
		CompanyID:   companyID,
		CompanyName: "Mock Company Name",
		Created:     strfmt.DateTime(time.Now().UTC()),
		Updated:     strfmt.DateTime(time.Now().UTC()),
	}, nil
}

func (repo *mockRepository) GetUserByUserName(userName string, fullMatch bool) (*models.User, error) {
	return &models.User{
		Admin:          true,
		CompanyID:      "mock_company_id",
		DateCreated:    "",
		DateModified:   "",
		Emails:         []string{"foo@gmail.com"},
		GithubID:       "",
		GithubUsername: "bar",
		LfEmail:        "foo@gmail.com",
		LfUsername:     "foo",
		Note:           "note",
		UserExternalID: "external",
		UserID:         "mock_user_id",
		Username:       "username",
		Version:        "v1",
	}, nil
}

func (repo *mockRepository) GetRecentEvents(pageSize int64) (*models.EventList, error) {
	return &models.EventList{}, nil
}
