package project

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/labstack/gommon/log"
)

var (
	ErrProjectDoesNotExist = errors.New("project does not exist")
)

// Project data model
type Project struct {
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

// DynamoRepository defines functions of Project repository
type DynamoRepository interface {
	GetProject(projectID string) (*Project, error)
}

// NewDynamoRepository creates instance of project repository
func NewDynamoRepository(awsSession *session.Session, stage string) DynamoRepository {
	return &repo{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

type repo struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

func (repo *repo) GetProject(projectID string) (*Project, error) {
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
	var project Project
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
