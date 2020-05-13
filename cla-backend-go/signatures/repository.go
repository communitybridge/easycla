// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"

	"github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// constants
const (
	LoadACLDetails     = true
	DontLoadACLDetails = false
)

// SignatureRepository interface defines the functions for the github whitelist service
type SignatureRepository interface {
	GetGithubOrganizationsFromWhitelist(signatureID string) ([]models.GithubOrg, error)
	AddGithubOrganizationToWhitelist(signatureID, githubOrganizationID string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromWhitelist(signatureID, githubOrganizationID string) ([]models.GithubOrg, error)
	InvalidateProjectRecord(signatureID string, projectName string) error

	GetSignature(signatureID string) (*models.Signature, error)
	GetSignatureACL(signatureID string) ([]string, error)
	GetProjectSignatures(params signatures.GetProjectSignaturesParams, pageSize int64) (*models.Signatures, error)
	GetProjectCompanySignatures(companyID, projectID string, nextKey *string, pageSize int64) (*models.Signatures, error)
	GetProjectCompanyEmployeeSignatures(params signatures.GetProjectCompanyEmployeeSignaturesParams, pageSize int64) (*models.Signatures, error)
	GetCompanySignatures(params signatures.GetCompanySignaturesParams, pageSize int64, loadACL bool) (*models.Signatures, error)
	GetUserSignatures(params signatures.GetUserSignaturesParams, pageSize int64) (*models.Signatures, error)
	ProjectSignatures(projectID string) (*models.Signatures, error)

	AddCLAManager(signatureID, claManagerID string) (*models.Signature, error)
	RemoveCLAManager(signatureID, claManagerID string) (*models.Signature, error)
}

// repository data model
type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
	companyRepo    company.IRepository
	usersRepo      users.Service
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string, companyRepo company.IRepository, usersRepo users.Service) SignatureRepository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
		companyRepo:    companyRepo,
		usersRepo:      usersRepo,
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

	// Sort the array based on the ID
	sort.Slice(orgs, func(i, j int) bool {
		return *orgs[i].ID < *orgs[j].ID
	})

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

// GetSignature returns the signature for the specified signature id
func (repo repository) GetSignature(signatureID string) (*models.Signature, error) {
	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// This is the key we want to match
	condition := expression.Key("signature_id").Equal(expression.Value(signatureID))

	// Use the builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for signature ID query, signatureID: %s, error: %v",
			signatureID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	results, queryErr := repo.dynamoDBClient.Query(queryInput)
	if queryErr != nil {
		log.Warnf("error retrieving signature ID: %s, error: %v", signatureID, queryErr)
		return nil, queryErr
	}

	// No match, didn't find it
	if *results.Count == 0 {
		return nil, nil
	}

	// Convert the list of DB models to a list of response models - should have zero or 1 given that we query by ID
	signatureList, modelErr := repo.buildProjectSignatureModels(results, "", LoadACLDetails)
	if modelErr != nil {
		log.Warnf("error converting DB model to response model for signature: %s, error: %v",
			signatureID, modelErr)
		return nil, modelErr
	}

	if len(signatureList) == 0 {
		return nil, nil
	}

	return signatureList[0], nil
}

// GetSignatureACL returns the signature ACL for the specified signature id
func (repo repository) GetSignatureACL(signatureID string) ([]string, error) {
	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().
		WithProjection(buildSignatureACLProjection()).
		Build()
	if err != nil {
		log.Warnf("error building expression for signature ID query, signatureID: %s, error: %v",
			signatureID, err)
		return nil, err
	}

	// Assemble the query input parameters
	itemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {S: aws.String(signatureID)},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, queryErr := repo.dynamoDBClient.GetItem(itemInput)
	if queryErr != nil {
		log.Warnf("error retrieving signature ID: %s, error: %v", signatureID, queryErr)
		return nil, queryErr
	}

	// No match, didn't find it
	if result.Item == nil {
		return nil, nil
	}

	var dbModel DBManagersModel
	// Unmarshall the DB response
	unmarshallErr := dynamodbattribute.UnmarshalMap(result.Item, &dbModel)
	if unmarshallErr != nil {
		log.Warnf("error converting DB model signature query using siganture ID: %s, error: %v",
			signatureID, unmarshallErr)
		return nil, unmarshallErr
	}

	return dbModel.SignatureACL, nil
}

