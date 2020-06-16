// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/gofrs/uuid"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
)

// errors
var (
	ErrUserIDRequired    = errors.New("UserID cannot be empty")    //nolint
	ErrEventTypeRequired = errors.New("EventType cannot be empty") //nolint
)

// Repository interface defines methods of event repository service
type Repository interface {
	CreateEvent(event *models.Event) error
	SearchEvents(params *eventOps.SearchEventsParams, pageSize int64) (*models.EventList, error)
	GetRecentEvents(pageSize int64) (*models.EventList, error)
	GetRecentEventsForCompanyProject(companyID, projectID string, pageSize int64) (*models.EventList, error)
	GetFoundationSFDCEvents(foundationSFDC string, paramPageSize *int64) (*models.EventList, error)
	GetProjectSFDCEvents(projectSFDC string, paramPageSize *int64) (*models.EventList, error)
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the event repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func toDateFormat(t time.Time) string {
	//DD-MM-YYYY format
	return t.Format("02-01-2006")
}

// Create event will create event in database.
func (repo *repository) CreateEvent(event *models.Event) error {
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

	currentTime, currentTimeString := utils.CurrentTime()
	input := &dynamodb.PutItemInput{
		Item:      map[string]*dynamodb.AttributeValue{},
		TableName: aws.String(fmt.Sprintf("cla-%s-events", repo.stage)),
	}
	eventDateAndContainsPII := fmt.Sprintf("%s#%t", toDateFormat(currentTime), event.ContainsPII)
	addAttribute(input.Item, "event_id", eventID.String())
	addAttribute(input.Item, "event_type", event.EventType)
	addAttribute(input.Item, "event_user_id", event.UserID)
	addAttribute(input.Item, "event_user_name", event.UserName)
	addAttribute(input.Item, "event_lf_username", event.LfUsername)
	addAttribute(input.Item, "event_user_name_lower", strings.ToLower(event.UserName))
	addAttribute(input.Item, "event_time", currentTimeString)
	addAttribute(input.Item, "event_data", event.EventData)
	addAttribute(input.Item, "event_company_id", event.EventCompanyID)
	addAttribute(input.Item, "event_company_name", event.EventCompanyName)
	addAttribute(input.Item, "event_company_name_lower", strings.ToLower(event.EventCompanyName))
	addAttribute(input.Item, "event_project_id", event.EventProjectID)
	addAttribute(input.Item, "event_project_name", event.EventProjectName)
	addAttribute(input.Item, "event_project_name_lower", strings.ToLower(event.EventProjectName))
	addAttribute(input.Item, "event_date", toDateFormat(currentTime))
	addAttribute(input.Item, "event_project_external_id", event.EventProjectExternalID)
	addAttribute(input.Item, "event_date_and_contains_pii", eventDateAndContainsPII)
	input.Item["contains_pii"] = &dynamodb.AttributeValue{BOOL: &event.ContainsPII}
	input.Item["event_time_epoch"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(currentTime.Unix(), 10))}
	if event.EventCompanyID != "" && event.EventProjectExternalID != "" {
		companyIDexternalProjectID := fmt.Sprintf("%s#%s", event.EventCompanyID, event.EventProjectExternalID)
		addAttribute(input.Item, "company_id_external_project_id", companyIDexternalProjectID)
	}

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Unable to create a new event, error: %v", err)
		return err
	}
	log.Printf("added event : %s", eventID.String())

	return nil
}

