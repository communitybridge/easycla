// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/users"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	organization_service "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
)

type service struct {
	repo                IRepository
	userDynamoRepo      user.RepositoryService
	corporateConsoleURL string
	userService         users.Service
}

const (
	// StatusPending indicates the invitation status is pending
	StatusPending = "pending"
)

// IService interface defining the functions for the company service
type IService interface { // nolint
	GetCompanies() (*models.Companies, error)
	GetCompany(companyID string) (*models.Company, error)
	GetCompanyByExternalID(companySFID string) (*models.Company, error)
	SearchCompanyByName(companyName string, nextKey string) (*models.Companies, error)
	GetCompaniesByUserManager(userID string) (*models.Companies, error)
	GetCompaniesByUserManagerWithInvites(userID string) (*models.CompaniesWithInvites, error)

	AddUserToCompanyAccessList(companyID, lfid string) error
	GetCompanyInviteRequests(companyID string, status *string) ([]models.CompanyInviteUser, error)
	GetCompanyUserInviteRequests(companyID string, userID string) (*models.CompanyInviteUser, error)
	AddPendingCompanyInviteRequest(companyID string, userID string) (*InviteModel, error)
	ApproveCompanyAccessRequest(companyInviteID string) (*InviteModel, error)
	RejectCompanyAccessRequest(companyInviteID string) (*InviteModel, error)

	// calls org service
	SearchOrganizationByName(orgName string) (*models.OrgList, error)

	sendRequestAccessEmail(companyModel *models.Company, requesterName, requesterEmail, recipientName, recipientAddress string)
	sendRequestApprovedEmailToRecipient(companyModel *models.Company, recipientName, recipientAddress string)
	sendRequestRejectedEmailToRecipient(companyModel *models.Company, recipientName, recipientAddress string)
	getPreferredNameAndEmail(lfid string) (string, string, error)
}

// NewService creates a new company service object
func NewService(repo IRepository, corporateConsoleURL string, userDynamoRepo user.RepositoryService, userService users.Service) IService {
	return service{
		repo:                repo,
		userDynamoRepo:      userDynamoRepo,
		corporateConsoleURL: corporateConsoleURL,
		userService:         userService,
	}
}

// GetCompanies returns all the companies
func (s service) GetCompanies() (*models.Companies, error) {
	return s.repo.GetCompanies()
}

// GetCompany returns the company associated with the company ID
func (s service) GetCompany(companyID string) (*models.Company, error) {
	return s.repo.GetCompany(companyID)
}

// SearchCompanyByName locates companies by the matching name and return any potential matches
func (s service) SearchCompanyByName(companyName string, nextKey string) (*models.Companies, error) {
	companies, err := s.repo.SearchCompanyByName(companyName, nextKey)
	if err != nil {
		log.Warnf("Error searching company by company name: %s, error: %v", companyName, err)
		return nil, err
	}

	return companies, nil
}

// GetCompanyUserManager the get a list of companies when provided the company id and user manager
func (s service) GetCompaniesByUserManager(userID string) (*models.Companies, error) {
	userModel, err := s.userDynamoRepo.GetUser(userID)
	if err != nil {
		log.Warnf("Unable to lookup user by user id: %s, error: %v", userID, err)
		return nil, err
	}

	return s.repo.GetCompaniesByUserManager(userID, userModel)
}

// GetCompanyUserManagerWithInvites the get a list of companies including status when provided the company id and user manager
func (s service) GetCompaniesByUserManagerWithInvites(userID string) (*models.CompaniesWithInvites, error) {
	userModel, err := s.userDynamoRepo.GetUser(userID)
	if err != nil {
		log.Warnf("Unable to lookup user by user id: %s, error: %v", userID, err)
		return nil, err
	}

	return s.repo.GetCompaniesByUserManagerWithInvites(userID, userModel)
}

// GetCompanyInviteRequests returns a list of company invites when provided the company ID
func (s service) GetCompanyInviteRequests(companyID string, status *string) ([]models.CompanyInviteUser, error) {
	companyInvites, err := s.repo.GetCompanyInviteRequests(companyID, status)
	if err != nil {
		return nil, err
	}

	var users []models.CompanyInviteUser
	for _, invite := range companyInvites {

		dbUserModel, err := s.userDynamoRepo.GetUser(invite.UserID)
		if err != nil {
			log.Warnf("Error fetching user with userID: %s, error: %v", invite.UserID, err)
			continue
		}

		// Default status is pending if there's a record but no status
		if invite.Status == "" {
			invite.Status = StatusPending
		}

		users = append(users, models.CompanyInviteUser{
			InviteID:  invite.CompanyInviteID,
			UserName:  dbUserModel.UserName,
			UserEmail: dbUserModel.LFEmail,
			UserLFID:  dbUserModel.LFUsername,
			Status:    invite.Status,
		})
	}

	return users, nil
}

