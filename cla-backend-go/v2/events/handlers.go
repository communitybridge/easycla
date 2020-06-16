// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	v1Events "github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/events"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
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
func Configure(api *operations.EasyclaAPI, service v1Events.Service, v1CompanyRepo v1Company.IRepository) {
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

			result, err := service.GetFoundationSFDCEvents(params.FoundationSFID, params.PageSize)
			if err != nil {
				return events.NewGetFoundationEventsAsCSVBadRequest().WithPayload(errorResponse(err))
			}

			filename := fmt.Sprintf("foundation-events-%s.csv", params.FoundationSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
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

			result, err := service.GetProjectSFDCEvents(params.ProjectSFID, params.PageSize)
			if err != nil {
				return events.NewGetProjectEventsAsCSVBadRequest().WithPayload(errorResponse(err))
			}

			filename := fmt.Sprintf("project-events-%s.csv", params.ProjectSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
		})

	api.EventsGetRecentCompanyProjectEventsHandler = events.GetRecentCompanyProjectEventsHandlerFunc(
		func(params events.GetRecentCompanyProjectEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return events.NewGetRecentCompanyProjectEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to GetRecentCompanyProject Events with Organization scope of %s",
						authUser.UserName, params.CompanySFID),
				})
			}

			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return events.NewGetRecentCompanyProjectEventsNotFound()
				}
			}

			result, err := service.GetRecentEventsForCompanyProject(comp.CompanyID, params.ProjectSFID, params.PageSize)
			if err != nil {
				return events.NewGetRecentCompanyProjectEventsBadRequest().WithPayload(errorResponse(err))
			}

			resp, err := v2EventList(result)
			if err != nil {
				return events.NewGetRecentCompanyProjectEventsInternalServerError().WithPayload(errorResponse(err))
			}
			return events.NewGetRecentCompanyProjectEventsOK().WithPayload(resp)
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
