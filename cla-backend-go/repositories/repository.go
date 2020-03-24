// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// index
const (
	ProjectRepositoryIndex = "project-repository-index"
)

// Repository defines functions of Repositories
type Repository interface {
	GetMetrics() (*models.RepositoryMetrics, error)
	GetProjectRepositoriesGroupByOrgs(projectID string) ([]*models.GithubRepositoriesGroupByOrgs, error)
}

// NewRepository create new Repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

type repo struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// GetMetrics returns the metrics for the github repository
func (repo repo) GetMetrics() (*models.RepositoryMetrics, error) {
	var out models.RepositoryMetrics
	tableName := fmt.Sprintf("cla-%s-repositories", repo.stage)
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count of repositories, error: %v", err)
		return nil, err
	}

	out.TotalCount = *describeTableResult.Table.ItemCount
	return &out, nil
}

// GetProjectRepositoriesGroupByOrgs returns a list of GH orgs by project id
func (repo repo) GetProjectRepositoriesGroupByOrgs(projectID string) ([]*models.GithubRepositoriesGroupByOrgs, error) {
	out := make([]*models.GithubRepositoriesGroupByOrgs, 0)
	outMap := make(map[string]*models.GithubRepositoriesGroupByOrgs)
	ghrepos, err := repo.getProjectRepositories(projectID)
	if err != nil {
		return nil, err
	}
	for _, ghrepo := range ghrepos {
		ghrepoGroup, ok := outMap[ghrepo.RepositoryOrganizationName]
		if !ok {
			ghrepoGroup = &models.GithubRepositoriesGroupByOrgs{
				OrganizationName: ghrepo.RepositoryOrganizationName,
			}
			out = append(out, ghrepoGroup)
			outMap[ghrepo.RepositoryOrganizationName] = ghrepoGroup
		}
		ghrepoGroup.List = append(ghrepoGroup.List, ghrepo)
	}
	return out, nil
}

// getProjectRepositories returns an array of GH repositories for the specified project ID
func (repo repo) getProjectRepositories(projectID string) ([]*models.GithubRepository, error) {
	var out []*models.GithubRepository
	tableName := fmt.Sprintf("cla-%s-repositories", repo.stage)
	var condition expression.KeyConditionBuilder
	builder := expression.NewBuilder()

	condition = expression.Key("repository_project_id").Equal(expression.Value(projectID))

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
		IndexName:                 aws.String(ProjectRepositoryIndex),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.Warnf("unable to get project github repositories. error = %s", err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		return out, nil
	}
	var result []*GithubRepository
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &result)
	if err != nil {
		return nil, err
	}
	for _, gr := range result {
		out = append(out, gr.toModel())
	}
	return out, nil
}
