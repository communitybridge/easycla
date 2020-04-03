package events

import (
	"log"
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	ini "github.com/communitybridge/easycla/cla-backend-go/init"
)

func Test_repository_CreateEvent(t *testing.T) {

	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Fatal(err)
	}
	repo := NewRepository(awsSession, "dev")
	err = repo.CreateEvent(&models.Event{
		EventCompanyID:   "company123",
		EventCompanyName: "Deal Company",
		EventData:        "This is data",
		EventProjectID:   "project123",
		EventProjectName: "Project Name",
		EventType:        "Event Type",
		UserID:           "ID of ther user",
		UserName:         "Prasanna Mahajan",
	})
	t.Log("err is")
	t.Log(err)

}
