// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"errors"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	v1Models "github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	mock_users "github.com/linuxfoundation/easycla/cla-backend-go/users/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the user ID in the path must be the signed-in user's own EasyCLA user (lfx-self-serve-ops#343)
func TestCallerOwnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_users.NewMockUserRepository(ctrl)
	jane := &v1Models.User{UserID: "u-jane", LfUsername: "jane", LfEmail: strfmt.Email("jane@example.invalid"), Emails: []string{"jane.alt@example.invalid"}}
	repo.EXPECT().GetUser("u-jane").Return(jane, nil).AnyTimes()
	repo.EXPECT().GetUser("u-ghost").Return(nil, nil)
	repo.EXPECT().GetUser("u-down").Return(nil, errors.New("dynamodb unavailable"))
	svc := &service{userRepo: repo}
	name := func(s string) *string { return &s }

	assert.NoError(t, svc.callerOwnsUser("u-jane", name("Jane"), "other@example.invalid"), "LF username match")
	assert.NoError(t, svc.callerOwnsUser("u-jane", name("mallory"), "JANE@example.invalid"), "LF e-mail match")
	assert.NoError(t, svc.callerOwnsUser("u-jane", nil, "jane.alt@example.invalid"), "alternate e-mail match")
	err := svc.callerOwnsUser("u-jane", name("mallory"), "mallory@example.invalid")
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "another user's ID is refused")
	err = svc.callerOwnsUser("u-ghost", name("jane"), "jane@example.invalid")
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "an unknown user ID is refused")
	err = svc.callerOwnsUser("u-down", name("jane"), "jane@example.invalid")
	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrCompanyUserMismatch), "a lookup failure is not a mismatch")
}
