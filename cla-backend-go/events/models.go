// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

// Event data model
type Event struct {
	EventID          string `dynamodbav:"event_id"`
	EventType        string `dynamodbav:"event_type"`
	UserID           string `dynamodbav:"user_id"`
	UserName         string `dynamodbav:"user_name"`
	EventProjectID   string `dynamodbav:"event_project_id"`
	EventProjectName string `dynamodbav:"event_project_name"`
	EventCompanyID   string `dynamodbav:"event_company_id"`
	EventCompanyName string `dynamodbav:"event_company_name"`
	EventTime        string `dynamodbav:"event_time"`
	EventTimeEpoch   int64  `dynamodbav:"event_time_epoch"`
	EventData        string `dynamodbav:"event_data"`
}

// DBUser data model
type DBUser struct {
	UserID             string   `json:"user_id"`
	UserExternalID     string   `json:"user_external_id"`
	LFEmail            string   `json:"lf_email"`
	Admin              bool     `json:"admin"`
	LFUsername         string   `json:"lf_username"`
	DateCreated        string   `json:"date_created"`
	DateModified       string   `json:"date_modified"`
	UserName           string   `json:"user_name"`
	Version            string   `json:"version"`
	UserEmails         []string `json:"user_emails"`
	UserGithubID       string   `json:"user_github_id"`
	UserCompanyID      string   `json:"user_company_id"`
	UserGithubUsername string   `json:"user_github_username"`
	Note               string   `json:"note"`
}

func (e *Event) toEvent() *models.Event { //nolint
	return &models.Event{
		EventCompanyID:   e.EventCompanyID,
		EventCompanyName: e.EventCompanyName,
		EventData:        e.EventData,
		EventID:          e.EventID,
		EventProjectID:   e.EventProjectID,
		EventProjectName: e.EventProjectName,
		EventTime:        e.EventTime,
		EventType:        e.EventType,
		UserID:           e.UserID,
		UserName:         e.UserName,
		EventTimeEpoch:   e.EventTimeEpoch,
	}
}
