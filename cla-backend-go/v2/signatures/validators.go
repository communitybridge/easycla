// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime/middleware"
)

// validateApprovalListInput is a helper function to validate the update approval list input parameters
func validateApprovalListInput(params signatures.UpdateApprovalListParams) middleware.Responder {
	if !hasApprovalListUpdates(params) {
		return signatures.NewUpdateApprovalListBadRequest().WithPayload(errorResponse(errors.New("missing approval list items")))
	}

	msg, valid := entriesAreValid(params)
	if !valid {
		return signatures.NewUpdateApprovalListBadRequest().WithPayload(errorResponse(errors.New(msg)))
	}
	return nil
}

// hasApprovalListUpdates returns true if we have something to update, otherwise returns false
func hasApprovalListUpdates(params signatures.UpdateApprovalListParams) bool {
	if len(params.Body.AddEmailApprovalList) > 0 || len(params.Body.RemoveEmailApprovalList) > 0 ||
		len(params.Body.AddDomainApprovalList) > 0 || len(params.Body.RemoveDomainApprovalList) > 0 ||
		len(params.Body.AddGithubUsernameApprovalList) > 0 || len(params.Body.RemoveGithubUsernameApprovalList) > 0 ||
		len(params.Body.AddGithubOrgApprovalList) > 0 || len(params.Body.RemoveGithubOrgApprovalList) > 0 {
		return true
	}

	return false
}

// entriesAreValid returns true if the values in the approval list are valid, returns false and a message otherwise
func entriesAreValid(params signatures.UpdateApprovalListParams) (string, bool) {
	// Ensure the email address are valid
	for _, email := range params.Body.AddEmailApprovalList {
		if !utils.ValidEmail(email) {
			return fmt.Sprintf("invalid add approval list email %s", email), false
		}
	}
	for _, email := range params.Body.RemoveEmailApprovalList {
		if !utils.ValidEmail(email) {
			return fmt.Sprintf("invalid remove approval list email %s", email), false
		}
	}

	// Ensure the domains are valid
	for _, domain := range params.Body.AddDomainApprovalList {
		msg, valid := utils.ValidDomain(domain)
		if !valid {
			return fmt.Sprintf("invalid add approval list domain %s - %s", domain, msg), false
		}
	}
	for _, domain := range params.Body.RemoveDomainApprovalList {
		msg, valid := utils.ValidDomain(domain)
		if !valid {
			return fmt.Sprintf("invalid remove approval list domain %s - %s", domain, msg), false
		}
	}

	// Ensure the github usernames are valid
	for _, githubUsername := range params.Body.AddGithubUsernameApprovalList {
		msg, valid := utils.ValidGitHubUsername(githubUsername)
		if !valid {
			return fmt.Sprintf("invalid add approval list GitHub Username %s - %s", githubUsername, msg), false
		}
	}
	for _, githubUsername := range params.Body.RemoveGithubUsernameApprovalList {
		msg, valid := utils.ValidGitHubUsername(githubUsername)
		if !valid {
			return fmt.Sprintf("invalid remove approval list GitHub Username %s - %s", githubUsername, msg), false
		}
	}

	// Ensure the github Organization values are valid
	for _, githubOrg := range params.Body.AddGithubOrgApprovalList {
		msg, valid := utils.ValidGitHubOrg(githubOrg)
		if !valid {
			return fmt.Sprintf("invalid add approval list GitHub Org %s - %s", githubOrg, msg), false
		}
	}
	for _, githubOrg := range params.Body.RemoveGithubOrgApprovalList {
		msg, valid := utils.ValidGitHubOrg(githubOrg)
		if !valid {
			return fmt.Sprintf("invalid remove approval list GitHub Org %s - %s", githubOrg, msg), false
		}
	}

	return "", true
}
