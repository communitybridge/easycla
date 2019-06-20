// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later
package authorizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertTokenToHTTPRequest(t *testing.T) {
	token := "Bearer 13412341"
	request, err := convertTokenToHTTPRequest(token)

	assert.Equal(t, token, request.Header.Get("Authorization"))
	assert.Nil(t, err)
}

func TestExtractParsedTokenInfo(t *testing.T) {
	claims := map[string]interface{}{
		"email":          "user@example.com",
		"email_verified": true,
		"sub":            "google-oauth2|1111111111",
	}
	tokenInfo, err := extractParsedTokenInfo(claims)
	expected := TokenInfo{
		Email:         "user@example.com",
		EmailVerified: true,
		Subject:       "google-oauth2|1111111111",
	}
	assert.Equal(t, expected, tokenInfo)
	assert.Equal(t, err, nil)
}
