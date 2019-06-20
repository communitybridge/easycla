package company

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service service) {

	api.AddUsertoCompanyAccessListHandler = operations.AddUsertoCompanyAccessListHandlerFunc(func(params operations.AddUsertoCompanyAccessListParams, claUser *user.CLAUser) middleware.Responder {
		err := service.AddUserToCompanyAccessList(params.CompanyID, params.User.InviteID, params.User.UserLFID)
		if err != nil {
			return operations.NewAddGithubOrganizationFromClaBadRequest()
		}

		return operations.NewAddUsertoCompanyAccessListOK()
	})

	api.GetPendingInviteRequestsHandler = operations.GetPendingInviteRequestsHandlerFunc(func(params operations.GetPendingInviteRequestsParams, claUser *user.CLAUser) middleware.Responder {
		result, err := service.GetPendingCompanyInviteRequests(params.CompanyID)
		if err != nil {
			return operations.NewGetPendingInviteRequestsBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewGetPendingInviteRequestsOK().WithPayload(result)
	})

	api.SendInviteRequestHandler = operations.SendInviteRequestHandlerFunc(func(params operations.SendInviteRequestParams, claUser *user.CLAUser) middleware.Responder {

		err := service.SendRequestAccessEmail(params.CompanyID, claUser)
		if err != nil {
			return operations.NewSendInviteRequestBadRequest().WithPayload(errorResponse(err))
		}
		return operations.NewSendInviteRequestOK()
	})

	api.DeletePendingInviteHandler = operations.DeletePendingInviteHandlerFunc(func(params operations.DeletePendingInviteParams, claUser *user.CLAUser) middleware.Responder {
		err := service.DeletePendingCompanyInviteRequest(params.User.InviteID)
		if err != nil {
			return operations.NewDeletePendingInviteBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewDeletePendingInviteOK()
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
