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
func Configure(api *operations.EasyclaAPI, service Service, v1CompanyRepo v1Company.IRepository, LFXPortalURL string) { // nolint

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

	api.CompanyGetCompanyCLAGroupManagersHandler = company.GetCompanyCLAGroupManagersHandlerFunc(
		// No auth - invoked from Contributor Console
		func(params company.GetCompanyCLAGroupManagersParams) middleware.Responder {
			result, err := service.GetCompanyCLAGroupManagers(params.CompanyID, params.ClaGroupID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyCLAGroupManagersNotFound().WithPayload(errorResponse(err))
				}
			}

			return company.NewGetCompanyCLAGroupManagersOK().WithPayload(result)
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
			// No permissions needed - anyone can create a company

			// Quick validation of the input parameters
			if !utils.ValidCompanyName(*params.Input.CompanyName) {
				return company.NewCreateCompanyBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - Company Name is not valid",
				})
			}

			companyModel, err := service.CreateCompany(*params.Input.CompanyName, *params.Input.CompanyWebsite, params.UserID, LFXPortalURL)
			if err != nil {
				log.Warnf("error returned from create company api: %+v", err)
				switch err := err; err.(type) {
				case *organizations.CreateOrgConflict:
					return company.NewCreateCompanyConflict().WithPayload(errorResponse(err))
				}
				return company.NewCreateCompanyBadRequest().WithPayload(errorResponse(err))
			}
			return company.NewCreateCompanyOK().WithPayload(companyModel)
		})

	api.CompanyGetCompanyByNameHandler = company.GetCompanyByNameHandlerFunc(
		func(params company.GetCompanyByNameParams, authUser *auth.User) middleware.Responder {
			// Anyone can query for a company by name

			companyModel, err := service.GetCompanyByName(params.CompanyName)
			if err != nil {
				log.Warnf("unable to locate company by name: %s, error: %+v", params.CompanyName, err)
				return company.NewGetCompanyByNameBadRequest().WithPayload(errorResponse(err))
			}

			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by name: %s", params.CompanyName)
				log.Warn(msg)
				return company.NewGetCompanyByNameNotFound().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", msg),
				})
			}

			return company.NewGetCompanyByNameOK().WithPayload(companyModel)
		})

	api.CompanyDeleteCompanyByIDHandler = company.DeleteCompanyByIDHandlerFunc(
		func(params company.DeleteCompanyByIDParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

			// Attempt to locate the company by ID
			companyModel, getErr := service.GetCompanyByID(params.CompanyID)
			if getErr != nil {
				msg := fmt.Sprintf("error returned from get company by ID: %s, error: %+v",
					params.CompanyID, getErr)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// Didn't find the company
			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by ID: %s", params.CompanyID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDNotFound().WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", msg),
				})
			}

			// No external ID assigned - unable to check permissions
			if companyModel.CompanyExternalID == "" {
				msg := fmt.Sprintf("company %s does not have an external SFID assigned - unable to validate permissions, company ID: %s",
					companyModel.CompanyName, params.CompanyID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// finally, we can check permissions for the delete operation
			if !utils.IsUserAuthorizedForOrganization(authUser, companyModel.CompanyExternalID) {
				msg := fmt.Sprintf(" user %s does not have access to company %s with Organization scope of %s",
					authUser.UserName, companyModel.CompanyName, companyModel.CompanyExternalID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDForbidden().WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - %s", msg),
				})
			}

			err := service.DeleteCompanyByID(params.CompanyID)
			if err != nil {
				log.Warnf("unable to delete company by ID: %s, error: %+v", params.CompanyID, err)
				return company.NewDeleteCompanyByIDBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewDeleteCompanyByIDNoContent()
		})

	api.CompanyDeleteCompanyBySFIDHandler = company.DeleteCompanyBySFIDHandlerFunc(
		func(params company.DeleteCompanyBySFIDParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

			// Attempt to locate the company by external SFID
			companyModel, getErr := service.GetCompanyBySFID(params.CompanySFID)
			if getErr != nil {
				msg := fmt.Sprintf("error returned from get company by SFID: %s, error: %+v",
					params.CompanySFID, getErr)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// Didn't find the company
			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by SFID: %s", params.CompanySFID)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDNotFound().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", msg),
				})
			}

			// This should never ever happen given we searched by this key, keep it here anyway
			// No external ID assigned - unable to check permissions
			if companyModel.CompanyExternalID == "" {
				msg := fmt.Sprintf("company %s does not have an external SFID assigned - unable to validate permissions, company ID: %s",
					companyModel.CompanyName, params.CompanySFID)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// finally, we can check permissions for the delete operation
			if !utils.IsUserAuthorizedForOrganization(authUser, companyModel.CompanyExternalID) {
				msg := fmt.Sprintf(" user %s does not have access to company %s with Organization scope of %s",
					authUser.UserName, companyModel.CompanyName, companyModel.CompanyExternalID)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDForbidden().WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - %s", msg),
				})
			}

			err := service.DeleteCompanyBySFID(params.CompanySFID)
			if err != nil {
				log.Warnf("unable to delete company by SFID: %s, error: %+v", params.CompanySFID, err)
				return company.NewDeleteCompanyBySFIDBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewDeleteCompanyBySFIDNoContent()
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
