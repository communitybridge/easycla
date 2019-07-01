// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package authorizer

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAuthResponseHasContext(t *testing.T) {
	token := TokenInfo{
		Email:         "user@example.com",
		EmailVerified: true,
		Subject:       "google-oauth2|1111111111",
	}
	methodArn := "arn:aws:execute-api:us-east-1:xxxxx:xxxxx/stage/some-method"
	result := generateAuthResponse(methodArn, &token)
	expected := map[string]interface{}{
		"subject":       token.Subject,
		"email":         token.Email,
		"emailVerified": token.EmailVerified,
	}
	assert.Equal(t, expected, result.Context)
}

func TestGenerateAuthResponseHasPrincipleID(t *testing.T) {
	token := TokenInfo{
		Email:         "user@example.com",
		EmailVerified: true,
		Subject:       "google-oauth2|1111111111",
	}
	methodArn := "arn:aws:execute-api:us-east-1:xxxxx:xxxxx/stage/some-method"
	result := generateAuthResponse(methodArn, &token)
	expected := "google-oauth2|1111111111"
	assert.Equal(t, expected, result.PrincipalID)
}

func TestGeneratePolicy(t *testing.T) {
	methodArn := "arn:aws:execute-api:us-east-1:xxxxx:xxxxx/stage/some-method"
	result := generatePolicy(methodArn)
	expected := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{"arn:aws:execute-api:us-east-1:xxxxx:xxxxx/stage/*"},
			},
		},
	}
	assert.Equal(t, expected, result)
}
