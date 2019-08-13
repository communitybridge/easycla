// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type service struct {
	repo       Repository
	httpClient *http.Client
}

// NewService creates a new whitelist service
func NewService(repo Repository, httpClient *http.Client) service {
	return service{
		repo:       repo,
		httpClient: httpClient,
	}
}

// DeleteGithubOrganizationFromWhitelist deletes the specified GH organization from the whitelist
func (s service) DeleteGithubOrganizationFromWhitelist(ctx context.Context, claGroupID, githubOrganizationID string) error {
	err := s.repo.DeleteGithubOrganizationFromWhitelist(claGroupID, githubOrganizationID)
	if err != nil {
		return err
	}

	return nil
}

// AddGithubOrganizationToWhitelist adds the GH organization to the whitelist
func (s service) AddGithubOrganizationToWhitelist(ctx context.Context, claGroupID, githubOrganizationID, githubAccessToken string) error {
	// Verify the authenticated github user has access to the github organization being added.
	if githubAccessToken == "" {
		log.Warnf("unable to add github organization, not logged in using claGroupID: %s, github organization id: %s",
			claGroupID, githubOrganizationID)
		return errors.New("unable to add github organization, not logged in")
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
		msg := fmt.Sprintf("user is not authorized for github organization id: %s", githubOrganizationID)
		log.Warnf(msg)
		return errors.New(msg)
	}

	err = s.repo.AddGithubOrganizationToWhitelist(claGroupID, githubOrganizationID)
	if err != nil {
		return err
	}

	return nil
}

// GetGithubOrganizationsFromWhitelist retrieves the organization from the whitelist
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
