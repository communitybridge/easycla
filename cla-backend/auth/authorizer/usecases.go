// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later
package authorizer

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"gopkg.in/square/go-jose.v2/jwt"
)

type (
	// Usecases is an interface for interacting with the authorizer's
	// usecases layer.
	Usecases interface {
		ValidateToken(token string) (TokenInfo, error)
	}

	// TokenValidator defines the interface necessary for validating a token
	// and extracting claims from the token
	TokenValidator interface {
		ValidateRequest(*http.Request) (*jwt.JSONWebToken, error)
		Claims(*http.Request, *jwt.JSONWebToken, ...interface{}) error
	}

	// usecasesContainer holds initialized dependencies
	usecasesContainer struct {
		validator TokenValidator
	}
)

// NewUsecases create a new usecases
func NewUsecases(validator TokenValidator) Usecases {
	result := usecasesContainer{
		validator: validator,
	}
	return &result
}

// ValidateToken verifies a token
func (uc *usecasesContainer) ValidateToken(token string) (TokenInfo, error) {
	fields := make(map[string]interface{})
	fields["function_name"] = "authorizer.usecases.ValidateToken"
	log.Print(fields, "Entered function")

	r, err := convertTokenToHTTPRequest(token)
	if err != nil {
		log.Print(fields, err.Error())
		return TokenInfo{}, err
	}

	parsedToken, err := uc.validator.ValidateRequest(r)
	if err != nil {
		return TokenInfo{}, err
	}

	claims := map[string]interface{}{}
	err = uc.validator.Claims(r, parsedToken, &claims)
	if err != nil {
		return TokenInfo{}, err
	}

	return extractParsedTokenInfo(claims)
}

func convertTokenToHTTPRequest(token string) (*http.Request, error) {
	// Leave the body and url empty, they aren't used by the http request validator.
	body := strings.NewReader("")
	url := ""

	r, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", token)

	return r, nil
}

func extractParsedTokenInfo(claims map[string]interface{}) (TokenInfo, error) {
	email, emailExists := claims["email"].(string)
	emailVerified, emailVerifiedExists := claims["email_verified"].(bool)
	subject, subjectExists := claims["sub"].(string)

	if !emailExists {
		return TokenInfo{}, errors.New("token missing email claim")
	}
	if !emailVerifiedExists {
		return TokenInfo{}, errors.New("token missing email_verified claim")
	}
	if !subjectExists {
		return TokenInfo{}, errors.New("token missing subject claim")
	}

	return TokenInfo{
		Email:         email,
		EmailVerified: emailVerified,
		Subject:       subject,
	}, nil
}
