// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github_organizations

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/github_organizations"
	v1GithubOrganizations "github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1GithubOrganizations.Service) {
	api.GithubOrganizationsGetProjectGithubOrganizationsHandler = github_organizations.GetProjectGithubOrganizationsHandlerFunc(
		func(params github_organizations.GetProjectGithubOrganizationsParams, claUser *auth.User) middleware.Responder {
			result, err := service.GetGithubOrganizations(params.ProjectSFID)
			if err != nil {
				return github_organizations.NewGetProjectGithubOrganizationsBadRequest().WithPayload(errorResponse(err))
			}
			return github_organizations.NewGetProjectGithubOrganizationsOK().WithPayload(*result)
		})
	api.GithubOrganizationsAddProjectGithubOrganizationHandler = github_organizations.AddProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.AddProjectGithubOrganizationParams, claUser *auth.User) middleware.Responder {
			result, err := service.AddGithubOrganization(params.ProjectSFID, &v1Models.CreateGithubOrganization{
				OrganizationName: params.Body.OrganizationName,
			})
			if err != nil {
				return github_organizations.NewAddProjectGithubOrganizationBadRequest().WithPayload(errorResponse(err))
			}
			return github_organizations.NewAddProjectGithubOrganizationOK().WithPayload(*result)
		})

	api.GithubOrganizationsDeleteProjectGithubOrganizationHandler = github_organizations.DeleteProjectGithubOrganizationHandlerFunc(
		func(params github_organizations.DeleteProjectGithubOrganizationParams, claUser *auth.User) middleware.Responder {
			err := service.DeleteGithubOrganization(params.ProjectSFID, params.OrgName)
			if err != nil {
				return github_organizations.NewDeleteProjectGithubOrganizationBadRequest().WithPayload(errorResponse(err))
			}
			return github_organizations.NewDeleteProjectGithubOrganizationOK()
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
