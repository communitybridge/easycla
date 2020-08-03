// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

import (
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service handles gerrit Repository service
type Service interface {
	DeleteClaGroupGerrits(claGroupID string) (int, error)
	DeleteGerrit(gerritID string) error
	GetGerrit(gerritID string) (*models.Gerrit, error)
	AddGerrit(claGroupID string, projectSFID string, input *models.AddGerritInput) (*models.Gerrit, error)
	GetClaGroupGerrits(claGroupID string, projectSFID *string) (*models.GerritList, error)
}

type service struct {
	repo    Repository
	lfGroup *LFGroup
}

// NewService creates a new gerrit service
func NewService(repo Repository, lfg *LFGroup) Service {
	return service{
		repo:    repo,
		lfGroup: lfg,
	}
}

func (s service) DeleteClaGroupGerrits(claGroupID string) (int, error) {
	gerrits, err := s.repo.GetClaGroupGerrits(claGroupID, nil)
	if err != nil {
		return 0, err
	}
	if len(gerrits.List) > 0 {
		log.Debugf(fmt.Sprintf("Deleting gerrits for cla-group :%s ", claGroupID))
		for _, gerrit := range gerrits.List {
			err = s.repo.DeleteGerrit(gerrit.GerritID)
			if err != nil {
				return 0, err
			}
		}
	}
	return len(gerrits.List), nil
}

func (s service) DeleteGerrit(gerritID string) error {
	return s.repo.DeleteGerrit(gerritID)
}

func (s service) GetGerrit(gerritID string) (*models.Gerrit, error) {
	return s.repo.GetGerrit(gerritID)
}

func (s service) AddGerrit(claGroupID string, projectSFID string, params *models.AddGerritInput) (*models.Gerrit, error) {
	if params.GroupIDIcla == "" && params.GroupIDCcla == "" {
		return nil, errors.New("should specify at least a LDAP group for ICLA or CCLA")
	}
	if params.GerritName == nil {
		return nil, errors.New("gerrit_name required")
	}
	if params.GerritURL == nil {
		return nil, errors.New("gerrit_url required")
	}
	var groupNameCcla, groupNameIcla string
	if params.GroupIDIcla != "" {
		group, err := s.lfGroup.GetGroup(params.GroupIDIcla)
		if err != nil {
			log.WithError(err).Warnf("unable to get LDAP ICLA Group: %s", params.GroupIDIcla)
			return nil, err
		}
		groupNameIcla = group.Title
	}
	if params.GroupIDCcla != "" {
		group, err := s.lfGroup.GetGroup(params.GroupIDCcla)
		if err != nil {
			log.WithError(err).Warnf("unable to get LDAP CCLA Group: %s", params.GroupIDCcla)
			return nil, err
		}
		groupNameCcla = group.Title
	}
	input := &models.Gerrit{
		GerritName:    utils.StringValue(params.GerritName),
		GerritURL:     utils.StringValue(params.GerritURL),
		GroupIDCcla:   params.GroupIDCcla,
		GroupIDIcla:   params.GroupIDIcla,
		GroupNameCcla: groupNameCcla,
		GroupNameIcla: groupNameIcla,
		ProjectID:     claGroupID,
		ProjectSFID:   projectSFID,
	}
	return s.repo.AddGerrit(input)
}

func (s service) GetClaGroupGerrits(claGroupID string, projectSFID *string) (*models.GerritList, error) {
	return s.repo.GetClaGroupGerrits(claGroupID, projectSFID)
}
