// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/swag"

	githubsdk "github.com/google/go-github/github"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/github"

	"github.com/aws/aws-sdk-go/aws"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	v2Models "github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	v1Repositories "github.com/communitybridge/easycla/cla-backend-go/repositories"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
)

// Service contains functions of Github Repository service
type Service interface {
	AddGithubRepository(projectSFID string, input *models.GithubRepositoryInput) (*v1Models.GithubRepository, error)
	EnableRepository(repositoryID string) error
	DisableRepository(repositoryID string) error
	ListProjectRepositories(projectSFID string) (*v1Models.ListGithubRepositories, error)
	GetRepository(repositoryID string) (*v1Models.GithubRepository, error)
	DisableCLAGroupRepositories(claGroupID string) error
	GetProtectedBranch(repositoryID string) (*v2Models.GithubRepositoryBranchProtection, error)
	UpdateProtectedBranch(repositoryID string, input *v2Models.GithubRepositoryBranchProtectionInput) (*v2Models.GithubRepositoryBranchProtection, error)
}

// GithubOrgRepo provide method to get github organization by name
type GithubOrgRepo interface {
	GetGithubOrganizationByName(githubOrganizationName string) (*v1Models.GithubOrganizations, error)
	GetGithubOrganization(githubOrganizationName string) (*v1Models.GithubOrganization, error)
}

type service struct {
	repo                  v1Repositories.Repository
	projectsClaGroupsRepo projects_cla_groups.Repository
	ghOrgRepo             GithubOrgRepo
}

var requiredBranchProtectionChecks = []string{"EasyCLA"}

// NewService creates a new githubOrganizations service
func NewService(repo v1Repositories.Repository, pcgRepo projects_cla_groups.Repository, ghOrgRepo GithubOrgRepo) Service {
	return &service{
		repo:                  repo,
		projectsClaGroupsRepo: pcgRepo,
		ghOrgRepo:             ghOrgRepo,
	}
}

func (s *service) AddGithubRepository(projectSFID string, input *models.GithubRepositoryInput) (*v1Models.GithubRepository, error) {
	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(projectSFID)
	if err != nil {
		return nil, err
	}
	var externalProjectID string
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		externalProjectID = projectSFID
	} else {
		externalProjectID = project.Parent
	}
	allMappings, err := s.projectsClaGroupsRepo.GetProjectsIdsForClaGroup(aws.StringValue(input.ClaGroupID))
	if err != nil {
		return nil, err
	}
	var valid bool
	for _, cgm := range allMappings {
		if cgm.ProjectSFID == projectSFID || cgm.FoundationSFID == projectSFID {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("provided cla group id %s is not linked to project sfid %s", utils.StringValue(input.ClaGroupID), projectSFID)
	}
	org, err := s.ghOrgRepo.GetGithubOrganizationByName(utils.StringValue(input.GithubOrganizationName))
	if err != nil {
		return nil, err
	}
	if len(org.List) == 0 {
		return nil, errors.New("github app not installed on github organization")
	}
	repoGithubID, err := strconv.ParseInt(utils.StringValue(input.RepositoryGithubID), 10, 64)
	if err != nil {
		return nil, err
	}
	ghRepo, err := github.GetRepositoryByExternalID(org.List[0].OrganizationInstallationID, repoGithubID)
	if err != nil {
		return nil, err
	}
	in := &v1Models.GithubRepositoryInput{
		RepositoryExternalID:       input.RepositoryGithubID,
		RepositoryName:             ghRepo.FullName,
		RepositoryOrganizationName: input.GithubOrganizationName,
		RepositoryProjectID:        input.ClaGroupID,
		RepositoryType:             aws.String("github"),
		RepositoryURL:              ghRepo.HTMLURL,
	}
	return s.repo.AddGithubRepository(externalProjectID, projectSFID, in)
}

func (s *service) EnableRepository(repositoryID string) error {
	return s.repo.EnableRepository(repositoryID)
}

func (s *service) DisableRepository(repositoryID string) error {
	return s.repo.DisableRepository(repositoryID)
}

func (s *service) ListProjectRepositories(projectSFID string) (*v1Models.ListGithubRepositories, error) {
	return s.repo.ListProjectRepositories("", projectSFID, true)
}

func (s *service) GetRepository(repositoryID string) (*v1Models.GithubRepository, error) {
	return s.repo.GetRepository(repositoryID)
}

func (s *service) GetProtectedBranch(repositoryID string) (*v2Models.GithubRepositoryBranchProtection, error) {
	githubRepository, err := s.GetRepository(repositoryID)
	if err != nil {
		log.Warnf("fetching repository failed : %s : %v", repositoryID, err)
		return nil, err
	}

	githubOrgName := githubRepository.RepositoryOrganizationName
	githubRepoName := githubRepository.RepositoryName
	githubRepoName = s.cleanGithubRepoName(githubRepoName)

	githubClient, err := s.getGithubClientForOrgName(githubOrgName)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	owner, branchName, err := s.getGithubOwnerBranchName(githubClient, githubOrgName, githubRepoName)
	if err != nil {
		return nil, err
	}

	result := &v2Models.GithubRepositoryBranchProtection{
		BranchName: &branchName,
	}
	branchProtection, err := github.GetProtectedBranch(ctx, githubClient, owner, githubRepoName, branchName)
	if err != nil {
		if errors.Is(err, github.ErrBranchNotProtected) {
			return result, nil
		}
		log.Warnf("getting the github protected branch for owner : %s, repo : %s and branch : %s failed : %v", owner, githubRepoName, branchName, err)
		return nil, err
	}

	result.ProtectionEnabled = true
	if github.IsEnforceAdminEnabled(branchProtection) {
		result.EnforceAdmin = true
	}

	requiredChecks := requiredBranchProtectionChecks
	requiredChecksResult := s.getRequiredProtectedBranchCheckStatus(branchProtection, requiredChecks)
	result.StatusChecks = requiredChecksResult

	return result, nil
}