func addConditionToFilter(filter expression.ConditionBuilder, cond expression.ConditionBuilder, filterAdded *bool) expression.ConditionBuilder {
	if !(*filterAdded) {
		*filterAdded = true
		filter = cond
	} else {
		filter = filter.And(cond)
	}
	return filter
}

const projIndexName = "project-signature-index"

// GetProjectSignatures returns a list of signatures for the specified project
func (repo repository) GetProjectSignatures(params signatures.GetProjectSignaturesParams, pageSize int64) (*models.Signatures, error) {

	queryStartTime := time.Now()

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	indexName := projIndexName

	// This is the key we want to match
	condition := expression.Key("signature_project_id").Equal(expression.Value(params.ProjectID))

	builder := expression.NewBuilder().WithProjection(buildProjection())
	var filter expression.ConditionBuilder
	var filterAdded bool

	if params.SearchField != nil {
		searchFieldExpression := expression.Name("signature_reference_type").Equal(expression.Value(params.SearchField))
		filter = addConditionToFilter(filter, searchFieldExpression, &filterAdded)
	}

	if params.SignatureType != nil {
		if params.SearchTerm != nil && (params.FullMatch != nil && !*params.FullMatch) {
			indexName = "signature-project-id-type-index"
			condition = condition.And(expression.Key("signature_type").Equal(expression.Value(strings.ToLower(*params.SignatureType))))
		} else {
			signatureTypeExpression := expression.Name("signature_type").Equal(expression.Value(params.SignatureType))
			filter = addConditionToFilter(filter, signatureTypeExpression, &filterAdded)
		}
		if *params.SignatureType == "ccla" {
			signatureReferenceIDExpression := expression.Name("signature_reference_id").AttributeExists()
			signatureUserCclaCompanyIDExpression := expression.Name("signature_user_ccla_company_id").AttributeNotExists()
			filter = addConditionToFilter(filter, signatureReferenceIDExpression, &filterAdded)
			filter = addConditionToFilter(filter, signatureUserCclaCompanyIDExpression, &filterAdded)
		}
	}

	if params.SearchTerm != nil {
		if *params.FullMatch {
			indexName = "reference-signature-search-index"
			condition = condition.And(expression.Key("signature_reference_name_lower").Equal(expression.Value(strings.ToLower(*params.SearchTerm))))
		} else {
			searchTermExpression := expression.Name("signature_reference_name_lower").Contains(strings.ToLower(*params.SearchTerm))
			filter = addConditionToFilter(filter, searchTermExpression, &filterAdded)
		}
	}

	// Filter condition to cater for approved and signed signatures
	signatureApprovedExpression := expression.Name("signature_approved").Equal(expression.Value(true))
	filter = addConditionToFilter(filter, signatureApprovedExpression, &filterAdded)

	signatureSignedExpression := expression.Name("signature_signed").Equal(expression.Value(true))
	filter = addConditionToFilter(filter, signatureSignedExpression, &filterAdded)

	if filterAdded {
		builder = builder.WithFilter(filter)
	}
	builder = builder.WithKeyCondition(condition)

	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		log.Warnf("error building expression for project signature query, projectID: %s, error: %v",
			params.ProjectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		Limit:                     aws.Int64(pageSize),   // The maximum number of items to evaluate (not necessarily the number of matching items)
		IndexName:                 aws.String(indexName), // Name of a secondary index to scan
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
		if params.FullMatch != nil && *params.FullMatch && params.SearchTerm != nil {
			queryInput.ExclusiveStartKey["signature_reference_name_lower"] = &dynamodb.AttributeValue{
				S: params.SearchTerm,
			}
		}
	}

	signatures := make([]*models.Signature, 0)
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		log.Debugf("Running signature project query using queryInput: %+v", queryInput)
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project signature ID for project: %s, error: %v",
				params.ProjectID, errQuery)
			return nil, errQuery
		}

		log.Debugf("Signature project query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

		// Convert the list of DB models to a list of response models
		signatureList, modelErr := repo.buildProjectSignatureModels(results, params.ProjectID, LoadACLDetails)
		if modelErr != nil {
			log.Warnf("error converting DB model to response model for signatures with project %s, error: %v",
				params.ProjectID, modelErr)
			return nil, modelErr
		}

		// Add to the signatures response model to the list
		signatures = append(signatures, signatureList...)

		//log.Debugf("LastEvaluatedKey: %+v", results.LastEvaluatedKey)
		if results.LastEvaluatedKey["signature_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["signature_id"].S
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(signatures)) >= pageSize {
			break
		}
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
	totalCount := *describeTableResult.Table.ItemCount
	if int64(len(signatures)) > pageSize {
		signatures = signatures[0:pageSize]
		lastEvaluatedKey = signatures[pageSize-1].SignatureID
	}

	return &models.Signatures{
		ProjectID:      params.ProjectID,
		ResultCount:    int64(len(signatures)),
		TotalCount:     totalCount,
		LastKeyScanned: lastEvaluatedKey,
		Signatures:     signatures,
	}, nil
}

