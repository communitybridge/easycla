// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package user

// CLAUser data model
type CLAUser struct {
	UserID         string
	Name           string
	Emails         []string
	LFEmail        string
	LFUsername     string
	LfidProvider   Provider
	GithubProvider Provider
	ProjectIDs     []string
	ClaIDs         []string
	CompanyIDs     []string
}

// Provider data model
type Provider struct {
	ProviderUserID string
}

// IsAuthorizedForProject checks if user have access of the project {
func (claUser *CLAUser) IsAuthorizedForProject(projectSFID string) bool {
	for _, v := range claUser.ProjectIDs {
		if v == projectSFID {
			return true
		}
	}
	return false
}
