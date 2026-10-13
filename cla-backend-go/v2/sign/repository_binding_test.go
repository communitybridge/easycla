// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"testing"

	"github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryBoundToCLAGroup(t *testing.T) {
	repo := &models.GithubRepository{Enabled: true, RepositoryClaGroupID: "group-a", RepositoryExternalID: 42}
	org := &models.GithubOrganization{OrganizationName: "acme", OrganizationInstallationID: 7}
	assert.NoError(t, repositoryBoundToCLAGroup(repo, org, 42, 7, "group-a"))
	assert.Error(t, repositoryBoundToCLAGroup(nil, org, 42, 7, "group-a"), "unknown repository")
	assert.Error(t, repositoryBoundToCLAGroup(&models.GithubRepository{Enabled: false, RepositoryClaGroupID: "group-a"}, org, 42, 7, "group-a"), "disabled repository")
	assert.Error(t, repositoryBoundToCLAGroup(repo, org, 42, 7, "group-b"), "a signature of another CLA group must not stamp this repository")
	assert.Error(t, repositoryBoundToCLAGroup(repo, nil, 42, 7, "group-a"), "unknown organization")
	assert.Error(t, repositoryBoundToCLAGroup(repo, &models.GithubOrganization{OrganizationName: "acme", OrganizationInstallationID: 8}, 42, 7, "group-a"), "foreign installation")
}
