// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gofrs/uuid"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	// Version is version of CclaWhitelistRequest
	Version = "v1"
	// StatusPending is status of CclaWhitelistRequest
	StatusPending = "pending"
)

// IRepository interface defines the functions for the whitelist service
type IRepository interface {
	AddCclaWhitelistRequest(company *models.Company, project *models.Project, user *models.User, requesterName, requesterEmail string) (string, error)
	GetCclaWhitelistRequest(requestID string) (*CLARequestModel, error)
	ApproveCclaWhitelistRequest(requestID string) error
	RejectCclaWhitelistRequest(requestID string) error
	ListCclaWhitelistRequest(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error)
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string) IRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// CclaWhitelistRequest data model
type CclaWhitelistRequest struct {
	RequestID          string   `dynamodbav:"request_id"`
	RequestStatus      string   `dynamodbav:"request_status"`
	CompanyID          string   `dynamodbav:"company_id"`
	CompanyName        string   `dynamodbav:"company_name"`
	ProjectID          string   `dynamodbav:"project_id"`
	ProjectName        string   `dynamodbav:"project_name"`
	UserID             string   `dynamodbav:"user_id"`
	UserEmails         []string `dynamodbav:"user_emails"`
	UserName           string   `dynamodbav:"user_name"`
	UserGithubID       string   `dynamodbav:"user_github_id"`
	UserGithubUsername string   `dynamodbav:"user_github_username"`
	DateCreated        string   `dynamodbav:"date_created"`
	DateModified       string   `dynamodbav:"date_modified"`
	Version            string   `dynamodbav:"version"`
}

// AddCclaWhitelistRequest adds the specified request
func (repo repository) AddCclaWhitelistRequest(company *models.Company, project *models.Project, user *models.User, requesterName, requesterEmail string) (string, error) {
	requestID, err := uuid.NewV4()
	status := "status:fail"

	if err != nil {
		log.Warnf("AddCclaWhitelistRequest - unable to generate a UUID for a whitelist request, error: %v", err)
		return status, err
	}

	currentTime := currentTime()
	input := &dynamodb.PutItemInput{
		Item:      map[string]*dynamodb.AttributeValue{},
		TableName: aws.String(fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)),
	}
	addStringAttribute(input.Item, "request_id", requestID.String())
	addStringAttribute(input.Item, "request_status", StatusPending)
	addStringAttribute(input.Item, "company_id", company.CompanyID)
	addStringAttribute(input.Item, "company_name", company.CompanyName)
	addStringAttribute(input.Item, "project_id", project.ProjectID)
	addStringAttribute(input.Item, "project_name", project.ProjectName)
	addStringAttribute(input.Item, "user_id", user.UserID)
	addStringSliceAttribute(input.Item, "user_emails", []string{requesterEmail})
	addStringAttribute(input.Item, "user_name", requesterName)
	addStringAttribute(input.Item, "user_github_id", user.GithubID)
	addStringAttribute(input.Item, "user_github_username", user.GithubUsername)
	addStringAttribute(input.Item, "date_created", currentTime)
	addStringAttribute(input.Item, "date_modified", currentTime)
	addStringAttribute(input.Item, "version", Version)

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("AddCclaWhitelistRequest - unable to create a new ccla whitelist request, error: %v", err)
		return status, err
	}

	// Load the new record - should be able to find it quickly
	record, readErr := repo.ListCclaWhitelistRequest(company.CompanyID, &project.ProjectID, &user.UserID)
	if readErr != nil || record == nil || record.List == nil {
		log.Warnf("AddCclaWhitelistRequest - unable to read newly created invite record, error: %v", readErr)
		return status, err
	}

	return record.List[0].RequestID, nil
}

// GetCclaWhitelistRequest fetches the specified request by ID
func (repo repository) GetCclaWhitelistRequest(requestID string) (*CLARequestModel, error) {
	tableName := fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)

	response, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
	})

	if err != nil {
		log.Warnf("error fetching request by ID: %s, error: %v", requestID, err)
		return nil, err
	}

	requestModel := CLARequestModel{}
	err = dynamodbattribute.UnmarshalMap(response.Item, &requestModel)
	if err != nil {
		log.Warnf("error unmarshalling %s table response model data, error: %v", tableName, err)
		return nil, err
	}

	return &requestModel, nil
}

// ApproveCclaWhitelistRequest approves the specified request
func (repo repository) ApproveCclaWhitelistRequest(requestID string) error {
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("request_status"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String("approved"),
			},
			":m": {
				S: aws.String(currentTime()),
			},
		},
		UpdateExpression: aws.String("SET #S = :s, #M = :m"),
		TableName:        aws.String(fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)),
	}

	_, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("ApproveCclaWhitelistRequest - unable to update approval request with approved status, error: %v",
			err)
		return err
	}

	return nil
}

