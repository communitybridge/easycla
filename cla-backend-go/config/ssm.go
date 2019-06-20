package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func loadSSMConfig(awsSession *session.Session, stage string) (Config, error) {
	config := Config{}

	ssmClient := ssm.New(awsSession)

	// Auth0
	auth0Domain, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-auth0-domain-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	auth0ClientID, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-auth0-clientId-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	auth0Username, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-auth0-username-claim-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	auth0Algorithm, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-auth0-algorithm-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	config.Auth0 = Auth0{
		Domain:        *auth0Domain.Parameter.Value,
		ClientID:      *auth0ClientID.Parameter.Value,
		UsernameClaim: *auth0Username.Parameter.Value,
		Algorithm:     *auth0Algorithm.Parameter.Value,
	}

	// SFDC

	// GitHub
	githubClientID, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-gh-oauth-client-id-go-backend-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}
	githubSecret, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-gh-oauth-secret-go-backend-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	config.Github = Github{
		ClientID:     *githubClientID.Parameter.Value,
		ClientSecret: *githubSecret.Parameter.Value,
	}

	//Corporate Console Link
	corporateConsoleURL, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-corporate-base-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}
	corporateConsoleURLValue := *corporateConsoleURL.Parameter.Value
	if corporateConsoleURLValue == "corporate.prod.lfcla.com" {
		corporateConsoleURLValue = "corporate.lfcla.com"
	}
	config.CorporateConsoleURL = corporateConsoleURLValue

	// Docusign

	// Docraptor
	docraptorAPIKey, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-doc-raptor-api-key-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	config.Docraptor.APIKey = *docraptorAPIKey.Parameter.Value
	config.Docraptor.TestMode = stage != "prod" && stage != "staging"

	// LF Identity

	// AWS
	config.AWS.Region = "us-east-1"

	// Session Store Table Name
	sessionStoreTableName, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-session-store-table-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}
	config.SessionStoreTableName = *sessionStoreTableName.Parameter.Value

	senderEmailAddress, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-ses-sender-email-address-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	allowedOrigins, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("cla-allowed-origins-%s", stage)),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return Config{}, err
	}

	config.SenderEmailAddress = *senderEmailAddress.Parameter.Value

	config.AllowedOriginsCommaSeparated = *allowedOrigins.Parameter.Value

	return config, nil
}
