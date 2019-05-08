package config

import (
	"log"

	"github.com/lytics/logrus"
)

type Config struct {
	// Auth0
	Auth0 Auth0 `json:"auth0"`

	// SFDC

	// GitHub

	// Docusign

	// LF Identity

	// AWS

}

type Auth0 struct {
	Domain        string `json:"auth0-domain"`
	ClientID      string `json:"auth0-clientId"`
	UsernameClaim string `json:"auth0-username-claim"`
	Algorithm     string `json:"auth0-algorithm"`
}

func LoadConfig(configFilePath string, awsStage string) (Config, error) {
	if configFilePath != "" {
		// Read from local env.json
		localConfig, err := loadLocalConfig(configFilePath)
		if err != nil {
			return Config{}, err
		}

		return localConfig, nil
	}

	logrus.Info("Local config not found")

	// Read from SSM
	log.Fatal("SSM Config not supported")

	return Config{}, nil
}
