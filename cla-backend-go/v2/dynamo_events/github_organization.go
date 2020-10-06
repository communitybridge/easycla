// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	"github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// GitHubOrgUpdatedEvent github repository updated event
func (s *service) GitHubOrgUpdatedEvent(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "GitHubOrgUpdatedEvent",
	}

	log.WithFields(f).Debug("GitHubOrgUpdatedEvent called")
	var newGitHubOrg, oldGitHubOrg github_organizations.GithubOrganization
	err := unmarshalStreamImage(event.Change.NewImage, &newGitHubOrg)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the new github organization model from the updated event, error: %+v", err)
		return err
	}
	err = unmarshalStreamImage(event.Change.OldImage, &oldGitHubOrg)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling the old github organization model from the updated event, error: %+v", err)
		return err
	}

	// If the branch protection value was updated from false to true....
	if !oldGitHubOrg.BranchProtectionEnabled && newGitHubOrg.BranchProtectionEnabled {
		log.WithFields(f).Debug("transition of branchProtectionEnabled false => true - processing...")

		// Locate the repositories already saved under this organization
		log.WithFields(f).Debug("loading repositories under the organization...")
		repos, err := s.repositoryService.GetRepositoriesByOrganizationName(context.Background(), newGitHubOrg.OrganizationName)
		if err != nil {
			log.WithFields(f).Warnf("problem locating repositories by organization name, error: %+v", err)
			return err
		}

		ctx := context.Background()
		log.WithFields(f).Debugf("creating a new GitHub client object for org: %s...", newGitHubOrg.OrganizationName)
		gitHubClient, clientErr := github.NewGithubAppClient(newGitHubOrg.OrganizationInstallationID)
		if clientErr != nil {
			return clientErr
		}

		var eg errgroup.Group
		for _, repo := range repos {
			// this is for goroutine local variables
			repo := repo
			// Update the branch protection in a go routine...
			eg.Go(func() error {
				log.WithFields(f).Debugf("enabling branch protection for repository: %s", repo.RepositoryName)

				log.WithFields(f).Debugf("looking up the default branch for the GitHub repository: %s...", repo.RepositoryName)
				defaultBranch, branchErr := github.GetDefaultBranchForRepo(ctx, gitHubClient.Repositories, newGitHubOrg.OrganizationName, repo.RepositoryName)
				if branchErr != nil {
					return branchErr
				}

				log.WithFields(f).Debugf("enabling branch protection on the default branch %s for the GitHub repository: %s...",
					defaultBranch, repo.RepositoryName)
				return github.EnableBranchProtection(ctx, gitHubClient.Repositories, newGitHubOrg.OrganizationName, repo.RepositoryName,
					defaultBranch, true, []string{utils.GitHubBotName}, []string{})
			})
		}

		// Wait for the go routines to finish
		log.WithFields(f).Debugf("waiting for %d repositories to complete...", len(repos))
		var branchProtectionErr error
		if loadErr := eg.Wait(); loadErr != nil {
			log.WithFields(f).Warnf("encountered branch protection setup error: %+v", loadErr)
			branchProtectionErr = loadErr
		}

		return branchProtectionErr
	}

	log.WithFields(f).Debug("no transition of branchProtectionEnabled false => true - ignoring...")
	return nil
}
