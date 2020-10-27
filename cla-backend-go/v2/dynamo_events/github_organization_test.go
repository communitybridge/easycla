// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"fmt"
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/golang/mock/gomock"
)

func Test_AutoEnableService(t *testing.T) {

	externalProjectID := "sfd12343"
	claGroupID := "da04291f-75d1-4e84-8275-9bc008205837"

	testCases := []struct {
		name              string
		githubOrg         github_organizations.GithubOrganization
		repositoryService func(m *repositories.MockService)
		errStr            string
	}{
		{
			name:      "fail missing installation id",
			githubOrg: github_organizations.GithubOrganization{},
			errStr:    "missing installation id",
		},
		{
			name: "fail repos fetching error",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(nil, fmt.Errorf("fetch error"))
			},
			errStr: "fetch error",
		},
		{
			name: "fail no repos",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{}, nil)
			},
		},
		{
			name: "fail no cla group found",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{
						List: []*models.GithubRepository{
							{
								RepositoryID: "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4",
								ProjectSFID:  externalProjectID,
							},
							{
								RepositoryID: "b42216b4-8f6d-41c0-8cde-7b2acbf0656a",
								ProjectSFID:  externalProjectID,
							},
						},
					}, nil)
			},
			errStr: "can't determine the cla group",
		},
		{
			name: "fail multiple cla groups found",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{
						List: []*models.GithubRepository{
							{
								RepositoryID:        "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4",
								RepositoryProjectID: claGroupID,
								ProjectSFID:         externalProjectID,
							},
							{
								RepositoryID:        "b42216b4-8f6d-41c0-8cde-7b2acbf0656a",
								RepositoryProjectID: "anotherclagroup",
								ProjectSFID:         externalProjectID,
							},
						},
					}, nil)
			},
			errStr: "can't determine main cla group",
		},
		{
			name: "success cla set",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
				AutoEnabledClaGroupID:      claGroupID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{
						List: []*models.GithubRepository{
							{
								RepositoryID: "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4",
								ProjectSFID:  externalProjectID,
							},
							{
								RepositoryID: "b42216b4-8f6d-41c0-8cde-7b2acbf0656a",
								ProjectSFID:  externalProjectID,
							},
						},
					}, nil)
				m.
					EXPECT().
					UpdateClaGroupID(gomock.Any(), "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4", claGroupID).
					Return(nil)
				m.
					EXPECT().
					UpdateClaGroupID(gomock.Any(), "b42216b4-8f6d-41c0-8cde-7b2acbf0656a", claGroupID).
					Return(nil)
			},
		},
		{
			name: "success no cla set, no call to update",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{
						List: []*models.GithubRepository{
							{
								RepositoryID:        "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4",
								ProjectSFID:         externalProjectID,
								RepositoryProjectID: claGroupID,
							},
							{
								RepositoryID:        "b42216b4-8f6d-41c0-8cde-7b2acbf0656a",
								ProjectSFID:         externalProjectID,
								RepositoryProjectID: claGroupID,
							},
						},
					}, nil)
			},
		},
		{
			name: "success no cla set determine from repos update",
			githubOrg: github_organizations.GithubOrganization{
				OrganizationInstallationID: 12354,
				ProjectSFID:                externalProjectID,
			},
			repositoryService: func(m *repositories.MockService) {
				m.
					EXPECT().
					ListProjectRepositories(gomock.Any(), externalProjectID).
					Return(&models.ListGithubRepositories{
						List: []*models.GithubRepository{
							{
								RepositoryID: "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4",
								ProjectSFID:  externalProjectID,
							},
							{
								RepositoryID:        "b42216b4-8f6d-41c0-8cde-7b2acbf0656a",
								ProjectSFID:         externalProjectID,
								RepositoryProjectID: claGroupID,
							},
						},
					}, nil)

				m.
					EXPECT().
					UpdateClaGroupID(gomock.Any(), "d7c1050b-2f32-44ea-bad2-3c8ff980ccd4", claGroupID).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			m := repositories.NewMockService(ctrl)
			if tc.repositoryService != nil {
				tc.repositoryService(m)
			}

			a := &AutoEnableService{repositoryService: m}
			err := a.autoEnabledForGithubOrg(logrus.Fields{
				"functionName": "TestAutoEnable",
			}, tc.githubOrg)

			if tc.errStr == "" {
				assert.NoError(tt, err)
			} else {
				assert.Contains(tt, err.Error(), tc.errStr)
			}
		})
	}
}
