package events

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1Company "github.com/communitybridge/easycla/cla-backend-go/company"
	v1Events "github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/company"
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
func isUserAuthorizedForOrganization(user *auth.User, externalCompanyID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorized(auth.Organization, externalCompanyID) {
			return false
		}
	}
	return true
}

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1Events.Service, v1CompanyRepo v1Company.IRepository) {
	api.EventsGetRecentEventsHandler = events.GetRecentEventsHandlerFunc(
		func(params events.GetRecentEventsParams, user *auth.User) middleware.Responder {
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

	api.EventsGetRecentCompanyProjectEventsHandler = events.GetRecentCompanyProjectEventsHandlerFunc(
		func(params events.GetRecentCompanyProjectEventsParams, authUser *auth.User) middleware.Responder {
			utils.SetAuthUserProperties(authUser, params.XUSERNAME, params.XEMAIL)
			if !isUserAuthorizedForOrganization(authUser, params.CompanySFID) {
				return events.NewGetRecentCompanyProjectEventsUnauthorized()
			}
			comp, err := v1CompanyRepo.GetCompanyByExternalID(params.CompanySFID)
			if err != nil {
				if err == v1Company.ErrCompanyDoesNotExist {
					return company.NewGetCompanyClaManagersNotFound()
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
