// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/LF-Engineering/lfx-kit/auth"

	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/go-openapi/runtime/middleware"
)

// Configure is the API handler routine for CLA Manager routes
func Configure(api *operations.EasyclaAPI, service Service, LfxPortalURL string) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerForbidden()
		}

		compCLAManager, errorResponse := service.CreateCLAManager(params, authUser.Email)

		if errorResponse != nil {
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse)
		}

		return cla_manager.NewCreateCLAManagerOK().WithPayload(compCLAManager)

	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {

		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewDeleteCLAManagerForbidden()
		}

		errResponse := service.DeleteCLAManager(params)

		if errResponse != nil {
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(errResponse)
		}

		return cla_manager.NewDeleteCLAManagerNoContent()

	})

	api.ClaManagerCreateCLAManagerDesigneeHandler = cla_manager.CreateCLAManagerDesigneeHandlerFunc(func(params cla_manager.CreateCLAManagerDesigneeParams, authUser *auth.User) middleware.Responder {
		if !authUser.IsUserAuthorizedForOrganizationScope(params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerDesigneeForbidden()
		}

		claManagerDesignee, err := service.CreateCLAManagerDesignee(params.CompanySFID, params.ProjectSFID, params.Body.UserEmail)

		if err != nil {
			msg := fmt.Sprintf("Problem creating cla Manager Designee for user :%s, error: %+v ", authUser.Email, err)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		log.Debugf("CLA Manager designee created : %+v", claManagerDesignee)
		return cla_manager.NewCreateCLAManagerDesigneeOK().WithPayload(claManagerDesignee)

	})

	api.ClaManagerCreateCLAManagerRequestHandler = cla_manager.CreateCLAManagerRequestHandlerFunc(func(params cla_manager.CreateCLAManagerRequestParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerRequestForbidden()
		}
		orgService := v2OrgService.GetClient()

		// GetSFProject
		ps := v2ProjectService.GetClient()
		projectSF, projectErr := ps.GetProject(params.ProjectSFID)
		if projectErr != nil {
			msg := fmt.Sprintf("Project service lookup error for SFID: %s, error : %+v", params.ProjectSFID, projectErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		userService := v2UserService.GetClient()
		user, err := userService.SearchUserByEmail(params.Body.UserEmail)
		if err != nil {
			msg := fmt.Sprintf("Problem getting user for userEmail: %s , error: %+v", params.Body.UserEmail, err)
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Check if sending cla manager request to company admin
		if params.Body.ContactAdmin {
			log.Debugf("Sending email to company Admin")
			scopes, listScopeErr := orgService.ListOrgUserAdminScopes(params.CompanySFID)
			if listScopeErr != nil {
				msg := fmt.Sprintf("Admin lookup error for organisation SFID: %s ", params.CompanySFID)
				return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    "400",
					})
			}
			for _, admin := range scopes.Userroles {
				sendEmailToOrgAdmin(admin.Contact.EmailAddress, admin.Contact.Name, projectSF.Name, authUser.Email, authUser.UserName, LfxPortalURL)
			}
			return cla_manager.NewCreateCLAManagerRequestNoContent()
		}

		claManagerDesignee, err := service.CreateCLAManagerDesignee(params.CompanySFID, params.ProjectSFID, params.Body.UserEmail)

		if err != nil {
			msg := fmt.Sprintf("Problem creating cla Manager Designee for user :%s, error: %+v ", params.Body.UserEmail, err)
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Search for salesForce Company aka external Company
		companyModel, companyErr := orgService.GetOrganization(params.CompanySFID)
		if companyErr != nil || companyModel == nil {
			msg := fmt.Sprintf("Problem getting company by SFID: %s", params.CompanySFID)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		log.Debugf("Sending Email to CLA Manager Designee email: %s ", params.Body.UserEmail)

		sendEmailToCLAManagerDesignee(LfxPortalURL, projectSF.Name, params.Body.UserEmail, user.Name, authUser.Email, authUser.UserName)

		log.Debugf("CLA Manager designee created : %+v", claManagerDesignee)
		return cla_manager.NewCreateCLAManagerRequestOK().WithPayload(claManagerDesignee)
	})
}

// buildErrorMessageCreate helper function to build an error message
func buildErrorMessageCreate(params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company SFID: %s, project SFID: %s, firstName: %s, lastName: %s, user email: %s, error: %+v",
		params.CompanySFID, params.ProjectSFID, *params.Body.FirstName, *params.Body.LastName, *params.Body.UserEmail, err)
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
