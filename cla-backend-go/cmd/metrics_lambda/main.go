package main

import (
	"context"
	"os"

	"github.com/communitybridge/easycla/cla-backend-go/config"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/metrics"
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

var awsSession = session.Must(session.NewSession(&aws.Config{}))
var metricsRepo metrics.Repository
var stage string

func init() {
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	log.Infof("STAGE set to %s\n", stage)
	configFile, err := config.LoadConfig("", awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}

	metricsRepo = metrics.NewRepository(awsSession, stage, configFile.APIGatewayURL)
}

func handler(ctx context.Context, event events.CloudWatchEvent) {
	err := metricsRepo.CalculateAndSaveMetrics()
	if err != nil {
		log.Fatalf("Unable to save metrics in dynamodb. error = %s", err)
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
	lambda.Start(handler)
	log.Infof("Lambda shutting down...")
}
