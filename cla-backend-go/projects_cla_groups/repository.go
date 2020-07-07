package projects_cla_groups

import (
	"errors"
	"fmt"
	"strings"

	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// constants
const (
	CLAGroupIDIndex     = "cla-group-id-index"
	FoundationSFIDIndex = "foundation-sfid-index"
	NotDefined          = "Not Defined"
	NotFound            = "Not Found"
)

// errors
var (
	ErrProjectNotAssociatedWithClaGroup = errors.New("provided project is not associated with cla_group")
	ErrAssociationAlreadyExist          = errors.New("cla_group project association already exist")
	ErrCLAGroupDoesNotExist             = errors.New("cla group does not exist")
)

// ProjectClaGroup is database model for projects_cla_group table
type ProjectClaGroup struct {
	ProjectSFID       string `json:"project_sfid"`
	ProjectName       string `json:"project_name"`
	ClaGroupID        string `json:"cla_group_id"`
	ClaGroupName      string `json:"cla_group_name"`
	FoundationSFID    string `json:"foundation_sfid"`
	FoundationName    string `json:"foundation_name"`
	RepositoriesCount int64  `json:"repositories_count"`
	Version           string `json:"version"`
}

// Repository provides interface for interacting with project_cla_groups table
type Repository interface {
	GetClaGroupIDForProject(projectSFID string) (*ProjectClaGroup, error)
	GetProjectsIdsForClaGroup(claGroupID string) ([]*ProjectClaGroup, error)
	GetProjectsIdsForFoundation(foundationSFID string) ([]*ProjectClaGroup, error)
	GetProjectsIdsForAllFoundation() ([]*ProjectClaGroup, error)
	AssociateClaGroupWithProject(claGroupID string, projectSFID string, foundationSFID string) error
	RemoveProjectAssociatedWithClaGroup(claGroupID string, projectSFIDList []string, all bool) error
	getCLAGroupNameByID(claGroupID string) (string, error)
}

type repo struct {
	tableName      string
	dynamoDBClient *dynamodb.DynamoDB
	stage          string
}

// NewRepository provides implementation of projects_cla_group repository
func NewRepository(awsSession *session.Session, stage string) Repository {
	return &repo{
		tableName:      fmt.Sprintf("cla-%s-projects-cla-groups", stage),
		dynamoDBClient: dynamodb.New(awsSession),
		stage:          stage,
	}
}

func (repo *repo) queryClaGroupsProjects(keyCondition expression.KeyConditionBuilder, indexName *string) ([]*ProjectClaGroup, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		log.Warnf("error building expression for project cla groups, error: %v", err)
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(repo.tableName),
		IndexName:                 indexName,
	}

	var projectClaGroups []*ProjectClaGroup
	for {
		results, errQuery := repo.dynamoDBClient.Query(queryInput)
		if errQuery != nil {
			log.Warnf("error retrieving project cla-groups, error: %v", errQuery)
			return nil, errQuery
		}

		var projectClaGroupsTmp []*ProjectClaGroup

		err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &projectClaGroupsTmp)
		if err != nil {
			log.Warnf("error unmarshalling project cla-groups from database. error: %v", err)
			return nil, err
		}
		projectClaGroups = append(projectClaGroups, projectClaGroupsTmp...)

		if len(results.LastEvaluatedKey) != 0 {
			queryInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	return projectClaGroups, nil
}

func (repo *repo) GetClaGroupIDForProject(projectSFID string) (*ProjectClaGroup, error) {
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"project_sfid": {
				S: aws.String(projectSFID),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, ErrProjectNotAssociatedWithClaGroup
	}
	var out ProjectClaGroup
	err = dynamodbattribute.UnmarshalMap(result.Item, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *repo) GetProjectsIdsForClaGroup(claGroupID string) ([]*ProjectClaGroup, error) {
	keyCondition := expression.Key("cla_group_id").Equal(expression.Value(claGroupID))
	return repo.queryClaGroupsProjects(keyCondition, aws.String(CLAGroupIDIndex))
}

func (repo *repo) GetProjectsIdsForFoundation(foundationSFID string) ([]*ProjectClaGroup, error) {
	keyCondition := expression.Key("foundation_sfid").Equal(expression.Value(foundationSFID))
	return repo.queryClaGroupsProjects(keyCondition, aws.String(FoundationSFIDIndex))
}

func (repo *repo) GetProjectsIdsForAllFoundation() ([]*ProjectClaGroup, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}
	var resultList []map[string]*dynamodb.AttributeValue
	for {
		results, err := repo.dynamoDBClient.Scan(scanInput) //nolint
		if err != nil {
			log.Warnf("error retrieving %s, error: %v", repo.tableName, err)
			return nil, err
		}
		resultList = append(resultList, results.Items...)
		if len(results.LastEvaluatedKey) != 0 {
			scanInput.ExclusiveStartKey = results.LastEvaluatedKey
		} else {
			break
		}
	}
	var output []*ProjectClaGroup
	err := dynamodbattribute.UnmarshalListOfMaps(resultList, &output)
	if err != nil {
		log.Warnf("error unmarshalling %s from database. error: %v", repo.tableName, err)
		return nil, err
	}
	return output, nil
}

// AssociateClaGroupWithProject creates entry in db to track cla_group association with project/foundation
func (repo *repo) AssociateClaGroupWithProject(claGroupID string, projectSFID string, foundationSFID string) error {
	var foundationName = NotDefined
	if foundationSFID == "" {
		log.Debugf("AssociateClaGroupWithProject - foundationSFID is empty - setting foundation name to 'LF Supported'")
		// Special case - "LF Supported" project
		foundationName = "LF Supported"
	} else {
		// Lookup the foundation name
		projectServiceModel, projErr := v2ProjectService.GetClient().GetProject(foundationSFID)
		if projErr != nil {
			log.Warnf("unable to lookup foundation SFID: %s - error: %+v - using '%s'",
				foundationSFID, projErr, NotDefined)
		} else {
			foundationName = projectServiceModel.Name
		}
	}

	// Lookup the project name
	var projectName = NotDefined
	projectServiceModel, projErr := v2ProjectService.GetClient().GetProject(projectSFID)
	if projErr != nil {
		log.Warnf("unable to lookup project SFID: %s - error: %+v - using '%s'",
			projectSFID, projErr, NotDefined)
	} else {
		projectName = projectServiceModel.Name
	}

	// Lookup the CLA Group name/Project Name
	claGroupName, claGroupLookupErr := repo.getCLAGroupNameByID(claGroupID)
	if claGroupLookupErr != nil {
		claGroupName = NotDefined
		log.Warnf("unable to lookup CLA Group ID/Project ID: %s - error: %+v - using '%s'",
			claGroupID, claGroupLookupErr, NotDefined)
	}

	input := &ProjectClaGroup{
		ProjectSFID:    projectSFID,
		ProjectName:    projectName,
		ClaGroupID:     claGroupID,
		ClaGroupName:   claGroupName,
		FoundationSFID: foundationSFID,
		FoundationName: foundationName,
		Version:        "v1",
	}

	av, err := dynamodbattribute.MarshalMap(input)
	if err != nil {
		return err
	}
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(repo.tableName),
		ConditionExpression: aws.String("attribute_not_exists(project_sfid)"),
	})
	if err != nil {
		log.Error(fmt.Sprintf("cannot put association entry of cla_group_id: %s, project_sfid: %s in dynamodb",
			claGroupID, projectSFID), err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return ErrAssociationAlreadyExist
			}
			return err
		}
	}
	return nil
}

