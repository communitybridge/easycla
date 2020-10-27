// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"log"

	"github.com/communitybridge/easycla/cla-backend-go/cmd"
	"github.com/communitybridge/easycla/cla-backend-go/config"
	ini "github.com/communitybridge/easycla/cla-backend-go/init"
	token "github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/spf13/viper"
)

var (
	// version the application version
	version string

	// build/Commit the application build number
	commit string

	// branch the build branch
	branch string

	// build date
	buildDate string

	configFile = ""
)

func init() {
	stage := viper.GetString("STAGE")

	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Panicf("Unable to load AWS session - Error: %v", err)
	}

	configFile, err := config.LoadConfig(configFile, awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}
	token.Init(configFile.Auth0Platform.ClientID, configFile.Auth0Platform.ClientSecret, configFile.Auth0Platform.URL, configFile.Auth0Platform.Audience)
}

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Branch = branch
	cmd.BuildDate = buildDate

	cmd.Execute()
}
