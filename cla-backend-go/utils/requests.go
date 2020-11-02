// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

// GetRequestID helper function to get the request ID from the session
func GetRequestID(reqID *string) string {
	if reqID != nil {
		return *reqID
	}

	return ""
}

// GetGithubEvent helper function to get the github Event Type from the session
func GetGithubEvent(xGithubEvent *string) string {
	if xGithubEvent != nil {
		return *xGithubEvent
	}

	return ""
}

// GetGithubSignature helper function to get the github Event Type from the session
func GetGithubSignature(xGithubSignature *string) string {
	if xGithubSignature != nil {
		return *xGithubSignature
	}

	return ""
}
