// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import "github.com/LF-Engineering/lfx-kit/auth"

// IsUserAuthorizedForOrganization helper function for determining if the user is authorized for this company
func IsUserAuthorizedForOrganization(user *auth.User, companySFID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorizedForOrganizationScope(companySFID) {
			return false
		}
	}
	return true
}

// IsUserAuthorizedForProject helper function for determining if the user is authorized for this project
func IsUserAuthorizedForProject(user *auth.User, projectSFID string) bool {
	if !user.Admin {
		if !user.Allowed || !user.IsUserAuthorizedForProjectScope(projectSFID) {
			return false
		}
	}
	return true
}

// IsUserAuthorizedForProjectOrganization helper function for determining if the user is authorized for this project organization scope
func IsUserAuthorizedForProjectOrganization(user *auth.User, projectSFID, companySFID string) bool {
	if !user.Allowed || !user.IsUserAuthorizedByProject(projectSFID, companySFID) {
		return false
	}
	return true
}
