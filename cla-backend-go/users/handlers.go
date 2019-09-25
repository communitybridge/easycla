// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/users"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/gommon/log"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	// Get User by ID handler
	api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams, claUser *user.CLAUser) middleware.Responder {
		userModel, err := service.GetUser(params.UserID)
		if err != nil {
			log.Warnf("error retrieving user for user_id: %s, error: %+v", params.UserID, err)
			return users.NewGetUserBadRequest().WithPayload(errorResponse(err))
		}

		return users.NewGetUserOK().WithPayload(userModel)
	})

	// Get User by name handler
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams, claUser *user.CLAUser) middleware.Responder {
		userModel, err := service.GetUserByUserName(params.UserName)
		if err != nil {
			log.Warnf("error retrieving user for user name: %s, error: %+v", params.UserName, err)
			return users.NewGetUserByUserNameBadRequest().WithPayload(errorResponse(err))
		}

		if userModel == nil {
			return users.NewGetUserByUserNameNotFound()
		}

		return users.NewGetUserByUserNameOK().WithPayload(userModel)
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
