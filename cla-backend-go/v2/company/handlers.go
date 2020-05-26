// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"fmt"

	"github.com/jinzhu/copier"

	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/organization"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

func isUserAuthorizedForOrganization(user *auth.User, externalCompanyID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorizedForOrganizationScope(externalCompanyID) {
			return false
		}
	}
	return true
}

// Configure sets up the middleware handlers
func Configure(api *operations.EasyclaAPI, service Service, v1CompanyRepo v1Company.IRepository, v1CompanyService v1Company.IService) {
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
			if !isUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectContributorsUnauthorized()
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
			if !isUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectClaUnauthorized()
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

	api.CompanyGetCompanyByExternalIDHandler = company.GetCompanyByExternalIDHandlerFunc(func(params company.GetCompanyByExternalIDParams) middleware.Responder {
		companyModel, err := v1CompanyService.GetCompanyByExternalID(params.CompanySFID)
		if err != nil {
			msg := fmt.Sprintf("Bad Request - unable to query company by ExternalID: %s, error: %v", params.CompanySFID, err)
			log.Warnf(msg)
			return company.NewGetCompanyByExternalIDBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}
		result := &models.Company{}
		err = copier.Copy(result, companyModel)
		if err != nil {
			return company.NewGetCompanyByExternalIDInternalServerError().WithPayload(errorResponse(err))
		}
		return company.NewGetCompanyByExternalIDOK().WithPayload(result)
	})
	api.OrganizationSearchOrganizationHandler = organization.SearchOrganizationHandlerFunc(func(params organization.SearchOrganizationParams) middleware.Responder {
		orgs, err := v1CompanyService.SearchOrganizationByName(params.CompanyName)
		if err != nil {
			log.Warnf("error occured while search org %s. error = %s", params.CompanyName, err.Error())
			return organization.NewSearchOrganizationInternalServerError().WithPayload(errorResponse(err))
		}
		result := &models.OrgList{}
		err = copier.Copy(result, orgs)
		if err != nil {
			return organization.NewSearchOrganizationInternalServerError().WithPayload(errorResponse(err))
		}
		return organization.NewSearchOrganizationOK().WithPayload(result)
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
