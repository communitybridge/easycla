// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	v1Events "github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/events"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
)

func v2EventList(eventList *v1Models.EventList) (*models.EventList, error) {
	var dst models.EventList
	err := copier.Copy(&dst, eventList)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1Events.Service, v1CompanyRepo v1Company.IRepository, projectsClaGroupsRepo projects_cla_groups.Repository) {
	api.EventsGetRecentEventsHandler = events.GetRecentEventsHandlerFunc(
		func(params events.GetRecentEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAdmin(authUser) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Recent Events - only Admins allowed to see all events.",
						authUser.UserName),
				})
			}

			result, err := service.GetRecentEvents(params.PageSize)
			if err != nil {
				return events.NewGetRecentEventsBadRequest().WithPayload(errorResponse(err))
			}

			resp, err := v2EventList(result)
			if err != nil {
				return events.NewGetRecentEventsInternalServerError().WithPayload(errorResponse(err))
			}

			return events.NewGetRecentEventsOK().WithPayload(resp)
		})

	api.EventsGetFoundationEventsAsCSVHandler = events.GetFoundationEventsAsCSVHandlerFunc(
		func(params events.GetFoundationEventsAsCSVParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.FoundationSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Foundation Events for foundation %s.",
						authUser.UserName, params.FoundationSFID),
				})
			}

			result, err := service.GetFoundationEvents(params.FoundationSFID, nil, nil, v1Events.ReturnAllEvents, nil)
			if err != nil {
				return events.NewGetFoundationEventsAsCSVBadRequest().WithPayload(errorResponse(err))
			}

			filename := fmt.Sprintf("foundation-events-%s.csv", params.FoundationSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
		})
	api.EventsGetFoundationEventsHandler = events.GetFoundationEventsHandlerFunc(
		func(params events.GetFoundationEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.FoundationSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Foundation Events for foundation %s.",
						authUser.UserName, params.FoundationSFID),
				})
			}

			result, err := service.GetFoundationEvents(params.FoundationSFID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents), params.SearchTerm)
			if err != nil {
				return events.NewGetFoundationEventsBadRequest().WithPayload(errorResponse(err))
			}
			resp, err := v2EventList(result)
			if err != nil {
				return events.NewGetFoundationEventsInternalServerError().WithPayload(errorResponse(err))
			}
			return events.NewGetFoundationEventsOK().WithPayload(resp)
		})

	api.EventsGetProjectEventsAsCSVHandler = events.GetProjectEventsAsCSVHandlerFunc(
		func(params events.GetProjectEventsAsCSVParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Events for project %s.",
						authUser.UserName, params.ProjectSFID),
				})
			}
			pm, err := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
			if err != nil {
				if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
					return events.NewGetProjectEventsAsCSVNotFound().WithPayload(&models.ErrorResponse{
						Code: "404",
						Message: fmt.Sprintf("EasyCLA - 403 Forbidden - project %s not found in cla",
							params.ProjectSFID),
					})
				}
				return events.NewGetProjectEventsAsCSVInternalServerError().WithPayload(errorResponse(err))
			}
			result, err := service.GetClaGroupEvents(pm.ClaGroupID, nil, nil, v1Events.ReturnAllEvents, nil)
			if err != nil {
				return events.NewGetProjectEventsAsCSVBadRequest().WithPayload(errorResponse(err))
			}

			filename := fmt.Sprintf("project-events-%s.csv", params.ProjectSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
		})

	api.EventsGetProjectEventsHandler = events.GetProjectEventsHandlerFunc(
		func(params events.GetProjectEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProject(authUser, params.ProjectSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Events for foundation %s.",
						authUser.UserName, params.ProjectSFID),
				})
			}

			pm, err := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
			if err != nil {
				if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
					return events.NewGetProjectEventsNotFound().WithPayload(&models.ErrorResponse{
						Code: "404",
						Message: fmt.Sprintf("EasyCLA - 404 Not found - project %s not found in cla",
							params.ProjectSFID),
					})
				}
				return events.NewGetProjectEventsInternalServerError().WithPayload(errorResponse(err))
			}
			result, err := service.GetClaGroupEvents(pm.ClaGroupID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents), params.SearchTerm)
			if err != nil {
				return events.NewGetProjectEventsBadRequest().WithPayload(errorResponse(err))
			}
			resp, err := v2EventList(result)
			if err != nil {
				return events.NewGetProjectEventsInternalServerError().WithPayload(errorResponse(err))
			}
			return events.NewGetProjectEventsOK().WithPayload(resp)
		})

	api.EventsGetCompanyProjectEventsHandler = events.GetCompanyProjectEventsHandlerFunc(
		func(params events.GetCompanyProjectEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return events.NewGetCompanyProjectEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to GetCompanyProject Events with Organization scope of %s",
						authUser.UserName, params.CompanySFID),
				})
			}

			var err error
			psc := v2ProjectService.GetClient()
			projectDetails, err := psc.GetProject(params.ProjectSFID)
			if err != nil {
				return events.NewGetCompanyProjectEventsBadRequest().WithPayload(errorResponse(err))
			}
			var result *v1Models.EventList
			if projectDetails.ProjectType == "Foundation" {
				result, err = service.GetCompanyFoundationEvents(params.CompanySFID, params.ProjectSFID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents))
			} else {
				pm, perr := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
				if perr != nil {
					if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
						return events.NewGetCompanyProjectEventsNotFound().WithPayload(&models.ErrorResponse{
							Code: "404",
							Message: fmt.Sprintf("EasyCLA - 404 Not found - project %s not found in cla",
								params.ProjectSFID),
						})
					}
					return events.NewGetCompanyProjectEventsInternalServerError().WithPayload(errorResponse(err))
				}
				result, err = service.GetCompanyClaGroupEvents(params.CompanySFID, pm.ClaGroupID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents))
			}
			if err != nil {
				return events.NewGetCompanyProjectEventsBadRequest().WithPayload(errorResponse(err))
			}
			resp, err := v2EventList(result)
			if err != nil {
				return events.NewGetCompanyProjectEventsInternalServerError().WithPayload(errorResponse(err))
			}
			return events.NewGetCompanyProjectEventsOK().WithPayload(resp)
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
