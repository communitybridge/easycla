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
