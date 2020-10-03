// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"context"
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
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		// No auth checks - anyone can request the list of projects
		projects, err := service.GetCLAGroups(ctx, &v1ProjectOps.GetProjectsParams{
			HTTPRequest: params.HTTPRequest,
			FullMatch:   params.FullMatch,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			SearchField: params.SearchField,
			SearchTerm:  params.SearchTerm,
		})
		if err != nil {
			return project.NewGetProjectsBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		result := &models.ClaGroups{}
		err = copier.Copy(result, projects)
		if err != nil {
			return project.NewGetProjectsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return project.NewGetProjectsOK().WithXRequestID(reqID).WithPayload(result)
	})

	// Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(params project.GetProjectByIDParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		claGroupModel, err := service.GetCLAGroupByID(ctx, params.ProjectSfdcID)
		if err != nil {

			if err.Error() == "project does not exist" {
				return project.NewGetProjectByIDNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return project.NewGetProjectByIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		if claGroupModel == nil {
			return project.NewGetProjectByIDNotFound().WithXRequestID(reqID)
		}

		if !utils.IsUserAuthorizedForProjectTree(user, claGroupModel.ProjectExternalID) {
			return project.NewGetProjectByIDForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project By ID with Project scope of %s",
					user.UserName, claGroupModel.ProjectExternalID),
			})
		}

		result, err := v2ProjectModel(claGroupModel)
		if err != nil {
			return project.NewGetProjectByIDInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		return project.NewGetProjectByIDOK().WithXRequestID(reqID).WithPayload(result)
	})

	api.ProjectGetProjectsByExternalIDHandler = project.GetProjectsByExternalIDHandlerFunc(func(params project.GetProjectsByExternalIDParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectTree(user, params.ExternalID) {
			return project.NewGetProjectsByExternalIDForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Projects By External ID with Project scope of %s",
					user.UserName, params.ExternalID),
			})
		}

		claGroupModel, err := service.GetCLAGroupsByExternalID(ctx, &v1ProjectOps.GetProjectsByExternalIDParams{
			HTTPRequest: params.HTTPRequest,
			ProjectSFID: params.ExternalID,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
		})
		if err != nil {
			return project.NewGetProjectsByExternalIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		results := &models.ClaGroups{}
		err = copier.Copy(results, claGroupModel)
		if err != nil {
			return project.NewGetProjectsByExternalIDInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		if results.Projects == nil {
			return project.NewGetProjectsByExternalIDNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "404",
				Message: fmt.Sprintf("project not found with id. [%s]", params.ExternalID),
			})
		}
		return project.NewGetProjectsByExternalIDOK().WithXRequestID(reqID).WithPayload(results)
	})

	// Get Project By Name
	api.ProjectGetProjectByNameHandler = project.GetProjectByNameHandlerFunc(func(params project.GetProjectByNameParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)

		claGroupModel, err := service.GetCLAGroupByName(ctx, params.ProjectName)
		if err != nil {
			return project.NewGetProjectByNameBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		if claGroupModel == nil {
			return project.NewGetProjectByNameNotFound().WithXRequestID(reqID)
		}

		if !utils.IsUserAuthorizedForProjectTree(user, claGroupModel.ProjectExternalID) {
			return project.NewGetProjectByNameForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project By Name with Project scope of %s",
					user.UserName, claGroupModel.ProjectExternalID),
			})
		}

		result, err := v2ProjectModel(claGroupModel)
		if err != nil {
			return project.NewGetProjectByNameInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return project.NewGetProjectByNameOK().WithXRequestID(reqID).WithPayload(result)
	})

	// Delete Project By ID
	api.ProjectDeleteProjectByIDHandler = project.DeleteProjectByIDHandlerFunc(func(params project.DeleteProjectByIDParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		f := logrus.Fields{
			"functionName":   "ProjectDeleteProjectByIDHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"projectSFID":    params.ProjectSfdcID,
			"userEmail":      user.Email,
			"userName":       user.UserName,
		}
		log.WithFields(f).Debug("Processing delete request")
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		claGroupModel, err := service.GetCLAGroupByID(ctx, params.ProjectSfdcID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewDeleteProjectByIDNotFound().WithXRequestID(reqID)
			}
			return project.NewDeleteProjectByIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		if !utils.IsUserAuthorizedForProjectTree(user, claGroupModel.ProjectExternalID) {
			return project.NewDeleteProjectByIDForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Delete Project By ID with Project scope of %s",
					user.UserName, claGroupModel.ProjectExternalID),
			})
		}

		err = service.DeleteCLAGroup(ctx, params.ProjectSfdcID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupDeleted,
			ClaGroupModel: claGroupModel,
			LfUsername:    user.UserName,
			EventData:     &events.CLAGroupDeletedEventData{},
		})

		return project.NewDeleteProjectByIDNoContent().WithXRequestID(reqID)
	})

	// Update Project By ID
	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(func(params project.UpdateProjectParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		claGroupModel, err := service.GetCLAGroupByID(ctx, params.Body.ProjectID)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		if !utils.IsUserAuthorizedForProjectTree(user, claGroupModel.ProjectExternalID) {
			return project.NewUpdateProjectForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Update Project By ID with Project scope of %s",
					user.UserName, claGroupModel.ProjectExternalID),
			})
		}

		in, err := v1ProjectModel(&params.Body)
		if err != nil {
			return project.NewUpdateProjectInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		claGroupModel, err = service.UpdateCLAGroup(ctx, in)
		if err != nil {
			if err == ErrCLAGroupDoesNotExist {
				return project.NewUpdateProjectNotFound().WithXRequestID(reqID)
			}
			return project.NewUpdateProjectBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:     events.CLAGroupUpdated,
			ClaGroupModel: claGroupModel,
			LfUsername:    user.UserName,
			EventData:     &events.CLAGroupUpdatedEventData{},
		})

		result, err := v2ProjectModel(claGroupModel)
		if err != nil {
			return project.NewUpdateProjectInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return project.NewUpdateProjectOK().WithXRequestID(reqID).WithPayload(result)
	})

	// Get CLA enabled projects
	api.ProjectGetCLAProjectsByIDHandler = project.GetCLAProjectsByIDHandlerFunc(func(params project.GetCLAProjectsByIDParams, user *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		// No auth checks - anyone including contributors can request
		claProjects, getErr := v2Service.GetCLAProjectsByID(ctx, params.FoundationSFID)
		if getErr != nil {
			return project.NewGetCLAProjectsByIDBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(getErr))
		}

		return project.NewGetCLAProjectsByIDOK().WithXRequestID(reqID).WithPayload(claProjects)
	})
}
