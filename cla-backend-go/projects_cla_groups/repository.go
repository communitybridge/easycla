package projects_cla_groups

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/labstack/gommon/log"
)

// ProjectClaGroup is database model for projects_cla_group table
type ProjectClaGroup struct {
	ProjectSFID string `json:"project_sfid"`
	ClaGroupID  string `json:"cla_group_id"`
}

// Repository provides interface for interacting with project_cla_groups table
type Repository interface {
	GetClaGroupsIdsForProject(projectSFID string) ([]*ProjectClaGroup, error)
	AssociateClaGroupWithProject(claGroupID string, projectSFID string) error
}

type repo struct {
	tableName      string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository provides implementation of projects_cla_group repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		tableName:      fmt.Sprintf("cla-%s-projects-cla-groups", stage),
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func (repo *repo) GetClaGroupsIdsForProject(projectSFID string) ([]*ProjectClaGroup, error) {
	keyCondition := expression.Key("project_sfid").Equal(expression.Value(projectSFID))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		log.Warnf("error building expression for project cla groups, error: %v", err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.tableName),
	}

	var projectClaGroups []*ProjectClaGroup
	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project cla-groups, error: %v", errQuery)
			return nil, errQuery
		}

		var projectClaGroupsTmp []*ProjectClaGroup

		err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &projectClaGroupsTmp)
		if err != nil {
			log.Warnf("error unmarshalling project cla-groups from database. error: %v", err)
			return nil, err
		}
		projectClaGroups = append(projectClaGroups, projectClaGroupsTmp...)

		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return projectClaGroups, nil
}

// AssociateClaGroupWithProject creates entry in db to track cla_group association with project/foundation
func (repo *repo) AssociateClaGroupWithProject(claGroupID string, projectSFID string) error {
	input := &ProjectClaGroup{
		ProjectSFID: projectSFID,
		ClaGroupID:  claGroupID,
	}
	av, err := dynamodbattribute.MarshalMap(input)
	if err != nil {
		return err
	}
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.tableName),
	})
	if err != nil {
		log.Errorf("cannot put association entry of cla_group_id: %s, project_sfid: %s in dynamodb. error = %s",
			claGroupID, projectSFID, err)
		return err
	}
	return nil
}
