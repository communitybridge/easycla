// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gofrs/uuid"
)

// RepositoryService interface methods
type RepositoryService interface {
	GetPendingCompanyInviteRequests(companyID string) ([]Invite, error)
	GetCompany(companyID string) (Company, error)
	DeletePendingCompanyInviteRequest(InviteID string) error
	AddPendingCompanyInviteRequest(companyID string, userID string) error
	UpdateCompanyAccessList(companyID string, companyACL []string) error
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// Company data model
type Company struct {
	CompanyID   string   `dynamodbav:"company_id"`
	CompanyName string   `dynamodbav:"company_name"`
	CompanyACL  []string `dynamodbav:"company_acl"`
	Created     string   `dynamodbav:"date_created"`
	Updated     string   `dynamodbav:"date_modified"`
}

// Invite data model
type Invite struct {
	CompanyInviteID    string `dynamodbav:"company_invite_id"`
	RequestedCompanyID string `dynamodbav:"requested_company_id"`
	UserID             string `dynamodbav:"user_id"`
}

// NewRepository creates a new company repository instance
func NewRepository(awsSession *session.Session, stage string) RepositoryService {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// GetCompany returns a company based on the company ID
func (repo repository) GetCompany(companyID string) (Company, error) {
	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)
	companyTableData, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"company_id": {
				S: aws.String(companyID),
			},
		},
	})
	if err != nil {
		log.Warnf(err.Error())
		log.Warnf("error fetching company table data using company id: %s, error: %v", companyID, err)
		return Company{}, err
	}

	company := Company{}
	err = dynamodbattribute.UnmarshalMap(companyTableData.Item, &company)
	if err != nil {
		log.Warnf("error unmarshalling company table data, error: %v", err)
		return Company{}, err
	}

	return company, nil
}

// GetPendingCompanyInviteRequests returns a list of company invites when provided the company ID
func (repo repository) GetPendingCompanyInviteRequests(companyID string) ([]Invite, error) {
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
		log.Warnf("Unable to retrieve data from Company-Invites table, error: %v", err)
		return nil, err
	}

	companyInvites := []Invite{}
	err = dynamodbattribute.UnmarshalListOfMaps(companyInviteAV.Items, &companyInvites)
	if err != nil {
		log.Warnf("error unmarshalling company invite data, error: %v", err)
		return nil, err
	}

	return companyInvites, nil
}

// DeletePendingCompanyInviteRequest deletes the spending invite
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
		log.Warnf("Unable to delete Company Invite Request, error: %v", err)
		return err
	}

	return nil
}

// AddPendingCompanyInviteRequest adds a pending company invite when provided the company ID and user ID
func (repo repository) AddPendingCompanyInviteRequest(companyID string, userID string) error {
	companyInviteID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("Unable to generate a UUID for a pending invite, error: %v", err)
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
		log.Warnf("Unable to create a new pending invite, error: %v", err)
		return err
	}

	return nil
}

// UpdateCompanyAccessList updates the company ACL when provided the company ID and ACL list
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
		log.Warnf("Error updating Company Access List, error: %v", err)
		return err
	}

	return nil
}
