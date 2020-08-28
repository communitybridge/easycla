// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

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
	ProjectRepositoryIndex                    = "project-repository-index"
	SFDCRepositoryIndex                       = "sfdc-repository-index"
	ExternalRepositoryIndex                   = "external-repository-index"
	ProjectSFIDRepositoryOrgnizationNameIndex = "project-sfid-repository-organization-name-index"
)

// errors
var (
	ErrGithubRepositoryNotFound = errors.New("github repository not found")
)

// Repository defines functions of Repositories
type Repository interface {
	AddGithubRepository(externalProjectID string, projectSFID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error)
	EnableRepository(repositoryID string) error
	DisableRepository(repositoryID string) error
	DisableRepositoriesByProjectID(projectID string) error
	DisableRepositoriesOfGithubOrganization(externalProjectID, githubOrgName string) error
	GetRepository(repositoryID string) (*models.GithubRepository, error)
	GetRepositoriesByCLAGroup(claGroup string, enabled bool) ([]*models.GithubRepository, error)
	GetCLAGroupRepositoriesGroupByOrgs(projectID string, enabled bool) ([]*models.GithubRepositoriesGroupByOrgs, error)
	ListProjectRepositories(externalProjectID string, projectSFID string, enabled bool) (*models.ListGithubRepositories, error)
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

// AddGithubRepository adds the specified repository
func (repo repo) AddGithubRepository(externalProjectID string, projectSFID string, input *models.GithubRepositoryInput) (*models.GithubRepository, error) {
	f := logrus.Fields{
		"functionName":               "AddGithubRepository",
		"externalProjectID":          externalProjectID,
		"projectSFID":                projectSFID,
		"repositoryName":             *input.RepositoryName,
		"repositoryOrganizationName": *input.RepositoryOrganizationName,
		"repositoryType":             *input.RepositoryType,
		"repositoryURL":              *input.RepositoryURL,
	}

	// Check first to see if the repository already exists
	_, err := repo.getRepositoryByGithubID(utils.StringValue(input.RepositoryExternalID), true)
	if err != nil {
		// Expecting Not found - no issue if not found - all other error we throw
		if err != ErrGithubRepositoryNotFound {
			return nil, err
		}
	} else {
		return nil, errors.New("github repository already exist")
	}

	_, currentTime := utils.CurrentTime()
	repoID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	repository := &RepositoryDBModel{
		DateCreated:                currentTime,
		DateModified:               currentTime,
		RepositoryExternalID:       utils.StringValue(input.RepositoryExternalID),
		RepositoryID:               repoID.String(),
		RepositoryName:             utils.StringValue(input.RepositoryName),
		RepositoryOrganizationName: utils.StringValue(input.RepositoryOrganizationName),
		RepositoryProjectID:        utils.StringValue(input.RepositoryProjectID),
		RepositorySfdcID:           externalProjectID,
		RepositoryType:             utils.StringValue(input.RepositoryType),
		RepositoryURL:              utils.StringValue(input.RepositoryURL),
		Enabled:                    true, // default is enabled
		Note:                       fmt.Sprintf("created on %s", currentTime),
		ProjectSFID:                projectSFID,
		Version:                    "v1",
	}
	av, err := dynamodbattribute.MarshalMap(repository)
	if err != nil {
		log.WithFields(f).Warnf("problem marshalling the input, error: %+v", err)
		return nil, err
	}

	log.WithFields(f).Debug("creating repository entry")
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.repositoryTableName),
	})
	if err != nil {
		log.WithFields(f).Warnf("cannot add github repository, error: %+v", err)
		return nil, err
	}

	return repository.toModel(), nil
}

// EnableRepository enables the repository entry
func (repo *repo) EnableRepository(repositoryID string) error {
	return repo.enableGithubRepository(repositoryID)
}

// DisableRepository disables the repository entry (we don't delete)
func (repo *repo) DisableRepository(repositoryID string) error {
	return repo.disableGithubRepository(repositoryID)
}

func (repo *repo) DisableRepositoriesByProjectID(projectID string) error {
	repoModels, err := repo.getProjectRepositories(projectID, true)
	if err != nil {
		return err
	}

	// For each model...
	for _, repoModel := range repoModels {
		disableErr := repo.DisableRepository(repoModel.RepositoryID)
		if disableErr != nil {
			return disableErr
		}
	}

	return nil
}

