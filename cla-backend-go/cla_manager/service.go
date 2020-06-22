// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	sigAPI "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

// IService interface defining the functions for the company service
type IService interface {
	CreateRequest(reqModel *CLAManagerRequest) (*models.ClaManagerRequest, error)
	GetRequests(companyID, projectID string) (*models.ClaManagerRequestList, error)
	GetRequestsByUserID(companyID, projectID, userID string) (*models.ClaManagerRequestList, error)
	GetRequest(requestID string) (*models.ClaManagerRequest, error)

	ApproveRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error)
	DenyRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error)
	PendingRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error)
	DeleteRequest(requestID string) error

	AddClaManager(companyID string, projectID string, LFID string) (*models.Signature, error)
	RemoveClaManager(companyID string, projectID string, LFID string) (*models.Signature, error)
}

type service struct {
	repo                IRepository
	companyService      company.IService
	projectService      project.Service
	usersService        users.Service
	sigService          signatures.SignatureService
	eventsService       events.Service
	corporateConsoleURL string
}

// NewService creates a new service object
func NewService(repo IRepository, companyService company.IService, projectService project.Service, usersService users.Service, sigService signatures.SignatureService, eventsService events.Service, corporateConsoleURL string) IService {
	return service{
		repo:                repo,
		companyService:      companyService,
		projectService:      projectService,
		usersService:        usersService,
		sigService:          sigService,
		eventsService:       eventsService,
		corporateConsoleURL: corporateConsoleURL,
	}
}

// CreateRequest creates a request based on the specified parameters
func (s service) CreateRequest(reqModel *CLAManagerRequest) (*models.ClaManagerRequest, error) {
	request, err := s.repo.CreateRequest(reqModel)
	if err != nil {
		log.Warnf("problem with approving request for company ID: %s, project ID: %s, user ID: %s, user name: %s, error :%+v",
			reqModel.CompanyID, reqModel.ProjectID, reqModel.UserID, reqModel.UserName, err)
		return nil, err
	}

	respModel := dbModelToServiceModel(*request)

	return &respModel, err
}

// GetRequests returns a requests object based on the specified parameters
func (s service) GetRequests(companyID, projectID string) (*models.ClaManagerRequestList, error) {
	requests, err := s.repo.GetRequests(companyID, projectID)
	if err != nil {
		log.Warnf("problem with fetching request for company ID: %s, project ID: %s, error :%+v",
			companyID, projectID, err)
		return nil, err
	}

	// Convert to a service response model
	responseModel := models.ClaManagerRequestList{}
	for _, request := range requests.Requests {
		responseModel.Requests = append(responseModel.Requests, dbModelToServiceModel(request))
	}

	return &responseModel, nil
}

// GetRequestsByUserID returns a requests object based on the specified parameters
func (s service) GetRequestsByUserID(companyID, projectID, userID string) (*models.ClaManagerRequestList, error) {
	requests, err := s.repo.GetRequestsByUserID(companyID, projectID, userID)
	if err != nil {
		log.Warnf("problem with fetching request for company ID: %s, project ID: %s, user ID: %s, error :%+v",
			companyID, projectID, userID, err)
		return nil, err
	}

	// Convert to a service response model
	responseModel := models.ClaManagerRequestList{}
	for _, request := range requests.Requests {
		responseModel.Requests = append(responseModel.Requests, dbModelToServiceModel(request))
	}

	return &responseModel, nil
}

// GetRequest returns the request object based on the specified parameters
func (s service) GetRequest(requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.GetRequest(requestID)
	if err != nil {
		log.Warnf("problem with fetching request for request ID: %s, error :%+v",
			requestID, err)
		return nil, err
	}

	if request == nil {
		log.Debugf("request not found for Request ID: %s", requestID)
		return nil, nil
	}

	respModel := dbModelToServiceModel(*request)

	return &respModel, err
}

// ApproveRequest approves the request based on the specified parameters
func (s service) ApproveRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.ApproveRequest(companyID, projectID, requestID)
	if err != nil {
		log.Warnf("problem with approving request for company ID: %s, project ID: %s, request ID: %s, error :%+v",
			companyID, projectID, requestID, err)
		return nil, err
	}

	respModel := dbModelToServiceModel(*request)

	return &respModel, err
}

// PendingRequest updates the specified request to the pending state
func (s service) PendingRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.PendingRequest(companyID, projectID, requestID)
	if err != nil {
		log.Warnf("problem with setting the pending status for company ID: %s, project ID: %s, request ID: %s, error :%+v",
			companyID, projectID, requestID, err)
		return nil, err
	}

	respModel := dbModelToServiceModel(*request)

	return &respModel, err
}

// DenyRequest denies the request based on the specified parameters
func (s service) DenyRequest(companyID, projectID, requestID string) (*models.ClaManagerRequest, error) {
	request, err := s.repo.DenyRequest(companyID, projectID, requestID)
	if err != nil {
		log.Warnf("problem with denying request for company ID: %s, project ID: %s, request ID: %s, error :%+v",
			companyID, projectID, requestID, err)
		return nil, err
	}

	respModel := dbModelToServiceModel(*request)

	return &respModel, err
}

