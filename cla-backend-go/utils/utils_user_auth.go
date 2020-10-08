// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"github.com/LF-Engineering/lfx-kit/auth"
)

// IsUserAdmin helper function for determining if the user is an admin
func IsUserAdmin(user *auth.User) bool {
	return user.Admin
}

// IsUserAuthorizedForOrganization helper function for determining if the user is authorized for this company
func IsUserAuthorizedForOrganization(user *auth.User, companySFID string) bool {
	// Previously, we checked for user.Admin - admins should be in a separate role
	// Previously, we checked for user.Allowed, which is currently not used (future flag that is currently not implemented)
	return user.IsUserAuthorizedForOrganizationScope(companySFID)
}

// IsUserAuthorizedForProject helper function for determining if the user is authorized for this project
func IsUserAuthorizedForProject(user *auth.User, projectSFID string) bool {
	// Previously, we checked for user.Admin - admins should be in a separate role
	// Previously, we checked for user.Allowed, which is currently not used (future flag that is currently not implemented)
	return user.IsUserAuthorizedForProjectScope(projectSFID)
}

// IsUserAuthorizedForProjectTree helper function for determining if the user is authorized for this project hierarchy/tree
func IsUserAuthorizedForProjectTree(user *auth.User, projectSFID string) bool {
	// Previously, we checked for user.Admin - admins should be in a separate role
	// Previously, we checked for user.Allowed, which is currently not used (future flag that is currently not implemented)
	return user.IsUserAuthorized(auth.Project, projectSFID, true)
}

// IsUserAuthorizedForProjectOrganization helper function for determining if the user is authorized for this project organization scope
func IsUserAuthorizedForProjectOrganization(user *auth.User, projectSFID, companySFID string) bool {
	// Previously, we checked for user.Allowed, which is currently not used (future flag that is currently not implemented)
	return user.IsUserAuthorizedByProject(projectSFID, companySFID)
}

// IsUserAuthorizedForProjectOrganizationTree helper function for determining if the user is authorized for this project organization scope and nested projects/orgs
func IsUserAuthorizedForProjectOrganizationTree(user *auth.User, projectSFID, companySFID string) bool {
	// Previously, we checked for user.Admin - admins should be in a separate role
	// Previously, we checked for user.Allowed, which is currently not used (future flag that is currently not implemented)
	return user.IsUserAuthorized(auth.ProjectOrganization, projectSFID+"|"+companySFID, true)
}
