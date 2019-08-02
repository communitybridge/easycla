// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Config data model
type Config struct {
	// Auth0
	Auth0 Auth0 `json:"auth0"`

	// SFDC

	// GitHub

	// Docusign

	// Docraptor
	Docraptor Docraptor `json:"docraptor"`

	// LF Identity

	// AWS
	AWS AWS `json:"aws"`

	// Github Application
	Github Github `json:"github"`

	// Dynamo Session Store
	SessionStoreTableName string `json:"sessionStoreTableName"`

	// Sender Email Address
	SenderEmailAddress string `json:"senderEmailAddress"`

	AllowedOriginsCommaSeparated string   `json:"allowedOriginsCommaSeparated"`
	AllowedOrigins               []string `json:"-"`

	CorporateConsoleURL string `json:"corporateConsoleURL"`
}

// Auth0 model
type Auth0 struct {
	Domain        string `json:"auth0-domain"`
	ClientID      string `json:"auth0-clientId"`
	UsernameClaim string `json:"auth0-username-claim"`
	Algorithm     string `json:"auth0-algorithm"`
}

// Docraptor model
type Docraptor struct {
	APIKey   string `json:"apiKey"`
	TestMode bool   `json:"testMode"`
}

// AWS model
type AWS struct {
	Region string `json:"region"`
}

// Github model
type Github struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// LoadConfig loads the configuration
func LoadConfig(configFilePath string, awsSession *session.Session, awsStage string) (Config, error) {
	var config Config
	var err error

	if configFilePath != "" {
		// Read from local env.jso
		log.Info("Loading local config...")
		config, err = loadLocalConfig(configFilePath)

	} else if awsSession != nil {
		// Read from SSM
		log.Info("Loading SSM config...")
		config, err = loadSSMConfig(awsSession, awsStage)

	} else {
		return Config{}, errors.New("config not found")
	}

	if err != nil {
		return Config{}, err
	}

	// Convert the allowed origins into an array of values
	config.AllowedOrigins = strings.Split(config.AllowedOriginsCommaSeparated, ",")

	return config, nil
}
