// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package onboard

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/onboard"
	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/gommon/log"
)

// Configure sets the response handlers for the onboarding API calls
func Configure(api *operations.ClaAPI, service Service) {
	api.OnboardCreateCLAManagerRequestHandler = onboard.CreateCLAManagerRequestHandlerFunc(
		func(params onboard.CreateCLAManagerRequestParams) middleware.Responder {
			responseModel, err := service.CreateCLAManagerRequest(
				*params.Body.LfID,
				*params.Body.ProjectName,
				*params.Body.CompanyName,
				*params.Body.UserFullName,
				*params.Body.UserEmail)

			if err != nil {
				msg := fmt.Sprintf("Bad Request - unable to create CLA Manager request using LFID: %s, error: %v", *params.Body.LfID, err)
				log.Warnf(msg)
				return onboard.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: msg,
				})
			}

			return onboard.NewCreateCLAManagerRequestOK().WithPayload(responseModel)
		})

	api.OnboardGetCLAManagerRequestsByLFIDHandler = onboard.GetCLAManagerRequestsByLFIDHandlerFunc(
		func(params onboard.GetCLAManagerRequestsByLFIDParams) middleware.Responder {
			responseModels, err := service.GetCLAManagerRequestsByLFID(params.LfID)

			if err != nil {
				msg := fmt.Sprintf("Bad Request - unable to query CLA Manager requests using LFID: %s, error: %v", params.LfID, err)
				log.Warnf(msg)
				return onboard.NewGetCLAManagerRequestsByLFIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: msg,
				})
			}

			if responseModels.Requests == nil || len(responseModels.Requests) == 0 {
				msg := fmt.Sprintf("No requests found for lfid: %s", params.LfID)
				log.Warnf(msg)
				return onboard.NewGetCLAManagerRequestsByLFIDNotFound().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: msg,
				})
			}

			return onboard.NewGetCLAManagerRequestsByLFIDOK().WithPayload(responseModels)
		})

	api.OnboardDeleteCLAManagerRequestsByRequestIDHandler = onboard.DeleteCLAManagerRequestsByRequestIDHandlerFunc(
		func(params onboard.DeleteCLAManagerRequestsByRequestIDParams) middleware.Responder {
			err := service.DeleteCLAManagerRequestsByRequestID(params.RequestID)

			if err != nil {
				msg := fmt.Sprintf("Bad Request - unable to delete CLA Manager requests using request id: %s, error: %v", params.RequestID, err)
				log.Warnf(msg)
				return onboard.NewDeleteCLAManagerRequestsByRequestIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: msg,
				})
			}

			return onboard.NewDeleteCLAManagerRequestsByRequestIDNoContent()
		})
}
