// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package test_models

// Auth0Config is the configuration for Auth0
type Auth0Config struct {
	Auth0UserName string
	Auth0Password string
	Auth0ClientID string
}

// Auth0Response is the response model from an auth0 token exchange
type Auth0Response struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}
