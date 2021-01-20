// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	v1User "github.com/communitybridge/easycla/cla-backend-go/user"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/go-openapi/runtime/middleware"
)

const (
	//BadRequest error Response code
	BadRequest = "400"
	//Conflict error Response code
	Conflict = "409"
	// NotFound error Response code
	NotFound = "404"
	//Accepted Response code
	Accepted = "202"
)

// Configure is the API handler routine for CLA Manager routes
func Configure(api *operations.EasyclaAPI, service Service, LfxPortalURL string, projectClaGroupRepo projects_cla_groups.Repository, easyCLAUserRepo v1User.RepositoryService) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectOrganizationTree(authUser, params.ProjectSFID, params.CompanySFID, utils.DISALLOW_ADMIN_SCOPE) {
			return cla_manager.NewCreateCLAManagerForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
				XRequestID: reqID,
			})
		}
		cginfo, err := projectClaGroupRepo.GetClaGroupIDForProject(params.ProjectSFID)
		if err != nil {
			if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
				return cla_manager.NewCreateCLAManagerBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       BadRequest,
					Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - No cla group associated with this project: %s", params.ProjectSFID),
					XRequestID: reqID,
				})
			}
			return cla_manager.NewCreateCLAManagerInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "500",
				Message:    fmt.Sprintf("EasyCLA - 500 Internal server error. error = %s", err.Error()),
				XRequestID: reqID,
			})
		}
		compCLAManager, errorResponse := service.CreateCLAManager(ctx, cginfo.ClaGroupID, params, authUser.UserName)
		if errorResponse != nil {
			if errorResponse.Code == BadRequest {
				return cla_manager.NewCreateCLAManagerBadRequest().WithXRequestID(reqID).WithPayload(errorResponse)
			} else if errorResponse.Code == Conflict {
				return cla_manager.NewCreateCLAManagerConflict().WithXRequestID(reqID).WithPayload(errorResponse)
			}
		}

		return cla_manager.NewCreateCLAManagerOK().WithXRequestID(reqID).WithPayload(compCLAManager)
	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectOrganizationTree(authUser, params.ProjectSFID, params.CompanySFID, utils.DISALLOW_ADMIN_SCOPE) {
			return cla_manager.NewDeleteCLAManagerForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAManager with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
				XRequestID: reqID,
			})
		}
		cginfo, err := projectClaGroupRepo.GetClaGroupIDForProject(params.ProjectSFID)
		if err != nil {
			return cla_manager.NewDeleteCLAManagerBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       BadRequest,
				Message:    fmt.Sprintf("EasyCLA - Bad Request. No Cla Group associated with ProjectSFID: %s ", params.ProjectSFID),
				XRequestID: reqID,
			})
		}

		errResponse := service.DeleteCLAManager(ctx, cginfo.ClaGroupID, params)
		if errResponse != nil {
			return cla_manager.NewDeleteCLAManagerBadRequest().WithXRequestID(reqID).WithPayload(errResponse)
		}

		return cla_manager.NewDeleteCLAManagerNoContent().WithXRequestID(reqID)
	})

	api.ClaManagerCreateCLAManagerDesigneeHandler = cla_manager.CreateCLAManagerDesigneeHandlerFunc(func(params cla_manager.CreateCLAManagerDesigneeParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		f := logrus.Fields{
			"functionName":   "ClaManagerCreateCLAManagerDesigneeHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"CompanySFID":    params.CompanySFID,
			"ProjectSFID":    params.ProjectSFID,
			"authUser":       *params.XUSERNAME,
		}

		// Note: anyone create assign a CLA manager designee...no permissions checks
		log.WithFields(f).Debugf("processing CLA Manager Desginee request")
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		claManagerDesignee, err := service.CreateCLAManagerDesignee(ctx, params.CompanySFID, params.ProjectSFID, params.Body.UserEmail.String())
		if err != nil {
			if err == ErrCLAManagerDesigneeConflict {
				msg := fmt.Sprintf("Conflict assigning cla manager role for Project SFID: %s ", params.ProjectSFID)
				return cla_manager.NewCreateCLAManagerDesigneeByGroupConflict().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       Conflict,
						XRequestID: reqID,
					})
			}
			msg := fmt.Sprintf("user :%s, error: %+v ", authUser.Email, err)
			return cla_manager.NewCreateCLAManagerBadRequest().WithXRequestID(reqID).WithPayload(
				&models.ErrorResponse{
					Message:    msg,
					Code:       BadRequest,
					XRequestID: reqID,
				})
		}

		log.Debugf("CLA Manager designee created : %+v", claManagerDesignee)
		return cla_manager.NewCreateCLAManagerDesigneeOK().WithXRequestID(reqID).WithPayload(claManagerDesignee)
	})

	api.ClaManagerCreateCLAManagerDesigneeByGroupHandler = cla_manager.CreateCLAManagerDesigneeByGroupHandlerFunc(
		func(params cla_manager.CreateCLAManagerDesigneeByGroupParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			f := logrus.Fields{
				"functionName":   "ClaManagerCreateCLAManagerDesigneeByGroupHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"CompanySFID":    params.CompanySFID,
				"ClaGroupID":     params.ClaGroupID,
				"Email":          params.Body.UserEmail.String(),
				"authUser":       *params.XUSERNAME,
			}

			// Note: anyone create assign a CLA manager designee...no permissions checks
			log.WithFields(f).Debugf("processing CLA Manager Designee by group request")

			log.WithFields(f).Debugf("getting project IDs for CLA group")
			projectCLAGroups, getErr := projectClaGroupRepo.GetProjectsIdsForClaGroup(params.ClaGroupID)
			if getErr != nil {
				msg := fmt.Sprintf("Error getting SF projects for claGroup: %s ", params.ClaGroupID)
				log.WithFields(f).Warn(msg)
				return cla_manager.NewCreateCLAManagerDesigneeByGroupBadRequest().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       BadRequest,
						XRequestID: reqID,
					})
			}
			log.WithFields(f).Debugf("found %d project IDs for CLA group", len(projectCLAGroups))
			if len(projectCLAGroups) == 0 {
				msg := fmt.Sprintf("no projects associated with CLA Group: %s", params.ClaGroupID)
				log.WithFields(f).Warn(msg)
				return cla_manager.NewCreateCLAManagerDesigneeByGroupNotFound().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       BadRequest,
						XRequestID: reqID,
					})

			}

			designeeScopes, msg, err := service.CreateCLAManagerDesigneeByGroup(ctx, params, projectCLAGroups, f)
			if err != nil {
				if err == ErrCLAManagerDesigneeConflict {
					return cla_manager.NewCreateCLAManagerDesigneeByGroupConflict().WithXRequestID(reqID).WithPayload(
						&models.ErrorResponse{
							Message:    msg,
							Code:       Conflict,
							XRequestID: reqID,
						})
				}
				log.WithFields(f).Warn(msg)
				return cla_manager.NewCreateCLAManagerDesigneeByGroupBadRequest().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       BadRequest,
						XRequestID: reqID,
					})
			}
			return cla_manager.NewCreateCLAManagerDesigneeByGroupOK().WithXRequestID(reqID).WithPayload(&models.ClaManagerDesignees{
				List: designeeScopes,
			})
		})
	api.ClaManagerIsCLAManagerDesigneeHandler = cla_manager.IsCLAManagerDesigneeHandlerFunc(func(params cla_manager.IsCLAManagerDesigneeParams) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

		userRoleStatus, err := service.IsCLAManagerDesignee(ctx, params.CompanySFID, params.ClaGroupID, params.UserLFID)
		if err != nil {
			log.Debugf("Problem checking cla-manager-designee role status for user: %s, error: %+v  ", params.UserLFID, err)
			return cla_manager.NewIsCLAManagerDesigneeBadRequest().WithXRequestID(reqID)
		}

		return cla_manager.NewIsCLAManagerDesigneeOK().WithXRequestID(reqID).WithPayload(userRoleStatus)
	})

	api.ClaManagerInviteCompanyAdminHandler = cla_manager.InviteCompanyAdminHandlerFunc(func(params cla_manager.InviteCompanyAdminParams) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		// Get Contributor details
		user, userErr := easyCLAUserRepo.GetUser(params.UserID)
		if userErr != nil {
			msg := fmt.Sprintf("Problem getting user by ID : %s, error: %+v ", params.UserID, userErr)
			return cla_manager.NewInviteCompanyAdminBadRequest().WithXRequestID(reqID).WithPayload(
				&models.ErrorResponse{
					Code:       BadRequest,
					Message:    msg,
					XRequestID: reqID,
				})
		}

		claManagerDesignees, err := service.InviteCompanyAdmin(ctx, params.Body.ContactAdmin, params.Body.CompanyID, *params.Body.ClaGroupID, params.Body.UserEmail.String(), params.Body.Name, &user, LfxPortalURL)

		if err != nil {
			statusCode := buildErrorStatusCode(err)
			if statusCode == NotFound {
				return cla_manager.NewInviteCompanyAdminNotFound().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    err.Error(),
						Code:       NotFound,
						XRequestID: reqID,
					})
			}
			if statusCode == Conflict {
				msg := fmt.Sprintf("User %s already has role scope assigned ", params.Body.UserEmail)
				return cla_manager.NewInviteCompanyAdminConflict().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       Conflict,
						XRequestID: reqID,
					})
			}
			if statusCode == Accepted {
				msg := fmt.Sprintf("User %s has no LF Login account ", params.Body.UserEmail)
				return cla_manager.NewInviteCompanyAdminAccepted().WithXRequestID(reqID).WithPayload(
					&models.SuccessResponse{
						Message:    msg,
						Code:       Accepted,
						XRequestID: reqID,
					})
			}

			return cla_manager.NewInviteCompanyAdminBadRequest().WithXRequestID(reqID).WithPayload(
				&models.ErrorResponse{
					Message:    err.Error(),
					Code:       BadRequest,
					XRequestID: reqID,
				})
			//return cla_manager.NewInviteCompanyAdminBadRequest().WithPayload(err)
		}

		// successfully created cla manager designee and sent invite
		return cla_manager.NewInviteCompanyAdminOK().WithXRequestID(reqID).WithPayload(&models.ClaManagerDesignees{
			List: claManagerDesignees,
		})
	})

	api.ClaManagerCreateCLAManagerRequestHandler = cla_manager.CreateCLAManagerRequestHandlerFunc(func(params cla_manager.CreateCLAManagerRequestParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		f := logrus.Fields{
			"functionName":   "ClaManagerCreateCLAManagerRequestHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"CompanySFID":    params.CompanySFID,
			"ProjectSFID":    params.ProjectSFID,
			"authUser":       *params.XUSERNAME,
		}
		log.WithFields(f).Debugf("processing CLA Manager request")
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID, utils.ALLOW_ADMIN_SCOPE) {
			return cla_manager.NewCreateCLAManagerRequestForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManagerRequest with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
				XRequestID: reqID,
			})
		}

		claManagerDesignee, err := service.CreateCLAManagerRequest(ctx, params.Body.ContactAdmin, params.CompanySFID, params.ProjectSFID, params.Body.UserEmail.String(),
			*params.Body.FullName, authUser, LfxPortalURL)

		if err != nil {
			statusCode := buildErrorStatusCode(err)
			if statusCode == NotFound {
				return cla_manager.NewCreateCLAManagerRequestNotFound().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    err.Error(),
						Code:       NotFound,
						XRequestID: reqID,
					})
			}
			if statusCode == Conflict {
				msg := fmt.Sprintf("User %s already has role scope assigned ", *params.Body.FullName)
				return cla_manager.NewCreateCLAManagerRequestConflict().WithXRequestID(reqID).WithPayload(
					&models.ErrorResponse{
						Message:    msg,
						Code:       Conflict,
						XRequestID: reqID,
					})
			}

			// Return Bad Request
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithXRequestID(reqID).WithPayload(
				&models.ErrorResponse{
					Message:    err.Error(),
					Code:       BadRequest,
					XRequestID: reqID,
				})
		}
		if params.Body.ContactAdmin {
			return cla_manager.NewCreateCLAManagerRequestNoContent().WithXRequestID(reqID)
		}
		return cla_manager.NewCreateCLAManagerRequestOK().WithXRequestID(reqID).WithPayload(claManagerDesignee)
	})

	api.ClaManagerNotifyCLAManagersHandler = cla_manager.NotifyCLAManagersHandlerFunc(
		func(params cla_manager.NotifyCLAManagersParams) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			err := service.NotifyCLAManagers(ctx, params.Body, LfxPortalURL)
			if err != nil {
				if err == ErrCLAUserNotFound {
					return cla_manager.NewNotifyCLAManagersNotFound().WithXRequestID(reqID).WithPayload(
						utils.ErrorResponseNotFound(
							reqID,
							fmt.Sprintf("unable to notify cla managers - user not found: %s", params.Body.UserID)))
				}
				return cla_manager.NewNotifyCLAManagersBadRequest().WithXRequestID(reqID).WithPayload(
					utils.ErrorResponseBadRequestWithError(reqID, fmt.Sprintf("unable to notify cla managers - cla group: %s, company: %s", params.Body.ClaGroupName, params.Body.CompanyName), err))
			}

			return cla_manager.NewNotifyCLAManagersNoContent().WithXRequestID(reqID)
		})

}

// buildErrorMessageCreate helper function to build an error message
func buildErrorMessageCreate(params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company SFID: %s, project SFID: %s, firstName: %s, lastName: %s, user email: %s, error: %+v",
		params.CompanySFID, params.ProjectSFID, *params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, err)
}

// buildErrorMessage helper function to build an error message
func buildErrorMessageDelete(params cla_manager.DeleteCLAManagerParams, err error) string {
	return fmt.Sprintf("problem deleting new CLA Manager Request using company SFID: %s, project SFID: %s, user ID: %s, error: %+v",
		params.CompanySFID, params.ProjectSFID, params.UserLFID, err)
}

// buildErrorStatusCode helper function to build an error statusCodes
func buildErrorStatusCode(err error) string {
	if err == ErrNoOrgAdmins || err == ErrCLACompanyNotFound || err == ErrClaGroupNotFound || err == ErrCLAUserNotFound {
		return NotFound
	}
	// Check if user is already assigned scope/role
	if err == ErrRoleScopeConflict {
		return Conflict
	}
	// Check if user does not exiss
	if err == ErrNoLFID {
		return Accepted
	}
	// Return Bad Request
	return BadRequest
}
