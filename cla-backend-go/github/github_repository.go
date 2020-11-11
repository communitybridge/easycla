// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/google/go-github/v32/github"
)

// errors
var (
	ErrGithubRepositoryNotFound = errors.New("github repository not found")
)

// GetRepositoryByExternalID finds github repository by github repository id
func GetRepositoryByExternalID(ctx context.Context, installationID, id int64) (*github.Repository, error) {
	client, err := NewGithubAppClient(installationID)
	if err != nil {
		return nil, err
	}
	org, resp, err := client.Repositories.GetByID(ctx, id)
	if err != nil {
		logging.Warnf("GetRepository %v failed. error = %s", id, err.Error())
		if resp.StatusCode == 404 {
			return nil, ErrGithubRepositoryNotFound
		}
		return nil, err
	}
	return org, nil
}
