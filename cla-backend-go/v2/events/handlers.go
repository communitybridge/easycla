// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/go-openapi/runtime"

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
)

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1Events.Service, v1CompanyRepo v1Company.IRepository, projectsClaGroupsRepo projects_cla_groups.Repository) { // nolint
	api.EventsGetRecentEventsHandler = events.GetRecentEventsHandlerFunc(
		func(params events.GetRecentEventsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetRecentEventsHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
			}

			if !utils.IsUserAdmin(authUser) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Recent Events - only Admins allowed to see all events.",
						authUser.UserName),
					XRequestID: reqID,
				})
			}

			result, err := service.GetRecentEvents(params.PageSize)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem fetching recent events")
				return events.NewGetRecentEventsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			// Return an empty list
			if result == nil || len(result.Events) == 0 {
				return events.NewGetRecentEventsOK().WithPayload(&models.EventList{
					Events: []*models.Event{},
				})
			}

			resp, err := v2EventList(result)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem converting events to a v2 object")
				return events.NewGetRecentEventsInternalServerError().WithPayload(errorResponse(reqID, err))
			}

			return events.NewGetRecentEventsOK().WithPayload(resp)
		})

	api.EventsGetFoundationEventsAsCSVHandler = events.GetFoundationEventsAsCSVHandlerFunc(
		func(params events.GetFoundationEventsAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetFoundationEventsAsCSVHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
				"foundationSFID": params.FoundationSFID,
			}

			if !utils.IsUserAuthorizedForProjectTree(authUser, params.FoundationSFID) {
				return WriteResponse(http.StatusForbidden, runtime.JSONMime, runtime.JSONProducer(), &models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Foundation Events for foundation %s.",
						authUser.UserName, params.FoundationSFID),
					XRequestID: reqID,
				})
			}

			result, err := service.GetFoundationEvents(params.FoundationSFID, nil, nil, v1Events.ReturnAllEvents, nil)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem fetching foundation events")
				return WriteResponse(http.StatusBadRequest, runtime.JSONMime, runtime.JSONProducer(), errorResponse(reqID, err))
			}

			filename := fmt.Sprintf("foundation-events-%s.csv", params.FoundationSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
		})

	api.EventsGetFoundationEventsHandler = events.GetFoundationEventsHandlerFunc(
		func(params events.GetFoundationEventsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetFoundationEventsHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
				"foundationSFID": params.FoundationSFID,
			}
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.FoundationSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Foundation Events for foundation %s.",
						authUser.UserName, params.FoundationSFID),
					XRequestID: reqID,
				})
			}

			result, err := service.GetFoundationEvents(params.FoundationSFID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents), params.SearchTerm)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem fetching foundation events")
				return events.NewGetFoundationEventsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			// Return an empty list
			if result == nil || len(result.Events) == 0 {
				return events.NewGetFoundationEventsOK().WithPayload(&models.EventList{
					Events: []*models.Event{},
				})
			}

			resp, err := v2EventList(result)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem converting events to a v2 object")
				return events.NewGetFoundationEventsInternalServerError().WithPayload(errorResponse(reqID, err))
			}
			return events.NewGetFoundationEventsOK().WithPayload(resp)
		})

	api.EventsGetProjectEventsAsCSVHandler = events.GetProjectEventsAsCSVHandlerFunc(
		func(params events.GetProjectEventsAsCSVParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetProjectEventsAsCSVHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
				"projectSFID":    params.ProjectSFID,
			}
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return WriteResponse(http.StatusForbidden, runtime.JSONMime, runtime.JSONProducer(), &models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Events for project %s.",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			pm, err := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
			if err != nil {
				if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
					return WriteResponse(http.StatusBadRequest, runtime.JSONMime, runtime.JSONProducer(), &models.ErrorResponse{
						Code:       "400",
						Message:    fmt.Sprintf("EasyCLA - 400 Bad Request - No cla group associated with this project: %s", params.ProjectSFID),
						XRequestID: reqID,
					})
				}
				return WriteResponse(http.StatusInternalServerError, runtime.JSONMime, runtime.JSONProducer(), errorResponse(reqID, err))
			}

			result, err := service.GetClaGroupEvents(pm.ClaGroupID, nil, nil, v1Events.ReturnAllEvents, nil)
			if err != nil {
				log.WithFields(f).Warnf("problem loading events for CLA Group: %s with ID: %s", pm.ClaGroupName, pm.ClaGroupID)
				return events.NewGetProjectEventsAsCSVBadRequest().WithPayload(errorResponse(reqID, err))
			}

			filename := fmt.Sprintf("project-events-%s.csv", params.ProjectSFID)
			csvResponder := CSVEventsResponse(filename, result)
			return csvResponder
		})

	api.EventsGetProjectEventsHandler = events.GetProjectEventsHandlerFunc(
		func(params events.GetProjectEventsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetProjectEventsHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
				"projectSFID":    params.ProjectSFID,
			}
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !utils.IsUserAuthorizedForProjectTree(authUser, params.ProjectSFID) {
				return events.NewGetRecentEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to Get Project Events for foundation %s.",
						authUser.UserName, params.ProjectSFID),
					XRequestID: reqID,
				})
			}

			// Lookup the CLA Group associated with this Project SFID...
			pm, err := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
			if err != nil {
				if err == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
					// Although the API should view this as a bad request since the project doesn't seem to belong to a
					// CLA Group...just return a successful 200 with an empty list to the caller - nothing to see here, move along.
					return events.NewGetProjectEventsOK().WithPayload(&models.EventList{
						Events: []*models.Event{},
					})
				}
				// Not an error that we are expecting - return an error and give up...
				return events.NewGetProjectEventsInternalServerError().WithPayload(errorResponse(reqID, err))
			}

			// Lookup any events for this CLA Group....
			result, err := service.GetClaGroupEvents(pm.ClaGroupID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents), params.SearchTerm)
			if err != nil {
				log.WithFields(f).Warnf("problem loading events for CLA Group: %s with ID: %s", pm.ClaGroupName, pm.ClaGroupID)
				return events.NewGetProjectEventsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			// Return an empty list
			if result == nil || len(result.Events) == 0 {
				return events.NewGetProjectEventsOK().WithPayload(&models.EventList{
					Events: []*models.Event{},
				})
			}

			resp, err := v2EventList(result)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem converting events to a v2 object")
				return events.NewGetProjectEventsInternalServerError().WithPayload(errorResponse(reqID, err))
			}

			return events.NewGetProjectEventsOK().WithPayload(resp)
		})

	api.EventsGetCompanyProjectEventsHandler = events.GetCompanyProjectEventsHandlerFunc(
		func(params events.GetCompanyProjectEventsParams, authUser *auth.User) middleware.Responder {
			reqID := utils.GetRequestID(params.XREQUESTID)
			ctx := context.WithValue(context.Background(), utils.XREQUESTID, reqID) // nolint
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			f := logrus.Fields{
				"functionName":   "EventsGetCompanyProjectEventsHandler",
				utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
				"authUserName":   authUser.UserName,
				"authUserEmail":  authUser.Email,
				"projectSFID":    params.ProjectSFID,
				"companySFID":    params.CompanySFID,
			}
			if !utils.IsUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return events.NewGetCompanyProjectEventsForbidden().WithPayload(&models.ErrorResponse{
					Code: "403",
					Message: fmt.Sprintf("EasyCLA - 403 Forbidden - user %s does not have access to GetCompanyProject Events with Organization scope of %s",
						authUser.UserName, params.CompanySFID),
					XRequestID: reqID,
				})
			}

			var err error
			psc := v2ProjectService.GetClient()
			projectDetails, err := psc.GetProject(params.ProjectSFID)
			if err != nil {
				log.WithFields(f).Warnf("problem loading project by SFID: %s", params.ProjectSFID)
				return events.NewGetCompanyProjectEventsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			var result *v1Models.EventList
			if projectDetails.ProjectType == utils.ProjectTypeProjectGroup {
				result, err = service.GetCompanyFoundationEvents(params.CompanySFID, params.ProjectSFID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents))
			} else {
				pm, perr := projectsClaGroupsRepo.GetClaGroupIDForProject(params.ProjectSFID)
				if perr != nil {
					if perr == projects_cla_groups.ErrProjectNotAssociatedWithClaGroup {
						// Although the API should view this as a bad request since the project doesn't seem to belong to a
						// CLA Group...just return a successful 200 with an empty list to the caller - nothing to see here, move along.
						return events.NewGetCompanyProjectEventsOK().WithPayload(&models.EventList{
							Events: []*models.Event{},
						})
					}
					log.WithFields(f).WithError(perr).Warnf("problem determining CLA Group for project SFID: %s", params.ProjectSFID)
					return events.NewGetCompanyProjectEventsInternalServerError().WithPayload(errorResponse(reqID, perr))
				}
				result, err = service.GetCompanyClaGroupEvents(params.CompanySFID, pm.ClaGroupID, params.NextKey, params.PageSize, aws.BoolValue(params.ReturnAllEvents))
			}
			if err != nil {
				log.WithFields(f).WithError(err).Warn("problem loading events")
				return events.NewGetCompanyProjectEventsBadRequest().WithPayload(errorResponse(reqID, err))
			}

			// Return an empty list
			if result == nil || len(result.Events) == 0 {
				return events.NewGetCompanyProjectEventsOK().WithPayload(&models.EventList{
					Events: []*models.Event{},
				})
			}

			resp, err := v2EventList(result)
			if err != nil {
				log.WithFields(f).WithError(err).Warnf("problem converting events to a v2 object")
				return events.NewGetCompanyProjectEventsInternalServerError().WithPayload(errorResponse(reqID, err))
			}
			return events.NewGetCompanyProjectEventsOK().WithPayload(resp)
		})
}

// WriteResponse function writes http response.
func WriteResponse(httpStatus int, contentType string, contentProducer runtime.Producer, data interface{}) middleware.Responder {
	return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
		rw.Header().Set(runtime.HeaderContentType, contentType)
		rw.WriteHeader(httpStatus)
		err := contentProducer.Produce(rw, data)
		if err != nil {
			log.Warnf("failed to write data. error = %v", err)
		}
	})
}