// GetProjectCompanySignatures returns a list of signatures for the specified project and specified company
func (repo repository) GetProjectCompanySignatures(companyID, projectID string, nextKey *string, pageSize int64) (*models.Signatures, error) {

	queryStartTime := time.Now()

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// These are the keys we want to match
	condition := expression.Key("signature_project_id").Equal(expression.Value(projectID))
	filter := expression.Name("signature_reference_id").Equal(expression.Value(companyID))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter).WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for project signature ID query, project: %s, error: %v",
			projectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("project-signature-index"), // Name of a secondary index to scan
		//Limit:                     aws.Int64(pageSize),                   // The maximum number of items to evaluate (not necessarily the number of matching items)
	}

	// If we have the next key, set the exclusive start key value
	if nextKey != nil {
		log.Debugf("Received a nextKey, value: %s", *nextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: nextKey,
			},
			"signature_project_id": {
				S: &projectID,
			},
		}
	}

	var signatures []*models.Signature
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		//log.Debugf("Running signature project company query using queryInput: %+v", queryInput)
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project signature ID for project: %s with company: %s, error: %v",
				projectID, companyID, errQuery)
			return nil, errQuery
		}

		log.Debugf("Signature project company query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

		// Convert the list of DB models to a list of response models
		signatureList, modelErr := repo.buildProjectSignatureModels(results, projectID, LoadACLDetails)
		if modelErr != nil {
			log.Warnf("error converting DB model to response model for signatures with project %s with company: %s, error: %v",
				projectID, companyID, modelErr)
			return nil, modelErr
		}

		// Add to the signatures response model to the list
		signatures = append(signatures, signatureList...)

		// log.Debugf("LastEvaluatedKey: %+v", results.LastEvaluatedKey["signature_id"])
		if results.LastEvaluatedKey["signature_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["signature_id"].S
			queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"signature_id": {
					S: aws.String(lastEvaluatedKey),
				},
				"signature_project_id": {
					S: &projectID,
				},
			}
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(signatures)) >= pageSize {
			break
		}
	}

	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count for project: %s, error: %v", projectID, err)
		return nil, err
	}

	// Meta-data for the response
	totalCount := *describeTableResult.Table.ItemCount

	return &models.Signatures{
		ProjectID:      projectID,
		ResultCount:    int64(len(signatures)),
		TotalCount:     totalCount,
		LastKeyScanned: lastEvaluatedKey,
		Signatures:     signatures,
	}, nil
}

// Get project signatures with no pagination
func (repo repository) ProjectSignatures(projectID string) (*models.Signatures, error) {

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	indexName := projIndexName

	// This is the key we want to match
	condition := expression.Key("signature_project_id").Equal(expression.Value(projectID))

	builder := expression.NewBuilder().WithProjection(buildProjection())
	var filter expression.ConditionBuilder
	var filterAdded bool

	// Filter condition to cater for approved and signed signatures
	signatureApprovedExpression := expression.Name("signature_approved").Equal(expression.Value(true))
	filter = addConditionToFilter(filter, signatureApprovedExpression, &filterAdded)

	signatureSignedExpression := expression.Name("signature_signed").Equal(expression.Value(true))
	filter = addConditionToFilter(filter, signatureSignedExpression, &filterAdded)

	if filterAdded {
		builder = builder.WithFilter(filter)
	}
	builder = builder.WithKeyCondition(condition)

	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		log.Warnf("error building expression for project signature query, projectID: %s, error: %v",
			projectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName), // Name of a secondary index to scan
	}

	results, errQuery := repo.dynamoDBClient.Query(queryInput)

	if errQuery != nil {
		log.Warnf("error retrieving project signature ID for project: %s, error: %v",
			projectID, errQuery)
		return nil, errQuery
	}

	// Convert the list of DB models to a list of response models
	sigs, modelErr := repo.buildProjectSignatureModels(results, projectID, LoadACLDetails)
	if modelErr != nil {
		log.Warnf("error converting DB model to response model for signatures with project %s, error: %v",
			projectID, modelErr)
		return nil, modelErr
	}

	return &models.Signatures{
		ProjectID:  projectID,
		Signatures: sigs,
	}, nil
}

