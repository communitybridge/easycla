// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"

	v1ClaManager "github.com/communitybridge/easycla/cla-backend-go/cla_manager"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1User "github.com/communitybridge/easycla/cla-backend-go/user"
	easyCLAUser "github.com/communitybridge/easycla/cla-backend-go/users"
	v2AcsService "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"
	v2Company "github.com/communitybridge/easycla/cla-backend-go/v2/company"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
)

// Lead representing type of user
const Lead = "lead"

var (
	//ErrSalesForceProjectNotFound returned error if salesForce Project not found
	ErrSalesForceProjectNotFound = errors.New("salesforce Project not found")
	//ErrCLACompanyNotFound returned if EasyCLA company not found
	ErrCLACompanyNotFound = errors.New("company not found")
	//ErrGitHubRepoNotFound returned if GH Repos is not found
	ErrGitHubRepoNotFound = errors.New("gH Repo not found")
	//ErrCLAUserNotFound returned if EasyCLA User is not found
	ErrCLAUserNotFound = errors.New("cLA User not found")
	//ErrCLAManagersNotFound when cla managers arent found for given  project and company
	ErrCLAManagersNotFound = errors.New("cla Managers not found")
)

type service struct {
	companyService      company.IService
	projectService      project.Service
	repositoriesService repositories.Service
	managerService      v1ClaManager.IService
	easyCLAUserService  easyCLAUser.Service
	v2CompanyService    v2Company.Service
}

// Service interface
type Service interface {
	CreateCLAManager(claGroupID string, params cla_manager.CreateCLAManagerParams, authUsername, authEmail string) (*models.CompanyClaManager, *models.ErrorResponse)
	DeleteCLAManager(claGroupID string, params cla_manager.DeleteCLAManagerParams) *models.ErrorResponse
	InviteCompanyAdmin(contactAdmin bool, companyID string, projectID string, userEmail string, contributor *v1User.User, lFxPortalURL string) (*models.ClaManagerDesignee, *models.ErrorResponse)
	CreateCLAManagerDesignee(companyID string, projectID string, userEmail string) (*models.ClaManagerDesignee, error)
	CLAManagersRequest(companyID string, projectID string, userID string) error
}

// NewService returns instance of CLA Manager service
func NewService(compService company.IService, projService project.Service, mgrService v1ClaManager.IService, claUserService easyCLAUser.Service, repoService repositories.Service, v2CompService v2Company.Service) Service {
	return &service{
		companyService:      compService,
		projectService:      projService,
		repositoriesService: repoService,
		managerService:      mgrService,
		easyCLAUserService:  claUserService,
		v2CompanyService:    v2CompService,
	}
}

