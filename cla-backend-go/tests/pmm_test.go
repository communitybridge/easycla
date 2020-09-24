// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package tests

import (
	"testing"

	project_service "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"

	"github.com/communitybridge/easycla/cla-backend-go/config"
	ini "github.com/communitybridge/easycla/cla-backend-go/init"
	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetLinuxFoundationProject(t *testing.T) {

	if false {

		viper.AutomaticEnv()
		defaults := map[string]interface{}{
			"PORT":                    8080,
			"APP_ENV":                 "local",
			"USE_MOCK":                "False",
			"DB_MAX_CONNECTIONS":      1,
			"STAGE":                   "dev",
			"GH_ORG_VALIDATION":       "false",
			"COMPANY_USER_VALIDATION": "false",
			"LOG_FORMAT":              "text",
		}

		for key, value := range defaults {
			viper.SetDefault(key, value)
		}

		ini.Init()
		ini.AWSInit()
		awsSession, err := ini.GetAWSSession()
		assert.Nil(t, err, "aws session returns success")
		assert.NotNil(t, awsSession, "valid AWS session")
		stage := viper.GetString("STAGE")
		configFile, err := config.LoadConfig("", awsSession, stage)
		assert.Nil(t, err, "get config returns success")
		assert.NotNil(t, configFile, "valid config file")

		token.Init(configFile.Auth0Platform.ClientID, configFile.Auth0Platform.ClientSecret, configFile.Auth0Platform.URL, configFile.Auth0Platform.Audience)
		project_service.InitClient(configFile.APIGatewayURL)
		client := project_service.GetClient()
		projectDetails, err := client.GetLinuxFoundationProject()
		assert.Nil(t, err, "get linux foundation returns success")
		assert.NotNil(t, projectDetails, "project details is not nil")
	}
}
