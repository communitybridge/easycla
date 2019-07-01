// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package authorizer

import (
	"errors"
	"log"
	"os"

	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
)

// ValidatorMaker is the interface for interacting with the validator layer.
type ValidatorMaker interface {
	NewTokenValidator() (TokenValidator, error)
}

// validatorContainer holds initialized dependencies for the validator.
type validatorContainer struct{}

// NewValidatorMaker create a new ValidatorMaker.
func NewValidatorMaker() *validatorContainer {
	return &validatorContainer{}
}

// NewTokenValidator creates a token validator, reading configuration from the environment.
func (vc *validatorContainer) NewTokenValidator() (TokenValidator, error) {
	fields := make(map[string]interface{})
	fields["function_name"] = "authorizer.validator.NewTokenValidator"
	log.Print(fields, "Entered function")

	authDomain := os.Getenv("AUTH0_DOMAIN")
	if len(authDomain) == 0 {
		errMsg := "couldn't find auth0 URI"
		log.Print(fields, errMsg)
		return nil, errors.New(errMsg)
	}

	url := "https://" + authDomain + "/"
	aud := os.Getenv("AUTH0_CLIENT_ID")
	if len(aud) == 0 {
		errMsg := "couldn't find auth0 client id"
		log.Print(fields, errMsg)
		return nil, errors.New(errMsg)
	}
	return vc.createTokenValidator(url, aud), nil
}

// newTokenValidator creates a token validator,
func (vc *validatorContainer) createTokenValidator(domain string, audience string) TokenValidator {
	fields := make(map[string]interface{})
	fields["function_name"] = "authorizer.validator.NewTokenValidator"
	log.Print(fields, "Entered function")

	uri := domain + ".well-known/jwks.json"
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: uri}, nil)
	log.Print(fields, "Initialized JWKClient")

	configuration := auth0.NewConfiguration(client, []string{audience}, domain, jose.RS256)
	validator := auth0.NewValidator(configuration, nil)

	log.Print(fields, "Successfully created validator")
	return validator
}
