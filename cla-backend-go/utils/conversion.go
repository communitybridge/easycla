// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

// StringValue function convert string pointer to string
func StringValue(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

// Int64Value function convert int64 pointer to string
func Int64Value(input *int64) int64 {
	if input == nil {
		return 0
	}
	return *input
}
