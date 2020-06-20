// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/events"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/LF-Engineering/lfx-kit/auth"

	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	v1User "github.com/communitybridge/easycla/cla-backend-go/user"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2OrgService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/go-openapi/runtime/middleware"
)

const (
	//BadRequest error Response code
	BadRequest = "400"
	//Conflict error Response code
	Conflict = "409"
)

// Configure is the API handler routine for CLA Manager routes
func Configure(api *operations.EasyclaAPI, service Service, LfxPortalURL string, projectClaGroupRepo projects_cla_groups.Repository, easyCLAUserRepo v1User.RepositoryService, eventsService events.Service) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManager with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
			})
		}
		cginfo, err := projectClaGroupRepo.GetClaGroupIDForProject(params.ProjectSFID)
		if err != nil {
			if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
				return cla_manager.NewCreateCLAManagerInternalServerError().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - Bad Request. error = %s", "No cla group is associated with this project"),
				})
			}
			return cla_manager.NewCreateCLAManagerInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error. error = %s", err.Error()),
			})
		}
		compCLAManager, errorResponse := service.CreateCLAManager(cginfo.ClaGroupID, params, authUser.UserName, authUser.Email)
		if errorResponse != nil {
			if errorResponse.Code == BadRequest {
				return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse)
			} else if errorResponse.Code == Conflict {
				return cla_manager.NewCreateCLAManagerConflict().WithPayload(errorResponse)
			}
		}

		return cla_manager.NewCreateCLAManagerOK().WithPayload(compCLAManager)
	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProjectOrganization(authUser, params.ProjectSFID, params.CompanySFID) {
			return cla_manager.NewDeleteCLAManagerForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAManager with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
			})
		}
		cginfo, err := projectClaGroupRepo.GetClaGroupIDForProject(params.ProjectSFID)
		if err != nil {
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - Bad Request. No Cla Group associated with ProjectSFID: %s ", params.ProjectSFID),
			})
		}

		errResponse := service.DeleteCLAManager(cginfo.ClaGroupID, params)
		if errResponse != nil {
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(errResponse)
		}

		return cla_manager.NewDeleteCLAManagerNoContent()
	})

	api.ClaManagerCreateCLAManagerDesigneeHandler = cla_manager.CreateCLAManagerDesigneeHandlerFunc(func(params cla_manager.CreateCLAManagerDesigneeParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerDesigneeForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManagerDesignee with Organization scope of %s",
					authUser.UserName, params.CompanySFID),
			})
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

	api.ClaManagerInviteCompanyAdminHandler = cla_manager.InviteCompanyAdminHandlerFunc(func(params cla_manager.InviteCompanyAdminParams) middleware.Responder {

		// Get Contributor details
		user, userErr := easyCLAUserRepo.GetUser(params.UserID)

		if userErr != nil {
			msg := fmt.Sprintf("Problem getting user by ID : %s, error: %+v ", params.UserID, userErr)
			return cla_manager.NewInviteCompanyAdminBadRequest().WithPayload(
				&models.ErrorResponse{
					Code:    "400",
					Message: msg,
				})
		}

		claManagerDesignee, err := service.InviteCompanyAdmin(params.Body.ContactAdmin, params.Body.CompanySFID, params.Body.ProjectSFID, params.Body.UserEmail, &user, LfxPortalURL)

		if err != nil {
			return cla_manager.NewInviteCompanyAdminBadRequest().WithPayload(err)
		}
		// Check if admins succcessfully sent email
		if claManagerDesignee == nil && err == nil {
			return cla_manager.NewInviteCompanyAdminNoContent()
		}

		// successfully created cla manager designee and sent invite
		return cla_manager.NewInviteCompanyAdminOK().WithPayload(claManagerDesignee)
	})

	api.ClaManagerCreateCLAManagerRequestHandler = cla_manager.CreateCLAManagerRequestHandlerFunc(func(params cla_manager.CreateCLAManagerRequestParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
			return cla_manager.NewCreateCLAManagerRequestForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAManagerRequest with Project|Organization scope of %s | %s",
					authUser.UserName, params.ProjectSFID, params.CompanySFID),
			})
		}
		orgService := v2OrgService.GetClient()

		// GetSFProject
		ps := v2ProjectService.GetClient()
		projectSF, projectErr := ps.GetProject(params.ProjectSFID)
		if projectErr != nil {
			msg := fmt.Sprintf("EasyCLA - 400 Bad Request - Project service lookup error for SFID: %s, error : %+v",
				params.ProjectSFID, projectErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		// Search for salesForce Company aka external Company
		companyModel, companyErr := orgService.GetOrganization(params.CompanySFID)
		if companyErr != nil || companyModel == nil {
			msg := fmt.Sprintf("EasyCLA - 400 Bad Request - Problem getting company by SFID: %s, error: %+v",
				params.CompanySFID, companyErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		userService := v2UserService.GetClient()
		user, userErr := userService.SearchUserByEmail(params.Body.UserEmail)
		if userErr != nil {
			msg := fmt.Sprintf("EasyCLA - 400 Bad Request - Problem looking up user with Email: %s , error: %+v", params.Body.UserEmail, userErr)
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
				msg := fmt.Sprintf("EasyCLA - 400 Bad Request - Admin lookup error for organisation SFID: %s, error: %+v ",
					params.CompanySFID, listScopeErr)
				return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    "400",
					})
			}

			if len(scopes.Userroles) == 0 {
				msg := fmt.Sprintf("EasyCLA - 400 Bad Request - No admins for organization SFID: %s",
					params.CompanySFID)
				return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
					&models.ErrorResponse{
						Message: msg,
						Code:    "400",
					})
			}

			for _, admin := range scopes.Userroles {
				sendEmailToOrgAdmin(admin.Contact.EmailAddress, admin.Contact.Name, companyModel.Name, projectSF.Name, authUser.Email, authUser.UserName, LfxPortalURL)
				// Make a note in the event log
				eventsService.LogEvent(&events.LogEventArgs{
					EventType:         events.ContributorNotifyCompanyAdminType,
					LfUsername:        authUser.UserName,
					ExternalProjectID: params.ProjectSFID,
					CompanyID:         companyModel.ID,
					EventData: &events.ContributorNotifyCompanyAdminData{
						AdminName:  admin.Contact.Name,
						AdminEmail: admin.Contact.EmailAddress,
					},
				})
			}

			return cla_manager.NewCreateCLAManagerRequestNoContent()
		}

		claManagerDesignee, err := service.CreateCLAManagerDesignee(params.CompanySFID, params.ProjectSFID, params.Body.UserEmail)
		if err != nil {
			msg := fmt.Sprintf("EasyCLA - 400 Bad Request - problem creating CLA Manager Designee for user :%s, error: %+v ",
				params.Body.UserEmail, err)
			return cla_manager.NewCreateCLAManagerRequestBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		// Make a note in the event log
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:         events.ContributorAssignCLADesigneeType,
			LfUsername:        authUser.UserName,
			ExternalProjectID: params.ProjectSFID,
			CompanyID:         companyModel.ID,
			EventData: &events.ContributorAssignCLADesignee{
				DesigneeName:  claManagerDesignee.LfUsername,
				DesigneeEmail: claManagerDesignee.Email,
			},
		})

		log.Debugf("Sending Email to CLA Manager Designee email: %s ", params.Body.UserEmail)
		sendEmailToCLAManagerDesignee(LfxPortalURL, companyModel.Name, projectSF.Name, params.Body.UserEmail, user.Name, authUser.Email, authUser.UserName)
		// Make a note in the event log
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:         events.ContributorNotifyCLADesigneeType,
			LfUsername:        authUser.UserName,
			ExternalProjectID: params.ProjectSFID,
			CompanyID:         companyModel.ID,
			EventData: &events.ContributorNotifyCLADesignee{
				DesigneeName:  claManagerDesignee.LfUsername,
				DesigneeEmail: claManagerDesignee.Email,
			},
		})

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

