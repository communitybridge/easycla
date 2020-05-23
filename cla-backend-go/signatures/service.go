// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/signatures"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	githubpkg "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// SignatureService interface
type SignatureService interface {
	GetSignature(signatureID string) (*models.Signature, error)
	GetProjectSignatures(params signatures.GetProjectSignaturesParams) (*models.Signatures, error)
	GetProjectCompanySignature(companyID, projectID string, signed, approved *bool, nextKey *string, pageSize *int64) (*models.Signature, error)
	GetProjectCompanySignatures(params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error)
	GetProjectCompanyEmployeeSignatures(params signatures.GetProjectCompanyEmployeeSignaturesParams) (*models.Signatures, error)
	GetCompanySignatures(params signatures.GetCompanySignaturesParams) (*models.Signatures, error)
	GetUserSignatures(params signatures.GetUserSignaturesParams) (*models.Signatures, error)
	InvalidateProjectRecords(projectID string, projectName string) error

	GetGithubOrganizationsFromWhitelist(signatureID string, githubAccessToken string) ([]models.GithubOrg, error)
	AddGithubOrganizationToWhitelist(signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	DeleteGithubOrganizationFromWhitelist(signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error)
	UpdateApprovalList(projectID, companyID string, params *models.ApprovalList) (*models.Signature, error)

	AddCLAManager(signatureID, claManagerID string) (*models.Signature, error)
	RemoveCLAManager(signatureID, claManagerID string) (*models.Signature, error)
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

// GetSignature returns the signature associated with the specified signature ID
func (s service) GetSignature(signatureID string) (*models.Signature, error) {
	return s.repo.GetSignature(signatureID)
}

// GetProjectSignatures returns the list of signatures associated with the specified project
func (s service) GetProjectSignatures(params signatures.GetProjectSignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	projectSignatures, err := s.repo.GetProjectSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanySignature returns the signature associated with the specified project and company
func (s service) GetProjectCompanySignature(companyID, projectID string, signed, approved *bool, nextKey *string, pageSize *int64) (*models.Signature, error) {
	return s.repo.GetProjectCompanySignature(companyID, projectID, signed, approved, nextKey, pageSize)
}

// GetProjectCompanySignatures returns the list of signatures associated with the specified project
func (s service) GetProjectCompanySignatures(params signatures.GetProjectCompanySignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	signed := true
	approved := true

	projectSignatures, err := s.repo.GetProjectCompanySignatures(
		params.CompanyID, params.ProjectID, &signed, &approved, params.NextKey, &pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetProjectCompanyEmployeeSignatures returns the list of employee signatures associated with the specified project
func (s service) GetProjectCompanyEmployeeSignatures(params signatures.GetProjectCompanyEmployeeSignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	projectSignatures, err := s.repo.GetProjectCompanyEmployeeSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return projectSignatures, nil
}

// GetCompanySignatures returns the list of signatures associated with the specified company
func (s service) GetCompanySignatures(params signatures.GetCompanySignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 50
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	companySignatures, err := s.repo.GetCompanySignatures(params, pageSize, LoadACLDetails)
	if err != nil {
		return nil, err
	}

	return companySignatures, nil
}

// GetUserSignatures returns the list of user signatures associated with the specified user
func (s service) GetUserSignatures(params signatures.GetUserSignaturesParams) (*models.Signatures, error) {

	const defaultPageSize int64 = 10
	var pageSize = defaultPageSize
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	userSignatures, err := s.repo.GetUserSignatures(params, pageSize)
	if err != nil {
		return nil, err
	}

	return userSignatures, nil
}

// GetGithubOrganizationsFromWhitelist retrieves the organization from the whitelist
func (s service) GetGithubOrganizationsFromWhitelist(signatureID string, githubAccessToken string) ([]models.GithubOrg, error) {

	if signatureID == "" {
		msg := "unable to get GitHub organizations whitelist - signature ID is nil"
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
		tc := oauth2.NewClient(context.Background(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		orgs, _, err := client.Organizations.List(context.Background(), "", opt)
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
func (s service) AddGithubOrganizationToWhitelist(signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {
	organizationID := whiteListParams.OrganizationID

	if signatureID == "" {
		msg := "unable to add GitHub organization from whitelist - signature ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := "unable to add GitHub organization from whitelist - organization ID is nil"
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
		tc := oauth2.NewClient(context.Background(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		log.Debugf("querying for user's github organizations...")
		orgs, _, err := client.Organizations.List(context.Background(), "", opt)
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
func (s service) DeleteGithubOrganizationFromWhitelist(signatureID string, whiteListParams models.GhOrgWhitelist, githubAccessToken string) ([]models.GithubOrg, error) {

	// Extract the payload values
	organizationID := whiteListParams.OrganizationID

	if signatureID == "" {
		msg := "unable to delete GitHub organization from whitelist - signature ID is nil"
		log.Warn(msg)
		return nil, errors.New(msg)
	}

	if organizationID == nil {
		msg := "unable to delete GitHub organization from whitelist - organization ID is nil"
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
		tc := oauth2.NewClient(context.Background(), ts)
		client := githubpkg.NewClient(tc)

		opt := &githubpkg.ListOptions{
			PerPage: 100,
		}

		log.Debugf("querying for user's github organizations...")
		orgs, _, err := client.Organizations.List(context.Background(), "", opt)
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

// UpdateApprovalList service method
func (s service) UpdateApprovalList(projectID, companyID string, params *models.ApprovalList) (*models.Signature, error) {
	return s.repo.UpdateApprovalList(projectID, companyID, params)
}

// Disassociate project signatures
func (s service) InvalidateProjectRecords(projectID string, projectName string) error {
	result, err := s.repo.ProjectSignatures(projectID)
	if err != nil {
		log.Warnf(fmt.Sprintf("Unable to get signatures for project : %s", projectID))
		return err
	}
	if len(result.Signatures) > 0 {
		log.Debugf(fmt.Sprintf("Invalidating signatures for project : %s ", projectID))
		for _, signature := range result.Signatures {
			updateErr := s.repo.InvalidateProjectRecord(signature.SignatureID, projectName)
			if updateErr != nil {
				log.Warnf("Unable to update signature :%s , error: %v", signature.SignatureID, updateErr)
			}
		}
	}
	return nil
}

// AddCLAManager adds the specified manager to the signature ACL list
func (s service) AddCLAManager(signatureID, claManagerID string) (*models.Signature, error) {
	return s.repo.AddCLAManager(signatureID, claManagerID)
}

// RemoveCLAManager removes the specified manager from the signature ACL list
func (s service) RemoveCLAManager(signatureID, claManagerID string) (*models.Signature, error) {
	return s.repo.RemoveCLAManager(signatureID, claManagerID)
}