// DeleteRequest deletes the request based on the specified parameters
func (s service) DeleteRequest(requestID string) error {
	err := s.repo.DeleteRequest(requestID)
	if err != nil {
		log.Warnf("problem deleting request ID: %s, error :%+v",
			requestID, err)
		return err
	}
	return nil
}

// AddClaManager Adds LFID to Signature Access Control List list
func (s service) AddClaManager(companyID string, projectID string, LFID string) (*models.Signature, error) {

	userModel, userErr := s.usersService.GetUserByLFUserName(LFID)
	if userErr != nil || userModel == nil {
		return nil, userErr
	}
	companyModel, companyErr := s.companyService.GetCompany(companyID)
	if companyErr != nil || companyModel == nil {
		return nil, companyErr
	}

	projectModel, projectErr := s.projectService.GetProjectByID(projectID)
	if projectErr != nil || projectModel == nil {
		return nil, projectErr
	}

	// Look up signature ACL to ensure the user can add cla manager
	sigModels, sigErr := s.sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
		HTTPRequest: nil,
		CompanyID:   companyID,
		ProjectID:   projectID,
		NextKey:     nil,
		PageSize:    aws.Int64(5),
	})
	if sigErr != nil || sigModels == nil {
		return nil, sigErr
	}

	if len(sigModels.Signatures) > 1 {
		log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s",
			companyID, projectID)
	}

	sigModel := sigModels.Signatures[0]
	claManagers := sigModel.SignatureACL

	log.Debugf("Got Company signatures - Company: %s , Project: %s , signatureID: %s ",
		companyID, projectID, sigModel.SignatureID)

	// Update the signature ACL
	addedSignature, aclErr := s.sigService.AddCLAManager(sigModel.SignatureID, LFID)
	if aclErr != nil {
		return nil, aclErr
	}

	// Notify CLA Managers - send email to each manager
	for _, manager := range claManagers {
		sendClaManagerAddedEmailToCLAManagers(companyModel, projectModel, userModel.Username, userModel.LfEmail,
			manager.Username, manager.LfEmail)
	}
	// Notify the added user
	sendClaManagerAddedEmailToUser(companyModel, projectModel, userModel.Username, userModel.LfEmail, s.corporateConsoleURL)

	// Send an event
	s.eventsService.LogEvent(&events.LogEventArgs{
		EventType:         events.ClaManagerCreated,
		ProjectID:         projectID,
		ProjectModel:      projectModel,
		CompanyID:         companyID,
		CompanyModel:      companyModel,
		LfUsername:        LFID,
		UserID:            LFID,
		UserModel:         userModel,
		ExternalProjectID: projectModel.ProjectExternalID,
		EventData: &events.CLAManagerCreatedEventData{
			UserName:  userModel.Username,
			UserEmail: userModel.LfEmail,
			UserLFID:  userModel.LfUsername,
		},
	})

	return addedSignature, nil
}

// Utility function that returns company signature
func (s service) getCompanySignature(companyID string, projectID string) (*models.Signature, error) {
	// Look up signature ACL to ensure the user can remove given cla manager
	sigModels, sigErr := s.sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
		HTTPRequest: nil,
		CompanyID:   companyID,
		ProjectID:   projectID,
		NextKey:     nil,
		PageSize:    aws.Int64(5),
	})
	if sigErr != nil || sigModels == nil {
		log.Warnf("Unable to lookup project company signature using Project ID: %s, Company ID: %s, error: %+v",
			projectID, companyID, sigErr)
		return nil, sigErr
	}

	if len(sigModels.Signatures) > 1 {
		log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s",
			companyID, projectID)
	}
	return sigModels.Signatures[0], nil
}

