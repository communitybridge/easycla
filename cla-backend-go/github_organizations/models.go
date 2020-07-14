package github_organizations

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

// GithubOrganization is data model for github organizations
type GithubOrganization struct {
	DateCreated                string `json:"date_created,omitempty"`
	DateModified               string `json:"date_modified,omitempty"`
	OrganizationInstallationID int64  `json:"organization_installation_id,omitempty"`
	OrganizationName           string `json:"organization_name,omitempty"`
	OrganizationNameLower      string `json:"organization_name_lower,omitempty"`
	OrganizationSfid           string `json:"organization_sfid,omitempty"`
	ProjectSFID                string `json:"project_sfid"`
	AutoEnabled                bool   `json:"auto_enabled"`
	Version                    string `json:"version,omitempty"`
}

func toModel(in *GithubOrganization) *models.GithubOrganization {
	return &models.GithubOrganization{
		DateCreated:                in.DateCreated,
		DateModified:               in.DateModified,
		OrganizationInstallationID: in.OrganizationInstallationID,
		OrganizationName:           in.OrganizationName,
		OrganizationSfid:           in.OrganizationSfid,
		Version:                    in.Version,
		AutoEnabled:                in.AutoEnabled,
		ProjectSFID:                in.ProjectSFID,
	}
}

func toModels(input []*GithubOrganization) []*models.GithubOrganization {
	out := make([]*models.GithubOrganization, 0)
	for _, in := range input {
		out = append(out, toModel(in))
	}
	return out
}
