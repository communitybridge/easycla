// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/gerrits"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/gofrs/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// errors
var (
	ErrProjectDoesNotExist = errors.New("project does not exist")
	ErrProjectIDMissing    = errors.New("project id is missing")
)

// constants
const (
	LoadRepoDetails     = true
	DontLoadRepoDetails = false
)

// ProjectRepository defines functions of Project repository
type ProjectRepository interface { //nolint
	CreateCLAGroup(project *models.Project) (*models.Project, error)
	GetCLAGroupByID(projectID string, loadRepoDetails bool) (*models.Project, error)
	GetCLAGroupsByExternalID(params *project.GetProjectsByExternalIDParams, loadRepoDetails bool) (*models.Projects, error)
	GetCLAGroupByName(projectName string) (*models.Project, error)
	GetExternalCLAGroup(projectExternalID string) (*models.Project, error)
	GetCLAGroups(params *project.GetProjectsParams) (*models.Projects, error)
	DeleteCLAGroup(projectID string) error
	UpdateCLAGroup(projectModel *models.Project) (*models.Project, error)

	GetClaGroupsByFoundationSFID(foundationSFID string, loadRepoDetails bool) (*models.Projects, error)
	GetClaGroupByProjectSFID(projectSFID string, loadRepoDetails bool) (*models.Project, error)
	UpdateRootCLAGroupRepositoriesCount(claGroupID string, diff int64) error
}

// NewRepository creates instance of project repository
func NewRepository(awsSession *session.Session, stage string, ghRepo repositories.Repository, gerritRepo gerrits.Repository, projectClaGroupRepo projects_cla_groups.Repository) ProjectRepository {
	return &repo{
		dynamoDBClient:      dynamodb.New(awsSession),
		stage:               stage,
		ghRepo:              ghRepo,
		gerritRepo:          gerritRepo,
		projectClaGroupRepo: projectClaGroupRepo,
		claGroupTable:       fmt.Sprintf("cla-%s-projects", stage),
	}
}

type repo struct {
	stage               string
	dynamoDBClient      *dynamodb.DynamoDB
	ghRepo              repositories.Repository
	gerritRepo          gerrits.Repository
	projectClaGroupRepo projects_cla_groups.Repository
	claGroupTable       string
}

// CreateCLAGroup creates a new project
func (repo *repo) CreateCLAGroup(projectModel *models.Project) (*models.Project, error) {
	f := logrus.Fields{
		"function":          "CreateCLAGroup",
		"projectName":       projectModel.ProjectName,
		"projectExternalID": projectModel.ProjectExternalID,
		"foundationSFID":    projectModel.FoundationSFID,
		"tableName":         repo.claGroupTable}
	// Generate a new project ID
	projectID, err := uuid.NewV4()
	if err != nil {
		log.WithFields(f).Warnf("Unable to generate a UUID for a new project request, error: %v", err)
		return nil, err
	}
	f["projectID"] = projectID

	_, currentTimeString := utils.CurrentTime()
	input := &dynamodb.PutItemInput{
		Item:      map[string]*dynamodb.AttributeValue{},
		TableName: aws.String(repo.claGroupTable),
	}

	//var individualDocs []*dynamodb.AttributeValue
	//var corporateDocs []*dynamodb.AttributeValue
	addStringAttribute(input.Item, "project_id", projectID.String())
	addStringAttribute(input.Item, "project_external_id", projectModel.ProjectExternalID)
	addStringAttribute(input.Item, "foundation_sfid", projectModel.FoundationSFID)
	addStringAttribute(input.Item, "project_description", projectModel.ProjectDescription)
	addStringAttribute(input.Item, "project_name", projectModel.ProjectName)
	addStringAttribute(input.Item, "project_name_lower", strings.ToLower(projectModel.ProjectName))
	addStringSliceAttribute(input.Item, "project_acl", projectModel.ProjectACL)
	addBooleanAttribute(input.Item, "project_icla_enabled", projectModel.ProjectICLAEnabled)
	addBooleanAttribute(input.Item, "project_ccla_enabled", projectModel.ProjectCCLAEnabled)
	addBooleanAttribute(input.Item, "project_ccla_requires_icla_signature", projectModel.ProjectCCLARequiresICLA)

	// Empty documents for now - will add the template details later
	addListAttribute(input.Item, "project_corporate_documents", []*dynamodb.AttributeValue{})
	addListAttribute(input.Item, "project_individual_documents", []*dynamodb.AttributeValue{})
	addListAttribute(input.Item, "project_member_documents", []*dynamodb.AttributeValue{})

	addStringAttribute(input.Item, "date_created", currentTimeString)
	addStringAttribute(input.Item, "date_modified", currentTimeString)
	// Set the version attribute if not already set
	if projectModel.Version == "" {
		projectModel.Version = "v1" // default value
	}
	addStringAttribute(input.Item, "version", projectModel.Version)

	_, err = repo.dynamoDBClient.PutItem(input)
	if err != nil {
		log.WithFields(f).Warnf("Unable to create a new project record, error: %v", err)
		return nil, err
	}

	// Re-use the provided model - just update the dynamically assigned values
	projectModel.ProjectID = projectID.String()
	projectModel.DateCreated = currentTimeString
	projectModel.DateModified = currentTimeString

	return projectModel, nil
}

