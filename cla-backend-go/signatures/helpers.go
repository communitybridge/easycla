// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"github.com/communitybridge/easycla/cla-backend-go/gen/v1/models"
	"github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/go-openapi/strfmt"
)

// getAddEmailContributors is a helper function to lookup the contributors impacted by the Approval List update
func (s service) getAddEmailContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.AddEmailApprovalList {
		userModel, err := s.usersService.GetUserByEmail(value)
		if err != nil {
			logging.Warnf("unable to lookup user by LF email: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// getRemoveEmailContributors is a helper function to lookup the contributors impacted by the Approval List update
func (s service) getRemoveEmailContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.RemoveEmailApprovalList {
		userModel, err := s.usersService.GetUserByEmail(value)
		if err != nil {
			logging.Warnf("unable to lookup user by LF email: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// getAddGitHubContributors is a helper function to lookup the contributors impacted by the Approval List update
func (s service) getAddGitHubContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.AddGithubUsernameApprovalList {
		userModel, err := s.usersService.GetUserByGitHubUsername(value)
		if err != nil {
			logging.Warnf("unable to lookup user by GitHub username: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// getRemoveGitHubContributors is a helper function to lookup the contributors impacted by the Approval List update
func (s service) getRemoveGitHubContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.RemoveGithubUsernameApprovalList {
		userModel, err := s.usersService.GetUserByGitHubUsername(value)
		if err != nil {
			logging.Warnf("unable to lookup user by GitHub username: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// getAddGitlabContributors is a helper function to look up the Gitlab contributors impacted by the Approval List update
func (s service) getAddGitlabContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.AddGitlabUsernameApprovalList {
		userModel, err := s.usersService.GetUserByGitHubUsername(value)
		if err != nil {
			logging.Warnf("unable to lookup user by Gitlab username: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// getRemoveGitlabContributors is a helper function to look up the Gitlab contributors impacted by the Approval List update
func (s service) getRemoveGitlabContributors(approvalList *models.ApprovalList) []*models.User {
	var userModelList []*models.User
	for _, value := range approvalList.RemoveGitlabUsernameApprovalList {
		userModel, err := s.usersService.GetUserByGitHubUsername(value)
		if err != nil {
			logging.Warnf("unable to lookup user by Gitlab username: %s, error: %+v", value, err)
		} else {
			userModelList = append(userModelList, userModel)
		}
	}

	return userModelList
}

// createUserModelFromEmail is a helper function to create a user model in our local database based on the specified email only
func (s service) createUserModelFromEmail(email string) (*models.User, error) {
	return s.usersService.CreateUser(&models.User{
		Admin:          false,
		CompanyID:      "",
		DateCreated:    "",
		DateModified:   "",
		Emails:         []string{email},
		GithubID:       "",
		GithubUsername: "",
		GitlabID:       "",
		GitlabUsername: "",
		LfEmail:        strfmt.Email(email),
		LfUsername:     "",
		Note:           "",
		UserExternalID: "",
		UserID:         "",
		Username:       "",
		Version:        "v1",
	}, &user.CLAUser{
		UserID:         "",
		Name:           "",
		Emails:         nil,
		LFEmail:        "",
		LFUsername:     "",
		LfidProvider:   user.Provider{},
		GithubProvider: user.Provider{},
		ProjectIDs:     nil,
		ClaIDs:         nil,
		CompanyIDs:     nil,
	})
}

// createUserModelFromGitHubUsername is a helper function to create a user model in our local database based on the specified GitHub username
func (s service) createUserModelFromGitHubUsername(gitHubUsername, gitHubUserID, email string) (*models.User, error) {
	return s.usersService.CreateUser(&models.User{
		Admin:          false,
		CompanyID:      "",
		DateCreated:    "",
		DateModified:   "",
		Emails:         []string{email},
		GithubID:       gitHubUserID,
		GithubUsername: gitHubUsername,
		GitlabID:       "",
		GitlabUsername: "",
		LfEmail:        strfmt.Email(email),
		LfUsername:     "",
		Note:           "",
		UserExternalID: "",
		UserID:         "",
		Username:       "",
		Version:        "v1",
	}, &user.CLAUser{
		UserID:         "",
		Name:           "",
		Emails:         nil,
		LFEmail:        "",
		LFUsername:     "",
		LfidProvider:   user.Provider{},
		GithubProvider: user.Provider{},
		ProjectIDs:     nil,
		ClaIDs:         nil,
		CompanyIDs:     nil,
	})
}

// createUserModelFromGitLabUsername is a helper function to create a user model in our local database based on the specified GitLab username
func (s service) createUserModelFromGitLabUsername(gitLabUsername, gitLabUserID, email string) (*models.User, error) {
	return s.usersService.CreateUser(&models.User{
		Admin:          false,
		CompanyID:      "",
		DateCreated:    "",
		DateModified:   "",
		Emails:         []string{email},
		GithubID:       "",
		GithubUsername: "",
		GitlabID:       gitLabUserID,
		GitlabUsername: gitLabUsername,
		LfEmail:        strfmt.Email(email),
		LfUsername:     "",
		Note:           "",
		UserExternalID: "",
		UserID:         "",
		Username:       "",
		Version:        "v1",
	}, &user.CLAUser{
		UserID:         "",
		Name:           "",
		Emails:         nil,
		LFEmail:        "",
		LFUsername:     "",
		LfidProvider:   user.Provider{},
		GithubProvider: user.Provider{},
		ProjectIDs:     nil,
		ClaIDs:         nil,
		CompanyIDs:     nil,
	})
}
