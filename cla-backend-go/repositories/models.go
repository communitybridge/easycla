// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import "github.com/communitybridge/easycla/cla-backend-go/gen/models"

// GithubRepository represent repositories table
type GithubRepository struct {
	DateCreated                string `json:"date_created,omitempty"`
	DateModified               string `json:"date_modified,omitempty"`
	RepositoryExternalID       string `json:"repository_external_id,omitempty"`
	RepositoryID               string `json:"repository_id,omitempty"`
	RepositoryName             string `json:"repository_name,omitempty"`
	RepositoryOrganizationName string `json:"repository_organization_name,omitempty"`
	RepositoryProjectID        string `json:"repository_project_id,omitempty"`
	RepositorySfdcID           string `json:"repository_sfdc_id,omitempty"`
	RepositoryType             string `json:"repository_type,omitempty"`
	RepositoryURL              string `json:"repository_url,omitempty"`
	Version                    string `json:"version,omitempty"`
}

func (gr *GithubRepository) toModel() *models.GithubRepository {
	return &models.GithubRepository{
		DateCreated:                gr.DateCreated,
		DateModified:               gr.DateModified,
		RepositoryExternalID:       gr.RepositoryExternalID,
		RepositoryID:               gr.RepositoryID,
		RepositoryName:             gr.RepositoryName,
		RepositoryOrganizationName: gr.RepositoryOrganizationName,
		RepositoryProjectID:        gr.RepositoryProjectID,
		RepositorySfdcID:           gr.RepositorySfdcID,
		RepositoryType:             gr.RepositoryType,
		RepositoryURL:              gr.RepositoryURL,
		Version:                    gr.Version,
	}
}