func (repo repository) InvalidateProjectRecord(signatureID string, projectName string) error {
	// Update project signatures for signature_approved and notes attributes
	signatureTableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	expressionAttributeNames := map[string]*string{}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{}
	updateExpression := "SET "

	expressionAttributeNames["#A"] = aws.String("signature_approved")
	expressionAttributeValues[":a"] = &dynamodb.AttributeValue{BOOL: aws.Bool(false)}
	updateExpression = updateExpression + " #A = :a,"

	expressionAttributeNames["#S"] = aws.String("note")
	note := fmt.Sprintf("Signature invalidated (approved set to false) due to CLA Group/Project: %s deletion", projectName)
	expressionAttributeValues[":s"] = &dynamodb.AttributeValue{S: aws.String(note)}
	updateExpression = updateExpression + " #S = :s"

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		UpdateExpression:          &updateExpression,
		TableName:                 aws.String(signatureTableName),
	}

	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.Warnf("error updating signature_approved for signature_id : %s error : %v ", signatureID, updateErr)
		return updateErr
	}

	return nil
}

// GetProjectCompanyEmployeeSignatures returns a list of employee signatures for the specified project and specified company
func (repo repository) GetProjectCompanyEmployeeSignatures(params signatures.GetProjectCompanyEmployeeSignaturesParams, pageSize int64) (*models.Signatures, error) {

	queryStartTime := time.Now()

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// This is the keys we want to match
	condition := expression.Key("signature_project_id").Equal(expression.Value(params.ProjectID))
	filter := expression.Name("signature_user_ccla_company_id").Equal(expression.Value(params.CompanyID))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter).WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for project signature ID query, project: %s, error: %v",
			params.ProjectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
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

	signatures := make([]*models.Signature, 0)
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		//log.Debugf("Running signature project company query using queryInput: %+v", queryInput)
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project company employee signature ID for project: %s with company: %s, error: %v",
				params.ProjectID, params.CompanyID, errQuery)
			return nil, errQuery
		}

		log.Debugf("Signature project company employee query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

		// Convert the list of DB models to a list of response models
		signatureList, modelErr := repo.buildProjectSignatureModels(results, params.ProjectID, LoadACLDetails)
		if modelErr != nil {
			log.Warnf("error converting DB model to response model for employee signatures with project %s with company: %s, error: %v",
				params.ProjectID, params.CompanyID, modelErr)
			return nil, modelErr
		}

		// Add to the signatures response model to the list
		signatures = append(signatures, signatureList...)

		// log.Debugf("LastEvaluatedKey: %+v", results.LastEvaluatedKey["signature_id"])
		if results.LastEvaluatedKey["signature_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["signature_id"].S
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(signatures)) >= pageSize {
			break
		}
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
	totalCount := *describeTableResult.Table.ItemCount
	if int64(len(signatures)) > pageSize {
		signatures = signatures[0:pageSize]
		lastEvaluatedKey = signatures[pageSize-1].SignatureID
	}

	return &models.Signatures{
		ProjectID:      params.ProjectID,
		ResultCount:    int64(len(signatures)),
		TotalCount:     totalCount,
		LastKeyScanned: lastEvaluatedKey,
		Signatures:     signatures,
	}, nil
}

