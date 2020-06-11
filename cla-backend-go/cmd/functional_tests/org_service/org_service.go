// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package org_service

import (
	"github.com/communitybridge/easycla/cla-backend-go/cmd/functional_tests/test_models"
	"github.com/communitybridge/easycla/cla-backend-go/config"
	ini "github.com/communitybridge/easycla/cla-backend-go/init"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"
	organization_service "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"
	"github.com/spf13/viper"
)

// TestBehaviour data model
type TestBehaviour struct {
	apiURL           string
	auth0User1Config test_models.Auth0Config
	auth0User2Config test_models.Auth0Config
}

// NewTestBehaviour creates a new test behavior model
func NewTestBehaviour(apiURL string, auth0User1Config, auth0User2Config test_models.Auth0Config) *TestBehaviour {
	return &TestBehaviour{
		apiURL + "/v3",
		auth0User1Config,
		auth0User2Config,
	}
}

// RunIsUserHaveRoleScope test
func (t *TestBehaviour) RunIsUserHaveRoleScope() {
	stage := viper.GetString("STAGE")

	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Panicf("Unable to load AWS session - Error: %v", err)
	}

	configFile, err := config.LoadConfig("", awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}

	organization_service.InitClient(configFile.APIGatewayURL)
	acs_service.InitClient(configFile.APIGatewayURL, configFile.AcsAPIKey)
	acsClient := acs_service.GetClient()
	roleID, roleErr := acsClient.GetRoleID("cla-manager")
	if roleErr != nil {
		log.Fatalf("unable to read role: cla-manager from ACS Client, error: %+v", roleErr)
	}

	log.Debugf("Role ID for cla-manager-role : %s", roleID)
	orgClient := organization_service.GetClient()
	userSFID := "clamanager1devintel"
	hasScope, err := orgClient.IsUserHaveRoleScope(roleID, userSFID, "00117000015vpjXAAQ", "a092M00001IfVmKQAV")
	if err != nil {
		log.Fatalf("unable to invoke org client IsUserHaveRoleScope, error: %+v", err)
	}
	if !hasScope {
		log.Fatalf("user: %s does not have scope for a092M00001IfVmKQAV|00117000015vpjXAAQ", userSFID)
	}
}

// RunAllTests runs all the Organization Service Tests
func (t *TestBehaviour) RunAllTests() {
	t.RunIsUserHaveRoleScope()
}
