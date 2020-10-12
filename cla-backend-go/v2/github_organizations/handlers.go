// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/github_organizations"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service Service, eventService events.Service) {
	api.GithubOrganizationsGetProjectGithubOrganizationsHandler = github_organizations.GetProjectGithubOrganizationsHandlerFunc(
		func(params github_organizations.GetProjectGithubOrganizationsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewGetProjectGithubOrganizationsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Github Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			result, err := service.GetGithubOrganizations(ctx, params.ProjectSFID)
			if err != nil {
				if strings.ContainsAny(err.Error(), "getProjectNotFound") {
					return github_organizations.NewGetProjectGithubOrganizationsNotFound().WithPayload(&models.ErrorResponse{
						Code:       "404",
						Message:    fmt.Sprintf("project not found with given ID. [%s]", params.ProjectSFID),
						XRequestID: reqID,
					})
				}
				return github_organizations.NewGetProjectGithubOrganizationsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			return github_organizations.NewGetProjectGithubOrganizationsOK().WithPayload(result)
		})

	api.GithubOrganizationsAddProjectGithubOrganizationHandler = github_organizations.AddProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.AddProjectGithubOrganizationParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewAddProjectGithubOrganizationForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Add Project GitHub Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			if params.Body.OrganizationName == nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(&models.ErrorResponse{
					Code:       "400",
					Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - missing organization name in body: %+v", params.Body),
					XRequestID: reqID,
				})
			}

			_, err := github.GetOrganization(ctx, *params.Body.OrganizationName)
			if err != nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(errorResponse(reqID, err))
			}

			result, err := service.AddGithubOrganization(ctx, params.ProjectSFID, params.Body)
			if err != nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(errorResponse(reqID, err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				LfUsername:        authUser.UserName,
				EventType:         events.GithubOrganizationAdded,
				ExternalProjectID: params.ProjectSFID,
				EventData: &events.GithubOrganizationAddedEventData{
					GithubOrganizationName: *params.Body.OrganizationName,
				},
			})

			return github_organizations.NewAddProjectGithubOrganizationOK().WithPayload(result)
		})

	api.GithubOrganizationsDeleteProjectGithubOrganizationHandler = github_organizations.DeleteProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.DeleteProjectGithubOrganizationParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewDeleteProjectGithubOrganizationForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete Project GitHub Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			err := service.DeleteGithubOrganization(ctx, params.ProjectSFID, params.OrgName)
			if err != nil {
				if strings.Contains(err.Error(), "getProjectNotFound") {
					return github_organizations.NewDeleteProjectGithubOrganizationNotFound().WithPayload(&models.ErrorResponse{
						Code:       "404",
						Message:    fmt.Sprintf("project not found with given ID. [%s]", params.ProjectSFID),
						XRequestID: reqID,
					})
				}
				return github_organizations.NewDeleteProjectGithubOrganizationBadRequest().WithPayload(errorResponse(reqID, err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				LfUsername:        authUser.UserName,
				EventType:         events.GithubOrganizationDeleted,
				ExternalProjectID: params.ProjectSFID,
				EventData: &events.GithubOrganizationDeletedEventData{
					GithubOrganizationName: params.OrgName,
				},
			})

			return github_organizations.NewDeleteProjectGithubOrganizationNoContent()
		})

	api.GithubOrganizationsUpdateProjectGithubOrganizationConfigHandler = github_organizations.UpdateProjectGithubOrganizationConfigHandlerFunc(
		func(params github_organizations.UpdateProjectGithubOrganizationConfigParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewUpdateProjectGithubOrganizationConfigForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Update Project GitHub Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			if params.Body.AutoEnabled == nil {
				return github_organizations.NewUpdateProjectGithubOrganizationConfigBadRequest().WithPayload(&models.ErrorResponse{
					Code:       "400",
					Message:    "EasyCLA - 400 Bad Request - missing auto enable value in body",
					XRequestID: reqID,
				})
			}

			err := service.UpdateGithubOrganization(ctx, params.ProjectSFID, params.OrgName, *params.Body.AutoEnabled, params.Body.BranchProtectionEnabled)
			if err != nil {
				return github_organizations.NewUpdateProjectGithubOrganizationConfigBadRequest().WithPayload(errorResponse(reqID, err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				LfUsername:        authUser.UserName,
				EventType:         events.GithubOrganizationUpdated,
				ExternalProjectID: params.ProjectSFID,
				EventData: &events.GithubOrganizationUpdatedEventData{
					GithubOrganizationName: params.OrgName,
					AutoEnabled:            *params.Body.AutoEnabled,
				},
			})

			return github_organizations.NewUpdateProjectGithubOrganizationConfigOK()
		})
}

type codedResponse interface {
	Code() string
}

func errorResponse(reqID string, err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:       code,
		Message:    err.Error(),
		XRequestID: reqID,
	}

	return &e
}
