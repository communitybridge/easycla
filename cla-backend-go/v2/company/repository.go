// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"

	"errors"

	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/gofrs/uuid"
)

var (
	//ErrDuplicateCompany error thrown in case company Name already exists
	ErrDuplicateCompany = errors.New("company already exists in CLA companies list")
)

// IRepository interface methods for company interface
type IRepository interface {
	CreateCompany(companyName string, companySFID string, userID string) error
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
	v1CompanyRepo  v1Company.IRepository
}

// NewRepository creates a new company repository instance
func NewRepository(awsSession *session.Session, stage string, v1CompRepo v1Company.IRepository) IRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
		v1CompanyRepo:  v1CompRepo,
	}
}

func (repo repository) CreateCompany(companyName string, companySFID string, userID string) error {
	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)

	companyID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate UUID for company creatiob, error: %v", err)
		return err
	}

	if err != nil {
		log.Warnf("Error building expression for Company input")
		return err
	}

	_, currentTimeString := utils.CurrentTime()

	item := &v1Company.Company{
		CompanyID:         companyID.String(),
		CompanyExternalID: companySFID,
		CompanyName:       companyName,
		CompanyManagerID:  userID,
		Updated:           currentTimeString,
		Created:           currentTimeString,
		Version:           "v1",
	}

	attributeValue, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new company item:")
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(tableName),
	}
	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.Warnf("Error adding company for :%s ", companyName)
		return err
	}

	log.Debugf("Successfully added easyCLA company for :%s", companyName)

	return nil
}
