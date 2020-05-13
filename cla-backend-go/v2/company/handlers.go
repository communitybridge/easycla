// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

func isUserAuthorizedForOrganization(user *auth.User, externalCompanyID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorized(auth.Organization, externalCompanyID) {
			return false
		}
	}
	return true
}

func isUserAuthorizedForProjectOrganization(user *auth.User, externalProjectID, externalCompanyID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorizedByProject(externalProjectID, externalCompanyID) {
			return false
		}
	}
	return true
}

// Configure sets up the middleware handlers
func Configure(api *operations.EasyclaAPI, service Service, v1CompanyRepo v1Company.IRepository) {
	api.CompanyGetCompanyClaManagersHandler = company.GetCompanyClaManagersHandlerFunc(
		func(params company.GetCompanyClaManagersParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !isUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyClaManagersUnauthorized()
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyClaManagersNotFound()
				}
			}

			result, err := service.GetCompanyCLAManagers(comp.CompanyID)
			if err != nil {
				return company.NewGetCompanyClaManagersBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyClaManagersOK().WithPayload(result)
		})
	api.CompanyGetCompanyActiveClaHandler = company.GetCompanyActiveClaHandlerFunc(
		func(params company.GetCompanyActiveClaParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !isUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyActiveClaUnauthorized()
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyActiveClaNotFound()
				}
			}
			result, err := service.GetCompanyActiveCLAs(comp.CompanyID)
			if err != nil {
				return company.NewGetCompanyActiveClaBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyActiveClaOK().WithPayload(result)
		})
	api.CompanyGetCompanyProjectContributorsHandler = company.GetCompanyProjectContributorsHandlerFunc(
		func(params company.GetCompanyProjectContributorsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
				return company.NewGetCompanyProjectContributorsUnauthorized()
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectContributorsNotFound()
				}
			}
			log.WithField("company_id", comp.CompanyID).Debugf("searching corporate contributors for company")
			result, err := service.GetCompanyProjectContributors(params.ProjectSFID, comp.CompanyID, utils.StringValue(params.SearchTerm))
			if err != nil {
				return company.NewGetCompanyProjectContributorsBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectContributorsOK().WithPayload(result)
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