// GetCompanySignatures returns a list of company signatures for the specified company
func (repo repository) GetCompanySignatures(params signatures.GetCompanySignaturesParams, pageSize int64, loadACL bool) (*models.Signatures, error) {

	queryStartTime := time.Now()

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// This is the keys we want to match
	condition := expression.Key("signature_reference_id").Equal(expression.Value(params.CompanyID))

	// Check for approved signatures
	filter := expression.Name("signature_approved").Equal(expression.Value(aws.Bool(true))).
		And(expression.Name("signature_signed").Equal(expression.Value(aws.Bool(true))))

	if params.SignatureType != nil {
		filter = filter.And(expression.Name("signature_type").Equal(expression.Value(*params.SignatureType)))
	}

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter).WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for company signature query, companyID: %s, error: %v",
			params.CompanyID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("reference-signature-index"), // Name of a secondary index to scan
		//Limit:                     aws.Int64(pageSize),                   // The maximum number of items to evaluate (not necessarily the number of matching items)
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
			"signature_reference_id": {
				S: &params.CompanyID,
			},
		}
	}

	signatures := make([]*models.Signature, 0)
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		//log.Debugf("Running signature project company query using queryInput: %+v", queryInput)
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving company signature ID for company: %s with company: %s, error: %v",
				params.CompanyID, params.CompanyID, errQuery)
			return nil, errQuery
		}

		log.Debugf("Signature company query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

		// Convert the list of DB models to a list of response models
		signatureList, modelErr := repo.buildProjectSignatureModels(results, "", loadACL)
		if modelErr != nil {
			log.Warnf("error converting DB model to response model for signatures with company: %s, error: %v",
				params.CompanyID, modelErr)
			return nil, modelErr
		}

		// Add to the signatures response model to the list
		signatures = append(signatures, signatureList...)

		// log.Debugf("LastEvaluatedKey: %+v", results.LastEvaluatedKey["signature_id"])
		if results.LastEvaluatedKey["signature_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["signature_id"].S
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(signatures)) >= pageSize {
			break
		}
	}

	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count for company: %s/%s, error: %v",
			params.CompanyID, *params.CompanyName, err)
		return nil, err
	}
	if int64(len(signatures)) > pageSize {
		signatures = signatures[0:pageSize]
		lastEvaluatedKey = signatures[pageSize-1].SignatureID
	}

	// Meta-data for the response
	totalCount := *describeTableResult.Table.ItemCount

	return &models.Signatures{
		ProjectID:      "",
		ResultCount:    int64(len(signatures)),
		TotalCount:     totalCount,
		LastKeyScanned: lastEvaluatedKey,
		Signatures:     signatures,
	}, nil
}

// GetUserSignatures returns a list of user signatures for the specified user
func (repo repository) GetUserSignatures(params signatures.GetUserSignaturesParams, pageSize int64) (*models.Signatures, error) {

	queryStartTime := time.Now()

	// The table we're interested in
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)

	// This is the keys we want to match
	condition := expression.Key("signature_reference_id").Equal(expression.Value(params.UserID))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.Warnf("error building expression for user signature query, userID: %s, error: %v",
			params.UserID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("reference-signature-index"), // Name of a secondary index to scan
		Limit:                     aws.Int64(pageSize),                     // The maximum number of items to evaluate (not necessarily the number of matching items)
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
			"signature_reference_id": {
				S: &params.UserID,
			},
		}
	}

	signatures := make([]*models.Signature, 0)
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		// Make the DynamoDB Query API call
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving user signatures for user: %s/%s, error: %v",
				params.UserID, *params.UserName, errQuery)
			return nil, errQuery
		}

		log.Debugf("Signature user query took: %v resulting in %d results",
			utils.FmtDuration(time.Since(queryStartTime)), len(results.Items))

		// Convert the list of DB models to a list of response models
		signatureList, modelErr := repo.buildProjectSignatureModels(results, "", LoadACLDetails)
		if modelErr != nil {
			log.Warnf("error converting DB model to response model for signatures for user %s/%s, error: %v",
				params.UserID, *params.UserName, modelErr)
			return nil, modelErr
		}

		// Add to the signatures response model to the list
		signatures = append(signatures, signatureList...)

		// log.Debugf("LastEvaluatedKey: %+v", results.LastEvaluatedKey["signature_id"])
		if results.LastEvaluatedKey["signature_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["signature_id"].S
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(signatures)) >= pageSize {
			break
		}
	}

	// How many total records do we have - may not be up-to-date as this value is updated only periodically
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}
	describeTableResult, err := repo.dynamoDBClient.DescribeTable(describeTableInput)
	if err != nil {
		log.Warnf("error retrieving total record count for user: %s/%s, error: %v",
			params.UserID, *params.UserName, err)
		return nil, err
	}

	// Meta-data for the response
	totalCount := *describeTableResult.Table.ItemCount

	return &models.Signatures{
		ProjectID:      "",
		ResultCount:    int64(len(signatures)),
		TotalCount:     totalCount,
		LastKeyScanned: lastEvaluatedKey,
		Signatures:     signatures,
	}, nil
}

