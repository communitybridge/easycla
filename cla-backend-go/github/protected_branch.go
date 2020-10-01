// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	githubpkg "github.com/google/go-github/github"
)

var (
	// ErrBranchNotProtected indicates the situation where the branch is not enabled for protection on github side
	ErrBranchNotProtected = errors.New("not protected")
)

// GetOwnerName retrieves the owner name of the given org and repo name
func GetOwnerName(ctx context.Context, client *githubpkg.Client, orgName, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	log.Debugf("GetOwnerName : getting owner name for org %s and repoName : %s", orgName, repoName)
	listOpt := &githubpkg.RepositoryListByOrgOptions{
		ListOptions: githubpkg.ListOptions{
			PerPage: 30,
		},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, listOpt)
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
func GetDefaultBranchForRepo(ctx context.Context, client *githubpkg.Client, owner, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	repo, resp, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
			return "", wErr
		}
		return "", err
	}

	var defaultBranch string
	if repo.DefaultBranch == nil {
		defaultBranch = "master"
	} else {
		defaultBranch = *repo.DefaultBranch
	}

	return defaultBranch, nil
}

// GetProtectedBranch fetches the protected branch details
func GetProtectedBranch(ctx context.Context, client *githubpkg.Client, owner, repoName, protectedBranchName string) (*githubpkg.Protection, error) {
	repoName = CleanGithubRepoName(repoName)
	protection, resp, err := client.Repositories.GetBranchProtection(ctx, owner, repoName, protectedBranchName)

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
func EnableBranchProtection(ctx context.Context, client *githubpkg.Client, owner, repoName, branchName string, enforceAdmin bool, enableStatusChecks, disableStatusChecks []string) error {
	repoName = CleanGithubRepoName(repoName)
	protectedBranch, err := GetProtectedBranch(ctx, client, owner, repoName, branchName)
	if err != nil && !errors.Is(err, ErrBranchNotProtected) {
		return fmt.Errorf("fetching the protected branch for repo : %s : %w", repoName, err)
	}

	var currentChecks *githubpkg.RequiredStatusChecks
	if protectedBranch != nil {
		currentChecks = protectedBranch.RequiredStatusChecks
	}
	requiredStatusChecks := mergeStatusChecks(currentChecks, enableStatusChecks, disableStatusChecks)

	branchProtection := &githubpkg.ProtectionRequest{
		EnforceAdmins:        enforceAdmin,
		RequiredStatusChecks: requiredStatusChecks,
	}
	_, resp, err := client.Repositories.UpdateBranchProtection(ctx, owner, repoName, branchName, branchProtection)

	if ok, wErr := checkAndWrapForKnownErrors(resp, err); ok {
		return wErr
	}
	return err
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

//AreStatusChecksEnabled checks if all of the status checks are enabled
func AreStatusChecksEnabled(protection *githubpkg.Protection, checks []string) bool {
	if len(checks) == 0 {
		return false
	}

	currentChecks := protection.RequiredStatusChecks
	if currentChecks == nil || !protection.RequiredStatusChecks.Strict {
		return false
	}

	if len(currentChecks.Contexts) < len(checks) {
		return false
	}

	var found []string
	for _, cc := range currentChecks.Contexts {
		for _, c := range checks {
			if c == cc {
				found = append(found, cc)
			}
		}
	}

	return len(found) == len(checks)
}

// CleanGithubRepoName removes the orgname if existing in the string
func CleanGithubRepoName(githubRepoName string) string {
	if strings.Contains(githubRepoName, "/") {
		parts := strings.Split(githubRepoName, "/")
		githubRepoName = parts[len(parts)-1]
	}
	return githubRepoName
}
