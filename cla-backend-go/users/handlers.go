// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/users"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {

	// Create user handler
	api.UsersAddUserHandler = users.AddUserHandlerFunc(func(params users.AddUserParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("AddUserHandlerFunc - claUser: %+v", claUser)
		userModel, err := service.CreateUser(&params.Body)
		if err != nil {
			log.Warnf("error creating user from user: %+v, error: %+v", params.Body, err)
			return users.NewAddUserBadRequest().WithPayload(errorResponse(err))
		}

		return users.NewAddUserOK().WithPayload(userModel)
	})

	// Get User by ID handler
	api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("GetUserHandlerFunc - claUser: %+v", claUser)
		userModel, err := service.GetUser(params.UserID)
		if err != nil {
			log.Warnf("error retrieving user for user_id: %s, error: %+v", params.UserID, err)
			return users.NewGetUserBadRequest().WithPayload(errorResponse(err))
		}

		return users.NewGetUserOK().WithPayload(userModel)
	})

	// Get User by name handler
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("GetUserByUserNameHandlerFunc - claUser: %+v", claUser)
		userModel, err := service.GetUserByUserName(params.UserName)
		if err != nil {
			log.Warnf("error retrieving user for user name: '%s', error: %+v", params.UserName, err)
			return users.NewGetUserByUserNameBadRequest().WithPayload(errorResponse(err))
		}

		if userModel == nil {
			log.Warnf("Get User By User Name - '%s' was not found", params.UserName)
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
