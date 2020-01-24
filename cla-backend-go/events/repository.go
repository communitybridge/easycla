package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

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
	SearchEvents(ctx context.Context, params *events.SearchEventsParams, pageSize int64) (*models.EventList, error)
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

func addConditionToFilter(filter expression.ConditionBuilder, cond expression.ConditionBuilder, filterAdded *bool) expression.ConditionBuilder {
	if !(*filterAdded) {
		*filterAdded = true
		filter = cond
	} else {
		filter = filter.And(cond)
	}
	return filter
}

func createSearchEventFilter(pk string, sk string, params *events.SearchEventsParams) *expression.ConditionBuilder {
	var filter expression.ConditionBuilder
	var filterAdded bool
	if params.ProjectID != nil && "event_project_id" != pk && "event_project_id" != sk {
		filterExpression := expression.Name("event_project_id").Equal(expression.Value(params.ProjectID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.CompanyID != nil && "event_company_id" != pk && "event_company_id" != sk {
		filterExpression := expression.Name("event_company_id").Equal(expression.Value(params.CompanyID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.UserID != nil && "user_id" != pk && "user_id" != sk {
		filterExpression := expression.Name("user_id").Equal(expression.Value(params.UserID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.EventType != nil && "event_type" != pk && "event_type" != sk {
		filterExpression := expression.Name("event_type").Equal(expression.Value(params.EventType))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.After != nil && "event_time" != pk && "event_time" != sk {
		filterExpression := expression.Name("event_time").GreaterThanEqual(expression.Value(params.After))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.Before != nil && "event_time" != pk && "event_time" != sk {
		filterExpression := expression.Name("event_time").LessThanEqual(expression.Value(params.Before))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if filterAdded {
		return &filter
	}
	return nil
}

func addTimeExpression(keyCond expression.KeyConditionBuilder, params *events.SearchEventsParams) expression.KeyConditionBuilder {
	if params.Before != nil && params.After != nil {
		exp := expression.Key("event_time").Between(expression.Value(params.After), expression.Value(params.Before))
		return keyCond.And(exp)
	}
	if params.After != nil {
		exp := expression.Key("event_time").GreaterThanEqual(expression.Value(params.After))
		return keyCond.And(exp)
	}
	if params.Before != nil {
		exp := expression.Key("event_time").LessThanEqual(expression.Value(params.Before))
		return keyCond.And(exp)
	}
	return keyCond
}

// SearchEvents returns list of events matching with filter criteria.
func (r *repository) SearchEvents(ctx context.Context, params *events.SearchEventsParams, pageSize int64) (*models.EventList, error) {
	if params.ProjectID == nil && params.CompanyID == nil && params.EventType == nil && params.UserID == nil {
		return nil, errors.New("invalid request.")
	}
	var condition expression.KeyConditionBuilder
	var indexName, pk, sk string
	builder := expression.NewBuilder().WithProjection(buildProjection())
	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-events", r.stage)

	switch {
	case params.CompanyID != nil:
		// search by companyID
		indexName = "event-company-id-event-time-index"
		condition = expression.Key("event_company_id").Equal(expression.Value(params.CompanyID))
		pk = "event_company_id"
		condition = addTimeExpression(condition, params)
		sk = "event_time"
	case params.ProjectID != nil:
		// search by projectID
		indexName = "event-project-id-event-time-index"
		condition = expression.Key("event_project_id").Equal(expression.Value(params.ProjectID))
		pk = "event_project_id"
		condition = addTimeExpression(condition, params)
		sk = "event_time"
	case params.UserID != nil:
		// search by userID
		condition = expression.Key("user_id").Equal(expression.Value(params.UserID))
		indexName = "user-id-index"
		pk = "user_id"
	case params.EventType != nil:
		// serach by eventType
		indexName = "event-type-index"
		condition = expression.Key("event_type").Equal(expression.Value(params.EventType))
		pk = "event_type"
	}
	filter := createSearchEventFilter(pk, sk, params)
	if filter != nil {
		builder = builder.WithFilter(*filter)
	}

	builder = builder.WithKeyCondition(condition)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	// Assemble the query input parameters
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		Limit:                     aws.Int64(pageSize), // The maximum number of items to evaluate (not necessarily the number of matching items)
	}
	if params.SortOrder != nil && *params.SortOrder == "desc" {
		input.ScanIndexForward = aws.Bool(false)
	}

	if params.NextKey != nil {
		log.Debugf("Received a nextKey, value: %s", *params.NextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: params.NextKey,
			},
		}
	}

	log.Debugf("query = %s", input.String())
	queryOutput, err := r.dynamoDBClient.Query(input)
	if err != nil {
		return nil, err
	}
	log.Debugf("last_evaluated_key = %#v", queryOutput.LastEvaluatedKey)
	return buildEventListModels(queryOutput)
}

func buildEventListModels(results *dynamodb.QueryOutput) (*models.EventList, error) {
	var out models.EventList
	out.Events = make([]*models.Event, 0)

	var items []Event

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &items)
	if err != nil {
		log.Warnf("error unmarshalling events from database, error: %v",
			err)
		return nil, err
	}
	for _, e := range items {
		out.Events = append(out.Events, e.toEvent())
	}
	return &out, nil
}

func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("event_id"),
		expression.Name("event_type"),
		expression.Name("user_id"),
		expression.Name("event_project_id"),
		expression.Name("event_company_id"),
		expression.Name("event_time"),
		expression.Name("event_data"),
	)
}
