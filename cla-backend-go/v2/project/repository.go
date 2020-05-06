package project

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gerrits"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
)

// Repository defines functions of v2 Project repository
type Repository interface {
	GetCCLAProjectsByExternalID(params *project.GetCCLAProjectsByExternalIDParams) (*models.Projects, error)
}

type repo struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
	ghRepo         repositories.Repository
	gerritRepo     gerrits.Repository
}

// NewRepository creates instance of project repository
func NewRepository(awsSession *session.Session, stage string, ghRepo repositories.Repository, gerritRepo gerrits.Repository) Repository {
	return &repo{
		dynamoDBClient: dynamodb.New(awsSession),
		stage:          stage,
		ghRepo:         ghRepo,
		gerritRepo:     gerritRepo,
	}
}

func (repo *repo) GetCCLAProjectsByExternalID(params *project.GetCCLAProjectsByExternalIDParams) (*models.Projects, error) {
	log.Debugf("Project - Repository Service - GetCCLAProjectsByExternalID - ExternalID: %s", params.ExternalID)
	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)

	condition := expression.Key("project_external_id").Equal(expression.Value(params.ExternalID))

	filter := expression.Name("project_ccla_enabled").Equal(expression.Value(true))

	expr, err := expression.NewBuilder().WithKeyCondition(condition).WithFilter(filter).WithProjection(buildProjection()).Build()

	if err != nil {
		log.Warnf("error building expression for Project query, projectExternalID : %s, error: %s ", params.ExternalID, err)
		return nil, err
	}

	// Assemble the query input
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("external-project-index"),
	}

	// If we have the next key, set the exclusive start key value
	if params.NextKey != nil && *params.NextKey != "" {
		log.Debugf("Received a nextKey, value: %s", *params.NextKey)
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
		log.Debugf("Received a pageSize parameter, value: %d", *params.PageSize)
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
			log.Warnf("error retrieving projects, error: %v", errQuery)
			return nil, errQuery
		}

		// Convert the list of DB models to a list of response models
		projectList, modelErr := repo.buildProjectModels(results.Items)
		if modelErr != nil {
			log.Warnf("error converting project DB model to response model, error: %v",
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

// buildProject is a helper function to build a common set of projection/columns for the query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("project_id"),
		expression.Name("project_external_id"),
		expression.Name("project_name"),
		expression.Name("project_name_lower"),
		expression.Name("project_acl"),
		expression.Name("project_ccla_enabled"),
		expression.Name("project_icla_enabled"),
		expression.Name("project_ccla_requires_icla_signature"),
		expression.Name("project_member_documents"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
	)
}

// buildProjectModels converts the database response model into an API response data model
func (repo *repo) buildProjectModels(results []map[string]*dynamodb.AttributeValue) ([]models.Project, error) {
	var projects []models.Project

	// The DB project model
	var dbProjects []v1Project.DBProjectModel

	err := dynamodbattribute.UnmarshalListOfMaps(results, &dbProjects)
	if err != nil {
		log.Warnf("error unmarshalling projects from database, error: %v", err)
		return nil, err
	}

	// Create an output channel to receive the results
	responseChannel := make(chan *models.Project)

	// For each project, convert to a response model - using a go routine
	for _, dbProject := range dbProjects {
		go func(dbProject v1Project.DBProjectModel) {
			// Send the results to the output channel
			responseChannel <- repo.buildProjectModel(dbProject)
		}(dbProject)
	}

	// Append all the responses to our list
	for i := 0; i < len(dbProjects); i++ {
		projects = append(projects, *<-responseChannel)
	}

	return projects, nil
}

// buildProjectModel maps the database model to the API response model
func (repo *repo) buildProjectModel(dbModel v1Project.DBProjectModel) *models.Project {
	var ghOrgs []*models.GithubRepositoriesGroupByOrgs
	var gerrits []*models.Gerrit

	if dbModel.ProjectID != "" {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			var err error
			ghOrgs, err = repo.ghRepo.GetProjectRepositoriesGroupByOrgs(dbModel.ProjectID)
			if err != nil {
				log.Warnf("buildProjectModel - unable to load GH organizations by project ID: %s, error: %+v",
					dbModel.ProjectID, err)
				// Reset to empty array
				ghOrgs = make([]*models.GithubRepositoriesGroupByOrgs, 0)
			}
		}()

		go func() {
			defer wg.Done()
			var err error
			gerrits, err = repo.gerritRepo.GetProjectGerrits(dbModel.ProjectID)
			if err != nil {
				log.Warnf("buildProjectModel - unable to load Gerrit repositories by project ID: %s, error: %+v",
					dbModel.ProjectID, err)
				// Reset to empty array
				gerrits = make([]*models.Gerrit, 0)
			}
		}()
		wg.Wait()
	} else {
		log.Warnf("buildProjectModel - project ID missing for project '%s' - ID: %s - unable to load GH and Gerrit repository details",
			dbModel.ProjectName, dbModel.ProjectID)
	}

	return &models.Project{
		ProjectID:                  dbModel.ProjectID,
		ProjectExternalID:          dbModel.ProjectExternalID,
		ProjectName:                dbModel.ProjectName,
		ProjectACL:                 dbModel.ProjectACL,
		ProjectCCLAEnabled:         dbModel.ProjectCclaEnabled,
		ProjectICLAEnabled:         dbModel.ProjectIclaEnabled,
		ProjectCCLARequiresICLA:    dbModel.ProjectCclaRequiresIclaSignature,
		ProjectCorporateDocuments:  repo.buildProjectDocumentModels(dbModel.ProjectCorporateDocuments),
		ProjectIndividualDocuments: repo.buildProjectDocumentModels(dbModel.ProjectIndividualDocuments),
		ProjectMemberDocuments:     repo.buildProjectDocumentModels(dbModel.ProjectMemberDocuments),
		GithubRepositories:         ghOrgs,
		Gerrits:                    gerrits,
		DateCreated:                dbModel.DateCreated,
		DateModified:               dbModel.DateModified,
		Version:                    dbModel.Version,
	}
}

// buildProjectDocumentModels builds response models based on the array of db models
func (repo *repo) buildProjectDocumentModels(dbDocumentModels []v1Project.DBProjectDocumentModel) []models.ProjectDocument {
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
