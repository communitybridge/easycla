// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"testing"

	"github.com/bmizerany/assert"
	githubsdk "github.com/google/go-github/github"
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
