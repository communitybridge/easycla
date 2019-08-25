// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package health

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/labstack/gommon/log"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/health"
	ini "github.com/communitybridge/easycla/cla-backend-go/init"
)

// Service provides an API to the health API
type Service struct {
	version   string
	commit    string
	branch    string
	buildDate string
}

// New is a simple helper function to create a health service instance
func New(version, commit, branch, buildDate string) Service {
	return Service{
		version:   version,
		commit:    commit,
		branch:    branch,
		buildDate: buildDate,
	}
}

// HealthCheck API call returns the current health of the service
func (s Service) HealthCheck(ctx context.Context, in health.HealthCheckParams) (*models.Health, error) {

	t := time.Now()
	duration := time.Since(t)
	hs := models.HealthStatus{TimeStamp: time.Now().UTC().Format(time.RFC3339), Healthy: true, Name: "CLA", Duration: duration.String()}

	// Do a quick check to see if we have a database connection
	dynamoNow := time.Now()
	dynamoAlive := isDynamoAlive()
	dynamoDuration := time.Since(dynamoNow)
	dy := models.HealthStatus{TimeStamp: time.Now().UTC().Format(time.RFC3339), Healthy: dynamoAlive, Name: "CLA - Dynamodb", Duration: dynamoDuration.String()}

	var status = "healthy"
	if !dynamoAlive {
		status = "not healthy"
	}

	response := models.Health{
		Status:         status,
		TimeStamp:      time.Now().UTC().Format(time.RFC3339),
		Version:        s.version,
		Githash:        s.commit,
		Branch:         s.branch,
		BuildTimeStamp: s.buildDate,
		Healths: []*models.HealthStatus{
			&hs,
			&dy,
		},
	}

	return &response, nil
}

// isDynamoAlive runs a check to see if we have connectivity to the database - returns true if successful, false otherwise
func isDynamoAlive() bool {
	// Grab the AWS session
	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Warnf("Unable to acquire AWS session - returning failed health, error: %v", err)
		return false
	}

	// Known table that we can query
	tableName := "cla-" + ini.GetStage() + "-projects"

	// Create a client and make a query - don't wory about the result - just check the error response
	dynamoDBClient := dynamodb.New(awsSession)
	_, err = dynamoDBClient.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &tableName,
	})

	// No error is success
	return err == nil
}
