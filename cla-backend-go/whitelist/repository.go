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

// Repository interface defines the functions for the whitelist service
type Repository interface {
	DeleteGithubOrganizationFromWhitelist(claGroupID, githubOrganizationID string) error
	AddGithubOrganizationToWhitelist(claGroupID, githubOrganizationID string) error
	GetGithubOrganizationsFromWhitelist(claGroupID string) ([]models.GithubOrg, error)

	AddCclaWhitelistRequest(company *models.Company, project *models.Project, user *models.User) (string, error)
	DeleteCclaWhitelistRequest(requestID string) error
	ListCclaWhitelistRequest(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error)
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string) repository {
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

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// AddGithubOrganizationToWhitelist adds the specified GH organization to the whitelist
func (repo repository) AddGithubOrganizationToWhitelist(claGroupID, GithubOrganizationID string) error {
	// get item from dynamoDB table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	log.Debugf("querying database for github organization whitelist using claGroupID: %s", claGroupID)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(claGroupID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s and GH Org: %s, error: %v",
			claGroupID, GithubOrganizationID, err)
		return err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		log.Debugf("claGroupID: %s is missing the 'github_org_whitelist' column - will add", claGroupID)
		itemFromMap = &dynamodb.AttributeValue{}
	}

	// generate new List L without element to be deleted
	// if we find a org with the same id just return without updating the record
	var newList []*dynamodb.AttributeValue
	for _, element := range itemFromMap.L {
		newList = append(newList, element)
		if *element.S == GithubOrganizationID {
			log.Debugf("github organization already in the list - nothing to do, org id: %s",
				GithubOrganizationID)
			return nil
		}
	}

	// Add the organization to list
	log.Debugf("adding github organization to the list, org id: %s", GithubOrganizationID)
	newList = append(newList, &dynamodb.AttributeValue{
		S: aws.String(GithubOrganizationID),
	})

	// Update dynamoDB table
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(claGroupID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#L": aws.String("github_org_whitelist"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				L: newList,
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
	}

	log.Warnf("updating database record using claGroupID: %s with values: %v", claGroupID, newList)
	_, err = repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating white list, error: %v", err)
		return err
	}

	return nil
}

// DeleteGithubOrganizationFromWhitelist removes the specified GH organization from the whitelist
func (repo repository) DeleteGithubOrganizationFromWhitelist(CLAGroupID, GithubOrganizationID string) error {
	// get item from dynamoDB table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s and GH Org: %s, error: %v", CLAGroupID, GithubOrganizationID, err)
		return err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		return errors.New("no github_org_whitelist column")
	}

	// generate new List L without element to be deleted
	newList := []*dynamodb.AttributeValue{}
	for _, element := range itemFromMap.L {
		if *element.S != GithubOrganizationID {
			newList = append(newList, element)
		}
	}

	// update dynamodb table
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#L": aws.String("github_org_whitelist"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				L: newList,
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
	}

	_, err = repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating white list, error: %v", err)
		return err
	}

	return nil
}

// GetGithubOrganizationsFromWhitelist returns a list of GH organizations stored in the whitelist
func (repo repository) GetGithubOrganizationsFromWhitelist(CLAGroupID string) ([]models.GithubOrg, error) {
	// get item from dynamodb table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s, error: %v", CLAGroupID, err)
		return nil, err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		return nil, nil
	}

	orgs := []models.GithubOrg{}
	for _, org := range itemFromMap.L {
		selected := true
		orgs = append(orgs, models.GithubOrg{
			ID:       org.S,
			Selected: &selected,
		})
	}

	return orgs, nil
}

func (repo repository) AddCclaWhitelistRequest(company *models.Company, project *models.Project, user *models.User) (string, error) {
	requestID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate a UUID for a whitelist request, error: %v", err)
		return "", err
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
	addStringSliceAttribute(input.Item, "user_emails", user.Emails)
	addStringAttribute(input.Item, "user_name", user.Username)
	addStringAttribute(input.Item, "user_github_id", user.GithubID)
	addStringAttribute(input.Item, "user_github_username", user.GithubUsername)
	addStringAttribute(input.Item, "date_created", currentTime)
	addStringAttribute(input.Item, "date_modified", currentTime)
	addStringAttribute(input.Item, "version", Version)

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Unable to create a new ccla whitelist request, error: %v", err)
		return "", err
	}

	return requestID.String(), nil
}

func (repo repository) DeleteCclaWhitelistRequest(requestID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
		TableName: aws.String(fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)),
	}

	_, err := repo.dynamoDBClient.DeleteItem(input)
	if err != nil {
		log.Warnf("Unable to delete ccla whitelist request, error: %v", err)
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

func (repo repository) ListCclaWhitelistRequest(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error) {
	tableName := fmt.Sprintf("cla-%s-ccla-whitelist-requests", repo.stage)

	indexName := "company-id-project-id-index"

	condition := expression.Key("company_id").Equal(expression.Value(companyID))

	builder := expression.NewBuilder().WithProjection(buildProjection())

	var filter expression.ConditionBuilder
	var filterAdded bool

	if userID != nil {
		userFilterExpression := expression.Name("user_id").Equal(expression.Value(userID))
		filter = addConditionToFilter(filter, userFilterExpression, &filterAdded)
	}

	if projectID != nil {
		projectExpression := expression.Key("project_id").Equal(expression.Value(projectID))
		condition = condition.And(projectExpression)
	}

	if filterAdded {
		builder = builder.WithFilter(filter)
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
	}
	queryOutput, err := repo.dynamoDBClient.Query(input)
	if err != nil {
		return nil, err
	}
	list, err := buildCclaWhitelistRequestsModels(queryOutput)
	if err != nil {
		return nil, err
	}
	return &models.CclaWhitelistRequestList{List: list}, nil
}

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
func buildCclaWhitelistRequestsModels(results *dynamodb.QueryOutput) ([]models.CclaWhitelistRequest, error) {
	requests := make([]models.CclaWhitelistRequest, 0)

	var itemRequests []CclaWhitelistRequest

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &itemRequests)
	if err != nil {
		log.Warnf("error unmarshalling ccla_whitelist_requests from database, error: %v",
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

func addStringAttribute(item map[string]*dynamodb.AttributeValue, key string, value string) {
	if value != "" {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
}
func addStringSliceAttribute(item map[string]*dynamodb.AttributeValue, key string, value []string) {
	if len(value) > 0 {
		item[key] = &dynamodb.AttributeValue{SS: aws.StringSlice(value)}
	}
}
