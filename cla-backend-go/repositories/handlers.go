// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package repositories

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/github_repositories"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure establishes the middleware handlers for the repository service
func Configure(api *operations.ClaAPI, service Service) {
	api.GithubRepositoriesGetProjectGithubRepositoriesHandler = github_repositories.GetProjectGithubRepositoriesHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoriesParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.ListProjectRepositories(params.ProjectSFID)
			if err != nil {
				return github_repositories.NewGetProjectGithubRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			return github_repositories.NewGetProjectGithubRepositoriesOK().WithPayload(result)
		})

	api.GithubRepositoriesAddProjectGithubRepositoryHandler = github_repositories.AddProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.AddProjectGithubRepositoryParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.AddGithubRepository(params.ProjectSFID, params.GithubRepositoryInput)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			return github_repositories.NewAddProjectGithubRepositoryOK().WithPayload(result)
		})

	api.GithubRepositoriesDeleteProjectGithubRepositoryHandler = github_repositories.DeleteProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.DeleteProjectGithubRepositoryParams, claUser *user.CLAUser) middleware.Responder {
			err := service.DeleteGithubRepository(params.ProjectSFID, params.RepositoryID)
			if err != nil {
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			return github_repositories.NewDeleteProjectGithubRepositoryOK()
		})
}

// codedResponse interface
type codedResponse interface {
	Code() string
}

// errorResponse is a helper to wrap the specified error into an error response model
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
