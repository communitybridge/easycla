// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/utils"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// errors
var (
	ErrCclaWhitelistRequestAlreadyExists = errors.New("CCLA whiltelist request already exist")
)

// IService interface defines the service methods/functions
type IService interface {
	AddCclaWhitelistRequest(companyID string, projectID string, args models.CclaWhitelistRequestInput) (string, error)
	DeleteCclaWhitelistRequest(requestID string) error
	ListCclaWhitelistRequest(companyID string, projectID *string) (*models.CclaWhitelistRequestList, error)
	ListCclaWhitelistRequestByCompanyProjectUser(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error)

	sendRequestSentEmail(companyModel *models.Company, projectModel *models.Project, signature models.Signature, contributorName, contributorEmail, recipientName, recipientEmail, message string)
	sendRequestEmailToRecipient(companyModel *models.Company, projectModel *models.Project, contributorName, contributorEmail, recipientName, recipientAddress, message string)
}

type service struct {
	repo           Repository
	userRepo       users.UserRepository
	companyRepo    company.CompanyRepository
	projectRepo    project.ProjectRepository
	signatureRepo  signatures.SignatureRepository
	corpConsoleURL string
	httpClient     *http.Client
}

// NewService creates a new whitelist service
func NewService(repo Repository, userRepo users.UserRepository, companyRepo company.CompanyRepository, projectRepo project.ProjectRepository, signatureRepo signatures.SignatureRepository, corpConsoleURL string, httpClient *http.Client) IService {
	return service{
		repo:           repo,
		userRepo:       userRepo,
		companyRepo:    companyRepo,
		projectRepo:    projectRepo,
		signatureRepo:  signatureRepo,
		corpConsoleURL: corpConsoleURL,
		httpClient:     httpClient,
	}
}

func (s service) AddCclaWhitelistRequest(companyID string, projectID string, args models.CclaWhitelistRequestInput) (string, error) {
	list, err := s.ListCclaWhitelistRequestByCompanyProjectUser(companyID, &projectID, &args.ContributorID)
	if err != nil {
		log.Warnf("add CLA approve request - error looking up existing contributor invite requests for company: %s, project: %s, user by id: %s with name: %s, email: %s, error: %+v",
			companyID, projectID, args.ContributorID, args.ContributorName, args.ContributorEmail, err)
		return "", err
	}
	if len(list.List) > 0 {
		log.Warnf("add CLA approve request - found existing contributor invite - id: %s, request for company: %s, project: %s, user by id: %s with name: %s, email: %s",
			list.List[0].RequestID, companyID, projectID, args.ContributorID, args.ContributorName, args.ContributorEmail)
		return "", ErrCclaWhitelistRequestAlreadyExists
	}
	companyModel, err := s.companyRepo.GetCompany(companyID)
	if err != nil {
		log.Warnf("add CLA approve request - unable to lookup company by id: %s, error: %+v", companyID, err)
		return "", err
	}
	projectModel, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		log.Warnf("add CLA approve request - unable to lookup project by id: %s, error: %+v", projectID, err)
		return "", err
	}
	userModel, err := s.userRepo.GetUser(args.ContributorID)
	if err != nil {
		log.Warnf("add CLA approve request - unable to lookup user by id: %s with name: %s, email: %s, error: %+v",
			args.ContributorID, args.ContributorName, args.ContributorEmail, err)
		return "", err
	}
	if userModel == nil {
		log.Warnf("add CLA approve request - unable to lookup user by id: %s with name: %s, email: %s, error: user object not found",
			args.ContributorID, args.ContributorName, args.ContributorEmail)
		return "", errors.New("invalid user")
	}

	sig, sigErr := s.signatureRepo.GetProjectCompanySignatures(companyID, projectID, nil, 5)
	if sigErr != nil || sig == nil || sig.Signatures == nil {
		log.Warnf("add CLA approve request - unable to lookup signature by company id: %s project id: %s - (or no managers), sig: %+v, error: %+v",
			companyID, projectID, sig, err)
		return "", err
	}

	requestID, addErr := s.repo.AddCclaWhitelistRequest(companyModel, projectModel, userModel, args.ContributorName, args.ContributorEmail)
	if addErr != nil {
		log.Warnf("add CLA approve request - unable to add Approval Request for id: %s with name: %s, email: %s, error: %+v",
			args.ContributorID, args.ContributorName, args.ContributorEmail, addErr)
	}

	// Send the emails to the CLA managers for this CCLA Signature which includes the managers in the ACL list
	s.sendRequestSentEmail(companyModel, projectModel, sig.Signatures[0], args.ContributorName, args.ContributorEmail, args.RecipientName, args.RecipientEmail, args.Message)

	return requestID, nil
}

