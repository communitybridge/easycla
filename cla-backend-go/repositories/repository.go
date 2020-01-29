package repositories

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/labstack/gommon/log"
)

// Repository defines functions of Repositories
type Repository interface {
	GetMetrics() (*models.RepositoryMetrics, error)
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
