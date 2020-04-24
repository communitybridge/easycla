// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service IService, sessionStore *dynastore.Store, signatureService signatures.SignatureService, eventsService events.Service) {

	api.CompanyAddCclaWhitelistRequestHandler = company.AddCclaWhitelistRequestHandlerFunc(
		func(params company.AddCclaWhitelistRequestParams) middleware.Responder {
			requestID, err := service.AddCclaWhitelistRequest(params.CompanyID, params.ProjectID, params.Body)
			if err != nil {
				return company.NewAddCclaWhitelistRequestBadRequest().WithPayload(errorResponse(err))
			}

			eventsService.LogEvent(&events.LogEventArgs{
				EventType: events.CCLAWhitelistRequestCreated,
				ProjectID: params.ProjectID,
				CompanyID: params.CompanyID,
				UserID:    params.Body.ContributorID,
				EventData: &events.CCLAWhitelistRequestCreatedEventData{RequestID: requestID},
			})

			return company.NewAddCclaWhitelistRequestOK()
		})

	api.CompanyDeleteCclaWhitelistRequestHandler = company.DeleteCclaWhitelistRequestHandlerFunc(
		func(params company.DeleteCclaWhitelistRequestParams, claUser *user.CLAUser) middleware.Responder {
			err := service.DeleteCclaWhitelistRequest(params.RequestID)
			if err != nil {
				return company.NewDeleteCclaWhitelistRequestBadRequest().WithPayload(errorResponse(err))
			}

			eventsService.LogEvent(&events.LogEventArgs{
				EventType: events.CCLAWhitelistRequestDeleted,
				ProjectID: params.ProjectID,
				CompanyID: params.CompanyID,
				UserID:    claUser.UserID,
				EventData: &events.CCLAWhitelistRequestDeletedEventData{RequestID: params.RequestID},
			})

			return company.NewDeleteCclaWhitelistRequestOK()
		})

	api.CompanyListCclaWhitelistRequestsHandler = company.ListCclaWhitelistRequestsHandlerFunc(
		func(params company.ListCclaWhitelistRequestsParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.ListCclaWhitelistRequest(params.CompanyID, params.ProjectID)
			if err != nil {
				return company.NewListCclaWhitelistRequestsBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewListCclaWhitelistRequestsOK().WithPayload(result)
		})

	api.CompanyListCclaWhitelistRequestsByCompanyAndProjectHandler = company.ListCclaWhitelistRequestsByCompanyAndProjectHandlerFunc(
		func(params company.ListCclaWhitelistRequestsByCompanyAndProjectParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.ListCclaWhitelistRequest(params.CompanyID, &params.ProjectID)
			if err != nil {
				return company.NewListCclaWhitelistRequestsByCompanyAndProjectBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewListCclaWhitelistRequestsByCompanyAndProjectOK().WithPayload(result)
		})

	api.CompanyListCclaWhitelistRequestsByCompanyAndProjectAndUserHandler = company.ListCclaWhitelistRequestsByCompanyAndProjectAndUserHandlerFunc(
		func(params company.ListCclaWhitelistRequestsByCompanyAndProjectAndUserParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.ListCclaWhitelistRequestByCompanyProjectUser(params.CompanyID, &params.ProjectID, &claUser.LFUsername)
			if err != nil {
				return company.NewListCclaWhitelistRequestsByCompanyAndProjectAndUserBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewListCclaWhitelistRequestsByCompanyAndProjectAndUserOK().WithPayload(result)
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
