package whitelist

import (
	"encoding/json"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// events
const (
	CclaWhitelistRequestAdded   = "ccla_whitelist_request_added"
	CclaWhitelistRequestDeleted = "ccla_whitelist_request_deleted"
)

// CclaWhitelistRequestAddedData is event data for event ccla_whitelist_request_added
type CclaWhitelistRequestAddedData struct {
	CompanyID string `json:"company_id"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
}

// CclaWhitelistRequestDeletedData is event data for event ccla_whitelist_request_deleted
type CclaWhitelistRequestDeletedData struct {
	CompanyID string `json:"company_id"`
	ProjectID string `json:"project_id"`
	RequestID string `json:"request_id"`
}

func createEvent(eventService events.Service, eventType string, companyID string, projectID string, userID string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Debugf("unable to log event. error = %s", err.Error())
		return
	}
	event := models.Event{
		EventCompanyID: companyID,
		EventData:      string(b),
		EventProjectID: projectID,
		EventType:      eventType,
		UserID:         userID,
	}
	err = eventService.CreateEvent(event)
	if err != nil {
		log.Debugf("unable to log event. event = %#v . error = %s", event, err.Error())
		return
	}
}