// DeleteCclaWhitelistRequest is the handler for the Delete CLA request
func (s service) DeleteCclaWhitelistRequest(requestID string) error {
	return s.repo.DeleteCclaWhitelistRequest(requestID)
}

// ListCclaWhitelistRequest is the handler for the list CLA request
func (s service) ListCclaWhitelistRequest(companyID string, projectID *string) (*models.CclaWhitelistRequestList, error) {
	return s.repo.ListCclaWhitelistRequest(companyID, projectID, nil)
}

// ListCclaWhitelistRequestByCompanyProjectUser is the handler for the list CLA request
func (s service) ListCclaWhitelistRequestByCompanyProjectUser(companyID string, projectID *string, userID *string) (*models.CclaWhitelistRequestList, error) {
	return s.repo.ListCclaWhitelistRequest(companyID, projectID, userID)
}

// sendRequestSentEmail sends emails to the CLA managers specified in the signature record
func (s service) sendRequestSentEmail(companyModel *models.Company, projectModel *models.Project, signature models.Signature, contributorName, contributorEmail, recipientName, recipientEmail, message string) {

	// If we have an override name and email from the request - possibly from the web form where the user selected the
	// CLA Manager Name/Email from a list, send this to this recipient (CLA Manager) - otherwise we will send to all
	// CLA Managers on the Signature ACL
	if recipientName != "" && recipientEmail != "" {
		s.sendRequestEmailToRecipient(companyModel, projectModel, contributorName, contributorEmail, recipientName, recipientEmail, message)
		return
	}

	// Send an email to each manager
	for _, manager := range signature.SignatureACL {

		// Need to determine which email...
		var whichEmail = ""
		if manager.LfEmail != "" {
			whichEmail = manager.LfEmail
		}

		// If no LF Email try to grab the first other email in their email list
		if manager.LfEmail == "" && manager.Emails != nil {
			whichEmail = manager.Emails[0]
		}
		if whichEmail == "" {
			log.Warnf("unable to send email to manager: %+v - no email on file...", manager)
		} else {
			// Send the email
			s.sendRequestEmailToRecipient(companyModel, projectModel, contributorName, contributorEmail, manager.Username, whichEmail, message)
		}
	}
}

// sendRequestEmailToRecipient generates and sends an email to the specified recipient
func (s service) sendRequestEmailToRecipient(companyModel *models.Company, projectModel *models.Project, contributorName, contributorEmail, recipientName, recipientAddress, message string) {
	companyName := companyModel.CompanyName
	projectName := projectModel.ProjectName

	var optionalMessage = ""
	if message != "" {
		optionalMessage = fmt.Sprintf("<p>%s included the following message in the request:</p>", contributorName)
		optionalMessage += fmt.Sprintf("<br/><p>%s</p><br/", message)
	}

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: Request to authorize %s for %s", contributorName, projectName)
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
<p>%s (%s) has requested to be added to the Allow List as an authorized contributor from
%s to the project %s. You are receiving this message as a CLA Manager from %s for
%s.</p>
%s
<p>If you want to add them to the Allow List, please
<a href="https://%s#/company/%s" target="_blank">log into the EasyCLA Corporate
Console</a>, where you can approve this user's request by selecting the 'Manage Approved List' and adding the
contributor's email, the contributor's entire email domain, their github ID or the entire GitHub Organization for the
repository. This will permit them to begin contributing to %s on behalf of %s.</p>
<p>If you are not certain whether to add them to the Allow List, please reach out to them directly to discuss.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/easycla/getting-started" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, projectName, contributorName, contributorEmail,
		companyName, projectName, companyName, projectName,
		optionalMessage, s.corpConsoleURL,
		companyModel.CompanyID, projectName, companyName)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
