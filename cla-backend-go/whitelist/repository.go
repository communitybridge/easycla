// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Repository interface defines the functions for the whitelist service
type Repository interface {
	DeleteGithubOrganizationFromWhitelist(claGroupID, githubOrganizationID string) error
	AddGithubOrganizationToWhitelist(claGroupID, githubOrganizationID string) error
	GetGithubOrganizationsFromWhitelist(claGroupID string) ([]models.GithubOrg, error)
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

// NewRepository creates a new instance of the whitelist service
func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// AddGithubOrganizationToWhitelist adds the specified GH organization to the whitelist
func (repo repository) AddGithubOrganizationToWhitelist(CLAGroupID, GithubOrganizationID string) error {
	// get item from dynamodb table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
	})
	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s and GH Org: %s, error: %v", CLAGroupID, GithubOrganizationID, err)
		return err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		itemFromMap = &dynamodb.AttributeValue{}
	}

	// generate new List L without element to be deleted
	// if we find a org with the same id just return without updating
	newList := []*dynamodb.AttributeValue{}
	for _, element := range itemFromMap.L {
		newList = append(newList, element)
		if *element.S == GithubOrganizationID {
			return nil
		}
	}

	// add element to list
	newList = append(newList, &dynamodb.AttributeValue{
		S: aws.String(GithubOrganizationID),
	})

	// update dynamodb table
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
				S: aws.String(CLAGroupID),
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
	}

	_, err = repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating white list, error: %v", err)
		return err
	}

	return nil
}

// DeleteGithubOrganizationFromWhitelist removes the specified GH organization from the whitelist
func (repo repository) DeleteGithubOrganizationFromWhitelist(CLAGroupID, GithubOrganizationID string) error {
	// get item from dynamodb table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s and GH Org: %s, error: %v", CLAGroupID, GithubOrganizationID, err)
		return err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		return errors.New("no github_org_whitelist column")
	}

	// generate new List L without element to be deleted
	newList := []*dynamodb.AttributeValue{}
	for _, element := range itemFromMap.L {
		if *element.S != GithubOrganizationID {
			newList = append(newList, element)
		}
	}

	// update dynamodb table
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
				S: aws.String(CLAGroupID),
			},
		},
		UpdateExpression: aws.String("SET #L = :l"),
	}

	_, err = repo.dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Warnf("Error updating white list, error: %v", err)
		return err
	}

	return nil
}

// GetGithubOrganizationsFromWhitelist returns a list of GH organizations stored in the whitelist
func (repo repository) GetGithubOrganizationsFromWhitelist(CLAGroupID string) ([]models.GithubOrg, error) {
	// get item from dynamodb table
	tableName := fmt.Sprintf("cla-%s-signatures", repo.stage)
	result, err := repo.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"signature_id": {
				S: aws.String(CLAGroupID),
			},
		},
	})

	if err != nil {
		log.Warnf("Error retrieving GH organization whitelist for CLAGroupID: %s, error: %v", CLAGroupID, err)
		return nil, err
	}

	itemFromMap, ok := result.Item["github_org_whitelist"]
	if !ok {
		return nil, nil
	}

	orgs := []models.GithubOrg{}
	for _, org := range itemFromMap.L {
		selected := true
		orgs = append(orgs, models.GithubOrg{
			ID:       org.S,
			Selected: &selected,
		})
	}

	return orgs, nil
}
