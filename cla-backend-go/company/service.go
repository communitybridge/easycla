// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type service struct {
	repo                RepositoryService
	userDynamoRepo      user.RepositoryService
	sesClient           *ses.SES
	senderEmailAddress  string
	corporateConsoleURL string
}

// Service interface defining the public functions
type Service interface { // nolint
	GetCompany(companyID string) (*models.Company, error)
	SearchCompanyByName(companyName string, nextKey string) (*models.Companies, error)
	AddUserToCompanyAccessList(companyID string, inviteID string, lfid string) error
	SendApprovalEmail(companyName, recipientAddress, senderAddress string, user *user.CLAUser) error
	SendRequestAccessEmail(companyID string, user *user.CLAUser) error
	GetPendingCompanyInviteRequests(companyID string) ([]models.CompanyInviteUser, error)
	DeletePendingCompanyInviteRequest(inviteID string) error
}

// NewService creates a new company service object
func NewService(repo RepositoryService, awsSession *session.Session, senderEmailAddress, corporateConsoleURL string, userDynamoRepo user.RepositoryService) Service {
	return service{
		repo:                repo,
		userDynamoRepo:      userDynamoRepo,
		sesClient:           ses.New(awsSession),
		senderEmailAddress:  senderEmailAddress,
		corporateConsoleURL: corporateConsoleURL,
	}
}

// GetCompany returns the company associated with the company ID
func (s service) GetCompany(companyID string) (*models.Company, error) {
	dbCompanyModel, err := s.repo.GetCompany(companyID)
	if err != nil {
		log.Warnf("Error retrieving company by company ID: %s, error: %v", companyID, err)
		return nil, err
	}

	const timeFormat = "2006-01-02T15:04:05.999999+0000"
	// Convert the "string" date time
	createdDateTime, err := time.Parse(timeFormat, dbCompanyModel.Created)
	if err != nil {
		log.Warnf("Error converting created date time for company: %s, error: %v", companyID, err)
		return nil, err
	}
	updateDateTime, err := time.Parse(timeFormat, dbCompanyModel.Updated)
	if err != nil {
		log.Warnf("Error converting updated date time for company: %s, error: %v", companyID, err)
		return nil, err
	}

	// Convert the local DB model to a public swagger model
	return &models.Company{
		CompanyACL:  dbCompanyModel.CompanyACL,
		CompanyID:   dbCompanyModel.CompanyID,
		CompanyName: dbCompanyModel.CompanyName,
		Created:     strfmt.DateTime(createdDateTime),
		Updated:     strfmt.DateTime(updateDateTime),
	}, nil
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

// AddUserToCompanyAccessList adds a user to the specified company
func (s service) AddUserToCompanyAccessList(companyID string, inviteID string, lfid string) error {
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
			err = s.repo.DeletePendingCompanyInviteRequest(inviteID)
			if err != nil {
				log.Warnf("Error deleting pending company invite request with inviteID: %s, error: %v", inviteID, err)
				return fmt.Errorf("failed to delete pending invite")
			}
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

	userProfile, err := s.userDynamoRepo.GetUserAndProfilesByLFID(lfid)
	if err != nil {
		log.Warnf("Error getting user profile by LFID: %s, error: %v", lfid, err)
		return nil
	}

	recipientEmailAddress := userProfile.LFEmail

	err = s.SendApprovalEmail(company.CompanyName, recipientEmailAddress, s.senderEmailAddress, &userProfile)
	if err != nil {
		return errors.New("failed to send notification email")
	}

	// Remove pending invite ID once approval emails are sent
	err = s.repo.DeletePendingCompanyInviteRequest(inviteID)
	if err != nil {
		return fmt.Errorf("failed to delete pending invite")
	}

	return nil
}

// SendApprovalEmail sends the approval email when provided the company name, address and user object
func (s service) SendApprovalEmail(companyName, recipientAddress, senderAddress string, user *user.CLAUser) error {
	var (
		Sender    = senderAddress
		Recipient = recipientAddress
		Subject   = "CLA: Approval of Access for Corporate CLA"

		//The email body for recipients with non-HTML email clients.
		TextBody = fmt.Sprintf(`Hello %s,

You have now been granted access to the organization: %s

	%s <%s>

- Linux Foundation CLA System`, user.Name, companyName, user.LFUsername, user.LFEmail)
		// The character encoding for the email.
		CharSet = "UTF-8"
	)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	_, err := s.sesClient.SendEmail(input)
	if err != nil {
		log.Warnf("Error sending mail, error: %v", err)
		return err
	}

	return nil
}

// SendRequestAccessEmail sends the request access e-mail when provided the company ID and user object
func (s service) SendRequestAccessEmail(companyID string, user *user.CLAUser) error {

	// Get Company
	company, err := s.repo.GetCompany(companyID)
	if err != nil {
		log.Warnf("Error fetching company by company ID: %s, error: %v", companyID, err)
		return err
	}

	// Add a pending request to the company-invites table
	err = s.repo.AddPendingCompanyInviteRequest(companyID, user.UserID)
	if err != nil {
		log.Warnf("Error adding pending company invite request using company ID: %s, user ID: %s, error: %v", companyID, user.UserID, err)
		return err
	}

	// Send Email to every CLA Manager in the Company ACL
	Subject := "CLA: Request of Access for Corporate CLA Manager"

	for _, admin := range company.CompanyACL {
		// Retrieve admin's user profile for email and name
		adminUser, err := s.userDynamoRepo.GetUserAndProfilesByLFID(admin)
		if err != nil {
			log.Warnf("Error fetching user profile using admin: %s, error: %v", admin, err)
			return err
		}

		TextBody := fmt.Sprintf(`Hello %s, 

The following user is requesting access to your organization: %s

	%s <%s>

Please navigate to the Corporate Console using the link below, where you can approve this user's request.

%s

- Linux Foundation CLA System`, adminUser.Name, company.CompanyName, user.LFUsername, user.LFEmail, s.corporateConsoleURL)

		CharSet := "UTF-8"

		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(adminUser.LFEmail),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Text: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(TextBody),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
			},
			Source: aws.String(s.senderEmailAddress),
		}

		_, err = s.sesClient.SendEmail(input)
		if err != nil {
			log.Warnf("Error sending mail, error: %v", err)
			return err
		}
	}

	return nil
}

