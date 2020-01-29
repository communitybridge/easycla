// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

const (
	// CreateUser event type
	CreateUser = "CreateUser"
	// UpdateUser event type
	UpdateUser = "UpdateUser"
	// DeleteUser event type
	DeleteUser = "DeleteUser"

	// CreateTemplate event type
	CreateTemplate = "CreateTemplate"

	// AddGithubOrgToWL event type
	AddGithubOrgToWL = "AddGithubOrganizationToWhitelist"
	// DeleteGithubOrgFromWL event type
	DeleteGithubOrgFromWL = "DeleteGithubOrganizationFromWhitelist"

	// CreateCCLAWhitelistRequest event type
	CreateCCLAWhitelistRequest = "CreateCCLAWhitelistRequest"
	// DeleteCCLAWhitelistRequest event type
	DeleteCCLAWhitelistRequest = "DeleteCCLAWhitelistRequest"

	// AddUserToCompanyACL event type
	AddUserToCompanyACL = "AddUserToCompanyACL"
	// DeleteUserFromCompanyACL event type
	//DeleteUserFromCompanyACL = "DeleteUserFromCompanyACL"

	// DeletePendingInvite event type
	DeletePendingInvite = "DeletePendingInvite"
)
