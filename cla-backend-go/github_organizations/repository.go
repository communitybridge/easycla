// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/github"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// indexes
const (
	GithubOrgSFIDIndex               = "github-org-sfid-index"
	GithubOrgLowerNameIndex          = "organization-name-lower-search-index"
	ProjectSFIDOrganizationNameIndex = "project-sfid-organization-name-index"
)

// errors
var (
	ErrOrganizationDoesNotExist = errors.New("github organization does not exist in cla")
)

// Repository interface defines the functions for the github organizations data model
type Repository interface {
	GetGithubOrganizations(ctx context.Context, externalProjectID string, projectSFID string) (*models.GithubOrganizations, error)
	GetGithubOrganizationByName(ctx context.Context, githubOrganizationName string) (*models.GithubOrganizations, error)
	AddGithubOrganization(ctx context.Context, externalProjectID string, projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(ctx context.Context, externalProjectID string, projectSFID string, githubOrgName string) error
	UpdateGithubOrganization(ctx context.Context, projectSFID string, organizationName string, autoEnabled bool, branchProtectionEnabled bool) error
	GetGithubOrganization(ctx context.Context, githubOrganizationName string) (*models.GithubOrganization, error)
}

type repository struct {
	stage              string
	dynamoDBClient     *dynamodb.DynamoDB
	githubOrgTableName string
}

// NewRepository creates a new instance of the githubOrganizations repository
func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:              stage,
		dynamoDBClient:     dynamodb.New(awsSession),
		githubOrgTableName: fmt.Sprintf("cla-%s-github-orgs", stage),
	}
}
func (repo repository) AddGithubOrganization(ctx context.Context, externalProjectID string, projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error) {
	f := logrus.Fields{
		"functionName":            "AddGithubOrganization",
		utils.XREQUESTID:          ctx.Value(utils.XREQUESTID),
		"externalProjectID":       externalProjectID,
		"projectSFID":             projectSFID,
		"organizationName":        aws.StringValue(input.OrganizationName),
		"autoEnabled":             aws.BoolValue(input.AutoEnabled),
		"branchProtectionEnabled": aws.BoolValue(input.BranchProtectionEnabled),
	}

	_, currentTime := utils.CurrentTime()
	githubOrg := &GithubOrganization{
		DateCreated:                currentTime,
		DateModified:               currentTime,
		OrganizationInstallationID: 0,
		OrganizationName:           *input.OrganizationName,
		OrganizationNameLower:      strings.ToLower(*input.OrganizationName),
		OrganizationSfid:           externalProjectID,
		ProjectSFID:                projectSFID,
		AutoEnabled:                aws.BoolValue(input.AutoEnabled),
		BranchProtectionEnabled:    aws.BoolValue(input.BranchProtectionEnabled),
		Version:                    "v1",
	}

	log.WithFields(f).Debug("Encoding github organization record for adding to the database...")
	av, err := dynamodbattribute.MarshalMap(githubOrg)
	if err != nil {
		return nil, err
	}

	log.WithFields(f).Debug("Adding github organization record to the database...")
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(repo.githubOrgTableName),
		ConditionExpression: aws.String("attribute_not_exists(organization_name)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.WithFields(f).Debug("github organization already exists")
				return nil, errors.New("github organization already exists")
			}
		}
		log.WithFields(f).Error("cannot put github organization in dynamodb", err)
		return nil, err
	}

	return toModel(githubOrg), nil
}

