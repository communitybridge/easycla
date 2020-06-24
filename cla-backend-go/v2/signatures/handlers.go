// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/communitybridge/easycla/cla-backend-go/project"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/events"
	v1Signatures "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	signatureService "github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"github.com/savaki/dynastore"
)

func v2Signature(src *v1Models.Signature) (*models.Signature, error) {
	var dst models.Signature
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func v2Signatures(src *v1Models.Signatures) (*models.Signatures, error) {
	var dst models.Signatures
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func v2SignaturesReplaceCompanyID(src *v1Models.Signatures, internalID, externalID string) (*models.Signatures, error) {
	var dst models.Signatures
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}

	// Resplace the internal ID with the External ID
	for _, sig := range dst.Signatures {
		if sig.SignatureReferenceID == internalID {
			sig.SignatureReferenceID = externalID
		}
	}

	return &dst, nil
}

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, projectService project.Service, companyService company.IService, v1SignatureService signatureService.SignatureService, sessionStore *dynastore.Store, eventsService events.Service, v2service Service) { //nolint

	// Get Signature
	api.SignaturesGetSignatureHandler = signatures.GetSignatureHandlerFunc(func(params signatures.GetSignatureParams, authUser *auth.User) middleware.Responder {

		signature, err := v1SignatureService.GetSignature(params.SignatureID)
		if err != nil {
			log.Warnf("error retrieving signature metrics, error: %+v", err)
			return signatures.NewGetSignatureBadRequest().WithPayload(errorResponse(err))
		}

		if signature == nil {
			return signatures.NewGetSignatureNotFound()
		}
		resp, err := v2Signature(signature)
		if err != nil {
			return signatures.NewGetSignatureBadRequest()
		}

		return signatures.NewGetSignatureOK().WithPayload(resp)
	})

	api.SignaturesUpdateApprovalListHandler = signatures.UpdateApprovalListHandlerFunc(func(params signatures.UpdateApprovalListParams, authUser *auth.User) middleware.Responder {
		if params.XEMAIL == nil || params.XUSERNAME == nil || params.XACL == "" {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - unknown user is not authorized to update project company signature approval list for project ID: %s, company ID: %s",
				params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewUpdateApprovalListForbidden().WithPayload(errorResponse(errors.New(msg)))
		}

		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		// Must be in the Project|Organization Scope to see this
		if !utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to update Project Company Approval List with Project|Organization scope of %s | %s",
				authUser.UserName, params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewUpdateApprovalListForbidden().WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: msg,
			})
		}

		// Valid the payload input - the validator will return a middleware.Responder response/error type
		validationError := validateApprovalListInput(params)
		if validationError != nil {
			return validationError
		}

		// Lookup the internal company ID when provided the external ID via the v1SignatureService call
		companyModel, compErr := companyService.GetCompanyByExternalID(params.CompanySFID)
		if compErr != nil || companyModel == nil {
			log.Warnf("unable to locate company by external company ID: %s", params.CompanySFID)
			return signatures.NewUpdateApprovalListNotFound().WithPayload(errorResponse(compErr))
		}

		projectModels, projsErr := projectService.GetProjectsByExternalSFID(params.ProjectSFID)
		if projsErr != nil || projectModels == nil {
			log.Warnf("unable to locate projects by Project SFID: %s", params.ProjectSFID)
			return signatures.NewUpdateApprovalListNotFound().WithPayload(errorResponse(projsErr))
		}

		// Lookup the internal project ID when provided the external ID via the v1SignatureService call
		projectModel, projErr := projectService.GetProjectByID(params.ClaGroupID)
		if projErr != nil || projectModel == nil {
			log.Warnf("unable to locate project by CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListNotFound().WithPayload(errorResponse(projErr))
		}

		// Convert the v2 input parameters to a v1 model
		v1ApprovalList := v1Models.ApprovalList{}
		err := copier.Copy(&v1ApprovalList, params.Body)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithPayload(errorResponse(err))
		}

		// Invoke the update v1SignatureService function
		updatedSig, updateErr := v1SignatureService.UpdateApprovalList(authUser, projectModel, companyModel, params.ClaGroupID, &v1ApprovalList)
		if updateErr != nil || updatedSig == nil {
			if err, ok := err.(*signatureService.ForbiddenError); ok {
				return signatures.NewUpdateApprovalListForbidden().WithPayload(errorResponse(err))
			}

			log.Warnf("unable to update signature approval list using CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListBadRequest().WithPayload(errorResponse(updateErr))
		}

		// Convert the v1 output model to a v2 response model
		v2Sig := models.Signature{}
		err = copier.Copy(&v2Sig, updatedSig)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithPayload(errorResponse(err))
		}

		return signatures.NewUpdateApprovalListOK().WithPayload(&v2Sig)
	})

	// Retrieve GitHub Approval Entries
	api.SignaturesGetGitHubOrgWhitelistHandler = signatures.GetGitHubOrgWhitelistHandlerFunc(func(params signatures.GetGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := v1SignatureService.GetGithubOrganizationsFromWhitelist(params.SignatureID, githubAccessToken)
		if err != nil {
			log.Warnf("error fetching github organization whitelist entries v using signature_id: %s, error: %+v",
				params.SignatureID, err)
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}
		response := []models.GithubOrg{}
		err = copier.Copy(&response, ghWhiteList)
		if err != nil {
			return signatures.NewGetGitHubOrgWhitelistInternalServerError().WithPayload(errorResponse(err))
		}
		return signatures.NewGetGitHubOrgWhitelistOK().WithPayload(response)
	})

	// Add GitHub Approval Entries
	api.SignaturesAddGitHubOrgWhitelistHandler = signatures.AddGitHubOrgWhitelistHandlerFunc(func(params signatures.AddGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithPayload(errorResponse(err))
		}

		ghApprovalList, err := v1SignatureService.AddGithubOrganizationToWhitelist(params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error adding github organization %s using signature_id: %s to the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		// Create an event
		signatureModel, getSigErr := v1SignatureService.GetSignature(params.SignatureID)
		var projectID = ""
		var companyID = ""
		if getSigErr != nil || signatureModel == nil {
			log.Warnf("error looking up signature using signature_id: %s, error: %+v",
				params.SignatureID, getSigErr)
		}
		if signatureModel != nil {
			projectID = signatureModel.ProjectID
			companyID = signatureModel.SignatureReferenceID
		}

		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.ApprovalListGithubOrganizationAdded,
			ProjectID:  projectID,
			CompanyID:  companyID,
			LfUsername: authUser.UserName,
			EventData: &events.ApprovalListGithubOrganizationAddedEventData{
				GithubOrganizationName: utils.StringValue(params.Body.OrganizationID),
			},
		})
		response := []models.GithubOrg{}
		err = copier.Copy(&response, ghApprovalList)
		if err != nil {
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithPayload(errorResponse(err))
		}
		return signatures.NewAddGitHubOrgWhitelistOK().WithPayload(response)
	})

	// Delete GitHub Approval List Entries
	api.SignaturesDeleteGitHubOrgWhitelistHandler = signatures.DeleteGitHubOrgWhitelistHandlerFunc(func(params signatures.DeleteGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithPayload(errorResponse(err))
		}

		ghApprovalList, err := v1SignatureService.DeleteGithubOrganizationFromWhitelist(params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error deleting github organization %s using signature_id: %s from the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithPayload(errorResponse(err))
		}

		// Create an event
		signatureModel, getSigErr := v1SignatureService.GetSignature(params.SignatureID)
		var projectID = ""
		var companyID = ""
		if getSigErr != nil || signatureModel == nil {
			log.Warnf("error looking up signature using signature_id: %s, error: %+v",
				params.SignatureID, getSigErr)
		}
		if signatureModel != nil {
			projectID = signatureModel.ProjectID
			companyID = signatureModel.SignatureReferenceID
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.ApprovalListGithubOrganizationDeleted,
			ProjectID:  projectID,
			CompanyID:  companyID,
			LfUsername: authUser.UserName,
			EventData: &events.ApprovalListGithubOrganizationDeletedEventData{
				GithubOrganizationName: utils.StringValue(params.Body.OrganizationID),
			},
		})
		var response []models.GithubOrg
		err = copier.Copy(&response, ghApprovalList)
		if err != nil {
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithPayload(errorResponse(err))
		}
		return signatures.NewDeleteGitHubOrgWhitelistNoContent().WithPayload(response)
	})

	// Get Project Signatures
	api.SignaturesGetProjectSignaturesHandler = signatures.GetProjectSignaturesHandlerFunc(func(params signatures.GetProjectSignaturesParams, authUser *auth.User) middleware.Responder {
		projectSignatures, err := v1SignatureService.GetProjectSignatures(v1Signatures.GetProjectSignaturesParams{
			HTTPRequest:   params.HTTPRequest,
			FullMatch:     params.FullMatch,
			NextKey:       params.NextKey,
			PageSize:      params.PageSize,
			ProjectID:     params.ProjectID,
			SearchField:   params.SearchField,
			SearchTerm:    params.SearchTerm,
			SignatureType: params.SignatureType,
		})
		if err != nil {
			log.Warnf("error retrieving project signatures for projectID: %s, error: %+v",
				params.ProjectID, err)
			return signatures.NewGetProjectSignaturesBadRequest().WithPayload(errorResponse(err))
		}
		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			return signatures.NewGetProjectSignaturesBadRequest()
		}

		return signatures.NewGetProjectSignaturesOK().WithPayload(resp)
	})

	// Get Project Company Signatures
	api.SignaturesGetProjectCompanySignaturesHandler = signatures.GetProjectCompanySignaturesHandlerFunc(func(params signatures.GetProjectCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		// Must be in the Organization Scope to see this
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s is not authorized to view project company signatures with Organization scope: %s",
				authUser.UserName, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: msg,
			})
		}

		projectSignatures, err := v2service.GetProjectCompanySignatures(params.CompanySFID, params.ProjectSFID)
		if err != nil {
			log.Warnf("error retrieving project signatures for project: %s, company: %s, error: %+v",
				params.ProjectSFID, params.CompanySFID, err)
			return signatures.NewGetProjectCompanySignaturesBadRequest().WithPayload(errorResponse(err))
		}
		return signatures.NewGetProjectCompanySignaturesOK().WithPayload(projectSignatures)
	})

	// Get Employee Project Company Signatures
	api.SignaturesGetProjectCompanyEmployeeSignaturesHandler = signatures.GetProjectCompanyEmployeeSignaturesHandlerFunc(func(params signatures.GetProjectCompanyEmployeeSignaturesParams, authUser *auth.User) middleware.Responder {
		projectSignatures, err := v1SignatureService.GetProjectCompanyEmployeeSignatures(v1Signatures.GetProjectCompanyEmployeeSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			CompanyID:   params.CompanyID,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			ProjectID:   params.ProjectID,
		})
		if err != nil {
			log.Warnf("error retrieving employee project signatures for project: %s, company: %s, error: %+v",
				params.ProjectID, params.CompanyID, err)
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest()
		}
		return signatures.NewGetProjectCompanyEmployeeSignaturesOK().WithPayload(resp)
	})

	// Get Company Signatures
	api.SignaturesGetCompanySignaturesHandler = signatures.GetCompanySignaturesHandlerFunc(func(params signatures.GetCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanyID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s is not authorized to view company signatures with Organization scope: %s",
				authUser.UserName, params.CompanyID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: msg,
			})
		}

		// Lookup the internal company ID when provided the external ID via the v1SignatureService call
		companyModel, compErr := companyService.GetCompanyByExternalID(params.CompanyID)
		if compErr != nil || companyModel == nil {
			log.Warnf("unable to locate company by external company ID: %s", params.CompanyID)
			return signatures.NewGetCompanySignaturesNotFound()
		}

		companySignatures, err := v1SignatureService.GetCompanySignatures(v1Signatures.GetCompanySignaturesParams{
			HTTPRequest:   params.HTTPRequest,
			CompanyID:     companyModel.CompanyID, // need to internal company ID here
			CompanyName:   params.CompanyName,
			NextKey:       params.NextKey,
			PageSize:      params.PageSize,
			SignatureType: params.SignatureType,
		})
		if err != nil {
			log.Warnf("error retrieving company signatures for companyID: %s/%s, error: %+v",
				params.CompanyID, companyModel.CompanyID, err)
			return signatures.NewGetCompanySignaturesBadRequest().WithPayload(errorResponse(err))
		}

		resp, err := v2SignaturesReplaceCompanyID(companySignatures, companyModel.CompanyID, companyModel.CompanyExternalID)
		if err != nil {
			log.Warnf("error converting company signatures for companyID: %s/%s, error: %+v",
				params.CompanyID, companyModel.CompanyID, err)
			return signatures.NewGetCompanySignaturesBadRequest()
		}
		return signatures.NewGetCompanySignaturesOK().WithPayload(resp)
	})

	// Get User Signatures
	api.SignaturesGetUserSignaturesHandler = signatures.GetUserSignaturesHandlerFunc(func(params signatures.GetUserSignaturesParams, authUser *auth.User) middleware.Responder {
		userSignatures, err := v1SignatureService.GetUserSignatures(v1Signatures.GetUserSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			UserName:    params.UserName,
			UserID:      params.UserID,
		})
		if err != nil {
			log.Warnf("error retrieving user signatures for userID: %s, error: %+v", params.UserID, err)
			return signatures.NewGetUserSignaturesBadRequest().WithPayload(errorResponse(err))
		}

		resp, err := v2Signatures(userSignatures)
		if err != nil {
			return signatures.NewGetUserSignaturesBadRequest().WithPayload(errorResponse(err))
		}
		return signatures.NewGetUserSignaturesOK().WithPayload(resp)
	})
	api.SignaturesDownloadProjectSignatureICLAAsCSVHandler = signatures.DownloadProjectSignatureICLAAsCSVHandlerFunc(
		func(params signatures.DownloadProjectSignatureICLAAsCSVParams, authUser *auth.User) middleware.Responder {
			claGroupModel, err := projectService.GetProjectByID(params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProject(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}

			result, err := v2service.GetProjectIclaSignaturesCsv(params.ClaGroupID)
			if err != nil {
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithPayload(errorResponse(err))
			}
			return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
				rw.Header().Set("Content-Type", "text/csv")
				rw.WriteHeader(http.StatusOK)
				_, err := rw.Write(result)
				if err != nil {
					log.Warnf("Error writing csv file, error: %v", err)
				}
			})
		})
	api.SignaturesListClaGroupIclaSignatureHandler = signatures.ListClaGroupIclaSignatureHandlerFunc(
		func(params signatures.ListClaGroupIclaSignatureParams, authUser *auth.User) middleware.Responder {
			claGroupModel, err := projectService.GetProjectByID(params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProject(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}
			result, err := v2service.GetProjectIclaSignatures(params.ClaGroupID, params.SearchTerm)
			if err != nil {
				return signatures.NewListClaGroupIclaSignatureInternalServerError().WithPayload(errorResponse(err))
			}
			return signatures.NewListClaGroupIclaSignatureOK().WithPayload(result)
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
