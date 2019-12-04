// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/gofrs/uuid"
)

// OnboardRepository interface defines the functions for the github whitelist service
type OnboardRepository interface { // nolint
	CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail string) (*models.OnboardClaManagerRequest, error)
	GetCLAManagerRequestsByLFID(lfid string) (*models.OnboardClaManagerRequests, error)
	DeleteCLAManagerRequestsByRequestID(requestID string) error
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the onboard repository service
func NewRepository(awsSession *session.Session, stage string) OnboardRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func (repo repository) CreateCLAManagerRequest(lfid, projectName, companyName, userFullName, userEmail string) (*models.OnboardClaManagerRequest, error) {
	requestID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate a UUID for a new CLA Manager Request, error: %v", err)
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID.String()),
			},
			"lf_id": {
				S: aws.String(lfid),
			},
			"project_name": {
				S: aws.String(projectName),
			},
			"company_name": {
				S: aws.String(companyName),
			},
			"user_full_name": {
				S: aws.String(userFullName),
			},
			"user_email": {
				S: aws.String(userEmail),
			},
			"date_created": {
				S: aws.String(now),
			},
			"date_modified": {
				S: aws.String(now),
			},
		},
		TableName: aws.String(fmt.Sprintf("cla-%s-cla-manager-onboard-requests", repo.stage)),
	}

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Unable to create a new CLA Manager onboard request, error: %v", err)
		return nil, err
	}

	// Return the data model with the request id and times
	response := models.OnboardClaManagerRequest{
		RequestID:    requestID.String(),
		LfID:         aws.String(lfid),
		ProjectName:  aws.String(projectName),
		CompanyName:  aws.String(companyName),
		UserFullName: aws.String(userFullName),
		UserEmail:    aws.String(userEmail),
		DateCreated:  now,
		DateModified: now,
	}

	return &response, nil
}

func (repo repository) GetCLAManagerRequestsByLFID(lfid string) (*models.OnboardClaManagerRequests, error) {

	queryStartTime := time.Now()

	tableName := fmt.Sprintf("cla-%s-cla-manager-onboard-requests", repo.stage)
	condition := expression.Key("lf_id").Equal(expression.Value(lfid))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().
		WithKeyCondition(condition).
		WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for cla manager search by lfid: %s, error: %v", lfid, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("cla-manager-requests-lfid-index"), // Name of a secondary index
		//FilterExpression:          expr.Filter(),
	}

	var lastEvaluatedKey string
	var requests models.OnboardClaManagerRequests
	requests.Requests = []models.OnboardClaManagerRequest{}

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		queryResults, err := repo.dynamoDBClient.Query(queryInput)
		if err != nil {
			log.Warnf("Unable to retrieve data from %s using lfid: %s, error: %v", tableName, lfid, err)
			return nil, err
		}

		var managerRequests []models.OnboardClaManagerRequest
		err = dynamodbattribute.UnmarshalListOfMaps(queryResults.Items, &managerRequests)
		if err != nil {
			log.Warnf("error unmarshalling CLA manager requests for table %s using lfid: %s, error: %v",
				tableName, lfid, err)
			return nil, err
		}

		// Add to our response model
		requests.Requests = append(requests.Requests, managerRequests...)

		// Determine if we have more records - if so, update the start key and loop again
		if queryResults.LastEvaluatedKey["request_id"] != nil {
			lastEvaluatedKey = *queryResults.LastEvaluatedKey["request_id"].S
			queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"request_id": {
					S: aws.String(lastEvaluatedKey),
				},
			}
		} else {
			lastEvaluatedKey = ""
		}
	}

	log.Debugf("CLA Manager Request query with lfid %s took: %v with %d results", lfid,
		utils.FmtDuration(time.Since(queryStartTime)), len(requests.Requests))

	return &requests, nil
}

func (repo repository) DeleteCLAManagerRequestsByRequestID(requestID string) error {
	tableName := fmt.Sprintf("cla-%s-cla-manager-onboard-requests", repo.stage)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"request_id": {
				S: aws.String(requestID),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := repo.dynamoDBClient.DeleteItem(input)
	if err != nil {
		log.Warnf("Unable to delete CLA Manager Request by request id: %s, error: %v", requestID, err)
		return err
	}

	return nil
}

// buildProjection returns the list of columns for the query/scan projection
func buildProjection() expression.ProjectionBuilder {

	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("request_id"),
		expression.Name("lf_id"),
		expression.Name("project_name"),
		expression.Name("company_name"),
		expression.Name("user_full_name"),
		expression.Name("user_email"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
	)
}
