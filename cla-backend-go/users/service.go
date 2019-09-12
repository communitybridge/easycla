// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// Service interface
type Service interface {
	GetUser(userID string) (*models.User, error)
}

type service struct {
	repo Repository
}

// NewService creates a new whitelist service
func NewService(repo Repository) Service {
	return service{
		repo,
	}
}

func (s service) GetUser(userID string) (*models.User, error) {
	userModel, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}
