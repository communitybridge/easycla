package template

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func Configure(api *operations.ClaAPI, service Service) {
	// Retrieve a list of available templates
	api.TemplateGetTemplateHandler = template.GetTemplateHandlerFunc(func(template template.GetTemplateParams) middleware.Responder {

		template, err := service.getTemplate(template.HTTPRequest.Context(), params)
		if err != nil {
			return template.NewGetTemplateBadRequest().WithPayload(errorResponse(err))
		}

		return template.NewGetTemplateOK().WithPayload(template)
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
