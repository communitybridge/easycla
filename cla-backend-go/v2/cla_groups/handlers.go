package cla_groups

import (
	"fmt"
	"strings"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/cla_group"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// Configure configures the cla group api
func Configure(api *operations.EasyclaAPI, service Service, v1ProjectService v1Project.Service) {
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

	api.ClaGroupDeleteClaGroupHandler = cla_group.DeleteClaGroupHandlerFunc(func(params cla_group.DeleteClaGroupParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		cg, err := v1ProjectService.GetProjectByID(params.ClaGroupID)
		if err != nil {
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewDeleteClaGroupNotFound().WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewDeleteClaGroupInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		if !utils.IsUserAuthorizedForProject(authUser, cg.FoundationSFID) {
			return cla_group.NewDeleteClaGroupForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAGroup with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
			})
		}

		err = service.DeleteCLAGroup(params.ClaGroupID)
		if err != nil {
			return cla_group.NewDeleteClaGroupInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		return cla_group.NewDeleteClaGroupOK()
	})
	api.ClaGroupEnrollProjectsHandler = cla_group.EnrollProjectsHandlerFunc(func(params cla_group.EnrollProjectsParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		cg, err := v1ProjectService.GetProjectByID(params.ClaGroupID)
		if err != nil {
			if err == v1Project.ErrProjectDoesNotExist {
				return cla_group.NewEnrollProjectsNotFound().WithPayload(&models.ErrorResponse{
					Code: "404",
					Message: fmt.Sprintf("EasyCLA - 404 Not Found - cla_group %s not found",
						params.ClaGroupID),
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "400",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		if !utils.IsUserAuthorizedForProject(authUser, cg.FoundationSFID) {
			return cla_group.NewEnrollProjectsForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to DeleteCLAGroup with Project scope of %s",
					authUser.UserName, cg.FoundationSFID),
			})
		}

		err = service.EnrollProjectsInClaGroup(params.ClaGroupID, cg.FoundationSFID, params.ProjectSFIDList)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return cla_group.NewEnrollProjectsBadRequest().WithPayload(&models.ErrorResponse{
					Code:    "400",
					Message: fmt.Sprintf("EasyCLA - 400 %s", err.Error()),
				})
			}
			return cla_group.NewEnrollProjectsInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		return cla_group.NewEnrollProjectsOK()
	})
	api.ClaGroupListClaGroupsUnderFoundationHandler = cla_group.ListClaGroupsUnderFoundationHandlerFunc(func(params cla_group.ListClaGroupsUnderFoundationParams, authUser *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
		if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
			return cla_group.NewListClaGroupsUnderFoundationForbidden().WithPayload(&models.ErrorResponse{
				Code: "403",
				Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to ListCLAGroupsUnderFoundation with Project scope of %s",
					authUser.UserName, params.ProjectSFID),
			})
		}

		result, err := service.ListClaGroupsUnderFoundation(params.ProjectSFID)
		if err != nil {
			return cla_group.NewListClaGroupsUnderFoundationInternalServerError().WithPayload(&models.ErrorResponse{
				Code:    "500",
				Message: fmt.Sprintf("EasyCLA - 500 Internal server error - error = %s", err.Error()),
			})
		}
		return cla_group.NewListClaGroupsUnderFoundationOK().WithPayload(result)
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
