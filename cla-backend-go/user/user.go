// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package user

// CLAUser data model
type CLAUser struct {
	UserID string
	Name   string

	Emails []string

	LFEmail    string
	LFUsername string

	LfidProvider   Provider
	GithubProvider Provider

	ProjectIDs []string
	ClaIDs     []string
	CompanyIDs []string
}

// Provider data model
type Provider struct {
	ProviderUserID string
}
