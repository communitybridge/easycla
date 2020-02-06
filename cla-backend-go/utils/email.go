package utils

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/labstack/gommon/log"
)

// EmailSender contains method to send email
type EmailSender interface {
	SendEmail(subject string, body string, recipients []string) error
}

var emailSender EmailSender

// SetEmailSender sets up default email sender
func SetEmailSender(es EmailSender) {
	emailSender = es
}

type snsEmail struct {
	snsClient          *sns.SNS
	snsEventTopicARN   string
	senderEmailAddress string
}

// SetSnsEmailSender set sns as mechanism to send email
func SetSnsEmailSender(awsSession *session.Session, snsEventTopicARN string, senderEmailAddress string) {
	emailSender = &snsEmail{
		snsClient:          sns.New(awsSession),
		snsEventTopicARN:   snsEventTopicARN,
		senderEmailAddress: senderEmailAddress,
	}
}

func (s *snsEmail) SendEmail(subject string, body string, recipients []string) error {
	var awsRecipients = make([]*string, len(recipients))
	for i, recipient := range recipients {
		awsRecipients[i] = &recipient
	}

	event := CreateEventWrapper("cla-email-event")
	event.Data = ToEmailEvent(&s.senderEmailAddress, recipients, &subject, &body)

	b, err := event.MarshalBinary()
	if err != nil {
		log.Warnf("Unable to marshal event model")
		return err
	}

	log.Debugf("Sending SNS message '%s' to topic: '%s'", b, s.snsEventTopicARN)
	input := &sns.PublishInput{
		Message:  aws.String(string(b)),          // Required
		TopicArn: aws.String(s.snsEventTopicARN), // Required
	}

	sendResp, err := s.snsClient.Publish(input)
	if err != nil {
		log.Warnf("Error publishing message to topic: %s, Error: %v", s.snsEventTopicARN, err)
		return err
	}

	log.Debugf("Successfully sent SNS message. Response: %v", sendResp)
	return nil
}

// SendEmail function send email. It uses emailSender interface.
func SendEmail(subject string, body string, recipients []string) error {
	if emailSender == nil {
		return errors.New("email sender not set")
	}
	return emailSender.SendEmail(subject, body, recipients)
}
