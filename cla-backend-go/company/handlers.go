// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers
func Configure(api *operations.ClaAPI, service service) {

	api.CompanyAddUsertoCompanyAccessListHandler = company.AddUsertoCompanyAccessListHandlerFunc(func(params company.AddUsertoCompanyAccessListParams, claUser *user.CLAUser) middleware.Responder {
		err := service.AddUserToCompanyAccessList(params.CompanyID, params.User.InviteID, params.User.UserLFID)
		if err != nil {
			return company.NewAddGithubOrganizationFromClaBadRequest()
		}

		return company.NewAddUsertoCompanyAccessListOK()
	})

	api.CompanyGetPendingInviteRequestsHandler = company.GetPendingInviteRequestsHandlerFunc(func(params company.GetPendingInviteRequestsParams, claUser *user.CLAUser) middleware.Responder {
		result, err := service.GetPendingCompanyInviteRequests(params.CompanyID)
		if err != nil {
			return company.NewGetPendingInviteRequestsBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewGetPendingInviteRequestsOK().WithPayload(result)
	})

	api.CompanySendInviteRequestHandler = company.SendInviteRequestHandlerFunc(func(params company.SendInviteRequestParams, claUser *user.CLAUser) middleware.Responder {

		err := service.SendRequestAccessEmail(params.CompanyID, claUser)
		if err != nil {
			return company.NewSendInviteRequestBadRequest().WithPayload(errorResponse(err))
		}
		return company.NewSendInviteRequestOK()
	})

	api.CompanyDeletePendingInviteHandler = company.DeletePendingInviteHandlerFunc(func(params company.DeletePendingInviteParams, claUser *user.CLAUser) middleware.Responder {
		err := service.DeletePendingCompanyInviteRequest(params.User.InviteID)
		if err != nil {
			return company.NewDeletePendingInviteBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewDeletePendingInviteOK()
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
