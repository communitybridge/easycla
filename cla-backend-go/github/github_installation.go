package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/github"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
)

// GetInstallationRepositories returns list of repositories for github app installation
func GetInstallationRepositories(installationID int64) ([]*github.Repository, error) {
	client, err := newGithubAppClient(installationID)
	if err != nil {
		return nil, errors.New("cannot create github client")
	}
	repos, _, err := client.Apps.ListRepos(context.TODO(), nil)
	if err != nil {
		logging.Error("error while getting installation repositories", err)
		err = fmt.Errorf("unable to get repositories for installation id : %d", installationID)
		return nil, err
	}
	return repos, nil
}
