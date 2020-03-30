package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
)

// GetInstallationRepositories returns list of repositories for github app installation
func GetInstallationRepositories(installationID int64) ([]string, error) {
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
	result := make([]string, 0)
	for _, repo := range repos {
		if repo.FullName != nil {
			result = append(result, *repo.FullName)
		}
	}
	return result, nil
}
