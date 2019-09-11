// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// SignatureRepository interface defines the functions for the github whitelist service
type SignatureRepository interface {
	GetGithubOrganizationsFromWhitelist(signatureID string) ([]models.GithubOrg, error)
	AddGithubOrganizationToWhitelist(signatureID, githubOrganizationID string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromWhitelist(signatureID, githubOrganizationID string) ([]models.GithubOrg, error)
	GetProjectSignatures(params signatures.GetProjectSignaturesParams, pageSize int64) (*models.ProjectSignatures, error)
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
	companyRepo    company.RepositoryService
	userRepo       user.RepositoryService
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string, companyRepo company.RepositoryService, userRepo user.RepositoryService) SignatureRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
		companyRepo:    companyRepo,
		userRepo:       userRepo,
	}
}

// GetGithubOrganizationsFromWhitelist returns a list of GH organizations stored in the whitelist
func (repo repository) GetGithubOrganizationsFromWhitelist(signatureID string) ([]models.GithubOrg, error) {
	// get item from dynamoDB table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for signatureID: %s, error: %v", signatureID, err)
		return nil, err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		return nil, nil
	}

	var orgs []models.GithubOrg
	for _, org := range itemFromMap.L {
		selected := true
		orgs = append(orgs, models.GithubOrg{
			ID:       org.S,
			Selected: &selected,
		})
	}

	return orgs, nil
}

// AddGithubOrganizationToWhitelist adds the specified GH organization to the whitelist
func (repo repository) AddGithubOrganizationToWhitelist(signatureID, GithubOrganizationID string) ([]models.GithubOrg, error) {
	// get item from dynamoDB table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	log.Debugf("querying database for github organization whitelist using signatureID: %s", signatureID)

	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for signatureID: %s and GH Org: %s, error: %v",
			signatureID, GithubOrganizationID, err)
		return nil, err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		log.Debugf("signatureID: %s is missing the 'github_org_whitelist' column - will add", signatureID)
		itemFromMap = &dynamodb.AttributeValue{}
	}

	// generate new List L without element to be deleted
	// if we find a org with the same id just return without updating the record
	var newList []*dynamodb.AttributeValue
	for _, element := range itemFromMap.L {
		newList = append(newList, element)
		if *element.S == GithubOrganizationID {
			log.Debugf("github organization for signature: %s already in the list - nothing to do, org id: %s",
				signatureID, GithubOrganizationID)
			return buildResponse(itemFromMap.L), nil
		}
	}

	// Add the organization to list
	log.Debugf("adding github organization for signature: %s to the list, org id: %s",
		signatureID, GithubOrganizationID)
	newList = append(newList, &dynamodb.AttributeValue{
		S: aws.String(GithubOrganizationID),
	})

	// return values flag - Returns all of the attributes of the item, as they appear after the UpdateItem operation.
	addReturnValues := "ALL_NEW" // nolint

	// Update dynamoDB table
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#L": aws.String("github_org_whitelist"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				L: newList,
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
		ReturnValues:     &addReturnValues,
	}

	log.Warnf("updating database record using signatureID: %s with values: %v", signatureID, newList)
	updatedValues, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating white list, error: %v", err)
		return nil, err
	}

	updatedItemFromMap, ok := updatedValues.Attributes["github_org_whitelist"]
	if !ok {
		msg := fmt.Sprintf("unable to fetch updated whitelist organization values for "+
			"organization id: %s for signature: %s - list is empty - returning empty list",
			GithubOrganizationID, signatureID)
		log.Debugf(msg)
		return []models.GithubOrg{}, nil
	}

	return buildResponse(updatedItemFromMap.L), nil
}

