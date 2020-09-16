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
