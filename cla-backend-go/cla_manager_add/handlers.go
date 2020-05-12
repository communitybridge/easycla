// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager_add

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/cla_manager_add"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure is the API handler routine for the CLA manager routes
func Configure(api *operations.ClaAPI, service IService) {
	api.ClaManagerAddCreateCLAManagerHandler = cla_manager_add.CreateCLAManagerHandlerFunc(func(params cla_manager_add.CreateCLAManagerParams, claUser *user.CLAUser) middleware.Responder {

		signature, addErr := service.AddClaManager(params.CompanyID, params.ProjectID, params.Body.UserLFID)

		if addErr != nil {
			msg := buildErrorMessage(params, addErr)
			log.Warn(msg)
			return cla_manager_add.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		return cla_manager_add.NewCreateCLAManagerOK().WithPayload(signature)

	})

	// Delete CLA Manager
	api.ClaManagerAddDeleteCLAManagerHandler = cla_manager_add.DeleteCLAManagerHandlerFunc(func(params cla_manager_add.DeleteCLAManagerParams, claUser *user.CLAUser) middleware.Responder {

		signature, deleteErr := service.RemoveClaManager(params.CompanyID, params.ProjectID, params.UserLFID)

		if deleteErr != nil {
			msg := buildErrorMessageDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager_add.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		return cla_manager_add.NewDeleteCLAManagerOK().WithPayload(signature)

	})
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(params cla_manager_add.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company ID: %s, project ID: %s, user ID: %s, user name: %s, user email: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.Body.UserLFID, params.Body.UserName, params.Body.UserEmail, err)
}

// buildErrorMessage helper function to build an error message
func buildErrorMessageDelete(params cla_manager_add.DeleteCLAManagerParams, err error) string {
	return fmt.Sprintf("problem deleting new CLA Manager Request using company ID: %s, project ID: %s, user ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.UserLFID, err)
}
