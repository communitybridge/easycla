// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/labstack/gommon/log"
)

// errors
var (
	ErrProjectDoesNotExist = errors.New("project does not exist")
)

// DBProjectModel data model
type DBProjectModel struct {
	DateCreated                      string   `dynamodbav:"date_created"`
	DateModified                     string   `dynamodbav:"date_modified"`
	ProjectExternalID                string   `dynamodbav:"project_external_id"`
	ProjectID                        string   `dynamodbav:"project_id"`
	ProjectName                      string   `dynamodbav:"project_name"`
	Version                          string   `dynamodbav:"version"`
	ProjectCclaEnabled               bool     `dynamodbav:"project_ccla_enabled"`
	ProjectCclaRequiresIclaSignature bool     `dynamodbav:"project_ccla_requires_icla_signature"`
	ProjectIclaEnabled               bool     `dynamodbav:"project_icla_enabled"`
	ProjectACL                       []string `dynamodbav:"project_acl"`
}

// Repository defines functions of Project repository
type Repository interface {
	GetMetrics() (*models.ProjectMetrics, error)
	GetProject(projectID string) (*models.Project, error)
	GetProjects() ([]models.Project, error)
	buildProjectModel(dbModel DBProjectModel) *models.Project
	buildProjectModels(results *dynamodb.ScanOutput) ([]models.Project, error)
}

// NewDynamoRepository creates instance of project repository
func NewDynamoRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		dynamoDBClient: dynamodb.New(awsSession),
		stage:          stage,
	}
}

type repo struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// GetMetrics returns the metrics for the projects
func (repo repo) GetMetrics() (*models.ProjectMetrics, error) {
	var out models.ProjectMetrics
	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count of projects, error: %v", err)
		return nil, err
	}

	out.TotalCount = *describeTableResult.Table.ItemCount
	return &out, nil
}

// GetProject returns the project model associated for the specified projectID
func (repo *repo) GetProject(projectID string) (*models.Project, error) {
	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(projectID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving project having ID : %s, error: %v", projectID, err)
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, ErrProjectDoesNotExist
	}
	var dbModel DBProjectModel
	err = dynamodbattribute.UnmarshalMap(result.Item, &dbModel)
	if err != nil {
		log.Warnf("error unmarshalling db project model, error: %+v", err)
		return nil, err
	}

	// Convert the database model to an API response model
	return repo.buildProjectModel(dbModel), nil
}

// GetProjects queries the database and returns a list of the projects
func (repo *repo) GetProjects() ([]models.Project, error) {
	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for project scan, error: %v", err)
	}

	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	var projects []models.Project
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		results, errQuery := repo.dynamoDBClient.Scan(scanInput)
		if errQuery != nil {
			log.Warnf("error retrieving projects, error: %v", errQuery)
			return nil, errQuery
		}

		// Convert the list of DB models to a list of response models
		projectList, modelErr := repo.buildProjectModels(results)
		if modelErr != nil {
			log.Warnf("error converting project DB model to response model, error: %v",
				modelErr)
			return nil, modelErr
		}

		// Add to the project response models to the list
		projects = append(projects, projectList...)

		if results.LastEvaluatedKey["project_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["project_id"].S
			scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"project_id": {
					S: aws.String(lastEvaluatedKey),
				},
			}
		} else {
			lastEvaluatedKey = ""
		}
	}

	return projects, nil
}

// buildProjectModels converts the database response model into an API response data model
func (repo *repo) buildProjectModels(results *dynamodb.ScanOutput) ([]models.Project, error) {
	var projects = make([]models.Project, *results.Count)

	// The DB project model
	var dbProjects []DBProjectModel

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &dbProjects)
	if err != nil {
		log.Warnf("error unmarshalling projects from database, error: %v", err)
		return nil, err
	}

	for _, dbProject := range dbProjects {
		projects = append(projects, models.Project{
			ProjectID:               dbProject.ProjectID,
			ProjectExternalID:       dbProject.ProjectExternalID,
			ProjectName:             dbProject.ProjectName,
			ProjectACL:              dbProject.ProjectACL,
			ProjectCCLAEnabled:      dbProject.ProjectCclaEnabled,
			ProjectICLAEnabled:      dbProject.ProjectIclaEnabled,
			ProjectCCLARequiresICLA: dbProject.ProjectCclaRequiresIclaSignature,
			DateCreated:             dbProject.DateCreated,
			DateModified:            dbProject.DateModified,
			Version:                 dbProject.Version,
		})
	}

	return projects, nil
}

// buildProjectModel maps the database model to the API response model
func (repo *repo) buildProjectModel(dbModel DBProjectModel) *models.Project {
	return &models.Project{
		ProjectID:          dbModel.ProjectID,
		ProjectName:        dbModel.ProjectName,
		ProjectACL:         dbModel.ProjectACL,
		ProjectCCLAEnabled: dbModel.ProjectCclaEnabled,
		ProjectICLAEnabled: dbModel.ProjectIclaEnabled,
		ProjectExternalID:  dbModel.ProjectExternalID,
		DateCreated:        dbModel.DateCreated,
		DateModified:       dbModel.DateModified,
		Version:            dbModel.Version,
	}
}

// buildProject is a helper function to build a common set of projection/columns for the query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("project_id"),
		expression.Name("project_external_id"),
		expression.Name("project_name"),
		expression.Name("project_acl"),
		expression.Name("project_ccla_enabled"),
		expression.Name("project_icla_enabled"),
		expression.Name("project_ccla_requires_icla_signature"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
	)
}
