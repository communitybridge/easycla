// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

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
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		log.Infof("Processing message %s for event source %s\n", message.MessageId, message.EventSource)

		userData := EventSchemaData{}
		err := json.Unmarshal([]byte(message.Body), &userData)
		if err != nil {
			log.Warnf("Error: %v, JSON unmarshal failed - unable to process message: %s", err, message.MessageId)
		}
		Write(userData)
	}
	return nil
}

// Write saves the user data model to persistent storage
func Write(user EventSchemaData) {
	var userDetails *models.User
	var userErr error
	var awsSession = session.Must(session.NewSession(&aws.Config{}))
	stage := os.Getenv("STAGE")
	if stage == "" {
		log.Fatal("stage not set")
	}
	usersRepo := users.NewRepository(awsSession, stage)

	userDetails, userErr = usersRepo.GetUserByLFUserName(user.Data.UserName)
	if userErr != nil {
		log.Warnf("Error - unable to locate user by LfUsername: %s, error: %+v", user.Data.UserName, userErr)
		log.Error("", userErr)
		return
	}

	if userDetails == nil {
		userDetails, userErr = usersRepo.GetUserByEmail(user.Data.Email)
		if userErr != nil {
			log.Warnf("Error - unable to locate user by LfUsername: %s, error: %+v", user.Data.UserName, userErr)
			return
		}
	}

	if userDetails == nil {
		userDetails, userErr = usersRepo.GetUserByExternalID(user.Data.UserID)
		if userErr != nil {
			log.Warnf("Error - unable to locate user by UserExternalID: %s, error: %+v", user.Data.UserID, userErr)
			return
		}
	}

	if userDetails == nil {
		log.Debugf("User model is nil so skipping user %s", user.Data.UserName)
		return
	}

	userServiceClient := user_service.GetClient()

	sfdcUserObject, err := userServiceClient.GetUser(user.Data.UserID)
	if err != nil {
		log.Warnf("Error - unable to locate user by SFID: %s, error: %+v", user.Data.UserID, userErr)
		log.Error("", userErr)
		return
	}

	if sfdcUserObject == nil {
		log.Debugf("User-service model is nil so skipping user %s with SFID %s", user.Data.UserName, user.Data.UserID)
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

	_, updateErr := usersRepo.Save(updateUserModel)
	if updateErr != nil {
		log.Warnf("Error - unable to update user by LfUsername: %s, error: %+v", user.Data.UserName, updateErr)
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

// UserData . . .
type UserData struct {
	UserID    string `json:"userId,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	Title     string `json:"title,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Type      string `json:"type,omitempty"`
	UserName  string `json:"username,omitempty"`
	Picture   string `json:"picture,omitempty"`
}

// EventSchemaData . . .
type EventSchemaData struct {
	Schema struct {
		Type       string `json:"type"`
		Properties struct {
			UserID struct {
				Type string `json:"type"`
			} `json:"userID"`
			FirstName struct {
				Type string `json:"type"`
			} `json:"firstName"`
			LastName struct {
				Type string `json:"type"`
			} `json:"lastName"`
			Email struct {
				Type string `json:"type"`
			} `json:"email"`
			Title struct {
				Type string `json:"type"`
			} `json:"title"`
			Username struct {
				Type string `json:"type"`
			} `json:"username"`
			Picture struct {
				Type string `json:"type"`
			} `json:"picture"`
			LinkedInID struct {
				Type string `json:"type"`
			} `json:"linkedInId"`
			GithubID struct {
				Type string `json:"type"`
			} `json:"githubId"`
		} `json:"properties"`
		Required []string `json:"required"`
	} `json:"$schema"`
	Version  string    `json:"version"`
	Type     string    `json:"type"`
	Created  time.Time `json:"created"`
	ID       string    `json:"id"`
	SourceID struct {
		ClientID    string `json:"client_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"source_id"`
	Data UserData `json:"data,omitempty"`
}
