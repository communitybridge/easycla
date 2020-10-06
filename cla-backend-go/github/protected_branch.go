// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	githubpkg "github.com/google/go-github/github"
)

const (
	defaultBranchName = "master"
)

var (
	// ErrBranchNotProtected indicates the situation where the branch is not enabled for protection on github side
	ErrBranchNotProtected = errors.New("not protected")
)

// Repositories is part of the interface working with github repositories, it's inside of the github client
// It's extracted here as interface so we can mock that functionality.
type Repositories interface {
	ListByOrg(ctx context.Context, org string, opt *githubpkg.RepositoryListByOrgOptions) ([]*githubpkg.Repository, *githubpkg.Response, error)
	Get(ctx context.Context, owner, repo string) (*githubpkg.Repository, *githubpkg.Response, error)
	GetBranchProtection(ctx context.Context, owner, repo, branch string) (*githubpkg.Protection, *githubpkg.Response, error)
	UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *githubpkg.ProtectionRequest) (*githubpkg.Protection, *githubpkg.Response, error)
}

// GetOwnerName retrieves the owner name of the given org and repo name
func GetOwnerName(ctx context.Context, githubRepo Repositories, orgName, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	log.Debugf("GetOwnerName : getting owner name for org %s and repoName : %s", orgName, repoName)
	listOpt := &githubpkg.RepositoryListByOrgOptions{
		ListOptions: githubpkg.ListOptions{
			PerPage: 30,
		},
	}
	for {
		repos, resp, err := githubRepo.ListByOrg(ctx, orgName, listOpt)
		if err != nil {
			if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
				return "", wErr
			}
			return "", err
		}

		if len(repos) == 0 {
			log.Warnf("GetOwnerName : no repos found under orgName : %s (maybe no access ?)", orgName)
			return "", nil
		}

		for _, repo := range repos {
			if *repo.Name == repoName {
				if repo.Owner != nil {
					owner := *repo.Owner
					return *owner.Login, nil
				}
			}
		}

		// means we're at the end of it so exit
		if resp.NextPage == 0 {
			log.Warnf("GetOwnerName : owner name not found for org : %s and repo : %s", orgName, repoName)
			return "", nil
		}

		// set it to the next page
		listOpt.Page = resp.NextPage
	}
}

// GetDefaultBranchForRepo helps with pulling the default branch for the given repo
func GetDefaultBranchForRepo(ctx context.Context, githubRepo Repositories, owner, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	repo, resp, err := githubRepo.Get(ctx, owner, repoName)
	if err != nil {
		if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
			return "", wErr
		}
		return "", err
	}

	var defaultBranch string
	if repo.DefaultBranch == nil {
		defaultBranch = defaultBranchName
	} else {
		defaultBranch = *repo.DefaultBranch
	}

	return defaultBranch, nil
}

// GetProtectedBranch fetches the protected branch details
func GetProtectedBranch(ctx context.Context, githubRepo Repositories, owner, repoName, protectedBranchName string) (*githubpkg.Protection, error) {
	repoName = CleanGithubRepoName(repoName)
	protection, resp, err := githubRepo.GetBranchProtection(ctx, owner, repoName, protectedBranchName)

	if err != nil {
		if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
			return nil, wErr
		}
		if resp != nil && resp.StatusCode == 404 {
			if gErr, ok := err.(*githubpkg.ErrorResponse); ok {
				if strings.Contains(strings.ToLower(gErr.Message), "not protected") {
					return nil, ErrBranchNotProtected
				}
			}
		}

		return nil, fmt.Errorf("fetching branch proteciton : %w", err)
	}
	return protection, err
}

//EnableBranchProtection enables branch protection if not enabled and makes sure passed arguments such as enforceAdmin
//statusChecks are applied. The operation makes sure it doesn't override the existing checks.
func EnableBranchProtection(ctx context.Context, githubRepo Repositories, owner, repoName, branchName string, enforceAdmin bool, enableStatusChecks, disableStatusChecks []string) error {
	repoName = CleanGithubRepoName(repoName)
	protectedBranch, err := GetProtectedBranch(ctx, githubRepo, owner, repoName, branchName)
	if err != nil && !errors.Is(err, ErrBranchNotProtected) {
		return fmt.Errorf("fetching the protected branch for repo : %s : %w", repoName, err)
	}

	branchProtectionRequest, err := createBranchProtectionRequest(protectedBranch, enableStatusChecks, disableStatusChecks, enforceAdmin)
	if err != nil {
		return fmt.Errorf("creating branch protection request failed : %v", err)
	}

	_, resp, err := githubRepo.UpdateBranchProtection(ctx, owner, repoName, branchName, branchProtectionRequest)

	if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
		return wErr
	}
	return err
}

