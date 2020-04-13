// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"errors"
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	v1ProjectOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/go-openapi/runtime/middleware"
)

// errors
var (
	ErrProjectDoesNotExist = errors.New("project does not exist")
	ErrProjectIDMissing    = errors.New("project id is missing")
)

// Configure establishes the middleware handlers for the project service
func Configure(api *operations.EasyclaAPI, service v1Project.Service, eventsService events.Service) {
	api.ProjectCreateProjectHandler = project.CreateProjectHandlerFunc(func(params project.CreateProjectParams, user *auth.User) middleware.Responder {
		if !user.IsUserAuthorized(auth.Project, params.Body.ProjectExternalID) {
			return project.NewCreateProjectUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}
		if params.Body.ProjectName == "" || params.Body.ProjectACL == nil {
			msg := "Missing Project Name or Project ACL parameter."
			log.Warnf("Create Project Failed - %s", msg)
			return project.NewCreateProjectBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		exitingModel, getErr := service.GetProjectByName(params.Body.ProjectName)
		if getErr != nil {
			msg := fmt.Sprintf("Error querying the project by name, error: %+v", getErr)
			log.Warnf("Create Project Failed - %s", msg)
			return project.NewCreateProjectBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: msg,
			})
		}

		// If the project with the same name exists...
		if exitingModel != nil {
			msg := fmt.Sprintf("Project with same name exists: %s", params.Body.ProjectName)
			log.Warnf("Create Project Failed - %s", msg)
			return project.NewCreateProjectConflict().WithPayload(&models.ErrorResponse{
				Code:    "409",
				Message: msg,
			})
		}

		// Ok, safe to create now
		projectModel, err := service.CreateProject(&params.Body)
		if err != nil {
			log.Warnf("Create Project Failed - %+v", err)
			return project.NewCreateProjectBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.ProjectCreated,
			ProjectModel: projectModel,
			LfUsername:   user.UserName,
			EventData:    &events.ProjectCreatedEventData{},
		})

		log.Infof("Create Project Succeeded, project name: %s, project external ID: %s",
			params.Body.ProjectName, params.Body.ProjectExternalID)
		return project.NewCreateProjectOK().WithPayload(projectModel)
	})

	// Get Projects
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams, user *auth.User) middleware.Responder {

		projects, err := service.GetProjects(&v1ProjectOps.GetProjectsParams{
			HTTPRequest: params.HTTPRequest,
			FullMatch:   params.FullMatch,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			SearchField: params.SearchField,
			SearchTerm:  params.SearchTerm,
		})
		if err != nil {
			return project.NewGetProjectsBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectsOK().WithPayload(projects)
	})

	// Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(projectParams project.GetProjectByIDParams, user *auth.User) middleware.Responder {
		projectModel, err := service.GetProjectByID(projectParams.ProjectSfdcID)
		if err != nil {
			return project.NewGetProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByIDNotFound()
		}
		if !user.IsUserAuthorized(auth.Project, projectModel.ProjectExternalID) {
			return project.NewGetProjectByIDUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}

		return project.NewGetProjectByIDOK().WithPayload(projectModel)
	})

	api.ProjectGetProjectsByExternalIDHandler = project.GetProjectsByExternalIDHandlerFunc(func(projectParams project.GetProjectsByExternalIDParams, user *auth.User) middleware.Responder {
		if !user.IsUserAuthorized(auth.Project, projectParams.ExternalID) {
			return project.NewGetProjectsByExternalIDUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}

		projectModel, err := service.GetProjectsByExternalID(&v1ProjectOps.GetProjectsByExternalIDParams{
			HTTPRequest: projectParams.HTTPRequest,
			ExternalID:  projectParams.ExternalID,
			NextKey:     projectParams.NextKey,
			PageSize:    projectParams.PageSize,
		})
		if err != nil {
			return project.NewGetProjectsByExternalIDBadRequest().WithPayload(errorResponse(err))
		}
		return project.NewGetProjectsByExternalIDOK().WithPayload(projectModel)
	})

	// Get Project By Name
	api.ProjectGetProjectByNameHandler = project.GetProjectByNameHandlerFunc(func(projectParams project.GetProjectByNameParams, user *auth.User) middleware.Responder {

		projectModel, err := service.GetProjectByName(projectParams.ProjectName)
		if err != nil {
			return project.NewGetProjectByNameBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByNameNotFound()
		}
		if !user.IsUserAuthorized(auth.Project, projectModel.ProjectExternalID) {
			return project.NewGetProjectByNameUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}

		return project.NewGetProjectByNameOK().WithPayload(projectModel)
	})

	// Delete Project By ID
	api.ProjectDeleteProjectByIDHandler = project.DeleteProjectByIDHandlerFunc(func(projectParams project.DeleteProjectByIDParams, user *auth.User) middleware.Responder {
		log.Debugf("Processing delete request with project id: %s", projectParams.ProjectSfdcID)
		projectModel, err := service.GetProjectByID(projectParams.ProjectSfdcID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if !user.IsUserAuthorized(auth.Project, projectModel.ProjectExternalID) {
			return project.NewDeleteProjectByIDUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}
		err = service.DeleteProject(projectParams.ProjectSfdcID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.ProjectDeleted,
			ProjectModel: projectModel,
			LfUsername:   user.UserName,
			EventData:    &events.ProjectDeletedEventData{},
		})

		return project.NewDeleteProjectByIDNoContent()
	})

	// Update Project By ID
	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(func(projectParams project.UpdateProjectParams, user *auth.User) middleware.Responder {
		projectModel, err := service.GetProjectByID(projectParams.Body.ProjectID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if !user.IsUserAuthorized(auth.Project, projectModel.ProjectExternalID) {
			return project.NewUpdateProjectUnauthorized().WithPayload(&models.ErrorResponse{
				Code:    "401",
				Message: "user does not have access to this project",
			})
		}
		projectModel, err = service.UpdateProject(&projectParams.Body)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.ProjectUpdated,
			ProjectModel: projectModel,
			LfUsername:   user.UserName,
			EventData:    &events.ProjectUpdatedEventData{},
		})

		return project.NewUpdateProjectOK().WithPayload(projectModel)
	})

	// Project metrics
	api.ProjectGetProjectMetricsHandler = project.GetProjectMetricsHandlerFunc(func(projectParams project.GetProjectMetricsParams, user *auth.User) middleware.Responder {
		projectMetrics, err := service.GetMetrics()
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewGetProjectMetricsNotFound()
			}
			return project.NewGetProjectMetricsBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectMetricsOK().WithPayload(projectMetrics)
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
