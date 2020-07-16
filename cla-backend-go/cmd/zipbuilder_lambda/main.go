// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"

	"github.com/communitybridge/easycla/cla-backend-go/v2/signatures"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/communitybridge/easycla/cla-backend-go/config"

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

var zipBuilder signatures.ZipBuilder

func init() {
	var awsSession = session.Must(session.NewSession(&aws.Config{}))
	stage := os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	log.Infof("STAGE set to %s\n", stage)
	configFile, err := config.LoadConfig("", awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}
	zipBuilder = signatures.NewZipBuilder(awsSession, configFile.SignatureFilesBucket)
}

func handler(ctx context.Context) {
	err := zipBuilder.BuildICLAZip("2319870d-0a10-4d97-9b49-1241071e27fa")
	if err != nil {
		log.Error("failed", err)
	}
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
		handler(context.Background())
	} else {
		lambda.Start(handler)
	}
	log.Infof("Lambda shutting down...")
}