func (repo *repo) getCLAGroupByID(projectID string, loadCLAGroupDetails bool) (*models.Project, error) {
	f := logrus.Fields{
		"function":           "getCLAGroupByID",
		"projectID":          projectID,
		"loadProjectDetails": loadCLAGroupDetails,
		"tableName":          repo.claGroupTable}
	log.WithFields(f).Debugf("loading project")
	// This is the key we want to match
	condition := expression.Key("project_id").Equal(expression.Value(projectID))

	// Use the builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for CLA Group query, projectID: %s, error: %v",
			projectID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
	}

	// Make the DynamoDB Query API call
	results, queryErr := repo.dynamoDBClient.Query(queryInput)
	if queryErr != nil {
		log.WithFields(f).Warnf("error retrieving project by projectID: %s, error: %v", projectID, queryErr)
		return nil, queryErr
	}

	if len(results.Items) < 1 {
		return nil, ErrProjectDoesNotExist
	}
	var dbModel DBProjectModel
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &dbModel)
	if err != nil {
		log.WithFields(f).Warnf("error unmarshalling db project model, error: %+v", err)
		return nil, err
	}

	// Convert the database model to an API response model
	return repo.buildCLAGroupModel(dbModel, loadCLAGroupDetails), nil
}

// GetCLAGroupByID returns the project model associated for the specified projectID
func (repo *repo) GetCLAGroupByID(projectID string, loadRepoDetails bool) (*models.Project, error) {
	return repo.getCLAGroupByID(projectID, loadRepoDetails)
}

