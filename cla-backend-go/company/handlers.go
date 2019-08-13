// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/labstack/gommon/log"

	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers
func Configure(api *operations.ClaAPI, service service) {

	api.CompanyAddUsertoCompanyAccessListHandler = company.AddUsertoCompanyAccessListHandlerFunc(func(params company.AddUsertoCompanyAccessListParams, claUser *user.CLAUser) middleware.Responder {
		err := service.AddUserToCompanyAccessList(params.CompanyID, params.User.InviteID, params.User.UserLFID)
		if err != nil {
			log.Warnf("error adding user to company access list using company id: %s, invite id: %s, and user LFID: %s, error: %v",
				params.CompanyID, params.User.InviteID, params.User.UserLFID, err)
			return company.NewAddGithubOrganizationFromClaBadRequest()
		}

		return company.NewAddUsertoCompanyAccessListOK()
	})

	api.CompanyGetPendingInviteRequestsHandler = company.GetPendingInviteRequestsHandlerFunc(func(params company.GetPendingInviteRequestsParams, claUser *user.CLAUser) middleware.Responder {
		result, err := service.GetPendingCompanyInviteRequests(params.CompanyID)
		if err != nil {
			log.Warnf("error getting pending company invite using company id: %s, error: %v", params.CompanyID, err)
			return company.NewGetPendingInviteRequestsBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewGetPendingInviteRequestsOK().WithPayload(result)
	})

	api.CompanySendInviteRequestHandler = company.SendInviteRequestHandlerFunc(func(params company.SendInviteRequestParams, claUser *user.CLAUser) middleware.Responder {

		err := service.SendRequestAccessEmail(params.CompanyID, claUser)
		if err != nil {
			log.Warnf("error sending request access email using company id: %s with user: %v, error: %v", params.CompanyID, claUser, err)
			return company.NewSendInviteRequestBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewSendInviteRequestOK()
	})

	api.CompanyDeletePendingInviteHandler = company.DeletePendingInviteHandlerFunc(func(params company.DeletePendingInviteParams, claUser *user.CLAUser) middleware.Responder {
		err := service.DeletePendingCompanyInviteRequest(params.User.InviteID)
		if err != nil {
			log.Warnf("error deleting pending company invite using id: %s, error: %v", params.User.InviteID, err)
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
