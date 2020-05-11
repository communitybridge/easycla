// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	v1ClaManager "github.com/communitybridge/easycla/cla-backend-go/cla_manager_add"

	"github.com/LF-Engineering/lfx-kit/auth"

	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
)

// Configure is the API handler routune for CLA Manager routes
func Configure(api *operations.EasyclaAPI, managerService v1ClaManager.IService) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		// Get user by firstname,lastname and email parameters
		userServiceClient := v2UserService.GetClient()
		user, userErr := userServiceClient.SearchUsers(params.Body.FirstName, params.Body.LastName, params.Body.UserEmail)

		if userErr != nil {
			msg := fmt.Sprintf("Failed to get user when searching by firstname : %s, lastname: %s , email: %s , error: %v ",
				params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, userErr)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		signature, addErr := managerService.AddClaManager(params.CompanyID, params.ProjectID, user.Username)

		if addErr != nil {
			msg := buildErrorMessageCreate(params, addErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		v2Signature, err := convertTov2(signature)
		if err != nil {
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse(err))
		}

		return cla_manager.NewCreateCLAManagerOK().WithPayload(v2Signature)
	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		signature, deleteErr := managerService.RemoveClaManager(params.CompanyID, params.ProjectID, params.UserLFID)

		if deleteErr != nil {
			msg := buildErrorMessageDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		v2Signature, err := convertTov2(signature)
		if err != nil {
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse(err))
		}

		return cla_manager.NewDeleteCLAManagerOK().WithPayload(v2Signature)

	})
}

func convertTov2(sig *v1Models.Signature) (*models.Signature, error) {
	var dst models.Signature
	err := copier.Copy(&dst, sig)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// buildErrorMessageCreate helper function to build an error message
func buildErrorMessageCreate(params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company ID: %s, project ID: %s, firstName: %s, lastName: %s, user email: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, err)
}

// buildErrorMessage helper function to build an error message
func buildErrorMessageDelete(params cla_manager.DeleteCLAManagerParams, err error) string {
	return fmt.Sprintf("problem deleting new CLA Manager Request using company ID: %s, project ID: %s, user ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.UserLFID, err)
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
