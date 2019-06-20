// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

package config

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/lytics/logrus"
)

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

	AllowedOriginsCommaSeparated string              `json:"allowedOriginsCommaSeparated"`
	AllowedOrigins               map[string]struct{} `json:"-"`

	CorporateConsoleURL string `json:"corporateConsoleURL"`
}

type Auth0 struct {
	Domain        string `json:"auth0-domain"`
	ClientID      string `json:"auth0-clientId"`
	UsernameClaim string `json:"auth0-username-claim"`
	Algorithm     string `json:"auth0-algorithm"`
}

type Docraptor struct {
	APIKey   string `json:"apiKey"`
	TestMode bool   `json:"testMode"`
}

type AWS struct {
	Region string `json:"region"`
}

type Github struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func LoadConfig(configFilePath string, awsSession *session.Session, awsStage string) (Config, error) {
	var config Config
	var err error

	if configFilePath != "" {
		// Read from local env.jso
		logrus.Info("Loading local config")
		config, err = loadLocalConfig(configFilePath)

	} else if awsSession != nil {
		// Read from SSM
		logrus.Info("Loading SSM config")
		config, err = loadSSMConfig(awsSession, awsStage)

	} else {
		return Config{}, errors.New("config not found")
	}

	if err != nil {
		return Config{}, err
	}

	config.AllowedOrigins = map[string]struct{}{}
	for _, origin := range strings.Split(config.AllowedOriginsCommaSeparated, ",") {
		config.AllowedOrigins[origin] = struct{}{}
	}

	return config, nil
}
