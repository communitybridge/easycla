// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import "strings"

// TrimRemoveTrailingComma trims the whitespace on the specified string and removes the trailing comma
func TrimRemoveTrailingComma(input string) string {
	if input == "" {
		return input
	}

	s := strings.TrimSpace(input)
	return strings.TrimSuffix(s, ",")
}

// TrimSpaceFromItems is a helper function to trim space on array items
func TrimSpaceFromItems(arr []string) []string {
	newArr := make([]string, len(arr))
	for i := range arr {
		newArr[i] = strings.TrimSpace(arr[i])
	}

	return newArr
}