// GetCLAGroupsByExternalID queries the database and returns a list of the projects
func (repo *repo) GetCLAGroupsByExternalID(params *project.GetProjectsByExternalIDParams, loadRepoDetails bool) (*models.Projects, error) {
	f := logrus.Fields{
		"functionName":    "GetCLAGroupsByExternalID",
		"ProjectSFID":     params.ProjectSFID,
		"NextKey":         params.NextKey,
		"PageSize":        params.PageSize,
		"loadRepoDetails": loadRepoDetails,
		"tableName":       repo.claGroupTable}
	log.WithFields(f).Debugf("loading project")

	// This is the key we want to match
	condition := expression.Key("project_external_id").Equal(expression.Value(params.ProjectSFID))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for project query, error: %v", err)
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
		IndexName:                 aws.String("external-project-index"),
	}

	// If we have the next key, set the exclusive start key value
	if params.NextKey != nil && *params.NextKey != "" {
		log.WithFields(f).Debugf("Received a nextKey, value: %s", *params.NextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: params.NextKey,
			},
		}
	}

	var pageSize *int64
	// If we have a page size, set the limit value - make sure it's a positive value
	if params.PageSize != nil && *params.PageSize > 0 {
		log.WithFields(f).Debugf("Received a pageSize parameter, value: %d", *params.PageSize)
		pageSize = params.PageSize
	} else {
		// Default page size
		pageSize = aws.Int64(50)
	}
	queryInput.Limit = pageSize

	var projects []models.Project
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.WithFields(f).Warnf("error retrieving projects, error: %v", errQuery)
			return nil, errQuery
		}

		// Convert the list of DB models to a list of response models
		projectList, modelErr := repo.buildCLAGroupModels(results.Items, loadRepoDetails)
		if modelErr != nil {
			log.WithFields(f).Warnf("error converting project DB model to response model, error: %v",
				modelErr)
			return nil, modelErr
		}

		// Add to the project response models to the list
		projects = append(projects, projectList...)

		if results.LastEvaluatedKey["project_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["project_id"].S
			queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"project_id": {
					S: aws.String(lastEvaluatedKey),
				},
			}
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(projects)) >= *pageSize {
			break
		}
	}

	return &models.Projects{
		LastKeyScanned: lastEvaluatedKey,
		PageSize:       *pageSize,
		ResultCount:    int64(len(projects)),
		Projects:       projects,
	}, nil
}

// GetClaGroupsByFoundationID queries the database and returns a list of all cla_groups associated with foundation
func (repo *repo) GetClaGroupsByFoundationSFID(foundationSFID string, loadRepoDetails bool) (*models.Projects, error) {
	f := logrus.Fields{
		"functionName":    "GetClaGroupsByFoundationSFID",
		"foundationSFID":  foundationSFID,
		"loadRepoDetails": loadRepoDetails,
		"tableName":       repo.claGroupTable}
	log.WithFields(f).Debugf("loading project by foundation SFID")

	// This is the key we want to match
	condition := expression.Key("foundation_sfid").Equal(expression.Value(foundationSFID))

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for project scan, error: %v", err)
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
		IndexName:                 aws.String("foundation-sfid-project-name-index"),
	}

	var projects []models.Project
	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.WithFields(f).Warnf("error retrieving projects, error: %v", errQuery)
			return nil, errQuery
		}

		// Convert the list of DB models to a list of response models
		projectList, modelErr := repo.buildCLAGroupModels(results.Items, loadRepoDetails)
		if modelErr != nil {
			log.WithFields(f).Warnf("error converting project DB model to response model, error: %v",
				modelErr)
			return nil, modelErr
		}

		// Add to the project response models to the list
		projects = append(projects, projectList...)

		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}

	return &models.Projects{
		ResultCount: int64(len(projects)),
		Projects:    projects,
	}, nil
}

// GetClaGroupsByProjectSFID returns cla_group associated with project
func (repo *repo) GetClaGroupByProjectSFID(projectSFID string, loadRepoDetails bool) (*models.Project, error) {
	f := logrus.Fields{
		"functionName":    "GetClaGroupByProjectSFID",
		"projectSFID":     projectSFID,
		"loadRepoDetails": loadRepoDetails,
		"tableName":       repo.claGroupTable}
	log.WithFields(f).Debugf("loading project")

	claGroupProject, err := repo.projectClaGroupRepo.GetClaGroupIDForProject(projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("error fetching CLA Group ID for project, error: %v", err)
		return nil, err
	}

	return repo.getCLAGroupByID(claGroupProject.ClaGroupID, loadRepoDetails)
}

