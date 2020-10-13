// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"

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
			return signatures.NewGetSignatureBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			msg := fmt.Sprintf("%s - user %s does not have access to update Project Company Approval List with Project|Organization scope of %s | %s",
				utils.EasyCLA403Forbidden, authUser.UserName, params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewUpdateApprovalListForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    msg,
				XRequestID: reqID,
			})
		}

		// Valid the payload input - the validator will return a middleware.Responder response/error type
		validationError := validateApprovalListInput(reqID, params)
		if validationError != nil {
			return validationError
		}

		// Lookup the internal company ID when provided the external ID via the v1SignatureGService call
		companyModel, compErr := companyService.GetCompanyByExternalID(ctx, params.CompanySFID)
		if compErr != nil || companyModel == nil {
			log.Warnf("unable to locate company by external company ID: %s", params.CompanySFID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, compErr))
		}

		projectModels, projsErr := projectService.GetCLAGroupsByExternalSFID(ctx, params.ProjectSFID)
		if projsErr != nil || projectModels == nil {
			log.Warnf("unable to locate projects by Project SFID: %s", params.ProjectSFID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, projsErr))
		}

		// Lookup the internal project ID when provided the external ID via the v1SignatureService call
		projectModel, projErr := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
		if projErr != nil || projectModel == nil {
			log.Warnf("unable to locate project by CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, projErr))
		}

		// Convert the v2 input parameters to a v1 model
		v1ApprovalList := v1Models.ApprovalList{}
		err := copier.Copy(&v1ApprovalList, params.Body)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		// Invoke the update v1SignatureService function
		updatedSig, updateErr := v1SignatureService.UpdateApprovalList(ctx, authUser, projectModel, companyModel, params.ClaGroupID, &v1ApprovalList)
		if updateErr != nil || updatedSig == nil {
			if err, ok := err.(*signatureService.ForbiddenError); ok {
				return signatures.NewUpdateApprovalListForbidden().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}

			log.Warnf("unable to update signature approval list using CLA Group ID: %s", params.ClaGroupID)
			return signatures.NewUpdateApprovalListBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, updateErr))
		}

		// Convert the v1 output model to a v2 response model
		v2Sig := models.Signature{}
		err = copier.Copy(&v2Sig, updatedSig)
		if err != nil {
			return signatures.NewUpdateApprovalListInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewGetGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		var response []models.GithubOrg
		err = copier.Copy(&response, ghWhiteList)
		if err != nil {
			return signatures.NewGetGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		ghApprovalList, err := v1SignatureService.AddGithubOrganizationToWhitelist(ctx, params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error adding github organization %s using signature_id: %s to the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewAddGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewAddGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		githubAccessToken, ok := session.Values["github_access_token"].(string)
		if !ok {
			log.Debugf("no github access token in the session - initializing to empty string")
			githubAccessToken = ""
		}

		input := v1Models.GhOrgWhitelist{}
		err = copier.Copy(&input, &params.Body)
		if err != nil {
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		ghApprovalList, err := v1SignatureService.DeleteGithubOrganizationFromWhitelist(ctx, params.SignatureID, input, githubAccessToken)
		if err != nil {
			log.Warnf("error deleting github organization %s using signature_id: %s from the approval list, error: %+v",
				*params.Body.OrganizationID, params.SignatureID, err)
			return signatures.NewDeleteGitHubOrgWhitelistBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewDeleteGitHubOrgWhitelistInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
			return signatures.NewGetProjectSignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}
		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			return signatures.NewGetProjectSignaturesBadRequest().WithXRequestID(reqID)
		}

		return signatures.NewGetProjectSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get Project Company Signatures
	api.SignaturesGetProjectCompanySignaturesHandler = signatures.GetProjectCompanySignaturesHandlerFunc(func(params signatures.GetProjectCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		// Must be in the one of the above scopes to see this
		// - if project scope (like a PM)
		// - if project|organization scope (like CLA Manager, CLA Signatory)
		// - if organization scope (like company admin)
		if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) &&
			!utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) &&
			!utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			msg := fmt.Sprintf("%s - user %s is not authorized to view project company signatures any scope of project: %s, organization %s",
				utils.EasyCLA403Forbidden, authUser.UserName, params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    msg,
				XRequestID: reqID,
			})
		}

		projectSignatures, err := v2service.GetProjectCompanySignatures(ctx, params.CompanySFID, params.ProjectSFID)
		if err != nil {
			log.Warnf("error retrieving project signatures for project: %s, company: %s, error: %+v",
				params.ProjectSFID, params.CompanySFID, err)
			return signatures.NewGetProjectCompanySignaturesBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}
		return signatures.NewGetProjectCompanySignaturesOK().WithXRequestID(reqID).WithPayload(projectSignatures)
	})

	// Get Employee Project Company Signatures
	api.SignaturesGetProjectCompanyEmployeeSignaturesHandler = signatures.GetProjectCompanyEmployeeSignaturesHandlerFunc(func(params signatures.GetProjectCompanyEmployeeSignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		f := logrus.Fields{
			"functionName":   "SignaturesGetProjectCompanyEmployeeSignaturesHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"projectSFID":    params.ProjectSFID,
			"companySFID":    params.CompanySFID,
			"nextKey":        aws.StringValue(params.NextKey),
			"pageSize":       aws.Int64Value(params.PageSize),
		}

		// Try to load the company model - use both approaches - internal and external
		var companyModel *v1Models.Company
		var err error
		// Internal IDs are UUIDv4 - external are not
		if utils.IsUUIDv4(params.CompanySFID) {
			// Oops - not provided a SFID - but an internal ID - that's ok, we'll lookup via the internal ID
			log.WithFields(f).Debug("companySFID provided as internal ID - looking up record by internal ID")
			// Lookup the company model by internal ID
			companyModel, err = companyService.GetCompany(ctx, params.CompanySFID)
			if companyModel != nil && companyModel.CompanyExternalID == "" {
				log.WithFields(f).WithError(err).Warnf("problem loading company - company external ID not defined")
				return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseNotFound(
					reqID, fmt.Sprintf("problem loading company - company external ID not defined - comapny ID: %s", params.CompanySFID)))
			}
		} else {
			// Lookup the company model by external ID
			log.WithFields(f).Debug("companySFID provided as external ID - looking up record by external ID")
			companyModel, err = companyService.GetCompanyByExternalID(ctx, params.CompanySFID)
		}
		if err != nil {
			var companyDoesNotExistErr utils.CompanyDoesNotExist
			if errors.Is(err, &companyDoesNotExistErr) {
				log.WithFields(f).WithError(err).Warnf("problem loading company by ID")
				return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String404,
					Message:    fmt.Sprintf("%s - problem loading company by ID: %s - error: %+v", utils.EasyCLA404NotFound, params.CompanySFID, err),
					XRequestID: reqID,
				})
			}
			log.WithFields(f).WithError(err).Warnf("problem loading company by ID")
			return signatures.NewGetProjectCompanyEmployeeSignaturesInternalServerError().WithXRequestID(reqID).WithPayload(utils.ErrorResponseInternalServerErrorWithError(
				reqID, fmt.Sprintf("problem loading company by ID: %s", params.CompanySFID), err))
		}
		if companyModel == nil {
			log.WithFields(f).WithError(err).Warnf("problem loading loading company by ID")
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseNotFound(
				reqID, fmt.Sprintf("problem loading company by ID: %s", params.CompanySFID)))
		}

		log.WithFields(f).Debug("checking access control permissions...")
		if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) &&
			!utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) &&
			!utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			msg := fmt.Sprintf("%s - user %s is not authorized to view project company signatures any scope of project: %s, organization %s",
				utils.EasyCLA403Forbidden, authUser.UserName, params.ProjectSFID, params.CompanySFID)
			log.Warn(msg)
			return signatures.NewGetProjectCompanyEmployeeSignaturesForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    msg,
				XRequestID: reqID,
			})
		}

		// Locate the CLA Group for the provided project SFID
		projectCLAGroupModel, err := projectClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
		if err != nil {
			log.WithFields(f).WithError(err).Warnf("problem loading project -> cla group mapping")
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(
				reqID, fmt.Sprintf("problem loading project -> cla group mapping using project id: %s", params.ProjectSFID), err))
		}
		if projectCLAGroupModel == nil {
			log.WithFields(f).WithError(err).Warnf("problem loading project -> cla group mapping - no mapping found")
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequest(
				reqID, fmt.Sprintf("unable to locate cla group for project ID: %s", params.ProjectSFID)))
		}

		projectSignatures, err := v1SignatureService.GetProjectCompanyEmployeeSignatures(ctx, v1Signatures.GetProjectCompanyEmployeeSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			ProjectID:   projectCLAGroupModel.ClaGroupID, // cla group ID
			CompanyID:   companyModel.CompanyID,          // internal company id
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
		})
		if err != nil {
			log.WithFields(f).WithError(err).Warnf("error retrieving employee project signatures for project: %s, company: %s, error: %+v",
				params.ProjectSFID, params.CompanySFID, err)
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(
				reqID, fmt.Sprintf("unable to fetch employee signatures for project ID: %s and company: %s", params.ProjectSFID, params.CompanySFID), err))
		}

		resp, err := v2Signatures(projectSignatures)
		if err != nil {
			msg := fmt.Sprintf("error converting project company signatures for project: %s, company name: %s, companyID: %s, company external ID: %s",
				params.ProjectSFID, companyModel.CompanyName, companyModel.CompanyID, companyModel.CompanyExternalID)
			log.WithFields(f).WithError(err).Warn(msg)
			return signatures.NewGetProjectCompanyEmployeeSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(reqID, msg, err))
		}

		return signatures.NewGetProjectCompanyEmployeeSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get Company Signatures
	api.SignaturesGetCompanySignaturesHandler = signatures.GetCompanySignaturesHandlerFunc(func(params signatures.GetCompanySignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":   "SignaturesGetCompanySignaturesHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"companySFID":    params.CompanySFID,
			"companyName":    aws.StringValue(params.CompanyName),
			"signatureType":  aws.StringValue(params.SignatureType),
			"nextKey":        aws.StringValue(params.NextKey),
			"pageSize":       aws.Int64Value(params.PageSize),
		}

		// Lookup the internal company ID
		companyModel, err := companyService.GetCompanyByExternalID(ctx, params.CompanySFID)
		if err != nil {
			log.WithFields(f).WithError(err).Warnf("problem loading company by SFID - returning empty response")
			// Not sure this is the correct response as the LFX UI/Admin console wants 200 empty lists instead of non-200 status back
			return signatures.NewGetCompanySignaturesOK().WithXRequestID(reqID).WithPayload(&models.Signatures{
				Signatures:  []*models.Signature{},
				ResultCount: 0,
				TotalCount:  0,
			})
		}
		if companyModel == nil {
			log.WithFields(f).WithError(err).Warnf("problem loading company model by ID - returning empty response")
			// Not sure this is the correct response as the LFX UI/Admin console wants 200 empty lists instead of non-200 status back
			return signatures.NewGetCompanySignaturesOK().WithXRequestID(reqID).WithPayload(&models.Signatures{
				Signatures:  []*models.Signature{},
				ResultCount: 0,
				TotalCount:  0,
			})
		}

		if !utils.IsUserAuthorizedForOrganization(authUser, companyModel.CompanyExternalID) {
			msg := fmt.Sprintf("%s - user %s is not authorized to view company signatures with Organization scope: %s",
				utils.EasyCLA403Forbidden, authUser.UserName, companyModel.CompanyExternalID)
			log.WithFields(f).Warn(msg)
			return signatures.NewGetProjectCompanySignaturesForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    msg,
				XRequestID: reqID,
			})
		}

		log.WithFields(f).Debug("loading company signatures...")
		companySignatures, err := v1SignatureService.GetCompanySignatures(ctx, v1Signatures.GetCompanySignaturesParams{
			HTTPRequest:   params.HTTPRequest,
			CompanyID:     companyModel.CompanyID, // need to internal company ID here
			CompanyName:   params.CompanyName,
			NextKey:       params.NextKey,
			PageSize:      params.PageSize,
			SignatureType: params.SignatureType,
		})
		if err != nil {
			msg := fmt.Sprintf("%s - error retrieving company signatures for company name: %s, companyID: %s, company external ID: %s, error: %+v",
				utils.EasyCLA403Forbidden, companyModel.CompanyName, companyModel.CompanyID, companyModel.CompanyExternalID, err)
			log.WithFields(f).WithError(err).Warn(msg)
			return signatures.NewGetCompanySignaturesBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    msg,
				XRequestID: reqID,
			})
		}

		// Nothing in the query response - return a empty model
		if companySignatures == nil || len(companySignatures.Signatures) == 0 {
			return signatures.NewGetCompanySignaturesOK().WithXRequestID(reqID).WithPayload(&models.Signatures{
				Signatures:  []*models.Signature{},
				ResultCount: 0,
				TotalCount:  0,
			})
		}

		log.WithFields(f).Debug("updating company IDs...")
		resp, err := v2SignaturesReplaceCompanyID(companySignatures, companyModel.CompanyID, companyModel.CompanyExternalID)
		if err != nil {
			msg := fmt.Sprintf("error converting company signatures for company name: %s, companyID: %s, company external ID: %s",
				companyModel.CompanyName, companyModel.CompanyID, companyModel.CompanyExternalID)
			log.WithFields(f).WithError(err).Warn(msg)
			return signatures.NewGetCompanySignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(reqID, msg, err))
		}

		return signatures.NewGetCompanySignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Get User Signatures
	api.SignaturesGetUserSignaturesHandler = signatures.GetUserSignaturesHandlerFunc(func(params signatures.GetUserSignaturesParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		f := logrus.Fields{
			"functionName":   "SignaturesGetUserSignaturesHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"userID":         params.UserID,
			"userName":       aws.StringValue(params.UserName),
			"nextKey":        aws.StringValue(params.NextKey),
			"pageSize":       aws.Int64Value(params.PageSize),
		}

		userSignatures, err := v1SignatureService.GetUserSignatures(ctx, v1Signatures.GetUserSignaturesParams{
			HTTPRequest: params.HTTPRequest,
			NextKey:     params.NextKey,
			PageSize:    params.PageSize,
			UserName:    params.UserName,
			UserID:      params.UserID,
		})
		if err != nil {
			msg := fmt.Sprintf("error retrieving user signatures for userID: %s", params.UserID)
			log.WithFields(f).WithError(err).Warn(msg)
			return signatures.NewGetUserSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(reqID, msg, err))
		}

		resp, err := v2Signatures(userSignatures)
		if err != nil {
			msg := "problem converting signatures from v1 to v2"
			log.WithFields(f).WithError(err).Warn(msg)
			return signatures.NewGetUserSignaturesBadRequest().WithXRequestID(reqID).WithPayload(utils.ErrorResponseBadRequestWithError(reqID, msg, err))
		}

		return signatures.NewGetUserSignaturesOK().WithXRequestID(reqID).WithPayload(resp)
	})

	// Download ECLAs as a CSV document
	api.SignaturesDownloadProjectSignatureEmployeeAsCSVHandler = signatures.DownloadProjectSignatureEmployeeAsCSVHandlerFunc(
		func(params signatures.DownloadProjectSignatureEmployeeAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			f := logrus.Fields{
				"functionName":   "SignaturesDownloadProjectSignatureEmployeeAsCSVHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"claGroupID":     params.ClaGroupID,
				"companySFID":    params.CompanySFID,
			}
			log.WithFields(f).Debug("processing request...")

			claGroupModel, err := projectService.GetCLAGroupByID(ctx, params.ClaGroupID)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}

			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVForbidden().WithXRequestID(reqID).WithPayload(utils.ErrorResponseForbidden(
					reqID, fmt.Sprintf("user %s does not have access to DownloadProjectSignatureEmployeeAsCSV with project scope of %s", authUser.UserName, claGroupModel.FoundationSFID)))
			}

			result, err := v2service.GetClaGroupCorporateContributorsCsv(ctx, params.ClaGroupID, params.CompanySFID)
			if err != nil {
				if _, ok := err.(*organizations.GetOrgNotFound); ok {
					formatErr := errors.New("error retrieving company using companySFID")
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, formatErr))
				}
				if ok := err.Error() == "not Found"; ok {
					msg := fmt.Sprintf("request not found for Company ID: %s, Cla Group ID: %s", params.CompanySFID, params.ClaGroupID)
					return signatures.NewDownloadProjectSignatureEmployeeAsCSVNotFound().WithXRequestID(reqID).WithPayload(utils.ErrorResponseNotFoundWithError(reqID, msg, err))
				}
				return signatures.NewDownloadProjectSignatureEmployeeAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(utils.ErrorResponseInternalServerErrorWithError(
					reqID, fmt.Sprintf("problem getting corporate contributors CSV for CLA Group: %s with company: %s", params.ClaGroupID, params.CompanySFID), err))
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
					return signatures.NewListClaGroupIclaSignatureBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewListClaGroupIclaSignatureInternalServerError().WithPayload(errorResponse(reqID, err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewListClaGroupIclaSignatureForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: utils.String403,
					Message: fmt.Sprintf("%s - user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						utils.EasyCLA403Forbidden, authUser.UserName, claGroupModel.FoundationSFID),
					XRequestID: reqID,
				})
			}
			if !claGroupModel.ProjectICLAEnabled {
				return signatures.NewListClaGroupIclaSignatureBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - individual contribution is not supported for this project", utils.EasyCLA403Forbidden),
					XRequestID: reqID,
				})
			}
			result, err := v2service.GetProjectIclaSignatures(ctx, params.ClaGroupID, params.SearchTerm)
			if err != nil {
				return signatures.NewListClaGroupIclaSignatureInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - missing companySFID as input", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}

			// Lookup the CLA Group by ID - make sure it's valid
			claGroupModel, err := projectRepo.GetCLAGroupByID(ctx, params.ClaGroupID, project.DontLoadRepoDetails)
			if err != nil {
				if err == project.ErrProjectDoesNotExist {
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}

			// Authorized to view this? Allow project scope and project|org scope for matching IDs
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) && !utils.IsUserAuthorizedForProjectOrganization(authUser, claGroupModel.FoundationSFID, *params.CompanySFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: utils.String403,
					Message: fmt.Sprintf("%s - user %s does not have access to ListClaGroupCorporateContributors with project scope of %s or project|organization scope of %s|%s",
						utils.EasyCLA403Forbidden, authUser.UserName, claGroupModel.FoundationSFID, claGroupModel.FoundationSFID, *params.CompanySFID),
					XRequestID: reqID,
				})
			}

			// Make sure CCLA is enabled for this CLA Group
			if !claGroupModel.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureICLABadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - This project does not support corporate contribution", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}

			result, err := v2service.GetClaGroupCorporateContributors(ctx, params.ClaGroupID, params.CompanySFID, params.SearchTerm)
			if err != nil {
				return signatures.NewListClaGroupCorporateContributorsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			return signatures.NewListClaGroupCorporateContributorsOK().WithXRequestID(reqID).WithPayload(result)
		})

	api.SignaturesGetSignatureSignedDocumentHandler = signatures.GetSignatureSignedDocumentHandlerFunc(func(params signatures.GetSignatureSignedDocumentParams, authUser *auth.User) middleware.Responder {
		reqID := utils.GetRequestID(params.XREQUESTID)
		ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		f := logrus.Fields{
			"functionName":   "SignaturesGetSignatureSignedDocumentHandler",
			utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
			"signatureID":    params.SignatureID,
		}

		signatureModel, err := v1SignatureService.GetSignature(ctx, params.SignatureID)
		if err != nil {
			log.WithFields(f).WithError(err).Warn("problem loading signature")
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}
		if signatureModel == nil {
			log.WithFields(f).Warn("problem loading signature - signature not found")
			return signatures.NewGetSignatureSignedDocumentNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, errors.New("signature not found")))
		}

		haveAccess, err := isUserHaveAccessOfSignedSignaturePDF(ctx, authUser, signatureModel, companyService, projectClaGroupsRepo)
		if err != nil {
			log.WithFields(f).WithError(err).Warn("problem determining signature access")
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
		}

		if !haveAccess {
			return signatures.NewGetSignatureSignedDocumentForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
				Code:       utils.String403,
				Message:    fmt.Sprintf("%s - user %s does not have access to the specified signature", utils.EasyCLA403Forbidden, authUser.UserName),
				XRequestID: reqID,
			})
		}

		doc, err := v2service.GetSignedDocument(ctx, signatureModel.SignatureID.String())
		if err != nil {
			log.WithFields(f).WithError(err).Warn("problem fetching signed document")
			if strings.Contains(err.Error(), "bad request") {
				return signatures.NewGetSignatureSignedDocumentBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			return signatures.NewGetSignatureSignedDocumentInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
					return signatures.NewDownloadProjectSignatureICLAsNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureICLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroup.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String403,
					Message:    fmt.Sprintf("%s User does not have permission to access project : %s", utils.EasyCLA403Forbidden, claGroup.FoundationSFID),
					XRequestID: reqID,
				})
			}
			if !claGroup.ProjectICLAEnabled {
				return signatures.NewDownloadProjectSignatureICLAsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s icla is not enabled on this project", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}
			result, err := v2service.GetSignedIclaZipPdf(params.ClaGroupID)
			if err != nil {
				if err == ErrZipNotPresent {
					return signatures.NewDownloadProjectSignatureICLAsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:       utils.String404,
						Message:    fmt.Sprintf("%s no icla signatures found for this cla-group", utils.EasyCLA404NotFound),
						XRequestID: reqID,
					})
				}
				return signatures.NewDownloadProjectSignatureICLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
					return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureICLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: utils.String403,
					Message: fmt.Sprintf("user %s does not have access to DownloadProjectSignatureICLAAsCSV with project scope of %s",
						authUser.UserName, claGroupModel.FoundationSFID),
					XRequestID: reqID,
				})
			}
			if !claGroupModel.ProjectICLAEnabled {
				return signatures.NewDownloadProjectSignatureICLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - individual contribution is not supported for this project", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}

			result, err := v2service.GetProjectIclaSignaturesCsv(ctx, params.ClaGroupID)
			if err != nil {
				return signatures.NewDownloadProjectSignatureICLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
					return signatures.NewDownloadProjectSignatureCCLAsNotFound().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureCCLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroup.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureCCLAsForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String403,
					Message:    fmt.Sprintf("%s - User does not have permission to access project : %s", utils.EasyCLA403Forbidden, claGroup.FoundationSFID),
					XRequestID: reqID,
				})
			}
			if !claGroup.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureCCLAsBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - ccla is not enabled on this project", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}
			result, err := v2service.GetSignedCclaZipPdf(params.ClaGroupID)
			if err != nil {
				if err == ErrZipNotPresent {
					return signatures.NewDownloadProjectSignatureCCLAsNotFound().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
						Code:       utils.String404,
						Message:    fmt.Sprintf("%s - no ccla signatures found for this cla-group", utils.EasyCLA404NotFound),
						XRequestID: reqID,
					})
				}
				return signatures.NewDownloadProjectSignatureCCLAsInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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
					return signatures.NewDownloadProjectSignatureCCLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
				}
				return signatures.NewDownloadProjectSignatureCCLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
			}
			if !utils.IsUserAuthorizedForProjectTree(authUser, claGroupModel.FoundationSFID) {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVForbidden().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code: utils.String403,
					Message: fmt.Sprintf("%s - user %s does not have access to DownloadProjectSignatureCCLAAsCSV with project scope of %s",
						utils.EasyCLA403Forbidden, authUser.UserName, claGroupModel.FoundationSFID),
					XRequestID: reqID,
				})
			}
			if !claGroupModel.ProjectCCLAEnabled {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVBadRequest().WithXRequestID(reqID).WithPayload(&models.ErrorResponse{
					Code:       utils.String400,
					Message:    fmt.Sprintf("%s - corporate contribution is not supported for this project", utils.EasyCLA400BadRequest),
					XRequestID: reqID,
				})
			}

			result, err := v2service.GetProjectCclaSignaturesCsv(ctx, params.ClaGroupID)
			if err != nil {
				return signatures.NewDownloadProjectSignatureCCLAAsCSVInternalServerError().WithXRequestID(reqID).WithPayload(errorResponse(reqID, err))
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

// isUserHaveAccessOfSignedSignaturePDF returns true if the specified user has access to the provided signature, false otherwise
func isUserHaveAccessOfSignedSignaturePDF(ctx context.Context, authUser *auth.User, signature *v1Models.Signature, companyService company.IService, projectClaGroupRepo projects_cla_groups.Repository) (bool, error) {
	f := logrus.Fields{
		"functionName":           "isUserHaveAccessOfSignedSignaturePDF",
		utils.XREQUESTID:         ctx.Value(utils.XREQUESTID),
		"authUserName":           authUser.UserName,
		"authUserEmail":          authUser.Email,
		"signatureID":            signature.SignatureID,
		"claGroupID":             signature.ProjectID,
		"signatureType":          signature.SignatureType,
		"signatureReferenceType": signature.SignatureReferenceType,
	}

	projects, err := projectClaGroupRepo.GetProjectsIdsForClaGroup(signature.ProjectID)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("error loading load project IDs for CLA Group")
		return false, err
	}
	if len(projects) == 0 {
		log.WithFields(f).Warn("unable to locate any project IDs for CLA Group")
		return false, fmt.Errorf("cannot find project(s) associated with CLA group cla_group_id: %s - please update the database with the foundation/project mapping",
			signature.ProjectID)
	}

	// Foundation ID's should be all the same for each project ID - just grab the first one
	foundationID := projects[0].FoundationSFID

	// First, check for PM access
	if utils.IsUserAuthorizedForProjectTree(authUser, foundationID) {
		log.WithFields(f).Debugf("user is authorized for %s scope for foundation ID: %s", utils.ProjectScope, foundationID)
		return true, nil
	}

	// In case the project tree didn't pass, let's check the project list individually - if any has access, we return true
	for _, proj := range projects {
		if utils.IsUserAuthorizedForProject(authUser, proj.ProjectSFID) {
			log.WithFields(f).Debugf("user is authorized for %s scope for project ID: %s", utils.ProjectScope, proj.ProjectSFID)
			return true, nil
		}
	}

	// Corporate signature...we can check the company details
	if signature.SignatureType == CclaSignatureType {
		comp, err := companyService.GetCompany(ctx, signature.SignatureReferenceID.String())
		if err != nil {
			log.WithFields(f).WithError(err).Warnf("failed to load company record using signature reference id: %s", signature.SignatureReferenceID.String())
			return false, err
		}

		// No company SFID? Then, we can't check permissions...
		if comp == nil || comp.CompanyExternalID == "" {
			log.WithFields(f).Warnf("failed to load company record with external SFID using signature reference id: %s", signature.SignatureReferenceID.String())
			return false, err
		}

		// Check the project|org tree starting with the foundation
		if utils.IsUserAuthorizedForProjectOrganizationTree(authUser, foundationID, comp.CompanyExternalID) {
			return true, nil
		}

		// In case the project organization tree didn't pass, let's check the project list individually - if any has access, we return true
		for _, proj := range projects {
			if utils.IsUserAuthorizedForProjectOrganization(authUser, proj.ProjectSFID, comp.CompanyExternalID) {
				log.WithFields(f).Debugf("user is authorized for %s scope for project ID: %s, org iD: %s", utils.ProjectOrgScope, proj.ProjectSFID, comp.CompanyExternalID)
				return true, nil
			}
		}
	}

	log.WithFields(f).Debug("tried everything - user doesn't have access with project or project|org scope")
	return false, nil
}

type codedResponse interface {
	Code() string
}

func errorResponse(reqID string, err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:       code,
		Message:    err.Error(),
		XRequestID: reqID,
	}

	return &e
}
