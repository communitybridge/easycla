// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

import (
	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/signatures"
	"github.com/communitybridge/easycla/cla-backend-go/users"
	v2ProjectServiceModels "github.com/communitybridge/easycla/cla-backend-go/v2/project-service/models"
)

type service struct {
	signatureRepo signatures.SignatureRepository
	projectRepo   ProjectRepo
	userRepo      users.UserRepository
	companyRepo   company.IRepository
	repo          IRepository
}

type signatureResponse struct {
	companyID  string
	projectID  string
	signatures *v1Models.Signatures
	err        error
}

type projectDetailedModel struct {
	v1ProjectModel       *v1Models.Project
	v2ProjectOutputModel *v2ProjectServiceModels.ProjectOutput
}