// GetCLAGroupByName returns the project model associated for the specified project name
func (repo *repo) GetCLAGroupByName(projectName string) (*models.Project, error) {
	f := logrus.Fields{
		"functionName": "GetCLAGroupByName",
		"projectName":  projectName,
		"tableName":    repo.claGroupTable}
	log.WithFields(f).Debugf("loading project")

	// This is the key we want to match
	condition := expression.Key("project_name_lower").Equal(expression.Value(strings.ToLower(projectName)))

	// Use the builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for CLAGroup query, projectName: %s, error: %v",
			projectName, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
		IndexName:                 aws.String("project-name-lower-search-index"),
	}

	// Make the DynamoDB Query API call
	results, queryErr := repo.dynamoDBClient.Query(queryInput)
	if queryErr != nil {
		log.WithFields(f).Warnf("error retrieving project by projectName: %s, error: %v", projectName, queryErr)
		return nil, queryErr
	}

	// Should only have one result
	if *results.Count > 1 {
		log.WithFields(f).Warnf("CLAGroup scan by name returned more than one result using projectName: %s", projectName)
	}

	// Didn't find it...
	if *results.Count == 0 {
		log.WithFields(f).Debugf("CLAGroup scan by name returned no results using projectName: %s", projectName)
		return nil, nil
	}

	// Found it...
	var dbModel DBProjectModel
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &dbModel)
	if err != nil {
		log.WithFields(f).Warnf("error unmarshalling db project model, error: %+v", err)
		return nil, err
	}

	// Convert the database model to an API response model
	return repo.buildCLAGroupModel(dbModel, LoadRepoDetails), nil
}

// GetExternalCLAGroup returns the project model associated for the specified external project ID
func (repo *repo) GetExternalCLAGroup(projectExternalID string) (*models.Project, error) {
	f := logrus.Fields{
		"functionName":      "GetExternalCLAGroup",
		"projectExternalID": projectExternalID,
		"tableName":         repo.claGroupTable}
	log.WithFields(f).Debugf("loading project")
	// This is the key we want to match
	condition := expression.Key("project_external_id").Equal(expression.Value(projectExternalID))

	// Use the builder to create the expression
	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for CLAGroup query, projectExternalID: %s, error: %v",
			projectExternalID, err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
		IndexName:                 aws.String("external-project-index"),
	}

	// Make the DynamoDB Query API call
	results, queryErr := repo.dynamoDBClient.Query(queryInput)
	if queryErr != nil {
		log.WithFields(f).Warnf("error retrieving project by projectExternalID: %s, error: %v", projectExternalID, queryErr)
		return nil, queryErr
	}

	// No match, didn't find it
	if *results.Count == 0 {
		return nil, nil
	}

	// Should only have one result
	if *results.Count > 1 {
		log.WithFields(f).Warnf("CLAGroup query returned more than one result using projectExternalID: %s", projectExternalID)
	}

	var dbModel DBProjectModel
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &dbModel)
	if err != nil {
		log.WithFields(f).Warnf("error unmarshalling db project model, error: %+v", err)
		return nil, err
	}

	// Convert the database model to an API response model
	return repo.buildCLAGroupModel(dbModel, LoadRepoDetails), nil
}

