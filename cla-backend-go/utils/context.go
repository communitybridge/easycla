// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"context"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/gofrs/uuid"
)

// NewContext returns a new context object with a new request ID
func NewContext() context.Context {
	requestID, err := uuid.NewV4()
	if err != nil {
		log.Warnf("unable to generate a UUID for x-request-id, error: %v", err)
		return context.Background()
	}

	return context.WithValue(context.Background(), XREQUESTID, requestID.String()) // nolint
}
