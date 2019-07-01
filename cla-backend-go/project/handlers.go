// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"github.com/communitybridge/easy-cla/cla-backend-go/gen/models"
	"github.com/communitybridge/easy-cla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easy-cla/cla-backend-go/gen/restapi/operations/project"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service service) {
	//Get Project By ID
	api.ProjectGetProjectsHandler = project.GetProjectsHandlerFunc(func(params project.GetProjectsParams) middleware.Responder {

		projects, err := service.GetProjects(params.HTTPRequest.Context())
		if err != nil {
			return project.NewGetProjectsBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectsOK().WithPayload(projects)
	})

	//Get Project By ID
	api.ProjectGetProjectByIDHandler = project.GetProjectByIDHandlerFunc(func(projectParams project.GetProjectByIDParams) middleware.Responder {

		sfdcProject, err := service.GetProjectByID(projectParams.HTTPRequest.Context(), projectParams.ProjectSfdcID)
		if err != nil {
			return project.NewGetProjectByIDBadRequest().WithPayload(errorResponse(err))
		}

		return project.NewGetProjectByIDOK().WithPayload(sfdcProject)
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