// RemoveClaManager removes lfid from signature acl with given company and project
func (s service) RemoveClaManager(companyID string, projectID string, LFID string) (*models.Signature, error) {

	userModel, userErr := s.usersService.GetUserByLFUserName(LFID)
	if userErr != nil || userModel == nil {
		return nil, userErr
	}
	companyModel, companyErr := s.companyService.GetCompany(companyID)
	if companyErr != nil || companyModel == nil {
		return nil, companyErr
	}

	projectModel, projectErr := s.projectService.GetProjectByID(projectID)
	if projectErr != nil || projectModel == nil {
		return nil, projectErr
	}

	// Look up signature ACL to ensure the user can remove given cla manager
	sigModels, sigErr := s.sigService.GetProjectCompanySignatures(sigAPI.GetProjectCompanySignaturesParams{
		HTTPRequest: nil,
		CompanyID:   companyID,
		ProjectID:   projectID,
		NextKey:     nil,
		PageSize:    aws.Int64(5),
	})
	if sigErr != nil || sigModels == nil {
		log.Warnf("Unable to lookup project company signature using Project ID: %s, Company ID: %s, error: %+v",
			projectID, companyID, sigErr)
		return nil, sigErr
	}

	if len(sigModels.Signatures) > 1 {
		log.Warnf("returned multiple CCLA signature models for company ID: %s, project ID: %s",
			companyID, projectID)
	}

	sigModel, sigErr := s.getCompanySignature(companyID, projectID)
	if sigErr != nil {
		return nil, sigErr
	}

	// Update the signature ACL
	updatedSignature, aclErr := s.sigService.RemoveCLAManager(sigModel.SignatureID, LFID)
	if aclErr != nil || updatedSignature == nil {
		log.Warnf("remove CLA Manager returned an error or empty signature model using Signature ID: %s, error: %+v",
			sigModel.SignatureID, sigErr)
		return nil, aclErr
	}

	// Get Updated cla manager list with removed manager for email purporses
	sigModel, sigErr = s.getCompanySignature(companyID, projectID)
	if sigErr != nil {
		return nil, sigErr
	}
	claManagers := sigModel.SignatureACL
	// Notify CLA Managers - send email to each manager
	for _, manager := range claManagers {
		sendClaManagerDeleteEmailToCLAManagers(companyModel, projectModel, userModel.LfUsername,
			manager.Username, manager.LfEmail)
	}

	// Notify the removed manager
	sendRemovedClaManagerEmailToRecipient(companyModel, projectModel, userModel.LfUsername, userModel.LfEmail, claManagers)

	// Send an event
	s.eventsService.LogEvent(&events.LogEventArgs{
		EventType:         events.ClaManagerDeleted,
		ProjectID:         projectID,
		ProjectModel:      projectModel,
		CompanyID:         companyID,
		CompanyModel:      companyModel,
		LfUsername:        userModel.LfUsername,
		UserID:            LFID,
		UserModel:         userModel,
		ExternalProjectID: projectModel.ProjectExternalID,
		EventData: &events.CLAManagerDeletedEventData{
			CompanyName: companyModel.CompanyName,
			ProjectName: projectModel.ProjectName,
			UserName:    userModel.LfUsername,
			UserEmail:   userModel.LfEmail,
			UserLFID:    LFID,
		},
	})

	return updatedSignature, nil
}

func sendClaManagerAddedEmailToUser(companyModel *models.Company, projectModel *models.Project, requesterName, requesterEmail, corporateConsoleURL string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: Added as CLA Manager for Project :%s", projectName)
	recipients := []string{requesterEmail}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>You have been added as a CLA Manager from %s for the project %s.  This means that you can now maintain the
list of employees allowed to contribute to %s on behalf of your company, as well as view and manage the list of your
company’s CLA Managers for %s.</p>
<p> To get started, please log into the EasyCLA Corporate Console at https://%s, and select your company and then the
project %s. From here you will be able to edit the list of approved employees and CLA Managers.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>`,
		requesterName, projectName,
		companyName, projectName, projectName, projectName,
		corporateConsoleURL, projectName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendClaManagerAddedEmailToCLAManagers(companyModel *models.Company, projectModel *models.Project, name, email, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: CLA Manager Added Notice for %s", projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>The following user has been added as a CLA Manager from %s for the project %s. This means that they can now
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
<p>EasyCLA support team</p>`,
		recipientName, projectName,
		companyName, projectName, projectName, projectName,
		name, email)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// sendRequestRejectedEmailToRecipient generates and sends an email to the specified recipient
func sendRemovedClaManagerEmailToRecipient(companyModel *models.Company, projectModel *models.Project, recipientName, recipientAddress string, claManagers []models.User) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	var companyManagerText = ""
	companyManagerText += "<ul>"

	// Build a fancy text string with CLA Manager name <email> as an HTML unordered list
	for _, companyAdmin := range claManagers {

		// Need to determine which email...
		var whichEmail = ""
		if companyAdmin.LfEmail != "" {
			whichEmail = companyAdmin.LfEmail
		}

		// If no LF Email try to grab the first other email in their email list
		if companyAdmin.LfEmail == "" && companyAdmin.Emails != nil {
			whichEmail = companyAdmin.Emails[0]
		}

		if whichEmail == "" {
			log.Warnf("unable to send email to manager: %+v - no email on file...", companyAdmin)
		} else {
			companyManagerText += fmt.Sprintf("<li>%s <%s></li>", companyAdmin.LfUsername, whichEmail)
		}
	}

	companyManagerText += "</ul>"

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: Removed as CLA Manager for Project %s", projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>You have been removed as a CLA Manager from %s for the project %s.</p>
<p>If you have further questions about this, please contact one of the existing managers from
%s:</p>
%s
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>`,
		recipientName, projectName, companyName, projectName, companyName, companyManagerText)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendClaManagerDeleteEmailToCLAManagers(companyModel *models.Company, projectModel *models.Project, name, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: CLA Manager Removed Notice for %s", projectName)
	recipients := []string{recipientAddress}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the project %s.</p>
<p>%s has been removed as a CLA Manager from %s for the project %s.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>`, recipientName, projectName, name, companyName, projectName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
