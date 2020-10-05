// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_groups

import (
	"context"
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
		if !utils.IsUserAuthorizedForProjectTree(authUser, *params.ClaGroupInput.FoundationSfid) {
			return cla_group.NewCreateClaGroupForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAGroup with Project scope of %s",
					authUser.UserName, *params.ClaGroupInput.FoundationSfid),
			})
		}

		claGroup, err := service.CreateCLAGroup(ctx, params.ClaGroupInput, utils.StringValue(params.XUSERNAME))
		if err != nil {
			return cla_group.NewCreateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
			})
		}

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
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "EasyCLA - 400 Bad Request - missing update parameters - body missing required values",
			})
		}

		claGroupModel, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			log.WithFields(f).Warn(err)
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewUpdateClaGroupNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewUpdateClaGroupNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "400",
				Message: fmt.Sprintf("EasyCLA - 400 Bad Request - unable to lookup CLA Group by ID: %s, error: %+v",
					params.ClaGroupID, err),
			})
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
			return cla_group.NewUpdateClaGroupForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to UpdateCLAGroup with Project scope of %s",
					authUser.UserName, claGroupModel.FoundationSFID),
			})
		}

		// Don't try to update values that are the same... that would be pointless
		if claGroupModel.ProjectName == params.Body.ClaGroupName {
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 400 Bad Request - cannot update CLA Group Name - existing CLA Group name matches updated name value: %s", claGroupModel.ProjectName),
			})
		}
		// Don't try to update values that are the same... that would be pointless
		if claGroupModel.ProjectDescription == params.Body.ClaGroupDescription {
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 400 Bad Request - cannot update CLA Group Description - existing CLA Group description matches updated description value: %s", claGroupModel.ProjectDescription),
			})
		}

		claGroup, err := service.UpdateCLAGroup(ctx, params.ClaGroupID, params.Body, utils.StringValue(params.XUSERNAME))
		if err != nil {
			return cla_group.NewUpdateClaGroupBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
			})
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.CLAGroupUpdated,
			ProjectID:  claGroup.ClaGroupID,
			LfUsername: authUser.UserName,
			EventData:  &events.CLAGroupUpdatedEventData{},
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
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewDeleteClaGroupNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewDeleteClaGroupInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - unable to lookup CLA Group by ID: %s, error: %+v",
					params.ClaGroupID, err),
			})
		}

		if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
			return cla_group.NewDeleteClaGroupForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAGroup with Project scope of %s",
					authUser.UserName, claGroupModel.FoundationSFID),
			})
		}

		err = service.DeleteCLAGroup(ctx, claGroupModel, authUser)
		if err != nil {
			log.WithFields(f).Warn(err)
			return cla_group.NewDeleteClaGroupInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error deleting CLA Group %s, error: %+v",
					params.ClaGroupID, err),
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
		cg, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewEnrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewEnrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		if !utils.IsUserAuthorizedForProjectTree(authUser, cg.FoundationSFID) {
			return cla_group.NewEnrollProjectsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to enroll with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
			})
		}

		err = service.EnrollProjectsInClaGroup(ctx, params.ClaGroupID, cg.FoundationSFID, params.ProjectSFIDList)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewEnrollProjectsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupUpdated,
			ClaGroupModel: cg,
			LfUsername:    authUser.UserName,
			EventData:     &events.CLAGroupUpdatedEventData{},
		})

		return cla_group.NewEnrollProjectsOK().WithXRequestID(reqID)
	})

	api.ClaGroupUnenrollProjectsHandler = cla_group.UnenrollProjectsHandlerFunc(func(params cla_group.UnenrollProjectsParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		cg, err := v1ProjectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if err != nil {
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewUnenrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
				})
			}
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewUnenrollProjectsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewUnenrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		if !utils.IsUserAuthorizedForProjectTree(authUser, cg.FoundationSFID) {
			return cla_group.NewUnenrollProjectsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to unenroll with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
			})
		}

		err = service.UnenrollProjectsInClaGroup(ctx, params.ClaGroupID, cg.FoundationSFID, params.ProjectSFIDList)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewUnenrollProjectsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 Bad Request - %s", err.Error()),
				})
			}
			return cla_group.NewUnenrollProjectsInternalServerError().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}

		// TODO: Project Service - remove CLA Enabled flag

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupUpdated,
			ClaGroupModel: cg,
			LfUsername:    authUser.UserName,
			EventData:     &events.CLAGroupUpdatedEventData{},
		})

		return cla_group.NewUnenrollProjectsOK().WithXRequestID(reqID)
	})

	api.ClaGroupListClaGroupsUnderFoundationHandler = cla_group.ListClaGroupsUnderFoundationHandlerFunc(func(params cla_group.ListClaGroupsUnderFoundationParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
			return cla_group.NewListClaGroupsUnderFoundationForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to ListCLAGroupsUnderFoundation with Project scope of %s",
					authUser.UserName, params.ProjectSFID),
			})
		}

		result, err := service.ListClaGroupsForFoundationOrProject(ctx, params.ProjectSFID)
		if err != nil {
			if err, ok := err.(*utils.SFProjectNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
				})
			}
			if _, ok := err.(*utils.ProjectCLAGroupMappingNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationOK().WithXRequestID(reqID).WithPayload(&models.ClaGroupListSummary{
					List: []*models.ClaGroupSummary{},
				})
			}
			if err, ok := err.(*utils.CLAGroupNotFound); ok {
				return cla_group.NewListClaGroupsUnderFoundationNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - %s", err.Error()),
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
