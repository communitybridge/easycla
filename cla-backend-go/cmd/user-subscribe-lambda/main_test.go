// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	mock_users "github.com/linuxfoundation/easycla/cla-backend-go/users/mocks"
	userServiceModels "github.com/linuxfoundation/easycla/cla-backend-go/v2/user-service/models"
)

func serviceUser() *userServiceModels.User {
	return &userServiceModels.User{
		ID:        "0032M00000AbCdEQAU",
		Username:  "new.user",
		FirstName: "New",
		LastName:  "User",
		Emails: []*userServiceModels.Email{
			{EmailAddress: aws.String("new.user@example.org"), IsPrimary: aws.Bool(true)},
			{EmailAddress: aws.String("new.user@example.com"), IsPrimary: aws.Bool(false)},
		},
	}
}

// user-service e-mail entries are pointers with an omitempty IsPrimary: entries without an address or without the
// IsPrimary flag must not take the Lambda down either
func TestCreateUserTolerantEmailEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_users.NewMockUserRepository(ctrl)
	repo.EXPECT().GetUserByLFUserName("new.user").Return(nil, nil)
	var created *models.User
	repo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
		created = u
		return u, nil
	})
	user := serviceUser()
	user.Emails = []*userServiceModels.Email{
		nil,
		{EmailAddress: nil, IsPrimary: aws.Bool(true)},
		{EmailAddress: aws.String("secondary@example.com")},
		{EmailAddress: aws.String("primary@example.org"), IsPrimary: aws.Bool(true)},
	}
	createUser(logrus.Fields{}, "new.user", repo, func(string) (*userServiceModels.User, error) { return user, nil })
	require.NotNil(t, created)
	assert.Equal(t, "primary@example.org", created.LfEmail.String())
	assert.Equal(t, []string{"secondary@example.com", "primary@example.org"}, created.Emails)
}

// ops#888: a brand-new LF user has no EasyCLA record yet (GetUserByLFUserName returns nil), so the record is built from the
// user-service object alone and the user id is left to the repository, which assigns a fresh UUID
func TestCreateUserNewUser(t *testing.T) {
	for name, lookupErr := range map[string]error{"no record": nil, "lookup error": errors.New("dynamo unavailable")} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_users.NewMockUserRepository(ctrl)
			repo.EXPECT().GetUserByLFUserName("new.user").Return(nil, lookupErr)
			var created *models.User
			repo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
				created = u
				return u, nil
			})
			lookups := 0
			createUser(logrus.Fields{}, "new.user", repo, func(username string) (*userServiceModels.User, error) {
				lookups++
				assert.Equal(t, "new.user", username)
				return serviceUser(), nil
			})
			require.NotNil(t, created, "CreateUser is called once for a new user")
			assert.Equal(t, 1, lookups)
			assert.Empty(t, created.UserID, "the repository assigns the user id")
			assert.Equal(t, "new.user", created.LfUsername)
			assert.Equal(t, "new.user@example.org", created.LfEmail.String())
			assert.Equal(t, []string{"new.user@example.org", "new.user@example.com"}, created.Emails)
			assert.Equal(t, "0032M00000AbCdEQAU", created.UserExternalID)
			assert.Equal(t, "New User", created.Username)
			assert.Equal(t, "Create via user-service event", created.Note)
			assert.Equal(t, "v1", created.Version)
			assert.False(t, created.Admin)
			assert.NotEmpty(t, created.DateCreated)
			assert.Equal(t, created.DateCreated, created.DateModified)
		})
	}
}

// a UserSignedUp for a username that already has an EasyCLA record adds nothing: profile changes arrive as
// UserUpdatedProfile and are applied by Update (the strict mock fails the test on any CreateUser call)
func TestCreateUserExistingUserIsNotDuplicated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_users.NewMockUserRepository(ctrl)
	repo.EXPECT().GetUserByLFUserName("known.user").Return(&models.User{UserID: "existing-uuid", LfUsername: "known.user"}, nil)
	createUser(logrus.Fields{}, "known.user", repo, func(string) (*userServiceModels.User, error) {
		t.Errorf("the user service must not be consulted for an existing user")
		return serviceUser(), nil
	})
}

// no user-service record (or a user-service error) means there is nothing to create
func TestCreateUserWithoutUserServiceRecord(t *testing.T) {
	for name, lookup := range map[string]userLookup{
		"nil record":    func(string) (*userServiceModels.User, error) { return nil, nil },
		"service error": func(string) (*userServiceModels.User, error) { return nil, errors.New("user-service 502") },
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_users.NewMockUserRepository(ctrl)
			repo.EXPECT().GetUserByLFUserName("new.user").Return(nil, nil)
			createUser(logrus.Fields{}, "new.user", repo, lookup)
		})
	}
}
