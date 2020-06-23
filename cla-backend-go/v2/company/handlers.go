// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"

	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure sets up the middleware handlers
func Configure(api *operations.EasyclaAPI, service Service, v1CompanyRepo v1Company.IRepository) {
	api.CompanyGetCompanyProjectClaManagersHandler = company.GetCompanyProjectClaManagersHandlerFunc(
		func(params company.GetCompanyProjectClaManagersParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectClaManagersForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Company Project CLA Managers with Organization scope of %s",
						authUser.UserName, params.CompanySFID),
				})
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectClaManagersNotFound()
				}
			}

			result, err := service.GetCompanyProjectCLAManagers(comp.CompanyID, params.ProjectSFID)
			if err != nil {
				return company.NewGetCompanyProjectClaManagersBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectClaManagersOK().WithPayload(result)
		})
	api.CompanyGetCompanyProjectActiveClaHandler = company.GetCompanyProjectActiveClaHandlerFunc(
		func(params company.GetCompanyProjectActiveClaParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectActiveClaForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.CompanySFID),
				})
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectActiveClaNotFound()
				}
			}
			result, err := service.GetCompanyProjectActiveCLAs(comp.CompanyID, params.ProjectSFID)
			if err != nil {
				return company.NewGetCompanyProjectActiveClaBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectActiveClaOK().WithPayload(result)
		})
	api.CompanyGetCompanyProjectContributorsHandler = company.GetCompanyProjectContributorsHandlerFunc(
		func(params company.GetCompanyProjectContributorsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectContributorsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.CompanySFID),
				})
			}
			result, err := service.GetCompanyProjectContributors(params.ProjectSFID, params.CompanySFID, utils.StringValue(params.SearchTerm))
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectContributorsNotFound()
				}
				return company.NewGetCompanyProjectContributorsBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectContributorsOK().WithPayload(result)
		})

	api.CompanyGetCompanyProjectClaHandler = company.GetCompanyProjectClaHandlerFunc(
		func(params company.GetCompanyProjectClaParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectClaForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.CompanySFID),
				})
			}
			result, err := service.GetCompanyProjectCLA(authUser, params.CompanySFID, params.ProjectSFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectClaNotFound().WithPayload(errorResponse(err))
				}
				if err == ErrProjectNotFound {
					return company.NewGetCompanyProjectClaNotFound().WithPayload(errorResponse(err))
				}
				return company.NewGetCompanyProjectClaBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectClaOK().WithPayload(result)
		})

	api.CompanyCreateCompanyHandler = company.CreateCompanyHandlerFunc(
		func(params company.CreateCompanyParams) middleware.Responder {
			// Quick validation of the input parameters
			if !utils.ValidCompanyName(*params.Input.CompanyName) {
				return company.NewCreateCompanyBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - Company Name is not valid",
				})
			}

			companyModel, err := service.CreateCompany(*params.Input.CompanyName, *params.Input.CompanyWebsite, params.UserID)
			if err != nil {
				log.Warnf("error returned from create company api: %+v", err)
				// If EasyCLA company conflict/duplicate or Platform Org Service conflict/duplicate
				if err == ErrDuplicateCompany || err == err.(*organizations.CreateOrgConflict) {
					return company.NewCreateCompanyConflict().WithPayload(errorResponse(err))
				}

				return company.NewCreateCompanyBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewCreateCompanyOK().WithPayload(companyModel)
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
