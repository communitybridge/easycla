// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/jinzhu/copier"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	githubpkg "github.com/google/go-github/v32/github"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
)

const (
	defaultBranchName = "master"
)

var (
	// ErrBranchNotProtected indicates the situation where the branch is not enabled for protection on github side
	ErrBranchNotProtected = errors.New("not protected")
)

// rate limiting variables
var (
	// blockingRateLimit is useful for background tasks where the interaction is more predictable
	blockingRateLimit = ratelimit.New(2)
	// nonBlockingRateLimit is preferred when the github methods would be called realtime
	// in this case we can call Allow method to check if can proceed or return error
	nonBlockingRateLimit = rate.NewLimiter(2, 5)
)

// Repositories is part of the interface working with github repositories, it's inside of the github client
// It's extracted here as interface so we can mock that functionality.
type Repositories interface {
	ListByOrg(ctx context.Context, org string, opt *githubpkg.RepositoryListByOrgOptions) ([]*githubpkg.Repository, *githubpkg.Response, error)
	Get(ctx context.Context, owner, repo string) (*githubpkg.Repository, *githubpkg.Response, error)
	GetBranchProtection(ctx context.Context, owner, repo, branch string) (*githubpkg.Protection, *githubpkg.Response, error)
	UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *githubpkg.ProtectionRequest) (*githubpkg.Protection, *githubpkg.Response, error)
}

type blockingRateLimitRepositories struct {
	Repositories
}

// NewBlockLimiterRepositories returns a new instance of Repositories interface with blocking rate limiting
// where when the limit is reached the next call blocks till the bucket is ready again
func NewBlockLimiterRepositories(repos Repositories) Repositories {
	return blockingRateLimitRepositories{repos}
}

func (b blockingRateLimitRepositories) ListByOrg(ctx context.Context, org string, opt *githubpkg.RepositoryListByOrgOptions) ([]*githubpkg.Repository, *githubpkg.Response, error) {
	blockingRateLimit.Take()
	return b.Repositories.ListByOrg(ctx, org, opt)
}

func (b blockingRateLimitRepositories) Get(ctx context.Context, owner, repo string) (*githubpkg.Repository, *githubpkg.Response, error) {
	blockingRateLimit.Take()
	return b.Repositories.Get(ctx, owner, repo)
}

func (b blockingRateLimitRepositories) GetBranchProtection(ctx context.Context, owner, repo, branch string) (*githubpkg.Protection, *githubpkg.Response, error) {
	blockingRateLimit.Take()
	return b.Repositories.GetBranchProtection(ctx, owner, repo, branch)
}

func (b blockingRateLimitRepositories) UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *githubpkg.ProtectionRequest) (*githubpkg.Protection, *githubpkg.Response, error) {
	blockingRateLimit.Take()
	return b.Repositories.UpdateBranchProtection(ctx, owner, repo, branch, preq)
}

type nonBlockingRateLimitRepositories struct {
	Repositories
}

// NewNonBlockLimiterRepositories returns a new instance of Repositories interface with non blocking rate limiting
func NewNonBlockLimiterRepositories(repos Repositories) Repositories {
	return nonBlockingRateLimitRepositories{repos}
}

func (nb nonBlockingRateLimitRepositories) ListByOrg(ctx context.Context, org string, opt *githubpkg.RepositoryListByOrgOptions) ([]*githubpkg.Repository, *githubpkg.Response, error) {
	if nonBlockingRateLimit.Allow() {
		return nb.Repositories.ListByOrg(ctx, org, opt)
	}
	return nil, nil, fmt.Errorf("too many requests : %w", ErrRateLimited)
}

func (nb nonBlockingRateLimitRepositories) Get(ctx context.Context, owner, repo string) (*githubpkg.Repository, *githubpkg.Response, error) {
	if nonBlockingRateLimit.Allow() {
		return nb.Repositories.Get(ctx, owner, repo)
	}
	return nil, nil, fmt.Errorf("too many requests : %w", ErrRateLimited)
}