func (repo repository) AddCLAManager(signatureID, claManagerID string) (*models.Signature, error) {
	aclEntries, err := repo.GetSignatureACL(signatureID)
	if err != nil {
		log.Warnf("unable to fetch signature by ID: %s, error: %+v", signatureID, err)
		return nil, err
	}

	if aclEntries == nil {
		log.Warnf("unable to fetch signature by ID: %s - record not found", signatureID)
		return nil, nil
	}

	for _, manager := range aclEntries {
		if claManagerID == manager {
			return nil, errors.New("manager already in signature ACL")
		}
	}

	aclEntries = append(aclEntries, claManagerID)

	_, now := utils.CurrentTime()

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("signature_acl"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				SS: aws.StringSlice(aclEntries),
			},
			":m": {
				S: aws.String(now),
			},
		},
		UpdateExpression: aws.String("SET #A = :a, #M = :m"),
		TableName:        aws.String(fmt.Sprintf("cla-%s-signatures", repo.stage)),
	}

	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.Warnf("add CLA manager - unable to update request with new ACL entry of '%s' for signature ID: %s, error: %v",
			claManagerID, signatureID, updateErr)
		return nil, updateErr
	}

	// Load the updated document and return it
	sigModel, err := repo.GetSignature(signatureID)
	if err != nil {
		log.Warnf("unable to fetch signature by ID: %s - record not found", signatureID)
		return nil, err
	}

	return sigModel, nil
}

func (repo repository) RemoveCLAManager(signatureID, claManagerID string) (*models.Signature, error) {
	aclEntries, err := repo.GetSignatureACL(signatureID)
	if err != nil {
		log.Warnf("unable to fetch signature by ID: %s, error: %+v", signatureID, err)
		return nil, err
	}

	if aclEntries == nil {
		log.Warnf("unable to fetch signature by ID: %s - record not found", signatureID)
		return nil, nil
	}

	// A bit of logic to determine if the manager is listed and to build the new list without the specified manager
	found := false
	var updateEntries []string
	for _, manager := range aclEntries {
		if claManagerID == manager {
			found = true
		} else {
			updateEntries = append(updateEntries, manager)
		}
	}

	if !found {
		return nil, fmt.Errorf("manager ID: %s not found in signature ACL", claManagerID)
	}

	_, now := utils.CurrentTime()

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(signatureID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("signature_acl"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				SS: aws.StringSlice(updateEntries),
			},
			":m": {
				S: aws.String(now),
			},
		},
		UpdateExpression: aws.String("SET #A = :a, #M = :m"),
		TableName:        aws.String(fmt.Sprintf("cla-%s-signatures", repo.stage)),
	}

	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.Warnf("remove CLA manager - unable to remove ACL entry of '%s' for signature ID: %s, error: %v",
			claManagerID, signatureID, updateErr)
		return nil, updateErr
	}

	// Load the updated document and return it
	sigModel, err := repo.GetSignature(signatureID)
	if err != nil {
		log.Warnf("unable to fetch signature by ID: %s - record not found", signatureID)
		return nil, err
	}

	return sigModel, nil
}

