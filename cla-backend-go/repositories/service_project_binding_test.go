// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"testing"

	"github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AddGithubRepository must not re-bind a repository record registered under another project (lfx-self-serve-ops#751)
func TestRepositoryBoundToOtherProject(t *testing.T) {
	const authorizedProject = "a0941000002wBz4AAE"

	t.Run("record registered under the authorized project is reusable", func(t *testing.T) {
		assert.NoError(t, repositoryBoundToOtherProject(&models.GithubRepository{
			RepositoryName:        "org/repo",
			RepositoryProjectSfid: authorizedProject,
			RepositorySfdcID:      authorizedProject,
		}, authorizedProject))
	})

	t.Run("legacy record without any project binding is reusable", func(t *testing.T) {
		assert.NoError(t, repositoryBoundToOtherProject(&models.GithubRepository{RepositoryName: "org/repo"}, authorizedProject))
		assert.NoError(t, repositoryBoundToOtherProject(nil, authorizedProject))
	})

	t.Run("record registered under another project is refused", func(t *testing.T) {
		err := repositoryBoundToOtherProject(&models.GithubRepository{
			RepositoryName:        "org/repo",
			RepositoryProjectSfid: "a0941000002wBz4ZZZ",
			RepositorySfdcID:      "a0941000002wBz4ZZZ",
		}, authorizedProject)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already registered under project a0941000002wBz4ZZZ")
	})

	t.Run("either stored binding alone is enough to refuse", func(t *testing.T) {
		require.Error(t, repositoryBoundToOtherProject(&models.GithubRepository{RepositorySfdcID: "a0941000002wBz4ZZZ"}, authorizedProject))
		require.Error(t, repositoryBoundToOtherProject(&models.GithubRepository{RepositoryProjectSfid: "a0941000002wBz4ZZZ"}, authorizedProject))
	})
}
