// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/users"
)

type service struct {
	signatureRepo        signatures.SignatureRepository
	projectRepo          ProjectRepo
	userRepo             users.UserRepository
	companyRepo          company.IRepository
	repo                 IRepository
	projectClaGroupsRepo projects_cla_groups.Repository
}

type claGroupModel struct {
	ProjectName  string
	ProjectLogo  string
	ProjectSFID  string
	ProjectType  string
	SubProjects  []string
	ClaGroupName string

	// For processing
	FoundationSFID string
	SubProjectIDs  []string
}
