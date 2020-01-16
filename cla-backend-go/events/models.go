package events

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

// Event data model
type Event struct {
	EventID        string `dynamodbav:"event_id"`
	EventType      string `dynamodbav:"event_type"`
	UserID         string `dynamodbav:"user_id"`
	EventProjectID string `dynamodbav:"event_project_id"`
	EventCompanyID string `dynamodbav:"event_company_id"`
	EventTime      string `dynamodbav:"event_time"`
	EventData      string `dynamodbav:"event_data"`
}

func (e *Event) toEvent() *models.Event { //nolint
	return &models.Event{
		EventCompanyID: e.EventCompanyID,
		EventData:      e.EventData,
		EventID:        e.EventID,
		EventProjectID: e.EventProjectID,
		EventTime:      e.EventTime,
		EventType:      e.EventType,
		UserID:         e.UserID,
	}
}
