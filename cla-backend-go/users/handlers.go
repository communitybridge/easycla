// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/gommon/log"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	// Get Project Signatures
	api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams) middleware.Responder {
		userModel, err := service.GetUser(params.UserID)
		if err != nil {
			log.Warnf("error retrieving user for user_id: %s, error: %+v", params.UserID, err)
			return users.NewGetUserBadRequest().WithPayload(errorResponse(err))
		}

		return users.NewGetUserOK().WithPayload(userModel)
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