// RejectCclaWhitelistRequest rejects the specified request
func (repo repository) RejectCclaWhitelistRequest(requestID string) error {
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("request_status"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String("rejected"),
			},
			":m": {
				S: aws.String(currentTime()),
			},
		},
		UpdateExpression: aws.String("SET #S = :s, #M = :m"),
		TableName:        aws.String(fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)),
	}

	_, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("RejectCclaWhitelistRequest - unable to update approval request with rejected status, error: %v",
			err)
		return err
	}

	return nil
}

// ListCclaWhitelistRequest list the requests for the specified query parameters
func (repo repository) ListCclaWhitelistRequest(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error) {
	if projectID == nil {
		return nil, errors.New("project ID can not be nil for ListCclaWhitelistRequest")
	}

	tableName := fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)

	// hashkey is company_id, range key is project_id
	indexName := "company-id-project-id-index"

	condition := expression.Key("company_id").Equal(expression.Value(companyID))
	projectExpression := expression.Key("project_id").Equal(expression.Value(projectID))
	condition = condition.And(projectExpression)

	builder := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection())

	var filter expression.ConditionBuilder
	var filterAdded bool

	// Add the user ID filter if provided
	if userID != nil {
		userFilterExpression := expression.Name("user_id").Equal(expression.Value(userID))
		filter = addConditionToFilter(filter, userFilterExpression, &filterAdded)
	}
	if filterAdded {
		builder = builder.WithFilter(filter)
	}

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
	}

	queryOutput, queryErr := repo.dynamoDBClient.Query(input)
	if queryErr != nil {
		log.Warnf("list requests error while querying, error: %+v", queryErr)
		return nil, queryErr
	}

	list, err := buildCclaWhitelistRequestsModels(queryOutput)
	if err != nil {
		log.Warnf("unmarshall requests error while decoding the response, error: %+v", err)
		return nil, err
	}

	return &models.CclaWhitelistRequestList{List: list}, nil
}

// buildProjects builds the response model projection for a given query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("request_id"),
		expression.Name("request_status"),
		expression.Name("company_id"),
		expression.Name("company_name"),
		expression.Name("project_id"),
		expression.Name("project_name"),
		expression.Name("user_id"),
		expression.Name("user_emails"),
		expression.Name("user_name"),
		expression.Name("user_github_id"),
		expression.Name("user_github_username"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
	)
}

// buildCclaWhitelistRequestsModels builds the request models
func buildCclaWhitelistRequestsModels(results *dynamodb.QueryOutput) ([]models.CclaWhitelistRequest, error) {
	requests := make([]models.CclaWhitelistRequest, 0)

	var itemRequests []CclaWhitelistRequest

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &itemRequests)
	if err != nil {
		log.Warnf("error unmarshalling CCLA Authorization Request from database, error: %v",
			err)
		return nil, err
	}
	for _, r := range itemRequests {
		requests = append(requests, models.CclaWhitelistRequest{
			CompanyID:          r.CompanyID,
			CompanyName:        r.CompanyName,
			DateCreated:        r.DateCreated,
			DateModified:       r.DateModified,
			ProjectID:          r.ProjectID,
			ProjectName:        r.ProjectName,
			RequestID:          r.RequestID,
			RequestStatus:      r.RequestStatus,
			UserEmails:         r.UserEmails,
			UserGithubID:       r.UserGithubID,
			UserGithubUsername: r.UserGithubUsername,
			UserID:             r.UserID,
			UserName:           r.UserName,
			Version:            r.Version,
		})
	}
	return requests, nil
}

// addStringAttribute adds the specified attribute as a string
func addStringAttribute(item map[string]*dynamodb.AttributeValue, key string, value string) {
	if value != "" {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
}

// addStringSliceAttribute adds the specified attribute as a string slice
func addStringSliceAttribute(item map[string]*dynamodb.AttributeValue, key string, value []string) {
	if len(value) > 0 {
		item[key] = &dynamodb.AttributeValue{SS: aws.StringSlice(value)}
	}
}

// addConditionToFilter - helper routine for adding a filter condition
func addConditionToFilter(filter expression.ConditionBuilder, cond expression.ConditionBuilder, filterAdded *bool) expression.ConditionBuilder {
	if !(*filterAdded) {
		*filterAdded = true
		filter = cond
	} else {
		filter = filter.And(cond)
	}
	return filter
}

// currentTime helper routine to return the date/time
func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}
