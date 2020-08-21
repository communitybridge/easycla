// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/copier"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

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

// Configure establishes the middleware handlers for the project service
func Configure(api *operations.EasyclaAPI, service v1Project.Service, v2Service Service, eventsService events.Service) { //nolint
	// Get Projects
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams, user *auth.User) middleware.Responder {

		// No auth checks - anyone can request the list of projects
		projects, err := service.GetCLAGroups(&v1ProjectOps.GetProjectsParams{
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
		projectModel, err := service.GetCLAGroupByID(params.ProjectSfdcID)
		if err != nil {

			if err.Error() == "project does not exist" {
				return project.NewGetProjectByIDNotFound().WithPayload(errorResponse(err))
			}
			return project.NewGetProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		if projectModel == nil {
			return project.NewGetProjectByIDNotFound()
		}

		if !utils.IsUserAuthorizedForProjectTree(user, projectModel.ProjectExternalID) {
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
		if !utils.IsUserAuthorizedForProjectTree(user, params.ExternalID) {
			return project.NewGetProjectsByExternalIDForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Projects By External ID with Project scope of %s",
					user.UserName, params.ExternalID),
			})
		}

		projectModel, err := service.GetCLAGroupsByExternalID(&v1ProjectOps.GetProjectsByExternalIDParams{
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
		if results.Projects == nil {
			return project.NewGetProjectsByExternalIDNotFound().WithPayload(&models.ErrorResponse{
				Code:    "404",
				Message: fmt.Sprintf("project not found with id. [%s]", params.ExternalID),
			})
		}
		return project.NewGetProjectsByExternalIDOK().WithPayload(results)
	})

	// Get Project By Name
	api.ProjectGetProjectByNameHandler = project.GetProjectByNameHandlerFunc(func(params project.GetProjectByNameParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)

		projectModel, err := service.GetCLAGroupByName(params.ProjectName)
		if err != nil {
			return project.NewGetProjectByNameBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByNameNotFound()
		}

		if !utils.IsUserAuthorizedForProjectTree(user, projectModel.ProjectExternalID) {
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
		f := logrus.Fields{
			"functionName": "ProjectDeleteProjectByIDHandler",
			"projectSFID":  projectParams.ProjectSfdcID,
			"userEmail":    user.Email,
			"userName":     user.UserName,
		}
		log.WithFields(f).Debug("Processing delete request")
		utils.SetAuthUserProperties(user, projectParams.XUSERNAME, projectParams.XEMAIL)
		projectModel, err := service.GetCLAGroupByID(projectParams.ProjectSfdcID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		if !utils.IsUserAuthorizedForProjectTree(user, projectModel.ProjectExternalID) {
			return project.NewDeleteProjectByIDForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete Project By ID with Project scope of %s",
					user.UserName, projectModel.ProjectExternalID),
			})
		}

		err = service.DeleteCLAGroup(projectParams.ProjectSfdcID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.CLAGroupDeleted,
			ProjectModel: projectModel,
			LfUsername:   user.UserName,
			EventData:    &events.CLAGroupDeletedEventData{},
		})

		return project.NewDeleteProjectByIDNoContent()
	})

	// Update Project By ID
	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(func(projectParams project.UpdateProjectParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, projectParams.XUSERNAME, projectParams.XEMAIL)
		projectModel, err := service.GetCLAGroupByID(projectParams.Body.ProjectID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectNotFound().WithPayload(errorResponse(err))
		}
		if !utils.IsUserAuthorizedForProjectTree(user, projectModel.ProjectExternalID) {
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

		projectModel, err = service.UpdateCLAGroup(in)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectBadRequest().WithPayload(errorResponse(err))
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.CLAGroupUpdated,
			ProjectModel: projectModel,
			LfUsername:   user.UserName,
			EventData:    &events.CLAGroupUpdatedEventData{},
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
		claProjects, getErr := v2Service.GetCLAProjectsByID(projectParams.FoundationSFID)
		if getErr != nil {
			return project.NewGetCLAProjectsByIDBadRequest().WithPayload(errorResponse(getErr))
		}

		return project.NewGetCLAProjectsByIDOK().WithPayload(claProjects)
	})
}
