// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// getSSMString is a generic routine to fetch the specified key value
func getSSMString(ssmClient *ssm.SSM, key string) (string, error) {
	log.Debugf("Loading SSM parameter: %s", key)
	value, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		log.Warnf("unable to read SSM parameter %s - error: %+v", key, err)
		return "", err
	}

	return strings.TrimSpace(*value.Parameter.Value), nil
}

func loadSSMConfig(awsSession *session.Session, stage string) (Config, error) {
	config := Config{}

	ssmClient := ssm.New(awsSession)

	// TODO: DAD - optimization: fetch the SSM values in parallel

	// Auth0
	auth0Domain, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-domain-%s", stage))
	if err != nil {
		return Config{}, err
	}

	auth0ClientID, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-clientId-%s", stage))
	if err != nil {
		return Config{}, err
	}

	auth0Username, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-username-claim-%s", stage))
	if err != nil {
		return Config{}, err
	}

	auth0Algorithm, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-algorithm-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.Auth0 = Auth0{
		Domain:        auth0Domain,
		ClientID:      auth0ClientID,
		UsernameClaim: auth0Username,
		Algorithm:     auth0Algorithm,
	}

	// SFDC

	// GitHub
	githubClientID, err := getSSMString(ssmClient, fmt.Sprintf("cla-gh-oauth-client-id-go-backend-%s", stage))
	if err != nil {
		return Config{}, err
	}
	githubSecret, err := getSSMString(ssmClient, fmt.Sprintf("cla-gh-oauth-secret-go-backend-%s", stage))
	if err != nil {
		return Config{}, err
	}
	githubAccessToken, err := getSSMString(ssmClient, fmt.Sprintf("cla-gh-access-token-%s", stage))
	if err != nil {
		return Config{}, err
	}
	ghAppID, err := getSSMString(ssmClient, fmt.Sprintf("cla-gh-app-id-%s", stage))
	if err != nil {
		return Config{}, err
	}
	githubAppID, err := strconv.Atoi(ghAppID)
	if err != nil {
		return Config{}, err
	}
	githubAppPrivateKey, err := getSSMString(ssmClient, fmt.Sprintf("cla-gh-app-private-key-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.Github = Github{
		ClientID:      githubClientID,
		ClientSecret:  githubSecret,
		AccessToken:   githubAccessToken,
		AppPrivateKey: githubAppPrivateKey,
		AppID:         githubAppID,
	}

	//Corporate Console Link
	corporateConsoleURL, err := getSSMString(ssmClient, fmt.Sprintf("cla-corporate-base-%s", stage))
	if err != nil {
		return Config{}, err
	}
	corporateConsoleURLValue := corporateConsoleURL
	if corporateConsoleURLValue == "corporate.prod.lfcla.com" {
		corporateConsoleURLValue = "corporate.lfcla.com"
	}
	config.CorporateConsoleURL = corporateConsoleURLValue

	// Docusign

	// Docraptor
	config.Docraptor.APIKey, err = getSSMString(ssmClient, fmt.Sprintf("cla-doc-raptor-api-key-%s", stage))
	if err != nil {
		return Config{}, err
	}
	config.Docraptor.TestMode = stage != "prod" && stage != "staging"

	// LF Identity

	// AWS
	config.AWS.Region = "us-east-1"

	// Session Store Table Name
	config.SessionStoreTableName, err = getSSMString(ssmClient, fmt.Sprintf("cla-session-store-table-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.SenderEmailAddress, err = getSSMString(ssmClient, fmt.Sprintf("cla-ses-sender-email-address-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.AllowedOriginsCommaSeparated, err = getSSMString(ssmClient, fmt.Sprintf("cla-allowed-origins-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.SNSEventTopicARN, err = getSSMString(ssmClient, fmt.Sprintf("cla-sns-event-topic-arn-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.SignatureFilesBucket, err = getSSMString(ssmClient, fmt.Sprintf("cla-signature-files-bucket-%s", stage))
	if err != nil {
		return Config{}, err
	}

	auth0PlatformClientID, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-platform-client-id-%s", stage))
	if err != nil {
		return Config{}, err
	}
	auth0PlatformSecret, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-platform-client-secret-%s", stage))
	if err != nil {
		return Config{}, err
	}
	auth0PlatformAudience, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-platform-audience-%s", stage))
	if err != nil {
		return Config{}, err
	}
	auth0PlatformURL, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-platform-url-%s", stage))
	if err != nil {
		return Config{}, err
	}
	apiGw, err := getSSMString(ssmClient, fmt.Sprintf("cla-auth0-platform-api-gw-%s", stage))
	if err != nil {
		return Config{}, err
	}
	config.Auth0Platform = Auth0Platform{
		ClientID:     auth0PlatformClientID,
		ClientSecret: auth0PlatformSecret,
		Audience:     auth0PlatformAudience,
		URL:          auth0PlatformURL,
	}
	config.APIGatewayURL = apiGw

	return config, nil
}
