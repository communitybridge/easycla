package company

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type service struct {
	repo                Repository
	userDynamoRepo      user.RepositoryDynamo
	sesClient           *ses.SES
	senderEmailAddress  string
	corporateConsoleURL string
}

func NewService(repo Repository, awsSession *session.Session, senderEmailAddress, corporateConsoleURL string, userDynamoRepo user.RepositoryDynamo) service {
	return service{
		repo:                repo,
		userDynamoRepo:      userDynamoRepo,
		sesClient:           ses.New(awsSession),
		senderEmailAddress:  senderEmailAddress,
		corporateConsoleURL: corporateConsoleURL,
	}
}

//
func (s service) AddUserToCompanyAccessList(companyID string, inviteID string, lfid string) error {
	// call getcompany function
	company, err := s.repo.GetCompany(companyID)
	if err != nil {
		return err
	}

	// perform ACL check
	// check if user already exists in the company acl
	for _, acl := range company.CompanyACL {
		if acl == lfid {
			fmt.Println(fmt.Sprintf("User %s has already been added to the company acl", lfid))
			err = s.repo.DeletePendingCompanyInviteRequest(inviteID)
			if err != nil {
				return fmt.Errorf("Failed to delete pending invite")
			}
			return nil
		}
	}
	// add user to string set
	company.CompanyACL = append(company.CompanyACL, lfid)

	err = s.repo.UpdateCompanyAccessList(companyID, company.CompanyACL)
	if err != nil {
		return err
	}

	user, err := s.userDynamoRepo.GetUserAndProfilesByLFID(lfid)
	if err != nil {
		return nil
	}

	recipientEmailAddress := user.LFEmail

	err = s.SendApprovalEmail(company.CompanyName, recipientEmailAddress, s.senderEmailAddress, &user)
	if err != nil {
		return errors.New("Failed to send notification email")
	}

	// Remove pending invite ID once approval emails are sent
	err = s.repo.DeletePendingCompanyInviteRequest(inviteID)
	if err != nil {
		return fmt.Errorf("Failed to delete pending invite")
	}

	return nil
}

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
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (s service) SendRequestAccessEmail(companyID string, user *user.CLAUser) error {

	// Get Company
	company, err := s.repo.GetCompany(companyID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Add a pending request to the company-invites table
	err = s.repo.AddPendingCompanyInviteRequest(companyID, user.UserID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send Email to every CLA Manager in the Company ACL
	Subject := "CLA: Request of Access for Corporate CLA Manager"

	for _, admin := range company.CompanyACL {
		// Retrieve admin's user profile for email and name
		adminUser, err := s.userDynamoRepo.GetUserAndProfilesByLFID(admin)
		if err != nil {
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
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (s service) GetPendingCompanyInviteRequests(companyID string) ([]models.CompanyInviteUser, error) {
	companyInvites, err := s.repo.GetPendingCompanyInviteRequests(companyID)
	if err != nil {
		return nil, err
	}

	users := []models.CompanyInviteUser{}
	for _, invite := range companyInvites {
		inviteID := invite.CompanyInviteID
		userID := invite.UserID
		user, err := s.userDynamoRepo.GetUser(userID)
		if err != nil {
			continue
		}

		users = append(users, models.CompanyInviteUser{
			InviteID:  inviteID,
			UserName:  user.UserName,
			UserEmail: user.LFEmail,
			UserLFID:  user.LFUsername,
		})
	}

	return users, nil

}
func (s service) DeletePendingCompanyInviteRequest(inviteID string) error {
	// When a CLA Manager Declines a pending invite, remove the invite from the table
	err := s.repo.DeletePendingCompanyInviteRequest(inviteID)
	if err != nil {
		return err
	}

	return nil

}
