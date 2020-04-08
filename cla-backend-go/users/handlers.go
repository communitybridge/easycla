// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/users"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service, eventsService events.Service) {

	// Create user handler
	api.UsersAddUserHandler = users.AddUserHandlerFunc(func(params users.AddUserParams, claUser *user.CLAUser) middleware.Responder {
		exitingModel, getErr := service.GetUserByUserName(claUser.LFUsername, true)
		if getErr != nil {
			msg := fmt.Sprintf("Error querying the user by username, error: %+v", getErr)
			log.Warnf("Create User Failed - %s", msg)
			return users.NewAddUserBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: msg,
			})
		}

		// If the user with the same name exists...
		if exitingModel != nil {
			msg := fmt.Sprintf("User with same username exists: %s", claUser.LFUsername)
			log.Warnf("Create User Failed - %s", msg)
			return users.NewAddUserConflict().WithPayload(&models.ErrorResponse{
				Code:    "409",
				Message: msg,
			})
		}

		newUser := &models.User{
			LfEmail:    claUser.LFEmail,
			LfUsername: claUser.LFUsername,
			Username:   claUser.Name,
		}
		userModel, err := service.CreateUser(newUser)
		if err != nil {
			log.Warnf("error creating user from user: %+v, error: %+v", newUser, err)
			return users.NewAddUserBadRequest().WithPayload(errorResponse(err))
		}
		// filling userID in claUser for logging event
		claUser.UserID = userModel.UserID

		// Create an event - run as a go-routine
		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.UserCreated,
			UserModel: userModel,
			EventData: &events.UserCreatedEventData{},
		})

		return users.NewAddUserOK().WithPayload(userModel)
	})

	// Save/Update User Handler
	api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams, claUser *user.CLAUser) middleware.Responder {
		// Make sure we have good non-empty parameters
		if claUser.LFUsername != params.Body.LfUsername {
			return users.NewUpdateUserUnauthorized().WithPayload(errorResponse(
				fmt.Errorf("user: %s not authorized to update user: %s", claUser.LFUsername, params.Body.LfUsername)))
		}

		userModel, err := service.Save(params.Body)
		if err != nil {
			log.Warnf("error updating user from user request with body: %+v, error: %+v", params.Body, err)
			return users.NewUpdateUserBadRequest().WithPayload(errorResponse(err))
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.UserUpdated,
			UserModel: userModel,
			EventData: &events.UserUpdatedEventData{},
		})

		return users.NewUpdateUserOK().WithPayload(userModel)
	})

	// Delete User Handler
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(func(params users.DeleteUserParams, claUser *user.CLAUser) middleware.Responder {
		/*
				if claUser.UserID == "" {
					return users.NewDeleteUserUnauthorized().WithPayload(errorResponse(
						fmt.Errorf("auth - UsersDeleteUserHandler - user %+v not authorized to delete users - missing UserID", claUser)))
				}


			// Let's lookup the authenticated user in our database - we need to see if they have admin access
			claUserModel, err := service.GetUser(claUser.UserID)
			if err != nil || claUserModel == nil {
				return users.NewUpdateUserUnauthorized().WithPayload(errorResponse(
					fmt.Errorf("error looking up current user permissions to determine if delete is allowed, id: %s, error: %+v",
						params.UserID, err)))
			}

			// Should be an admin to delete
			if !claUserModel.Admin {
				return users.NewUpdateUserUnauthorized().WithPayload(errorResponse(
					fmt.Errorf("user with id: %s is not authorized to delete users - must be admin",
						params.UserID)))
			}
		*/

		err := service.Delete(params.UserID)
		if err != nil {
			log.Warnf("error deleting user from user table with id: %s, error: %+v", params.UserID, err)
			return users.NewUpdateUserBadRequest().WithPayload(errorResponse(err))
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.UserDeleted,
			UserID:    claUser.UserID,
			EventData: &events.UserDeletedEventData{
				DeletedUserID: params.UserID,
			},
		})

		return users.NewDeleteUserNoContent()
	})

	// Get User by ID handler
	api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams, claUser *user.CLAUser) middleware.Responder {
		if claUser.UserID == "" {
			return users.NewGetUserUnauthorized().WithPayload(errorResponse(
				fmt.Errorf("auth - UsersGetUserHandler - user %+v not authorized to get users - missing UserID", claUser)))
		}

		userModel, err := service.GetUser(params.UserID)
		if err != nil {
			log.Warnf("error retrieving user for user_id: %s, error: %+v", params.UserID, err)
			return users.NewGetUserBadRequest().WithPayload(errorResponse(err))
		}

		return users.NewGetUserOK().WithPayload(userModel)
	})

	// Get User by name handler
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams, claUser *user.CLAUser) middleware.Responder {

		userModel, err := service.GetUserByUserName(params.UserName, true)
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

	// Get User by name handler
	api.UsersSearchUsersHandler = users.SearchUsersHandlerFunc(func(params users.SearchUsersParams, claUser *user.CLAUser) middleware.Responder {
		// Make sure we have good non-empty parameters
		if claUser.UserID == "" {
			return users.NewUpdateUserUnauthorized().WithPayload(errorResponse(
				fmt.Errorf("auth - UsersSearchUsersHandler - user %+v not authorized to search users - missing UserID", claUser)))
		}

		// No required params? Return empty result
		if params.SearchField == nil || params.SearchTerm == nil {
			return users.NewSearchUsersBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "Missing searchField and/or searchTerm parameters",
			})
		}

		// Were we passed the full match flag? If so, use it.
		var fullMatch = false
		if params.FullMatch != nil {
			fullMatch = *params.FullMatch
		}

		// Perform the search
		userModel, err := service.SearchUsers(*params.SearchField, *params.SearchTerm, fullMatch)
		if err != nil {
			log.Warnf("error retrieving user for user with name: '%s', error: %+v", *params.SearchTerm, err)
			return users.NewSearchUsersBadRequest().WithPayload(errorResponse(err))
		}

		if userModel == nil {
			log.Warnf("Get User By User Name - '%s' was not found", *params.SearchTerm)
			return users.NewSearchUsersNotFound()
		}

		return users.NewSearchUsersOK().WithPayload(userModel)

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
