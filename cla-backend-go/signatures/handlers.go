// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service SignatureService, sessionStore *dynastore.Store) {

	// Retrieve GitHub Whitelist Entries
	api.SignaturesGetGitHubOrgWhitelistHandler = signatures.GetGitHubOrgWhitelistHandlerFunc(func(params signatures.GetGitHubOrgWhitelistParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %v", err)
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := service.GetGithubOrganizationsFromWhitelist(params.HTTPRequest.Context(), params.SignatureID, githubAccessToken)
		if err != nil {
			log.Warnf("error fetching github organization whitelist entries v using signature_id: %s, error: %v",
				params.SignatureID, err)
			return company.NewGetGithubOrganizationfromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewGetGithubOrganizationfromClaOK().WithPayload(ghWhiteList)
	})

	// Add GitHub Whitelist Entries
	api.SignaturesAddGitHubOrgWhitelistHandler = signatures.AddGitHubOrgWhitelistHandlerFunc(func(params signatures.AddGitHubOrgWhitelistParams) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %v", err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := service.AddGithubOrganizationToWhitelist(params.HTTPRequest.Context(), params.SignatureID, params.Body, githubAccessToken)
		if err != nil {
			log.Warnf("error adding github organization %v using signature_id: %s to the whitelist, error: %v",
				params.Body.OrganizationID, params.SignatureID, err)
			return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewAddGithubOrganizationFromClaOK().WithPayload(ghWhiteList)
	})

	// Delete GitHub Whitelist Entries
	api.SignaturesDeleteGitHubOrgWhitelistHandler = signatures.DeleteGitHubOrgWhitelistHandlerFunc(func(params signatures.DeleteGitHubOrgWhitelistParams) middleware.Responder {

		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %v", err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := service.DeleteGithubOrganizationFromWhitelist(params.HTTPRequest.Context(), params.SignatureID, params.Body, githubAccessToken)
		if err != nil {
			log.Warnf("error deleting github organization %v using signature_id: %s from the whitelist, error: %v",
				params.Body.OrganizationID, params.SignatureID, err)
			return company.NewDeleteGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewDeleteGithubOrganizationFromClaOK().WithPayload(ghWhiteList)
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
