// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/LF-Engineering/lfx-models/models/event"
	usersModels "github.com/LF-Engineering/lfx-models/models/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/communitybridge/easycla/cla-backend-go/config"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/communitybridge/easycla/cla-backend-go/userSubscribeLambda/cmd"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	user_service "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"
	"github.com/mitchellh/mapstructure"
)

// Build and version variables defined and set during the build process
var (
	// version the application version
	version string

	// build/Commit the application build number
	commit string

	// build date
	buildDate string
)

func init() {
	var awsSession = session.Must(session.NewSession(&aws.Config{}))
	stage := os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	log.Infof("STAGE set to %s\n", stage)
	configFile, err := config.LoadConfig("", awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}

	token.Init(configFile.Auth0Platform.ClientID, configFile.Auth0Platform.ClientSecret, configFile.Auth0Platform.URL, configFile.Auth0Platform.Audience)
	user_service.InitClient(configFile.APIGatewayURL, configFile.AcsAPIKey)
}

// Handler is the user subscribe handler lambda entry function
func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	if len(snsEvent.Records) == 0 {
		log.Warn("SNS event contained 0 records - ignoring message.")
		return nil
	}

	for _, message := range snsEvent.Records {
		log.Infof("Processing message id: '%s' for event source '%s'\n", message.SNS.MessageID, message.EventSource)

		log.Debugf("Unmarshalling message body: '%s'", message.SNS.Message)

		// log.Debugf("Unmarshalling message body: '%s'", message.SNS.Message)
		var model event.Event
		err := model.UnmarshalBinary([]byte(message.SNS.Message))
		if err != nil {
			log.Warnf("Error: %v, JSON unmarshal failed - unable to process message: %s", err, message.SNS.MessageID)
			return err
		}

		switch model.Type {
		case "UserUpdatedProfile":
			Write(model)
		default:
			log.Warnf("unrecognized message type: %s - unable to process message ", model.Type)
		}

	}
	return nil
}

// Write saves the user data model to persistent storage
func Write(user event.Event) {

	uc := &usersModels.UserUpdated{}
	err := mapstructure.Decode(user.Data, uc)
	if err != nil {
		return
	}

	var userDetails *models.User
	var userErr error
	var awsSession = session.Must(session.NewSession(&aws.Config{}))

	stage := os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	usersRepo := users.NewRepository(awsSession, stage)

	userDetails, userErr = usersRepo.GetUserByLFUserName(*uc.Username)
	if userErr != nil {
		log.Warnf("Error - unable to locate user by LfUsername: %s, error: %+v", *uc.Username, userErr)
		log.Error("", userErr)
	}

	if userDetails == nil {
		for _, email := range uc.Emails {
			userDetails, userErr = usersRepo.GetUserByEmail(*email.EmailAddress)
			if userErr != nil {
				log.Warnf("Error - unable to locate user by LfUsername: %s, error: %+v", *uc.Username, userErr)
			}
		}
	}

	if userDetails == nil {
		userDetails, userErr = usersRepo.GetUserByExternalID(uc.UserID)
		if userErr != nil {
			log.Warnf("Error - unable to locate user by UserExternalID: %s, error: %+v", uc.UserID, userErr)
		}
	}

	if userDetails == nil {
		log.Debugf("User model is nil so skipping user %s", *uc.Username)
		return
	}

	userServiceClient := user_service.GetClient()

	sfdcUserObject, err := userServiceClient.GetUser(uc.UserID)
	if err != nil {
		log.Warnf("Error - unable to locate user by SFID: %s, error: %+v", uc.UserID, userErr)
		log.Error("", userErr)
		return
	}

	log.Debugf("Salesforce user-service object : %+v", sfdcUserObject)

	if sfdcUserObject == nil {
		log.Debugf("User-service model is nil so skipping user %s with SFID %s", *uc.Username, uc.UserID)
		return
	}

	var primaryEmail string
	var emails []string
	for _, email := range sfdcUserObject.Emails {
		if *email.IsPrimary {
			primaryEmail = *email.EmailAddress
		}
		emails = append(emails, *email.EmailAddress)
	}

	updateUserModel := &models.UserUpdate{
		LfEmail:        primaryEmail,
		LfUsername:     sfdcUserObject.Username,
		Note:           "Update via user-service event",
		UserExternalID: sfdcUserObject.ID,
		UserID:         userDetails.UserID,
		Username:       fmt.Sprintf("%s %s", sfdcUserObject.FirstName, sfdcUserObject.LastName),
		Emails:         emails,
	}

	log.Debugf("Updating user in Dynamo DB : %+v", updateUserModel)

	_, updateErr := usersRepo.Save(updateUserModel)
	if updateErr != nil {
		log.Warnf("Error - unable to update user by LfUsername: %s, error: %+v", *uc.Username, updateErr)
		return
	}
}

func main() {
	var err error

	// Show the version and build info
	log.Infof("Name                  : userSubscribe handler")
	log.Infof("Version               : %s", version)
	log.Infof("Git commit hash       : %s", commit)
	log.Infof("Build date            : %s", buildDate)
	log.Infof("Golang OS             : %s", runtime.GOOS)
	log.Infof("Golang Arch           : %s", runtime.GOARCH)

	err = cmd.Start(Handler)
	if err != nil {
		log.Fatal(err)
	}
}