// GetCLAGroups queries the database and returns a list of the projects
func (repo *repo) GetCLAGroups(params *project.GetProjectsParams) (*models.Projects, error) {
	f := logrus.Fields{
		"functionName": "GetCLAGroups",
		"searchField":  params.SearchField,
		"searchTerm":   params.SearchTerm,
		"nextKey":      params.NextKey,
		"pageSize":     params.PageSize,
		"fullMatch":    params.FullMatch,
		"tableName":    repo.claGroupTable}
	log.WithFields(f).Debugf("searching project")

	// Use the nice builder to create the expression
	expr, err := expression.NewBuilder().WithProjection(buildProjection()).Build()
	if err != nil {
		log.WithFields(f).Warnf("error building expression for project scan, error: %v", err)
	}

	// Assemble the query input parameters
	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(repo.claGroupTable),
	}

	// If we have the next key, set the exclusive start key value
	if params.NextKey != nil && *params.NextKey != "" {
		log.WithFields(f).Debugf("Received a nextKey, value: %s", *params.NextKey)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: params.NextKey,
			},
		}
	}
	var pageSize int64
	// If we have a page size, set the limit value - make sure it's a positive value
	if params.PageSize != nil && *params.PageSize > 0 {
		log.WithFields(f).Debugf("Received a pageSize parameter, value: %d", *params.PageSize)
		// The primary key of the first item that this operation will evaluate.
		// and the query key (if not the same)
		scanInput.Limit = params.PageSize
	} else {
		// Default page size
		pageSize = 50
		params.PageSize = &pageSize
	}

	var projects []models.Project
	var lastEvaluatedKey string

	// Loop until we have all the records
	for ok := true; ok; ok = lastEvaluatedKey != "" {
		results, errQuery := repo.dynamoDBClient.Scan(scanInput)
		if errQuery != nil {
			log.WithFields(f).Warnf("error retrieving projects, error: %v", errQuery)
			return nil, errQuery
		}

		// Convert the list of DB models to a list of response models
		projectList, modelErr := repo.buildCLAGroupModels(results.Items, LoadRepoDetails)
		if modelErr != nil {
			log.WithFields(f).Warnf("error converting project DB model to response model, error: %v",
				modelErr)
			return nil, modelErr
		}

		// Add to the project response models to the list
		projects = append(projects, projectList...)

		if results.LastEvaluatedKey["project_id"] != nil {
			lastEvaluatedKey = *results.LastEvaluatedKey["project_id"].S
			scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
				"project_id": {
					S: aws.String(lastEvaluatedKey),
				},
			}
		} else {
			lastEvaluatedKey = ""
		}

		if int64(len(projects)) >= *params.PageSize {
			break
		}
	}

	return &models.Projects{
		LastKeyScanned: lastEvaluatedKey,
		PageSize:       *params.PageSize,
		Projects:       projects,
	}, nil
}

// DeleteCLAGroup deletes the CLAGroup by projectID
func (repo *repo) DeleteCLAGroup(projectID string) error {
	f := logrus.Fields{
		"functionName": "DeleteCLAGroup",
		"projectID":    projectID,
		"tableName":    repo.claGroupTable}
	log.WithFields(f).Debugf("deleting CLA Group")

	existingCLAGroup, getErr := repo.GetCLAGroupByID(projectID, DontLoadRepoDetails)
	if getErr != nil {
		log.WithFields(f).Warnf("delete - error locating the CLA Group, error: %+v", getErr)
		return getErr
	}

	if existingCLAGroup == nil {
		log.WithFields(f).Warn("unable to locate CLA Group by ID - CLA Group does not exist")
		return ErrProjectDoesNotExist
	}

	var deleteErr error
	// Perform the delete
	_, deleteErr = repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(repo.claGroupTable),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(existingCLAGroup.ProjectID),
			},
		},
	})

	if deleteErr != nil {
		log.WithFields(f).Warnf("Error deleting project with CLA Group ID : %s, error: %v", projectID, deleteErr)
		return deleteErr
	}

	return nil
}