// buildProjectSignatureModels converts the response model into a response data model
func (repo repository) buildProjectSignatureModels(results *dynamodb.QueryOutput, projectID string, loadACLDetails bool) ([]*models.Signature, error) {
	var signatures []*models.Signature

	// The DB signature model
	var dbSignatures []ItemSignature

	err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &dbSignatures)
	if err != nil {
		log.Warnf("error unmarshalling signatures from database for project: %s, error: %v",
			projectID, err)
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(dbSignatures))
	for _, dbSignature := range dbSignatures {
		sig := &models.Signature{
			SignatureID:                 dbSignature.SignatureID,
			SignatureCreated:            dbSignature.DateCreated,
			SignatureModified:           dbSignature.DateModified,
			SignatureType:               dbSignature.SignatureType,
			SignatureReferenceID:        dbSignature.SignatureReferenceID,
			SignatureReferenceName:      dbSignature.SignatureReferenceName,
			SignatureReferenceNameLower: dbSignature.SignatureReferenceNameLower,
			SignatureSigned:             dbSignature.SignatureSigned,
			SignatureApproved:           dbSignature.SignatureApproved,
			SignatureMajorVersion:       dbSignature.SignatureDocumentMajorVersion,
			SignatureMinorVersion:       dbSignature.SignatureDocumentMinorVersion,
			Version:                     dbSignature.SignatureDocumentMajorVersion + "." + dbSignature.SignatureDocumentMinorVersion,
			SignatureReferenceType:      dbSignature.SignatureReferenceType,
			ProjectID:                   dbSignature.SignatureProjectID,
			Created:                     dbSignature.DateCreated,
			Modified:                    dbSignature.DateModified,
			EmailWhitelist:              dbSignature.EmailWhitelist,
			DomainWhitelist:             dbSignature.DomainWhitelist,
			GithubWhitelist:             dbSignature.GitHubWhitelist,
			GithubOrgWhitelist:          dbSignature.GitHubOrgWhitelist,
		}
		signatures = append(signatures, sig)
		go func(sigModel *models.Signature, signatureUserCompanyID string, sigACL []string) {
			defer wg.Done()
			var companyName = ""
			var userName = ""
			var userLFID = ""
			var userGHID = ""
			var swg sync.WaitGroup
			swg.Add(2)

			go func() {
				defer swg.Done()
				if sigModel.SignatureReferenceType == "user" {
					userModel, userErr := repo.usersRepo.GetUser(sigModel.SignatureReferenceID)
					if userErr != nil || userModel == nil {
						log.Warnf("unable to lookup user using id: %s, error: %v", sigModel.SignatureReferenceID, userErr)
					} else {
						userName = userModel.Username
						userLFID = userModel.LfUsername
						userGHID = userModel.GithubID
					}
					if signatureUserCompanyID != "" {
						dbCompanyModel, companyErr := repo.companyRepo.GetCompany(signatureUserCompanyID)
						if companyErr != nil {
							log.Warnf("unable to lookup company using id: %s, error: %v", signatureUserCompanyID, companyErr)
						} else {
							companyName = dbCompanyModel.CompanyName
						}
					}
				} else if sigModel.SignatureReferenceType == "company" {
					dbCompanyModel, companyErr := repo.companyRepo.GetCompany(sigModel.SignatureReferenceID)
					if companyErr != nil {
						log.Warnf("unable to lookup company using id: %s, error: %v", sigModel.SignatureReferenceID, companyErr)
					} else {
						companyName = dbCompanyModel.CompanyName
					}
				}
			}()

			var signatureACL []models.User
			go func() {
				defer swg.Done()
				for _, userName := range sigACL {
					if loadACLDetails {
						userModel, userErr := repo.usersRepo.GetUserByUserName(userName, true)
						if userErr != nil {
							log.Warnf("unable to lookup user using username: %s, error: %v", userName, userErr)
						} else {
							if userModel == nil {
								log.Warnf("User looking for username is null: %s for signature: %s", userName, sigModel.SignatureID)
							} else {
								signatureACL = append(signatureACL, *userModel)
							}
						}
					} else {
						signatureACL = append(signatureACL, models.User{LfUsername: userName})
					}
				}
			}()
			swg.Wait()
			sigModel.CompanyName = companyName
			sigModel.UserName = userName
			sigModel.UserLFID = userLFID
			sigModel.UserGHID = userGHID
			sigModel.SignatureACL = signatureACL
		}(sig, dbSignature.SignatureUserCompanyID, dbSignature.SignatureACL)
	}
	wg.Wait()
	return signatures, nil
}

// buildResponse is a helper function which converts a database model to a GitHub organization response model
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

// buildProject is a helper function to build a common set of projection/columns for the query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("signature_id"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("signature_acl"),
		expression.Name("signature_approved"),
		expression.Name("signature_document_major_version"),
		expression.Name("signature_document_minor_version"),
		expression.Name("signature_reference_id"),
		expression.Name("signature_reference_name"),       // Added to support simplified UX queries
		expression.Name("signature_reference_name_lower"), // Added to support case insensitive UX queries
		expression.Name("signature_project_id"),
		expression.Name("signature_reference_type"),       // user or company
		expression.Name("signature_signed"),               // T/F
		expression.Name("signature_type"),                 // ccla or cla
		expression.Name("signature_user_ccla_company_id"), // reference to the company
		expression.Name("email_whitelist"),
		expression.Name("domain_whitelist"),
		expression.Name("github_whitelist"),
		expression.Name("github_org_whitelist"),
	)
}

// buildSignatureACLProject is a helper function to build a signature ACL response/projection
func buildSignatureACLProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("signature_id"),
		expression.Name("signature_acl"),
	)
}
