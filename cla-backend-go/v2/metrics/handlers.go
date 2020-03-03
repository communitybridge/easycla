package metrics

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1MetricsOps "github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/metrics"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/metrics"
	v1Metrics "github.com/communitybridge/easycla/cla-backend-go/metrics"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.EasyclaAPI, service v1Metrics.Service) {
	api.MetricsGetMetricsHandler = metrics.GetMetricsHandlerFunc(
		func(params metrics.GetMetricsParams, user *auth.User) middleware.Responder {
			// TODO: Need to inspect the auth.User roles/permissions to confirm they can query the metrics
			result, err := service.GetMetrics(v1MetricsOps.GetMetricsParams{HTTPRequest: params.HTTPRequest})
			if err != nil {
				return metrics.NewGetMetricsBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetMetricsOK().WithPayload(result)
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
