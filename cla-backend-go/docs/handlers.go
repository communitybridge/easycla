// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package docs

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/v1/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v1/restapi/operations/docs"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with Service
func Configure(api *operations.ClaAPI) {
	api.DocsGetDocHandler = docs.GetDocHandlerFunc(func(params docs.GetDocParams) middleware.Responder {
		return NewGetDocOK()
	})

	api.DocsGetSwaggerHandler = docs.GetSwaggerHandlerFunc(func(params docs.GetSwaggerParams) middleware.Responder {
		return NewGetSwaggerOK()
	})
}
