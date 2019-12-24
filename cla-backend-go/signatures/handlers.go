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
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/go-openapi/runtime/middleware"
	"github.com/savaki/dynastore"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service SignatureService, sessionStore *dynastore.Store) {

	// Retrieve GitHub Whitelist Entries
	api.SignaturesGetGitHubOrgWhitelistHandler = signatures.GetGitHubOrgWhitelistHandlerFunc(func(params signatures.GetGitHubOrgWhitelistParams, claUser *user.CLAUser) middleware.Responder {
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
	api.SignaturesAddGitHubOrgWhitelistHandler = signatures.AddGitHubOrgWhitelistHandlerFunc(func(params signatures.AddGitHubOrgWhitelistParams, claUser *user.CLAUser) middleware.Responder {
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
			log.Warnf("error adding github organization %s using signature_id: %s to the whitelist, error: %v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return company.NewAddGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewAddGithubOrganizationFromClaOK().WithPayload(ghWhiteList)
	})

	// Delete GitHub Whitelist Entries
	api.SignaturesDeleteGitHubOrgWhitelistHandler = signatures.DeleteGitHubOrgWhitelistHandlerFunc(func(params signatures.DeleteGitHubOrgWhitelistParams, claUser *user.CLAUser) middleware.Responder {

		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %v", err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := service.DeleteGithubOrganizationFromWhitelist(params.HTTPRequest.Context(), params.SignatureID, params.Body, githubAccessToken)
		if err != nil {
			log.Warnf("error deleting github organization %s using signature_id: %s from the whitelist, error: %v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return company.NewDeleteGithubOrganizationFromClaBadRequest().WithPayload(errorResponse(err))
		}

		return company.NewDeleteGithubOrganizationFromClaOK().WithPayload(ghWhiteList)
	})

	// Get Signatures
	api.SignaturesGetSignaturesHandler = signatures.GetSignaturesHandlerFunc(func(params signatures.GetSignaturesParams, claUser *user.CLAUser) middleware.Responder {
		signatureList, err := service.GetSignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving user signatures for signatureID: %s, error: %v", params.SignatureID, err)
			return signatures.NewGetSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetSignaturesOK().WithPayload(signatureList)
	})

	// Get Project Signatures
	api.SignaturesGetProjectSignaturesHandler = signatures.GetProjectSignaturesHandlerFunc(func(params signatures.GetProjectSignaturesParams, claUser *user.CLAUser) middleware.Responder {
		projectSignatures, err := service.GetProjectSignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving project signatures for projectID: %s, error: %v",
				params.ProjectID, err)
			return signatures.NewGetProjectSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetProjectSignaturesOK().WithPayload(projectSignatures)
	})

	// Get Project Company Signatures
	api.SignaturesGetProjectCompanySignaturesHandler = signatures.GetProjectCompanySignaturesHandlerFunc(func(params signatures.GetProjectCompanySignaturesParams, claUser *user.CLAUser) middleware.Responder {
		projectSignatures, err := service.GetProjectCompanySignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving project signatures for project: %s, company: %s, error: %+v",
				params.ProjectID, params.CompanyID, err)
			return signatures.NewGetProjectCompanySignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetProjectCompanySignaturesOK().WithPayload(projectSignatures)
	})

	// Get Employee Project Company Signatures
	api.SignaturesGetProjectCompanyEmployeeSignaturesHandler = signatures.GetProjectCompanyEmployeeSignaturesHandlerFunc(func(params signatures.GetProjectCompanyEmployeeSignaturesParams, claUser *user.CLAUser) middleware.Responder {
		projectSignatures, err := service.GetProjectCompanyEmployeeSignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving employee project signatures for project: %s, company: %s, error: %+v",
				params.ProjectID, params.CompanyID, err)
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetProjectCompanyEmployeeSignaturesOK().WithPayload(projectSignatures)
	})

	// Get Company Signatures
	api.SignaturesGetCompanySignaturesHandler = signatures.GetCompanySignaturesHandlerFunc(func(params signatures.GetCompanySignaturesParams, claUser *user.CLAUser) middleware.Responder {
		companySignatures, err := service.GetCompanySignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving company signatures for companyID: %s, error: %v", params.CompanyID, err)
			return signatures.NewGetCompanySignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetCompanySignaturesOK().WithPayload(companySignatures)
	})

	// Get User Signatures
	api.SignaturesGetUserSignaturesHandler = signatures.GetUserSignaturesHandlerFunc(func(params signatures.GetUserSignaturesParams, claUser *user.CLAUser) middleware.Responder {
		userSignatures, err := service.GetUserSignatures(params.HTTPRequest.Context(), params)
		if err != nil {
			log.Warnf("error retrieving user signatures for userID: %s, error: %v", params.UserID, err)
			return signatures.NewGetUserSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		return signatures.NewGetUserSignaturesOK().WithPayload(userSignatures)
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
