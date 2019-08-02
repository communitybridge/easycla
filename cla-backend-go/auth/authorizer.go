// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package auth

import (
	"errors"

	"github.com/communitybridge/easycla/cla-backend-go/user"
)

const (
	projectScope Scope = "project"
	companyScope Scope = "company"
	adminScope   Scope = "admin"
)

// Scope string
type Scope string

// UserPermissioner interface methods
type UserPermissioner interface {
	GetUserAndProfilesByLFID(lfidUsername string) (user.CLAUser, error)
	GetUserProjectIDs(userID string) ([]string, error)
	GetClaManagerCorporateClaIDs(userID string) ([]string, error)
	GetUserCompanyIDs(userID string) ([]string, error)
}

// Authorizer data model
type Authorizer struct {
	authValidator    Validator
	userPermissioner UserPermissioner
}

// NewAuthorizer creates a new authorizer based on the specified parameters
func NewAuthorizer(authValidator Validator, userPermissioner UserPermissioner) Authorizer {
	return Authorizer{
		authValidator:    authValidator,
		userPermissioner: userPermissioner,
	}
}

// SecurityAuth creates a new CLA user based on the token and scopes
func (a Authorizer) SecurityAuth(token string, scopes []string) (*user.CLAUser, error) {
	// This handler is called by the runtime whenever a route needs authentication
	// against the 'OAuthSecurity' scheme.
	// It is passed a token extracted from the Authentication Bearer header, and
	// the list of scopes mentioned by the spec for this route.

	// Verify the token is valid
	claims, err := a.authValidator.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	// Get the username from the token claims
	usernameClaim, ok := claims[a.authValidator.usernameClaim]
	if !ok {
		return nil, errors.New("username not found")
	}

	username, ok := usernameClaim.(string)
	if !ok {
		return nil, errors.New("invalid username")
	}

	// Get User by LFID
	user, err := a.userPermissioner.GetUserAndProfilesByLFID(username)
	if err != nil {
		return nil, err
	}

	for _, scope := range scopes {
		switch Scope(scope) {
		case projectScope:
			projectIDs, err := a.userPermissioner.GetUserProjectIDs(user.UserID)
			if err != nil {
				return nil, err
			}

			user.ProjectIDs = projectIDs
		case companyScope:
			//TODO:  Get all companies for this user
			companies, err := a.userPermissioner.GetUserCompanyIDs(user.UserID)
			if err != nil {
				return nil, err
			}

			user.CompanyIDs = companies
		case adminScope:
		}
	}

	return &user, nil
}
