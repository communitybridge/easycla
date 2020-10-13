// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"context"

	"github.com/communitybridge/easycla/cla-backend-go/github"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
)

// GithubRepoAddedEvent github repository added event
func (s *service) GithubRepoAddedEvent(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "GithubRepoAddedEvent",
	}

	log.WithFields(f).Debug("GithubRepoAddedEvent called")
	var newRepoModel repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.NewImage, &newRepoModel)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the new repository model event, error: %+v", err)
		return err
	}

	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(newRepoModel.ProjectSFID)
	if err != nil {
		return err
	}

	var uerr error
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		log.Debugf("incrementing root_project_repositories_count of cla_group_id %s", newRepoModel.RepositoryProjectID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(context.Background(), newRepoModel.RepositoryProjectID, 1)
	} else {
		log.Debugf("incrementing repositories_count for project %s", newRepoModel.ProjectSFID)
		uerr = s.projectsClaGroupRepo.UpdateRepositoriesCount(newRepoModel.ProjectSFID, 1)
	}
	if uerr != nil {
		return err
	}

	return nil
}

// GithubRepoDeletedEvent github repository removed event
func (s *service) GithubRepoDeletedEvent(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "GithubRepoDeletedEvent",
	}

	log.WithFields(f).Debug("GithubRepoDeletedEvent called...")
	var oldRepoModel repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.OldImage, &oldRepoModel)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the old repository model event, error: %+v", err)
		return err
	}

	psc := v2ProjectService.GetClient()
	project, err := psc.GetProject(oldRepoModel.ProjectSFID)
	if err != nil {
		return err
	}

	var uerr error
	if project.Parent == "" || project.Parent == utils.TheLinuxFoundation {
		log.Debugf("decrementing root_project_repositories_count of cla_group_id %s", oldRepoModel.RepositoryProjectID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(context.Background(), oldRepoModel.RepositoryProjectID, -1)
	} else {
		log.Debugf("decrementing repositories_count for project %s", oldRepoModel.ProjectSFID)
		uerr = s.projectsClaGroupRepo.UpdateRepositoriesCount(oldRepoModel.ProjectSFID, -1)
	}
	if uerr != nil {
		return err
	}

	return nil
}

// EnableBranchProtectionServiceHandler handles enabling the CLA Service attribute from the project service
func (s *service) EnableBranchProtectionServiceHandler(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "EnableBranchProtectionServiceHandler",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("EnableBranchProtectionServiceHandler called...")
	var newRepoModel repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.NewImage, &newRepoModel)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the new github organization add event, error: %+v", err)
		return err
	}

	f["repositoryName"] = newRepoModel.RepositoryName
	f["repositoryOrganizationName"] = newRepoModel.RepositoryOrganizationName
	f["projectSFID"] = newRepoModel.ProjectSFID

	// Branch protection only available for GitHub
	if newRepoModel.RepositoryType == utils.GitHubType {
		log.WithFields(f).Debugf("repository type is: %s", utils.GitHubType)

		parentOrgName := newRepoModel.RepositoryOrganizationName
		log.WithFields(f).Warnf("problem locating github organization by name: %s, error: %+v", parentOrgName, err)
		gitHubOrg, err := s.githubOrgService.GetGithubOrganizationByName(context.Background(), parentOrgName)
		if err != nil {
			log.WithFields(f).Warnf("problem locating github organization by name: %s, error: %+v", parentOrgName, err)
			return nil
		}
		if gitHubOrg == nil {
			log.WithFields(f).Warnf("problem locating github organization by name: %s - record not found", parentOrgName)
			return nil
		}

		if gitHubOrg.BranchProtectionEnabled {
			log.WithFields(f).Debug("branch protection is enabled for this organization")

			ctx := context.Background()
			log.WithFields(f).Debug("creating a new GitHub client object...")
			gitHubClient, clientErr := github.NewGithubAppClient(gitHubOrg.OrganizationInstallationID)
			if clientErr != nil {
				return clientErr
			}

			branchProtectionRepository := github.NewBranchProtectionRepository(gitHubClient.Repositories, github.EnableBlockingLimiter())

			log.WithFields(f).Debug("looking up the default branch for the GitHub repository...")
			defaultBranch, branchErr := branchProtectionRepository.GetDefaultBranchForRepo(ctx, gitHubOrg.OrganizationName, newRepoModel.RepositoryName)
			if branchErr != nil {
				return branchErr
			}

			log.WithFields(f).Debugf("enabling branch protection on th default branch %s for the GitHub repository: %s...",
				defaultBranch, newRepoModel.RepositoryName)
			return branchProtectionRepository.EnableBranchProtection(ctx,
				parentOrgName, newRepoModel.RepositoryName,
				defaultBranch, true, []string{utils.GitHubBotName}, []string{})
		}

		log.WithFields(f).Debug("github organization branch protection is not enabled - no action required")
	}

	return nil
}

// DisableBranchProtectionServiceHandler handles disabling/removing the CLA Service attribute from the project service
func (s *service) DisableBranchProtectionServiceHandler(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "DisableBranchProtectionServiceHandler",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("DisableBranchProtectionServiceHandler called")
	var oldRepoModel repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.OldImage, &oldRepoModel)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the old github organization removed event, error: %+v", err)
		return err
	}

	// Branch protection only available for GitHub
	if oldRepoModel.RepositoryType == utils.GitHubType {
		parentOrgName := oldRepoModel.RepositoryOrganizationName
		gitHubOrg, err := s.githubOrgService.GetGithubOrganizationByName(context.Background(), parentOrgName)
		if err != nil {
			log.WithFields(f).Warnf("problem locating github organization by name: %s, error: %+v", parentOrgName, err)
			return nil
		}
		if gitHubOrg == nil {
			log.WithFields(f).Warnf("problem locating github organization by name: %s - record not found", parentOrgName)
			return nil
		}

		if gitHubOrg.BranchProtectionEnabled {
			log.Debug("github organization branch protection is enabled - no cleanup action required")
			return nil
		}
		log.Debug("github organization branch protection is not enabled - no action required")
	}

	return nil
}