// DisableRepositoriesOfGithubOrganization disables the repositories under the GitHub organization
func (repo repo) DisableRepositoriesOfGithubOrganization(externalProjectID, githubOrgName string) error {
	repoModels, err := repo.getRepositoriesByGithubOrg(githubOrgName)
	if err != nil {
		return err
	}
	for _, repoModel := range repoModels {
		if repoModel.RepositoryExternalID == externalProjectID || repoModel.RepositorySfdcID == externalProjectID {
			err = repo.disableGithubRepository(repoModel.RepositoryID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetRepository by repository id
func (repo *repo) GetRepository(repositoryID string) (*models.GithubRepository, error) {
	f := logrus.Fields{
		"functionName": "GetRepository",
		"repositoryID": repositoryID,
	}
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.repositoryTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"repository_id": {
				S: aws.String(repositoryID),
			},
		},
	})
	if err != nil {
		log.WithFields(f).Warn("problem querying using repository ID")
		return nil, err
	}
	if len(result.Item) == 0 {
		log.WithFields(f).Warn("repository with ID does not exist")
		return nil, ErrGithubRepositoryNotFound
	}

	var out RepositoryDBModel
	err = dynamodbattribute.UnmarshalMap(result.Item, &out)
	if err != nil {
		log.WithFields(f).Warn("problem unmarshalling response")
		return nil, err
	}

	return out.toModel(), nil
}

// GetRepositoryByCLAGroup gets the list of repositories based on the CLA Group ID
func (repo *repo) GetRepositoriesByCLAGroup(claGroupID string, enabled bool) ([]*models.GithubRepository, error) {
	f := logrus.Fields{
		"functionName": "GetRepositoryByCLAGroup",
		"claGroupID":   claGroupID,
		"enabled":      enabled,
	}
	builder := expression.NewBuilder()
	condition := expression.Key("repository_project_id").Equal(expression.Value(claGroupID))
	filter := expression.Name("enabled").Equal(expression.Value(enabled))
	builder = builder.WithKeyCondition(condition).WithFilter(filter)

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
		IndexName:                 aws.String(ProjectRepositoryIndex),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.WithFields(f).Warnf("unable to get project github repositories. error: %+v", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		log.WithFields(f).Warn("no repositories found matching the search criteria")
		return nil, ErrGithubRepositoryNotFound
	}

	var repositories []*RepositoryDBModel
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &repositories)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response, error: %+v", err)
		return nil, err
	}

	return convertModels(repositories), nil
}

// GetCLAGroupRepositoriesGroupByOrgs returns a list of GH organizations by CLA Group - enabled flag indicates that we search the enabled repositories list
func (repo repo) GetCLAGroupRepositoriesGroupByOrgs(projectID string, enabled bool) ([]*models.GithubRepositoriesGroupByOrgs, error) {
	out := make([]*models.GithubRepositoriesGroupByOrgs, 0)
	outMap := make(map[string]*models.GithubRepositoriesGroupByOrgs)
	ghrepos, err := repo.getProjectRepositories(projectID, enabled)
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

// List github repositories of project by external/salesforce project id
func (repo repo) ListProjectRepositories(externalProjectID string, projectSFID string, enabled bool) (*models.ListGithubRepositories, error) {
	f := logrus.Fields{
		"functionName":      "ListProjectRepositories",
		"externalProjectID": externalProjectID,
		"projectSFID":       projectSFID,
		"enabled":           enabled,
	}

	var indexName string
	out := &models.ListGithubRepositories{
		List: make([]*models.GithubRepository, 0),
	}
	var condition expression.KeyConditionBuilder

	if externalProjectID != "" {
		condition = expression.Key("repository_sfdc_id").Equal(expression.Value(externalProjectID))
		indexName = SFDCRepositoryIndex
	} else {
		condition = expression.Key("project_sfid").Equal(expression.Value(projectSFID))
		indexName = ProjectSFIDRepositoryOrgnizationNameIndex
	}

	// Add the enabled filter
	filter := expression.Name("enabled").Equal(expression.Value(enabled))

	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter).Build()
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
		IndexName:                 aws.String(indexName),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.WithFields(f).Warnf("unable to get project github repositories. error = %s", err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		return out, nil
	}
	var result []*RepositoryDBModel
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &result)
	if err != nil {
		return nil, err
	}
	for _, gr := range result {
		out.List = append(out.List, gr.toModel())
	}
	return out, nil
}

