// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package whitelist

import (
	"context"
	"errors"
	"net/http"

	"github.com/communitybridge/easy-cla/cla-backend-go/gen/models"

	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type service struct {
	repo       Repository
	httpClient *http.Client
}

func NewService(repo Repository, httpClient *http.Client) service {
	return service{
		repo:       repo,
		httpClient: httpClient,
	}
}
func (s service) DeleteGithubOrganizationFromWhitelist(ctx context.Context, claGroupID, githubOrganizationID string) error {
	err := s.repo.DeleteGithubOrganizationFromWhitelist(claGroupID, githubOrganizationID)
	if err != nil {
		return err
	}

	return nil
}

func (s service) AddGithubOrganizationToWhitelist(ctx context.Context, claGroupID, githubOrganizationID, githubAccessToken string) error {
	// Verify the authenticated github user has access to the github organization being added.
	if githubAccessToken == "" {
		return errors.New("unable to add github organizaiton, not logged in")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubpkg.NewClient(tc)

	opt := &githubpkg.ListOptions{
		PerPage: 100,
	}

	orgs, _, err := client.Organizations.List(ctx, "", opt)
	if err != nil {
		return err
	}

	found := false
	for _, org := range orgs {
		if *org.Login == githubOrganizationID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("user is not authorized")
	}

	err = s.repo.AddGithubOrganizationToWhitelist(claGroupID, githubOrganizationID)
	if err != nil {
		return err
	}

	return nil
}

func (s service) GetGithubOrganizationsFromWhitelist(ctx context.Context, claGroupID, githubAccessToken string) ([]models.GithubOrg, error) {
	orgIds, err := s.repo.GetGithubOrganizationsFromWhitelist(claGroupID)
	if err != nil {
		return nil, err
	}

	if githubAccessToken != "" {
		selectedOrgs := make(map[string]struct{}, len(orgIds))
		for _, selectedOrg := range orgIds {
			selectedOrgs[*selectedOrg.ID] = struct{}{}
		}

		// Since we're logged into github, lets get the list of organizaiton we can add.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		orgs, _, err := client.Organizations.List(ctx, "", opt)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			_, ok := selectedOrgs[*org.Login]
			if ok {
				continue
			}

			orgIds = append(orgIds, models.GithubOrg{ID: org.Login})
		}
	}

	return orgIds, nil
}
