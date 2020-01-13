package project

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/labstack/gommon/log"
)

type Project struct {
	DateCreated                      string   `json:"date_created"`
	DateModified                     string   `json:"date_modified"`
	ProjectAcl                       []string `json:"project_acl"`
	ProjectCclaEnabled               bool     `json:"project_ccla_enabled"`
	ProjectCclaRequiresIclaSignature bool     `json:"project_ccla_requires_icla_signature"`
	ProjectExternalID                string   `json:"project_external_id"`
	ProjectIclaEnabled               bool     `json:"project_icla_enabled"`
	ProjectID                        string   `json:"project_id"`
	ProjectName                      string   `json:"project_name"`
	Version                          string   `json:"version"`
}

type DynamoRepository interface {
	GetProject(projectID string) (*Project, error)
}

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
	var project Project
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