// CreateCLAManager creates Cla Manager
func (s *service) CreateCLAManager(claGroupID string, params cla_manager.CreateCLAManagerParams, authUsername, authEmail string) (*models.CompanyClaManager, *models.ErrorResponse) {
	if *params.Body.FirstName == "" || *params.Body.LastName == "" || *params.Body.UserEmail == "" {
		msg := "firstName, lastName and UserEmail cannot be empty"
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
		msg := buildErrorMessage("company lookup error", claGroupID, params, companyErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	claGroup, err := getCLAGroup(claGroupID, params.ProjectSFID, s.projectService)
	if err != nil || claGroup == nil {
		msg := buildErrorMessage("project cla lookup failure or project doesnt have CCLA", claGroupID, params, err)
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
		designeeEmail := *params.Body.UserEmail
		lfxUser, lfxUserErr := userServiceClient.SearchUserByEmail(designeeEmail)
		if lfxUserErr != nil || lfxUser == nil {
			msg := fmt.Sprintf("LFID not existing for %s", designeeName)
			log.Warn(msg)
			sendEmailToUserWithNoLFID(claGroup, authUsername, authEmail, designeeName, designeeEmail)
			return nil, &models.ErrorResponse{
				Message: msg,
				Code:    "400",
			}
		}

		msg := fmt.Sprintf("Failed search for User with firstname : %s, lastname: %s , email: %s , error: %v ",
			*params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, userErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	// Check if user exists in easyCLA DB, if not add User
	log.Debugf("Checking user: %s in easyCLA records", user.Username)
	claUser, claUserErr := s.easyCLAUserService.GetUserByLFUserName(user.Username)
	if claUserErr != nil {
		msg := fmt.Sprintf("Problem getting claUser by :%s, error: %+v ", user.Username, claUserErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}

	if claUser == nil {
		msg := fmt.Sprintf("User not found when searching by LFID: %s and shall be created", user.Username)
		log.Debug(msg)
		userName := fmt.Sprintf("%s %s", *params.Body.FirstName, *params.Body.LastName)
		_, currentTimeString := utils.CurrentTime()
		claUserModel := &v1Models.User{
			UserExternalID: params.CompanySFID,
			LfEmail:        *user.Emails[0].EmailAddress,
			Admin:          true,
			LfUsername:     user.Username,
			DateCreated:    currentTimeString,
			DateModified:   currentTimeString,
			Username:       userName,
			Version:        "v1",
		}
		newUserModel, userModelErr := s.easyCLAUserService.CreateUser(claUserModel)
		if userModelErr != nil {
			msg := fmt.Sprintf("Failed to create user : %+v", claUserModel)
			log.Warn(msg)
			return nil, &models.ErrorResponse{
				Message: msg,
				Code:    "400",
			}
		}
		log.Debugf("Created easyCLAUser %+v ", newUserModel)
	}

	// GetSFProject
	ps := v2ProjectService.GetClient()
	projectSF, projectErr := ps.GetProject(params.ProjectSFID)
	if projectErr != nil {
		msg := buildErrorMessage("project service lookup error", claGroupID, params, projectErr)
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
	hasScope, err := orgClient.IsUserHaveRoleScope("cla-manager", user.ID, params.CompanySFID, params.ProjectSFID)
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
			Code:    "409",
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
		sigMsg := fmt.Sprintf("Signature not found for project: %s and company: %s ", claGroupID, companyModel.CompanyID)
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

func (s *service) DeleteCLAManager(claGroupID string, params cla_manager.DeleteCLAManagerParams) *models.ErrorResponse {
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

	signature, deleteErr := s.managerService.RemoveClaManager(companyModel.CompanyID, claGroupID, params.UserLFID)

	if deleteErr != nil {
		msg := buildErrorMessageDelete(params, deleteErr)
		log.Warn(msg)
		return &models.ErrorResponse{
			Message: msg,
			Code:    "400",
		}
	}
	if signature == nil {
		msg := fmt.Sprintf("Not found signature for project: %s and company: %s ", claGroupID, companyModel.CompanyID)
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
		msg := "Problem getting role ID for cla-manager-designee"
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

func (s *service) InviteCompanyAdmin(contactAdmin bool, companyID string, projectID string, userEmail string, contributor *v1User.User, LfxPortalURL string) (*models.ClaManagerDesignee, *models.ErrorResponse) {
	orgService := v2OrgService.GetClient()
	projectService := v2ProjectService.GetClient()
	userService := v2UserService.GetClient()

	// Get repo instance (assist in getting salesforce project)
	log.Debugf("Get salesforce project by claGroupID: %s ", projectID)
	ghRepoModel, ghRepoErr := s.repositoriesService.GetGithubRepositoryByCLAGroup(projectID)
	if ghRepoErr != nil || ghRepoModel.RepositorySfdcID == "" {
		msg := fmt.Sprintf("Problem getting salesforce project by claGroupID : %s ", projectID)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Code:    "404",
			Message: msg,
		}
	}

	// Get company
	log.Debugf("Get company for companyID: %s ", companyID)
	companyModel, companyErr := s.companyService.GetCompany(companyID)
	if companyErr != nil || companyModel.CompanyExternalID == "" {
		msg := fmt.Sprintf("Problem getting company for companyID: %s ", companyID)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Code:    "404",
			Message: msg,
		}
	}

	project, projectErr := projectService.GetProject(ghRepoModel.RepositorySfdcID)
	if projectErr != nil {
		msg := fmt.Sprintf("Problem getting project by ID: %s ", projectID)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Code:    "400",
			Message: msg,
		}
	}

	organization, orgErr := orgService.GetOrganization(companyModel.CompanyExternalID)
	if orgErr != nil {
		msg := fmt.Sprintf("Problem getting company by ID: %s ", companyID)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Code:    "400",
			Message: msg,
		}
	}

	// Get suggested CLA Manager user details
	user, userErr := userService.SearchUserByEmail(userEmail)
	if userErr != nil {
		msg := fmt.Sprintf("Problem getting user for userEmail: %s , error: %+v", userEmail, userErr)
		log.Warn(msg)
		return nil, &models.ErrorResponse{
			Code:    "400",
			Message: msg,
		}
	}

	// Check if sending cla manager request to company admin
	if contactAdmin {
		log.Debugf("Sending email to company Admin")
		scopes, listScopeErr := orgService.ListOrgUserAdminScopes(companyID)
		if listScopeErr != nil {
			msg := fmt.Sprintf("Admin lookup error for organisation SFID: %s ", companyID)
			return nil, &models.ErrorResponse{
				Code:    "400",
				Message: msg,
			}
		}
		for _, admin := range scopes.Userroles {
			// Check if is Gerrit User or GH User
			if contributor.LFUsername != "" && contributor.LFEmail != "" {
				sendEmailToOrgAdmin(admin.Contact.EmailAddress, admin.Contact.Name, organization.Name, project.Name, contributor.LFEmail, contributor.LFUsername, LfxPortalURL)
			} else {
				sendEmailToOrgAdmin(admin.Contact.EmailAddress, admin.Contact.Name, organization.Name, project.Name, contributor.UserGithubID, contributor.UserGithubUsername, LfxPortalURL)
			}

		}
		return nil, nil
	}

	claManagerDesignee, err := s.CreateCLAManagerDesignee(organization.ID, project.ID, userEmail)

	if err != nil {
		msg := fmt.Sprintf("Problem creating cla Manager Designee for user :%s, error: %+v ", userEmail, err)
		return nil, &models.ErrorResponse{
			Code:    "400",
			Message: msg,
		}
	}

	log.Debugf("Sending Email to CLA Manager Designee email: %s ", userEmail)

	if contributor.LFUsername != "" && contributor.LFEmail != "" {
		sendEmailToCLAManagerDesignee(LfxPortalURL, organization.Name, project.Name, userEmail, user.Name, contributor.LFEmail, contributor.LFUsername)
	} else {
		sendEmailToCLAManagerDesignee(LfxPortalURL, organization.Name, project.Name, userEmail, user.Name, contributor.UserGithubID, contributor.UserGithubUsername)
	}

	log.Debugf("CLA Manager designee created : %+v", claManagerDesignee)
	return claManagerDesignee, nil

}

func (s *service) CLAManagersRequest(companyID string, projectID string, userID string) error {
	projectService := v2ProjectService.GetClient()
	//Get EasyCLA Company
	company, companyErr := s.companyService.GetCompany(companyID)
	if companyErr != nil {
		msg := fmt.Sprintf("EasyCLA Bad Request- Problem getting company :%s ", companyID)
		log.Warn(msg)
		return ErrCLACompanyNotFound
	}

	ghRepoModel, ghRepoErr := s.repositoriesService.GetGithubRepositoryByCLAGroup(projectID)
	if ghRepoErr != nil || ghRepoModel.RepositorySfdcID == "" {
		msg := fmt.Sprintf("Problem getting salesforce project by claGroupID : %s ", projectID)
		log.Warn(msg)
		return ErrGitHubRepoNotFound
	}

	project, projectErr := projectService.GetProject(ghRepoModel.RepositorySfdcID)
	if projectErr != nil {
		msg := fmt.Sprintf("Problem getting salesforce project for given repository association :%+v ", projectErr)
		log.Warn(msg)
		return ErrSalesForceProjectNotFound
	}

	// Search for Easy CLA User
	log.Debugf("Getting user by ID: %s", userID)
	userModel, userErr := s.easyCLAUserService.GetUser(userID)
	if userErr != nil {
		msg := fmt.Sprintf("Problem getting user by ID: %s ", userID)
		log.Warn(msg)
		return ErrCLAUserNotFound
	}

	// Get CLA Managers
	claManagers, managersErr := s.v2CompanyService.GetCompanyProjectCLAManagers(company.CompanyExternalID, ghRepoModel.RepositorySfdcID)
	if managersErr != nil || len(claManagers.List) == 0 {
		msg := fmt.Sprintf("Problem getting salesforce Company for for easyCLA company: %s ", company.CompanyName)
		log.Warn(msg)
		return ErrCLAManagersNotFound
	}

	// Send Emails to respective CLA Managers
	for _, claManager := range claManagers.List {
		sendEmailToCLAManager(claManager.Name, claManager.Email, userModel.GithubUsername, company.CompanyName, project.Name)
	}

	return nil
}

func sendEmailToCLAManager(manager string, managerEmail string, contributorName string, company string, project string) {
	subject := fmt.Sprintf("EasyCLA: Approval Request for contributor: %s  ", contributorName)
	recipients := []string{managerEmail}
	body := fmt.Sprintf(`
	<p>Hello %s,</p>
	<p>This is a notification email from EasyCLA regarding the organization %s.</p>
	<p>The following contributor would like to submit a contribution to %s 
	   and is requesting to be approved as a contributor for your organization: </p>
	<p> %s </p>
	<p>Please notify the contributor once they are added so that they may complete the contribution process.</p>
	<p>If you need help or have questions about EasyCLA, you can
	<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
	<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
	support</a>.</p>
	<p>Thanks,</p>
	<p>EasyCLA support team </p>`,
		manager, company, project, company)
	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendEmailToOrgAdmin(adminEmail string, admin string, company string, project string, contributorID string, contributorName string, corporateConsole string) {
	subject := fmt.Sprintf("EasyCLA:  Invitation to Sign the %s Corporate CLA and add to approved list %s ", company, contributorID)
	recipients := []string{adminEmail}
	body := fmt.Sprintf(`
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
<p>Thanks,</p>
<p>EasyCLA support team </p>`,
		admin, project, project, contributorName, contributorID, corporateConsole, project)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

func sendEmailToCLAManagerDesignee(corporateConsole string, company string, project string, designeeEmail string, designeeName string, contributorID string, contributorName string) {
	subject := fmt.Sprintf("EasyCLA:  Invitation to Sign the %s Corporate CLA and add to approved list %s ", company, contributorID)
	recipients := []string{designeeEmail}
	body := fmt.Sprintf(`
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
<p>Thanks,</p>
<p>EasyCLA support team </p>`,
		designeeName, project, project, contributorName, contributorID, corporateConsole, project)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(errPrefix string, claGroupID string, params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("%s - problem creating new CLA Manager Request using company SFID: %s, project ID: %s, first name: %s, last name: %s, user email: %s, error: %+v",
		errPrefix, params.CompanySFID, claGroupID, *params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, err)
}

func getCLAGroup(projectID string, projectSFID string, projectService project.Service) (*v1Models.Project, error) {
	var claGroup v1Models.Project
	// Search for projects by ProjectSFID
	projects, projectErr := projectService.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
		ProjectSFID: projectSFID,
	})
	// Get unique project by passed CLAGroup ID parameter
	for _, proj := range projects.Projects {
		if proj.ProjectID == projectID && proj.ProjectCCLAEnabled {
			claGroup = proj
			break
		}
	}

	if projectErr != nil {
		return nil, projectErr
	}

	return &claGroup, nil
}
