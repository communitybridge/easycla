// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package auth

import (
	"errors"

	"github.com/communitybridge/easy-cla/cla-backend-go/user"
)

const (
	projectScope Scope = "project"
	companyScope Scope = "company"
	adminScope   Scope = "admin"
)

type Scope string

type UserPermissioner interface {
	GetUserAndProfilesByLFID(lfidUsername string) (user.CLAUser, error)
	GetUserProjectIDs(userID string) ([]string, error)
	GetClaManagerCorporateClaIDs(userID string) ([]string, error)
	GetUserCompanyIDs(userID string) ([]string, error)
}

type Authorizer struct {
	auth0Validator   Auth0Validator
	userPermissioner UserPermissioner
}

func NewAuthorizer(auth0Validator Auth0Validator, userPermissioner UserPermissioner) Authorizer {
	return Authorizer{
		auth0Validator:   auth0Validator,
		userPermissioner: userPermissioner,
	}
}

func (a Authorizer) SecurityAuth(token string, scopes []string) (*user.CLAUser, error) {
	// This handler is called by the runtime whenever a route needs authentication
	// against the 'OAuthSecurity' scheme.
	// It is passed a token extracted from the Authentication Bearer header, and
	// the list of scopes mentioned by the spec for this route.

	// Verify the token is valid
	claims, err := a.auth0Validator.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	// Get the username from the token claims
	usernameClaim, ok := claims[a.auth0Validator.usernameClaim]
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