// GetCompanyUserInviteRequests returns a list of company invites when provided the company ID
func (s service) GetCompanyUserInviteRequests(companyID string, userID string) (*models.CompanyInviteUser, error) {
	invite, err := s.repo.GetCompanyUserInviteRequests(companyID, userID)
	if err != nil {
		return nil, err
	}

	if invite == nil {
		return nil, nil
	}

	//var users []models.CompanyInviteUser

	dbUserModel, err := s.userDynamoRepo.GetUser(invite.UserID)
	if err != nil {
		log.Warnf("Error fetching company invite user with company id: %s and user id: %s, error: %v",
			companyID, userID, err)
		return nil, err
	}

	// Default status is pending if there's a record but no status
	if invite.Status == "" {
		invite.Status = StatusPending
	}

	// Let's do a company lookup so we can grab the company name
	company, err := s.repo.GetCompany(companyID)
	if err != nil {
		log.Warnf("Error fetching company with company id: %s, error: %v", companyID, err)
		return nil, err
	}

	return &models.CompanyInviteUser{
		InviteID:    invite.CompanyInviteID,
		UserName:    dbUserModel.UserName,
		UserEmail:   dbUserModel.LFEmail,
		UserLFID:    dbUserModel.LFUsername,
		Status:      invite.Status,
		CompanyName: company.CompanyName,
	}, nil
}

// AddPendingCompanyInviteRequest adds a new company invite request
func (s service) AddPendingCompanyInviteRequest(companyID string, userID string) (*InviteModel, error) {

	newInvite, err := s.repo.AddPendingCompanyInviteRequest(companyID, userID)
	if err != nil {
		return nil, err
	}

	companyModel, companyErr := s.GetCompany(companyID)
	if companyErr != nil {
		log.Warnf("AddPendingCompanyInviteRequest - unable to locate company model by ID: %s, error: %+v",
			companyID, companyErr)
		return nil, companyErr
	}

	userModel, userErr := s.userDynamoRepo.GetUser(userID)
	if userErr != nil {
		log.Warnf("AddPendingCompanyInviteRequest - unable to locate user model by ID: %s, error: %+v",
			userID, userErr)
		return nil, userErr
	}

	// Need to determine which email...
	var requesterEmail = ""
	if userModel.LFEmail != "" {
		requesterEmail = userModel.LFEmail
	}

	// If no LF Email try to grab the first other email in their email list
	if userModel.LFEmail == "" && userModel.UserEmails != nil {
		requesterEmail = userModel.UserEmails[0]
	}

	// Send the email to each company manager
	for _, companyManagerLFID := range companyModel.CompanyACL {
		companyManagerName, companyManagerEmail, err := s.getPreferredNameAndEmail(companyManagerLFID)
		if err != nil {
			log.Warnf("unable to lookup company manager's name and email using LFID: %s - unable to send email, error: %+v",
				companyManagerLFID, err)
			continue
		}

		// Send an email to this company manager
		s.sendRequestAccessEmail(companyModel, userModel.UserName, requesterEmail, companyManagerName, companyManagerEmail)
	}

	return &InviteModel{
		CompanyInviteID:    newInvite.CompanyInviteID,
		RequestedCompanyID: newInvite.RequestedCompanyID,
		CompanyName:        companyModel.CompanyName,
		UserName:           userModel.UserName,
		UserEmail:          userModel.LFEmail,
		UserID:             newInvite.UserID,
		Status:             newInvite.Status,
		Created:            newInvite.Created,
		Updated:            newInvite.Updated,
	}, nil
}

