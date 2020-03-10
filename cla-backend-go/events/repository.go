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

	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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
	GetProject(projectID string) (*models.Project, error)
	GetCompany(companyID string) (*models.Company, error)
	GetUserByUserName(userName string, fullMatch bool) (*models.User, error)
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
	addAttribute(input.Item, "event_id", eventID.String())
	addAttribute(input.Item, "event_type", event.EventType)
	addAttribute(input.Item, "event_user_id", event.UserID)
	addAttribute(input.Item, "event_user_name", event.UserName)
	addAttribute(input.Item, "event_user_name_lower", strings.ToLower(event.UserName))
	addAttribute(input.Item, "event_time", currentTimeString)
	addAttribute(input.Item, "event_data", event.EventData)
	addAttribute(input.Item, "event_company_id", event.EventCompanyID)
	addAttribute(input.Item, "event_company_name", event.EventCompanyName)
	addAttribute(input.Item, "event_company_name_lower", strings.ToLower(event.EventCompanyName))
	addAttribute(input.Item, "event_project_id", event.EventProjectID)
	addAttribute(input.Item, "event_project_name", event.EventProjectName)
	addAttribute(input.Item, "event_project_name_lower", strings.ToLower(event.EventProjectName))
	input.Item["event_time_epoch"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(currentTime.Unix(), 10))}
	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Unable to create a new event, error: %v", err)
		return err
	}

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

// GetCompany returns a company based on the company ID
func (repo repository) GetCompany(companyID string) (*models.Company, error) {

	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)
	queryStartTime := time.Now()

	companyTableData, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"company_id": {
				S: aws.String(companyID),
			},
		},
	})

	if err != nil {
		log.Warnf(err.Error())
		log.Warnf("error fetching company table data using company id: %s, error: %v", companyID, err)
		return &models.Company{}, err
	}

	if len(companyTableData.Item) == 0 {
		return &models.Company{}, errors.New("company does not exist")
	}
	log.Debugf("Get company query took: %v", utils.FmtDuration(time.Since(queryStartTime)))

	company := models.Company{}
	err = dynamodbattribute.UnmarshalMap(companyTableData.Item, &company)
	if err != nil {
		log.Warnf("error unmarshalling company table data, error: %v", err)
		return &models.Company{}, err
	}

	return &company, nil
}

// GetProject looks up a project by id
func (repo repository) GetProject(projectID string) (*models.Project, error) {

	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(projectID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving project having ID : %s, error: %v", projectID, err)
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, errors.New("project does not exist")
	}

	var dbModel project.DBProjectModel
	err = dynamodbattribute.UnmarshalMap(result.Item, &dbModel)
	if err != nil {
		log.Warnf("error unmarshalling db project model, error: %+v", err)
		return nil, err
	}

	// Convert the database model to an API response model
	return buildProjectModel(dbModel), nil
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
		expression.Name("event_project_id"),
		expression.Name("event_project_name"),
		expression.Name("event_company_id"),
		expression.Name("event_company_name"),
		expression.Name("event_time"),
		expression.Name("event_time_epoch"),
		expression.Name("event_data"),
	)
}

// buildProjectModel maps the database model to the API response model
func buildProjectModel(dbModel project.DBProjectModel) *models.Project {
	return &models.Project{
		ProjectID:          dbModel.ProjectID,
		ProjectName:        dbModel.ProjectName,
		ProjectACL:         dbModel.ProjectACL,
		ProjectCCLAEnabled: dbModel.ProjectCclaEnabled,
		ProjectICLAEnabled: dbModel.ProjectIclaEnabled,
		ProjectExternalID:  dbModel.ProjectExternalID,
		DateCreated:        dbModel.DateCreated,
		DateModified:       dbModel.DateModified,
		Version:            dbModel.Version,
	}
}

// buildUserProject returns the list of columns we want from the result
func buildUserProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("user_id"),
		expression.Name("user_external_id"),
		expression.Name("user_company_id"),
		expression.Name("admin"),
		expression.Name("lf_email"),
		expression.Name("lf_username"),
		expression.Name("user_name"),
		expression.Name("user_emails"),
		expression.Name("user_github_username"),
		expression.Name("user_github_id"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
		expression.Name("note"),
	)
}

// convertDBUserModel translates a dyanamoDB data model into a service response model
func convertDBUserModel(user DBUser) *models.User {
	return &models.User{
		UserID:         user.UserID,
		UserExternalID: user.UserExternalID,
		Admin:          user.Admin,
		LfEmail:        user.LFEmail,
		LfUsername:     user.LFUsername,
		DateCreated:    user.DateCreated,
		DateModified:   user.DateModified,
		Username:       user.UserName,
		Version:        user.Version,
		Emails:         user.UserEmails,
		GithubID:       user.UserGithubID,
		CompanyID:      user.UserCompanyID,
		GithubUsername: user.UserGithubUsername,
		Note:           user.Note,
	}
}

// GetUser looks up a user by id
func (repo repository) GetUserByUserName(userName string, fullMatch bool) (*models.User, error) {
	tableName := fmt.Sprintf("cla-%s-users", repo.stage)
	var indexName string
	// This is the filter we want to match
	var condition expression.KeyConditionBuilder
	if strings.Contains(userName, "github:") {
		indexName = "github-user-index"
		// Username for Github comes in as github:123456, so we want to remove the initial string
		githubID, err := strconv.Atoi(strings.Replace(userName, "github:", "", 1))
		if err != nil {
			log.Warnf("Unable to convert Github ID to number: %s", err)
			return nil, err
		}
		condition = expression.Key("user_github_id").Equal(expression.Value(githubID))
	} else {
		indexName = "lf-username-index"
		condition = expression.Key("lf_username").Equal(expression.Value(userName))
	}
	// These are the columns we want returned
	projection := buildUserProjection()
	builder := expression.NewBuilder().WithProjection(projection)
	// This is the filter we want to match
	if fullMatch {
		filter := expression.Name("lf_username").Equal(expression.Value(userName))
		builder.WithFilter(filter)
	} else {
		filter := expression.Name("lf_username").Contains(userName)
		builder.WithFilter(filter)
	}
	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for user name: %s, error: %v", userName, err)
		return nil, err
	}
	// Assemble the scan input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
	}
	var lastEvaluatedKey string
	// The database user model
	var dbUserModels []DBUser
	// Loop until we find a match or exhausted all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		result, err := repo.dynamoDBClient.Query(queryInput)
		if err != nil {
			log.Warnf("Error retrieving user by user name: %s, error: %+v", userName, err)
			return nil, err
		}
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dbUserModels)
		if err != nil {
			log.Warnf("error unmarshalling user record from database for user name: %s, error: %+v", userName, err)
			return nil, err
		}
		if len(dbUserModels) == 1 {
			return convertDBUserModel(dbUserModels[0]), nil
		} else if len(dbUserModels) > 1 {
			log.Warnf("retrieved %d results for the getUser(id) query when we should return 0 or 1", len(dbUserModels))
			return convertDBUserModel(dbUserModels[0]), nil
		}
		// Didn't find a match so far...need to keep looking via the next page of data
		// If we have another page of results...
		if result.LastEvaluatedKey["user_id"] != nil {
			lastEvaluatedKey = *result.LastEvaluatedKey["user_id"].S
			queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"user_id": {
					S: aws.String(lastEvaluatedKey),
				},
			}
		} else {
			lastEvaluatedKey = ""
		}
	}
	return nil, nil
}