func (nb nonBlockingRateLimitRepositories) GetBranchProtection(ctx context.Context, owner, repo, branch string) (*githubpkg.Protection, *githubpkg.Response, error) {
	if nonBlockingRateLimit.Allow() {
		return nb.Repositories.GetBranchProtection(ctx, owner, repo, branch)
	}
	return nil, nil, fmt.Errorf("too many requests : %w", ErrRateLimited)
}

func (nb nonBlockingRateLimitRepositories) UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *githubpkg.ProtectionRequest) (*githubpkg.Protection, *githubpkg.Response, error) {
	if nonBlockingRateLimit.Allow() {
		return nb.Repositories.UpdateBranchProtection(ctx, owner, repo, branch, preq)
	}
	return nil, nil, fmt.Errorf("too many requests : %w", ErrRateLimited)
}

type branchProtectionRepositoryConfig struct {
	enableBlockingLimiter    bool
	enableNonBlockingLimiter bool
}

// BranchProtectionRepositoryOption enables optional parameters to BranchProtectionRepository
type BranchProtectionRepositoryOption func(config *branchProtectionRepositoryConfig)

// EnableBlockingLimiter enables the blocking limiter
func EnableBlockingLimiter() BranchProtectionRepositoryOption {
	return func(config *branchProtectionRepositoryConfig) {
		config.enableBlockingLimiter = true
	}
}

// EnableNonBlockingLimiter enables the non-blocking limiter
func EnableNonBlockingLimiter() BranchProtectionRepositoryOption {
	return func(config *branchProtectionRepositoryConfig) {
		config.enableNonBlockingLimiter = true
	}
}

// BranchProtectionRepository contains helper methods interacting with github api related to branch protection
type BranchProtectionRepository struct {
	githubRepo Repositories
}

// NewBranchProtectionRepository creates a new BranchProtectionRepository
func NewBranchProtectionRepository(githubRepo Repositories, opts ...BranchProtectionRepositoryOption) *BranchProtectionRepository {
	config := &branchProtectionRepositoryConfig{}
	for _, o := range opts {
		o(config)
	}

	if config.enableNonBlockingLimiter {
		githubRepo = NewNonBlockLimiterRepositories(githubRepo)
	} else if config.enableBlockingLimiter {
		githubRepo = NewBlockLimiterRepositories(githubRepo)
	}

	return &BranchProtectionRepository{
		githubRepo: githubRepo,
	}
}

