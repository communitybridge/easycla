// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	"github.com/stretchr/testify/assert"
)

func TestEventsService(t *testing.T) {

	mockRepo := NewMockRepository()
	eventsMockRepo := mockRepo
	combinedMockRepo := mockRepo
	eventsService := NewService(eventsMockRepo, combinedMockRepo)

	eventsService.LogEvent(&LogEventArgs{
		EventType: GithubOrganizationAdded,
		ProjectID: "project-1234",
		CompanyID: "company-1234",
		UserID:    "testUserID",
		EventData: &GithubOrganizationAddedEventData{GithubOrganizationName: "testorg"},
	})

	eventsSearch, err := eventsService.SearchEvents(&eventOps.SearchEventsParams{
		ProjectID: aws.String("project-1234"),
		CompanyID: aws.String("company-1234"),
	})
	assert.Nil(t, err, "Error is nil")
	assert.Equal(t, len(eventsSearch.Events), 1)
}
