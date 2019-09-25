// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service service, sessionStore *dynastore.Store) {

	api.CompanyAddGithubOrganizationFromClaHandler = company.AddGithubOrganizationFromClaHandlerFunc(
		func(params company.AddGithubOrganizationFromClaParams, claUser *user.CLAUser) middleware.Responder {
			session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
			if err != nil {
				log.Warnf("error retrieving session from the session store, error: %v", err)
				return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
			}

			githubAccessToken, ok := session.Values["github_access_token"].(string)
			if !ok {
				log.Debugf("no github access token in the session - initializing to empty string")
				githubAccessToken = ""
			}

			err = service.AddGithubOrganizationToWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID, githubAccessToken)
			if err != nil {
				log.Warnf("error adding github organization %v using company id: %s to the whitelist, error: %v",
					params.GithubOrganizationID.ID, params.CorporateClaID, err)
				return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewAddGithubOrganizationFromClaOK()
		})

	api.CompanyGetGithubOrganizationfromClaHandler = company.GetGithubOrganizationfromClaHandlerFunc(
		func(params company.GetGithubOrganizationfromClaParams, claUser *user.CLAUser) middleware.Responder {
			session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
			if err != nil {
				log.Warnf("error retrieving session from the session store, error: %v", err)
				return company.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
			}

			githubAccessToken, ok := session.Values["github_access_token"].(string)
			if !ok {
				log.Debugf("no github access token in the session - initializing to empty string")
				githubAccessToken = ""
			}

			result, err := service.GetGithubOrganizationsFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, githubAccessToken)
			if err != nil {
				log.Warnf("error fetching the github organization %v from the whitelist using company: %s, error: %v",
					params.GithubOrganizationID.ID, params.CorporateClaID, err)
				return company.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewGetGithubOrganizationfromClaOK().WithPayload(result)
		})

	api.CompanyDeleteGithubOrganizationFromClaHandler = company.DeleteGithubOrganizationFromClaHandlerFunc(
		func(params company.DeleteGithubOrganizationFromClaParams, claUser *user.CLAUser) middleware.Responder {
			err := service.DeleteGithubOrganizationFromWhitelist(params.HTTPRequest.Context(), params.CorporateClaID, *params.GithubOrganizationID.ID)
			if err != nil {
				log.Warnf("error deleting the github organization %v using company id: %s from the whitelist, error: %v",
					params.GithubOrganizationID.ID, params.CorporateClaID, err)
				return company.NewDeleteGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
			}

			return company.NewDeleteGithubOrganizationFromClaOK()
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
