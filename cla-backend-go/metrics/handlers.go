package metrics

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/metrics"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.MetricsGetMetricsHandler = metrics.GetMetricsHandlerFunc(
		func(params metrics.GetMetricsParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetMetrics(params)
			if err != nil {
				return metrics.NewGetMetricsBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetMetricsOK().WithPayload(result)
		})

	api.MetricsGetClaManagerDistributionHandler = metrics.GetClaManagerDistributionHandlerFunc(
		func(params metrics.GetClaManagerDistributionParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetCLAManagerDistribution(params)
			if err != nil {
				return metrics.NewGetClaManagerDistributionBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetClaManagerDistributionOK().WithPayload(result)
		})

	api.MetricsGetTotalCountHandler = metrics.GetTotalCountHandlerFunc(
		func(params metrics.GetTotalCountParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetTotalCountMetrics()
			if err != nil {
				return metrics.NewGetTotalCountBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetTotalCountOK().WithPayload(result)
		})

	api.MetricsGetCompanyMetricHandler = metrics.GetCompanyMetricHandlerFunc(
		func(params metrics.GetCompanyMetricParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetCompanyMetric(params.CompanyID)
			if err != nil {
				return metrics.NewGetCompanyMetricBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetCompanyMetricOK().WithPayload(result)
		})

	api.MetricsGetProjectMetricHandler = metrics.GetProjectMetricHandlerFunc(
		func(params metrics.GetProjectMetricParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetProjectMetric(params.ProjectID)
			if err != nil {
				return metrics.NewGetProjectMetricBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetProjectMetricOK().WithPayload(result)
		})
	api.MetricsGetTopCompaniesHandler = metrics.GetTopCompaniesHandlerFunc(
		func(params metrics.GetTopCompaniesParams, claUser *user.CLAUser) middleware.Responder {
			result, err := service.GetTopCompanies()
			if err != nil {
				return metrics.NewGetTopCompaniesBadRequest().WithPayload(errorResponse(err))
			}
			return metrics.NewGetTopCompaniesOK().WithPayload(result)
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
