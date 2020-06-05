package repositories

import (
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/github_repositories"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
)

// Configure establishes the middleware handlers for the repository service
func Configure(api *operations.EasyclaAPI, service repositories.Service, eventService events.Service) {
	api.GithubRepositoriesGetProjectGithubRepositoriesHandler = github_repositories.GetProjectGithubRepositoriesHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoriesParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
				return github_repositories.NewGetProjectGithubRepositoriesForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get GitHub Repositories with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			result, err := service.ListProjectRepositories(params.ProjectSFID)
			if err != nil {
				return github_repositories.NewGetProjectGithubRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			response := &models.ListGithubRepositories{}
			err = copier.Copy(response, result)
			if err != nil {
				return github_repositories.NewGetProjectGithubRepositoriesInternalServerError().WithPayload(errorResponse(err))
			}
			return github_repositories.NewGetProjectGithubRepositoriesOK().WithPayload(response)
		})

	api.GithubRepositoriesAddProjectGithubRepositoryHandler = github_repositories.AddProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.AddProjectGithubRepositoryParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
				return github_repositories.NewAddProjectGithubRepositoryForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Add GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			input := &v1Models.GithubRepositoryInput{}
			err := copier.Copy(input, &params.GithubRepositoryInput)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryInternalServerError().WithPayload(errorResponse(err))
			}

			result, err := service.AddGithubRepository(params.ProjectSFID, input)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				EventType:         events.GithubRepositoryAdded,
				ProjectID:         utils.StringValue(params.GithubRepositoryInput.RepositoryProjectID),
				ExternalProjectID: params.ProjectSFID,
				LfUsername:        authUser.UserName,
				EventData: &events.GithubRepositoryAddedEventData{
					RepositoryName: utils.StringValue(params.GithubRepositoryInput.RepositoryName),
				},
			})

			response := &models.GithubRepository{}
			err = copier.Copy(response, result)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryInternalServerError().WithPayload(errorResponse(err))
			}

			return github_repositories.NewAddProjectGithubRepositoryOK().WithPayload(response)
		})

	api.GithubRepositoriesDeleteProjectGithubRepositoryHandler = github_repositories.DeleteProjectGithubRepositoryHandlerFunc(
		func(params github_repositories.DeleteProjectGithubRepositoryParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
				return github_repositories.NewDeleteProjectGithubRepositoryForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

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
