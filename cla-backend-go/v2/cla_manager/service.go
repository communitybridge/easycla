// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	"github.com/communitybridge/easycla/cla-backend-go/project"

	v1ClaManager "github.com/communitybridge/easycla/cla-backend-go/cla_manager"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2AcsService "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
)

// Lead representing type of user
const Lead = "lead"

type service struct {
	companyService company.IService
	projectService project.Service
	managerService v1ClaManager.IService
}

// Service interface
type Service interface {
	CreateCLAManager(params cla_manager.CreateCLAManagerParams, authEmail string) (*models.CompanyClaManager, *models.ErrorResponse)
	DeleteCLAManager(params cla_manager.DeleteCLAManagerParams) *models.ErrorResponse
	CreateCLAManagerDesignee(companyID string, projectID string, userEmail string) (*models.ClaManagerDesignee, error)
}

// NewService returns instance of CLA Manager service
func NewService(compService company.IService, projService project.Service, mgrService v1ClaManager.IService) Service {
	return &service{
		companyService: compService,
		projectService: projService,
		managerService: mgrService,
	}
}

// CreateCLAManager creates Cla Manager
func (s *service) CreateCLAManager(params cla_manager.CreateCLAManagerParams, authEmail string) (*models.CompanyClaManager, *models.ErrorResponse) {
	if *params.Body.FirstName == "" || *params.Body.LastName == "" || *params.Body.UserEmail == "" {
		msg := fmt.Sprintf("firstName, lastName and UserEmail cannot be empty")
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	// Search for salesForce Company aka external Company
	log.Debugf("Getting company by external ID : %s", params.CompanySFID)
	companyModel, companyErr := s.companyService.GetCompanyByExternalID(params.CompanySFID)
	if companyErr != nil || companyModel == nil {
		msg := buildErrorMessage("company lookup error", params, companyErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	claGroup, err := getCLAGroup(params.ProjectID, params.ProjectSFID, s.projectService)
	if err != nil || claGroup == nil {
		msg := buildErrorMessage("project cla lookup failure or project doesnt have CCLA", params, err)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	// Get user by firstname,lastname and email parameters
	userServiceClient := v2UserService.GetClient()
	user, userErr := userServiceClient.SearchUsers(*params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail)

	if userErr != nil {
		designeeName := fmt.Sprintf("%s %s", *params.Body.FirstName, *params.Body.LastName)
		log.Debugf("LFID not existing for %s", designeeName)
		sendEmailToUserWithNoLFID(claGroup, authEmail, designeeName, *params.Body.UserEmail)
		msg := fmt.Sprintf("Failed search for User with firstname : %s, lastname: %s , email: %s , error: %v ",
			*params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, userErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	// GetSFProject
	ps := v2ProjectService.GetClient()
	projectSF, projectErr := ps.GetProject(params.ProjectSFID)
	if projectErr != nil {
		msg := buildErrorMessage("project service lookup error", params, projectErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	log.Warn("Getting role")
	// Get RoleID for cla-manager
	acsClient := v2AcsService.GetClient()

	roleID, roleErr := acsClient.GetRoleID("cla-manager")
	if roleErr != nil {
		msg := buildErrorMessageCreate(params, roleErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	log.Debugf("Role ID for cla-manager-role : %s", roleID)
	log.Debugf("Creating user role Scope for user : %s ", *params.Body.UserEmail)

	orgClient := v2OrgService.GetClient()
	hasScope, err := orgClient.IsUserHaveRoleScope(roleID, user.ID, params.CompanySFID, params.ProjectSFID)
	if err != nil {
		msg := buildErrorMessageCreate(params, err)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	if hasScope {
		msg := fmt.Sprintf("User %s is already cla-manager for Company: %s and Project: %s", user.Username, params.CompanySFID, params.ProjectSFID)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	scopeErr := orgClient.CreateOrgUserRoleOrgScopeProjectOrg(*params.Body.UserEmail, params.ProjectSFID, params.CompanySFID, roleID)
	if scopeErr != nil {
		msg := buildErrorMessageCreate(params, scopeErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	if user.Type == Lead {
		// convert user to contact
		log.Debug("converting lead to contact")
		err := userServiceClient.ConvertToContact(user.ID)
		if err != nil {
			msg := fmt.Sprintf("converting lead to contact failed: %v", err)
			log.Warn(msg)
			return nil, &models.ErrorResponse{
				Message: msg,
				Code:    "400",
			}
		}
	}

	// Add CLA Manager to Database
	signature, addErr := s.managerService.AddClaManager(companyModel.CompanyID, claGroup.ProjectID, user.Username)
	if addErr != nil {
		msg := buildErrorMessageCreate(params, addErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	if signature == nil {
		sigMsg := fmt.Sprintf("Signature not found for project: %s and company: %s ", params.ProjectID, companyModel.CompanyID)
		log.Warn(sigMsg)
		return nil, &models.ErrorResponse{
			Message: sigMsg,
			Code:    "400",
		}
	}

	claCompanyManager := &models.CompanyClaManager{
		LfUsername:       user.Username,
		Email:            *params.Body.UserEmail,
		UserSfid:         user.ID,
		ApprovedOn:       time.Now().String(),
		ProjectSfid:      params.ProjectSFID,
		ClaGroupName:     claGroup.ProjectName,
		ProjectID:        claGroup.ProjectID,
		ProjectName:      projectSF.Name,
		OrganizationName: companyModel.CompanyName,
		OrganizationSfid: params.CompanySFID,
		Name:             fmt.Sprintf("%s %s", user.FirstName, user.LastName),
	}
	return claCompanyManager, nil
}

func (s *service) DeleteCLAManager(params cla_manager.DeleteCLAManagerParams) *models.ErrorResponse {
	// Get user by firstname,lastname and email parameters
	userServiceClient := v2UserService.GetClient()
	user, userErr := userServiceClient.GetUserByUsername(params.UserLFID)

	if userErr != nil {
		msg := fmt.Sprintf("Failed to get user when searching by username: %s , error: %v ", params.UserLFID, userErr)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	// Search for salesForce Company aka external Company
	companyModel, companyErr := s.companyService.GetCompanyByExternalID(params.CompanySFID)
	if companyErr != nil || companyModel == nil {
		msg := buildErrorMessageDelete(params, companyErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	acsClient := v2AcsService.GetClient()

	roleID, roleErr := acsClient.GetRoleID("cla-manager")
	if roleErr != nil {
		msg := buildErrorMessageDelete(params, roleErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	log.Debugf("Role ID for cla-manager-role : %s", roleID)

	// Get Scope ID
	orgClient := v2OrgService.GetClient()
	scopeID, scopeErr := orgClient.GetScopeID(params.CompanySFID, "cla-manager", "project|organization", params.UserLFID)
	if scopeErr != nil {
		msg := buildErrorMessageDelete(params, scopeErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	email := *user.Emails[0].EmailAddress
	deleteErr := orgClient.DeleteOrgUserRoleOrgScopeProjectOrg(params.CompanySFID, roleID, scopeID, user.Username, email)
	if deleteErr != nil {
		msg := buildErrorMessageDelete(params, deleteErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	signature, deleteErr := s.managerService.RemoveClaManager(companyModel.CompanyID, params.ProjectID, params.UserLFID)

	if deleteErr != nil {
		msg := buildErrorMessageDelete(params, deleteErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	if signature == nil {
		msg := fmt.Sprintf("Not found signature for project: %s and company: %s ", params.ProjectID, companyModel.CompanyID)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	return nil
}

//CreateCLAManagerDesignee creates designee for cla manager prospect
func (s *service) CreateCLAManagerDesignee(companyID string, projectID string, userEmail string) (*models.ClaManagerDesignee, error) {
	// integrate user,acs,org and project services
	userClient := v2UserService.GetClient()
	acServiceClient := v2AcsService.GetClient()
	orgClient := v2OrgService.GetClient()
	projectClient := v2ProjectService.GetClient()

	user, userErr := userClient.SearchUserByEmail(userEmail)
	if userErr != nil {
		log.Debugf("Failed to get user by email: %s , error: %+v", userEmail, userErr)
		return nil, userErr
	}

	if user.Type == Lead {
		log.Debugf("Converting user: %s from lead to contact ", userEmail)
		contactErr := userClient.ConvertToContact(user.ID)
		if contactErr != nil {
			log.Debugf("failed to convert user: %s to contact ", userEmail)
			return nil, contactErr
		}
	}

	roleID, designeeErr := acServiceClient.GetRoleID("cla-manager-designee")
	if designeeErr != nil {
		msg := fmt.Sprintf("Problem getting role ID for cla-manager-designee")
		log.Warn(msg)
		return nil, designeeErr
	}

	scopeErr := orgClient.CreateOrgUserRoleOrgScopeProjectOrg(userEmail, projectID, companyID, roleID)
	if scopeErr != nil {
		msg := fmt.Sprintf("Problem creating projectOrg scope for email: %s , projectID: %s, companyID: %s", userEmail, projectID, companyID)
		log.Warn(msg)
		return nil, scopeErr
	}

	projectSF, projectErr := projectClient.GetProject(projectID)
	if projectErr != nil {
		msg := fmt.Sprintf("Problem getting project :%s ", projectID)
		log.Debug(msg)
		return nil, projectErr
	}

	claManagerDesignee := &models.ClaManagerDesignee{
		LfUsername:  user.Username,
		UserSfid:    user.ID,
		Type:        user.Type,
		AssignedOn:  time.Now().String(),
		Email:       userEmail,
		ProjectSfid: projectID,
		CompanySfid: companyID,
		ProjectName: projectSF.Name,
	}
	return claManagerDesignee, nil
}

func sendEmailToOrgAdmin(adminEmail string, admin string, company string, project string, contributorEmail string, contributorName string, corporateConsole string) {
	subject := fmt.Sprintf("EasyCLA:  Invitation to Sign the %s Corporate CLA and add to approved list %s ", company, contributorEmail)
	recipients := []string{adminEmail}
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
<p>The following contributor would like to submit a contribution to %s 
   and is requesting to be whitelisted as a contributor for your organization: </p>
<p> %s %s </p>
<p>Before the contribution can be accepted, your organization must sign a CLA.
<p>Kindly login to this portal %s and sign the CLA for this project %s. </p>
<p>Please notify the contributor once they are added so that they may complete the contribution process.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>

<p>Thanks,

<p>EasyCLA support team </p>
</body>
</html>
	
	`, admin, project, project, contributorName, contributorEmail, corporateConsole, project)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendEmailToCLAManagerDesignee(corporateConsole string, company string, project string, designeeEmail string, designeeName string, contributorEmail string, contributorName string) {
	subject := fmt.Sprintf("EasyCLA:  Invitation to Sign the %s Corporate CLA and add to approved list %s ", company, contributorEmail)
	recipients := []string{designeeEmail}
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
<p>The following contributor would like to submit a contribution to %s 
   and is requesting to be whitelisted as a contributor for your organization: </p>
<p> %s %s </p>
<p>Before the contribution can be accepted, your organization must sign a CLA.
<p>Kindly login to this portal %s and sign the CLA for this project %s. </p>
<p>Please notify the contributor once they are added so that they may complete the contribution process.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>

<p>Thanks,

<p>EasyCLA support team </p>
</body>
</html>
	
	`, designeeName, project, project, contributorName, contributorEmail, corporateConsole, project)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(errPrefix string, params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("%s - problem creating new CLA Manager Request using company SFID: %s, project ID: %s, first name: %s, last name: %s, user email: %s, error: %+v",
		errPrefix, params.CompanySFID, params.ProjectID, *params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, err)
}

func getCLAGroup(projectID string, projectSFID string, projectService project.Service) (*v1Models.Project, error) {
	var claGroup *v1Models.Project
	// Search for projects by ProjectSFID
	projects, projectErr := projectService.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
		ProjectSFID: projectSFID,
	})
	// Get unique project by passed CLAGroup ID parameter
	for _, proj := range projects.Projects {
		if proj.ProjectID == projectID && proj.ProjectCCLAEnabled {
			claGroup = &proj
			break
		}
	}

	if projectErr != nil {
		return nil, projectErr
	}

	return claGroup, nil
}