func (repo repository) DeleteGithubOrganization(ctx context.Context, externalProjectID string, projectSFID string, githubOrgName string) error {
	f := logrus.Fields{
		"functionName":      "DeleteGithubOrganization",
		utils.XREQUESTID:    ctx.Value(utils.XREQUESTID),
		"externalProjectID": externalProjectID,
		"projectSFID":       projectSFID,
		"githubOrgName":     githubOrgName,
	}

	var githubOrganizationName string
	orgs, orgErr := repo.GetGithubOrganizations(ctx, externalProjectID, projectSFID)
	if orgErr != nil {
		errMsg := fmt.Sprintf("github organization is not found using externalProjectID %s or projectSFID %s error: - %+v", externalProjectID, projectSFID, orgErr)
		log.WithFields(f).Warn(errMsg)
		return errors.New(errMsg)
	}

	for _, githubOrg := range orgs.List {
		if strings.EqualFold(githubOrg.OrganizationName, githubOrgName) {
			githubOrganizationName = githubOrg.OrganizationName
		}
	}

	log.WithFields(f).Debug("Deleting GitHub organization...")
	_, err := repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"organization_name": {
				S: aws.String(githubOrganizationName),
			},
		},
		TableName: aws.String(repo.githubOrgTableName),
	})
	if err != nil {
		errMsg := fmt.Sprintf("error deleting github organization: %s - %+v", githubOrgName, err)
		log.WithFields(f).Warnf(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// UpdateGithubOrganization updates the specified GitHub organization based on the update model provided
func (repo repository) UpdateGithubOrganization(ctx context.Context, projectSFID string, organizationName string, autoEnabled bool, branchProtectionEnabled bool) error {
	f := logrus.Fields{
		"functionName":            "UpdateGithubOrganization",
		utils.XREQUESTID:          ctx.Value(utils.XREQUESTID),
		"projectSFID":             projectSFID,
		"organizationName":        organizationName,
		"autoEnabled":             autoEnabled,
		"branchProtectionEnabled": branchProtectionEnabled,
		"tableName":               repo.githubOrgTableName,
	}

	_, currentTime := utils.CurrentTime()
	githubOrg, lookupErr := repo.GetGithubOrganization(ctx, organizationName)
	if lookupErr != nil {
		log.WithFields(f).Warnf("error looking up github organization by name, error: %+v", lookupErr)
		return lookupErr
	}
	if githubOrg == nil {
		lookupErr := errors.New("unable to lookup github organization by name")
		log.WithFields(f).Warnf("error looking up github organization, error: %+v", lookupErr)
		return lookupErr
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"organization_name": {
				S: aws.String(githubOrg.OrganizationName),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("auto_enabled"),
			"#B": aws.String("branch_protection_enabled"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				BOOL: aws.Bool(autoEnabled),
			},
			":b": {
				BOOL: aws.Bool(branchProtectionEnabled),
			},
			":m": {
				S: aws.String(currentTime),
			},
		},
		UpdateExpression: aws.String("SET #A = :a, #B = :b, #M = :m"),
		TableName:        aws.String(repo.githubOrgTableName),
	}

	log.WithFields(f).Debug("updating github organization record...")
	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.WithFields(f).Warnf("unable to update github organization record, error: %+v", updateErr)
		return updateErr
	}

	return nil
}

func (repo repository) GetGithubOrganizations(ctx context.Context, externalProjectID string, projectSFID string) (*models.GithubOrganizations, error) {
	f := logrus.Fields{
		"functionName":      "GetGithubOrganizations",
		utils.XREQUESTID:    ctx.Value(utils.XREQUESTID),
		"projectSFID":       projectSFID,
		"externalProjectID": externalProjectID,
	}

	var condition expression.KeyConditionBuilder
	var indexName string
	builder := expression.NewBuilder()

	if externalProjectID != "" {
		condition = expression.Key("organization_sfid").Equal(expression.Value(externalProjectID))
		indexName = GithubOrgSFIDIndex
	} else {
		condition = expression.Key("project_sfid").Equal(expression.Value(projectSFID))
		indexName = ProjectSFIDOrganizationNameIndex
	}

	builder = builder.WithKeyCondition(condition)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.githubOrgTableName),
		IndexName:                 aws.String(indexName),
	}

	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.WithFields(f).Warnf("error retrieving github_organizations using organization_sfid = %s, project_sfid = %s. error = %s", externalProjectID, projectSFID, err.Error())
		return nil, err
	}

	if len(results.Items) == 0 {
		log.WithFields(f).Debug("no results from query")
		return &models.GithubOrganizations{
			List: []*models.GithubOrganization{},
		}, nil
	}

	var resultOutput []*GithubOrganization
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &resultOutput)
	if err != nil {
		return nil, err
	}

	log.WithFields(f).Debug("building response model...")
	ghOrgList := buildGithubOrganizationListModels(ctx, resultOutput)
	return &models.GithubOrganizations{List: ghOrgList}, nil
}

