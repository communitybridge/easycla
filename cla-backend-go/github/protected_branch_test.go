// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-openapi/swag"

	"github.com/golang/mock/gomock"

	"github.com/bmizerany/assert"
	githubsdk "github.com/google/go-github/v32/github"
)

// TestMergeStatusChecks tests the functionality of where we enable/disable checks
func TestMergeStatusChecks(t *testing.T) {

	testCases := []struct {
		Name            string
		currentChecks   *githubsdk.RequiredStatusChecks
		expectedChecks  *githubsdk.RequiredStatusChecks
		enableContexts  []string
		disableContexts []string
	}{
		{
			Name: "all empty",
			expectedChecks: &githubsdk.RequiredStatusChecks{
				Strict:   true,
				Contexts: []string{},
			},
		},
		{
			Name: "empty state enable",
			expectedChecks: &githubsdk.RequiredStatusChecks{
				Strict:   true,
				Contexts: []string{"EasyCLA"},
			},
			enableContexts: []string{"EasyCLA"},
		},
		{
			Name: "preserve existing enable more",
			currentChecks: &githubsdk.RequiredStatusChecks{
				Contexts: []string{"travis-ci"},
			},
			expectedChecks: &githubsdk.RequiredStatusChecks{
				Strict:   true,
				Contexts: []string{"travis-ci", "EasyCLA"},
			},
			enableContexts: []string{"EasyCLA"},
		},
		{
			Name: "preserve existing disable some",
			currentChecks: &githubsdk.RequiredStatusChecks{
				Contexts: []string{"travis-ci", "EasyCLA"},
			},
			expectedChecks: &githubsdk.RequiredStatusChecks{
				Strict:   true,
				Contexts: []string{"travis-ci"},
			},
			disableContexts: []string{"EasyCLA"},
		},
		{
			Name: "add and remove in same operation",
			currentChecks: &githubsdk.RequiredStatusChecks{
				Contexts: []string{"travis-ci", "DCO", "EasyCLA"},
			},
			expectedChecks: &githubsdk.RequiredStatusChecks{
				Strict:   true,
				Contexts: []string{"travis-ci", "EasyCLA", "CodeQL"},
			},
			enableContexts:  []string{"CodeQL"},
			disableContexts: []string{"DCO"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			result := mergeStatusChecks(tc.currentChecks, tc.enableContexts, tc.disableContexts)
			assert.Equal(tt, tc.expectedChecks, result)
		})
	}
}

func TestEnableBranchProtection(t *testing.T) {
	owner := "johnenable"
	repo := "johnsrepoenable"
	branchName := defaultBranchName

	testCases := []struct {
		Name              string
		Checks            []string
		Protection        *githubsdk.Protection
		ProtectionRequest *githubsdk.ProtectionRequest
		Err               error
	}{
		{
			Name:       "success",
			Checks:     []string{"easyCLA"},
			Protection: &githubsdk.Protection{},
			ProtectionRequest: &githubsdk.ProtectionRequest{
				EnforceAdmins: true,
				RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
					Strict:   true,
					Contexts: []string{"easyCLA"},
				},
			},
		},
		{
			Name:   "preserve existing checks",
			Checks: []string{"easyCLA"},
			Protection: &githubsdk.Protection{
				RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
					Strict:   false,
					Contexts: []string{"circle/ci"},
				},
			},
			ProtectionRequest: &githubsdk.ProtectionRequest{
				EnforceAdmins: true,
				RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
					Strict:   true,
					Contexts: []string{"circle/ci", "easyCLA"},
				},
			},
		},
		{
			Name:   "preserve existing settings",
			Checks: []string{"easyCLA"},
			Protection: &githubsdk.Protection{
				RequiredPullRequestReviews: &githubsdk.PullRequestReviewsEnforcement{
					RequireCodeOwnerReviews:      true,
					RequiredApprovingReviewCount: 2,
					DismissalRestrictions: &githubsdk.DismissalRestrictions{
						Users: []*githubsdk.User{
							{Login: swag.String("alex")},
						},
						Teams: []*githubsdk.Team{
							{Slug: swag.String("alpha")},
						},
					},
				},
				Restrictions: &githubsdk.BranchRestrictions{
					Users: []*githubsdk.User{
						{Login: swag.String("john")},
					},
					Teams: []*githubsdk.Team{
						{Slug: swag.String("easyCLA-Team")},
					},
				},
				RequireLinearHistory: &githubsdk.RequireLinearHistory{Enabled: true},
				AllowForcePushes:     &githubsdk.AllowForcePushes{Enabled: true},
				AllowDeletions:       &githubsdk.AllowDeletions{Enabled: true},
			},
			ProtectionRequest: &githubsdk.ProtectionRequest{
				RequiredPullRequestReviews: &githubsdk.PullRequestReviewsEnforcementRequest{
					RequireCodeOwnerReviews:      true,
					RequiredApprovingReviewCount: 2,
					DismissalRestrictionsRequest: &githubsdk.DismissalRestrictionsRequest{
						Users: &[]string{"alex"},
						Teams: &[]string{"alpha"},
					},
				},
				Restrictions: &githubsdk.BranchRestrictionsRequest{
					Users: []string{"john"},
					Teams: []string{"easyCLA-Team"},
					Apps:  []string{},
				},
				EnforceAdmins: true,
				RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
					Strict:   true,
					Contexts: []string{"easyCLA"},
				},
				RequireLinearHistory: swag.Bool(true),
				AllowForcePushes:     swag.Bool(true),
				AllowDeletions:       swag.Bool(true),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			// Assert that Bar() is invoked.
			defer ctrl.Finish()

			m := NewMockRepositories(ctrl)
			m.
				EXPECT().
				GetBranchProtection(gomock.Any(), owner, repo, branchName).
				Return(tc.Protection, nil, nil)
			m.
				EXPECT().
				UpdateBranchProtection(gomock.Any(), owner, repo, branchName, tc.ProtectionRequest).
				Return(nil, nil, nil)

			branchProtectionRepo := NewBranchProtectionRepository(m)
			err := branchProtectionRepo.EnableBranchProtection(context.Background(), owner, repo, branchName, true, tc.Checks, nil)
			if err != nil {
				tt.Errorf("enable branch proteciton failed : %v", err)
			}
		})
	}

}

