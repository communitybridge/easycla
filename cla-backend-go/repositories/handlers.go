package repositories

import (
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/github_repositories"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure establishes the middleware handlers for the repository service
func Configure(api *operations.ClaAPI, service Service, eventService events.Service) {
	api.GithubRepositoriesGetProjectGithubRepositoriesHandler = github_repositories.GetProjectGithubRepositoriesHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoriesParams, claUser *user.CLAUser) middleware.Responder {
			if !claUser.IsAuthorizedForProject(params.ProjectSFID) {
				return github_repositories.NewGetProjectGithubRepositoriesUnauthorized()
			}
			result, err := service.ListProjectRepositories(params.ProjectSFID)
			if err != nil {
				return github_repositories.NewGetProjectGithubRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			return github_repositories.NewGetProjectGithubRepositoriesOK().WithPayload(result)
		})

	api.GithubRepositoriesAddProjectGithubRepositoryHandler = github_repositories.AddProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.AddProjectGithubRepositoryParams, claUser *user.CLAUser) middleware.Responder {
			if !claUser.IsAuthorizedForProject(params.ProjectSFID) {
				return github_repositories.NewAddProjectGithubRepositoryUnauthorized()
			}
			result, err := service.AddGithubRepository(params.ProjectSFID, params.GithubRepositoryInput)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			addGithubRepositoryEvent(eventService, claUser, params.GithubRepositoryInput)
			return github_repositories.NewAddProjectGithubRepositoryOK().WithPayload(result)
		})

	api.GithubRepositoriesDeleteProjectGithubRepositoryHandler = github_repositories.DeleteProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.DeleteProjectGithubRepositoryParams, claUser *user.CLAUser) middleware.Responder {
			if !claUser.IsAuthorizedForProject(params.ProjectSFID) {
				return github_repositories.NewDeleteProjectGithubRepositoryUnauthorized()
			}
			ghRepo, err := service.GetGithubRepository(params.RepositoryID)
			if err != nil {
				if err == ErrGithubRepositoryNotFound {
					return github_repositories.NewDeleteProjectGithubRepositoryNotFound()
				}
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			err = service.DeleteGithubRepository(params.ProjectSFID, params.RepositoryID)
			if err != nil {
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			deleteGithubRepositoryEvent(eventService, claUser, ghRepo.RepositoryName, ghRepo.RepositoryProjectID)
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
