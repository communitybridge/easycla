// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"context"

	"github.com/sirupsen/logrus"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/gofrs/uuid"
)

// NewContext returns a new context object with a new request ID
func NewContext() context.Context {
	f := logrus.Fields{
		"functionName": "utils.NewContext",
	}
	requestID, err := uuid.NewV4()
	if err != nil {
		log.WithFields(f).WithError(err).Warn("unable to generate a UUID for x-request-id")
		return context.Background()
	}

	return context.WithValue(context.Background(), XREQUESTID, requestID.String()) // nolint
}
