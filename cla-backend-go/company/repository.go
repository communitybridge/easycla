// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"
	"strings"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/strfmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

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
	SearchCompanyByName(companyName string, nextKey string) (*models.Companies, error)
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

// SearchCompanyByName locates companies by the matching name and return any potential matches
func (repo repository) SearchCompanyByName(companyName string, nextKey string) (*models.Companies, error) {
	// Sorry, no results if empty company name
	if strings.TrimSpace(companyName) == "" {
		return &models.Companies{
			Companies:      []models.Company{},
			LastKeyScanned: "",
			ResultCount:    0,
			SearchTerms:    companyName,
			TotalCount:     0,
		}, nil

	}

	queryStartTime := time.Now()

	tableName := fmt.Sprintf("cla-%s-companies", repo.stage)

	// This is the company name we want to match
	filter := expression.Name("company_name").Contains(companyName)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(buildCompanyProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for company scan, companyName: %s, error: %v",
			companyName, err)
		return nil, err
	}

	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// If we have the next key, set the exclusive start key value
	if nextKey != "" {
		log.Debugf("Received a nextKey, value: %s", nextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"company_id": {
				S: aws.String(nextKey),
			},
		}
	}

	log.Debugf("Running company search scan using queryInput: %+v", scanInput)

	// Make the DynamoDB Query API call
	results, err := repo.dynamoDBClient.Scan(scanInput)
	if err != nil {
		log.Warnf("error retrieving companies for search term: %s, error: %v", companyName, err)
		return nil, err
	}

	log.Debugf("User signatures query took: %v resulting in %d results",
		utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total company record count for companyName: %s, error: %v", companyName, err)
		return nil, err
	}

	// Meta-data for the response
	resultCount := *results.Count
	totalCount := *describeTableResult.Table.ItemCount
	var lastEvaluatedKey string
	if results.LastEvaluatedKey["company_id"] != nil {
		//log.Debugf("LastEvaluatedKey: %+v", result.LastEvaluatedKey["signature_id"])
		lastEvaluatedKey = *results.LastEvaluatedKey["company_id"].S
	}

	response, err := repo.buildCompanyModels(results, resultCount, totalCount, lastEvaluatedKey)
	log.Debugf("Total company search took: %v resulting in %d results",
		utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

	return response, err
}

// buildCompanyModels converts the response model into a response data model
func (repo repository) buildCompanyModels(results *dynamodb.ScanOutput, resultCount int64, totalCount int64, lastKey string) (*models.Companies, error) {
	var companies []models.Company

	type ItemSignature struct {
		CompanyID   string   `json:"company_id"`
		CompanyName string   `json:"company_name"`
		CompanyACL  []string `json:"company_acl"`
		Created     string   `json:"date_created"`
		Modified    string   `json:"date_modified"`
	}

	// The DB company model
	var dbCompanies []ItemSignature

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &dbCompanies)
	if err != nil {
		log.Warnf("error unmarshalling companies from database, error: %v", err)
		return nil, err
	}

	dateTimeFormat := "2006-01-02T15:04:05.000000+0000"

	for _, dbCompany := range dbCompanies {
		createdDateTime, err := time.Parse(dateTimeFormat, dbCompany.Created)
		if err != nil {
			log.Warnf("Unable to parse company created date time: %s, error: %v - using current time",
				dbCompany.Created, err)
			createdDateTime = time.Now()
		}

		modifiedDateTime, err := time.Parse(dateTimeFormat, dbCompany.Modified)
		if err != nil {
			log.Warnf("Unable to parse company modified date time: %s, error: %v - using current time",
				dbCompany.Created, err)
			modifiedDateTime = time.Now()
		}

		companies = append(companies, models.Company{
			CompanyACL:  dbCompany.CompanyACL,
			CompanyID:   dbCompany.CompanyID,
			CompanyName: dbCompany.CompanyName,
			Created:     strfmt.DateTime(createdDateTime),
			Updated:     strfmt.DateTime(modifiedDateTime),
		})
	}

	return &models.Companies{
		ResultCount:    resultCount,
		TotalCount:     totalCount,
		LastKeyScanned: lastKey,
		Companies:      companies,
	}, nil
}

// buildCompanyProjection creates a ProjectionBuilds with the columns we are interested in
func buildCompanyProjection() expression.ProjectionBuilder {

	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("company_id"),
		expression.Name("company_name"),
		expression.Name("company_acl"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
	)
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
