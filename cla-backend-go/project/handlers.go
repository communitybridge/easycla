// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gerrits"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

const defaultPageSize int64 = 50

// isValidUser is a helper function to determine if the user object is valid
func isValidUser(claUser *user.CLAUser) bool {
	return claUser != nil && claUser.UserID != "" && claUser.LFUsername != "" && claUser.LFEmail != ""
}

// Configure establishes the middleware handlers for the project service
func Configure(api *operations.ClaAPI, service Service, eventsService events.Service, gerritService gerrits.Service, repositoryService repositories.Service, signatureService signatures.SignatureService) {
	// Create CLA Group/Project Handler
	api.ProjectCreateProjectHandler = project.CreateProjectHandlerFunc(func(params project.CreateProjectParams, claUser *user.CLAUser) middleware.Responder {
		if params.Body.ProjectName == "" || params.Body.ProjectACL == nil {
			msg := "Missing Project Name or Project ACL parameter."
			log.Warnf("Create Project Failed - %s", msg)
			return project.NewCreateProjectBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: msg,
			})
		}

		exitingModel, getErr := service.GetCLAGroupByName(params.Body.ProjectName)
		if getErr != nil {
			msg := fmt.Sprintf("Error querying the project by name, error: %+v", getErr)
			log.Warnf("Create Project Failed - %s", msg)
			return project.NewCreateProjectBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
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
		projectModel, err := service.CreateCLAGroup(&params.Body)
		if err != nil {
			log.Warnf("Create Project Failed - %+v", err)
			return project.NewCreateProjectBadRequest().WithPayload(errorResponse(err))
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.CLAGroupCreated,
			ProjectModel: projectModel,
			UserID:       claUser.UserID,
			EventData:    &events.CLAGroupCreatedEventData{},
		})

		log.Infof("Create Project Succeeded, project name: %s, project external ID: %s",
			params.Body.ProjectName, params.Body.ProjectExternalID)
		return project.NewCreateProjectOK().WithPayload(projectModel)
	})

	// Get Projects
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams, claUser *user.CLAUser) middleware.Responder {
		if !isValidUser(claUser) {
			return project.NewGetProjectsUnauthorized().WithPayload(&models.ErrorResponse{
				Message: "unauthorized",
				Code:    "401",
			})
		}
		projects, err := service.GetCLAGroups(&params)
		if err != nil {
			return project.NewGetProjectsBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectsOK().WithPayload(projects)
	})

	// Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(projectParams project.GetProjectByIDParams, claUser *user.CLAUser) middleware.Responder {

		projectModel, err := service.GetCLAGroupByID(projectParams.ProjectID)
		if err != nil {
			return project.NewGetProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByIDNotFound()
		}

		return project.NewGetProjectByIDOK().WithPayload(projectModel)
	})

	// Get Project By External ID Handler
	api.ProjectGetProjectsByExternalIDHandler = project.GetProjectsByExternalIDHandlerFunc(func(projectParams project.GetProjectsByExternalIDParams, claUser *user.CLAUser) middleware.Responder {

		log.Debugf("Project Handler - GetProjectsByExternalID")
		if projectParams.ProjectSFID == "" {
			return project.NewGetProjectsByExternalIDBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: "External ID is empty",
			})
		}

		// Set the default page size
		if projectParams.PageSize == nil {
			projectParams.PageSize = aws.Int64(defaultPageSize)
		}

		log.Debugf("Project Handler - GetProjectsByExternalID - invoking service")
		projectsModel, err := service.GetCLAGroupsByExternalID(&projectParams)
		if err != nil {
			return project.NewGetProjectsByExternalIDBadRequest().WithPayload(errorResponse(err))
		}
		if projectsModel == nil {
			return project.NewGetProjectsByExternalIDNotFound()
		}

		return project.NewGetProjectsByExternalIDOK().WithPayload(projectsModel)
	})

	// Get Project By Name
	api.ProjectGetProjectByNameHandler = project.GetProjectByNameHandlerFunc(func(projectParams project.GetProjectByNameParams, claUser *user.CLAUser) middleware.Responder {

		projectModel, err := service.GetCLAGroupByName(projectParams.ProjectName)
		if err != nil {
			return project.NewGetProjectByNameBadRequest().WithPayload(errorResponse(err))
		}
		if projectModel == nil {
			return project.NewGetProjectByNameNotFound()
		}

		return project.NewGetProjectByNameOK().WithPayload(projectModel)
	})

	// Delete Project By ID
	api.ProjectDeleteProjectByIDHandler = project.DeleteProjectByIDHandlerFunc(func(projectParams project.DeleteProjectByIDParams, claUser *user.CLAUser) middleware.Responder {
		f := logrus.Fields{
			"functionName":                "ProjectDeleteProjectByIDHandler",
			"claGroupID":                  projectParams.ProjectID,
			"authenticatedUserLFUsername": claUser.LFUsername,
			"authenticatedUserLFEmail":    claUser.LFEmail,
			"authenticatedUserUserID":     claUser.UserID,
			"authenticatedUserName":       claUser.Name,
		}
		log.WithFields(f).Debug("Processing delete request")
		projectModel, err := service.GetCLAGroupByID(projectParams.ProjectID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		// Delete gerrit repositories
		log.WithFields(f).Debug("Processing gerrit delete")
		howMany, err := gerritService.DeleteClaGroupGerrits(projectParams.ProjectID)
		if err != nil {
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		// Log gerrit event
		if howMany > 0 {
			log.WithFields(f).Debugf("Deleted %d gerrit groups", howMany)
			eventsService.LogEvent(&events.LogEventArgs{
				EventType:    events.GerritRepositoryDeleted,
				ProjectModel: projectModel,
				UserID:       claUser.UserID,
				EventData: &events.GerritProjectDeletedEventData{
					DeletedCount: howMany,
				},
			})
		}

		// Delete github repositories
		log.WithFields(f).Debug("Processing github repository delete")
		howMany, err = repositoryService.DeleteProject(projectParams.ProjectID)
		if err != nil {
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if howMany > 0 {
			log.WithFields(f).Debugf("Deleted %d github repositories", howMany)

			// Log github delete event
			eventsService.LogEvent(&events.LogEventArgs{
				EventType:    events.GithubRepositoryDeleted,
				ProjectModel: projectModel,
				UserID:       claUser.UserID,
				EventData: &events.GithubProjectDeletedEventData{
					DeletedCount: howMany,
				},
			})
		}

		// Invalidate project signatures
		log.WithFields(f).Debug("Invalidating signatures")
		howMany, err = signatureService.InvalidateProjectRecords(projectParams.ProjectID, projectModel.ProjectName)
		if err != nil {
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		if howMany > 0 {
			log.WithFields(f).Debugf("Invalidated %d signatures", howMany)
			// Log invalidate signatures
			eventsService.LogEvent(&events.LogEventArgs{
				EventType:    events.InvalidatedSignature,
				ProjectModel: projectModel,
				UserID:       claUser.UserID,
				EventData: &events.SignatureProjectInvalidatedEventData{
					InvalidatedCount: howMany,
				},
			})
		}

		err = service.DeleteCLAGroup(projectParams.ProjectID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.CLAGroupDeleted,
			ProjectModel: projectModel,
			UserID:       claUser.UserID,
			EventData:    &events.CLAGroupDeletedEventData{},
		})

		return project.NewDeleteProjectByIDNoContent()
	})

	// Update Project By Name
	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(func(projectParams project.UpdateProjectParams, claUser *user.CLAUser) middleware.Responder {

		exitingModel, getErr := service.GetCLAGroupByID(projectParams.Body.ProjectID)
		if getErr != nil {
			msg := fmt.Sprintf("Error querying the project by ID, error: %+v", getErr)
			log.Warnf("Update Project Failed - %s", msg)
			return project.NewCreateProjectBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: msg,
			})
		}

		// If the project with the same name exists...
		if exitingModel == nil {
			msg := fmt.Sprintf("unable to locate project with ID: %s", projectParams.Body.ProjectID)
			log.Warn(msg)
			return project.NewUpdateProjectNotFound().WithPayload(&models.ErrorResponse{
				Code:    "404",
				Message: msg,
			})
		}

		projectModel, err := service.UpdateCLAGroup(&projectParams.Body)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:    events.CLAGroupUpdated,
			ProjectModel: projectModel,
			UserID:       claUser.UserID,
			EventData:    &events.CLAGroupUpdatedEventData{},
		})

		return project.NewUpdateProjectOK().WithPayload(projectModel)
	})
}
