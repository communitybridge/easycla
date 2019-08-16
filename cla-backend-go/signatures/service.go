// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// SignatureService interface
type SignatureService interface {
	GetGithubOrganizationsFromWhitelist(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error)
	AddGithubOrganizationToWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist) error
}

type service struct {
	repo SignatureRepository
}

// NewService creates a new whitelist service
func NewService(repo SignatureRepository) SignatureService {
	return service{
		repo,
	}
}

// GetGithubOrganizationsFromWhitelist retrieves the organization from the whitelist
func (s service) GetGithubOrganizationsFromWhitelist(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error) {
	orgIds, err := s.repo.GetGithubOrganizationsFromWhitelist(signatureID)
	if err != nil {
		log.Warnf("error loading github organization from whitelist using signatureID: %s, error: %v",
			signatureID, err)
		return nil, err
	}

	if githubAccessToken != "" {
		log.Debugf("already authenticated with github - scanning for user's orgs...")

		selectedOrgs := make(map[string]struct{}, len(orgIds))
		for _, selectedOrg := range orgIds {
			selectedOrgs[*selectedOrg.ID] = struct{}{}
		}

		// Since we're logged into github, lets get the list of organization we can add.
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

// AddGithubOrganizationToWhitelist adds the GH organization to the whitelist
func (s service) AddGithubOrganizationToWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {
	organizationID := *whiteListParams.OrganizationID
	// Verify the authenticated github user has access to the github organization being added.
	if githubAccessToken == "" {
		log.Warnf("unable to add github organization, not logged in using signatureID: %s, github organization id: %s",
			signatureID, organizationID)
		return nil, errors.New("unable to add github organization, not logged in")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubpkg.NewClient(tc)

	opt := &githubpkg.ListOptions{
		PerPage: 100,
	}

	log.Debugf("querying for user's github organizations...")
	orgs, _, err := client.Organizations.List(ctx, "", opt)
	if err != nil {
		return nil, err
	}

	found := false
	for _, org := range orgs {
		if *org.Login == organizationID {
			found = true
			break
		}
	}

	if !found {
		msg := fmt.Sprintf("user is not authorized for github organization id: %s", organizationID)
		log.Warnf(msg)
		return nil, errors.New(msg)
	}

	gitHubWhiteList, err := s.repo.AddGithubOrganizationToWhitelist(signatureID, organizationID)
	if err != nil {
		log.Warnf("issue adding github organization to white list using signatureID: %s, gh org id: %s, error: %v",
			signatureID, organizationID, err)
		return nil, err
	}

	return gitHubWhiteList, nil
}

// DeleteGithubOrganizationFromWhitelist deletes the specified GH organization from the whitelist
func (s service) DeleteGithubOrganizationFromWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist) error {
	organizationID := *whiteListParams.OrganizationID
	err := s.repo.DeleteGithubOrganizationFromWhitelist(signatureID, organizationID)
	if err != nil {
		return err
	}

	return nil
}
