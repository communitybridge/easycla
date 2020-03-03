package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/communitybridge/easycla/cla-backend-go/metrics"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	Insert = "INSERT"
	Modify = "MODIFY"
	Remove = "REMOVE"
)

var awsSession = session.Must(session.NewSession(&aws.Config{}))
var metricsRepo metrics.Repository
var stage string

func init() {
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	metricsRepo = metrics.NewRepository(awsSession, stage)
}

func handler(e events.DynamoDBEvent) error {
	return nil
}

func main() {
	metricsRepo.CalculateAndSaveMetrics()
	os.Exit(0)
	lambda.Start(handler)
}
