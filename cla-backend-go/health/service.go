// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package health

import (
	"context"
	"time"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
)

// Service provides an API to the health API
type Service struct {
	version   string
	commit    string
	branch    string
	buildDate string
}

// New is a simple helper function to create a health service instance
func New(version, commit, branch, buildDate string) Service {
	return Service{
		version:   version,
		commit:    commit,
		branch:    branch,
		buildDate: buildDate,
	}
}

// HealthCheck API call returns the current health of the service
func (s Service) HealthCheck(ctx context.Context, in operations.HealthCheckParams) (*models.Health, error) {
	t := time.Now()

	duration := time.Since(t)

	hs := models.HealthStatus{TimeStamp: time.Now().String(), Healthy: true, Name: "CLA", Duration: duration.String()}

	response := models.Health{
		Status:         "healthy",
		TimeStamp:      time.Now().String(),
		Version:        s.version,
		Githash:        s.commit,
		Branch:         s.branch,
		BuildTimeStamp: s.buildDate,
		Healths:        []*models.HealthStatus{&hs},
	}

	return &response, nil
}
