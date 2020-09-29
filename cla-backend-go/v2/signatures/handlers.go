// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"

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

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, projectService project.Service, projectRepo project.ProjectRepository, companyService company.IService, v1SignatureService signatureService.SignatureService, sessionStore *dynastore.Store, eventsService events.Service, v2service Service, projectClaGroupsRepo projects_cla_groups.Repository) { //nolint

	// Get Signature
	api.SignaturesGetSignatureHandler = signatures.GetSignatureHandlerFunc(func(params signatures.GetSignatureParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		signature, err := v1SignatureService.GetSignature(ctx, params.SignatureID)
		if err != nil {
			log.Warnf("error retrieving signature metrics, error: %+v", err)
			return signatures.NewGetSignatureBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		if signature == nil {
			return signatures.NewGetSignatureNotFound().WithXRequestID(reqID)
		}
		resp, err := v2Signature(signature)
		if err != nil {
			return signatures.NewGetSignatureBadRequest().WithXRequestID(reqID)
		}

		return signatures.NewGetSignatureOK().WithXRequestID(reqID).WithPayload(resp)
	})

	api.SignaturesUpdateApprovalListHandler = signatures.UpdateApprovalListHandlerFunc(func(params signatures.UpdateApprovalListParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		// Must be in the Project|Organization Scope to see this
		if !utils.IsUserAuthorizedForProjectOrganizationTree(authUser, params.ProjectSFID, params.CompanySFID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to update Project Company Approval List with Project|Organization scope of %s | %s",
				authUser.UserName, params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewUpdateApprovalListForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
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
		companyModel, compErr := companyService.GetCompanyByExternalID(ctx, params.CompanySFID)
		if compErr != nil || companyModel == nil {
			log.Warnf("unable to locate company by external company ID: %s", params.CompanySFID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(compErr))
		}

		projectModels, projsErr := projectService.GetCLAGroupsByExternalSFID(ctx, params.ProjectSFID)
		if projsErr != nil || projectModels == nil {
			log.Warnf("unable to locate projects by Project SFID: %s", params.ProjectSFID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(projsErr))
		}

		// Lookup the internal project ID when provided the external ID via the v1SignatureService call
		projectModel, projErr := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if projErr != nil || projectModel == nil {
			log.Warnf("unable to locate project by CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(projErr))
		}

		// Convert the v2 input parameters to a v1 model
		v1ApprovalList := v1Models.ApprovalList{}
		err := copier.Copy(&v1ApprovalList, params.Body)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		// Invoke the update v1SignatureService function
		updatedSig, updateErr := v1SignatureService.UpdateApprovalList(ctx, authUser, projectModel, companyModel, params.ClaGroupID, &v1ApprovalList)
		if updateErr != nil || updatedSig == nil {
			if err, ok := err.(*signatureService.ForbiddenError); ok {
				return signatures.NewUpdateApprovalListForbidden().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			log.Warnf("unable to update signature approval list using CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(updateErr))
		}

		// Convert the v1 output model to a v2 response model
		v2Sig := models.Signature{}
		err = copier.Copy(&v2Sig, updatedSig)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		return signatures.NewUpdateApprovalListOK().WithXRequestID(reqID).WithPayload(&v2Sig)
	})

	// Retrieve GitHub Approval Entries
	api.SignaturesGetGitHubOrgWhitelistHandler = signatures.GetGitHubOrgWhitelistHandlerFunc(func(params signatures.GetGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		ghWhiteList, err := v1SignatureService.GetGithubOrganizationsFromWhitelist(ctx, params.SignatureID, githubAccessToken)
		if err != nil {
			log.Warnf("error fetching github organization whitelist entries v using signature_id: %s, error: %+v",
				params.SignatureID, err)
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		var response []models.GithubOrg
		err = copier.Copy(&response, ghWhiteList)
		if err != nil {
			return signatures.NewGetGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewGetGitHubOrgWhitelistOK().WithXRequestID(reqID).WithPayload(response)
	})

	// Add GitHub Approval Entries
	api.SignaturesAddGitHubOrgWhitelistHandler = signatures.AddGitHubOrgWhitelistHandlerFunc(func(params signatures.AddGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		ghApprovalList, err := v1SignatureService.AddGithubOrganizationToWhitelist(ctx, params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error adding github organization %s using signature_id: %s to the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		// Create an event
		signatureModel, getSigErr := v1SignatureService.GetSignature(ctx, params.SignatureID)
		var projectID = ""
		var companyID = ""
		if getSigErr != nil || signatureModel == nil {
			log.Warnf("error looking up signature using signature_id: %s, error: %+v",
				params.SignatureID, getSigErr)
		}
		if signatureModel != nil {
			projectID = signatureModel.ProjectID
			companyID = signatureModel.SignatureReferenceID.String()
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

		var response []models.GithubOrg
		err = copier.Copy(&response, ghApprovalList)
		if err != nil {
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewAddGitHubOrgWhitelistOK().WithXRequestID(reqID).WithPayload(response)
	})

	// Delete GitHub Approval List Entries
	api.SignaturesDeleteGitHubOrgWhitelistHandler = signatures.DeleteGitHubOrgWhitelistHandlerFunc(func(params signatures.DeleteGitHubOrgWhitelistParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		session, err := sessionStore.Get(params.HTTPRequest, github.SessionStoreKey)
		if err != nil {
			log.Warnf("error retrieving session from the session store, error: %+v", err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		ghApprovalList, err := v1SignatureService.DeleteGithubOrganizationFromWhitelist(ctx, params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error deleting github organization %s using signature_id: %s from the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		// Create an event
		signatureModel, getSigErr := v1SignatureService.GetSignature(ctx, params.SignatureID)
		var projectID = ""
		var companyID = ""
		if getSigErr != nil || signatureModel == nil {
			log.Warnf("error looking up signature using signature_id: %s, error: %+v",
				params.SignatureID, getSigErr)
		}
		if signatureModel != nil {
			projectID = signatureModel.ProjectID
			companyID = signatureModel.SignatureReferenceID.String()
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
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewDeleteGitHubOrgWhitelistNoContent().WithXRequestID(reqID).WithPayload(response)
	})

	// Get Project Signatures
	api.SignaturesGetProjectSignaturesHandler = signatures.GetProjectSignaturesHandlerFunc(func(params signatures.GetProjectSignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		projectSignatures, err := v1SignatureService.GetProjectSignatures(ctx, v1Signatures.GetProjectSignaturesParams{
			HTTPRequest:   params.HTTPRequest,
			FullMatch:     params.FullMatch,
			NextKey:       params.NextKey,
			PageSize:      params.PageSize,
			ProjectID:     params.ClaGroupID,
			SearchField:   params.SearchField,
			SearchTerm:    params.SearchTerm,
			SignatureType: params.SignatureType,
			ClaType:       params.ClaType,
		})
		if err != nil {
			log.Warnf("error retrieving project signatures for projectID: %s, error: %+v",
				params.ClaGroupID, err)
			return signatures.NewGetProjectSignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			return signatures.NewGetProjectSignaturesBadRequest().WithXRequestID(reqID)
		}

		if len(resp.Signatures) == 0 {
			return signatures.NewGetProjectSignaturesNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "404",
				Message: fmt.Sprintf("Signatures not found with given clagroupID. [%s]", params.ClaGroupID),
			})
		}

		return signatures.NewGetProjectSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get Project Company Signatures
	api.SignaturesGetProjectCompanySignaturesHandler = signatures.GetProjectCompanySignaturesHandlerFunc(func(params signatures.GetProjectCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		// Must be in the Organization Scope to see this
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s is not authorized to view project company signatures with Organization scope: %s",
				authUser.UserName, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: msg,
			})
		}

		projectSignatures, err := v2service.GetProjectCompanySignatures(ctx, params.CompanySFID, params.ProjectSFID)
		if err != nil {
			log.Warnf("error retrieving project signatures for project: %s, company: %s, error: %+v",
				params.ProjectSFID, params.CompanySFID, err)
			return signatures.NewGetProjectCompanySignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewGetProjectCompanySignaturesOK().WithXRequestID(reqID).WithPayload(projectSignatures)
	})

	// Get Employee Project Company Signatures
	api.SignaturesGetProjectCompanyEmployeeSignaturesHandler = signatures.GetProjectCompanyEmployeeSignaturesHandlerFunc(func(params signatures.GetProjectCompanyEmployeeSignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		projectSignatures, err := v1SignatureService.GetProjectCompanyEmployeeSignatures(ctx, v1Signatures.GetProjectCompanyEmployeeSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			CompanyID:   params.CompanyID,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			ProjectID:   params.ProjectID,
		})
		if err != nil {
			log.Warnf("error retrieving employee project signatures for project: %s, company: %s, error: %+v",
				params.ProjectID, params.CompanyID, err)
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID)
		}
		return signatures.NewGetProjectCompanyEmployeeSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get Company Signatures
	api.SignaturesGetCompanySignaturesHandler = signatures.GetCompanySignaturesHandlerFunc(func(params signatures.GetCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanyID) {
			msg := fmt.Sprintf("EasyCLA - 403 Forbidden - user %s is not authorized to view company signatures with Organization scope: %s",
				authUser.UserName, params.CompanyID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: msg,
			})
		}

		// Lookup the internal company ID when provided the external ID via the v1SignatureService call
		companyModel, compErr := companyService.GetCompanyByExternalID(ctx, params.CompanyID)
		if compErr != nil || companyModel == nil {
			log.Warnf("unable to locate company by external company ID: %s", params.CompanyID)
			return signatures.NewGetCompanySignaturesNotFound().WithXRequestID(reqID)
		}

		companySignatures, err := v1SignatureService.GetCompanySignatures(ctx, v1Signatures.GetCompanySignaturesParams{
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
			return signatures.NewGetCompanySignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		resp, err := v2SignaturesReplaceCompanyID(companySignatures, companyModel.CompanyID, companyModel.CompanyExternalID)
		if err != nil {
			log.Warnf("error converting company signatures for companyID: %s/%s, error: %+v",
				params.CompanyID, companyModel.CompanyID, err)
			return signatures.NewGetCompanySignaturesBadRequest().WithXRequestID(reqID)
		}
		return signatures.NewGetCompanySignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get User Signatures
	api.SignaturesGetUserSignaturesHandler = signatures.GetUserSignaturesHandlerFunc(func(params signatures.GetUserSignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		userSignatures, err := v1SignatureService.GetUserSignatures(ctx, v1Signatures.GetUserSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			UserName:    params.UserName,
			UserID:      params.UserID,
		})
		if err != nil {
			log.Warnf("error retrieving user signatures for userID: %s, error: %+v", params.UserID, err)
			return signatures.NewGetUserSignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}

		resp, err := v2Signatures(userSignatures)
		if err != nil {
			return signatures.NewGetUserSignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewGetUserSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Download ECLAs as a CSV document
	api.SignaturesDownloadProjectSignatureEmployeeAsCSVHandler = signatures.DownloadProjectSignatureEmployeeAsCSVHandlerFunc(
		func(params signatures.DownloadProjectSignatureEmployeeAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroupModel, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureEmployeeAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}

			result, err := v2service.GetClaGroupCorporateContributorsCsv(ctx, params.ClaGroupID, params.CompanySFID)
			if err != nil {
				if _, ok := err.(*organizations.GetOrgNotFound); ok {
					formatErr := errors.New("error retrieving company using companySFID")
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVNotFound().WithXRequestID(reqID).WithPayload(errorResponse(formatErr))
				}
				if ok := err.Error() == "not Found"; ok {
					message := fmt.Sprintf("request not found for Company ID: %s, Cla Group ID: %s",
						params.CompanySFID, params.ClaGroupID)
					formatErr := errors.New(message)
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVNotFound().WithXRequestID(reqID).WithPayload(errorResponse(formatErr))
				}
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
				rw.Header().Set("Content-Type", "text/csv")
				rw.Header().Set(utils.XREQUESTID, reqID)
				rw.WriteHeader(http.StatusOK)
				_, err := rw.Write(result)
				if err != nil {
					log.Warnf("Error writing csv file, error: %v", err)
				}
			})
		})

	api.SignaturesListClaGroupIclaSignatureHandler = signatures.ListClaGroupIclaSignatureHandlerFunc(
		func(params signatures.ListClaGroupIclaSignatureParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroupModel, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewListClaGroupIclaSignatureBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewListClaGroupIclaSignatureInternalServerError().WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewListClaGroupIclaSignatureForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}
			if !claGroupModel.ProjectICLAEnabled {
				return signatures.NewListClaGroupIclaSignatureBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - individual contribution is not supported for this project",
				})
			}
			result, err := v2service.GetProjectIclaSignatures(ctx, params.ClaGroupID, params.SearchTerm)
			if err != nil {
				return signatures.NewListClaGroupIclaSignatureInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return signatures.NewListClaGroupIclaSignatureOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.SignaturesListClaGroupCorporateContributorsHandler = signatures.ListClaGroupCorporateContributorsHandlerFunc(
		func(params signatures.ListClaGroupCorporateContributorsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint

			// Make sure the user has provided the companySFID
			if params.CompanySFID == nil {
				return signatures.NewDownloadProjectSignatureICLABadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - missing companySFID as input",
				})
			}

			// Lookup the CLA Group by ID - make sure it's valid
			claGroupModel, err := projectRepo.GetCLAGroupByID(params.ClaGroupID, project.DontLoadRepoDetails)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}

			// Authorized to view this? Allow project scope and project|org scope for matching IDs
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) && !utils.IsUserAuthorizedForProjectOrganization(authUser, claGroupModel.FoundationSFID, *params.CompanySFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to ListClaGroupCorporateContributors with project scope of %s or project|organization scope of %s|%s",
						authUser.UserName, claGroupModel.FoundationSFID, claGroupModel.FoundationSFID, *params.CompanySFID),
				})
			}

			// Make sure CCLA is enabled for this CLA Group
			if !claGroupModel.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureICLABadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - This project does not support corporate contribution",
				})
			}

			result, err := v2service.GetClaGroupCorporateContributors(ctx, params.ClaGroupID, params.CompanySFID, params.SearchTerm)
			if err != nil {
				return signatures.NewListClaGroupCorporateContributorsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return signatures.NewListClaGroupCorporateContributorsOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.SignaturesGetSignatureSignedDocumentHandler = signatures.GetSignatureSignedDocumentHandlerFunc(func(params signatures.GetSignatureSignedDocumentParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		signature, err := v1SignatureService.GetSignature(ctx, params.SignatureID)
		if err != nil {
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		if signature == nil {
			return signatures.NewGetSignatureSignedDocumentNotFound().WithXRequestID(reqID).WithPayload(errorResponse(errors.New("signature not found")))
		}
		haveAccess, err := isUserHaveAccessOfSignedSignaturePDF(ctx, authUser, signature, companyService, projectClaGroupsRepo)
		if err != nil {
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		if !haveAccess {
			return signatures.NewGetSignatureSignedDocumentForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:    "403",
				Message: "EasyCLA - 403 Forbidden : user does not have access of signature",
			})
		}
		doc, err := v2service.GetSignedDocument(ctx, signature.SignatureID.String())
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return signatures.NewGetSignatureSignedDocumentBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
		}
		return signatures.NewGetSignatureSignedDocumentOK().WithXRequestID(reqID).WithPayload(doc)
	})

	api.SignaturesDownloadProjectSignatureICLAsHandler = signatures.DownloadProjectSignatureICLAsHandlerFunc(
		func(params signatures.DownloadProjectSignatureICLAsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroup, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAsNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureICLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroup.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA: 403 Forbidden : User does not have permission to access project : %s", claGroup.FoundationSFID),
				})
			}
			if !claGroup.ProjectICLAEnabled {
				return signatures.NewDownloadProjectSignatureICLAsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA : 400 Bad Request : icla is not enabled on this project",
				})
			}
			result, err := v2service.GetSignedIclaZipPdf(params.ClaGroupID)
			if err != nil {
				if err == ErrZipNotPresent {
					return signatures.NewDownloadProjectSignatureICLAsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:    "404",
						Message: "EasyCLA: 404 Not found : no icla signatures found for this cla-group",
					})
				}
				return signatures.NewDownloadProjectSignatureICLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return signatures.NewDownloadProjectSignatureICLAsOK().WithXRequestID(reqID).WithPayload(result)

		})

	// Download ICLAs as a CSV document
	api.SignaturesDownloadProjectSignatureICLAAsCSVHandler = signatures.DownloadProjectSignatureICLAAsCSVHandlerFunc(
		func(params signatures.DownloadProjectSignatureICLAAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroupModel, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}
			if !claGroupModel.ProjectICLAEnabled {
				return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - individual contribution is not supported for this project",
				})
			}

			result, err := v2service.GetProjectIclaSignaturesCsv(ctx, params.ClaGroupID)
			if err != nil {
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
				rw.Header().Set("Content-Type", "text/csv")
				rw.Header().Set(utils.XREQUESTID, reqID)
				rw.WriteHeader(http.StatusOK)
				_, err := rw.Write(result)
				if err != nil {
					log.Warnf("Error writing csv file, error: %v", err)
				}
			})
		})

	api.SignaturesDownloadProjectSignatureCCLAsHandler = signatures.DownloadProjectSignatureCCLAsHandlerFunc(
		func(params signatures.DownloadProjectSignatureCCLAsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroup, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureCCLAsNotFound().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureCCLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroup.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureCCLAsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "403",
					Message: fmt.Sprintf("EasyCLA: 403 Forbidden : User does not have permission to access project : %s", claGroup.FoundationSFID),
				})
			}
			if !claGroup.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureCCLAsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA : 400 Bad Request : ccla is not enabled on this project",
				})
			}
			result, err := v2service.GetSignedCclaZipPdf(params.ClaGroupID)
			if err != nil {
				if err == ErrZipNotPresent {
					return signatures.NewDownloadProjectSignatureCCLAsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:    "404",
						Message: "EasyCLA: 404 Not found : no ccla signatures found for this cla-group",
					})
				}
				return signatures.NewDownloadProjectSignatureCCLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return signatures.NewDownloadProjectSignatureCCLAsOK().WithXRequestID(reqID).WithPayload(result)
		})

	// Download CCLAs as a CSV document
	api.SignaturesDownloadProjectSignatureCCLAAsCSVHandler = signatures.DownloadProjectSignatureCCLAAsCSVHandlerFunc(
		func(params signatures.DownloadProjectSignatureCCLAAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			claGroupModel, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureCCLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(err))
				}
				return signatures.NewDownloadProjectSignatureCCLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DownloadProjectSignatureCCLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
				})
			}
			if !claGroupModel.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: "EasyCLA - 400 Bad Request - corporate contribution is not supported for this project",
				})
			}

			result, err := v2service.GetProjectCclaSignaturesCsv(ctx, params.ClaGroupID)
			if err != nil {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(err))
			}
			return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
				rw.Header().Set("Content-Type", "text/csv")
				rw.Header().Set(utils.XREQUESTID, reqID)
				rw.WriteHeader(http.StatusOK)
				_, err := rw.Write(result)
				if err != nil {
					log.Warnf("Error writing csv file, error: %v", err)
				}
			})
		})
}

