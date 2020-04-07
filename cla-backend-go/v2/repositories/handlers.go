package repositories

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/github_repositories"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure establishes the middleware handlers for the repository service
func Configure(api *operations.EasyclaAPI, service repositories.Service, eventService events.Service) {
	api.GithubRepositoriesGetProjectGithubRepositoriesHandler = github_repositories.GetProjectGithubRepositoriesHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoriesParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			result, err := service.ListProjectRepositories(params.ProjectSFID)
			if err != nil {
				return github_repositories.NewGetProjectGithubRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			return github_repositories.NewGetProjectGithubRepositoriesOK().WithPayload(*result)
		})

	api.GithubRepositoriesAddProjectGithubRepositoryHandler = github_repositories.AddProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.AddProjectGithubRepositoryParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			result, err := service.AddGithubRepository(params.ProjectSFID, &params.GithubRepositoryInput)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			eventService.LogEvent(&events.LogEventArgs{
				EventType:         events.GithubRepositoryAdded,
				ProjectID:         utils.StringValue(params.GithubRepositoryInput.RepositoryProjectID),
				ExternalProjectID: params.ProjectSFID,
				UserName:          authUser.UserName,
				EventData: &events.GithubRepositoryAddedEventData{
					RepositoryName: utils.StringValue(params.GithubRepositoryInput.RepositoryName),
				},
			})
			return github_repositories.NewAddProjectGithubRepositoryOK().WithPayload(*result)
		})

	api.GithubRepositoriesDeleteProjectGithubRepositoryHandler = github_repositories.DeleteProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.DeleteProjectGithubRepositoryParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ghRepo, err := service.GetGithubRepository(params.RepositoryID)
			if err != nil {
				if err == repositories.ErrGithubRepositoryNotFound {
					return github_repositories.NewDeleteProjectGithubRepositoryNotFound()
				}
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			err = service.DeleteGithubRepository(params.ProjectSFID, params.RepositoryID)
			if err != nil {
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}
			eventService.LogEvent(&events.LogEventArgs{
				EventType:         events.GithubRepositoryDeleted,
				ExternalProjectID: params.ProjectSFID,
				ProjectID:         ghRepo.RepositoryProjectID,
				LfUsername:        authUser.UserName,
				EventData: &events.GithubRepositoryDeletedEventData{
					RepositoryName: ghRepo.RepositoryName,
				},
			})
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
