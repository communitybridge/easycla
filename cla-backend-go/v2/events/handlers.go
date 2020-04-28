package events

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1Events "github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/events"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1Events.Service) {
	api.EventsGetRecentEventsHandler = events.GetRecentEventsHandlerFunc(
		func(params events.GetRecentEventsParams, user *auth.User) middleware.Responder {
			result, err := service.GetRecentEvents(params.PageSize)
			if err != nil {
				return events.NewGetRecentEventsBadRequest().WithPayload(errorResponse(err))
			}
			return events.NewGetRecentEventsOK().WithPayload(*result)
		})

	api.EventsGetRecentCompanyProjectEventsHandler = events.GetRecentCompanyProjectEventsHandlerFunc(
		func(params events.GetRecentCompanyProjectEventsParams, user *auth.User) middleware.Responder {
			result, err := service.GetRecentEventsForCompanyProject(params.CompanyID, params.ProjectSFID, params.PageSize)
			if err != nil {
				return events.NewGetRecentEventsBadRequest().WithPayload(errorResponse(err))
			}
			return events.NewGetRecentEventsOK().WithPayload(*result)
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
