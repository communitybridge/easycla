// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/events"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/users"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// errors
var (
	ErrCclaWhitelistRequestAlreadyExists = errors.New("CCLA whiltelist request already exist")
)

type service struct {
	repo          Repository
	userRepo      users.Repository
	companyRepo   company.RepositoryService
	projectRepo   project.DynamoRepository
	eventsService events.Service
	httpClient    *http.Client
}

// NewService creates a new whitelist service
func NewService(repo Repository, userRepo users.Repository, companyRepo company.RepositoryService, projectRepo project.DynamoRepository, eventsService events.Service, httpClient *http.Client) service {
	return service{
		repo:          repo,
		userRepo:      userRepo,
		companyRepo:   companyRepo,
		projectRepo:   projectRepo,
		eventsService: eventsService,
		httpClient:    httpClient,
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

	log.Debugf("querying for user's github organizations...")
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
		log.Warnf("issue adding github organization to white list using claGroupID: %s, gh org id: %s, error: %v",
			claGroupID, githubOrganizationID, err)
		return err
	}

	return nil
}

// GetGithubOrganizationsFromWhitelist retrieves the organization from the whitelist
func (s service) GetGithubOrganizationsFromWhitelist(ctx context.Context, claGroupID, githubAccessToken string) ([]models.GithubOrg, error) {
	orgIds, err := s.repo.GetGithubOrganizationsFromWhitelist(claGroupID)
	if err != nil {
		log.Warnf("error loading github organization from whitelist using id: %s, error: %v",
			claGroupID, err)
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

func (s service) AddCclaWhitelistRequest(companyID string, projectID string, args models.CclaWhitelistRequestInput) (string, error) {
	list, err := s.repo.ListCclaWhitelistRequest(companyID, &projectID, &args.UserID)
	if err != nil {
		return "", err
	}
	if len(list.List) > 0 {
		return "", ErrCclaWhitelistRequestAlreadyExists
	}
	company, err := s.companyRepo.GetCompany(companyID)
	if err != nil {
		return "", err
	}
	project, err := s.projectRepo.GetProject(projectID)
	if err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUser(args.UserID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid user")
	}
	return s.repo.AddCclaWhitelistRequest(company, project, user)
}

func (s service) DeleteCclaWhitelistRequest(requestID string) error {
	return s.repo.DeleteCclaWhitelistRequest(requestID)
}

func (s service) ListCclaWhitelistRequest(companyID string, projectID *string) (*models.CclaWhitelistRequestList, error) {
	return s.repo.ListCclaWhitelistRequest(companyID, projectID, nil)
}
