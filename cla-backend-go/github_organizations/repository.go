// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"fmt"
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/github"

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
	ghOrgList := buildGithubOrganizationListModels(resultOutput)
	return &models.GithubOrganizations{List: ghOrgList}, nil
}

func buildGithubOrganizationListModels(githubOrganizations []*GithubOrganization) []*models.GithubOrganization {
	ghOrgList := toModels(githubOrganizations)
	if len(ghOrgList) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(ghOrgList))
		for _, ghorganization := range ghOrgList {
			go func(ghorg *models.GithubOrganization) {
				defer wg.Done()
				ghorg.GithubInfo = &models.GithubOrganizationGithubInfo{}
				user, err := github.GetUserDetails(ghorg.OrganizationName)
				if err != nil {
					ghorg.GithubInfo.Error = err.Error()
				} else {
					ghorg.GithubInfo.Details = &models.GithubOrganizationGithubInfoDetails{
						Bio:     user.Bio,
						HTMLURL: user.HTMLURL,
						ID:      user.ID,
					}
				}
				ghorg.Repositories = &models.GithubOrganizationRepositories{
					List: make([]string, 0),
				}
				if ghorg.OrganizationInstallationID != 0 {
					list, err := github.GetInstallationRepositories(ghorg.OrganizationInstallationID)
					if err != nil {
						log.Warnf("unable to get repositories for installation id : %d", ghorg.OrganizationInstallationID)
						ghorg.Repositories.Error = err.Error()
						return
					}
					ghorg.Repositories.List = list
				}
			}(ghorganization)
		}
		wg.Wait()
	}
	return ghOrgList
}