// createBranchProtectionRequest creates a branch protection request from existing protection
func createBranchProtectionRequest(protection *githubpkg.Protection, enableStatusChecks, disableStatusChecks []string, enforceAdmin bool) (*githubpkg.ProtectionRequest, error) {
	var currentChecks *githubpkg.RequiredStatusChecks
	if protection != nil {
		currentChecks = protection.RequiredStatusChecks
	}
	requiredStatusChecks := mergeStatusChecks(currentChecks, enableStatusChecks, disableStatusChecks)

	branchProtection := &githubpkg.ProtectionRequest{
		RequiredStatusChecks: requiredStatusChecks,
		EnforceAdmins:        enforceAdmin,
	}

	// don't have to check further in this case
	if protection == nil {
		return branchProtection, nil
	}

	if protection.RequiredPullRequestReviews != nil {
		var pullRequestReviewEnforcement githubpkg.PullRequestReviewsEnforcementRequest
		if err := copier.Copy(&pullRequestReviewEnforcement, protection.RequiredPullRequestReviews); err != nil {
			return nil, fmt.Errorf("copying from protected branch to request failed : requiredPullRequestReviews : %v", err)
		}

		if len(protection.RequiredPullRequestReviews.DismissalRestrictions.Users) > 0 {
			var users []string
			for _, user := range protection.RequiredPullRequestReviews.DismissalRestrictions.Users {
				users = append(users, *user.Login)
			}
			if pullRequestReviewEnforcement.DismissalRestrictionsRequest == nil {
				pullRequestReviewEnforcement.DismissalRestrictionsRequest = &githubpkg.DismissalRestrictionsRequest{}
			}
			pullRequestReviewEnforcement.DismissalRestrictionsRequest.Users = &users
		}

		if len(protection.RequiredPullRequestReviews.DismissalRestrictions.Teams) > 0 {
			var teams []string
			for _, team := range protection.RequiredPullRequestReviews.DismissalRestrictions.Teams {
				teams = append(teams, *team.Slug)
			}
			if pullRequestReviewEnforcement.DismissalRestrictionsRequest == nil {
				pullRequestReviewEnforcement.DismissalRestrictionsRequest = &githubpkg.DismissalRestrictionsRequest{}
			}
			pullRequestReviewEnforcement.DismissalRestrictionsRequest.Teams = &teams
		}

		branchProtection.RequiredPullRequestReviews = &pullRequestReviewEnforcement
	}

	if protection.Restrictions != nil {
		var restrictions githubpkg.BranchRestrictionsRequest
		if len(protection.Restrictions.Users) > 0 {
			var users []string
			for _, user := range protection.Restrictions.Users {
				users = append(users, *user.Login)
			}
			restrictions.Users = users
		}

		if len(protection.Restrictions.Teams) > 0 {
			var teams []string
			for _, team := range protection.Restrictions.Teams {
				teams = append(teams, *team.Slug)
			}
			restrictions.Teams = teams
		}

		branchProtection.Restrictions = &restrictions
	}

	return branchProtection, nil
}

//mergeStatusChecks merges the current checks with the new ones and disable the ones that are specified
func mergeStatusChecks(currentCheck *githubpkg.RequiredStatusChecks, enableContexts, disableContexts []string) *githubpkg.RequiredStatusChecks {

	// seems github api is not happy with nils for arrays ;)
	if len(enableContexts) == 0 {
		enableContexts = []string{}
	}

	if currentCheck == nil || len(currentCheck.Contexts) == 0 {
		return &githubpkg.RequiredStatusChecks{
			Strict:   true,
			Contexts: enableContexts,
		}
	}

	finalContexts := []string{}
	uniqueEnableContexts := map[string]bool{}

	for _, c := range currentCheck.Contexts {
		// first disable the ones we're not interested into
		found := false
		if len(disableContexts) > 0 {
			for _, disableContext := range disableContexts {
				if disableContext == c {
					found = true
					break
				}
			}
		}

		if found {
			continue
		}

		uniqueEnableContexts[c] = true
		finalContexts = append(finalContexts, c)
	}

	for _, c := range enableContexts {
		if uniqueEnableContexts[c] {
			continue
		}

		uniqueEnableContexts[c] = true
		finalContexts = append(finalContexts, c)
	}

	currentCheck.Contexts = finalContexts
	currentCheck.Strict = true

	return currentCheck
}

//IsEnforceAdminEnabled checks if enforce admin option is enabled for the branch protection
func IsEnforceAdminEnabled(protection *githubpkg.Protection) bool {
	if protection.EnforceAdmins == nil {
		return false
	}

	return protection.EnforceAdmins.Enabled
}

// CleanGithubRepoName removes the orgname if existing in the string
func CleanGithubRepoName(githubRepoName string) string {
	if strings.Contains(githubRepoName, "/") {
		parts := strings.Split(githubRepoName, "/")
		githubRepoName = parts[len(parts)-1]
	}
	return githubRepoName
}