func buildGithubOrganizationListModels(ctx context.Context, githubOrganizations []*GithubOrganization) []*models.GithubOrganization {
	f := logrus.Fields{
		"functionName":   "buildGithubOrganizationListModels",
		utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
	}
	ghOrgList := toModels(githubOrganizations)
	if len(ghOrgList) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(ghOrgList))
		for _, ghorganization := range ghOrgList {
			go func(ghorg *models.GithubOrganization) {
				defer wg.Done()
				ghorg.GithubInfo = &models.GithubOrganizationGithubInfo{}
				log.WithFields(f).Debugf("Loading GitHub organization details: %s...", ghorg.OrganizationName)
				user, err := github.GetUserDetails(ghorg.OrganizationName)
				if err != nil {
					ghorg.GithubInfo.Error = err.Error()
				} else {
					url := strfmt.URI(*user.HTMLURL)
					ghorg.GithubInfo.Details = &models.GithubOrganizationGithubInfoDetails{
						Bio:     user.Bio,
						HTMLURL: &url,
						ID:      user.ID,
					}
				}
				ghorg.Repositories = &models.GithubOrganizationRepositories{
					List: make([]*models.GithubRepositoryInfo, 0),
				}
				if ghorg.OrganizationInstallationID != 0 {
					log.WithFields(f).Debugf("Loading GitHub repository list based on installation id: %d...", ghorg.OrganizationInstallationID)
					list, err := github.GetInstallationRepositories(ghorg.OrganizationInstallationID)
					if err != nil {
						log.WithFields(f).Warnf("unable to get repositories for installation id : %d", ghorg.OrganizationInstallationID)
						ghorg.Repositories.Error = err.Error()
						return
					}

					log.WithFields(f).Debugf("Found %d GitHub repositories using installation id: %d...",
						len(list), ghorg.OrganizationInstallationID)
					for _, repoInfo := range list {
						ghorg.Repositories.List = append(ghorg.Repositories.List, &models.GithubRepositoryInfo{
							RepositoryGithubID: utils.Int64Value(repoInfo.ID),
							RepositoryName:     utils.StringValue(repoInfo.FullName),
							RepositoryURL:      utils.StringValue(repoInfo.URL),
							RepositoryType:     "github",
						})
					}
				}
			}(ghorganization)
		}
		wg.Wait()
	}
	return ghOrgList
}

func (repo repository) GetGithubOrganization(ctx context.Context, githubOrganizationName string) (*models.GithubOrganization, error) {
	f := logrus.Fields{
		"functionName":           "GetGithubOrganization",
		utils.XREQUESTID:         ctx.Value(utils.XREQUESTID),
		"githubOrganizationName": "githubOrganizationName",
	}

	log.WithFields(f).Debug("Querying for github organization by name")
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"organization_name": {
				S: aws.String(githubOrganizationName),
			},
		},
		TableName: aws.String(repo.githubOrgTableName),
	})
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		log.WithFields(f).Debug("Unable to find github organization by name")
		return nil, ErrOrganizationDoesNotExist
	}

	var org GithubOrganization
	err = dynamodbattribute.UnmarshalMap(result.Item, &org)
	if err != nil {
		log.WithFields(f).Warnf("error unmarshalling organization table data, error: %v", err)
		return nil, err
	}
	return toModel(&org), nil
}

func (repo repository) GetGithubOrganizationByName(ctx context.Context, githubOrganizationName string) (*models.GithubOrganizations, error) {
	f := logrus.Fields{
		"functionName":           "GetGithubOrganizationByName",
		utils.XREQUESTID:         ctx.Value(utils.XREQUESTID),
		"githubOrganizationName": "githubOrganizationName",
	}

	var condition expression.KeyConditionBuilder
	var indexName string
	builder := expression.NewBuilder()

	condition = expression.Key("organization_name_lower").Equal(expression.Value(strings.ToLower(githubOrganizationName)))
	indexName = GithubOrgLowerNameIndex

	builder = builder.WithKeyCondition(condition)
	// Use the nice builder to create the expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}
	// Assemble the query input parameters
	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(repo.githubOrgTableName),
		IndexName:                 aws.String(indexName),
	}

	log.WithFields(f).Debug("querying for github organization by name...")
	results, err := repo.dynamoDBClient.Query(queryInput)
	if err != nil {
		log.WithFields(f).Warnf("error retrieving github_organizations using githubOrganizationName = %s. error = %s", githubOrganizationName, err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		log.WithFields(f).Debug("unable to find github organization by name")
		return &models.GithubOrganizations{
			List: []*models.GithubOrganization{},
		}, nil
	}
	var resultOutput []*GithubOrganization
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &resultOutput)
	if err != nil {
		log.WithFields(f).Warnf("problem decoding database results, error: %+v", err)
		return nil, err
	}

	ghOrgList := buildGithubOrganizationListModels(ctx, resultOutput)
	return &models.GithubOrganizations{List: ghOrgList}, nil
}
