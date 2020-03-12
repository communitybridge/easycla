package events

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	eventOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.EventsSearchEventsHandler = eventOps.SearchEventsHandlerFunc(
		func(params eventOps.SearchEventsParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.SearchEvents(&params)
			if err != nil {
				log.Debugf("error retrieving events, error: %s", err.Error())
				return eventOps.NewSearchEventsBadRequest().WithPayload(errorResponse(err))
			}
			return eventOps.NewSearchEventsOK().WithPayload(result)
		})
	api.EventsGetRecentEventsHandler = eventOps.GetRecentEventsHandlerFunc(
		func(params eventOps.GetRecentEventsParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetRecentEvents(&params)
			if err != nil {
				log.Debugf("error retrieving events, error: %s", err.Error())
				return eventOps.NewGetRecentEventsBadRequest().WithPayload(errorResponse(err))
			}
			return eventOps.NewGetRecentEventsOK().WithPayload(result)
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