// ApproveCompanyAccessRequest approve access request service method
func (s service) ApproveCompanyAccessRequest(companyInviteID string) (*InviteModel, error) {
	err := s.repo.ApproveCompanyAccessRequest(companyInviteID)
	if err != nil {
		return nil, err
	}

	inviteModel, inviteErr := s.repo.GetCompanyInviteRequest(companyInviteID)
	if inviteErr != nil || inviteModel == nil {
		log.Warnf("ApproveCompanyAccessRequest - unable to locate company invite: %s, error: %+v",
			companyInviteID, inviteErr)
		return nil, inviteErr
	}

	companyModel, companyErr := s.GetCompany(inviteModel.RequestedCompanyID)
	if companyErr != nil {
		log.Warnf("ApproveCompanyAccessRequest - unable to locate company model by ID: %s, error: %+v",
			inviteModel.RequestedCompanyID, companyErr)
		return nil, companyErr
	}

	userModel, userErr := s.userDynamoRepo.GetUser(inviteModel.UserID)
	if userErr != nil {
		log.Warnf("ApproveCompanyAccessRequest - unable to locate user model by ID: %s, error: %+v",
			inviteModel.UserID, userErr)
		return nil, userErr
	}

	updatedUserModel, userUpdateErr := s.userDynamoRepo.SetCompanyID(userModel.UserID, companyModel.CompanyID)
	if userUpdateErr != nil {
		log.Warnf("ApproveCompanyAccessRequest - unable to update user model by ID: %s, with company ID: %s error: %+v",
			inviteModel.UserID, companyModel.CompanyID, userUpdateErr)
		return nil, userUpdateErr
	}

	// update the company ACL
	aclErr := s.AddUserToCompanyAccessList(inviteModel.RequestedCompanyID, updatedUserModel.LFUsername)
	if aclErr != nil {
		log.Warnf("ApproveCompanyAccessRequest - unable to add user to Company ACL, company ID: %s, user LFID: %s, error: %+v",
			inviteModel.RequestedCompanyID, updatedUserModel.UserName, err)
		return nil, aclErr
	}

	// Need to determine which email...
	var whichEmail = ""
	if updatedUserModel.LFEmail != "" {
		whichEmail = userModel.LFEmail
	}

	// If no LF Email try to grab the first other email in their email list
	if updatedUserModel.LFEmail == "" && updatedUserModel.UserEmails != nil {
		whichEmail = updatedUserModel.UserEmails[0]
	}

	s.sendRequestApprovedEmailToRecipient(companyModel, updatedUserModel.UserName, whichEmail)

	return &InviteModel{
		CompanyInviteID:    inviteModel.CompanyInviteID,
		RequestedCompanyID: inviteModel.RequestedCompanyID,
		CompanyName:        companyModel.CompanyName,
		UserName:           updatedUserModel.UserName,
		UserEmail:          updatedUserModel.LFEmail,
		UserID:             inviteModel.UserID,
		Status:             inviteModel.Status,
		Created:            inviteModel.Created,
		Updated:            inviteModel.Updated,
	}, nil
}

// RejectCompanyAccessRequest approve access request service method
func (s service) RejectCompanyAccessRequest(companyInviteID string) (*InviteModel, error) {
	err := s.repo.RejectCompanyAccessRequest(companyInviteID)
	if err != nil {
		return nil, err
	}

	inviteModel, inviteErr := s.repo.GetCompanyInviteRequest(companyInviteID)
	if inviteErr != nil || inviteModel == nil {
		log.Warnf("RejectCompanyAccessRequest - unable to locate company invite: %s, error: %+v",
			companyInviteID, inviteErr)
		return nil, inviteErr
	}

	companyModel, companyErr := s.GetCompany(inviteModel.RequestedCompanyID)
	if companyErr != nil {
		log.Warnf("RejectCompanyAccessRequest - unable to locate company model by ID: %s, error: %+v",
			inviteModel.RequestedCompanyID, companyErr)
		return nil, companyErr
	}

	userModel, userErr := s.userDynamoRepo.GetUser(inviteModel.UserID)
	if userErr != nil {
		log.Warnf("RejectCompanyAccessRequest - unable to locate user model by ID: %s, error: %+v",
			inviteModel.UserID, userErr)
		return nil, userErr
	}

	// Need to determine which email...
	var whichEmail = ""
	if userModel.LFEmail != "" {
		whichEmail = userModel.LFEmail
	}

	// If no LF Email try to grab the first other email in their email list
	if userModel.LFEmail == "" && userModel.UserEmails != nil {
		whichEmail = userModel.UserEmails[0]
	}

	s.sendRequestRejectedEmailToRecipient(companyModel, userModel.UserName, whichEmail)

	return &InviteModel{
		CompanyInviteID:    inviteModel.CompanyInviteID,
		RequestedCompanyID: inviteModel.RequestedCompanyID,
		CompanyName:        companyModel.CompanyName,
		UserName:           userModel.UserName,
		UserEmail:          userModel.LFEmail,
		UserID:             inviteModel.UserID,
		Status:             inviteModel.Status,
		Created:            inviteModel.Created,
		Updated:            inviteModel.Updated,
	}, nil
}

