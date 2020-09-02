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

	ExistsByName(gerritName string) ([]*models.Gerrit, error)
	GetGerritsByID(ID string, IDType string) (*models.GerritList, error)
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

func (repo *repo) ExistsByName(gerritName string) ([]*models.Gerrit, error) {
	resultList := make([]*models.Gerrit, 0)

	var condition expression.KeyConditionBuilder
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)

	// hashkey is gerrit-name
	indexName := "gerrit-name-index"

	builder := expression.NewBuilder().WithProjection(buildProjection())

	filter := expression.Name("gerrit_id").AttributeExists()
	condition = expression.Key("gerrit_name").Equal(expression.Value(gerritName))

	builder = builder.WithKeyCondition(condition).WithFilter(filter)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// Assemble the query input parameters
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		ScanIndexForward:          aws.Bool(false),
	}

	for {
		results, errQuery := repo.dynamoDBClient.Query(input)
		if errQuery != nil {
			log.Warnf("error retrieving Gerrit. error = %s", errQuery.Error())
			return nil, errQuery
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
			input.ExclusiveStartKey = results.LastEvaluatedKey

		} else {
			break
		}
	}
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].GerritName < resultList[j].GerritName
	})
	return resultList, nil
}

func (repo *repo) ExistsByID(gerritID string) ([]*models.Gerrit, error) {
	resultList := make([]*models.Gerrit, 0)

	var condition expression.KeyConditionBuilder
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)

	// hashkey is gerrit-name
	indexName := "gerrit-id-index"

	builder := expression.NewBuilder().WithProjection(buildProjection())

	filter := expression.Name("gerrit_id").AttributeExists()
	condition = expression.Key("group_id_icla").Equal(expression.Value(gerritID))
	//condition = expression.Key("group_id_ccla").Equal(expression.Value(gerritID))

	builder = builder.WithKeyCondition(condition).WithFilter(filter)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// Assemble the query input parameters
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		ScanIndexForward:          aws.Bool(false),
	}

	for {
		results, errQuery := repo.dynamoDBClient.Query(input)
		if errQuery != nil {
			log.Warnf("error retrieving Gerrit. error = %s", errQuery.Error())
			return nil, errQuery
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
			input.ExclusiveStartKey = results.LastEvaluatedKey

		} else {
			break
		}
	}
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].GerritName < resultList[j].GerritName
	})
	return resultList, nil
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

func (repo repo) GetGerritsByID(ID string, IDType string) (*models.GerritList, error) {
	var filter expression.ConditionBuilder
	resultList := make([]*models.Gerrit, 0)
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)

	if IDType == "ICLA" {
		filter = expression.Name("group_id_icla").Equal(expression.Value(ID))
	} else if IDType == "CCLA" {
		filter = expression.Name("group_id_ccla").Equal(expression.Value(ID))
	} else {
		return nil, errors.New("invalid IDType")
	}
	// filter = filter.And(expression.Name("project_sfid").Equal(expression.Value(*projectSFID)))

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

// buildProjection builds the query projection
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("gerrit_id"),
		expression.Name("gerrit_name"),
		expression.Name("gerrit_url"),
		expression.Name("group_id_ccla"),
		expression.Name("group_id_icla"),
		expression.Name("group_name_ccla"),
		expression.Name("group_name_icla"),
		expression.Name("project_id"),
		expression.Name("project_sfid"),
	)
}
