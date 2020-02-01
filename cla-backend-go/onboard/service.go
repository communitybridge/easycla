// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"github.com/aws/aws-sdk-go/service/sns"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// CharSet is the default character set for sending messages
const CharSet = "UTF-8"

type service struct {
	repo             OnboardRepository
	snsClient        *sns.SNS
	snsEventTopicARN string
}

// Service interface defining the public functions
type Service interface { // nolint
	CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail string) (*models.OnboardClaManagerRequest, error)
	GetCLAManagerRequestsByLFID(lfid string) (*models.OnboardClaManagerRequests, error)
	DeleteCLAManagerRequestsByRequestID(requestID string) error
	SendNotification(recipients []string, subject *string, emailBody *string) error
}

// NewService creates a new company service object
func NewService(repo OnboardRepository, awsSession *session.Session) Service {
	return service{
		repo:      repo,
		snsClient: sns.New(awsSession),
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

// SendNotification sends the notification to the specified recipients
func (s service) SendNotification(recipients []string, subject *string, emailBody *string) error {
	err := utils.SendEmail(*subject, *emailBody, recipients)
	if err != nil {
		log.Warnf("Error publishing message to topic: %s, Error: %v", s.snsEventTopicARN, err)
		return err
	}
	return nil
}
