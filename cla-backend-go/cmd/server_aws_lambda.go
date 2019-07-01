// +build aws_lambda

// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) {
	handler := server()

	lambdaHandler := httpadapter.New(handler)

	lambda.Start(lambdaHandler.Proxy)
}
