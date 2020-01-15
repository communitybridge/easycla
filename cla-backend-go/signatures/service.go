// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// SignatureService interface
type SignatureService interface {
	GetMetrics() (*models.SignatureMetrics, error)
	GetSignatures(ctx context.Context, params signatures.GetSignaturesParams) (*models.Signatures, error)
	GetProjectSignatures(ctx context.Context, params signatures.GetProjectSignaturesParams) (*models.Signatures, error)
	GetProjectCompanySignatures(ctx context.Context, params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error)
	GetProjectCompanyEmployeeSignatures(ctx context.Context, params signatures.GetProjectCompanyEmployeeSignaturesParams) (*models.Signatures, error)
	GetCompanySignatures(ctx context.Context, params signatures.GetCompanySignaturesParams) (*models.Signatures, error)
	GetUserSignatures(ctx context.Context, params signatures.GetUserSignaturesParams) (*models.Signatures, error)

	GetGithubOrganizationsFromWhitelist(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error)
	AddGithubOrganizationToWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
}

type service struct {
	repo                SignatureRepository
	githubOrgValidation bool
}

// NewService creates a new whitelist service
func NewService(repo SignatureRepository, githubOrgValidation bool) SignatureService {
	return service{
		repo,
		githubOrgValidation,
	}
}

// GetMetrics returns signature metrics
func (s service) GetMetrics() (*models.SignatureMetrics, error) {
	return s.repo.GetMetrics()
}

// GetSignatures returns the list of signatures associated with the specified signature ID
func (s service) GetSignatures(ctx context.Context, params signatures.GetSignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	signatureList, err := s.repo.GetSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return signatureList, nil
}

// GetProjectSignatures returns the list of signatures associated with the specified project
func (s service) GetProjectSignatures(ctx context.Context, params signatures.GetProjectSignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	projectSignatures, err := s.repo.GetProjectSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanySignatures returns the list of signatures associated with the specified project
func (s service) GetProjectCompanySignatures(ctx context.Context, params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	projectSignatures, err := s.repo.GetProjectCompanySignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanyEmployeeSignatures returns the list of employee signatures associated with the specified project
func (s service) GetProjectCompanyEmployeeSignatures(ctx context.Context, params signatures.GetProjectCompanyEmployeeSignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	projectSignatures, err := s.repo.GetProjectCompanyEmployeeSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetCompanySignatures returns the list of signatures associated with the specified company
func (s service) GetCompanySignatures(ctx context.Context, params signatures.GetCompanySignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	companySignatures, err := s.repo.GetCompanySignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return companySignatures, nil
}

// GetUserSignatures returns the list of user signatures associated with the specified user
func (s service) GetUserSignatures(ctx context.Context, params signatures.GetUserSignaturesParams) (*models.Signatures, error) {

	// Grab and convert the page size, if defined
	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize

	// If we have a value...attempt to parse it
	if params.PageSize != nil {
		//log.Debugf("page size is not null: %s", *params.PageSize)
		var err error
		pageSize, err = strconv.ParseInt(*params.PageSize, 10, 64)
		if err != nil {
			log.Warnf("error parsing pageSize parameter to int64 - using default size of %d - error: %v",
				defaultPageSize, err)
		}

		// Make sure it's positive
		if pageSize < 1 {
			log.Warnf("invalid page size of %d - must be a positive value - using default size of %d",
				pageSize, defaultPageSize)
			pageSize = defaultPageSize
		}
	}

	// log.Debugf("PageSize: %v", pageSize)

	userSignatures, err := s.repo.GetUserSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return userSignatures, nil
}

// GetGithubOrganizationsFromWhitelist retrieves the organization from the whitelist
func (s service) GetGithubOrganizationsFromWhitelist(ctx context.Context, signatureID string, githubAccessToken string) ([]models.GithubOrg, error) {

	if signatureID == "" {
		msg := fmt.Sprintf("unable to get GitHub organizations whitelist - signature ID is nil")
		log.Warn(msg)
		return nil, errors.New(msg)
	}

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
	organizationID := whiteListParams.OrganizationID

	if signatureID == "" {
		msg := fmt.Sprintf("unable to add GitHub organization from whitelist - signature ID is nil")
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := fmt.Sprintf("unable to add GitHub organization from whitelist - organization ID is nil")
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	// GH_ORG_VALIDATION environment - set to false to test locally which will by-pass the GH auth checks and
	// allow functional tests (e.g. with curl or postmon) - default is enabled

	if s.githubOrgValidation {
		// Verify the authenticated github user has access to the github organization being added.
		if githubAccessToken == "" {
			msg := fmt.Sprintf("unable to add github organization, not logged in using "+
				"signatureID: %s, github organization id: %s",
				signatureID, *organizationID)
			log.Warn(msg)
			return nil, errors.New(msg)
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
			if *org.Login == *organizationID {
				found = true
				break
			}
		}

		if !found {
			msg := fmt.Sprintf("user is not authorized for github organization id: %s", *organizationID)
			log.Warnf(msg)
			return nil, errors.New(msg)
		}
	}

	gitHubWhiteList, err := s.repo.AddGithubOrganizationToWhitelist(signatureID, *organizationID)
	if err != nil {
		log.Warnf("issue adding github organization to white list using signatureID: %s, gh org id: %s, error: %v",
			signatureID, *organizationID, err)
		return nil, err
	}

	return gitHubWhiteList, nil
}

// DeleteGithubOrganizationFromWhitelist deletes the specified GH organization from the whitelist
func (s service) DeleteGithubOrganizationFromWhitelist(ctx context.Context, signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {

	// Extract the payload values
	organizationID := whiteListParams.OrganizationID

	if signatureID == "" {
		msg := fmt.Sprintf("unable to delete GitHub organization from whitelist - signature ID is nil")
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := fmt.Sprintf("unable to delete GitHub organization from whitelist - organization ID is nil")
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	// GH_ORG_VALIDATION environment - set to false to test locally which will by-pass the GH auth checks and
	// allow functional tests (e.g. with curl or postmon) - default is enabled

	if s.githubOrgValidation {
		// Verify the authenticated github user has access to the github organization being added.
		if githubAccessToken == "" {
			msg := fmt.Sprintf("unable to delete github organization, not logged in using "+
				"signatureID: %s, github organization id: %s",
				signatureID, *organizationID)
			log.Warn(msg)
			return nil, errors.New(msg)
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
			if *org.Login == *organizationID {
				found = true
				break
			}
		}

		if !found {
			msg := fmt.Sprintf("user is not authorized for github organization id: %s", *organizationID)
			log.Warnf(msg)
			return nil, errors.New(msg)
		}
	}

	gitHubWhiteList, err := s.repo.DeleteGithubOrganizationFromWhitelist(signatureID, *organizationID)
	if err != nil {
		return nil, err
	}

	return gitHubWhiteList, nil
}
