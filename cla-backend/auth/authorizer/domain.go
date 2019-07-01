// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package authorizer

// TokenInfo represents claims present in the token
type TokenInfo struct {
	Email         string
	EmailVerified bool
	Subject       string
}
