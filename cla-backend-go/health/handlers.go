// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package health

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service service) {

	api.HealthCheckHandler = operations.HealthCheckHandlerFunc(func(params operations.HealthCheckParams) middleware.Responder {
		result, err := service.HealthCheck(params.HTTPRequest.Context(), params)
		if err != nil {
			return operations.NewHealthCheckBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewHealthCheckOK().WithPayload(result)
	})
}

type codedResponse interface {
	Code() string
}

func errorResponse(err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	return &e
}
