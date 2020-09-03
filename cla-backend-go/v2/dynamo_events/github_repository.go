// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"github.com/aws/aws-lambda-go/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
)

// GithubRepository is database model for repositories table table
type GithubRepository struct {
	ProjectSFID string `json:"project_sfid"`
	ClaGroupID  string `json:"repository_project_id"`
}

func (s *service) GithubRepoAddedEvent(event events.DynamoDBEventRecord) error {
	log.Debug("GithubRepoAddedEvent called")
	var newGithubOrg GithubRepository
	err := unmarshalStreamImage(event.Change.NewImage, &newGithubOrg)
	if err != nil {
		return err
	}
	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(newGithubOrg.ProjectSFID)
	if err != nil {
		return err
	}
	var uerr error
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		log.Debugf("incrementing root_project_repositories_count of cla_group_id %s", newGithubOrg.ClaGroupID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(newGithubOrg.ClaGroupID, 1)
	} else {
		log.Debugf("incrementing repositories_count for project %s", newGithubOrg.ProjectSFID)
		uerr = s.projectsClaGroupRepo.UpdateRepositoriesCount(newGithubOrg.ProjectSFID, 1)
	}
	if uerr != nil {
		return err
	}
	return nil
}

func (s *service) GithubRepoDeletedEvent(event events.DynamoDBEventRecord) error {
	log.Debug("GithubRepoDeletedEvent called")
	var oldGithubOrg GithubRepository
	err := unmarshalStreamImage(event.Change.OldImage, &oldGithubOrg)
	if err != nil {
		return err
	}
	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(oldGithubOrg.ProjectSFID)
	if err != nil {
		return err
	}
	var uerr error
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		log.Debugf("decrementing root_project_repositories_count of cla_group_id %s", oldGithubOrg.ClaGroupID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(oldGithubOrg.ClaGroupID, -1)
	} else {
		log.Debugf("decrementing repositories_count for project %s", oldGithubOrg.ProjectSFID)
		uerr = s.projectsClaGroupRepo.UpdateRepositoriesCount(oldGithubOrg.ProjectSFID, -1)
	}
	if uerr != nil {
		return err
	}
	return nil
}
