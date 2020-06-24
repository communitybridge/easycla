// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/communitybridge/easycla/cla-backend-go/v2/dynamo_events"

	"github.com/communitybridge/easycla/cla-backend-go/token"

	"github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/signatures"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/communitybridge/easycla/cla-backend-go/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

var (
	// version the application version
	version string

	// build/Commit the application build number
	commit string

	// branch the build branch
	branch string

	// build date
	buildDate string
)

var dynamoEventsService dynamo_events.Service

func init() {
	var awsSession = session.Must(session.NewSession(&aws.Config{}))
	var stage string
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	log.Infof("STAGE set to %s\n", stage)
	configFile, err := config.LoadConfig("", awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}
	usersRepo := users.NewRepository(awsSession, stage)
	companyRepo := company.NewRepository(awsSession, stage)
	signaturesRepo := signatures.NewRepository(awsSession, stage, companyRepo, usersRepo)
	projectClaGroupRepo := projects_cla_groups.NewRepository(awsSession, stage)
	token.Init(configFile.Auth0Platform.ClientID, configFile.Auth0Platform.ClientSecret, configFile.Auth0Platform.URL, configFile.Auth0Platform.Audience)
	dynamoEventsService = dynamo_events.NewService(stage, signaturesRepo, companyRepo, projectClaGroupRepo)
}

func handler(ctx context.Context, event events.DynamoDBEvent) {
	dynamoEventsService.ProcessEvents(event)
}

func printBuildInfo() {
	log.Infof("Version                 : %s", version)
	log.Infof("Git commit hash         : %s", commit)
	log.Infof("Branch                  : %s", branch)
	log.Infof("Build date              : %s", buildDate)
}

func main() {
	log.Info("Lambda server starting...")
	printBuildInfo()
	if os.Getenv("LOCAL_MODE") == "true" {
		handler(context.Background(), events.DynamoDBEvent{})
	} else {
		lambda.Start(handler)
	}
	log.Infof("Lambda shutting down...")
}
