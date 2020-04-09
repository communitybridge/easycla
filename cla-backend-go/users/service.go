// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package users

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

// Service interface for users
type Service interface {
	CreateUser(user *models.User) (*models.User, error)
	Save(user *models.UserUpdate) (*models.User, error)
	Delete(userID string) error
	GetUser(userID string) (*models.User, error)
	GetUserByUserName(userName string, fullMatch bool) (*models.User, error)
	SearchUsers(field string, searchTerm string, fullMatch bool) (*models.Users, error)
}

type service struct {
	repo UserRepository
}

// NewService creates a new whitelist service
func NewService(repo UserRepository) Service {
	return service{
		repo,
	}
}

// CreateUser attempts to create a new user based on the specified model
func (s service) CreateUser(user *models.User) (*models.User, error) {
	userModel, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

// Save saves/updates the user record
func (s service) Save(user *models.UserUpdate) (*models.User, error) {
	return s.repo.Save(user)
}

// Delete deletes the user record
func (s service) Delete(userID string) error {
	return s.repo.Delete(userID)
}

// GetUser attempts to locate the user by the user id field
func (s service) GetUser(userID string) (*models.User, error) {
	userModel, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

// GetUserByUserName attempts to locate the user by the user name field
func (s service) GetUserByUserName(userName string, fullMatch bool) (*models.User, error) {
	userModel, err := s.repo.GetUserByUserName(userName, fullMatch)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

// SearchUsers attempts to locate the user by the searchField and searchTerm fields
func (s service) SearchUsers(searchField string, searchTerm string, fullMatch bool) (*models.Users, error) {
	userModel, err := s.repo.SearchUsers(searchField, searchTerm, fullMatch)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}
