// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/gofrs/uuid"
)

// IRepository interface methods
type IRepository interface { //nolint
	CreateRequest(reqModel *CLAManagerRequest) (*CLAManagerRequest, error)
	GetRequests(companyID, projectID string) (*CLAManagerRequests, error)
	GetRequestsByUserID(companyID, projectID, userID string) (*CLAManagerRequests, error)
	GetRequest(requestID string) (*CLAManagerRequest, error)

	ApproveRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error)
	DenyRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error)
	PendingRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error)
	DeleteRequest(requestID string) error
	updateRequestStatus(companyID, projectID, requestID, status string) (*CLAManagerRequest, error)
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new company repository instance
func NewRepository(awsSession *session.Session, stage string) IRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// CreateRequest generates a new request
func (repo repository) CreateRequest(reqModel *CLAManagerRequest) (*CLAManagerRequest, error) {
	tableName := fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)

	requestID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate a UUID for a pending invite, error: %v", err)
		return nil, err
	}

	_, now := utils.CurrentTime()

	log.Debugf("request model: %+v", reqModel)

	itemMap := map[string]*dynamodb.AttributeValue{
		"request_id": {
			S: aws.String(requestID.String()),
		},
		"company_id": {
			S: aws.String(reqModel.CompanyID),
		},
		"company_name": {
			S: aws.String(reqModel.CompanyName),
		},
		"project_id": {
			S: aws.String(reqModel.ProjectID),
		},
		"project_name": {
			S: aws.String(reqModel.ProjectName),
		},
		"user_id": {
			S: aws.String(reqModel.UserID),
		},
		"user_name": {
			S: aws.String(reqModel.UserName),
		},
		"user_email": {
			S: aws.String(reqModel.UserEmail),
		},
		"status": {
			S: aws.String("pending"),
		},
		"date_created": {
			S: aws.String(now),
		},
		"date_modified": {
			S: aws.String(now),
		},
	}

	// If provided the project external ID - add it
	if reqModel.ProjectExternalID != "" {
		itemMap["project_external_id"] = &dynamodb.AttributeValue{
			S: aws.String(reqModel.ProjectExternalID),
		}
	}
	// If provided the company project external - add it
	if reqModel.CompanyExternalID != "" {
		itemMap["company_external_id"] = &dynamodb.AttributeValue{
			S: aws.String(reqModel.CompanyExternalID),
		}
	}
	// If provided the user external ID - add it
	if reqModel.UserExternalID != "" {
		itemMap["user_external_id"] = &dynamodb.AttributeValue{
			S: aws.String(reqModel.UserExternalID),
		}
	}

	input := &dynamodb.PutItemInput{
		Item:      itemMap,
		TableName: aws.String(tableName),
	}

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("unable to create a new CLA Manager request, error: %v", err)
		return nil, err
	}

	// Load the created record
	createdRequest, err := repo.GetRequest(requestID.String())
	if err != nil || createdRequest == nil {
		log.Warnf("unable to query newly created CLA Manager request by id: %s, error: %v",
			requestID.String(), err)
		return nil, err
	}

	return createdRequest, nil
}

// GetRequests returns the requests by Company ID and Project ID
func (repo repository) GetRequests(companyID, projectID string) (*CLAManagerRequests, error) {
	tableName := fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)

	condition := expression.Key("company_id").Equal(expression.Value(companyID)).And(
		expression.Key("project_id").Equal(expression.Value(projectID)))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().
		WithKeyCondition(condition).
		WithProjection(buildRequestProjection()).
		Build()
	if err != nil {
		log.Warnf("error building expression for cla manager request query using company ID: %s, project ID: %s, error: %v",
			companyID, projectID, err)
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
		IndexName:                 aws.String("cla-manager-requests-company-project-index"),
	}

	results, errQuery := repo.dynamoDBClient.Query(queryInput)
	if errQuery != nil {
		log.Warnf("error running query for cla manager request query using company ID: %s, project ID: %s, error: %v",
			companyID, projectID, err)
		return nil, errQuery
	}

	var requests []CLAManagerRequest

	// Unmarshall the DB response
	unmarshallErr := dynamodbattribute.UnmarshalListOfMaps(results.Items, &requests)
	if unmarshallErr != nil {
		log.Warnf("error converting DB model cla manager request query using company ID: %s, project ID: %s, error: %v",
			companyID, projectID, unmarshallErr)
		return nil, unmarshallErr
	}

	return &CLAManagerRequests{
		Requests: requests,
	}, nil
}

