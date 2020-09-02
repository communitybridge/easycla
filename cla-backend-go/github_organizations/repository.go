// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
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
	ProjectSFIDOrganizationNameIndex = "project-sfid-organization-name-index"
)

// errors
var (
	ErrOrganizationDoesNotExist = errors.New("github organization does not exist in cla")
)

// Repository interface defines the functions for the github organizations data model
type Repository interface {
	GetGithubOrganizations(externalProjectID string, projectSFID string) (*models.GithubOrganizations, error)
	AddGithubOrganization(externalProjectID string, projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error)
	DeleteGithubOrganization(externalProjectID string, projectSFID string, githubOrgName string) error
	UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error
	GetGithubOrganization(githubOrganizationName string) (*models.GithubOrganization, error)
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
func (repo repository) AddGithubOrganization(externalProjectID string, projectSFID string, input *models.CreateGithubOrganization) (*models.GithubOrganization, error) {
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
		Version:                    "v1",
	}
	av, err := dynamodbattribute.MarshalMap(githubOrg)
	if err != nil {
		return nil, err
	}
	_, err = repo.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(repo.githubOrgTableName),
		ConditionExpression: aws.String("attribute_not_exists(organization_name)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, errors.New("github organization already exist")
			}
		}
		log.Error("cannot put github organization in dynamodb", err)
		return nil, err
	}
	return toModel(githubOrg), nil
}

func (repo repository) DeleteGithubOrganization(externalProjectID string, projectSFID string, githubOrgName string) error {
	var attrName, attrValue string
	if externalProjectID != "" {
		attrName = "organization_sfid"
		attrValue = externalProjectID
	} else {
		attrName = "project_sfid"
		attrValue = projectSFID
	}
	_, err := repo.dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#id": aws.String(attrName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				S: aws.String(attrValue),
			},
		},
		ConditionExpression: aws.String("#id = :val"),
		Key: map[string]*dynamodb.AttributeValue{
			"organization_name": {S: aws.String(githubOrgName)},
		},
		TableName: aws.String(repo.githubOrgTableName),
	})
	if err != nil {
		errMsg := fmt.Sprintf("error deleting github organization: %s - %+v", githubOrgName, err)
		log.Warnf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// UpdateGithubOrganization updates the specified GitHub organization based on the update model provided
func (repo repository) UpdateGithubOrganization(projectSFID string, organizationName string, autoEnabled bool) error {
	f := logrus.Fields{
		"functionName":     "UpdateGithubOrganization",
		"projectSFID":      projectSFID,
		"organizationName": organizationName,
		"autoEnabled":      autoEnabled,
		"tableName":        repo.githubOrgTableName,
	}

	_, currentTime := utils.CurrentTime()
	githubOrg, lookupErr := repo.GetGithubOrganization(organizationName)
	if lookupErr != nil {
		log.WithFields(f).Warnf("error looking up github organization by name, error: %+v", lookupErr)
		return lookupErr
	}
	if githubOrg == nil {
		lookupErr := errors.New("unable to lookup github organization by name")
		log.WithFields(f).Warnf("error looking up github organization, error: %+v", lookupErr)
		return lookupErr
	}

	log.WithFields(f).Debug("updating github organization record")
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"organization_name": {
				S: aws.String(githubOrg.OrganizationName),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#A": aws.String("auto_enabled"),
			"#M": aws.String("date_modified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				BOOL: aws.Bool(autoEnabled),
			},
			":m": {
				S: aws.String(currentTime),
			},
		},
		UpdateExpression: aws.String("SET #A = :a, #M = :m"),
		TableName:        aws.String(repo.githubOrgTableName),
	}

	_, updateErr := repo.dynamoDBClient.UpdateItem(input)
	if updateErr != nil {
		log.Warnf("unable to update github organization record, error: %+v", updateErr)
		return updateErr
	}

	return nil
}

func (repo repository) GetGithubOrganizations(externalProjectID string, projectSFID string) (*models.GithubOrganizations, error) {
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
		log.Warnf("error retrieving github_organizations using organization_sfid = %s, project_sfid = %s. error = %s", externalProjectID, projectSFID, err.Error())
		return nil, err
	}
	if len(results.Items) == 0 {
		return &models.GithubOrganizations{
			List: []*models.GithubOrganization{},
		}, nil
	}
	var resultOutput []*GithubOrganization
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &resultOutput)
	if err != nil {
		return nil, err
	}
	ghOrgList := buildGithubOrganizationListModels(resultOutput)
	return &models.GithubOrganizations{List: ghOrgList}, nil
}

func buildGithubOrganizationListModels(githubOrganizations []*GithubOrganization) []*models.GithubOrganization {
	f := logrus.Fields{
		"functionName": "buildGithubOrganizationListModels",
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
						log.Warnf("unable to get repositories for installation id : %d", ghorg.OrganizationInstallationID)
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

func (repo repository) GetGithubOrganization(githubOrganizationName string) (*models.GithubOrganization, error) {
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
		return nil, ErrOrganizationDoesNotExist
	}
	var org GithubOrganization
	err = dynamodbattribute.UnmarshalMap(result.Item, &org)
	if err != nil {
		log.Warnf("error unmarshalling organization table data, error: %v", err)
		return nil, err
	}
	return toModel(&org), nil
}
