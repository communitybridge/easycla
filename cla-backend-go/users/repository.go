// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"fmt"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Repository interface defines the functions for the users service
type Repository interface {
	GetUser(userID string) (*models.User, error)
	GetUserByUserName(userName string) (*models.User, error)
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string) Repository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// DBUser data model
type DBUser struct {
	UserID             string   `json:"user_id"`
	LFEmail            string   `json:"lf_email"`
	LFUsername         string   `json:"lf_username"`
	DateCreated        string   `json:"date_created"`
	DateModified       string   `json:"date_modified"`
	UserName           string   `json:"user_name"`
	Version            string   `json:"version"`
	UserEmails         []string `json:"user_emails"`
	UserGithubID       string   `json:"user_github_id"`
	UserCompanyID      string   `json:"user_company_id"`
	UserGithubUsername string   `json:"user_github_username"`
}

func (repo repository) GetUser(userID string) (*models.User, error) {
	queryStartTime := time.Now()

	tableName := fmt.Sprintf("cla-%s-users", repo.stage)

	// This is the key we want to match
	condition := expression.Key("user_id").Equal(expression.Value(userID))

	// These are the columns we want returned
	projection := buildUserProjection()

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for user_id : %s, error: %v", userID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		//IndexName:                 aws.String(fmt.Sprintf("cla-%s-users", repo.stage)),
	}

	//log.Debugf("Running user query using queryInput: %+v", queryInput)

	// Make the DynamoDB Query API call
	result, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.Warnf("Error retrieving user by user_id: %s, error: %+v", userID, err)
		return nil, err
	}

	//log.Debugf("Result count : %d", *result.Count)
	//log.Debugf("Scanned count: %d", *result.ScannedCount)
	//log.Debugf("Result: %+v", *result)

	// The user model
	var dbUserModels []DBUser

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dbUserModels)
	if err != nil {
		log.Warnf("error unmarshalling user record from database for user_id: %s, error: %+v", userID, err)
		return nil, err
	}

	log.Debugf("GetUser by ID query took: %v resulting in %d results",
		utils.FmtDuration(time.Since(queryStartTime)), len(dbUserModels))

	if len(dbUserModels) == 0 {
		return nil, nil
	} else if len(dbUserModels) > 1 {
		log.Warnf("retrieved %d results for the getUser(id) query when we should return 0 or 1", len(dbUserModels))
	}

	//log.Debugf("%+v", dbUserModels[0])
	return convertDBUserModel(dbUserModels[0]), nil
}

func (repo repository) GetUserByUserName(userName string) (*models.User, error) {
	queryStartTime := time.Now()

	tableName := fmt.Sprintf("cla-%s-users", repo.stage)

	// This is the filter we want to match
	filter := expression.Name("lf_username").Equal(expression.Value(userName))

	// These are the columns we want returned
	projection := buildUserProjection()

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for user name: %s, error: %v", userName, err)
		return nil, err
	}

	// Assemble the scan input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	//log.Debugf("Running user query using scanInput: %+v", scanInput)

	var lastEvaluatedKey string
	// The user model
	var dbUserModels []DBUser

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		result, err := repo.dynamoDBClient.Scan(scanInput)
		if err != nil {
			log.Warnf("Error retrieving user by user name: %s, error: %+v", userName, err)
			return nil, err
		}

		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &dbUserModels)
		if err != nil {
			log.Warnf("error unmarshalling user record from database for user name: %s, error: %+v", userName, err)
			return nil, err
		}

		log.Debugf("GetUser by User Name query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(dbUserModels))

		if len(dbUserModels) == 1 {
			return convertDBUserModel(dbUserModels[0]), nil
		} else if len(dbUserModels) > 1 {
			log.Warnf("retrieved %d results for the getUser(id) query when we should return 0 or 1", len(dbUserModels))
			return convertDBUserModel(dbUserModels[0]), nil
		}

		// If we have another page of results...
		if result.LastEvaluatedKey["user_id"] != nil {
			lastEvaluatedKey = *result.LastEvaluatedKey["user_id"].S
			scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
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

// convertDBUserModel translates a dyanamoDB data model into a service response model
func convertDBUserModel(user DBUser) *models.User {
	return &models.User{
		UserID:         user.UserID,
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
	}
}

func buildUserProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("user_id"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("lf_email"),
		expression.Name("lf_username"),
		expression.Name("user_emails"),
		expression.Name("user_name"),
		expression.Name("user_company_id"),
		expression.Name("user_github_username"),
		expression.Name("user_github_id"),
	)
}
