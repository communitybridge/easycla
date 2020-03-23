// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package auth

import (
	"errors"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

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
		log.Warnf("SecurityAuth - verify token error: %+v", err)
		return nil, err
	}

	// Get the username from the token claims
	usernameClaim, ok := claims[a.authValidator.usernameClaim]
	if !ok {
		log.Warnf("SecurityAuth - username not found, error: %+v", err)
		return nil, errors.New("username not found")
	}

	username, ok := usernameClaim.(string)
	if !ok {
		log.Warnf("SecurityAuth - invalid username, error: %+v", err)
		return nil, errors.New("invalid username")
	}
	nameClaim, ok := claims[a.authValidator.nameClaim]
	if !ok {
		log.Warnf("SecurityAuth - name not found, error: %+v", err)
		return nil, errors.New("name not found")
	}
	name, ok := nameClaim.(string)
	if !ok {
		log.Warnf("SecurityAuth - invalid name, error: %+v", err)
		return nil, errors.New("invalid name")
	}
	emailClaim, ok := claims[a.authValidator.emailClaim]
	if !ok {
		log.Warnf("SecurityAuth - email not found, error: %+v", err)
		return nil, errors.New("email not found")
	}
	email, ok := emailClaim.(string)
	if !ok {
		log.Warnf("SecurityAuth - invalid email, error: %+v", err)
		return nil, errors.New("invalid email")
	}
	// Get User by LFID
	lfuser, err := a.userPermissioner.GetUserAndProfilesByLFID(username)
	if err != nil {
		log.Warnf("SecurityAuth - GetUserAndProfilesByLFID error for username: %s, error: %+v", username, err)
		if err.Error() == "user not found" {
			return &user.CLAUser{
				Name:       name,
				LFEmail:    email,
				LFUsername: username,
			}, nil
		}
		return nil, err
	}

	for _, scope := range scopes {
		switch Scope(scope) {
		case projectScope:
			projectIDs, err := a.userPermissioner.GetUserProjectIDs(lfuser.UserID)
			if err != nil {
				log.Warnf("SecurityAuth - GetUserProjectIDs error for user id: %s, error: %+v", lfuser.UserID, err)
				return nil, err
			}

			lfuser.ProjectIDs = projectIDs
		case companyScope:
			//TODO:  Get all companies for this user
			companies, err := a.userPermissioner.GetUserCompanyIDs(lfuser.UserID)
			if err != nil {
				log.Warnf("SecurityAuth - GetUserCompanyIDs error for user id: %s, error: %+v", lfuser.UserID, err)
				return nil, err
			}

			lfuser.CompanyIDs = companies
		case adminScope:
		}
	}

	return &lfuser, nil
}
