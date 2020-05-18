// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	v1ClaManager "github.com/communitybridge/easycla/cla-backend-go/cla_manager"

	"github.com/LF-Engineering/lfx-kit/auth"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_manager"
	"github.com/communitybridge/easycla/cla-backend-go/project"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2UserService "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
)

// Configure is the API handler routine for CLA Manager routes
func Configure(api *operations.EasyclaAPI, managerService v1ClaManager.IService, companyService company.IService, projectService project.Service) {
	api.ClaManagerCreateCLAManagerHandler = cla_manager.CreateCLAManagerHandlerFunc(func(params cla_manager.CreateCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !isUserAuthorizedForProjectOrganization(authUser, params.ProjectID, params.CompanyID) {
			return cla_manager.NewCreateCLAManagerUnauthorized()
		}

		// Get user by firstname,lastname and email parameters
		userServiceClient := v2UserService.GetClient()
		user, userErr := userServiceClient.SearchUsers(params.Body.FirstName, params.Body.LastName, params.Body.UserEmail)

		if userErr != nil {
			msg := fmt.Sprintf("Failed to get user when searching by firstname : %s, lastname: %s , email: %s , error: %v ",
				params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, userErr)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		companyModel, companyErr := companyService.GetCompany(params.CompanyID)
		if companyErr != nil || companyModel == nil {
			msg := buildErrorMessage("company lookup error", params, companyErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(&models.ErrorResponse{
				Message: msg,
				Code:    "400",
			})
		}

		projectModel, projectErr := projectService.GetProjectByID(params.ProjectID)
		if projectErr != nil || projectModel == nil {
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
			sendEmailToUserWithNoLFID(projectModel, authUser.Email, designeeName, params.Body.UserEmail)
		}

		// Add User to Signature ACL
		// Set authUser Scope to exact the cla-manager role
		authUser.Scopes = append(authUser.Scopes, auth.Scope{
			Type:  auth.ProjectOrganization,
			ID:    fmt.Sprintf("%s|%s", params.ProjectID, params.CompanyID),
			Role:  "cla-manager",
			Level: auth.Member,
		})

		log.Debugf("cla-manager role set for user: %+v ", authUser)

		signature, addErr := managerService.AddClaManager(params.CompanyID, params.ProjectID, user.Username)
		if addErr != nil {
			msg := buildErrorMessageCreate(params, addErr)
			log.Warn(msg)
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		v2Signature, err := convertTov2(signature)
		if err != nil {
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse(err))
		}
		return cla_manager.NewCreateCLAManagerOK().WithPayload(v2Signature)

	})

	api.ClaManagerDeleteCLAManagerHandler = cla_manager.DeleteCLAManagerHandlerFunc(func(params cla_manager.DeleteCLAManagerParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		signature, deleteErr := managerService.RemoveClaManager(params.CompanyID, params.ProjectID, params.UserLFID)

		if deleteErr != nil {
			msg := buildErrorMessageDelete(params, deleteErr)
			log.Warn(msg)
			return cla_manager.NewDeleteCLAManagerBadRequest().WithPayload(
				&models.ErrorResponse{
					Message: msg,
					Code:    "400",
				})
		}

		v2Signature, err := convertTov2(signature)
		if err != nil {
			return cla_manager.NewCreateCLAManagerBadRequest().WithPayload(errorResponse(err))
		}

		return cla_manager.NewDeleteCLAManagerOK().WithPayload(v2Signature)

	})
}

func convertTov2(sig *v1Models.Signature) (*models.Signature, error) {
	var dst models.Signature
	err := copier.Copy(&dst, sig)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// buildErrorMessageCreate helper function to build an error message
func buildErrorMessageCreate(params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("problem creating new CLA Manager using company ID: %s, project ID: %s, firstName: %s, lastName: %s, user email: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, err)
}

// buildErrorMessage helper function to build an error message
func buildErrorMessageDelete(params cla_manager.DeleteCLAManagerParams, err error) string {
	return fmt.Sprintf("problem deleting new CLA Manager Request using company ID: %s, project ID: %s, user ID: %s, error: %+v",
		params.CompanyID, params.ProjectID, params.UserLFID, err)
}

type codedResponse interface {
	Code() string
}

func errorResponse(err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	return &e
}

func isUserAuthorizedForProjectOrganization(user *auth.User, externalProjectID, externalCompanyID string) bool {
	if !user.Allowed || !user.IsUserAuthorizedByProject(externalProjectID, externalCompanyID) {
		return false
	}
	return true
}

// buildErrorMessage helper function to build an error message
func buildErrorMessage(errPrefix string, params cla_manager.CreateCLAManagerParams, err error) string {
	return fmt.Sprintf("%s - problem creating new CLA Manager Request using company ID: %s, project ID: %s, first name: %s, last name: %s, user email: %s, error: %+v",
		errPrefix, params.CompanyID, params.ProjectID, params.Body.FirstName, params.Body.LastName, params.Body.UserEmail, err)
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
