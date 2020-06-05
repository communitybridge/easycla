// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
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

func v1ProjectModel(in *models.Project) (*v1Models.Project, error) {
	out := &v1Models.Project{}
	err := copier.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func v2ProjectModel(in *v1Models.Project) (*models.Project, error) {
	out := &models.Project{}
	err := copier.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Configure establishes the middleware handlers for the project service
func Configure(api *operations.EasyclaAPI, service v1Project.Service, v2Service Service, eventsService events.Service) { //nolint
	api.ProjectCreateProjectHandler = project.CreateProjectHandlerFunc(func(params project.CreateProjectParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProject(user, params.Body.ProjectExternalID) {
			return project.NewCreateProjectForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Create CLA Group with Project scope of %s",
					user.UserName, params.Body.ProjectExternalID),
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

		input, err := v1ProjectModel(&params.Body)
		if err != nil {
			return project.NewCreateProjectInternalServerError().WithPayload(errorResponse(err))
		}

		// Ok, safe to create now
		projectModel, err := service.CreateProject(input)
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

		result, err := v2ProjectModel(projectModel)
		if err != nil {
			return project.NewCreateProjectInternalServerError().WithPayload(errorResponse(err))
		}

		log.Infof("Create Project Succeeded, project name: %s, project external ID: %s",
			params.Body.ProjectName, params.Body.ProjectExternalID)
		return project.NewCreateProjectOK().WithPayload(result)
	})

	// Get Projects
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams, user *auth.User) middleware.Responder {

		// No auth checks - anyone can request the list of projects
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

		result := &models.Projects{}
		err = copier.Copy(result, projects)
		if err != nil {
			return project.NewGetProjectsInternalServerError().WithPayload(errorResponse(err))
		}
		return project.NewGetProjectsOK().WithPayload(result)
	})

	// Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(params project.GetProjectByIDParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		projectModel, err := service.GetProjectByID(params.ProjectSfdcID)
		if err != nil {
			return project.NewGetProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		if projectModel == nil {
			return project.NewGetProjectByIDNotFound()
		}

		if !utils.IsUserAuthorizedForProject(user, projectModel.ProjectExternalID) {
			return project.NewGetProjectByIDForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project By ID with Project scope of %s",
					user.UserName, projectModel.ProjectExternalID),
			})
		}

		result, err := v2ProjectModel(projectModel)
		if err != nil {
			return project.NewGetProjectByIDInternalServerError().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectByIDOK().WithPayload(result)
	})

	api.ProjectGetProjectsByExternalIDHandler = project.GetProjectsByExternalIDHandlerFunc(func(params project.GetProjectsByExternalIDParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProject(user, params.ExternalID) {
			return project.NewGetProjectsByExternalIDForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Projects By External ID with Project scope of %s",
					user.UserName, params.ExternalID),
			})
		}

		projectModel, err := service.GetProjectsByExternalID(&v1ProjectOps.GetProjectsByExternalIDParams{
			HTTPRequest: params.HTTPRequest,
			ProjectSFID: params.ExternalID,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
		})
		if err != nil {
			return project.NewGetProjectsByExternalIDBadRequest().WithPayload(errorResponse(err))
		}

		results := &models.Projects{}
		err = copier.Copy(results, projectModel)
		if err != nil {
			return project.NewGetProjectsByExternalIDInternalServerError().WithPayload(errorResponse(err))
		}
		return project.NewGetProjectsByExternalIDOK().WithPayload(results)
	})

	// Get Project By Name
	api.ProjectGetProjectByNameHandler = project.GetProjectByNameHandlerFunc(func(params project.GetProjectByNameParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)

		projectModel, err := service.GetProjectByName(params.ProjectName)
		if err != nil {
			return project.NewGetProjectByNameBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByNameNotFound()
		}

		if !utils.IsUserAuthorizedForProject(user, projectModel.ProjectExternalID) {
			return project.NewGetProjectByNameForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project By Name with Project scope of %s",
					user.UserName, projectModel.ProjectExternalID),
			})
		}

		result, err := v2ProjectModel(projectModel)
		if err != nil {
			return project.NewGetProjectByNameInternalServerError().WithPayload(errorResponse(err))
		}
		return project.NewGetProjectByNameOK().WithPayload(result)
	})

	// Delete Project By ID
	api.ProjectDeleteProjectByIDHandler = project.DeleteProjectByIDHandlerFunc(func(projectParams project.DeleteProjectByIDParams, user *auth.User) middleware.Responder {
		log.Debugf("Processing delete request with project id: %s", projectParams.ProjectSfdcID)
		utils.SetAuthUserProperties(user, projectParams.XUSERNAME, projectParams.XEMAIL)
		projectModel, err := service.GetProjectByID(projectParams.ProjectSfdcID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		if !utils.IsUserAuthorizedForProject(user, projectModel.ProjectExternalID) {
			return project.NewDeleteProjectByIDForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete Project By ID with Project scope of %s",
					user.UserName, projectModel.ProjectExternalID),
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
		utils.SetAuthUserProperties(user, projectParams.XUSERNAME, projectParams.XEMAIL)
		projectModel, err := service.GetProjectByID(projectParams.Body.ProjectID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectNotFound().WithPayload(errorResponse(err))
		}
		if !utils.IsUserAuthorizedForProject(user, projectModel.ProjectExternalID) {
			return project.NewUpdateProjectForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Update Project By ID with Project scope of %s",
					user.UserName, projectModel.ProjectExternalID),
			})
		}

		in, err := v1ProjectModel(&projectParams.Body)
		if err != nil {
			return project.NewUpdateProjectInternalServerError().WithPayload(errorResponse(err))
		}

		projectModel, err = service.UpdateProject(in)
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

		result, err := v2ProjectModel(projectModel)
		if err != nil {
			return project.NewUpdateProjectInternalServerError().WithPayload(errorResponse(err))
		}
		return project.NewUpdateProjectOK().WithPayload(result)
	})

	// Get CLA enabled projects
	api.ProjectGetCLAProjectsByIDHandler = project.GetCLAProjectsByIDHandlerFunc(func(projectParams project.GetCLAProjectsByIDParams, user *auth.User) middleware.Responder {
		// No auth checks - anyone including contributors can request
		claProjects, getErr := v2Service.GetCLAProjectsByID(projectParams.ProjectSfdcID)
		if getErr != nil {
			return project.NewGetCLAProjectsByIDBadRequest().WithPayload(errorResponse(getErr))
		}

		return project.NewGetCLAProjectsByIDOK().WithPayload(claProjects)
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
