// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/sirupsen/logrus"

	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"
	"github.com/go-openapi/runtime/middleware"
)

const (
	//BadRequest error Response code
	BadRequest = "400"
	//Conflict error Response code
	Conflict = "409"
	// NotFound error Response code
	NotFound = "404"
)

// Configure sets up the middleware handlers
func Configure(api *operations.EasyclaAPI, service Service, v1CompanyRepo v1Company.IRepository, projectClaGroupRepo projects_cla_groups.Repository, LFXPortalURL string) { // nolint

	api.CompanyGetCompanyProjectClaManagersHandler = company.GetCompanyProjectClaManagersHandlerFunc(
		func(params company.GetCompanyProjectClaManagersParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectClaManagersForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Company Project CLA Managers with Organization scope of %s",
						authUser.UserName, params.CompanySFID),
				})
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(ctx, params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectClaManagersNotFound().WithXRequestID(reqID)
				}
			}

			result, err := service.GetCompanyProjectCLAManagers(ctx, comp.CompanyID, params.ProjectSFID)
			if err != nil {
				return company.NewGetCompanyProjectClaManagersBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectClaManagersOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.CompanyGetCompanyCLAGroupManagersHandler = company.GetCompanyCLAGroupManagersHandlerFunc(
		// No auth - invoked from Contributor Console
		func(params company.GetCompanyCLAGroupManagersParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			result, err := service.GetCompanyCLAGroupManagers(ctx, params.CompanyID, params.ClaGroupID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyCLAGroupManagersNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
			}

			return company.NewGetCompanyCLAGroupManagersOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.CompanyGetCompanyProjectActiveClaHandler = company.GetCompanyProjectActiveClaHandlerFunc(
		func(params company.GetCompanyProjectActiveClaParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectActiveClaForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.CompanySFID),
				})
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(ctx, params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectActiveClaNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:    "404",
						Message: fmt.Sprintf("Company not found with given ID. [%s]", params.CompanySFID),
					})
				}
			}
			result, err := service.GetCompanyProjectActiveCLAs(ctx, comp.CompanyID, params.ProjectSFID)
			if err != nil {
				if strings.ContainsAny(err.Error(), "getProjectNotFound") {
					return company.NewGetCompanyProjectActiveClaNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:    "404",
						Message: fmt.Sprintf("clagroup not found with given ID. [%s]", params.ProjectSFID),
					})
				}
				return company.NewGetCompanyProjectActiveClaBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectActiveClaOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.CompanyGetCompanyProjectContributorsHandler = company.GetCompanyProjectContributorsHandlerFunc(
		func(params company.GetCompanyProjectContributorsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

			// PM - check if authorized by project scope - allow if PM has project ID scope that matches
			// CLA Manager - check if authorized by project|organization scope - allow if CLA Manager (for example) has project ID + org DI scope that matches
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) && !utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
				return company.NewGetCompanyProjectContributorsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to get contributors with Project scope of %s or Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.ProjectSFID, params.CompanySFID),
				})
			}

			result, err := service.GetCompanyProjectContributors(ctx, params.ProjectSFID, params.CompanySFID, utils.StringValue(params.SearchTerm))
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectContributorsNotFound().WithXRequestID(reqID)
				}
				return company.NewGetCompanyProjectContributorsBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectContributorsOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.CompanyGetCompanyProjectClaHandler = company.GetCompanyProjectClaHandlerFunc(
		func(params company.GetCompanyProjectClaParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return company.NewGetCompanyProjectClaForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
						authUser.UserName, params.ProjectSFID, params.CompanySFID),
				})
			}
			result, err := service.GetCompanyProjectCLA(ctx, authUser, params.CompanySFID, params.ProjectSFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyProjectClaNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				if err == ErrProjectNotFound {
					return company.NewGetCompanyProjectClaNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return company.NewGetCompanyProjectClaBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return company.NewGetCompanyProjectClaOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.CompanyCreateCompanyHandler = company.CreateCompanyHandlerFunc(
		func(params company.CreateCompanyParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			// No permissions needed - anyone can create a company

			// Quick validation of the input parameters
			if !utils.ValidCompanyName(*params.Input.CompanyName) {
				return company.NewCreateCompanyBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - Company Name is not valid",
				})
			}

			companyModel, err := service.CreateCompany(ctx, *params.Input.CompanyName, *params.Input.CompanyWebsite, params.Input.UserEmail.String(), params.UserID)
			if err != nil {
				log.Warnf("error returned from create company api: %+v", err)
				if strings.Contains(err.Error(), "website already exists") {
					formatErr := errors.New("website already exists")
					return company.NewCreateCompanyConflict().WithXRequestID(reqID).WithPayload(errorResponse(formatErr))
				}
				if _, ok := err.(*organizations.CreateOrgConflict); ok {
					formatErr := errors.New("organization already exists")
					return company.NewCreateCompanyConflict().WithXRequestID(reqID).WithPayload(errorResponse(formatErr))
				}
				return company.NewCreateCompanyBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return company.NewCreateCompanyOK().WithXRequestID(reqID).WithPayload(companyModel)
		})

	api.CompanyGetCompanyByNameHandler = company.GetCompanyByNameHandlerFunc(
		func(params company.GetCompanyByNameParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			// Anyone can query for a company by name

			companyModel, err := service.GetCompanyByName(ctx, params.CompanyName)
			if err != nil {
				log.Warnf("unable to locate company by name: %s, error: %+v", params.CompanyName, err)
				return company.NewGetCompanyByNameBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by name: %s", params.CompanyName)
				log.Warn(msg)
				return company.NewGetCompanyByNameNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", msg),
				})
			}

			return company.NewGetCompanyByNameOK().WithXRequestID(reqID).WithPayload(companyModel)
		})

	api.CompanyDeleteCompanyByIDHandler = company.DeleteCompanyByIDHandlerFunc(
		func(params company.DeleteCompanyByIDParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

			// Attempt to locate the company by ID
			companyModel, getErr := service.GetCompanyByID(ctx, params.CompanyID)
			if getErr != nil {
				msg := fmt.Sprintf("error returned from get company by ID: %s, error: %+v",
					params.CompanyID, getErr)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// Didn't find the company
			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by ID: %s", params.CompanyID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", msg),
				})
			}

			// No external ID assigned - unable to check permissions
			if companyModel.CompanyExternalID == "" {
				msg := fmt.Sprintf("company %s does not have an external SFID assigned - unable to validate permissions, company ID: %s",
					companyModel.CompanyName, params.CompanyID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// finally, we can check permissions for the delete operation
			if !utils.IsUserAuthorizedForOrganization(authUser, companyModel.CompanyExternalID) {
				msg := fmt.Sprintf(" user %s does not have access to company %s with Organization scope of %s",
					authUser.UserName, companyModel.CompanyName, companyModel.CompanyExternalID)
				log.Warn(msg)
				return company.NewDeleteCompanyByIDForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - %s", msg),
				})
			}

			err := service.DeleteCompanyByID(ctx, params.CompanyID)
			if err != nil {
				log.Warnf("unable to delete company by ID: %s, error: %+v", params.CompanyID, err)
				return company.NewDeleteCompanyByIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			return company.NewDeleteCompanyByIDNoContent().WithXRequestID(reqID)
		})

	api.CompanyDeleteCompanyBySFIDHandler = company.DeleteCompanyBySFIDHandlerFunc(
		func(params company.DeleteCompanyBySFIDParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

			// Attempt to locate the company by external SFID
			companyModel, getErr := service.GetCompanyBySFID(ctx, params.CompanySFID)
			if getErr != nil {
				msg := fmt.Sprintf("error returned from get company by SFID: %s, error: %+v",
					params.CompanySFID, getErr)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// Didn't find the company
			if companyModel == nil {
				msg := fmt.Sprintf("unable to locate company by SFID: %s", params.CompanySFID)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
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
				return company.NewDeleteCompanyBySFIDBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", msg),
				})
			}

			// finally, we can check permissions for the delete operation
			if !utils.IsUserAuthorizedForOrganization(authUser, companyModel.CompanyExternalID) {
				msg := fmt.Sprintf(" user %s does not have access to company %s with Organization scope of %s",
					authUser.UserName, companyModel.CompanyName, companyModel.CompanyExternalID)
				log.Warn(msg)
				return company.NewDeleteCompanyBySFIDForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - %s", msg),
				})
			}

			err := service.DeleteCompanyBySFID(ctx, params.CompanySFID)
			if err != nil {
				log.Warnf("unable to delete company by SFID: %s, error: %+v", params.CompanySFID, err)
				return company.NewDeleteCompanyBySFIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			return company.NewDeleteCompanyBySFIDNoContent().WithXRequestID(reqID)
		})

	api.CompanyContributorAssociationHandler = company.ContributorAssociationHandlerFunc(
		func(params company.ContributorAssociationParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			contributor, contributorErr := service.AssociateContributor(ctx, params.CompanySFID, params.Body.UserEmail.String())
			if contributorErr != nil {
				if _, ok := contributorErr.(*organizations.CreateOrgUsrRoleScopesConflict); ok {
					formatErr := errors.New("user already assigned contributor role for company")
					return company.NewContributorAssociationConflict().WithXRequestID(reqID).WithPayload(errorResponse(formatErr))
				}
				return company.NewContributorAssociationBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(contributorErr))
			}
			return company.NewContributorAssociationOK().WithXRequestID(reqID).WithPayload(contributor)
		})

	api.CompanyGetCompanyAdminsHandler = company.GetCompanyAdminsHandlerFunc(
		func(params company.GetCompanyAdminsParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			adminList, adminErr := service.GetCompanyAdmins(ctx, params.CompanySFID)
			if adminErr != nil {
				return company.NewGetCompanyAdminsBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(adminErr))
			}
			return company.NewGetCompanyAdminsOK().WithXRequestID(reqID).WithPayload(adminList)
		})

	api.CompanyContributorRoleScopAssociationHandler = company.ContributorRoleScopAssociationHandlerFunc(
		func(params company.ContributorRoleScopAssociationParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			f := logrus.Fields{
				"functionName":   "CompanyContributorRoleScopAssociationHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"CompanySFID":    params.CompanySFID,
				"ClaGroupID":     params.ClaGroupID,
				"Email":          params.Body.UserEmail.String(),
			}
			log.WithFields(f).Debugf("processing CLA Manager Desginee by group request")

			log.WithFields(f).Debugf("getting project IDs for CLA group")
			projectCLAGroups, getErr := projectClaGroupRepo.GetProjectsIdsForClaGroup(params.ClaGroupID)
			if getErr != nil {
				msg := fmt.Sprintf("Error getting SF projects for claGroup: %s ", params.ClaGroupID)
				log.WithFields(f).Warn(msg)
				return company.NewContributorRoleScopAssociationBadRequest().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    BadRequest,
					})
			}
			log.WithFields(f).Debugf("found %d project IDs for CLA group", len(projectCLAGroups))
			if len(projectCLAGroups) == 0 {
				msg := fmt.Sprintf("no projects associated with CLA Group: %s", params.ClaGroupID)
				log.WithFields(f).Warn(msg)
				return company.NewContributorRoleScopAssociationNotFound().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    BadRequest,
					})

			}

			contributor, msg, err := service.AssociateContributorByGroup(ctx, params.CompanySFID, params.Body.UserEmail.String(), projectCLAGroups, params.ClaGroupID)
			if err != nil {
				if err == ErrContributorConflict {
					return company.NewContributorRoleScopAssociationConflict().WithXRequestID(reqID).WithPayload(
						&models.ErrorResponse{
							Message: msg,
							Code:    Conflict,
						})
				}
				log.WithFields(f).Warn(msg)
				return company.NewContributorRoleScopAssociationBadRequest().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    BadRequest,
					})
			}
			return company.NewContributorRoleScopAssociationOK().WithXRequestID(reqID).WithPayload(&models.Contributors{
				List: contributor,
			})
		})

	api.CompanyAssignCompanyOwnerHandler = company.AssignCompanyOwnerHandlerFunc(
		func(params company.AssignCompanyOwnerParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			f := logrus.Fields{
				"functionName": "CompanyCompanyAssignCompanyOwnerHandler",
				"CompanySFID":  params.CompanySFID,
				"Email":        params.Body.UserEmail.String(),
			}
			log.WithFields(f).Debugf("processing Assigning Company owner role to user")
			companyOwner, ownerErr := service.AssignCompanyOwner(ctx, params.CompanySFID, params.Body.UserEmail.String(), LFXPortalURL)
			if ownerErr != nil {
				if _, ok := ownerErr.(*organizations.ListOrgUsrAdminScopesNotFound); !ok {
					log.WithFields(f).Debugf("Problem assigning company owner , error: %+v ", ownerErr)
					return company.NewAssignCompanyOwnerBadRequest().WithXRequestID(reqID)
				}
			}
			log.WithFields(f).Debugf("processed Assigning Company owner role to user")
			return company.NewAssignCompanyOwnerOK().WithXRequestID(reqID).WithPayload(companyOwner)
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
