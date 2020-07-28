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

// BuildZipEvent is argument to zipbuilder
type BuildZipEvent struct {
	ClaGroupID    string `json:"cla_group_id"`
	SignatureType string `json:"signature_type"`
}

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

func handler(ctx context.Context, event BuildZipEvent) error {
	var err error
	log.WithField("event", event).Debug("zip builder called")
	switch event.SignatureType {
	case signatures.ICLA:
		err = zipBuilder.BuildICLAZip(event.ClaGroupID)
	case signatures.CCLA:
		err = zipBuilder.BuildCCLAZip(event.ClaGroupID)
	default:
		log.WithField("event", event).Debug("Invalid event")
	}
	if err != nil {
		log.WithField("args", event).Error("failed to build zip", err)
	}
	return err
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
		if len(os.Args) != 3 {
			log.Fatal("invalid number of args. first arg should be icla or ccla and 2nd arg should be cla_group_id")
		}
		err := handler(context.Background(), BuildZipEvent{SignatureType: os.Args[1], ClaGroupID: os.Args[2]})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		lambda.Start(handler)
	}
	log.Infof("Lambda shutting down...")
}