func isUserHaveAccessOfSignedSignaturePDF(ctx context.Context, authUser *auth.User, signature *v1Models.Signature, companyService company.IService, projectClaGroupRepo projects_cla_groups.Repository) (bool, error) {
	if authUser.Admin {
		return true, nil
	}
	projects, err := projectClaGroupRepo.GetProjectsIdsForClaGroup(signature.ProjectID)
	if err != nil {
		return false, err
	}
	if len(projects) == 0 {
		return false, fmt.Errorf("cannot find project(s) associated with CLA group cla_group_id: %s - please update the database with the foundation/project mapping",
			signature.ProjectID)
	}
	foundationID := projects[0].FoundationSFID
	projectSFID := projects[0].ProjectSFID

	pmScope := authUser.ResourceIDsByTypeAndRole(auth.Project, "project-manager")
	if len(pmScope) > 0 && utils.NewStringSetFromStringArray(pmScope).Include(foundationID) {
		return true, nil
	}
	if signature.SignatureType == CclaSignatureType {
		comp, err := companyService.GetCompany(ctx, signature.SignatureReferenceID.String())
		if err != nil {
			return false, err
		}
		expectedScope := fmt.Sprintf("%s|%s", projectSFID, comp.CompanyExternalID)
		cmScope := authUser.ResourceIDsByTypeAndRole(auth.ProjectOrganization, "cla-manager")
		if len(cmScope) > 0 && utils.NewStringSetFromStringArray(cmScope).Include(expectedScope) {
			return true, nil
		}
	}
	return false, nil
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
