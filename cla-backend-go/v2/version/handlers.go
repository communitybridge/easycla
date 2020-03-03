// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package version

import (
	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/version"

	"github.com/go-openapi/runtime/middleware"
)

// Configure sets the handlers for the API
func Configure(api *operations.EasyclaAPI, Version, Commit, Branch, BuildDate string) {
	api.VersionGetVersionHandler = version.GetVersionHandlerFunc(func(params version.GetVersionParams) middleware.Responder {
		return version.NewGetVersionOK().WithPayload(&v1Models.Version{
			Version:   aws.String(Version),
			Commit:    aws.String(Commit),
			Branch:    aws.String(Branch),
			BuildDate: aws.String(BuildDate),
		})
	})
}
