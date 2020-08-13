// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"errors"
	"fmt"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/sign"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"
	"github.com/go-openapi/runtime/middleware"
)

// Configure API call
func Configure(api *operations.EasyclaAPI, service Service) {
	// Retrieve a list of available templates
	api.SignRequestCorporateSignatureHandler = sign.RequestCorporateSignatureHandlerFunc(
		func(params sign.RequestCorporateSignatureParams, user *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectOrganizationTree(user, utils.StringValue(params.Input.ProjectSfid), utils.StringValue(params.Input.CompanySfid)) {
				return sign.NewRequestCorporateSignatureForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Request Corporate Signature with Project|Organization scope of %s | %s",
						user.UserName, utils.StringValue(params.Input.ProjectSfid), utils.StringValue(params.Input.CompanySfid)),
				})
			}

			resp, err := service.RequestCorporateSignature(utils.StringValue(params.XUSERNAME), params.Authorization, params.Input)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					return sign.NewRequestCorporateSignatureNotFound().WithPayload(errorResponse(err))
				}
				if strings.Contains(err.Error(), "internal server error") {
					return sign.NewRequestCorporateSignatureInternalServerError().WithPayload(errorResponse(err))
				}
				if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
					return sign.NewRequestCorporateSignatureNotFound().WithPayload(errorResponse(err))
				}
				if err == ErrCCLANotEnabled || err == ErrTemplateNotConfigured {
					return sign.NewRequestCorporateSignatureBadRequest().WithPayload(errorResponse(err))
				}
				if _, ok := err.(*organizations.ListOrgUsrAdminScopesNotFound); ok {
					formatErr := errors.New("user role scopes not found for cla-signatory role ")
					return sign.NewRequestCorporateSignatureNotFound().WithPayload(errorResponse(formatErr))
				}
				if _, ok := err.(*organizations.CreateOrgUsrRoleScopesConflict); ok {
					formatErr := errors.New("user role scope conflict")
					return sign.NewRequestCorporateSignatureConflict().WithPayload(errorResponse(formatErr))
				}
				if err == ErrNotInOrg {
					return sign.NewRequestCorporateSignatureConflict().WithPayload(errorResponse(err))
				}
				return sign.NewRequestCorporateSignatureBadRequest().WithPayload(errorResponse(err))
			}
			return sign.NewRequestCorporateSignatureOK().WithPayload(resp)
		})
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