// DeleteGithubOrganizationFromWhitelist removes the specified GH organization from the whitelist
func (repo repository) DeleteGithubOrganizationFromWhitelist(signatureID, GithubOrganizationID string) ([]models.GithubOrg, error) {
	// get item from dynamoDB table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
	})

	if err != nil {
		log.Warnf("error retrieving GH organization whitelist for signatureID: %s and GH Org: %s, error: %v",
			signatureID, GithubOrganizationID, err)
		return nil, err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		log.Warnf("unable to remove whitelist organization: %s for signature: %s - list is empty",
			GithubOrganizationID, signatureID)
		return nil, errors.New("no github_org_whitelist column")
	}

	// generate new List L without element to be deleted
	var newList []*dynamodb.AttributeValue
	for _, element := range itemFromMap.L {
		if *element.S != GithubOrganizationID {
			newList = append(newList, element)
		}
	}

	if len(newList) == 0 {
		// Since we don't have any items in our list, we can't simply update dynamoDB with an empty list,
		// nooooo, that would be too easy. Instead:
		// We need to set the value to NULL to clear it out (otherwise we'll get a validation error like:)
		// ValidationException: ExpressionAttributeValues contains invalid value: Supplied AttributeValue
		// is empty, must contain exactly one of the supported datatypes for key)

		log.Debugf("clearing out github org whitelist for organization: %s for signature: %s - list is empty",
			GithubOrganizationID, signatureID)
		nullFlag := true

		// update dynamoDB table
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#L": aws.String("github_org_whitelist"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":l": {
					NULL: &nullFlag,
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"signature_id": {
					S: aws.String(signatureID),
				},
			},
			UpdateExpression: aws.String("SET #L = :l"),
		}

		_, err = repo.dynamoDBClient.UpdateItem(input)
		if err != nil {
			log.Warnf("error updating github org whitelist to NULL value, error: %v", err)
			return nil, err
		}

		// Return an empty list
		return []models.GithubOrg{}, nil
	}

	// return values flag - Returns all of the attributes of the item, as they appear after the UpdateItem operation.
	updatedReturnValues := "ALL_NEW" // nolint

	// update dynamoDB table
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#L": aws.String("github_org_whitelist"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				L: newList,
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
		ReturnValues:     &updatedReturnValues,
	}

	updatedValues, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating github org whitelist, error: %v", err)
		return nil, err
	}

	updatedItemFromMap, ok := updatedValues.Attributes["github_org_whitelist"]
	if !ok {
		msg := fmt.Sprintf("unable to fetch updated whitelist organization values for "+
			"organization id: %s for signature: %s - list is empty - returning empty list",
			GithubOrganizationID, signatureID)
		log.Debugf(msg)
		return []models.GithubOrg{}, nil
	}

	return buildResponse(updatedItemFromMap.L), nil

}

// GetProjectSignatures returns a list of signatures for the specified project
func (repo repository) GetProjectSignatures(params signatures.GetProjectSignaturesParams, pageSize int64) (*models.ProjectSignatures, error) {

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// This is the key we want to match
	condition := expression.Key("signature_project_id").Equal(expression.Value(params.ProjectID))

	// These are the columns we want returned
	projection := expression.NamesList(
		expression.Name("signature_id"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("signature_approved"),
		expression.Name("signature_document_major_version"),
		expression.Name("signature_document_minor_version"),
		expression.Name("signature_reference_id"),
		expression.Name("signature_reference_type"),       // user or company
		expression.Name("signature_signed"),               // T/F
		expression.Name("signature_type"),                 // ccla or cla
		expression.Name("signature_user_ccla_company_id"), // reference to the company
	)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(projection).Build()
	if err != nil {
		log.Warnf("error building expression for project signature ID scan, project: %s, error: %v",
			params.ProjectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		Limit:                     aws.Int64(pageSize),                   // The maximum number of items to evaluate (not necessarily the number of matching items)
		IndexName:                 aws.String("project-signature-index"), // Name of a secondary index to scan
	}

	// If we have the next key, set the exclusive start key value
	if params.NextKey != nil {
		log.Debugf("Received a nextKey, value: %s", *params.NextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: params.NextKey,
			},
			"signature_project_id": {
				S: &params.ProjectID,
			},
		}
	}

	//log.Debugf("Running scan using queryInput: %+v", queryInput)

	// Make the DynamoDB Query API call
	result, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.Warnf("error retrieving project signature ID for project: %s, error: %v",
			params.ProjectID, err)
		return nil, err
	}

	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count for project: %s, error: %v", params.ProjectID, err)
		return nil, err
	}

	// Meta-data for the response
	resultCount := *result.Count
	totalCount := *describeTableResult.Table.ItemCount
	var lastEvaluatedKey string
	if result.LastEvaluatedKey["signature_id"] != nil {
		lastEvaluatedKey = *result.LastEvaluatedKey["signature_id"].S
	}

	//log.Debugf("Result count : %d", *result.Count)
	//log.Debugf("Scanned count: %d", *result.ScannedCount)
	//log.Debugf("Total count  : %d", *describeTableResult.Table.ItemCount)
	//log.Debugf("Last key     : %s", lastEvaluatedKey)

	return repo.buildProjectSignatureModels(result, params.ProjectID, resultCount, totalCount, lastEvaluatedKey)
}

