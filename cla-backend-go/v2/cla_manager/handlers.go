// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	v1ClaManager "github.com/communitybridge/easycla/cla-backend-go/cla_manager"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	"github.com/communitybridge/easycla/cla-backend-go/project"

	v1ProjectParams "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2AcsService "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/go-openapi/runtime/middleware"
)

// Configure is the API handler routine for CLA Manager routes
func Configure(api *operations.EasyclaAPI, managerService v1ClaManager.IService, companyService company.IService, projectService project.Service) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerUnauthorized()
		}

		// Get user by firstname,lastname and email parameters
		userServiceClient := v2UserService.GetClient()
		user, userErr := userServiceClient.SearchUsers(params.Body.FirstName, params.Body.LastName, params.Body.UserEmail)

		if userErr != nil {
			msg := fmt.Sprintf("Failed to get user when searching by firstname : %s, lastname: %s , email: %s , error: %v ",
				params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, userErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Get SFProjectModel
		projectServiceClient := v2ProjectService.GetClient()
		projectSF, projectSFErr := projectServiceClient.GetProject(params.ProjectSFID)
		if projectSFErr != nil {
			msg := buildErrorMessage("projectSF lookup error", params, projectSFErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Search for salesForce Company aka external Company
		companyModel, companyErr := companyService.GetCompanyByExternalID(params.CompanySFID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessage("company lookup error", params, companyErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Search for projects by ProjectSFID
		projects, projectErr := projectService.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
			ProjectSFID: params.ProjectSFID,
		})

		var claGroup *v1Models.Project
		// Get unique project by passed CLAGroup ID parameter
		for _, proj := range projects.Projects {
			if proj.ProjectID == params.ProjectID {
				claGroup = &proj
				break
			}
		}

		if projectErr != nil || claGroup == nil {
			msg := buildErrorMessage("project lookup error", params, projectErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Send Email if LFID doesnt exist
		if authUser.UserName == "" {
			designeeName := fmt.Sprintf("%s %s", params.Body.FirstName, params.Body.LastName)
			log.Debugf("LFID not existing for %s", designeeName)
			sendEmailToUserWithNoLFID(claGroup, authUser.Email, designeeName, params.Body.UserEmail)
			msg := fmt.Sprintf("%s does not have LFID ", designeeName)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Get RoleID for cla-manager
		acsClient := v2AcsService.GetClient()

		roleID, roleErr := acsClient.GetRoleID("cla-manager")
		if roleErr != nil {
			msg := buildErrorMessageCreate(params, roleErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}
		log.Debugf("Role ID for cla-manager-role : %s", roleID)
		log.Debugf("Creating user role Scope for user : %v ", *user.Email)

		orgClient := v2OrgService.GetClient()

		scopeErr := orgClient.CreateOrgUserRoleOrgScopeProjectOrg(*user.Email, params.ProjectSFID, params.CompanySFID, roleID)
		if scopeErr != nil {
			msg := buildErrorMessageCreate(params, scopeErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Add CLA Manager to Database
		signature, addErr := managerService.AddClaManager(companyModel.CompanyID, claGroup.ProjectID, user.Username)
		if addErr != nil {
			msg := buildErrorMessageCreate(params, addErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}
		if signature == nil {
			sigMsg := fmt.Sprintf("Signature not found for project: %s and company: %s ", params.ProjectID, companyModel.CompanyID)
			log.Warn(sigMsg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: sigMsg,
					Code:    "400",
				})
		}

		claCompanyManager := &models.CompanyClaManager{
			LfUsername:   user.Username,
			LogoURL:      user.LogoURL,
			Email:        *user.Email,
			UserSfid:     user.ID,
			ApprovedOn:   time.Now().String(),
			ProjectSfid:  params.ProjectSFID,
			ClaGroupName: claGroup.ProjectName,
			ProjectID:    claGroup.ProjectID,
			ProjectName:  projectSF.Name,
		}
		return cla_manager.NewCreateCLAManagerOK().WithPayload(claCompanyManager)

	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerUnauthorized()
		}

		// Get user by firstname,lastname and email parameters
		userServiceClient := v2UserService.GetClient()
		user, userErr := userServiceClient.GetUserByUsername(params.UserLFID)

		if userErr != nil {
			msg := fmt.Sprintf("Failed to get user when searching by username: %s , error: %v ", params.UserLFID, userErr)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Search for salesForce Company aka external Company
		companyModel, companyErr := companyService.GetCompanyByExternalID(params.CompanySFID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessageDelete(params, companyErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Search for projects by ProjectSFID
		projects, projectErr := projectService.GetProjectsByExternalID(&v1ProjectParams.GetProjectsByExternalIDParams{
			ProjectSFID: params.ProjectSFID,
		})

		var claGroup *v1Models.Project
		// Get unique project by passed CLAGroup ID parameter
		for _, proj := range projects.Projects {
			if proj.ProjectID == params.ProjectID {
				claGroup = &proj
				break
			}
		}

		if projectErr != nil || claGroup == nil {
			msg := fmt.Sprintf("Project lookup failed for projectID: %s and/or projectSFID:%s ", params.ProjectID, params.ProjectSFID)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}
		acsClient := v2AcsService.GetClient()

		roleID, roleErr := acsClient.GetRoleID("cla-manager")
		if roleErr != nil {
			msg := buildErrorMessageDelete(params, roleErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}
		log.Debugf("Role ID for cla-manager-role : %s", roleID)

		// Get Scope ID
		orgClient := v2OrgService.GetClient()
		scopeID, scopeErr := orgClient.GetScopeID(params.CompanySFID, "cla-manager", "project|organization")
		if scopeErr != nil {
			msg := buildErrorMessageDelete(params, scopeErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		deleteErr := orgClient.DeleteOrgUserRoleOrgScopeProjectOrg(params.CompanySFID, roleID, scopeID, user.Username, *user.Email)
		if deleteErr != nil {
			msg := buildErrorMessageDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		signature, deleteErr := managerService.RemoveClaManager(companyModel.CompanyID, params.ProjectID, params.UserLFID)

		if deleteErr != nil {
			msg := buildErrorMessageDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}
		if signature == nil {
			msg := fmt.Sprintf("Not found signature for project: %s and company: %s ", params.ProjectID, companyModel.CompanyID)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}
		return cla_manager.NewDeleteCLAManagerOK()

	})
}

// buildErrorMessageCreate helper function to build an error message
func buildErrorMessageCreate(params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company SFID: %s, project SFID: %s, firstName: %s, lastName: %s, user email: %s, error: %+v",
		params.CompanySFID, params.ProjectSFID, params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, err)
}

// buildErrorMessage helper function to build an error message
func buildErrorMessageDelete(params cla_manager.DeleteCLAManagerParams, err error) string {
	return fmt.Sprintf("problem deleting new CLA Manager Request using company SFID: %s, project SFID: %s, user ID: %s, error: %+v",
		params.CompanySFID, params.ProjectSFID, params.UserLFID, err)
}

func isUserAuthorizedForProjectOrganization(user *auth.User, externalProjectID, externalCompanyID string) bool {
	if !user.Allowed || !user.IsUserAuthorizedByProject(externalProjectID, externalCompanyID) {
		return false
	}
	return true
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(errPrefix string, params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("%s - problem creating new CLA Manager Request using company SFID: %s, project ID: %s, first name: %s, last name: %s, user email: %s, error: %+v",
		errPrefix, params.CompanySFID, params.ProjectID, params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, err)
}

func sendEmailToUserWithNoLFID(projectModel *v1Models.Project, managerEmail string, designeeName string, designeeEmail string) {
	projectName := projectModel.ProjectName
	// subject string, body string, recipients []string
	subject := fmt.Sprint("EasyCLA: Invitation to create LFID and complete process of becoming CLA Manager")
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
<p>This is a notification email from EasyCLA regarding the Project %s in the EasyCLA system.</p>
<p>User %s was trying to add you as a CLA Manager for Project %s in the EasyCLA system </p>
<p>In order to become CLA Manager, you will need to create a LFID </p>
<p>Please create a LFID by following this link <a href="https://identity.linuxfoundation.org/" target="_blank"> and let the user <emailid> know your new LFID %s </p>
<p>Then user %s will be able to add you as a CLA Manager.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>`, designeeName, projectName, managerEmail, projectName, managerEmail, managerEmail)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
