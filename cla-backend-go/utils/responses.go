// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
)

// ErrorResponseBadRequest Helper function to generate a bad request error response
func ErrorResponseBadRequest(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String400,
		Message: fmt.Sprintf("%s - %s", EasyCLA400BadRequest, msg),
	}
}

// ErrorResponseBadRequestWithError Helper function to generate a bad request error response
func ErrorResponseBadRequestWithError(msg string, err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String400,
		Message: fmt.Sprintf("%s - %s - error: %+v", EasyCLA400BadRequest, msg, err),
	}
}

// ErrorResponseForbidden Helper function to generate a forbidden error response
func ErrorResponseForbidden(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String403,
		Message: fmt.Sprintf("%s - %s", EasyCLA403Forbidden, msg),
	}
}

// ErrorResponseForbiddenWithError Helper function to generate a forbidden error response
func ErrorResponseForbiddenWithError(msg string, err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String403,
		Message: fmt.Sprintf("%s - %s - error: %+v", EasyCLA403Forbidden, msg, err),
	}
}

// ErrorResponseNotFound Helper function to generate a not found error response
func ErrorResponseNotFound(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String404,
		Message: fmt.Sprintf("%s - %s", EasyCLA404NotFound, msg),
	}
}

// ErrorResponseNotFoundWithError Helper function to generate a not found error response
func ErrorResponseNotFoundWithError(msg string, err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String404,
		Message: fmt.Sprintf("%s - %s - error: %+v", EasyCLA404NotFound, msg, err),
	}
}

// ErrorResponseInternalServerError Helper function to generate an internal server error response
func ErrorResponseInternalServerError(msg string) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String500,
		Message: fmt.Sprintf("%s - %s", EasyCLA500InternalServerError, msg),
	}
}

// ErrorResponseInternalServerErrorWithError Helper function to generate an internal server error response
func ErrorResponseInternalServerErrorWithError(msg string, err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Code:    String500,
		Message: fmt.Sprintf("%s - %s - error: %+v", EasyCLA500InternalServerError, msg, err),
	}
}
