// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

import (
	"fmt"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service handles gerrit Repository service
type Service interface {
	DeleteProject(projectID string) error
}

type service struct {
	repo Repository
}

// NewService creates a new gerrit service
func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) DeleteProject(projectID string) error {
	gerrits, err := s.repo.GetProjectGerrits(projectID)
	if err != nil {
		return err
	}
	if len(gerrits) > 0 {
		log.Debugf(fmt.Sprintf("Deleting gerrit projects for project :%s ", projectID))
		for _, gerrit := range gerrits {
			err = s.repo.DeleteProject(gerrit.GerritID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
