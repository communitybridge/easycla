// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"

	"github.com/go-openapi/swag"

	"github.com/davecgh/go-spew/spew"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/communitybridge/easycla/cla-backend-go/token"

	"github.com/communitybridge/easycla/cla-backend-go/config"

	"github.com/LF-Engineering/lfx-models/models/stats"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsqs "github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/v2/metrics"
	project_service "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
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
	pcgRepo := projects_cla_groups.NewRepository(awsSession, stage)
	metricsRepo = metrics.NewRepository(awsSession, stage, configFile.APIGatewayURL, pcgRepo)
	token.Init(configFile.Auth0Platform.ClientID, configFile.Auth0Platform.ClientSecret, configFile.Auth0Platform.URL, configFile.Auth0Platform.Audience)
	project_service.InitClient(configFile.APIGatewayURL)
}

func handler(ctx context.Context, event events.CloudWatchEvent) {
	totalCountMetrics, err := metricsRepo.GetTotalCountMetrics()
	if err != nil {
		log.Fatalf("Unable to get totalCount metrics from dynamodb. error = %s", err)
	}

	req := stats.Request{
		Products: &stats.Products{
			EasyCLA: stats.Stats{
				"cla_contributors": stats.Stat{
					Action:    swag.String(stats.StatActionReplace),
					Frequency: swag.String(stats.StatFrequencyAllTime),
					Type:      swag.String(stats.StatTypeNumber),
					Value:     float64(totalCountMetrics.ContributorsCount),
				},
				"clas_signed": stats.Stat{
					Action:    swag.String(stats.StatActionReplace),
					Frequency: swag.String(stats.StatFrequencyAllTime),
					Type:      swag.String(stats.StatTypeNumber),
					Value:     float64(totalCountMetrics.CLAsSignedCount),
				},
				"live_projects": stats.Stat{
					Action:    swag.String(stats.StatActionReplace),
					Frequency: swag.String(stats.StatFrequencyAllTime),
					Type:      swag.String(stats.StatTypeNumber),
					Value:     float64(totalCountMetrics.ProjectsLiveCount),
				},
				"repositories_covered": stats.Stat{
					Action:    swag.String(stats.StatActionReplace),
					Frequency: swag.String(stats.StatFrequencyAllTime),
					Type:      swag.String(stats.StatTypeNumber),
					Value:     float64(totalCountMetrics.GerritRepositoriesEnabledCount + totalCountMetrics.GithubRepositoriesEnabledCount),
				},
			},
		},
		Project: "",
	}

	dataInBytes, err := req.MarshalBinary()
	if err != nil {
		log.Fatalf("marshall the sqs event failed : %v", err)
	}

	log.Debugf("going to send total count metrics to sqs queue %s", spew.Sdump(req))
	conf := config.GetConfig()

	var region, queueURL string
	if conf.MetricsReport.AwsSQSRegion == "" {
		log.Fatalf("aws sqs region config is missing")
	}
	region = conf.MetricsReport.AwsSQSRegion

	if conf.MetricsReport.AwsSQSQueueURL == "" {
		log.Fatalf("aws sqs queue url config is missing")
	}
	queueURL = conf.MetricsReport.AwsSQSQueueURL

	sqsSession := awsqs.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	})))

	input := &awsqs.SendMessageInput{
		MessageBody: aws.String(string(dataInBytes)),
		QueueUrl:    aws.String(queueURL),
	}

	_, err = sqsSession.SendMessage(input)
	if err != nil {
		log.Fatalf("sending the message to sqs failed : %v", err)
	}

	log.Infof("metrics report sent successfully to queue %s", queueURL)
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
		handler(utils.NewContext(), events.CloudWatchEvent{})
	} else {
		lambda.Start(handler)
	}
	log.Infof("Lambda shutting down...")
}
