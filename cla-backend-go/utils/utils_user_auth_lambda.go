// +build aws_lambda

// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"github.com/LF-Engineering/lfx-kit/auth"
)

const (
	// ALLOW_ADMIN_SCOPE indicates that a given permissions check allows for admins to access to that resource
	ALLOW_ADMIN_SCOPE = true // nolint
	// DISALLOW_ADMIN_SCOPE indicates that a given permissions check does not allow for admins to access to that resource
	DISALLOW_ADMIN_SCOPE = false // nolint
)

// IsUserAdmin helper function for determining if the user is an admin
func IsUserAdmin(user *auth.User) bool {
	return user.Admin
}

// IsUserAuthorizedForOrganization helper function for determining if the user is authorized for this company
func IsUserAuthorizedForOrganization(user *auth.User, companySFID string, adminScopeAllowed bool) bool {

	if adminScopeAllowed && user.Admin {
		return true
	}

	return user.IsUserAuthorizedForOrganizationScope(companySFID)
}

// IsUserAuthorizedForProjectTree helper function for determining if the user is authorized for this project hierarchy/tree
func IsUserAuthorizedForProjectTree(user *auth.User, projectSFID string, adminScopeAllowed bool) bool {

	if adminScopeAllowed && user.Admin {
		return true
	}

	return user.IsUserAuthorized(auth.Project, projectSFID, true)
}

// IsUserAuthorizedForProject helper function for determining if the user is authorized for this project
func IsUserAuthorizedForProject(user *auth.User, projectSFID string, adminScopeAllowed bool) bool {

	if adminScopeAllowed && user.Admin {
		return true
	}

	return user.IsUserAuthorizedForProjectScope(projectSFID)
}

// IsUserAuthorizedForAnyProjects helper function for determining if the user is authorized for any of the specified projects
func IsUserAuthorizedForAnyProjects(user *auth.User, projectSFIDs []string, adminScopeAllowed bool) bool {
	for _, projectSFID := range projectSFIDs {
		if IsUserAuthorizedForProjectTree(user, projectSFID, adminScopeAllowed) {
			return true
		}
		if IsUserAuthorizedForProject(user, projectSFID, adminScopeAllowed) {
			return true
		}
	}

	return false
}

// IsUserAuthorizedForProjectOrganization helper function for determining if the user is authorized for this project organization scope
func IsUserAuthorizedForProjectOrganization(user *auth.User, projectSFID, companySFID string, adminScopeAllowed bool) bool {

	if adminScopeAllowed && user.Admin {
		return true
	}

	return user.IsUserAuthorizedByProject(projectSFID, companySFID)
}

// IsUserAuthorizedForAnyProjectOrganization helper function for determining if the user is authorized for any of the specified projects with scope of project + organization
func IsUserAuthorizedForAnyProjectOrganization(user *auth.User, projectSFIDs []string, companySFID string, adminScopeAllowed bool) bool {
	for _, projectSFID := range projectSFIDs {
		if IsUserAuthorizedForProjectOrganizationTree(user, projectSFID, companySFID, adminScopeAllowed) {
			return true
		}
		if IsUserAuthorizedForProjectOrganization(user, projectSFID, companySFID, adminScopeAllowed) {
			return true
		}
	}

	return false
}

// IsUserAuthorizedForProjectOrganizationTree helper function for determining if the user is authorized for this project organization scope and nested projects/orgs
func IsUserAuthorizedForProjectOrganizationTree(user *auth.User, projectSFID, companySFID string, adminScopeAllowed bool) bool {

	if adminScopeAllowed && user.Admin {
		return true
	}

	return user.IsUserAuthorized(auth.ProjectOrganization, projectSFID+"|"+companySFID, true)
}