// UpdateCLAGroup updates the project by projectID
func (repo *repo) UpdateCLAGroup(projectModel *models.Project) (*models.Project, error) {
	f := logrus.Fields{
		"functionName":            "UpdateCLAGroup",
		"ProjectID":               projectModel.ProjectID,
		"ProjectName":             projectModel.ProjectName,
		"FoundationSFID":          projectModel.FoundationSFID,
		"ProjectExternalID":       projectModel.ProjectExternalID,
		"ProjectICLAEnabled":      projectModel.ProjectICLAEnabled,
		"ProjectCCLAEnabled":      projectModel.ProjectCCLAEnabled,
		"ProjectCCLARequiresICLA": projectModel.ProjectCCLARequiresICLA,
		"tableName":               repo.claGroupTable}
	log.WithFields(f).Debugf("updating CLA Group")

	if projectModel.ProjectID == "" {
		return nil, ErrProjectIDMissing
	}

	existingCLAGroup, getErr := repo.GetCLAGroupByID(projectModel.ProjectID, DontLoadRepoDetails)
	if getErr != nil {
		log.WithFields(f).Warnf("update - error locating the project id: %s, error: %+v", projectModel.ProjectID, getErr)
		return nil, getErr
	}

	if existingCLAGroup == nil {
		return nil, ErrProjectDoesNotExist
	}

	expressionAttributeNames := map[string]*string{}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{}
	updateExpression := "SET "

	if projectModel.ProjectName != "" {
		log.WithFields(f).Debugf("adding project_name: %s", projectModel.ProjectName)
		expressionAttributeNames["#N"] = aws.String("project_name")
		expressionAttributeValues[":n"] = &dynamodb.AttributeValue{S: aws.String(projectModel.ProjectName)}
		updateExpression = updateExpression + " #N = :n, "
		log.WithFields(f).Debugf("adding project name lower: %s", strings.ToLower(projectModel.ProjectName))
		expressionAttributeNames["#LOW"] = aws.String("project_name_lower")
		expressionAttributeValues[":low"] = &dynamodb.AttributeValue{S: aws.String(strings.ToLower(projectModel.ProjectName))}
		updateExpression = updateExpression + " #LOW = :low, "
	}

	if projectModel.ProjectDescription != "" {
		log.WithFields(f).Debugf("adding project_description: %s", projectModel.ProjectDescription)
		expressionAttributeNames["#DESC"] = aws.String("project_description")
		expressionAttributeValues[":desc"] = &dynamodb.AttributeValue{S: aws.String(projectModel.ProjectDescription)}
		updateExpression = updateExpression + " #DESC = :desc, "
	}

	if projectModel.ProjectACL != nil && len(projectModel.ProjectACL) > 0 {
		log.WithFields(f).Debugf("adding project_acl: %s", projectModel.ProjectACL)
		expressionAttributeNames["#A"] = aws.String("project_acl")
		expressionAttributeValues[":a"] = &dynamodb.AttributeValue{SS: aws.StringSlice(projectModel.ProjectACL)}
		updateExpression = updateExpression + " #A = :a, "
	}

	log.WithFields(f).Debugf("adding project_icla_enabled: %t", projectModel.ProjectICLAEnabled)
	expressionAttributeNames["#I"] = aws.String("project_icla_enabled")
	expressionAttributeValues[":i"] = &dynamodb.AttributeValue{BOOL: aws.Bool(projectModel.ProjectICLAEnabled)}
	updateExpression = updateExpression + " #I = :i, "

	log.WithFields(f).Debugf("adding project_ccla_enabled: %t", projectModel.ProjectCCLAEnabled)
	expressionAttributeNames["#C"] = aws.String("project_ccla_enabled")
	expressionAttributeValues[":c"] = &dynamodb.AttributeValue{BOOL: aws.Bool(projectModel.ProjectCCLAEnabled)}
	updateExpression = updateExpression + " #C = :c, "

	log.WithFields(f).Debugf("adding project_ccla_requires_icla_signature: %t", projectModel.ProjectCCLARequiresICLA)
	expressionAttributeNames["#CI"] = aws.String("project_ccla_requires_icla_signature")
	expressionAttributeValues[":ci"] = &dynamodb.AttributeValue{BOOL: aws.Bool(projectModel.ProjectCCLARequiresICLA)}
	updateExpression = updateExpression + " #CI = :ci, "

	_, currentTimeString := utils.CurrentTime()
	log.WithFields(f).Debugf("adding date_modified: %s", currentTimeString)
	expressionAttributeNames["#M"] = aws.String("date_modified")
	expressionAttributeValues[":m"] = &dynamodb.AttributeValue{S: aws.String(currentTimeString)}
	updateExpression = updateExpression + " #M = :m "

	// Assemble the query input parameters
	updateInput := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(existingCLAGroup.ProjectID),
			},
		},
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		UpdateExpression:          &updateExpression,
		TableName:                 aws.String(repo.claGroupTable),
	}
	//log.Debugf("Update input: %+V", updateInput.GoString())

	// Make the DynamoDB Update API call
	_, updateErr := repo.dynamoDBClient.UpdateItem(updateInput)
	if updateErr != nil {
		log.WithFields(f).Warnf("error updating CLAGroup by projectID: %s, error: %v", projectModel.ProjectID, updateErr)
		return nil, updateErr
	}

	// Read the updated record back from the DB and return - probably could
	// just create/update a new model in memory and return it to make it fast,
	// but this approach return exactly what the DB has
	return repo.GetCLAGroupByID(projectModel.ProjectID, LoadRepoDetails)
}

