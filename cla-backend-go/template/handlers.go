// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package template

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations/template"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service Service) {
	// Retrieve a list of available templates
	api.TemplateGetTemplatesHandler = template.GetTemplatesHandlerFunc(func(params template.GetTemplatesParams) middleware.Responder {

		templates, err := service.GetTemplates(params.HTTPRequest.Context())
		if err != nil {
			return template.NewGetTemplatesBadRequest().WithPayload(errorResponse(err))
		}
		return template.NewGetTemplatesOK().WithPayload(templates)
	})

	api.TemplateCreateCLAGroupTemplateHandler = template.CreateCLAGroupTemplateHandlerFunc(func(params template.CreateCLAGroupTemplateParams) middleware.Responder {
		pdfUrls, err := service.CreateCLAGroupTemplate(params.HTTPRequest.Context(), params.ClaGroupID, params.Body)
		if err != nil {
			return template.NewGetTemplatesBadRequest().WithPayload(errorResponse(err))
		}

		return template.NewCreateCLAGroupTemplateOK().WithPayload(&pdfUrls)
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
