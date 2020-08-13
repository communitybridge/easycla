// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"fmt"

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
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewGetProjectGithubOrganizationsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Github Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}
			result, err := service.GetGithubOrganizations(params.ProjectSFID)
			if err != nil {
				return github_organizations.NewGetProjectGithubOrganizationsBadRequest().WithPayload(errorResponse(err))
			}
			return github_organizations.NewGetProjectGithubOrganizationsOK().WithPayload(result)
		})
	api.GithubOrganizationsAddProjectGithubOrganizationHandler = github_organizations.AddProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.AddProjectGithubOrganizationParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewAddProjectGithubOrganizationForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Add Project GitHub Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			_, err := github.GetOrganization(params.Body.OrganizationName)
			if err != nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(errorResponse(err))
			}

			result, err := service.AddGithubOrganization(params.ProjectSFID, params.Body)
			if err != nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(errorResponse(err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				LfUsername:        authUser.UserName,
				EventType:         events.GithubOrganizationAdded,
				ExternalProjectID: params.ProjectSFID,
				EventData: &events.GithubOrganizationAddedEventData{
					GithubOrganizationName: params.Body.OrganizationName,
				},
			})

			return github_organizations.NewAddProjectGithubOrganizationOK().WithPayload(result)
		})

	api.GithubOrganizationsDeleteProjectGithubOrganizationHandler = github_organizations.DeleteProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.DeleteProjectGithubOrganizationParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_organizations.NewDeleteProjectGithubOrganizationForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete Project GitHub Organizations with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			err := service.DeleteGithubOrganization(params.ProjectSFID, params.OrgName)
			if err != nil {
				return github_organizations.NewDeleteProjectGithubOrganizationBadRequest().WithPayload(errorResponse(err))
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
