// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package company

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gofrs/uuid"
)

type Repository interface {
	GetPendingCompanyInviteRequests(companyID string) ([]CompanyInvite, error)
	GetCompany(CompanyID string) (Company, error)
	DeletePendingCompanyInviteRequest(InviteID string) error
	AddPendingCompanyInviteRequest(companyID string, userID string) error
	UpdateCompanyAccessList(companyID string, companyACL []string) error
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

type Company struct {
	CompanyID   string   `dynamodbav:"company_id"`
	CompanyName string   `dynamodbav:"company_name"`
	CompanyACL  []string `dynamodbav:"company_acl"`
}

type CompanyInvite struct {
	CompanyInviteID    string `dynamodbav:"company_invite_id"`
	RequestedCompanyID string `dynamodbav:"requested_company_id"`
	UserID             string `dynamodbav:"user_id"`
}

func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func (repo repository) GetCompany(CompanyID string) (Company, error) {
	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)
	companyTableData, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"company_id": {
				S: aws.String(CompanyID),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return Company{}, err
	}

	company := Company{}
	err = dynamodbattribute.UnmarshalMap(companyTableData.Item, &company)
	if err != nil {
		fmt.Println(err.Error())
		return Company{}, err
	}

	return company, nil
}

func (repo repository) GetPendingCompanyInviteRequests(companyID string) ([]CompanyInvite, error) {
	tableName := fmt.Sprintf("cla-%s-company-invites", repo.stage)
	input := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"requested_company_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(companyID),
					},
				},
			},
		},
		TableName: aws.String(tableName),
		IndexName: aws.String("requested-company-index"),
	}
	companyInviteAV, err := repo.dynamoDBClient.Query(input)
	if err != nil {
		fmt.Println("Unable to retrieve data from Company-Invites table", err)
		return nil, err
	}

	companyInvites := []CompanyInvite{}
	err = dynamodbattribute.UnmarshalListOfMaps(companyInviteAV.Items, &companyInvites)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return companyInvites, nil
}

func (repo repository) DeletePendingCompanyInviteRequest(inviteID string) error {
	tableName := fmt.Sprintf("cla-%s-company-invites", repo.stage)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"company_invite_id": {
				S: aws.String(inviteID),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := repo.dynamoDBClient.DeleteItem(input)
	if err != nil {
		fmt.Println("Unable to delete Company Invite Request", err)
		return err
	}

	return nil
}

func (repo repository) AddPendingCompanyInviteRequest(companyID string, userID string) error {
	companyInviteID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Unable to generate a UUID for a pending invite", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"company_invite_id": {
				S: aws.String(companyInviteID.String()),
			},
			"requested_company_id": {
				S: aws.String(companyID),
			},
			"user_id": {
				S: aws.String(userID),
			},
		},
		TableName: aws.String(fmt.Sprintf("cla-%s-company-invites", repo.stage)),
	}

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		fmt.Println("Unable to create a new pending invite", err)
	}

	return nil
}

func (repo repository) UpdateCompanyAccessList(companyID string, companyACL []string) error {
	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("company_acl"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				SS: aws.StringSlice(companyACL),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"company_id": {
				S: aws.String(companyID),
			},
		},
		UpdateExpression: aws.String("SET #S = :s"),
	}

	_, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		fmt.Println("Error updating Company Access List:", err)
		return err
	}

	return nil
}
