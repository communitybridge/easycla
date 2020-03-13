// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	ini "github.com/communitybridge/easycla/cla-backend-go/init"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/stretchr/testify/assert"
)

func TestEventsService(t *testing.T) {
	awsSession, err := ini.GetAWSSession()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Unable to load AWS session - Error: %v", err))
	}
	stage := "dev"

	claUser := &user.CLAUser{
		UserID:         "test",
		Name:           "Test User",
		Emails:         []string{"test@foo.com"},
		LFEmail:        "test@foo.com",
		LFUsername:     "testlf",
		LfidProvider:   user.Provider{},
		GithubProvider: user.Provider{},
		ProjectIDs:     []string{"project-1234"},
		ClaIDs:         []string{"1234"},
		CompanyIDs:     []string{"company-1234"},
	}
	eventsMockRepo := NewMockRepository(awsSession, stage)
	eventsService := NewService(eventsMockRepo)

	eventsService.CreateAuditEvent(CreateUser, claUser, "project-1234", "company-1234", "Audit event test", false)
	eventsSearch, err := eventsService.SearchEvents(&eventOps.SearchEventsParams{
		ProjectID: aws.String("project-1234"),
		CompanyID: aws.String("company-1234"),
	})
	assert.Nil(t, err, "Error is nil")
	assert.Equal(t, len(eventsSearch.Events), 1)
}
