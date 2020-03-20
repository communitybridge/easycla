// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
)

const defaultPageSize int64 = 50

// Configure establishes the middleware handlers for the project service
func Configure(api *operations.ClaAPI, service Service) {
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

		log.Infof("Create Project Succeeded, project name: %s, project external ID: %s",
			params.Body.ProjectName, params.Body.ProjectExternalID)
		return project.NewCreateProjectOK().WithPayload(projectModel)
	})

	// Get Projects
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams, claUser *user.CLAUser) middleware.Responder {

		projects, err := service.GetProjects(&params)
		if err != nil {
			return project.NewGetProjectsBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectsOK().WithPayload(projects)
	})

	// Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(projectParams project.GetProjectByIDParams, claUser *user.CLAUser) middleware.Responder {

		projectModel, err := service.GetProjectByID(projectParams.ProjectID)
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
		if projectParams.ExternalID == "" {
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
		projectsModel, err := service.GetProjectsByExternalID(&projectParams)
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

		projectModel, err := service.GetProjectByName(projectParams.ProjectName)
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
		log.Debugf("Processing delete request with project id: %s", projectParams.ProjectID)
		err := service.DeleteProject(projectParams.ProjectID)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewDeleteProjectByIDNotFound()
			}
			return project.NewDeleteProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewDeleteProjectByIDNoContent()
	})

	// Update Project By ID
	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(func(projectParams project.UpdateProjectParams, claUser *user.CLAUser) middleware.Responder {
		projectModel, err := service.UpdateProject(&projectParams.Body)
		if err != nil {
			if err == ErrProjectDoesNotExist {
				return project.NewUpdateProjectNotFound()
			}
			return project.NewUpdateProjectBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewUpdateProjectOK().WithPayload(projectModel)
	})

	// Project metrics
	api.ProjectGetProjectMetricsHandler = project.GetProjectMetricsHandlerFunc(func(projectParams project.GetProjectMetricsParams, claUser *user.CLAUser) middleware.Responder {
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