// buildProjectSignatureModels converts the response model into a response data model
func (repo repository) buildProjectSignatureModels(results *dynamodb.QueryOutput, projectID string, resultCount int64, totalCount int64, lastKey string) (*models.ProjectSignatures, error) {
	var signatures []models.Signature

	type ItemSignature struct {
		SignatureID                   string `json:"signature_id"`
		DateCreated                   string `json:"date_created"`
		DateModified                  string `json:"date_modified"`
		SignatureApproved             bool   `json:"signature_approved"`
		SignatureSigned               bool   `json:"signature_signed"`
		SignatureDocumentMajorVersion string `json:"signature_document_major_version"`
		SignatureDocumentMinorVersion string `json:"signature_document_minor_version"`
		SignatureReferenceID          string `json:"signature_reference_id"`
		SignatureReferenceType        string `json:"signature_reference_type"`
		SignatureType                 string `json:"signature_type"`
		SignatureUserCompanyID        string `json:"signature_user_ccla_company_id"`
	}

	// The DB signature model
	var dbSignature []ItemSignature

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &dbSignature)
	if err != nil {
		log.Warnf("error unmarshalling signatures from database for project: %s, error: %v",
			projectID, err)
		return nil, err
	}

	for _, dbSignature := range dbSignature {
		var companyName = ""
		var userName = ""
		var userLFID = ""
		var userGHID = ""

		var dbCompanyModel = company.Company{}
		var dbUserModel = user.User{}

		if dbSignature.SignatureReferenceType == "user" {
			dbUserModel, err = repo.userRepo.GetUser(dbSignature.SignatureReferenceID)
			if err != nil {
				log.Warnf("unable to lookup user using id: %s, error: %v", dbSignature.SignatureReferenceID, err)
			} else {
				userName = dbUserModel.UserName
				userLFID = dbUserModel.LFUsername
				userGHID = dbUserModel.UserGithubID
			}
		} else if dbSignature.SignatureReferenceType == "company" {
			dbCompanyModel, err = repo.companyRepo.GetCompany(dbSignature.SignatureReferenceID)
			if err != nil {
				log.Warnf("unable to lookup company using id: %s, error: %v", dbSignature.SignatureReferenceID, err)
			} else {
				companyName = dbCompanyModel.CompanyName
			}
		}

		signatures = append(signatures, models.Signature{
			SignatureID:            dbSignature.SignatureID,
			CompanyName:            companyName,
			SignatureCreated:       dbSignature.DateCreated,
			SignatureModified:      dbSignature.DateModified,
			SignatureType:          dbSignature.SignatureType,
			SignatureSigned:        dbSignature.SignatureSigned,
			SignatureApproved:      dbSignature.SignatureApproved,
			Version:                dbSignature.SignatureDocumentMajorVersion + "." + dbSignature.SignatureDocumentMinorVersion,
			SignatureReferenceType: dbSignature.SignatureReferenceType,
			UserName:               userName,
			UserLFID:               userLFID,
			UserGHID:               userGHID,
		})
	}

	return &models.ProjectSignatures{
		ProjectID:      projectID,
		ResultCount:    resultCount,
		TotalCount:     totalCount,
		LastKeyScanned: lastKey,
		Signatures:     signatures,
	}, nil
}

// buildResponse converts a database model to a GitHub organization response model
func buildResponse(items []*dynamodb.AttributeValue) []models.GithubOrg {
	// Convert to a response model
	var orgs []models.GithubOrg
	for _, org := range items {
		selected := true
		orgs = append(orgs, models.GithubOrg{
			ID:       org.S,
			Selected: &selected,
		})
	}

	return orgs
}