// sendEmailToUserWithNoLFID helper function to send email to a given user with no LFID
func sendEmailToUserWithNoLFID(projectModel *v1Models.Project, requesterUsername, requesterEmail, userWithNoLFIDName, userWithNoLFIDEmail string) {
	projectName := projectModel.ProjectName
	// subject string, body string, recipients []string
	subject := fmt.Sprint("EasyCLA: Invitation to create LFID and complete process of becoming CLA Manager")
	recipients := []string{userWithNoLFIDEmail}
	body := fmt.Sprintf(`
<p>Hello %s,</p>
<p>This is a notification email from EasyCLA regarding the Project %s in the EasyCLA system.</p>
<p>User %s (%s) was trying to add you as a CLA Manager for Project %s but was unable to identify your account details in
the EasyCLA system. In order to become a CLA Manager for Project %s, you will need to create a LFID by 
<a href="https://identity.linuxfoundation.org/" target="_blank">following this link</a> and establishing an account.
Once complete, notify the user %s and they will be able to add you as a CLA Manager.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>`,
		userWithNoLFIDName, projectName,
		requesterUsername, requesterEmail, projectName, projectName,
		requesterUsername)

	err := utils.SendEmail(subject, body, recipients)
	if err != nil {
		log.Warnf("problem sending email with subject: %s to recipients: %+v, error: %+v", subject, recipients, err)
	} else {
		log.Debugf("sent email with subject: %s to recipients: %+v", subject, recipients)
	}
}
