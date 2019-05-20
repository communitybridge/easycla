package config

import (
	"errors"

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

func LoadConfig(configFilePath string, awsSession *session.Session, awsStage string) (Config, error) {
	if configFilePath != "" {
		logrus.Info("Loading local config")

		// Read from local env.json
		localConfig, err := loadLocalConfig(configFilePath)
		if err != nil {
			return Config{}, err
		}

		return localConfig, nil
	}

	// Read from SSM
	if awsSession != nil {
		logrus.Info("Loading SSM config")

		ssmConfig, err := loadSSMConfig(awsSession, awsStage)
		if err != nil {
			return Config{}, err
		}

		return ssmConfig, nil
	}

	return Config{}, errors.New("config not found")
}