// AddUserToCompanyAccessList adds a user to the specified company
func (s service) AddUserToCompanyAccessList(companyID, lfid string) error {
	// call the get company function
	company, err := s.repo.GetCompany(companyID)
	if err != nil {
		log.Warnf("Error retrieving company by company ID: %s, error: %v", companyID, err)
		return err
	}

	// perform ACL check
	// check if user already exists in the company acl
	for _, acl := range company.CompanyACL {
		if acl == lfid {
			log.Warnf(fmt.Sprintf("User %s has already been added to the company acl", lfid))
			return nil
		}
	}
	// add user to string set
	company.CompanyACL = append(company.CompanyACL, lfid)

	err = s.repo.UpdateCompanyAccessList(companyID, company.CompanyACL)
	if err != nil {
		log.Warnf("Error updating company access list with company ID: %s, company ACL: %v, error: %v", companyID, company.CompanyACL, err)
		return err
	}

	return nil
}

// sendRequestAccessEmail sends the request access email
func (s service) sendRequestAccessEmail(companyModel *models.Company, requesterName, requesterEmail, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName

	requestedUserInfo := fmt.Sprintf("<ul><li>%s (%s)</li></ul>", requesterName, requesterEmail)

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: New Company Manager Access Request for %s", companyName)
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
<p>This is a notification email from EasyCLA regarding the company %s.</p>
<p>The following user has requested to join %s as a Company Manager. 
By approving this request the user could view and apply for CLA Manager
status on projects associated with your company. </p>
%s
<p>To get started, please log into the EasyCLA Corporate Console at 
https://%s, and select your company. From there you will
be able to view the list of projects which have EasyCLA configured and apply
for CLA Manager status.
</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, companyName, companyName, requestedUserInfo, s.corporateConsoleURL)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// sendRequestApprovedEmailToRecipient generates and sends an email to the specified recipient
func (s service) sendRequestApprovedEmailToRecipient(companyModel *models.Company, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: Company Manager Access Approved for %s", companyName)
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
<p>This is a notification email from EasyCLA regarding the company %s.</p>
<p>You have now been approved as a Company Manager for %s.
This means that you can now view and apply for CLA Manager status on
projects associated with your company.
</p>
<p>To get started, please log into the EasyCLA Corporate Console at 
https://%s, and select your company. From there you will
be able to view the list of projects which have EasyCLA configured and apply
for CLA Manager status.
</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, companyName, companyName, s.corporateConsoleURL)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// sendRequestRejectedEmailToRecipient generates and sends an email to the specified recipient
func (s service) sendRequestRejectedEmailToRecipient(companyModel *models.Company, recipientName, recipientAddress string) {
	companyName := companyModel.CompanyName

	var companyManagerText = ""
	companyManagerText += "<ul>"

	// Build a fancy text string with CLA Manager name <email> as an HTML unordered list
	for _, companyAdminLFID := range companyModel.CompanyACL {

		userModel, userErr := s.userDynamoRepo.GetUserAndProfilesByLFID(companyAdminLFID)
		if userErr != nil {
			log.Warnf("RejectCompanyAccessRequest - unable to locate user model by ID: %s, error: %+v",
				companyAdminLFID, userErr)
		}

		// Need to determine which email...
		var whichEmail = ""
		if userModel.LFEmail != "" {
			whichEmail = userModel.LFEmail
		}

		// If no LF Email try to grab the first other email in their email list
		if userModel.LFEmail == "" && userModel.Emails != nil {
			whichEmail = userModel.Emails[0]
		}

		if whichEmail == "" {
			log.Warnf("unable to send email to manager: %+v - no email on file...", userModel)
		} else {
			companyManagerText += fmt.Sprintf("<li>%s <%s></li>", userModel.Name, whichEmail)
		}
	}

	companyManagerText += "</ul>"

	// subject string, body string, recipients []string
	subject := fmt.Sprintf("EasyCLA: CLA Manager Access Denied for %s", companyName)
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
<p>This is a notification email from EasyCLA regarding the company %s.</p>
<p>Your request to become a Company Manager was denied by one of the existing Company Managers.
If you have further questions about this denial, please contact one of the existing managers from
%s:</p>
%s
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, recipientName, companyName, companyName, companyManagerText)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// getPreferredNameAndEmail when given the user LFID, this routine returns the user's name and preferred email
func (s service) getPreferredNameAndEmail(lfid string) (string, string, error) {
	userModel, userErr := s.userService.GetUserByLFUserName(lfid)
	if userErr != nil {
		log.Warnf("getPreferredNameAndEmail - unable to locate user model by ID: %s, error: %+v",
			lfid, userErr)
		return "", "", userErr
	}

	userName := userModel.Username
	// If no user name specified - then use the user's LF user name I guess
	if userName == "" {
		userName = userModel.LfUsername
	}

	userEmail := userModel.LfEmail
	if userEmail == "" && userModel.Emails != nil && len(userModel.Emails) > 0 {
		userEmail = userModel.Emails[0]
	}

	return userName, userEmail, nil
}

