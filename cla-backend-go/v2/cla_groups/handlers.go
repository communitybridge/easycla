// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_groups

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/foundation"

	"github.com/communitybridge/easycla/cla-backend-go/events"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_group"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure configures the cla group api
func Configure(api *operations.EasyclaAPI, service Service, v1ProjectService v1Project.Service, eventsService events.Service) { //nolint

	api.ClaGroupCreateClaGroupHandler = cla_group.CreateClaGroupHandlerFunc(func(params cla_group.CreateClaGroupParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":        "ClaGroupCreateClaGroupHandler",
			utils.XREQUESTID:      ctx.Value(utils.XREQUESTID),
			"claGroupName":        params.ClaGroupInput.ClaGroupName,
			"claGroupDescription": params.ClaGroupInput.ClaGroupDescription,
			"projectSFIDList":     strings.Join(params.ClaGroupInput.ProjectSfidList, ","),
			"authUsername":        params.XUSERNAME,
			"authEmail":           params.XEMAIL,
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, *params.ClaGroupInput.FoundationSfid) {
			log.WithFields(f).Warnf("user %s does not have access to Create CLA Group with project scope of %s", authUser.UserName, *params.ClaGroupInput.FoundationSfid)
			return cla_group.NewCreateClaGroupForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAGroup with Project scope of %s",
					authUser.UserName, *params.ClaGroupInput.FoundationSfid),
				XRequestID: reqID,
			})
		}

		claGroup, err := service.CreateCLAGroup(ctx, params.ClaGroupInput, utils.StringValue(params.XUSERNAME))
		if err != nil {
			log.WithFields(f).WithError(err).Warn("unable to create the CLA Group")
			return cla_group.NewCreateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "400",
				Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
				XRequestID: reqID,
			})
		}

		// Log the event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.CLAGroupCreated,
			ProjectID:  claGroup.ClaGroupID,
			LfUsername: authUser.UserName,
			EventData:  &events.CLAGroupCreatedEventData{},
		})

		return cla_group.NewCreateClaGroupOK().WithXRequestID(reqID).WithPayload(claGroup)
	})

	api.ClaGroupUpdateClaGroupHandler = cla_group.UpdateClaGroupHandlerFunc(func(params cla_group.UpdateClaGroupParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":   "ClaGroupUpdateClaGroupHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"claGroupID":     params.ClaGroupID,
			"authUsername":   params.XUSERNAME,
			"authEmail":      params.XEMAIL,
		}

		// Make sure we have some parameters to process...
		if params.Body == nil || (params.Body.ClaGroupName == "" && params.Body.ClaGroupDescription == "") {
			log.WithFields(f).Warn("missing CLA Group update parameters - body missing required values")
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(
				utils.ErrorResponseBadRequest(reqID, "missing update parameters - body missing required values"))
		}

		// Attempt to load the existing CLA Group model
		log.WithFields(f).Debug("Loading existing CLA Group by ID...")
		claGroupModel, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			log.WithFields(f).WithError(err).Warn("problem loading CLA group by ID")

			var e *utils.CLAGroupNotFound
			//if err, ok := err.(*utils.CLAGroupNotFound); ok {
			if errors.As(err, &e) {
				log.WithFields(f).WithError(err).Warn("problem loading CLA group by ID - cla group not found")
				return cla_group.NewUpdateClaGroupNotFound().WithXRequestID(reqID).WithPayload(
					utils.ErrorResponseNotFoundWithError(reqID, "CLA Group not found", err))
			}
			if errors.Is(err, v1Project.ErrProjectDoesNotExist) {
				return cla_group.NewUpdateClaGroupNotFound().WithXRequestID(reqID).WithPayload(
					utils.ErrorResponseNotFoundWithError(reqID, "CLA Group not found", err))
			}
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(
				utils.ErrorResponseBadRequestWithError(reqID, fmt.Sprintf("unable to lookup CLA Group by ID: %s", params.ClaGroupID), err))
		}

		// Check permissions now that we can identify the SF Foundation/Project details
		if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
			return cla_group.NewUpdateClaGroupForbidden().WithXRequestID(reqID).WithPayload(
				utils.ErrorResponseForbidden(reqID, fmt.Sprintf("user %s does not have access to UpdateCLAGroup with Project scope of %s", authUser.UserName, claGroupModel.FoundationSFID)))
		}

		// Only update if either the CLA Group Name or Description is changed - if both are the same, abort.
		if claGroupModel.ProjectName == params.Body.ClaGroupName && claGroupModel.ProjectDescription == params.Body.ClaGroupDescription {
			log.WithFields(f).Warn("unable to update the CLA Group Name or Description - provided values are the same as the existing record")
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(
				utils.ErrorResponseBadRequest(reqID, fmt.Sprintf("unable to update the CLA Group Name or Description - values are the same for CLA Group ID: %s", params.ClaGroupID)))
		}

		claGroup, err := service.UpdateCLAGroup(ctx, params.ClaGroupID, params.Body, utils.StringValue(params.XUSERNAME))
		if err != nil {
			log.WithFields(f).WithError(err).Warn("unable to update the CLA Group Name and/or Description - update failed")
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(
				utils.ErrorResponseBadRequestWithError(reqID, fmt.Sprintf("unable to update CLA Group by ID: %s", params.ClaGroupID), err))
		}

		// Log the event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.CLAGroupUpdated,
			ProjectID:  claGroup.ClaGroupID,
			LfUsername: authUser.UserName,
			EventData: &events.CLAGroupUpdatedEventData{
				ClaGroupName:        params.Body.ClaGroupName,
				ClaGroupDescription: params.Body.ClaGroupDescription,
			},
		})

		return cla_group.NewUpdateClaGroupOK().WithXRequestID(reqID).WithPayload(claGroup)
	})

	api.ClaGroupDeleteClaGroupHandler = cla_group.DeleteClaGroupHandlerFunc(func(params cla_group.DeleteClaGroupParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":   "ClaGroupDeleteClaGroupHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"claGroupID":     params.ClaGroupID,
			"authUsername":   params.XUSERNAME,
			"authEmail":      params.XEMAIL,
		}

		claGroupModel, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			log.WithFields(f).Warn(err)
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewDeleteClaGroupNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewDeleteClaGroupNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
					XRequestID: reqID,
				})
			}
			return cla_group.NewDeleteClaGroupInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - unable to lookup CLA Group by ID: %s, error: %+v",
					params.ClaGroupID, err),
				XRequestID: reqID,
			})
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
			return cla_group.NewDeleteClaGroupForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAGroup with Project scope of %s",
					authUser.UserName, claGroupModel.FoundationSFID),
				XRequestID: reqID,
			})
		}

		err = service.DeleteCLAGroup(ctx, claGroupModel, authUser)
		if err != nil {
			log.WithFields(f).Warn(err)
			return cla_group.NewDeleteClaGroupInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error deleting CLA Group %s, error: %+v",
					params.ClaGroupID, err),
				XRequestID: reqID,
			})
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupDeleted,
			ClaGroupModel: claGroupModel,
			LfUsername:    authUser.UserName,
			EventData:     &events.CLAGroupDeletedEventData{},
		})

		return cla_group.NewDeleteClaGroupNoContent().WithXRequestID(reqID)
	})

	api.ClaGroupEnrollProjectsHandler = cla_group.EnrollProjectsHandlerFunc(func(params cla_group.EnrollProjectsParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":    "ClaGroupEnrollProjectsHandler",
			utils.XREQUESTID:  ctx.Value(utils.XREQUESTID),
			"ClaGroupID":      params.ClaGroupID,
			"authUsername":    params.XUSERNAME,
			"authEmail":       params.XEMAIL,
			"projectSFIDList": strings.Join(params.ProjectSFIDList, ","),
		}

		cg, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewEnrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewEnrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
					XRequestID: reqID,
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "400",
				Message:    fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
				XRequestID: reqID,
			})
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, cg.FoundationSFID) {
			log.WithFields(f).Warnf("user %s does not have access with project scope of: %s", authUser.UserName, cg.FoundationSFID)
			return cla_group.NewEnrollProjectsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to enroll with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
				XRequestID: reqID,
			})
		}

		err = service.EnrollProjectsInClaGroup(ctx, params.ClaGroupID, cg.FoundationSFID, params.ProjectSFIDList)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewEnrollProjectsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "400",
					Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "500",
				Message:    fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
				XRequestID: reqID,
			})
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupUpdated,
			ClaGroupModel: cg,
			LfUsername:    authUser.UserName,
			EventData: &events.CLAGroupUpdatedEventData{
				ClaGroupName:        cg.ProjectName,
				ClaGroupDescription: cg.ProjectDescription,
			},
		})

		return cla_group.NewEnrollProjectsOK().WithXRequestID(reqID)
	})

	api.ClaGroupUnenrollProjectsHandler = cla_group.UnenrollProjectsHandlerFunc(func(params cla_group.UnenrollProjectsParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":    "ClaGroupUnenrollProjectsHandler",
			utils.XREQUESTID:  ctx.Value(utils.XREQUESTID),
			"ClaGroupID":      params.ClaGroupID,
			"authUsername":    params.XUSERNAME,
			"authEmail":       params.XEMAIL,
			"projectSFIDList": strings.Join(params.ProjectSFIDList, ","),
		}

		cg, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewUnenrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewUnenrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found", params.ClaGroupID),
					XRequestID: reqID,
				})
			}
			return cla_group.NewUnenrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "400",
				Message:    fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
				XRequestID: reqID,
			})
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, cg.FoundationSFID) {
			log.WithFields(f).Warnf("user %s does not have access with project scope of: %s", authUser.UserName, cg.FoundationSFID)
			return cla_group.NewUnenrollProjectsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to unenroll with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
				XRequestID: reqID,
			})
		}

		err = service.UnenrollProjectsInClaGroup(ctx, params.ClaGroupID, cg.FoundationSFID, params.ProjectSFIDList)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewUnenrollProjectsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "400",
					Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			return cla_group.NewUnenrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       "500",
				Message:    fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
				XRequestID: reqID,
			})
		}

		// TODO: Project Service - remove CLA Enabled flag

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupUpdated,
			ClaGroupModel: cg,
			LfUsername:    authUser.UserName,
			EventData: &events.CLAGroupUpdatedEventData{
				ClaGroupName:        cg.ProjectName,
				ClaGroupDescription: cg.ProjectDescription,
			},
		})

		return cla_group.NewUnenrollProjectsOK().WithXRequestID(reqID)
	})

	api.ClaGroupListClaGroupsUnderFoundationHandler = cla_group.ListClaGroupsUnderFoundationHandlerFunc(func(params cla_group.ListClaGroupsUnderFoundationParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":   "ClaGroupListClaGroupsUnderFoundationHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"projectSFID":    params.ProjectSFID,
			"authUsername":   params.XUSERNAME,
			"authEmail":      params.XEMAIL,
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
			log.WithFields(f).Warnf("user %s does not have access with project scope of: %s", authUser.UserName, params.ProjectSFID)
			return cla_group.NewListClaGroupsUnderFoundationForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to ListCLAGroupsUnderFoundation with Project scope of %s",
					authUser.UserName, params.ProjectSFID),
				XRequestID: reqID,
			})
		}

		result, err := service.ListClaGroupsForFoundationOrProject(ctx, params.ProjectSFID)
		if err != nil {
			if err, ok := err.(*utils.SFProjectNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			if _, ok := err.(*utils.ProjectCLAGroupMappingNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationOK().WithXRequestID(reqID).WithPayload(&models.ClaGroupListSummary{
					List: []*models.ClaGroupSummary{},
				})
			}
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       "404",
					Message:    fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
					XRequestID: reqID,
				})
			}
			return cla_group.NewListClaGroupsUnderFoundationBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 404 Bad Request - error = %s", err.Error()),
			})
		}

		// No results - empty OK response
		if result == nil {
			return cla_group.NewListClaGroupsUnderFoundationOK().WithXRequestID(reqID).WithPayload(&models.ClaGroupListSummary{
				List: []*models.ClaGroupSummary{},
			})
		}

		return cla_group.NewListClaGroupsUnderFoundationOK().WithXRequestID(reqID).WithPayload(result)
	})

	api.ClaGroupValidateClaGroupHandler = cla_group.ValidateClaGroupHandlerFunc(func(params cla_group.ValidateClaGroupParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		// No API user validation - anyone can confirm or use the validate API endpoint

		valid, validationErrors := service.ValidateCLAGroup(ctx, params.ValidationInputRequest)
		return cla_group.NewValidateClaGroupOK().WithXRequestID(reqID).WithPayload(&models.ClaGroupValidationResponse{
			Valid:            valid,
			ValidationErrors: validationErrors,
		})
	})

	api.FoundationListFoundationClaGroupsHandler = foundation.ListFoundationClaGroupsHandlerFunc(func(params foundation.ListFoundationClaGroupsParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		result, err := service.ListAllFoundationClaGroups(ctx, params.FoundationSFID)
		if err != nil {
			return foundation.NewListFoundationClaGroupsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: err.Error(),
			})
		}
		return foundation.NewListFoundationClaGroupsOK().WithXRequestID(reqID).WithPayload(result)
	})
}
