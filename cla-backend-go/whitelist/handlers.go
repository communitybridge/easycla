// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/github"

	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service service, sessionStore *dynastore.Store) {

	api.CompanyDeleteGithubOrganizationFromClaHandler = company.DeleteGithubOrganizationFromClaHandlerFunc(func(params company.DeleteGithubOrganizationFromClaParams) middleware.Responder {
		err := service.DeleteGithubOrganizationFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID)
		if err != nil {
			return company.NewDeleteGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewDeleteGithubOrganizationFromClaOK()
	})

	api.CompanyAddGithubOrganizationFromClaHandler = company.AddGithubOrganizationFromClaHandlerFunc(func(params company.AddGithubOrganizationFromClaParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			githubAccessToken = ""
		}

		err = service.AddGithubOrganizationToWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID, githubAccessToken)
		if err != nil {
			return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewAddGithubOrganizationFromClaOK()
	})

	api.CompanyGetGithubOrganizationfromClaHandler = company.GetGithubOrganizationfromClaHandlerFunc(func(params company.GetGithubOrganizationfromClaParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			return company.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			githubAccessToken = ""
		}

		result, err := service.GetGithubOrganizationsFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, githubAccessToken)
		if err != nil {
			return company.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewGetGithubOrganizationfromClaOK().WithPayload(result)
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
