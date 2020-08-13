// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// errors
var (
	ErrGerritNotFound = errors.New("gerrit not found")
)

// Repository defines functions of Repositories
type Repository interface {
	GetClaGroupGerrits(projectID string, projectSFID *string) (*models.GerritList, error)
	DeleteGerrit(gerritID string) error
	GetGerrit(gerritID string) (*models.Gerrit, error)
	AddGerrit(input *models.Gerrit) (*models.Gerrit, error)
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

func (repo repo) GetClaGroupGerrits(projectID string, projectSFID *string) (*models.GerritList, error) {
	resultList := make([]*models.Gerrit, 0)
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)
	filter := expression.Name("project_id").Equal(expression.Value(projectID))
	if projectSFID != nil {
		filter = filter.And(expression.Name("project_sfid").Equal(expression.Value(*projectSFID)))
	}
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		log.Warnf("error building expression for gerrit instances scan, error: %v", err)
		return nil, err
	}
	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	for {
		results, err := repo.dynamoDBClient.Scan(scanInput)
		if err != nil {
			log.Warnf("error retrieving gerrit instances, error: %v", err)
			return nil, err
		}

		var gerrits []*Gerrit

		err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &gerrits)
		if err != nil {
			log.Warnf("error unmarshalling gerrit from database. error: %v", err)
			return nil, err
		}

		for _, g := range gerrits {
			resultList = append(resultList, g.toModel())
		}

		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].GerritName < resultList[j].GerritName
	})
	return &models.GerritList{List: resultList}, nil
}

func (repo *repo) DeleteGerrit(gerritID string) error {
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"gerrit_id": {
				S: aws.String(gerritID),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := repo.dynamoDBClient.DeleteItem(input)
	if err != nil {
		log.Warnf("error updating gerrit repository : %s during delete project process ", gerritID)
		return err
	}
	return nil
}

func (repo *repo) GetGerrit(gerritID string) (*models.Gerrit, error) {
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"gerrit_id": {
				S: aws.String(gerritID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := repo.dynamoDBClient.GetItem(input)
	if err != nil {
		log.Warnf("error getting gerrit repository : %s. error = %s", gerritID, err)
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, ErrGerritNotFound
	}
	var gerrit Gerrit
	err = dynamodbattribute.UnmarshalMap(result.Item, &gerrit)
	if err != nil {
		log.Warnf("unable to read data from gerrit repository : %s. error = %s", gerritID, err)
		return nil, err
	}
	return gerrit.toModel(), nil
}

func (repo *repo) AddGerrit(input *models.Gerrit) (*models.Gerrit, error) {
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)
	gerritID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	_, currentTime := utils.CurrentTime()
	gerrit := &Gerrit{
		DateCreated:   currentTime,
		DateModified:  currentTime,
		GerritID:      gerritID.String(),
		GerritName:    input.GerritName,
		GerritURL:     input.GerritURL.String(),
		GroupIDCcla:   input.GroupIDCcla,
		GroupIDIcla:   input.GroupIDIcla,
		GroupNameCcla: input.GroupNameCcla,
		GroupNameIcla: input.GroupNameIcla,
		ProjectID:     input.ProjectID,
		ProjectSFID:   input.ProjectSFID,
		Version:       input.Version,
	}
	av, err := dynamodbattribute.MarshalMap(gerrit)
	if err != nil {
		return nil, err
	}
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Error("cannot put gerrit in dynamodb", err)
		return nil, err
	}
	return repo.GetGerrit(gerritID.String())
}