// GetRequestsByUserID returns the requests by Company ID and Project ID and User ID
func (repo repository) GetRequestsByUserID(companyID, projectID, userID string) (*CLAManagerRequests, error) {
	tableName := fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)

	condition := expression.Key("company_id").Equal(expression.Value(companyID)).And(
		expression.Key("project_id").Equal(expression.Value(projectID)))

	filter := expression.Name("user_id").Contains(userID)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().
		WithKeyCondition(condition).
		WithFilter(filter).
		WithProjection(buildRequestProjection()).
		Build()
	if err != nil {
		log.Warnf("error building expression for cla manager request query using company ID: %s, project ID: %s, user ID: %s, error: %v",
			companyID, projectID, userID, err)
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
		IndexName:                 aws.String("cla-manager-requests-company-project-index"),
	}

	results, errQuery := repo.dynamoDBClient.Query(queryInput)
	if errQuery != nil {
		log.Warnf("error running query for cla manager request query using company ID: %s, project ID: %s, user ID: %s, error: %v",
			companyID, projectID, userID, err)
		return nil, errQuery
	}

	var requests []CLAManagerRequest

	// Unmarshall the DB response
	unmarshallErr := dynamodbattribute.UnmarshalListOfMaps(results.Items, &requests)
	if unmarshallErr != nil {
		log.Warnf("error converting DB model cla manager request query using company ID: %s, project ID: %s, user ID: %s, error: %v",
			companyID, projectID, userID, unmarshallErr)
		return nil, unmarshallErr
	}

	return &CLAManagerRequests{
		Requests: requests,
	}, nil
}

// GetRequest returns the request by Request ID
func (repo repository) GetRequest(requestID string) (*CLAManagerRequest, error) {
	tableName := fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().
		WithProjection(buildRequestProjection()).
		Build()
	if err != nil {
		log.Warnf("error building expression for cla manager request query using request ID: %s, error: %v",
			requestID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {S: aws.String(requestID)},
		},
		ProjectionExpression:     expr.Projection(),
		ExpressionAttributeNames: expr.Names(),
		TableName:                aws.String(tableName),
	}

	result, errQuery := repo.dynamoDBClient.GetItem(queryInput)
	if errQuery != nil {
		log.Warnf("error running query for cla manager request query using request ID: %s, error: %v",
			requestID, err)
		return nil, errQuery
	}

	// If no response...
	if result.Item == nil {
		return nil, nil
	}

	var request CLAManagerRequest

	// Unmarshall the DB response
	unmarshallErr := dynamodbattribute.UnmarshalMap(result.Item, &request)
	if unmarshallErr != nil {
		log.Warnf("error converting DB model cla manager request query using request ID: %s, error: %v",
			requestID, unmarshallErr)
		return nil, unmarshallErr
	}

	return &request, nil
}

// DeleteRequest deletes the request by Request ID
func (repo repository) DeleteRequest(requestID string) error {
	tableName := fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)

	_, err := repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {S: aws.String(requestID)},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.New("request ID does not exist")
			}
		}
		log.Error(fmt.Sprintf("error deleting request with id: %s", requestID), err)
		return err
	}
	return nil
}

// ApproveRequest approves the specified request
func (repo repository) updateRequestStatus(companyID, projectID, requestID, status string) (*CLAManagerRequest, error) {
	// First, let's check if we already have a previous request
	requestModel, err := repo.GetRequest(requestID)
	if err != nil || requestModel == nil {
		log.Warnf("CLA Manager updateRequestStatus - unable to locate previous request with request ID: %s, company ID: %s, project ID: %s, error: %v",
			requestID, companyID, projectID, err)
		return nil, err
	}

	_, now := utils.CurrentTime()

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("status"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(status),
			},
			":m": {
				S: aws.String(now),
			},
		},
		UpdateExpression: aws.String("SET #S = :s, #M = :m"),
		TableName:        aws.String(fmt.Sprintf("cla-%s-cla-manager-requests", repo.stage)),
	}

	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.Warnf("CLA Manager ApproveRequest - unable to update request with '%s' status for request ID: %s, company ID: %s, project ID: %s, error: %v",
			status, requestID, companyID, projectID, updateErr)
		return nil, updateErr
	}

	// Load the updated document and return it
	updatedRequestModel, err := repo.GetRequest(requestID)
	if err != nil {
		log.Warnf("CLA Manager updateRequestStatus - unable to locate previous request with request ID: %s, company ID: %s, project ID: %s, error: %v",
			requestID, companyID, projectID, err)
		return nil, err
	}

	return updatedRequestModel, nil
}

// ApproveRequest approves the specified request
func (repo repository) ApproveRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error) {
	return repo.updateRequestStatus(companyID, projectID, requestID, "approved")
}

// DenyRequest denies the specified request
func (repo repository) DenyRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error) {
	return repo.updateRequestStatus(companyID, projectID, requestID, "denied")
}

// PendingRequest updates the status of an existing request to pending
func (repo repository) PendingRequest(companyID, projectID, requestID string) (*CLAManagerRequest, error) {
	return repo.updateRequestStatus(companyID, projectID, requestID, "pending")
}

// buildRequestProjection returns the database field projection for the table
func buildRequestProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("request_id"),
		expression.Name("company_id"),
		//expression.Name("company_external_id"),
		expression.Name("company_name"),
		expression.Name("project_id"),
		expression.Name("project_external_id"),
		expression.Name("project_name"),
		expression.Name("user_id"),
		//expression.Name("user_external_id"),
		expression.Name("user_name"),
		expression.Name("user_email"),
		expression.Name("status"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
	)
}
