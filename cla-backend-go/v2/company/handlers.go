// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers
func Configure(api *operations.EasyclaAPI, service Service) {
	api.CompanyGetCompanyClaManagersHandler = company.GetCompanyClaManagersHandlerFunc(
		func(params company.GetCompanyClaManagersParams, authUser *auth.User) middleware.Responder {
			result, err := service.GetCompanyCLAManagers(params.CompanyID)
			if err != nil {
				return company.NewGetCompanyClaManagersBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyClaManagersOK().WithPayload(result)
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
