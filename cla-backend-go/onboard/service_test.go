// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/sns"
)

type mockEmailSender struct {
	argSubject    string
	argBody       string
	argRecipients []string
	returnValue   error
}

func (s *mockEmailSender) SendEmail(subject string, body string, recipients []string) error {
	s.argSubject = subject
	s.argBody = body
	s.argRecipients = recipients
	return s.returnValue
}

func Test_service_SendNotification(t *testing.T) {
	type fields struct {
		repo             OnboardRepository
		snsClient        *sns.SNS
		snsEventTopicARN string
	}
	type args struct {
		recipients []string
		subject    *string
		emailBody  *string
	}
	su := "sub"
	bd := "body"
	tests := []struct {
		name        string
		fields      fields
		emailSender *mockEmailSender
		args        args
		want        *string
	}{
		{
			fields: fields{},
			emailSender: &mockEmailSender{
				returnValue: nil,
			},
			args: args{
				recipients: []string{"toaddr"},
				subject:    &su,
				emailBody:  &bd,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repo: tt.fields.repo,
			}
			utils.SetEmailSender(tt.emailSender)
			err := s.SendNotification(tt.args.recipients, tt.args.subject, tt.args.emailBody)
			if err != tt.emailSender.returnValue {
				t.Errorf("service.SendNotification() error = %v, wantErr %v", err, tt.emailSender.returnValue)
				return
			}
			assert.Equal(t, tt.args.recipients, tt.emailSender.argRecipients)
			assert.Equal(t, *tt.args.subject, tt.emailSender.argSubject)
			assert.Equal(t, *tt.args.emailBody, tt.emailSender.argBody)
		})
	}
}
