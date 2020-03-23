// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

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

// Repository defines functions of Repositories
type Repository interface {
	GetProjectGerrits(projectID string) ([]*models.Gerrit, error)
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

func (repo repo) GetProjectGerrits(projectID string) ([]*models.Gerrit, error) {
	out := make([]*models.Gerrit, 0)
	tableName := fmt.Sprintf("cla-%s-gerrit-instances", repo.stage)
	filter := expression.Name("project_id").Equal(expression.Value(projectID))
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
			out = append(out, g.toModel())
		}

		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return out, nil
}
