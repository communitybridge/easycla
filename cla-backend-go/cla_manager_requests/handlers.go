// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager_requests

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/aws/aws-sdk-go/aws"
	sigAPI "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/cla_manager_requests"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	"github.com/go-openapi/runtime/middleware"
)

// isValidUser is a helper function to determine if the user object is valid
func isValidUser(claUser *user.CLAUser) bool {
	return claUser != nil && claUser.UserID != "" && claUser.LFUsername != "" && claUser.LFEmail != ""
}

// Configure is the API handler routine for the CLA manager routes
func Configure(api *operations.ClaAPI, service IService, companyService company.IService, projectService project.Service, usersService users.Service, sigService signatures.SignatureService, eventsService events.Service, corporateConsoleURL string) { // nolint
	api.ClaManagerRequestsCreateCLAManagerRequestHandler = cla_manager_requests.CreateCLAManagerRequestHandlerFunc(func(params cla_manager_requests.CreateCLAManagerRequestParams, claUser *user.CLAUser) middleware.Responder {
		if !isValidUser(claUser) {
			return cla_manager_requests.NewCreateCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: "unauthorized",
				Code:    "401",
			})
		}

		existingRequests, getErr := service.GetRequests(params.CompanyID, params.ProjectID)
		if getErr != nil {
			msg := buildErrorMessage(params, getErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		companyModel, companyErr := companyService.GetCompany(params.CompanyID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessage(params, companyErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		projectModel, projectErr := projectService.GetProjectByID(params.ProjectID)
		if projectErr != nil || projectModel == nil {
			msg := buildErrorMessage(params, projectErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		userModel, userErr := usersService.GetUserByLFUserName(params.Body.UserLFID)
		if userErr != nil || userModel == nil {
			msg := buildErrorMessage(params, userErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Look up signature ACL to ensure the user can approve the request
		sigModels, sigErr := sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
			HTTPRequest: nil,
			CompanyID:   params.CompanyID,
			ProjectID:   params.ProjectID,
			NextKey:     nil,
			PageSize:    aws.Int64(5),
		})
		if sigErr != nil || sigModels == nil {
			msg := buildErrorMessage(params, sigErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: "CLA Manager Create Request - error reading CCA Signatures - " + msg,
				Code:    "400",
			})
		}
		if len(sigModels.Signatures) > 1 {
			log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s",
				params.CompanyID, params.ProjectID)
		}

		sigModel := sigModels.Signatures[0]
		claManagers := sigModel.SignatureACL

		var request *models.ClaManagerRequest
		// If no previous requests...
		if existingRequests == nil || existingRequests.Requests == nil {
			var createErr error
			request, createErr = service.CreateRequest(&CLAManagerRequest{
				CompanyID:         params.CompanyID,
				CompanyExternalID: companyModel.CompanyExternalID,
				CompanyName:       companyModel.CompanyName,
				ProjectID:         params.ProjectID,
				ProjectExternalID: projectModel.ProjectExternalID,
				ProjectName:       projectModel.ProjectName,
				UserID:            params.Body.UserLFID,
				UserExternalID:    userModel.UserExternalID,
				UserName:          params.Body.UserName,
				UserEmail:         params.Body.UserEmail,
				Status:            "pending",
			})
			if createErr != nil {
				msg := buildErrorMessage(params, createErr)
				log.Warn(msg)
				return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
			}

		} else {
			// Ok - we have an existing request with some state...

			// Check to see if we have an existing request in a pending or approved state - if so, don't allow
			for _, existingRequest := range existingRequests.Requests {
				if existingRequest.Status == "pending" || existingRequest.Status == "approved" {
					return cla_manager_requests.NewCreateCLAManagerRequestConflict().WithPayload(&models.ErrorResponse{
						Message: "an existing pending request exists for this user for this company and project",
						Code:    "409",
					})
				}
			}

			// Ok - existing state which is either denied or approved - allow them to create another request
			var updateErr error
			request, updateErr = service.PendingRequest(params.CompanyID, params.ProjectID, existingRequests.Requests[0].RequestID)
			if updateErr != nil {
				msg := buildErrorMessage(params, updateErr)
				log.Warn(msg)
				return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
			}
		}

		// Send an event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:         events.ClaManagerAccessRequestCreated,
			ProjectID:         params.ProjectID,
			ProjectModel:      projectModel,
			CompanyID:         params.CompanyID,
			CompanyModel:      companyModel,
			LfUsername:        params.Body.UserLFID,
			UserID:            params.Body.UserLFID,
			UserModel:         userModel,
			ExternalProjectID: projectModel.ProjectExternalID,
			EventData: &events.CLAManagerRequestCreatedEventData{
				RequestID:   request.RequestID,
				CompanyName: companyModel.CompanyName,
				ProjectName: projectModel.ProjectName,
				UserName:    params.Body.UserName,
				UserEmail:   params.Body.UserEmail,
				UserLFID:    params.Body.UserLFID,
			},
		})

		// Send email to each manager
		for _, manager := range claManagers {
			sendRequestAccessEmailToCLAManagers(companyModel, projectModel, params.Body.UserName, params.Body.UserEmail,
				manager.Username, manager.LfEmail, corporateConsoleURL)
		}

		return cla_manager_requests.NewCreateCLAManagerRequestOK().WithPayload(request)
	})

	// Get Requests
	api.ClaManagerRequestsGetCLAManagerRequestsHandler = cla_manager_requests.GetCLAManagerRequestsHandlerFunc(func(params cla_manager_requests.GetCLAManagerRequestsParams, claUser *user.CLAUser) middleware.Responder {
		if !isValidUser(claUser) {
			return cla_manager_requests.NewCreateCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: "unauthorized",
				Code:    "401",
			})
		}

		request, err := service.GetRequests(params.CompanyID, params.ProjectID)
		if err != nil {
			msg := buildErrorMessageForGetRequests(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewGetCLAManagerRequestsBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		return cla_manager_requests.NewGetCLAManagerRequestsOK().WithPayload(request)
	})

	// Get Request
	api.ClaManagerRequestsGetCLAManagerRequestHandler = cla_manager_requests.GetCLAManagerRequestHandlerFunc(func(params cla_manager_requests.GetCLAManagerRequestParams, claUser *user.CLAUser) middleware.Responder {
		if !isValidUser(claUser) {
			return cla_manager_requests.NewCreateCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: "unauthorized",
				Code:    "401",
			})
		}

		request, err := service.GetRequest(params.RequestID)
		if err != nil {
			msg := buildErrorMessageForGetRequest(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewGetCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// If we didn't find it
		if request == nil {
			return cla_manager_requests.NewGetCLAManagerRequestNotFound().WithPayload(&models.ErrorResponse{
				Message: fmt.Sprintf("request not found for Company ID: %s, Project ID: %s, Request ID: %s",
					params.CompanyID, params.ProjectID, params.RequestID),
				Code: "404",
			})
		}

		return cla_manager_requests.NewGetCLAManagerRequestOK().WithPayload(request)
	})

	// Approve Request
	api.ClaManagerRequestsApproveCLAManagerRequestHandler = cla_manager_requests.ApproveCLAManagerRequestHandlerFunc(func(params cla_manager_requests.ApproveCLAManagerRequestParams, claUser *user.CLAUser) middleware.Responder {

		companyModel, companyErr := companyService.GetCompany(params.CompanyID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessageForApprove(params, companyErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		projectModel, projectErr := projectService.GetProjectByID(params.ProjectID)
		if projectErr != nil || projectModel == nil {
			msg := buildErrorMessageForApprove(params, projectErr)
			log.Warn(msg)
			return cla_manager_requests.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Look up signature ACL to ensure the user can approve the request
		sigModels, sigErr := sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
			HTTPRequest: nil,
			CompanyID:   params.CompanyID,
			ProjectID:   params.ProjectID,
			NextKey:     nil,
			PageSize:    aws.Int64(5),
		})
		if sigErr != nil || sigModels == nil {
			msg := buildErrorMessageForApprove(params, sigErr)
			log.Warn(msg)
			return cla_manager_requests.NewApproveCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: "CLA Manager Approve Request - error reading CCA Signatures - " + msg,
				Code:    "400",
			})
		}
		if len(sigModels.Signatures) > 1 {
			log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s for request ID: %s",
				params.CompanyID, params.ProjectID, params.RequestID)
		}

		sigModel := sigModels.Signatures[0]
		claManagers := sigModel.SignatureACL
		if !currentUserInACL(claUser, claManagers) {
			return cla_manager_requests.NewApproveCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: fmt.Sprintf("CLA Manager %s / %s / %s not authorized to approve request for company ID: %s, project ID: %s",
					claUser.UserID, claUser.Name, claUser.LFEmail, params.CompanyID, params.ProjectID),
				Code: "401",
			})
		}

		// Approve the request
		request, err := service.ApproveRequest(params.CompanyID, params.ProjectID, params.RequestID)
		if err != nil {
			msg := buildErrorMessageForApprove(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewApproveCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Update the signature ACL
		_, aclErr := sigService.AddCLAManager(sigModel.SignatureID, request.UserID)
		if aclErr != nil {
			msg := buildErrorMessageForApprove(params, aclErr)
			log.Warn(msg)
			return cla_manager_requests.NewApproveCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Send an event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.ClaManagerAccessRequestApproved,
			ProjectID: params.ProjectID,
			CompanyID: params.CompanyID,
			UserID:    claUser.UserID,
			EventData: &events.CLAManagerRequestApprovedEventData{
				RequestID:    request.RequestID,
				CompanyName:  companyModel.CompanyName,
				ProjectName:  projectModel.ProjectName,
				UserName:     request.UserName,
				UserEmail:    request.UserEmail,
				ManagerName:  claUser.Name,    // from the request
				ManagerEmail: claUser.LFEmail, // from the request
			},
		})

		// Notify CLA Managers - send email to each manager
		for _, manager := range claManagers {
			sendRequestApprovedEmailToCLAManagers(companyModel, projectModel, request.UserName, request.UserEmail,
				manager.Username, manager.LfEmail)
		}

		// Notify the requester
		sendRequestApprovedEmailToRequester(companyModel, projectModel, request.UserName, request.UserEmail, corporateConsoleURL)

		return cla_manager_requests.NewCreateCLAManagerRequestOK().WithPayload(request)
	})

	// Deny Request
	api.ClaManagerRequestsDenyCLAManagerRequestHandler = cla_manager_requests.DenyCLAManagerRequestHandlerFunc(func(params cla_manager_requests.DenyCLAManagerRequestParams, claUser *user.CLAUser) middleware.Responder {

		companyModel, companyErr := companyService.GetCompany(params.CompanyID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessageForDeny(params, companyErr)
			log.Warn(msg)
			return cla_manager_requests.NewDenyCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		projectModel, projectErr := projectService.GetProjectByID(params.ProjectID)
		if projectErr != nil || projectModel == nil {
			msg := buildErrorMessageForDeny(params, projectErr)
			log.Warn(msg)
			return cla_manager_requests.NewDenyCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Look up signature ACL to ensure the user can deny the request
		sigModels, sigErr := sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
			HTTPRequest: nil,
			CompanyID:   params.CompanyID,
			ProjectID:   params.ProjectID,
			NextKey:     nil,
			PageSize:    aws.Int64(5),
		})
		if sigErr != nil || sigModels == nil {
			msg := buildErrorMessageForDeny(params, sigErr)
			log.Warn(msg)
			return cla_manager_requests.NewApproveCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: "CLA Manager Deny Request - error reading CCA Signatures - " + msg,
				Code:    "400",
			})
		}
		if len(sigModels.Signatures) > 1 {
			log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s for request ID: %s",
				params.CompanyID, params.ProjectID, params.RequestID)
		}

		sigModel := sigModels.Signatures[0]
		claManagers := sigModel.SignatureACL
		if !currentUserInACL(claUser, claManagers) {
			return cla_manager_requests.NewApproveCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: fmt.Sprintf("CLA Manager %s / %s / %s not authorized to approve request for company ID: %s, project ID: %s",
					claUser.UserID, claUser.Name, claUser.LFEmail, params.CompanyID, params.ProjectID),
				Code: "401",
			})
		}

		request, err := service.DenyRequest(params.CompanyID, params.ProjectID, params.RequestID)
		if err != nil {
			msg := buildErrorMessageForDeny(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewDenyCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Send an event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.ClaManagerAccessRequestDenied,
			ProjectID: params.ProjectID,
			CompanyID: params.CompanyID,
			UserID:    claUser.UserID,
			EventData: &events.CLAManagerRequestDeniedEventData{
				RequestID:    request.RequestID,
				CompanyName:  companyModel.CompanyName,
				ProjectName:  projectModel.ProjectName,
				UserName:     request.UserName,
				UserEmail:    request.UserEmail,
				ManagerName:  claUser.Name,    // from the request
				ManagerEmail: claUser.LFEmail, // from the request
			},
		})

		// Notify CLA Managers - send email to each manager
		for _, manager := range claManagers {
			sendRequestDeniedEmailToCLAManagers(companyModel, projectModel, request.UserName, request.UserEmail,
				manager.Username, manager.LfEmail)
		}

		// Notify the requester
		sendRequestDeniedEmailToRequester(companyModel, projectModel, request.UserName, request.UserEmail)

		return cla_manager_requests.NewCreateCLAManagerRequestOK().WithPayload(request)
	})

	// Delete Request
	api.ClaManagerRequestsDeleteCLAManagerRequestHandler = cla_manager_requests.DeleteCLAManagerRequestHandlerFunc(func(params cla_manager_requests.DeleteCLAManagerRequestParams, claUser *user.CLAUser) middleware.Responder {

		// Make sure the company id exists...
		companyModel, companyErr := companyService.GetCompany(params.CompanyID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessageForDelete(params, companyErr)
			return cla_manager_requests.NewDeleteCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Make sure the project id exists...
		projectModel, projectErr := projectService.GetProjectByID(params.ProjectID)
		if projectErr != nil || projectModel == nil {
			msg := buildErrorMessageForDelete(params, projectErr)
			log.Warn(msg)
			return cla_manager_requests.NewDenyCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Make sure the request exists...
		request, err := service.GetRequest(params.RequestID)
		if err != nil {
			msg := buildErrorMessageForDelete(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewDeleteCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		if request == nil {
			msg := buildErrorMessageForDelete(params, err)
			log.Warn(msg)
			return cla_manager_requests.NewDeleteCLAManagerRequestNotFound().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "404",
			})
		}

		// Look up signature ACL to ensure the user can delete the request
		sigModels, sigErr := sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
			HTTPRequest: nil,
			CompanyID:   params.CompanyID,
			ProjectID:   params.ProjectID,
			NextKey:     nil,
			PageSize:    aws.Int64(5),
		})
		if sigErr != nil || sigModels == nil {
			msg := buildErrorMessageForDelete(params, sigErr)
			log.Warn(msg)
			return cla_manager_requests.NewDeleteCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: "CLA Manager Delete Request - error reading CCA Signatures - " + msg,
				Code:    "400",
			})
		}
		if len(sigModels.Signatures) > 1 {
			log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s for request ID: %s",
				params.CompanyID, params.ProjectID, params.RequestID)
		}

		sigModel := sigModels.Signatures[0]
		claManagers := sigModel.SignatureACL
		if !currentUserInACL(claUser, claManagers) {
			return cla_manager_requests.NewDeleteCLAManagerRequestUnauthorized().WithPayload(&models.ErrorResponse{
				Message: fmt.Sprintf("CLA Manager %s / %s / %s not authorized to delete requests for company ID: %s, project ID: %s",
					claUser.UserID, claUser.Name, claUser.LFEmail, params.CompanyID, params.ProjectID),
				Code: "401",
			})
		}

		// Delete the request
		deleteErr := service.DeleteRequest(params.RequestID)
		if deleteErr != nil {
			msg := buildErrorMessageForDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager_requests.NewDeleteCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Send an event
		eventsService.LogEvent(&events.LogEventArgs{
			EventType: events.ClaManagerAccessRequestDeleted,
			ProjectID: params.ProjectID,
			CompanyID: params.CompanyID,
			UserID:    claUser.UserID,
			EventData: &events.CLAManagerRequestDeniedEventData{
				RequestID:    params.RequestID,
				CompanyName:  companyModel.CompanyName,
				ProjectName:  projectModel.ProjectName,
				UserName:     request.UserName,
				UserEmail:    request.UserEmail,
				ManagerName:  claUser.Name,    // from the request
				ManagerEmail: claUser.LFEmail, // from the request
			},
		})

		return cla_manager_requests.NewDeleteCLAManagerRequestOK()
	})
}

