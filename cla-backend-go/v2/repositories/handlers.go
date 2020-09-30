// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/github"

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
func Configure(api *operations.EasyclaAPI, service Service, eventService events.Service) {
	api.GithubRepositoriesGetProjectGithubRepositoriesHandler = github_repositories.GetProjectGithubRepositoriesHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoriesParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_repositories.NewGetProjectGithubRepositoriesForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get GitHub Repositories with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			result, err := service.ListProjectRepositories(ctx, params.ProjectSFID)
			if err != nil {
				if strings.ContainsAny(err.Error(), "getProjectNotFound") {
					return github_repositories.NewGetProjectGithubRepositoriesNotFound().WithPayload(&models.ErrorResponse{
						Code:    "404",
						Message: fmt.Sprintf("project not found with given ID. [%s]", params.ProjectSFID),
					})
				}
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
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_repositories.NewAddProjectGithubRepositoryForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Add GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			result, err := service.AddGithubRepository(ctx, params.ProjectSFID, params.GithubRepositoryInput)
			if err != nil {
				return github_repositories.NewAddProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				EventType:         events.RepositoryAdded,
				ProjectID:         utils.StringValue(params.GithubRepositoryInput.ClaGroupID),
				ExternalProjectID: params.ProjectSFID,
				LfUsername:        authUser.UserName,
				ProjectModel: &v1Models.Project{
					ProjectExternalID: params.ProjectSFID,
					ProjectID:         utils.StringValue(params.GithubRepositoryInput.ClaGroupID),
				},
				EventData: &events.RepositoryAddedEventData{
					RepositoryName: result.RepositoryName,
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
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_repositories.NewDeleteProjectGithubRepositoryForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			ghRepo, err := service.GetRepository(ctx, params.RepositoryID)
			if err != nil {
				if err == repositories.ErrGithubRepositoryNotFound {
					return github_repositories.NewDeleteProjectGithubRepositoryNotFound()
				}
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}

			err = service.DisableRepository(ctx, params.RepositoryID)
			if err != nil {
				return github_repositories.NewDeleteProjectGithubRepositoryBadRequest().WithPayload(errorResponse(err))
			}

			eventService.LogEvent(&events.LogEventArgs{
				EventType:         events.RepositoryDisabled,
				ExternalProjectID: params.ProjectSFID,
				ProjectID:         ghRepo.RepositoryProjectID,
				LfUsername:        authUser.UserName,
				EventData: &events.RepositoryDisabledEventData{
					RepositoryName: ghRepo.RepositoryName,
				},
			})

			return github_repositories.NewDeleteProjectGithubRepositoryNoContent()
		})

	api.GithubRepositoriesGetProjectGithubRepositoryBranchProtectionHandler = github_repositories.GetProjectGithubRepositoryBranchProtectionHandlerFunc(
		func(params github_repositories.GetProjectGithubRepositoryBranchProtectionParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_repositories.NewGetProjectGithubRepositoryBranchProtectionForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Query Protected Branch GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			protectedBranch, err := service.GetProtectedBranch(ctx, params.RepositoryID)
			if err != nil {
				if err == repositories.ErrGithubRepositoryNotFound {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionNotFound()
				}
				if errors.Is(err, github.ErrAccessDenied) {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionForbidden().WithPayload(errorResponse(err))
				}

				// shall we return the actual code for rate liming ?
				if errors.Is(err, github.ErrRateLimited) {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionInternalServerError()
				}

				return github_repositories.NewGetProjectGithubRepositoryBranchProtectionBadRequest().WithPayload(errorResponse(err))
			}

			return github_repositories.NewGetProjectGithubRepositoryBranchProtectionOK().WithPayload(protectedBranch)
		})

	api.GithubRepositoriesUpdateProjectGithubRepositoryBranchProtectionHandler = github_repositories.UpdateProjectGithubRepositoryBranchProtectionHandlerFunc(
		func(params github_repositories.UpdateProjectGithubRepositoryBranchProtectionParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return github_repositories.NewUpdateProjectGithubRepositoryBranchProtectionForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Update Protected Branch GitHub Repository with Project scope of %s",
						authUser.UserName, params.ProjectSFID),
				})
			}

			protectedBranch, err := service.UpdateProtectedBranch(ctx, params.RepositoryID, params.GithubRepositoryBranchProtectionInput)
			if err != nil {
				log.Warnf("UpdateProjectGithubRepositoryBranchProtectionHandler : failed for repo %s : %v", params.RepositoryID, err)
				if err == repositories.ErrGithubRepositoryNotFound {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionNotFound()
				}
				if errors.Is(err, github.ErrAccessDenied) {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionForbidden().WithPayload(errorResponse(err))
				}

				// shall we return the actual code for rate liming ?
				if errors.Is(err, github.ErrRateLimited) {
					return github_repositories.NewGetProjectGithubRepositoryBranchProtectionInternalServerError().WithPayload(errorResponse(err))
				}

				return github_repositories.NewGetProjectGithubRepositoryBranchProtectionBadRequest().WithPayload(errorResponse(err))
			}

			return github_repositories.NewGetProjectGithubRepositoryBranchProtectionOK().WithPayload(protectedBranch)
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
