// +build aws_lambda

// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type fn func(ctx context.Context, sqsEvent events.SQSEvent) error

func Start(handler fn) error {
	log.Info("starting lambda handler...")

	lambda.Start(handler)
	return nil
}
