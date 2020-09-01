package github

import (
	"context"
	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)


// NewClient creates a new client
func NewClient(githubAccessToken string) *githubpkg.Client {
	return NewClientWithCtx(context.Background(), githubAccessToken)
}

// NewClientWithCtx creates a new client with ctx and accessToken
func NewClientWithCtx(ctx context.Context, githubAccessToken string) *githubpkg.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := githubpkg.NewClient(tc)
	return client
}
