// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	"github.com/go-openapi/strfmt"
)

// mockRepository data model
type mockRepository struct{}

var events []*models.Event

// NewMockRepository creates a new instance of the mock event repository
func NewMockRepository() *mockRepository {
	return &mockRepository{}
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

// GetFoundationSFDCEvents returns the list of foundation events
func (repo *mockRepository) GetFoundationSFDCEvents(foundationSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return nil, nil
}

// GetProjectSFDCEvents returns the list of project events
func (repo *mockRepository) GetProjectSFDCEvents(projectSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return nil, nil
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

func (repo *mockRepository) GetProjectByID(projectID string) (*models.Project, error) {
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

func (repo *mockRepository) GetUser(userID string) (*models.User, error) {
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

func (repo *mockRepository) GetRecentEventsForCompanyProject(companyID, projectID string, pageSize int64) (*models.EventList, error) {
	return &models.EventList{}, nil
}
