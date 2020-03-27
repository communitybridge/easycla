// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

const (
	GithubOrgSFIDIndex = "github-org-sfid-index"
)

// Repository interface defines the functions for the whitelist service
type Repository interface {
	GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error)
}

type repository struct {
	stage              string
	dynamoDBClient     *dynamodb.DynamoDB
	githubOrgTableName string
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:              stage,
		dynamoDBClient:     dynamodb.New(awsSession),
		githubOrgTableName: fmt.Sprintf("cla-%s-github-orgs", stage),
	}
}
func (repo repository) GetGithubOrganizations(externalProjectID string) (*models.GithubOrganizations, error) {
	var condition expression.KeyConditionBuilder
	builder := expression.NewBuilder()

	condition = expression.Key("organization_sfid").Equal(expression.Value(externalProjectID))

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
		TableName:                 aws.String(repo.githubOrgTableName),
		IndexName:                 aws.String(GithubOrgSFIDIndex),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.Warnf("error retrieving github_organizations using organization_sfid. error = %s", err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		return &models.GithubOrganizations{
			List: []*models.GithubOrganization{},
		}, nil
	}
	var resultOutput []*GithubOrganization
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &resultOutput)
	if err != nil {
		return nil, err
	}
	return &models.GithubOrganizations{List: toModels(resultOutput)}, nil
}
