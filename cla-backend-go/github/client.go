package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubAppClient(installationID int64) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, int64(getGithubAppID()), installationID, []byte(getGithubAppPrivateKey()))
	if err != nil {
		return nil, err
	}
	return github.NewClient(&http.Client{Transport: itr}), nil
}

func newGithubOauthClient() *github.Client {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getSecretAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
