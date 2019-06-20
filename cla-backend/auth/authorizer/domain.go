// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later
package authorizer

// TokenInfo represents claims present in the token
type TokenInfo struct {
	Email         string
	EmailVerified bool
	Subject       string
}