// currentUserInACL is a helper function to determine if the current logged in user is in the CLA Manager list
func currentUserInACL(currentUser *user.CLAUser, managers []models.User) bool {
	var inACL = false
	for _, manager := range managers {
		if manager.UserID == currentUser.UserID {
			inACL = true
			break
		}
	}

	return inACL
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(params cla_manager_requests.CreateCLAManagerRequestParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager Request using company ID: %s, project ID: %s, user ID: %s, user name: %s, user email: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.Body.UserLFID, params.Body.UserName, params.Body.UserEmail, err)
}

// buildErrorMessageForApprove is a helper function to build an error message
func buildErrorMessageForApprove(params cla_manager_requests.ApproveCLAManagerRequestParams, err error) string {
	return fmt.Sprintf("problem approving the CLA Manager Request using company ID: %s, project ID: %s, request ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.RequestID, err)
}

// buildErrorMessageForDeny is a helper function to build an error message
func buildErrorMessageForDeny(params cla_manager_requests.DenyCLAManagerRequestParams, err error) string {
	return fmt.Sprintf("problem denying the CLA Manager Request using company ID: %s, project ID: %s, request ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.RequestID, err)
}

// buildErrorMessageForDelete is a helper function to build an error message
func buildErrorMessageForDelete(params cla_manager_requests.DeleteCLAManagerRequestParams, err error) string {
	return fmt.Sprintf("problem deleting the CLA Manager Request using company ID: %s, project ID: %s, request ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.RequestID, err)
}

// buildErrorMessageForGetRequests is a helper function to build an error message
func buildErrorMessageForGetRequests(params cla_manager_requests.GetCLAManagerRequestsParams, err error) string {
	return fmt.Sprintf("problem fetching the CLA Manager Requests using company ID: %s, project ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, err)
}

// buildErrorMessageForGetRequest is a helper function to build an error message
func buildErrorMessageForGetRequest(params cla_manager_requests.GetCLAManagerRequestParams, err error) string {
	return fmt.Sprintf("problem fetching the CLA Manager Requests using company ID: %s, project ID: %s, request ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.RequestID, err)
}

// sendRequestAccessEmailToCLAManagers sends the request access email to the specified CLA Managers
func sendRequestAccessEmailToCLAManagers(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail, recipientName, recipientAddress, corporateConsoleURL string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: New CLA Manager Access Request for %s on %s", companyName, projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>You are currently listed as a CLA Manager from %s for the project %s. This means that you are able to maintain the
list of employees allowed to contribute to %s on behalf of your company, as well as view and manage the list of
your company’s CLA Managers for %s.</p>
<p>%s (%s) has requested to be added as another CLA Manager from %s for %s. This would permit them to maintain the
lists of approved contributors and CLA Managers as well.</p>
<p>If you want to permit this, please log into the EasyCLA Corporate Console at https://%s, select your company, then
select the %s project. From the CLA Manager requests, you can approve this user as an additional CLA Manager.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, projectName,
		companyName, projectName, projectName, projectName,
		requesterName, requesterEmail, companyName, projectName,
		corporateConsoleURL, projectName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendRequestApprovedEmailToCLAManagers(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: CLA Manager Access Approval Notice for %s", projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>The following user has been approved as a CLA Manager from %s for the project %s. This means that they can now
maintain the list of employees allowed to contribute to %s on behalf of your company, as well as view and manage the
list of company’s CLA Managers for %s.</p>
<ul>
<li>%s (%s)</li>
</ul>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, projectName,
		companyName, projectName, projectName, projectName,
		requesterName, requesterEmail)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendRequestApprovedEmailToRequester(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail, corporateConsoleURL string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: New CLA Manager Access Approved for %s", projectName)
	recipients := []string{requesterEmail}
	body := fmt.Sprintf(`
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>You have now been approved as a CLA Manager from %s for the project %s.  This means that you can now maintain the
list of employees allowed to contribute to %s on behalf of your company, as well as view and manage the list of your
company’s CLA Managers for %s.</p>
<p> To get started, please log into the EasyCLA Corporate Console at https://%s, and select your company and then the
project %s. From here you will be able to edit the list of approved employees and CLA Managers.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, requesterName, projectName,
		companyName, projectName, projectName, projectName,
		corporateConsoleURL, projectName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendRequestDeniedEmailToCLAManagers(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: CLA Manager Access Denied Notice for %s", projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>The following user has been denied as a CLA Manager from %s for the project %s. This means that they will not
be able to maintain the list of employees allowed to contribute to %s on behalf of your company.</p>
<ul>
<li>%s (%s)</li>
</ul>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, projectName,
		companyName, projectName, projectName,
		requesterName, requesterEmail)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendRequestDeniedEmailToRequester(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: New CLA Manager Access Denied for %s", projectName)
	recipients := []string{requesterEmail}
	body := fmt.Sprintf(`
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>You have been denied as a CLA Manager from %s for the project %s. This means that you can not maintain the
list of employees allowed to contribute to %s on behalf of your company.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, requesterName, projectName,
		companyName, projectName, projectName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
