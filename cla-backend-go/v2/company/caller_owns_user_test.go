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

// the user ID in the path must be the signed-in user's own EasyCLA user, proven by the token identity (lfx-self-serve-ops#343)
func TestCallerOwnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_users.NewMockUserRepository(ctrl)
	jane := &v1Models.User{UserID: "u-jane", LfUsername: "jane", LfEmail: strfmt.Email("jane@example.invalid"), Emails: []string{"jane.alt@example.invalid"}}
	githubOnly := &v1Models.User{UserID: "u-gh", Emails: []string{"gh@example.invalid"}}
	repo.EXPECT().GetUser("u-jane").Return(jane, nil).AnyTimes()
	repo.EXPECT().GetUser("u-gh").Return(githubOnly, nil).AnyTimes()
	repo.EXPECT().GetUser("u-ghost").Return(nil, nil)
	repo.EXPECT().GetUser("u-down").Return(nil, errors.New("dynamodb unavailable"))
	svc := &service{userRepo: repo}
	str := func(s string) *string { return &s }

	assert.NoError(t, svc.callerOwnsUser("u-jane", str("Jane"), str("other@example.invalid")), "LF username proves ownership")
	assert.NoError(t, svc.callerOwnsUser("u-jane", str("mallory"), str("JANE@example.invalid")), "token e-mail matching lf_email proves ownership")
	assert.NoError(t, svc.callerOwnsUser("u-jane", nil, str("jane.alt@example.invalid")), "token e-mail matching an alternate e-mail proves ownership")
	assert.NoError(t, svc.callerOwnsUser("u-gh", nil, str("gh@example.invalid")), "a GitHub-only row is owned through its public e-mail")
	err := svc.callerOwnsUser("u-jane", str("mallory"), str("mallory@example.invalid"))
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "another user's ID is refused")
	err = svc.callerOwnsUser("u-gh", str("mallory"), nil)
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "a GitHub-only row is not owned by an unrelated token")
	err = svc.callerOwnsUser("u-jane", nil, nil)
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "a token without username or e-mail owns nothing")
	err = svc.callerOwnsUser("u-ghost", str("jane"), str("jane@example.invalid"))
	assert.True(t, errors.Is(err, ErrCompanyUserMismatch), "an unknown user ID is refused")
	err = svc.callerOwnsUser("u-down", str("jane"), str("jane@example.invalid"))
	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrCompanyUserMismatch), "a lookup failure is not a mismatch")
}
