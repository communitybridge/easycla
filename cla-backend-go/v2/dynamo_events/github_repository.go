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

func (s *service) GithubRepoAddedEvent(event events.DynamoDBEventRecord) error {
	log.Debug("GithubRepoAddedEvent called")
	var newGithubOrg repositories.RepositoryDBModel
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
		log.Debugf("incrementing root_project_repositories_count of cla_group_id %s", newGithubOrg.RepositoryProjectID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(context.Background(), newGithubOrg.RepositoryProjectID, 1)
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
	var oldGithubOrg repositories.RepositoryDBModel
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
		log.Debugf("decrementing root_project_repositories_count of cla_group_id %s", oldGithubOrg.RepositoryProjectID)
		uerr = s.projectRepo.UpdateRootCLAGroupRepositoriesCount(context.Background(), oldGithubOrg.RepositoryProjectID, -1)
	} else {
		log.Debugf("decrementing repositories_count for project %s", oldGithubOrg.ProjectSFID)
		uerr = s.projectsClaGroupRepo.UpdateRepositoriesCount(oldGithubOrg.ProjectSFID, -1)
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
	var newGithubRepositoryModel repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.NewImage, &newGithubRepositoryModel)
	if err != nil {
		return err
	}

	f["repositoryName"] = newGithubRepositoryModel.RepositoryName
	f["repositoryOrganizationName"] = newGithubRepositoryModel.RepositoryOrganizationName
	f["projectSFID"] = newGithubRepositoryModel.ProjectSFID

	// Branch protection only available for GitHub
	if newGithubRepositoryModel.RepositoryType == utils.GitHubType {
		log.WithFields(f).Debugf("repository type is: %s", utils.GitHubType)

		parentOrgName := newGithubRepositoryModel.RepositoryOrganizationName
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

			log.WithFields(f).Debug("looking up the default branch for the GitHub repository...")
			defaultBranch, branchErr := github.GetDefaultBranchForRepo(ctx, gitHubClient, gitHubOrg.OrganizationName, newGithubRepositoryModel.RepositoryName)
			if branchErr != nil {
				return branchErr
			}

			log.WithFields(f).Debugf("enabling branch protection on th default branch %s for the GitHub repository: %s...",
				defaultBranch, newGithubRepositoryModel.RepositoryName)
			return github.EnableBranchProtection(ctx, gitHubClient,
				parentOrgName, newGithubRepositoryModel.RepositoryName,
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
	var oldGithubOrg repositories.RepositoryDBModel
	err := unmarshalStreamImage(event.Change.OldImage, &oldGithubOrg)
	if err != nil {
		return err
	}

	// Branch protection only available for GitHub
	if oldGithubOrg.RepositoryType == utils.GitHubType {
		parentOrgName := oldGithubOrg.RepositoryOrganizationName
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
