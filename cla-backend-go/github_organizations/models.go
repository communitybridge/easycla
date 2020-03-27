package github_organizations

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

type GithubOrganization struct {
	DateCreated                string `json:"date_created,omitempty"`
	DateModified               string `json:"date_modified,omitempty"`
	OrganizationInstallationID int64  `json:"organization_installation_id,omitempty"`
	OrganizationName           string `json:"organization_name,omitempty"`
	OrganizationSfid           string `json:"organization_sfid,omitempty"`
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
	}

}

func toModels(input []*GithubOrganization) []*models.GithubOrganization {
	out := make([]*models.GithubOrganization, 0)
	for _, in := range input {
		out = append(out, toModel(in))
	}
	return out
}
