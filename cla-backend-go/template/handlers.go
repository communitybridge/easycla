package template

import (
	"fmt"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations/template"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service service) {
	// Retrieve a list of available templates
	api.TemplateGetTemplatesHandler = template.GetTemplatesHandlerFunc(func(params template.GetTemplatesParams) middleware.Responder {

		templates, err := service.GetTemplates(params.HTTPRequest.Context())
		if err != nil {
			return template.NewGetTemplatesBadRequest().WithPayload(errorResponse(err))
		}
		return template.NewGetTemplatesOK().WithPayload(templates)
	})

	// Add Templates to a CLA Group
	api.TemplateCreateCLAGroupTemplateHandler = template.CreateCLAGroupTemplateHandlerFunc(func(params template.CreateCLAGroupTemplateParams) middleware.Responder {

		err := service.AddContractGroupTemplates(params.HTTPRequest.Context(), params.ClaGroupID)
		if err != nil {
			fmt.Print(err)
		}
		dummy := models.Template{}
		return template.NewCreateCLAGroupTemplateOK().WithPayload(dummy)
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