func (repo *repo) UpdateRootCLAGroupRepositoriesCount(claGroupID string, diff int64) error {
	val := strconv.FormatInt(diff, 10)
	updateExp := "ADD root_project_repositories_count :val"
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":val": {N: aws.String(val)}},
		UpdateExpression:          aws.String(updateExp),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {S: aws.String(claGroupID)},
		},
		TableName: aws.String(repo.claGroupTable),
	}
	_, err := repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.WithField("cla_group_id", claGroupID).Error("unable to update repositories count", err)
	}
	return err
}

// buildCLAGroupModels converts the database response model into an API response data model
func (repo *repo) buildCLAGroupModels(results []map[string]*dynamodb.AttributeValue, loadRepoDetails bool) ([]models.Project, error) {
	var projects []models.Project

	// The DB project model
	var dbProjects []DBProjectModel

	err := dynamodbattribute.UnmarshalListOfMaps(results, &dbProjects)
	if err != nil {
		log.Warnf("error unmarshalling projects from database, error: %v", err)
		return nil, err
	}

	// Create an output channel to receive the results
	responseChannel := make(chan *models.Project)

	// For each project, convert to a response model - using a go routine
	for _, dbProject := range dbProjects {
		go func(dbProject DBProjectModel) {
			// Send the results to the output channel
			responseChannel <- repo.buildCLAGroupModel(dbProject, loadRepoDetails)
		}(dbProject)
	}

	// Append all the responses to our list
	for i := 0; i < len(dbProjects); i++ {
		projects = append(projects, *<-responseChannel)
	}

	return projects, nil
}

// buildCLAGroupModel maps the database model to the API response model
func (repo *repo) buildCLAGroupModel(dbModel DBProjectModel, loadRepoDetails bool) *models.Project {

	var ghOrgs []*models.GithubRepositoriesGroupByOrgs
	var gerrits []*models.Gerrit

	if loadRepoDetails {
		if dbModel.ProjectID != "" {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				var err error
				ghOrgs, err = repo.ghRepo.GetProjectRepositoriesGroupByOrgs(dbModel.ProjectID)
				if err != nil {
					log.Warnf("buildPCLAGroupModel - unable to load GH organizations by project ID: %s, error: %+v",
						dbModel.ProjectID, err)
					// Reset to empty array
					ghOrgs = make([]*models.GithubRepositoriesGroupByOrgs, 0)
				}
			}()

			go func() {
				defer wg.Done()
				var err error
				var gerritsList *models.GerritList
				gerritsList, err = repo.gerritRepo.GetClaGroupGerrits(dbModel.ProjectID, nil)
				if err != nil {
					log.Warnf("buildCLAGroupModel - unable to load Gerrit repositories by project ID: %s, error: %+v",
						dbModel.ProjectID, err)
					// Reset to empty array
					gerrits = make([]*models.Gerrit, 0)
					return
				}
				gerrits = gerritsList.List
			}()
			wg.Wait()
		} else {
			log.Warnf("buildCLAGroupModel - project ID missing for project '%s' - ID: %s - unable to load GH and Gerrit repository details",
				dbModel.ProjectName, dbModel.ProjectID)
		}
	}

	return &models.Project{
		ProjectID:                    dbModel.ProjectID,
		FoundationSFID:               dbModel.FoundationSFID,
		RootProjectRepositoriesCount: dbModel.RootProjectRepositoriesCount,
		ProjectExternalID:            dbModel.ProjectExternalID,
		ProjectName:                  dbModel.ProjectName,
		ProjectDescription:           dbModel.ProjectDescription,
		ProjectACL:                   dbModel.ProjectACL,
		ProjectCCLAEnabled:           dbModel.ProjectCclaEnabled,
		ProjectICLAEnabled:           dbModel.ProjectIclaEnabled,
		ProjectCCLARequiresICLA:      dbModel.ProjectCclaRequiresIclaSignature,
		ProjectCorporateDocuments:    repo.buildCLAGroupDocumentModels(dbModel.ProjectCorporateDocuments),
		ProjectIndividualDocuments:   repo.buildCLAGroupDocumentModels(dbModel.ProjectIndividualDocuments),
		ProjectMemberDocuments:       repo.buildCLAGroupDocumentModels(dbModel.ProjectMemberDocuments),
		GithubRepositories:           ghOrgs,
		Gerrits:                      gerrits,
		DateCreated:                  dbModel.DateCreated,
		DateModified:                 dbModel.DateModified,
		Version:                      dbModel.Version,
	}
}