func addAttribute(item map[string]*dynamodb.AttributeValue, key string, value string) {
	if value != "" {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
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

// createSearchEventFilter creates the search event filter
func createSearchEventFilter(pk string, sk string, params *eventOps.SearchEventsParams) *expression.ConditionBuilder {
	var filter expression.ConditionBuilder
	var filterAdded bool
	if params.ProjectID != nil && "event_project_id" != pk && "event_project_id" != sk { //nolint
		filterExpression := expression.Name("event_project_id").Equal(expression.Value(params.ProjectID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.CompanyID != nil && "event_company_id" != pk && "event_company_id" != sk { //nolint
		filterExpression := expression.Name("event_company_id").Equal(expression.Value(params.CompanyID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.UserID != nil && "event_user_id" != pk && "event_user_id" != sk { //nolint
		filterExpression := expression.Name("event_user_id").Equal(expression.Value(params.UserID))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.EventType != nil && "event_type" != pk && "event_type" != sk { //nolint
		filterExpression := expression.Name("event_type").Equal(expression.Value(params.EventType))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.After != nil && "event_time_epoch" != pk && "event_time_epoch" != sk { //nolint
		filterExpression := expression.Name("event_time_epoch").GreaterThanEqual(expression.Value(params.After))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.Before != nil && "event_time_epoch" != pk && "event_time_epoch" != sk { //nolint
		filterExpression := expression.Name("event_time_epoch").LessThanEqual(expression.Value(params.Before))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.UserName != nil && "user_name_lower" != pk && "user_name_lower" != sk { //nolint
		filterExpression := expression.Name("user_name_lower").Contains(strings.ToLower(*params.UserName))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.CompanyName != nil && "event_company_name_lower" != pk && "event_company_name_lower" != sk { //nolint
		filterExpression := expression.Name("event_company_name_lower").Contains(strings.ToLower(*params.CompanyName))
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if params.SearchTerm != nil {
		filterExpression := expression.Name("event_data").Contains(*params.SearchTerm)
		filter = addConditionToFilter(filter, filterExpression, &filterAdded)
	}
	if filterAdded {
		return &filter
	}
	return nil
}

// addTimeExpression adds the time expression to the query
func addTimeExpression(keyCond expression.KeyConditionBuilder, params *eventOps.SearchEventsParams) expression.KeyConditionBuilder {
	if params.Before != nil && params.After != nil {
		exp := expression.Key("event_time_epoch").Between(expression.Value(params.After), expression.Value(params.Before))
		return keyCond.And(exp)
	}
	if params.After != nil {
		exp := expression.Key("event_time_epoch").GreaterThanEqual(expression.Value(params.After))
		return keyCond.And(exp)
	}
	if params.Before != nil {
		exp := expression.Key("event_time_epoch").LessThanEqual(expression.Value(params.Before))
		return keyCond.And(exp)
	}
	return keyCond
}

// SearchEvents returns list of events matching with filter criteria.
func (repo *repository) SearchEvents(params *eventOps.SearchEventsParams, pageSize int64) (*models.EventList, error) {
	if params.ProjectID == nil {
		return nil, errors.New("invalid request. projectID is compulsory")
	}
	var condition expression.KeyConditionBuilder
	var indexName, pk, sk string
	builder := expression.NewBuilder().WithProjection(buildProjection())
	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-events", repo.stage)

	switch {
	case params.ProjectID != nil:
		// search by projectID
		indexName = "event-project-id-event-time-epoch-index"
		condition = expression.Key("event_project_id").Equal(expression.Value(params.ProjectID))
		pk = "event_project_id"
		condition = addTimeExpression(condition, params)
		sk = "event_time_epoch"
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
	queryInput := &dynamodb.QueryInput{
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
		queryInput.ScanIndexForward = aws.Bool(false)
	}

	if params.NextKey != nil {
		log.Debugf("Received a nextKey, value: %s", *params.NextKey)
		queryInput.ExclusiveStartKey, err = fromString(*params.NextKey)
		if err != nil {
			return nil, err
		}
	}

	var lastEvaluatedKey string
	events := make([]*models.Event, 0)

	for ok := true; ok; ok = lastEvaluatedKey != "" {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving events. error = %s", errQuery.Error())
			return nil, errQuery
		}

		eventsList, modelErr := buildEventListModels(results)
		if modelErr != nil {
			return nil, modelErr
		}

		events = append(events, eventsList...)
		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		}
		lastEvaluatedKey, err = toString(results.LastEvaluatedKey)
		if err != nil {
			return nil, err
		}
		if int64(len(events)) >= pageSize {
			break
		}
	}

	return &models.EventList{
		Events:  events,
		NextKey: lastEvaluatedKey,
	}, nil
}

// GetFoundationSFDCEvents returns the list of foundation events
func (repo *repository) GetFoundationSFDCEvents(foundationSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return nil, nil
}

// GetProjectSFDCEvents returns the list of project events
func (repo *repository) GetProjectSFDCEvents(projectSFDC string, paramPageSize *int64) (*models.EventList, error) {
	return nil, nil
}

// toString encodes the map as a string
func toString(in map[string]*dynamodb.AttributeValue) (string, error) {
	if len(in) == 0 {
		return "", nil
	}
	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// fromString converts the string to a map
func fromString(str string) (map[string]*dynamodb.AttributeValue, error) {
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	var m map[string]*dynamodb.AttributeValue
	err = json.Unmarshal(sDec, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// buildEventListModel converts the query results to a list event models
func buildEventListModels(results *dynamodb.QueryOutput) ([]*models.Event, error) {
	events := make([]*models.Event, 0)

	var items []Event

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &items)
	if err != nil {
		log.Warnf("error unmarshalling events from database, error: %v",
			err)
		return nil, err
	}
	for _, e := range items {
		events = append(events, e.toEvent())
	}
	return events, nil
}

// buildProjection builds the query projection
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("event_id"),
		expression.Name("event_type"),
		expression.Name("event_user_id"),
		expression.Name("event_user_name"),
		expression.Name("event_lf_username"),
		expression.Name("event_project_id"),
		expression.Name("event_project_name"),
		expression.Name("event_company_id"),
		expression.Name("event_company_name"),
		expression.Name("event_time"),
		expression.Name("event_time_epoch"),
		expression.Name("event_data"),
		expression.Name("event_project_external_id"),
	)
}

func (repo repository) GetRecentEvents(pageSize int64) (*models.EventList, error) {
	ctime := time.Now()
	maxQueryDays := 30
	events := make([]*models.Event, 0)
	for queriedDays := 0; queriedDays < maxQueryDays; queriedDays++ {
		day := toDateFormat(ctime)
		eventList, err := repo.getEventByDay(day, false, pageSize)
		if err != nil {
			return nil, err
		}
		events = append(events, eventList...)
		if int64(len(events)) >= pageSize {
			events = events[0:pageSize]
			break
		}
		ctime = ctime.Add(-(24 * time.Hour))
	}

	return &models.EventList{
		Events: events,
	}, nil

}

func (repo repository) getEventByDay(day string, containsPII bool, pageSize int64) ([]*models.Event, error) {
	tableName := fmt.Sprintf("cla-%s-events", repo.stage)
	var condition expression.KeyConditionBuilder
	builder := expression.NewBuilder().WithProjection(buildProjection())

	indexName := "event-date-and-contains-pii-event-time-epoch-index"
	eventDateAndContainsPII := fmt.Sprintf("%s#%t", day, containsPII)
	filter := expression.Name("event_project_id").AttributeExists()
	condition = expression.Key("event_date_and_contains_pii").Equal(expression.Value(eventDateAndContainsPII))

	builder = builder.WithKeyCondition(condition).WithFilter(filter)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		Limit:                     aws.Int64(pageSize), // The maximum number of items to evaluate (not necessarily the number of matching items)
		ScanIndexForward:          aws.Bool(false),
	}

	events := make([]*models.Event, 0)

	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving events. error = %s", errQuery.Error())
			return nil, errQuery
		}

		eventsList, modelErr := buildEventListModels(results)
		if modelErr != nil {
			return nil, modelErr
		}

		if len(eventsList) > 0 {
			events = append(events, eventsList...)
		}
		if int64(len(events)) >= pageSize {
			return events[0:pageSize], nil
		}
		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return events, nil
}

func (repo repository) GetRecentEventsForCompanyProject(companyID, projectSFID string, pageSize int64) (*models.EventList, error) {
	key := fmt.Sprintf("%s#%s", companyID, projectSFID)
	indexName := "company-id-external-project-id-event-epoch-time-index"
	condition := expression.Key("company_id_external_project_id").Equal(expression.Value(key))
	events, err := repo.eventIndexQuery(indexName, condition, nil, pageSize, false)
	if err != nil {
		return nil, err
	}
	var out models.EventList
	out.Events = events
	return &out, nil
}

func (repo repository) eventIndexQuery(indexName string, condition expression.KeyConditionBuilder, filter *expression.ConditionBuilder, pageSize int64, scanIndexForward bool) ([]*models.Event, error) {
	tableName := fmt.Sprintf("cla-%s-events", repo.stage)
	builder := expression.NewBuilder().WithProjection(buildProjection())
	builder = builder.WithKeyCondition(condition)
	if filter != nil {
		builder = builder.WithFilter(*filter)
	}
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		Limit:                     aws.Int64(pageSize), // The maximum number of items to evaluate (not necessarily the number of matching items)
		ScanIndexForward:          aws.Bool(scanIndexForward),
	}

	events := make([]*models.Event, 0)

	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving events. error = %s", errQuery.Error())
			return nil, errQuery
		}

		eventsList, modelErr := buildEventListModels(results)
		if modelErr != nil {
			return nil, modelErr
		}

		if len(eventsList) > 0 {
			events = append(events, eventsList...)
		}
		if int64(len(events)) >= pageSize {
			events = events[0:pageSize]
			break
		}
		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return events, nil
}
