package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/gommon/log"

	"github.com/gofrs/uuid"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
)

// errors
var (
	ErrUserIDRequired    = errors.New("UserID cannot be empty")    //nolint
	ErrEventTypeRequired = errors.New("EventType cannot be empty") //nolint
)

// Repository interface defines methods of event repository service
type Repository interface {
	CreateEvent(event *models.Event) error
	SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error)
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the event service
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Create event will create event in database.
func (r *repository) CreateEvent(event *models.Event) error {
	if event.UserID == "" {
		return ErrUserIDRequired
	}
	if event.EventType == "" {
		return ErrEventTypeRequired
	}
	eventID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate a UUID for a whitelist request, error: %v", err)
		return err
	}

	currentTime := currentTime()
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"event_id": {
				S: aws.String(eventID.String()),
			},
			"event_type": {
				S: aws.String(event.EventType),
			},
			"user_id": {
				S: aws.String(event.UserID),
			},
			"event_time": {
				S: aws.String(currentTime),
			},
			"event_data": {
				S: aws.String(event.EventData),
			},
		},
		TableName: aws.String(fmt.Sprintf("cla-%s-events", r.stage)),
	}
	if event.EventCompanyID != "" {
		input.Item["event_company_id"] = &dynamodb.AttributeValue{S: aws.String(event.EventCompanyID)}
	}
	if event.EventProjectID != "" {
		input.Item["event_project_id"] = &dynamodb.AttributeValue{S: aws.String(event.EventProjectID)}
	}

	_, err = r.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Unable to create a new event, error: %v", err)
		return err
	}

	return nil
}

// SearchEvents returns list of events matching with filter criteria.
func (r *repository) SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error) {
	return &models.EventList{}, nil
}