func (s service) GetCompanyByExternalID(companySFID string) (*models.Company, error) {
	comp, err := s.repo.GetCompanyByExternalID(companySFID)
	if err == nil {
		return comp, nil
	}
	if err == ErrCompanyDoesNotExist {
		comp, err = s.createOrgFromExternalID(companySFID)
		if err != nil {
			return comp, err
		}
		return comp, nil
	}
	return nil, err
}

func (s service) SearchOrganizationByName(orgName string) (*models.OrgList, error) {
	osc := organization_service.GetClient()
	orgs, err := osc.SearchOrganization(orgName)
	if err != nil {
		return nil, err
	}
	result := &models.OrgList{List: make([]*models.Org, 0, len(orgs))}
	for _, org := range orgs {
		result.List = append(result.List, &models.Org{
			OrganizationID:   org.ID,
			OrganizationName: org.Name,
		})
	}
	return result, nil
}

func (s service) createOrgFromExternalID(orgID string) (*models.Company, error) {
	f := logrus.Fields{"orgID": orgID}
	osc := organization_service.GetClient()
	log.WithFields(f).Debugf("getting organization details")
	org, err := osc.GetOrganization(orgID)
	if err != nil {
		log.WithFields(f).Errorf("getting organization details failed. error = %s", err.Error())
		return nil, err
	}
	log.WithFields(f).Debugf("getting company-admin information")
	companyAdmin, err := getCompanyAdmin(orgID)
	if err != nil {
		return nil, err
	}
	f["company-admin"] = companyAdmin
	log.WithFields(f).Debugf("getting user information from cla")
	claUser, err := s.userService.GetUserByLFUserName(companyAdmin.LfUsername)
	if err != nil {
		log.WithFields(f).Errorf("getting user information from cla failed. error = %s", err.Error())
		return nil, err
	}
	if claUser == nil {
		// create cla-user
		log.WithFields(f).Debugf("cla user not found. creating cla user.")
		claUser, err = s.userService.CreateUser(companyAdmin)
		if err != nil {
			log.WithFields(f).Debugf("creating cla user failed. error = %s", err.Error())
			return nil, err
		}
	}
	newComp := &models.Company{
		CompanyACL:        []string{companyAdmin.LfUsername},
		CompanyExternalID: org.ID,
		CompanyName:       org.Name,
		CompanyManagerID:  claUser.UserID,
	}
	f["company"] = newComp
	log.WithFields(f).Debugf("creating cla company")
	// create company
	comp, err := s.repo.CreateCompany(newComp)
	if err != nil {
		log.WithFields(f).Debugf("creating cla company failed. error = %s", err.Error())
		return nil, err
	}
	return comp, nil
}

// getCompanyAdmin is helper function which queries org-service to get first company-admin
func getCompanyAdmin(companySFID string) (*models.User, error) {
	osc := organization_service.GetClient()
	result, err := osc.ListOrgUserAdminScopes(companySFID)
	if err != nil {
		log.WithField("companySFID", companySFID).Errorf("getting company-admin failed. error = %s", err.Error())
		return nil, err
	}
	for _, usc := range result.Userroles {
		for _, rs := range usc.RoleScopes {
			if rs.RoleName == "company-admin" {
				companyAdmin := &models.User{
					LfEmail:        usc.Contact.EmailAddress,
					LfUsername:     usc.Contact.Username,
					UserExternalID: usc.Contact.ID,
					Username:       usc.Contact.Name,
				}
				log.WithFields(logrus.Fields{"companySFID": companySFID, "company-admin": companyAdmin}).Debug("company-admin found")
				return companyAdmin, nil
			}
		}
	}
	log.WithField("companySFID", companySFID).Errorf("no company-admin found. error = %s", err.Error())
	return nil, errors.New("no company-admin found")
}
