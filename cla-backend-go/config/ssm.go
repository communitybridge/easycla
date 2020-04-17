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

// configLookupResponse is a channel response model for the configuration lookup
type configLookupResponse struct {
	key   string
	value string
}

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

// loadSSMConfig fetches all the configuration values and populates the response Config model
func loadSSMConfig(awsSession *session.Session, stage string) Config { //nolint
	config := Config{}

	ssmClient := ssm.New(awsSession)

	// AWS Region
	config.AWS.Region = *awsSession.Config.Region

	// A channel for the responses from the go routines
	responseChannel := make(chan configLookupResponse)

	// Our list of keys to lookup
	ssmKeys := []string{
		fmt.Sprintf("cla-auth0-domain-%s", stage),
		fmt.Sprintf("cla-auth0-clientId-%s", stage),
		fmt.Sprintf("cla-auth0-username-claim-%s", stage),
		fmt.Sprintf("cla-auth0-algorithm-%s", stage),
		fmt.Sprintf("cla-gh-oauth-client-id-go-backend-%s", stage),
		fmt.Sprintf("cla-gh-oauth-secret-go-backend-%s", stage),
		fmt.Sprintf("cla-gh-access-token-%s", stage),
		fmt.Sprintf("cla-gh-app-id-%s", stage),
		fmt.Sprintf("cla-gh-app-private-key-%s", stage),
		fmt.Sprintf("cla-corporate-base-%s", stage),
		fmt.Sprintf("cla-doc-raptor-api-key-%s", stage),
		fmt.Sprintf("cla-session-store-table-%s", stage),
		fmt.Sprintf("cla-ses-sender-email-address-%s", stage),
		fmt.Sprintf("cla-allowed-origins-%s", stage),
		fmt.Sprintf("cla-sns-event-topic-arn-%s", stage),
		fmt.Sprintf("cla-signature-files-bucket-%s", stage),
		fmt.Sprintf("cla-auth0-platform-client-id-%s", stage),
		fmt.Sprintf("cla-auth0-platform-client-secret-%s", stage),
		fmt.Sprintf("cla-auth0-platform-audience-%s", stage),
		fmt.Sprintf("cla-auth0-platform-url-%s", stage),
		fmt.Sprintf("cla-auth0-platform-api-gw-%s", stage),
		fmt.Sprintf("cla-lf-group-client-id-%s", stage),
		fmt.Sprintf("cla-lf-group-client-secret-%s", stage),
		fmt.Sprintf("cla-lf-group-client-url-%s", stage),
		fmt.Sprintf("cla-lf-group-refresh-token-%s", stage),
	}

	// For each key to lookup
	for _, key := range ssmKeys {
		// Create a go routine to this concurrently
		go func(theKey string) {
			theValue, err := getSSMString(ssmClient, theKey)
			if err != nil {
				log.Fatalf("error looking up key: %s", theKey)
			}
			// Send the response back through the channel
			responseChannel <- configLookupResponse{
				key:   theKey,
				value: theValue,
			}
		}(key)
	}

	for i := 0; i < len(ssmKeys); i++ {
		resp := <-responseChannel
		switch resp.key {
		case fmt.Sprintf("cla-auth0-domain-%s", stage):
			config.Auth0.Domain = resp.value
		case fmt.Sprintf("cla-auth0-clientId-%s", stage):
			config.Auth0.ClientID = resp.value
		case fmt.Sprintf("cla-auth0-username-claim-%s", stage):
			config.Auth0.UsernameClaim = resp.value
		case fmt.Sprintf("cla-auth0-algorithm-%s", stage):
			config.Auth0.Algorithm = resp.value
		case fmt.Sprintf("cla-gh-oauth-client-id-go-backend-%s", stage):
			config.Github.ClientID = resp.value
		case fmt.Sprintf("cla-gh-oauth-secret-go-backend-%s", stage):
			config.Github.ClientSecret = resp.value
		case fmt.Sprintf("cla-gh-access-token-%s", stage):
			config.Github.AccessToken = resp.value
		case fmt.Sprintf("cla-gh-app-id-%s", stage):
			githubAppID, err := strconv.Atoi(resp.value)
			if err != nil {
				errMsg := fmt.Sprintf("invalid value of key: %s", fmt.Sprintf("cla-gh-app-id-%s", stage))
				log.Fatal(errMsg)
			}
			config.Github.AppID = githubAppID
		case fmt.Sprintf("cla-gh-app-private-key-%s", stage):
			config.Github.AppPrivateKey = resp.value

		case fmt.Sprintf("cla-corporate-base-%s", stage):
			corporateConsoleURLValue := resp.value
			if corporateConsoleURLValue == "corporate.prod.lfcla.com" {
				corporateConsoleURLValue = "corporate.lfcla.com"
			}
			config.CorporateConsoleURL = corporateConsoleURLValue
		case fmt.Sprintf("cla-doc-raptor-api-key-%s", stage):
			config.Docraptor.APIKey = resp.value
			config.Docraptor.TestMode = stage != "prod" && stage != "staging"
		case fmt.Sprintf("cla-session-store-table-%s", stage):
			config.SessionStoreTableName = resp.value
		case fmt.Sprintf("cla-ses-sender-email-address-%s", stage):
			config.SenderEmailAddress = resp.value
		case fmt.Sprintf("cla-allowed-origins-%s", stage):
			config.AllowedOriginsCommaSeparated = resp.value
		case fmt.Sprintf("cla-sns-event-topic-arn-%s", stage):
			config.SNSEventTopicARN = resp.value
		case fmt.Sprintf("cla-signature-files-bucket-%s", stage):
			config.SignatureFilesBucket = resp.value
		case fmt.Sprintf("cla-auth0-platform-client-id-%s", stage):
			config.Auth0Platform.ClientID = resp.value
		case fmt.Sprintf("cla-auth0-platform-client-secret-%s", stage):
			config.Auth0Platform.ClientSecret = resp.value
		case fmt.Sprintf("cla-auth0-platform-audience-%s", stage):
			config.Auth0Platform.Audience = resp.value
		case fmt.Sprintf("cla-auth0-platform-url-%s", stage):
			config.Auth0Platform.URL = resp.value
		case fmt.Sprintf("cla-auth0-platform-api-gw-%s", stage):
			config.APIGatewayURL = resp.value
		case fmt.Sprintf("cla-lf-group-client-id-%s", stage):
			config.LFGroup.ClientID = resp.value
		case fmt.Sprintf("cla-lf-group-client-secret-%s", stage):
			config.LFGroup.ClientSecret = resp.value
		case fmt.Sprintf("cla-lf-group-client-url-%s", stage):
			config.LFGroup.ClientURL = resp.value
		case fmt.Sprintf("cla-lf-group-refresh-token-%s", stage):
			config.LFGroup.RefreshToken = resp.value
		}
	}

	return config
}