func (s *service) UpdateProtectedBranch(repositoryID string, input *v2Models.GithubRepositoryBranchProtectionInput) (*v2Models.GithubRepositoryBranchProtection, error) {
	githubRepository, err := s.GetRepository(repositoryID)
	if err != nil {
		log.Warnf("fetching repository failed : %s : %v", repositoryID, err)
		return nil, err
	}

	githubOrgName := githubRepository.RepositoryOrganizationName
	githubRepoName := githubRepository.RepositoryName
	githubRepoName = s.cleanGithubRepoName(githubRepoName)

	githubClient, err := s.getGithubClientForOrgName(githubOrgName)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	owner, branchName, err := s.getGithubOwnerBranchName(githubClient, githubOrgName, githubRepoName)
	if err != nil {
		return nil, err
	}

	var requiredChecks []string
	var disabledChecks []string
	if input.StatusChecks != nil {
		for _, inputCheck := range input.StatusChecks {
			// we want to make sure we only mutate checks related to lf
			var found bool
			for _, rc := range requiredBranchProtectionChecks {
				if rc == *inputCheck.Name {
					found = true
					break
				}
			}

			// just ignore that check if it's something not in our options
			if !found {
				continue
			}

			if !*inputCheck.Enabled {
				disabledChecks = append(disabledChecks, *inputCheck.Name)
				continue
			}
			requiredChecks = append(requiredChecks, *inputCheck.Name)
		}
	}

	err = github.EnableBranchProtection(ctx, githubClient, owner, githubRepoName, branchName, *input.EnforceAdmin, requiredChecks, disabledChecks)
	if err != nil {
		return nil, err
	}

	return s.GetProtectedBranch(repositoryID)
}

func (s *service) cleanGithubRepoName(githubRepoName string) string {
	if strings.Contains(githubRepoName, "/") {
		parts := strings.Split(githubRepoName, "/")
		githubRepoName = parts[len(parts)-1]
	}
	return githubRepoName
}

func (s *service) getGithubClientForOrgName(githubOrgName string) (*githubsdk.Client, error) {
	githubOrg, err := s.ghOrgRepo.GetGithubOrganization(githubOrgName)
	if err != nil {
		log.Warnf("fetching githubOrg failed : %s : %v", githubOrgName, err)
		return nil, err
	}

	githubClient, err := github.NewGithubAppClient(githubOrg.OrganizationInstallationID)
	if err != nil {
		log.Warnf("creating the github client for installation id failed  : %d : %v", githubOrg.OrganizationInstallationID, err)
		return nil, err
	}

	return githubClient, nil
}

func (s *service) getGithubOwnerBranchName(githubClient *githubsdk.Client, githubOrgName, githubRepoName string) (string, string, error) {
	ctx := context.Background()
	owner, err := github.GetOwnerName(ctx, githubClient, githubOrgName, githubRepoName)
	if err != nil {
		log.Warnf("getting the owner name for org : %s and repo : %s failed : %v", githubOrgName, githubRepoName, err)
		return "", "", err
	}

	if owner == "" {
		log.Warnf("github returned empty owner name for org : %s and repo : %s", githubOrgName, githubRepoName)
		return "", "", fmt.Errorf("empty owner name")
	}

	log.Debugf("getGithubOwnerBranchName : owner of the repo : %s found : %s", owner, githubRepoName)
	branchName, err := github.GetDefaultBranchForRepo(ctx, githubClient, owner, githubRepoName)
	if err != nil {
		log.Warnf("getting default github branch failed for owner : %s and repo : %s : %v", owner, githubRepoName, err)
		return "", "", err
	}

	return owner, branchName, nil
}

// getRequiredProtectedBranchCheckStatus
func (s *service) getRequiredProtectedBranchCheckStatus(protectedBranch *githubsdk.Protection, requiredChecks []string) []*v2Models.GithubRepositoryBranchProtectionStatusChecks {
	var result []*v2Models.GithubRepositoryBranchProtectionStatusChecks
	resultMap := map[string]bool{}
	for _, rc := range requiredChecks {
		result = append(result, &v2Models.GithubRepositoryBranchProtectionStatusChecks{
			Name:    swag.String(rc),
			Enabled: swag.Bool(false),
		})
		resultMap[rc] = true
	}
	if protectedBranch.RequiredStatusChecks == nil || len(protectedBranch.RequiredStatusChecks.Contexts) == 0 {
		return result
	}

	for _, existingCheck := range protectedBranch.RequiredStatusChecks.Contexts {
		if !resultMap[existingCheck] {
			continue
		}

		// we mark it as enabled in this case
		for i := range result {
			if *result[i].Name == existingCheck {
				result[i].Enabled = swag.Bool(true)
			}
		}
	}

	return result
}

func (s *service) DisableCLAGroupRepositories(claGroupID string) error {
	var deleteErr error
	ghOrgs, err := s.repo.GetCLAGroupRepositoriesGroupByOrgs(claGroupID, true)
	if err != nil {
		return err
	}
	if len(ghOrgs) > 0 {
		log.Debugf("Deleting repositories for cla-group :%s", claGroupID)
		for _, ghOrg := range ghOrgs {
			for _, item := range ghOrg.List {
				deleteErr = s.repo.DisableRepository(item.RepositoryID)
				if deleteErr != nil {
					log.Warnf("Unable to remove repository: %s for project :%s error :%v", item.RepositoryID, claGroupID, deleteErr)
				}
			}
		}
	}
	return nil
}