// getProjectRepositories returns an array of GH repositories for the specified project ID
func (repo repo) getProjectRepositories(projectID string, enabled bool) ([]*models.GithubRepository, error) {
	var out []*models.GithubRepository

	condition := expression.Key("repository_project_id").Equal(expression.Value(projectID))
	filter := expression.Name("enabled").Equal(expression.Value(enabled))
	builder := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter)
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
		TableName:                 aws.String(repo.repositoryTableName),
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
	var result []*RepositoryDBModel
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
		TableName:                 aws.String(repo.repositoryTableName),
	}

	results, err := repo.dynamoDBClient.Scan(scanInput)
	if err != nil {
		log.Warnf("unable to get github organizations repositories. error = %s", err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		return out, nil
	}
	var result []*RepositoryDBModel
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &result)
	if err != nil {
		return nil, err
	}
	for _, gr := range result {
		out = append(out, gr.toModel())
	}
	return out, nil
}

func (repo repo) getRepositoryByGithubID(externalID string, enabled bool) (*models.GithubRepository, error) {
	var condition expression.KeyConditionBuilder
	builder := expression.NewBuilder()
	condition = expression.Key("repository_external_id").Equal(expression.Value(externalID))
	filter := expression.Name("enabled").Equal(expression.Value(enabled))

	builder = builder.WithKeyCondition(condition).WithFilter(filter)
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
		TableName:                 aws.String(repo.repositoryTableName),
		IndexName:                 aws.String(ExternalRepositoryIndex),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.Warnf("unable to get project github repositories. error = %s", err.Error())
		return nil, err
	}
	var result *RepositoryDBModel
	if len(results.Items) == 0 {
		return nil, ErrGithubRepositoryNotFound
	}
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &result)
	if err != nil {
		return nil, err
	}
	return result.toModel(), nil
}

func (repo repo) enableGithubRepository(repositoryID string) error {
	return repo.updateGithubRepository(repositoryID, true)
}

func (repo repo) disableGithubRepository(repositoryID string) error {
	return repo.updateGithubRepository(repositoryID, false)
}

// deleteGithubRepository updates the existing repository record by setting the enabled flag to false
func (repo repo) updateGithubRepository(repositoryID string, enabled bool) error {
	f := logrus.Fields{
		"functionName": "deleteGithubRepository",
		"repositoryID": repositoryID,
		"enabled":      enabled,
	}

	// Load the existing model - need to fetch the old note value, if available
	existingModel, getErr := repo.GetRepository(repositoryID)
	if getErr != nil {
		return getErr
	}
	if existingModel == nil {
		return fmt.Errorf("unable to locate existing repository entry by ID: %s", repositoryID)
	}

	// Enabled string for the note...
	var enabledString = "disabled"
	if enabled {
		enabledString = "enabled"
	}

	// If we have an old note - grab it/save it
	var existingNote = ""
	if existingModel.Note != "" {
		existingNote = existingModel.Note + ". "
	}

	_, now := utils.CurrentTime()
	log.WithFields(f).Debug("updating repository record")
	_, err := repo.dynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"repository_id": {S: aws.String(repositoryID)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#enabled":      aws.String("enabled"),
			"#note":         aws.String("note"),
			"#dateModified": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":enabledValue": {
				BOOL: aws.Bool(enabled),
			},
			":noteValue": {
				S: aws.String(fmt.Sprintf("%s%s on %s", existingNote, enabledString, now)), // Add to existing note, if set
			},
			":dateModifiedValue": {
				S: aws.String(now),
			},
		},
		UpdateExpression: aws.String("SET #enabled = :enabledValue, #note = :noteValue, #dateModified = :dateModifiedValue"),
		TableName:        aws.String(repo.repositoryTableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.New("github repository entry does not exist or repository_sfdc_id does not match with specified project id")
			}
		}
		log.WithFields(f).Warnf("error disabling github repository, error: %+v", err)
		return err
	}

	return nil
}
