// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers
func Configure(api *operations.ClaAPI, service Service, usersService users.Service, companyUserValidation bool) {

	api.CompanyGetCompaniesHandler = company.GetCompaniesHandlerFunc(func(params company.GetCompaniesParams, claUser *user.CLAUser) middleware.Responder {
		companiesModel, err := service.GetCompanies()
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query all companies, error: %v", err)
			log.Warnf(msg)
			return company.NewGetCompaniesBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		return company.NewGetCompaniesOK().WithPayload(companiesModel)
	})

	api.CompanyGetCompanyHandler = company.GetCompanyHandlerFunc(func(params company.GetCompanyParams, claUser *user.CLAUser) middleware.Responder {
		companyModel, err := service.GetCompany(params.CompanyID)
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query company by ID: %s, error: %v", params.CompanyID, err)
			log.Warnf(msg)
			return company.NewGetCompanyBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		if companyModel.CompanyID == "" || companyModel.CompanyName == "" {
			return company.NewGetCompanyNotFound().WithPayload(&models.ErrorResponse{
				Code:    "404",
				Message: fmt.Sprintf("Company Not Found with ID: %s", params.CompanyID),
			})
		}

		return company.NewGetCompanyOK().WithPayload(companyModel)
	})

	api.CompanySearchCompanyHandler = company.SearchCompanyHandlerFunc(func(params company.SearchCompanyParams, claUser *user.CLAUser) middleware.Responder {
		var nextKey = ""
		if params.NextKey != nil {
			nextKey = *params.NextKey
		}

		companiesModel, err := service.SearchCompanyByName(params.CompanyName, nextKey)
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query company by name: %s, error: %v", params.CompanyName, err)
			log.Warnf(msg)
			return company.NewSearchCompanyBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		return company.NewSearchCompanyOK().WithPayload(companiesModel)
	})

	api.CompanyGetCompaniesByUserManagerHandler = company.GetCompaniesByUserManagerHandlerFunc(func(params company.GetCompaniesByUserManagerParams, claUser *user.CLAUser) middleware.Responder {
		if companyUserValidation {
			log.Debugf("Company User Validation - GetUserByUserName() - claUser: %+v", claUser)
			userModel, userErr := usersService.GetUserByUserName(claUser.LFUsername, true)
			if userErr != nil {
				return company.NewGetCompaniesByUserManagerUnauthorized().WithPayload(&models.ErrorResponse{
					Code:    "401",
					Message: fmt.Sprintf("unauthorized - unable to find current logged in user by lf_username: %s", claUser.LFUsername),
				})
			}

			if params.UserID == "" || params.UserID != userModel.UserID {
				return company.NewGetCompaniesByUserManagerUnauthorized().WithPayload(&models.ErrorResponse{
					Code:    "401",
					Message: fmt.Sprintf("unauthorized - userID mismatch: param user id: %s, claUser id: %s", params.UserID, userModel.UserID),
				})
			}
		}

		companies, err := service.GetCompaniesByUserManager(params.UserID)
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query companies by user manager id: %s, error: %v", params.UserID, err)
			log.Warnf(msg)
			return company.NewGetCompaniesByUserManagerBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		return company.NewGetCompaniesByUserManagerOK().WithPayload(companies)
	})

	api.CompanyGetCompaniesByUserManagerWithInvitesHandler = company.GetCompaniesByUserManagerWithInvitesHandlerFunc(func(params company.GetCompaniesByUserManagerWithInvitesParams, claUser *user.CLAUser) middleware.Responder {
		if companyUserValidation {
			log.Debugf("Company User Validation - GetUserByUserName() - claUser: %+v", claUser)
			userModel, userErr := usersService.GetUserByUserName(claUser.LFUsername, true)
			if userErr != nil {
				return company.NewGetCompaniesByUserManagerUnauthorized().WithPayload(&models.ErrorResponse{
					Code:    "401",
					Message: fmt.Sprintf("unauthorized - unable to find current logged in user by lf_username: %s", claUser.LFUsername),
				})
			}

			if params.UserID == "" || params.UserID != userModel.UserID {
				return company.NewGetCompaniesByUserManagerUnauthorized().WithPayload(&models.ErrorResponse{
					Code:    "401",
					Message: fmt.Sprintf("unauthorized - userID mismatch: param user id: %s, claUser id: %s", params.UserID, userModel.UserID),
				})
			}
		}

		companies, err := service.GetCompaniesByUserManagerWithInvites(params.UserID)
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query companies by user manager id: %s, error: %v", params.UserID, err)
			log.Warnf(msg)
			return company.NewGetCompaniesByUserManagerWithInvitesBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		return company.NewGetCompaniesByUserManagerWithInvitesOK().WithPayload(companies)
	})

	api.CompanyAddUsertoCompanyAccessListHandler = company.AddUsertoCompanyAccessListHandlerFunc(func(params company.AddUsertoCompanyAccessListParams, claUser *user.CLAUser) middleware.Responder {
		err := service.AddUserToCompanyAccessList(params.CompanyID, params.User.InviteID, params.User.UserLFID)
		if err != nil {
			log.Warnf("error adding user to company access list using company id: %s, invite id: %s, and user LFID: %s, error: %v",
				params.CompanyID, params.User.InviteID, params.User.UserLFID, err)
			return company.NewAddGithubOrganizationFromClaBadRequest()
		}

		return company.NewAddUsertoCompanyAccessListOK()
	})

	api.CompanyGetCompanyInviteRequestsHandler = company.GetCompanyInviteRequestsHandlerFunc(func(params company.GetCompanyInviteRequestsParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("Processing get company invite request for company ID: %s", params.CompanyID)
		result, err := service.GetCompanyInviteRequests(params.CompanyID)
		if err != nil {
			log.Warnf("error getting company invite using company id: %s, error: %v", params.CompanyID, err)
			return company.NewGetCompanyInviteRequestsBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewGetCompanyInviteRequestsOK().WithPayload(result)
	})

	api.CompanyGetCompanyUserInviteRequestsHandler = company.GetCompanyUserInviteRequestsHandlerFunc(func(params company.GetCompanyUserInviteRequestsParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("Processing get company user invite request for company ID: %s and user ID: %s", params.CompanyID, params.UserID)
		result, err := service.GetCompanyUserInviteRequests(params.CompanyID, params.UserID)
		if err != nil {
			log.Warnf("error getting company user invite using company id: %s, user id: %s, error: %v", params.CompanyID, params.UserID, err)
			return company.NewGetCompanyUserInviteRequestsBadRequest().WithPayload(errorResponse(err))
		}

		if result == nil {
			return company.NewGetCompanyUserInviteRequestsNotFound()
		}

		return company.NewGetCompanyUserInviteRequestsOK().WithPayload(result)
	})

	api.CompanySendInviteRequestHandler = company.SendInviteRequestHandlerFunc(func(params company.SendInviteRequestParams, claUser *user.CLAUser) middleware.Responder {
		log.Debugf("Processing send invite request for company ID: %s, user: %+v", params.CompanyID, claUser)
		if claUser.UserID == "" {
			claUser.UserID = params.Body.UserID
			claUser.LFUsername = params.Body.UserID
			claUser.LFEmail = params.Body.LfEmail
			claUser.Name = params.Body.Username
			log.Debugf("Processing send invite request - added userID for company ID: %s, user: %+v", params.CompanyID, claUser)
		}
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
