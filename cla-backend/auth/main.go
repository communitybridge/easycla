// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package main

import (
	"auth/authorizer"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	validatorMaker := authorizer.NewValidatorMaker()
	tokenValidator, err := validatorMaker.NewTokenValidator()
	if err != nil {
		return
	}
	usecases := authorizer.NewUsecases(tokenValidator)
	interfaces := authorizer.NewInterfaces(usecases)

	lambda.Start(interfaces.Handler)
}
