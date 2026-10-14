// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

import (
	"context"
	"testing"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/golang/mock/gomock"
	mock_company "github.com/linuxfoundation/easycla/cla-backend-go/company/mocks"
	v1Models "github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	mock_project "github.com/linuxfoundation/easycla/cla-backend-go/project/mocks"
	"github.com/linuxfoundation/easycla/cla-backend-go/projects_cla_groups"
	mock_projects_cla_groups "github.com/linuxfoundation/easycla/cla-backend-go/projects_cla_groups/mocks"
	"github.com/linuxfoundation/easycla/cla-backend-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func projectScopedUser(scopeType, id string) *auth.User {
	return &auth.User{UserName: "scoped-user", Email: "scoped@example.com", ACL: auth.ACL{Allowed: true, Scopes: []auth.Scope{{Type: scopeType, ID: id}}}}
}

// a cla manager of company A on project P must not pass the project branch for company B (lfx-self-serve-ops#382)
func TestProjectScopeHelpersIgnoreProjectOrganizationScopes(t *testing.T) {
	ctx := context.Background()
	managerA := projectScopedUser(auth.ProjectOrganization, "proj-sfid|company-a")
	managerB := projectScopedUser(auth.ProjectOrganization, "proj-sfid|company-b")
	projectManager := projectScopedUser(auth.Project, "proj-sfid")
	mapping := &projects_cla_groups.ProjectClaGroup{ProjectSFID: "proj-sfid", ProjectExternalID: "proj-sfid", FoundationSFID: "foundation-sfid", ClaGroupID: "cla-group-1", ClaGroupName: "CLA Group 1"}

	t.Run("isUserHaveAccessToCLAProjectOrganization", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_projects_cla_groups.NewMockRepository(ctrl)
		repo.EXPECT().GetClaGroupIDForProject(gomock.Any(), "proj-sfid").Return(mapping, nil).AnyTimes()
		repo.EXPECT().GetProjectsIdsForClaGroup(gomock.Any(), "cla-group-1").Return([]*projects_cla_groups.ProjectClaGroup{mapping}, nil).AnyTimes()
		assert.False(t, isUserHaveAccessToCLAProjectOrganization(ctx, managerA, "proj-sfid", "company-b", repo))
		assert.True(t, isUserHaveAccessToCLAProjectOrganization(ctx, managerB, "proj-sfid", "company-b", repo))
		assert.True(t, isUserHaveAccessToCLAProjectOrganization(ctx, projectManager, "proj-sfid", "company-b", repo))
		// a foundation-level CLA manager of A (ACS lists the child projects in Related) reaches the child project only for A
		foundationManagerA := &auth.User{UserName: "foundation-manager", ACL: auth.ACL{Allowed: true, Scopes: []auth.Scope{{Type: auth.ProjectOrganization, ID: "foundation-sfid|company-a", Related: []string{"proj-sfid"}}}}}
		assert.False(t, isUserHaveAccessToCLAProjectOrganization(ctx, foundationManagerA, "proj-sfid", "company-b", repo))
		assert.True(t, isUserHaveAccessToCLAProjectOrganization(ctx, foundationManagerA, "proj-sfid", "company-a", repo))
	})

	t.Run("isUserHaveAccessToCLAGroupProjects", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_projects_cla_groups.NewMockRepository(ctrl)
		repo.EXPECT().GetProjectsIdsForClaGroup(gomock.Any(), "cla-group-1").Return([]*projects_cla_groups.ProjectClaGroup{mapping}, nil).AnyTimes()
		projectRepo := mock_project.NewMockProjectRepository(ctrl)
		assert.False(t, isUserHaveAccessToCLAGroupProjects(ctx, managerA, "cla-group-1", repo, projectRepo))
		assert.True(t, isUserHaveAccessToCLAGroupProjects(ctx, projectManager, "cla-group-1", repo, projectRepo))
	})

	t.Run("isUserHaveAccessOfSignedSignaturePDF", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_projects_cla_groups.NewMockRepository(ctrl)
		repo.EXPECT().GetProjectsIdsForClaGroup(gomock.Any(), "cla-group-1").Return([]*projects_cla_groups.ProjectClaGroup{mapping}, nil).AnyTimes()
		projectRepo := mock_project.NewMockProjectRepository(ctrl)
		companyService := mock_company.NewMockIService(ctrl)
		companyService.EXPECT().GetCompany(gomock.Any(), "company-b-id").Return(&v1Models.Company{CompanyID: "company-b-id", CompanyExternalID: "company-b"}, nil).AnyTimes()
		ccla := &v1Models.Signature{SignatureID: "sig-1", ProjectID: "cla-group-1", SignatureType: utils.SignatureTypeCCLA, SignatureReferenceType: utils.SignatureReferenceTypeCompany, SignatureReferenceID: "company-b-id"}
		ok, err := isUserHaveAccessOfSignedSignaturePDF(ctx, managerA, ccla, companyService, repo, projectRepo)
		require.NoError(t, err)
		assert.False(t, ok)
		ok, err = isUserHaveAccessOfSignedSignaturePDF(ctx, managerB, ccla, companyService, repo, projectRepo)
		require.NoError(t, err)
		assert.True(t, ok)
		ok, err = isUserHaveAccessOfSignedSignaturePDF(ctx, projectManager, ccla, companyService, repo, projectRepo)
		require.NoError(t, err)
		assert.True(t, ok)
		foundationManagerA := &auth.User{UserName: "foundation-manager", ACL: auth.ACL{Allowed: true, Scopes: []auth.Scope{{Type: auth.ProjectOrganization, ID: "foundation-sfid|company-a", Related: []string{"proj-sfid"}}}}}
		ok, err = isUserHaveAccessOfSignedSignaturePDF(ctx, foundationManagerA, ccla, companyService, repo, projectRepo)
		require.NoError(t, err)
		assert.False(t, ok, "a foundation-level CLA manager of another company must not read company B's signed CCLA")
	})
}
