package gerrits

import (
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/gerrits"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// ProjectService contains Project methods
type ProjectService interface {
	GetCLAGroupByID(projectID string) (*models.Project, error)
}

// Configure the gerrit api
func Configure(api *operations.ClaAPI, service Service, projectService ProjectService, eventService events.Service) {
	api.GerritsDeleteGerritHandler = gerrits.DeleteGerritHandlerFunc(
		func(params gerrits.DeleteGerritParams, claUser *user.CLAUser) middleware.Responder {
			projectModel, err := projectService.GetCLAGroupByID(params.ProjectID)
			if err != nil {
				return gerrits.NewDeleteGerritBadRequest().WithPayload(errorResponse(err))
			}
			// verify user have access to the project
			if !claUser.IsAuthorizedForProject(projectModel.ProjectExternalID) {
				return gerrits.NewDeleteGerritUnauthorized()
			}
			gerrit, err := service.GetGerrit(params.GerritID)
			if err != nil {
				return gerrits.NewDeleteGerritBadRequest().WithPayload(errorResponse(err))
			}
			// verify gerrit project is same as the request
			if gerrit.ProjectID != params.ProjectID {
				return gerrits.NewDeleteGerritBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "provided project id does not match with gerrit project id",
				})
			}
			// delete the gerrit
			err = service.DeleteGerrit(params.GerritID)
			if err != nil {
				return gerrits.NewDeleteGerritBadRequest().WithPayload(errorResponse(err))
			}
			// record the event
			eventService.LogEvent(&events.LogEventArgs{
				EventType:    events.GerritRepositoryDeleted,
				ProjectModel: projectModel,
				UserID:       claUser.UserID,
				EventData: &events.GerritDeletedEventData{
					GerritRepositoryName: gerrit.GerritName,
				},
			})
			return gerrits.NewDeleteGerritOK()
		})

	api.GerritsAddGerritHandler = gerrits.AddGerritHandlerFunc(
		func(params gerrits.AddGerritParams, claUser *user.CLAUser) middleware.Responder {
			projectModel, err := projectService.GetCLAGroupByID(params.ProjectID)
			if err != nil {
				return gerrits.NewAddGerritBadRequest().WithPayload(errorResponse(err))
			}
			// verify user have access to the project
			if !claUser.IsAuthorizedForProject(projectModel.ProjectExternalID) {
				return gerrits.NewAddGerritUnauthorized()
			}

			if len(strings.TrimSpace(*params.AddGerritInput.GerritName)) == 0 {
				return gerrits.NewGetGerritReposBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "invalid gerritName parameter - expecting gerrit hostname",
				})
			}
			params.AddGerritInput.Version = "v1"
			// add the gerrit
			result, err := service.AddGerrit(params.ProjectID, projectModel.ProjectExternalID, params.AddGerritInput)
			if err != nil {
				return gerrits.NewAddGerritBadRequest().WithPayload(errorResponse(err))
			}
			// record the event
			eventService.LogEvent(&events.LogEventArgs{
				EventType:    events.GerritRepositoryAdded,
				ProjectModel: projectModel,
				UserID:       claUser.UserID,
				EventData: &events.GerritAddedEventData{
					GerritRepositoryName: utils.StringValue(params.AddGerritInput.GerritName),
				},
			})
			return gerrits.NewAddGerritOK().WithPayload(result)
		})

	api.GerritsGetGerritReposHandler = gerrits.GetGerritReposHandlerFunc(
		func(params gerrits.GetGerritReposParams, authUser *user.CLAUser) middleware.Responder {

			// No specific permissions required

			// Validate input
			if params.GerritHost == nil {
				return gerrits.NewGetGerritReposBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "missing gerritHost query parameter",
				})
			}

			result, err := service.GetGerritRepos(params.GerritHost.String())
			if err != nil {
				return gerrits.NewGetGerritReposBadRequest().WithPayload(errorResponse(err))
			}

			return gerrits.NewGetGerritReposOK().WithPayload(result)
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
