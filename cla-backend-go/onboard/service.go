// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"github.com/aws/aws-sdk-go/service/sns"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
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
	SendNotification(sender *string, recipients []string, subject *string, emailBody *string) (*string, error)
}

// NewService creates a new company service object
func NewService(repo OnboardRepository, awsSession *session.Session, snsEventTopicARN string) Service {
	return service{
		repo:             repo,
		snsClient:        sns.New(awsSession),
		snsEventTopicARN: snsEventTopicARN,
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
func (s service) SendNotification(sender *string, recipients []string, subject *string, emailBody *string) (*string, error) {

	// Convert the string array to an array of string pointers
	var awsRecipients = make([]*string, len(recipients))
	for i, recipient := range recipients {
		awsRecipients[i] = &recipient
	}

	event := utils.CreateEventWrapper("cla-email-event")
	event.Data = utils.ToEmailEvent(sender, recipients, subject, emailBody)

	b, err := event.MarshalBinary()
	if err != nil {
		log.Warnf("Unable to marshal event model")
		return nil, err
	}

	log.Debugf("Sending SNS message '%s' to topic: '%s'", b, s.snsEventTopicARN)
	input := &sns.PublishInput{
		Message:  aws.String(string(b)),          // Required
		TopicArn: aws.String(s.snsEventTopicARN), // Required
	}

	sendResp, err := s.snsClient.Publish(input)
	if err != nil {
		log.Warnf("Error publishing message to topic: %s, Error: %v", s.snsEventTopicARN, err)
		return nil, err
	}

	log.Debugf("Successfully sent SNS message. Response: %v", sendResp)
	return sendResp.MessageId, nil
}
