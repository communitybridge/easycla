// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/github_organizations"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"
	"github.com/go-openapi/swag"
	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
)

var (
	// ErrAutoEnabledOff indicates the flag is disabled on github org
	ErrAutoEnabledOff = errors.New("autoEnabled is off")
	// ErrCantDetermineAutoEnableClaGroup indicates the cla group can't be determined for github org
	ErrCantDetermineAutoEnableClaGroup = errors.New("can't determine autoEnable cla-group")
)

// AutoEnableService holds logic about handling autoEnabled field for github Org and Repos
type AutoEnableService interface {
	CreateAutoEnabledRepository(repo *github.Repository) (*models.GithubRepository, error)
	AutoEnabledForGithubOrg(f logrus.Fields, gitHubOrg github_organizations.GithubOrganization) error
}

// NewAutoEnableService creates a new AutoEnableService
func NewAutoEnableService(repositoryService repositories.Service,
	githubRepo repositories.Repository,
	githubOrgRepo github_organizations.Repository,
	claRepository projects_cla_groups.Repository) AutoEnableService {
	return &autoEnableServiceProvider{
		repositoryService: repositoryService,
		githubRepo:        githubRepo,
		githubOrgRepo:     githubOrgRepo,
		claRepository:     claRepository,
	}
}

// autoEnableServiceProvider is an abstraction helping with managing autoEnabled flag for Github Organization
// having it separated in its own struct makes testing easier.
type autoEnableServiceProvider struct {
	repositoryService repositories.Service
	githubRepo        repositories.Repository
	githubOrgRepo     github_organizations.Repository
	claRepository     projects_cla_groups.Repository
}

func (a *autoEnableServiceProvider) CreateAutoEnabledRepository(repo *github.Repository) (*models.GithubRepository, error) {
	repositoryFullName := *repo.FullName
	repositoryExternalID := strconv.FormatInt(*repo.ID, 10)

	f := logrus.Fields{
		"functionName":       "handleRepositoryAddedAction",
		"repositoryFullName": repositoryFullName,
	}

	organizationName := strings.Split(repositoryFullName, "/")[0]
	ctx := context.Background()
	orgModel, err := a.githubOrgRepo.GetGithubOrganization(ctx, organizationName)
	if err != nil {
		log.Warnf("fetching github org failed : %v", err)
		return nil, err
	}

	if !orgModel.AutoEnabled {
		log.Warnf("skipping adding the repository, autoEnabled flag is off")
		return nil, ErrAutoEnabledOff
	}
	orgName := orgModel.OrganizationName

	claGroupID := orgModel.AutoEnabledClaGroupID
	if claGroupID == "" {
		repos, listErr := a.repositoryService.ListProjectRepositories(context.Background(), orgModel.ProjectSFID)
		if listErr != nil {
			log.WithFields(f).Warnf("problem fetching the repositories for orgName : %s for ProjectSFID : %s", orgName, orgModel.ProjectSFID)
			return nil, listErr
		}

		if len(repos.List) == 0 {
			log.WithFields(f).Warnf("no repositories found for orgName : %s, skipping autoEnabled", orgName)
			return nil, ErrCantDetermineAutoEnableClaGroup
		}

		claGroupID, listErr = DetermineClaGroupID(f, orgModel, repos)
		if listErr != nil {
			return nil, listErr
		}
	}
	claGroupModel, err := a.claRepository.GetCLAGroup(claGroupID)
	if err != nil {
		log.Warnf("fetching the cla group for cla group id : %s failed : %v", claGroupID, err)
		return nil, err
	}

	repoModel, err := a.githubRepo.AddGithubRepository(ctx, "", claGroupModel.ProjectSFID, &models.GithubRepositoryInput{
		RepositoryProjectID:        swag.String(claGroupID),
		RepositoryName:             swag.String(repositoryFullName),
		RepositoryType:             swag.String("github"),
		RepositoryURL:              swag.String("https://github.com/" + repositoryFullName),
		RepositoryOrganizationName: swag.String(organizationName),
		RepositoryExternalID:       swag.String(repositoryExternalID),
	})

	if err != nil {
		return nil, err
	}

	return repoModel, nil
}

func (a *autoEnableServiceProvider) AutoEnabledForGithubOrg(f logrus.Fields, gitHubOrg github_organizations.GithubOrganization) error {
	orgName := gitHubOrg.OrganizationName
	log.WithFields(f).Debugf("running AutoEnable for github org : %s", orgName)
	if gitHubOrg.OrganizationInstallationID == 0 {
		log.WithFields(f).Warnf("missing installation id for : %s", orgName)
		return fmt.Errorf("missing installation id")
	}

	repos, err := a.repositoryService.ListProjectRepositories(context.Background(), gitHubOrg.ProjectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem fetching the repositories for orgName : %s for ProjectSFID : %s", orgName, gitHubOrg.ProjectSFID)
		return err
	}

	if len(repos.List) == 0 {
		log.WithFields(f).Warnf("no repositories found for orgName : %s, skipping autoEnabled", orgName)
		return nil
	}

	claGroupID, err := DetermineClaGroupID(f, github_organizations.ToModel(&gitHubOrg), repos)
	if err != nil {
		return err
	}

	for _, repo := range repos.List {
		if repo.RepositoryProjectID == claGroupID {
			continue
		}

		repo.RepositoryProjectID = claGroupID
		if err := a.repositoryService.UpdateClaGroupID(context.Background(), repo.RepositoryID, claGroupID); err != nil {
			log.WithFields(f).Warnf("updating claGroupID for repository : %s failed : %v", repo.RepositoryID, err)
			return err
		}
	}

	return nil
}

// DetermineClaGroupID checks if AutoEnabledClaGroupID is set then returns it (high precedence) otherwise tries to determine
// the autoEnabled claGroupID by guessing from existing repos
func DetermineClaGroupID(f logrus.Fields, gitHubOrg *models.GithubOrganization, repos *models.ListGithubRepositories) (string, error) {
	if gitHubOrg.AutoEnabledClaGroupID != "" {
		return gitHubOrg.AutoEnabledClaGroupID, nil
	}

	// fallback to old way of checking the cla group id by guessing it from the existing repos which has cla group id set
	claGroupSet := map[string]bool{}
	sfidSet := map[string]bool{}
	// check if any of the repos is member to more than one cla group, in general shouldn't happen
	var claGroupID string
	for _, repo := range repos.List {
		if repo.RepositoryProjectID == "" || repo.ProjectSFID == "" {
			continue
		}
		claGroupSet[repo.RepositoryProjectID] = true
		sfidSet[repo.ProjectSFID] = true
		claGroupID = repo.RepositoryProjectID
	}

	if len(claGroupSet) == 0 && len(sfidSet) == 0 {
		return "", fmt.Errorf("none of the existing repos have the clagroup set, can't determine the cla group, please set the claGroupID on githubOrg : %w", ErrCantDetermineAutoEnableClaGroup)
	}

	if len(claGroupSet) != 1 || len(sfidSet) != 1 {
		log.WithFields(f).Errorf(`Auto Enabled set for Organization %s, '
                                but we found repositories or SFIDs that belong to multiple CLA Groups.
                                Auto Enable only works when all repositories under a given
                                GitHub Organization are associated with a single CLA Group. This
                                organization is associated with %d CLA Groups and %d SFIDs.`, gitHubOrg.OrganizationName, len(claGroupSet), len(sfidSet))
		return "", fmt.Errorf("project and its repos should be part of the same cla group, can't determine main cla group, please set the claGroupID on githubOrg : %w", ErrCantDetermineAutoEnableClaGroup)
	}

	return claGroupID, nil
}