// buildCLAGroupDocumentModels builds response models based on the array of db models
func (repo *repo) buildCLAGroupDocumentModels(dbDocumentModels []DBProjectDocumentModel) []models.ProjectDocument {
	if dbDocumentModels == nil {
		return nil
	}

	// Response model
	var response []models.ProjectDocument

	for _, dbDocumentModel := range dbDocumentModels {
		response = append(response, models.ProjectDocument{
			DocumentName:            dbDocumentModel.DocumentName,
			DocumentAuthorName:      dbDocumentModel.DocumentAuthorName,
			DocumentContentType:     dbDocumentModel.DocumentContentType,
			DocumentFileID:          dbDocumentModel.DocumentFileID,
			DocumentLegalEntityName: dbDocumentModel.DocumentLegalEntityName,
			DocumentPreamble:        dbDocumentModel.DocumentPreamble,
			DocumentS3URL:           dbDocumentModel.DocumentS3URL,
			DocumentMajorVersion:    dbDocumentModel.DocumentMajorVersion,
			DocumentMinorVersion:    dbDocumentModel.DocumentMinorVersion,
			DocumentCreationDate:    dbDocumentModel.DocumentCreationDate,
		})
	}

	return response
}

// buildProject is a helper function to build a common set of projection/columns for the query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("project_id"),
		expression.Name("foundation_sfid"),
		expression.Name("root_project_repositories_count"),
		expression.Name("project_external_id"),
		expression.Name("project_name"),
		expression.Name("project_name_lower"),
		expression.Name("project_description"),
		expression.Name("project_acl"),
		expression.Name("project_ccla_enabled"),
		expression.Name("project_icla_enabled"),
		expression.Name("project_ccla_requires_icla_signature"),
		expression.Name("project_corporate_documents"),
		expression.Name("project_individual_documents"),
		expression.Name("project_member_documents"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
	)
}

// addStringAttribute adds a new string attribute to the existing map
func addStringAttribute(item map[string]*dynamodb.AttributeValue, key string, value string) {
	if value != "" {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
}

// addBooleanAttribute adds a new boolean attribute to the existing map
func addBooleanAttribute(item map[string]*dynamodb.AttributeValue, key string, value bool) {
	item[key] = &dynamodb.AttributeValue{BOOL: aws.Bool(value)}
}

// addStringSliceAttribute adds a new string slice attribute to the existing map
func addStringSliceAttribute(item map[string]*dynamodb.AttributeValue, key string, value []string) {
	item[key] = &dynamodb.AttributeValue{SS: aws.StringSlice(value)}
}

// addListAttribute adds a list to the existing map
func addListAttribute(item map[string]*dynamodb.AttributeValue, key string, value []*dynamodb.AttributeValue) {
	item[key] = &dynamodb.AttributeValue{L: value}
}
