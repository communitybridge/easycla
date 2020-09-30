// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

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

// buildCLAGroupDocumentModels builds response models based on the array of db models
func buildCLAGroupDocumentModels(dbDocumentModels []DBProjectDocumentModel) []models.ProjectDocument {
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
