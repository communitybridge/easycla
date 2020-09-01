package github

import (
	"context"
	"errors"
	"fmt"
	githubpkg "github.com/google/go-github/github"
	"strings"
)

var (
	BranchNotProtectedError = errors.New("not protected")
	branchNotFoundError     = errors.New("not found")
)

func GetOrganization(ctx context.Context, client *githubpkg.Client, orgName string) (*githubpkg.Organization, error) {
	org, _, err := client.Organizations.Get(ctx, orgName)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func GetOwnerName(ctx context.Context, client *githubpkg.Client, orgName, repoName string) (string, error) {
	repos, _, err := client.Repositories.ListByOrg(ctx, orgName, nil)
	if err != nil {
		return "", err
	}

	for _, repo := range repos {
		if *repo.Name == repoName {
			if repo.Owner != nil {
				owner := *repo.Owner
				return *owner.Login, nil
			}
		}
	}

	return "", nil
}

// GetDefaultBranchForRepo helps with pulling the default branch for the given repo
func GetDefaultBranchForRepo(ctx context.Context, client *githubpkg.Client, owner, repoName string) (string, error) {
	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
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
	protection, resp, err := client.Repositories.GetBranchProtection(ctx, owner, repoName, protectedBranchName)

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			if gErr, ok := err.(*githubpkg.ErrorResponse); ok {
				if strings.Contains(strings.ToLower(gErr.Message), "not protected") {
					return nil, BranchNotProtectedError
				}
			}
		}

		return nil, fmt.Errorf("fetching branch proteciton : %w", err)
	}
	return protection, err
}

func EnableBranchProtection(ctx context.Context, client *githubpkg.Client, owner, repoName, branchName string, enforceAdmin bool, statusChecks []string) error {
	protectedBranch, err := GetProtectedBranch(ctx, client, owner, repoName, branchName)

	if err != nil && !errors.Is(err, BranchNotProtectedError) {
		return fmt.Errorf("fetching the protected branch : %w", err)
	}

	var currentChecks *githubpkg.RequiredStatusChecks
	if protectedBranch != nil {
		currentChecks = protectedBranch.RequiredStatusChecks
	}
	requiredStatusChecks := MergeStatusChecks(currentChecks, statusChecks)

	_, _, err = client.Repositories.UpdateBranchProtection(ctx, owner, repoName, branchName, &githubpkg.ProtectionRequest{
		EnforceAdmins:        enforceAdmin,
		RequiredStatusChecks: requiredStatusChecks,
	})
	return err
}

func MergeStatusChecks(currentCheck *githubpkg.RequiredStatusChecks, contexts []string) *githubpkg.RequiredStatusChecks {

	if currentCheck == nil || len(currentCheck.Contexts) == 0 {
		return &githubpkg.RequiredStatusChecks{
			Strict:   true,
			Contexts: contexts,
		}
	}

	var finalContexts []string
	uniqueContexts := map[string]bool{}

	for _, c := range currentCheck.Contexts {
		uniqueContexts[c] = true
		finalContexts = append(finalContexts, c)
	}

	for _, c := range contexts {
		if uniqueContexts[c] {
			continue
		}

		uniqueContexts[c] = true
		finalContexts = append(finalContexts, c)
	}

	currentCheck.Contexts = finalContexts
	return currentCheck
}

func IsEnforceAdminEnabled(protection *githubpkg.Protection) bool {
	if protection.EnforceAdmins == nil {
		return false
	}

	return protection.EnforceAdmins.Enabled
}

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
