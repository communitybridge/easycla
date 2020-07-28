// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/jinzhu/copier"
)

func v1ProjectModel(in *models.Project) (*v1Models.Project, error) {
	out := &v1Models.Project{}
	err := copier.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func v2ProjectModel(in *v1Models.Project) (*models.Project, error) {
	out := &models.Project{}
	err := copier.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}