// GetPendingCompanyInviteRequests returns a list of company invites when provided the company ID
func (s service) GetPendingCompanyInviteRequests(companyID string) ([]models.CompanyInviteUser, error) {
	companyInvites, err := s.repo.GetPendingCompanyInviteRequests(companyID)
	if err != nil {
		return nil, err
	}

	var users []models.CompanyInviteUser
	for _, invite := range companyInvites {
		inviteID := invite.CompanyInviteID
		userID := invite.UserID
		dbUserModel, err := s.userDynamoRepo.GetUser(userID)
		if err != nil {
			log.Warnf("Error fetching user with userID: %s, error: %v", userID, err)
			continue
		}

		users = append(users, models.CompanyInviteUser{
			InviteID:  inviteID,
			UserName:  dbUserModel.UserName,
			UserEmail: dbUserModel.LFEmail,
			UserLFID:  dbUserModel.LFUsername,
		})
	}

	return users, nil

}

// DeletePendingCompanyInviteRequest deletes the pending company invite request when provided the invite ID
func (s service) DeletePendingCompanyInviteRequest(inviteID string) error {
	// When a CLA Manager Declines a pending invite, remove the invite from the table
	err := s.repo.DeletePendingCompanyInviteRequest(inviteID)
	if err != nil {
		log.Warnf("Error deleting the pending company invite with invite ID: %s, error: %v", inviteID, err)
		return err
	}

	return nil
}
