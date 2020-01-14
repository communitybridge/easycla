package whitelist

import (
	"encoding/json"

	"github.com/communitybridge/easycla/cla-backend-go/events"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// events
const (
	CclaWhitelistRequestAdded   = "ccla_whitelist_request_added"
	CclaWhitelistRequestDeleted = "ccla_whitelist_request_deleted"
)

type CclaWhitelistRequestAddedData struct {
	CompanyID string `json:"company_id"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
}

type CclaWhitelistRequestDeletedData struct {
	CompanyID string `json:"company_id"`
	ProjectID string `json:"project_id"`
	RequestID string `json:"request_id"`
}

func createEvent(eventService events.Service, eventType string, companyID string, projectID string, userID string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = eventService.CreateEvent(models.Event{
		EventCompanyID: companyID,
		EventData:      string(b),
		EventProjectID: projectID,
		EventType:      eventType,
		UserID:         userID,
	})
	return err
}
