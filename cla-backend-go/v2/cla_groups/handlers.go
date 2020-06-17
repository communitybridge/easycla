package cla_groups

import (
	"fmt"
	"strings"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_group"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure configures the cla group api
func Configure(api *operations.EasyclaAPI, service Service) {
	api.ClaGroupCreateClaGroupHandler = cla_group.CreateClaGroupHandlerFunc(func(params cla_group.CreateClaGroupParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProject(authUser, params.ClaGroupInput.FoundationSfid) {
			return cla_group.NewCreateClaGroupForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to CreateCLAGroup with Project scope of %s",
					authUser.UserName, params.ClaGroupInput.FoundationSfid),
			})
		}

		claGroup, err := service.CreateCLAGroup(params.ClaGroupInput, utils.StringValue(params.XUSERNAME))
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewCreateClaGroupBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 %s", err.Error()),
				})
			}
			return cla_group.NewCreateClaGroupInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}

		return cla_group.NewCreateClaGroupOK().WithPayload(claGroup)
	})

	api.ClaGroupValidateClaGroupHandler = cla_group.ValidateClaGroupHandlerFunc(func(params cla_group.ValidateClaGroupParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)

		// No API user validation - anyone can confirm or use the validate API endpoint

		valid, validationErrors := service.ValidateCLAGroup(params.ValidationInputRequest)
		return cla_group.NewValidateClaGroupOK().WithPayload(&models.ClaGroupValidationResponse{
			Valid:            valid,
			ValidationErrors: validationErrors,
		})
	})
}
