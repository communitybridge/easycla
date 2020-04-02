// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/gofrs/uuid"

	"github.com/aws/aws-sdk-go/aws/awserr"

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
	SFDCRepositoryIndex    = "sfdc-repository-index"
)

// Repository defines functions of Repositories
type Repository interface {
	GetMetrics() (*models.RepositoryMetrics, error)
	GetProjectRepositoriesGroupByOrgs(projectID string) ([]*models.GithubRepositoriesGroupByOrgs, error)
	DeleteRepositoriesOfGithubOrganization(externalProjectID, githubOrgName string) error
	AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error)
	DeleteGithubRepository(externalProjectID string, repositoryID string) error
	ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error)
}

// NewRepository create new Repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		stage:               stage,
		dynamoDBClient:      dynamodb.New(awsSession),
		repositoryTableName: fmt.Sprintf("cla-%s-repositories", stage),
	}
}

type repo struct {
	stage               string
	dynamoDBClient      *dynamodb.DynamoDB
	repositoryTableName string
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

// getRepositoriesByGithubOrg returns an array of GH repositories for the specified project ID
func (repo repo) getRepositoriesByGithubOrg(githubOrgName string) ([]*models.GithubRepository, error) {
	var out []*models.GithubRepository
	tableName := fmt.Sprintf("cla-%s-repositories", repo.stage)
	builder := expression.NewBuilder()
	filter := expression.Name("repository_organization_name").Equal(expression.Value(githubOrgName))
	builder = builder.WithFilter(filter)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	results, err := repo.dynamoDBClient.Scan(scanInput)
	if err != nil {
		log.Warnf("unable to get github organizations repositories. error = %s", err.Error())
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

func (repo repo) deleteGithubRepository(externalProjectID, ghRepoID string) error {
	tableName := fmt.Sprintf("cla-%s-repositories", repo.stage)
	_, err := repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#projectSFID": aws.String("repository_sfdc_id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":externalProjectID": {
				S: aws.String(externalProjectID),
			},
		},
		ConditionExpression: aws.String("#projectSFID = :externalProjectID"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository_id": {S: aws.String(ghRepoID)},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.New("project id not matching with github repositories project id")
			}
		}
		log.Error(fmt.Sprintf("error deleting github repository with id: %s", ghRepoID), err)
		return err
	}
	return nil
}

func (repo repo) DeleteRepositoriesOfGithubOrganization(externalProjectID, githubOrgName string) error {
	ghrepos, err := repo.getRepositoriesByGithubOrg(githubOrgName)
	if err != nil {
		return err
	}
	for _, ghrepo := range ghrepos {
		err = repo.deleteGithubRepository(externalProjectID, ghrepo.RepositoryID)
		if err != nil {
			return err
		}
	}
	return nil
}

// List github repositories of project by external/salesforce project id
func (repo repo) ListProjectRepositories(externalProjectID string) (*models.ListGithubRepositories, error) {
	out := &models.ListGithubRepositories{
		List: make([]*models.GithubRepository, 0),
	}
	var condition expression.KeyConditionBuilder
	builder := expression.NewBuilder()

	condition = expression.Key("repository_sfdc_id").Equal(expression.Value(externalProjectID))

	builder = builder.WithKeyCondition(condition)
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
		TableName:                 aws.String(repo.repositoryTableName),
		IndexName:                 aws.String(SFDCRepositoryIndex),
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
		out.List = append(out.List, gr.toModel())
	}
	return out, nil
}

func (repo repo) AddGithubRepository(externalProjectID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	_, currentTime := utils.CurrentTime()
	repoID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	repository := &GithubRepository{
		DateCreated:                currentTime,
		DateModified:               currentTime,
		RepositoryExternalID:       *input.RepositoryExternalID,
		RepositoryID:               repoID.String(),
		RepositoryName:             *input.RepositoryName,
		RepositoryOrganizationName: *input.RepositoryOrganizationName,
		RepositoryProjectID:        *input.RepositoryProjectID,
		RepositorySfdcID:           externalProjectID,
		RepositoryType:             *input.RepositoryType,
		RepositoryURL:              *input.RepositoryURL,
		Version:                    "v1",
	}
	av, err := dynamodbattribute.MarshalMap(repository)
	if err != nil {
		return nil, err
	}
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.repositoryTableName),
	})
	if err != nil {
		log.Error("cannot put github repository in dynamodb", err)
		return nil, err
	}
	return repository.toModel(), nil
}

func (repo repo) DeleteGithubRepository(externalProjectID string, repositoryID string) error {
	return repo.deleteGithubRepository(externalProjectID, repositoryID)
}
