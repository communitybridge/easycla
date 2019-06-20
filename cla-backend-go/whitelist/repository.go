// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package whitelist

import (
	"errors"
	"fmt"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Repository interface {
	DeleteGithubOrganizationFromWhitelist(claGroupID, githubOrganizationID string) error
	AddGithubOrganizationToWhitelist(claGroupID, githubOrganizationID string) error
	GetGithubOrganizationsFromWhitelist(claGroupID string) ([]models.GithubOrg, error)
}

type repository struct {
	stage          string
	dynamoDBClient *dynamodb.DynamoDB
}

func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

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
		fmt.Println(err.Error())
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
		fmt.Println("Error updating white list : ", err)
		return err
	}

	return nil
}

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
		fmt.Println(err.Error())
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
		fmt.Println("Error updating white list : ", err)
		return err
	}

	return nil
}

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
		fmt.Println(err.Error())
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