// RemoveProjectAssociatedWithClaGroup removes all associated project with cla_group
func (repo *repo) RemoveProjectAssociatedWithClaGroup(claGroupID string, projectSFIDList []string, all bool) error {
	list, err := repo.GetProjectsIdsForClaGroup(claGroupID)
	if err != nil {
		return err
	}
	var projectFilter *utils.StringSet
	if !all {
		projectFilter = utils.NewStringSetFromStringArray(projectSFIDList)
	}
	var errs []string
	for _, pr := range list {
		if !all && !projectFilter.Include(pr.ProjectSFID) {
			// ignore project not present in projectSFIDList
			continue
		}
		_, err = repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"project_sfid": {S: aws.String(pr.ProjectSFID)},
			},
			TableName: aws.String(repo.tableName),
		})
		if err != nil {
			log.Warnf("unable to delete cla_group_association cla_group_id:%s project_sfid:%s", claGroupID, pr.ProjectSFID)
			errs = append(errs, err.Error())
		}
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ","))
	}
	return nil
}

// getCLAGroupNameByID helper function to fetch the CLA Group name
func (repo *repo) getCLAGroupNameByID(claGroupID string) (string, error) {
	tableName := fmt.Sprintf("cla-%s-projects", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"repository_id": {
				S: aws.String(claGroupID),
			},
		},
	})
	if err != nil {
		return NotFound, err
	}
	if len(result.Item) == 0 {
		return NotFound, ErrCLAGroupDoesNotExist
	}

	// Quick model to grab the bare minimum values
	type claGroupIDNameModel struct {
		ProjectID   string `dynamodbav:"project_id"`
		ProjectName string `dynamodbav:"project_name"`
	}

	var claGroupModel claGroupIDNameModel
	err = dynamodbattribute.UnmarshalMap(result.Item, &claGroupModel)
	if err != nil {
		return NotFound, err
	}

	return claGroupModel.ProjectName, nil
}