// GetOwnerName retrieves the owner name of the given org and repo name
func (bp *BranchProtectionRepository) GetOwnerName(ctx context.Context, orgName, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	log.Debugf("GetOwnerName : getting owner name for org %s and repoName : %s", orgName, repoName)
	listOpt := &githubpkg.RepositoryListByOrgOptions{
		ListOptions: githubpkg.ListOptions{
			PerPage: 30,
		},
	}
	for {
		repos, resp, err := bp.githubRepo.ListByOrg(ctx, orgName, listOpt)
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
func (bp *BranchProtectionRepository) GetDefaultBranchForRepo(ctx context.Context, owner, repoName string) (string, error) {
	repoName = CleanGithubRepoName(repoName)
	repo, resp, err := bp.githubRepo.Get(ctx, owner, repoName)
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
func (bp *BranchProtectionRepository) GetProtectedBranch(ctx context.Context, owner, repoName, protectedBranchName string) (*githubpkg.Protection, error) {
	repoName = CleanGithubRepoName(repoName)
	protection, resp, err := bp.githubRepo.GetBranchProtection(ctx, owner, repoName, protectedBranchName)

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
func (bp *BranchProtectionRepository) EnableBranchProtection(ctx context.Context, owner, repoName, branchName string, enforceAdmin bool, enableStatusChecks, disableStatusChecks []string) error {
	repoName = CleanGithubRepoName(repoName)
	protectedBranch, err := bp.GetProtectedBranch(ctx, owner, repoName, branchName)
	if err != nil && !errors.Is(err, ErrBranchNotProtected) {
		return fmt.Errorf("fetching the protected branch for repo : %s : %w", repoName, err)
	}

	branchProtectionRequest, err := createBranchProtectionRequest(protectedBranch, enableStatusChecks, disableStatusChecks, enforceAdmin)
	if err != nil {
		return fmt.Errorf("creating branch protection request failed : %v", err)
	}

	_, resp, err := bp.githubRepo.UpdateBranchProtection(ctx, owner, repoName, branchName, branchProtectionRequest)

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

	branchProtectionRequest := &githubpkg.ProtectionRequest{
		RequiredStatusChecks: requiredStatusChecks,
		EnforceAdmins:        enforceAdmin,
	}

	// don't have to check further in this case
	if protection == nil {
		return branchProtectionRequest, nil
	}

	if protection.RequireLinearHistory != nil {
		branchProtectionRequest.RequireLinearHistory = swag.Bool(protection.RequireLinearHistory.Enabled)
	}

	if protection.AllowForcePushes != nil {
		branchProtectionRequest.AllowForcePushes = swag.Bool(protection.AllowForcePushes.Enabled)
	}

	if protection.AllowDeletions != nil {
		branchProtectionRequest.AllowDeletions = swag.Bool(protection.AllowDeletions.Enabled)
	}

	if protection.RequiredPullRequestReviews != nil {
		var pullRequestReviewEnforcement githubpkg.PullRequestReviewsEnforcementRequest
		if err := copier.Copy(&pullRequestReviewEnforcement, protection.RequiredPullRequestReviews); err != nil {
			return nil, fmt.Errorf("copying from protected branch to request failed : requiredPullRequestReviews : %v", err)
		}

		// github is not happy about null arrays, prefers empty arrays ...
		//No subschema in "anyOf" matched.
		//For 'properties/teams', nil is not an array.
		//Not all subschemas of "allOf" matched.
		var anyEnabled bool
		if len(protection.RequiredPullRequestReviews.DismissalRestrictions.Users) > 0 {
			anyEnabled = true
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
			anyEnabled = true
			var teams []string
			for _, team := range protection.RequiredPullRequestReviews.DismissalRestrictions.Teams {
				teams = append(teams, *team.Slug)
			}
			if pullRequestReviewEnforcement.DismissalRestrictionsRequest == nil {
				pullRequestReviewEnforcement.DismissalRestrictionsRequest = &githubpkg.DismissalRestrictionsRequest{}
			}
			pullRequestReviewEnforcement.DismissalRestrictionsRequest.Teams = &teams
		}

		if anyEnabled {
			if pullRequestReviewEnforcement.DismissalRestrictionsRequest.Users == nil {
				pullRequestReviewEnforcement.DismissalRestrictionsRequest.Users = &[]string{}
			}

			if pullRequestReviewEnforcement.DismissalRestrictionsRequest.Teams == nil {
				pullRequestReviewEnforcement.DismissalRestrictionsRequest.Teams = &[]string{}
			}

		}

		branchProtectionRequest.RequiredPullRequestReviews = &pullRequestReviewEnforcement
	}

	if protection.Restrictions != nil {
		var restrictions githubpkg.BranchRestrictionsRequest
		var anyEnabled bool
		if len(protection.Restrictions.Users) > 0 {
			anyEnabled = true
			var users []string
			for _, user := range protection.Restrictions.Users {
				users = append(users, *user.Login)
			}
			restrictions.Users = users
		}

		if len(protection.Restrictions.Teams) > 0 {
			anyEnabled = true
			var teams []string
			for _, team := range protection.Restrictions.Teams {
				teams = append(teams, *team.Slug)
			}
			restrictions.Teams = teams
		}

		if len(protection.Restrictions.Apps) > 0 {
			anyEnabled = true
			var apps []string
			for _, app := range protection.Restrictions.Apps {
				apps = append(apps, *app.Slug)
			}
			restrictions.Apps = apps
		}

		// make sure we don't send nil arrays ...
		if anyEnabled {
			if restrictions.Users == nil {
				restrictions.Users = []string{}
			}

			if restrictions.Teams == nil {
				restrictions.Teams = []string{}
			}

			if restrictions.Apps == nil {
				restrictions.Apps = []string{}
			}
		}

		branchProtectionRequest.Restrictions = &restrictions
	}

	return branchProtectionRequest, nil
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
