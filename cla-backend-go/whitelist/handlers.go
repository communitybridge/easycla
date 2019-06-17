package whitelist

import (
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/github"
	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service service, sessionStore *dynastore.Store) {

	api.DeleteGithubOrganizationFromClaHandler = operations.DeleteGithubOrganizationFromClaHandlerFunc(func(params operations.DeleteGithubOrganizationFromClaParams) middleware.Responder {
		err := service.DeleteGithubOrganizationFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID)
		if err != nil {
			return operations.NewDeleteGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewDeleteGithubOrganizationFromClaOK()
	})

	api.AddGithubOrganizationFromClaHandler = operations.AddGithubOrganizationFromClaHandlerFunc(func(params operations.AddGithubOrganizationFromClaParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			return operations.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			githubAccessToken = ""
		}

		err = service.AddGithubOrganizationToWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID, githubAccessToken)
		if err != nil {
			return operations.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewAddGithubOrganizationFromClaOK()
	})

	api.GetGithubOrganizationfromClaHandler = operations.GetGithubOrganizationfromClaHandlerFunc(func(params operations.GetGithubOrganizationfromClaParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			return operations.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			githubAccessToken = ""
		}

		result, err := service.GetGithubOrganizationsFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, githubAccessToken)
		if err != nil {
			return operations.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
		}

		return operations.NewGetGithubOrganizationfromClaOK().WithPayload(result)
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