func TestNonBlockingRateLimitRepositories_GetBranchProtection(t *testing.T) {
	owner := "johnblocking"
	repo := "johnsrepoblocking"
	branchName := defaultBranchName

	t.Run("no limit reached", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		protection := &githubsdk.Protection{
			RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
				Strict:   false,
				Contexts: []string{"circle/ci"},
			},
		}

		m := NewMockRepositories(ctrl)
		m.
			EXPECT().
			GetBranchProtection(gomock.Any(), owner, repo, branchName).
			Return(protection, nil, nil)

		nonBlockLimitRepo := NewBranchProtectionRepository(m, EnableNonBlockingLimiter())
		p, err := nonBlockLimitRepo.GetProtectedBranch(context.Background(), owner, repo, branchName)
		if err != nil {
			tt.Errorf("no error expected : %v", err)
		}
		assert.Equal(tt, protection, p)
	})

	t.Run("limit reached", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		protection := &githubsdk.Protection{
			RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
				Strict:   false,
				Contexts: []string{"circle/ci"},
			},
		}

		m := NewMockRepositories(ctrl)
		m.
			EXPECT().
			GetBranchProtection(gomock.Any(), owner, repo, branchName).
			Return(protection, nil, nil).AnyTimes()

		nonBlockLimitRepo := NewBranchProtectionRepository(m, EnableNonBlockingLimiter())
		// call it 100 times in loop to make it fail
		var expectedErr error
		for i := 0; i < 100; i++ {
			_, err := nonBlockLimitRepo.GetProtectedBranch(context.Background(), owner, repo, branchName)
			if err != nil {
				expectedErr = err
				break
			}
		}

		if expectedErr == nil {
			tt.Fatalf("no error returned")
			return
		}

		if !errors.Is(expectedErr, ErrRateLimited) {
			tt.Fatalf("was expecting ErrRateLimited got : %v", expectedErr)
			return
		}
	})
}

func TestBlockingRateLimitRepositories_GetBranchProtection(t *testing.T) {
	owner := "john"
	repo := "johnsrepo"
	branchName := defaultBranchName

	t.Run("no limit reached", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		protection := &githubsdk.Protection{
			RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
				Strict:   false,
				Contexts: []string{"circle/ci"},
			},
		}

		m := NewMockRepositories(ctrl)
		m.
			EXPECT().
			GetBranchProtection(gomock.Any(), owner, repo, branchName).
			Return(protection, nil, nil)

		blockLimitRepo := NewBranchProtectionRepository(m, EnableBlockingLimiter())
		p, err := blockLimitRepo.GetProtectedBranch(context.Background(), owner, repo, branchName)
		if err != nil {
			tt.Errorf("no error expected : %v", err)
		}
		assert.Equal(tt, protection, p)
	})

	t.Run("limit reached", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		protection := &githubsdk.Protection{
			RequiredStatusChecks: &githubsdk.RequiredStatusChecks{
				Strict:   false,
				Contexts: []string{"circle/ci"},
			},
		}

		m := NewMockRepositories(ctrl)
		m.
			EXPECT().
			GetBranchProtection(gomock.Any(), owner, repo, branchName).
			Return(protection, nil, nil).AnyTimes()

		blockLimitRepo := NewBranchProtectionRepository(m, EnableBlockingLimiter())

		// call it 100 times in loop to make it fail
		var expectedErr error
		start := time.Now()
		for i := 0; i < 10; i++ {
			_, err := blockLimitRepo.GetProtectedBranch(context.Background(), owner, repo, branchName)
			if err != nil {
				expectedErr = err
				break
			}
		}
		elapsed := time.Since(start)

		if expectedErr != nil {
			tt.Fatalf("no error was expected got : %v", expectedErr)
			return
		}

		if elapsed < 4*time.Second {
			tt.Fatalf("is rate limit enabled")
		}

	})
}
