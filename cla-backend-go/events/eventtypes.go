// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

// event types
const (
	// CreateUser event type
	CreateUser = "Create User"
	// UpdateUser event type
	UpdateUser = "Update User"
	// DeleteUser event type
	DeleteUser = "Delete User"

	// CreateTemplate event type
	CreateTemplate = "Create Template"

	// AddGithubOrgToWL event type
	AddGithubOrgToWL = "Add GH Org To WL"
	// DeleteGithubOrgFromWL event type
	DeleteGithubOrgFromWL = "Delete GH Org From WL"

	// CreateCCLAWhitelistRequest event type
	CreateCCLAWhitelistRequest = "Create CCLA WL Request"
	// DeleteCCLAWhitelistRequest event type
	DeleteCCLAWhitelistRequest = "Delete CCLA WL Request"

	// AddUserToCompanyACL event type
	AddUserToCompanyACL = "Add User To Company ACL"
	// DeleteUserFromCompanyACL event type
	//DeleteUserFromCompanyACL = "Delete User From Company ACL"

	// DeletePendingInvite event type
	DeletePendingInvite    = "Delete Pending Invite"
	AddGithubRepository    = "Add Github Repository"
	DeleteGithubRepository = "Delete Github Repository"
)
