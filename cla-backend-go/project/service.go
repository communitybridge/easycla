// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"context"
	"sync"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	"github.com/communitybridge/easycla/cla-backend-go/repositories"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations/project"
	"github.com/communitybridge/easycla/cla-backend-go/gerrits"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service interface defines the project service methods/functions
type Service interface {
	CreateCLAGroup(ctx context.Context, project *models.Project) (*models.Project, error)
	GetCLAGroups(ctx context.Context, params *project.GetProjectsParams) (*models.Projects, error)
	GetCLAGroupByID(ctx context.Context, projectID string) (*models.Project, error)
	GetCLAGroupsByExternalSFID(ctx context.Context, projectSFID string) (*models.Projects, error)
	GetCLAGroupsByExternalID(ctx context.Context, params *project.GetProjectsByExternalIDParams) (*models.Projects, error)
	GetCLAGroupByName(ctx context.Context, projectName string) (*models.Project, error)
	DeleteCLAGroup(ctx context.Context, projectID string) error
	UpdateCLAGroup(ctx context.Context, projectModel *models.Project) (*models.Project, error)
	GetClaGroupsByFoundationSFID(ctx context.Context, foundationSFID string, loadRepoDetails bool) (*models.Projects, error)
	GetClaGroupByProjectSFID(ctx context.Context, projectSFID string, loadRepoDetails bool) (*models.Project, error)
	SignedAtFoundationLevel(ctx context.Context, foundationSFID string) (bool, error)
}

// service
type service struct {
	repo             ProjectRepository
	repositoriesRepo repositories.Repository
	gerritRepo       gerrits.Repository
	projectCGRepo    projects_cla_groups.Repository
}

// NewService returns an instance of the project service
func NewService(projectRepo ProjectRepository, repositoriesRepo repositories.Repository, gerritRepo gerrits.Repository, pcgRepo projects_cla_groups.Repository) Service {
	return service{
		repo:             projectRepo,
		repositoriesRepo: repositoriesRepo,
		gerritRepo:       gerritRepo,
		projectCGRepo:    pcgRepo,
	}
}

// CreateProject service method
func (s service) CreateCLAGroup(ctx context.Context, projectModel *models.Project) (*models.Project, error) {
	return s.repo.CreateCLAGroup(ctx, projectModel)
}

// GetCLAGroups service method
func (s service) GetCLAGroups(ctx context.Context, params *project.GetProjectsParams) (*models.Projects, error) {
	return s.repo.GetCLAGroups(ctx, params)
}

// GetProjectByID service method
func (s service) GetCLAGroupByID(ctx context.Context, projectID string) (*models.Project, error) {
	f := logrus.Fields{
		"functionName":    "GetCLAGroupByID",
		utils.XREQUESTID:  ctx.Value(utils.XREQUESTID),
		"projectID":       projectID,
		"loadRepoDetails": LoadRepoDetails,
	}

	log.WithFields(f).Debug("locating CLA Group by ID...")
	project, err := s.repo.GetCLAGroupByID(ctx, projectID, LoadRepoDetails)
	if err != nil {
		return nil, err
	}

	// No Foundation SFID value? Maybe this is a v1 CLA Group record...
	if project.FoundationSFID == "" {
		log.WithFields(f).Debug("CLA Group missing FoundationSFID...")
		// Most likely this is a CLA Group v1 record - use the external ID if available
		if project.ProjectExternalID != "" {
			log.WithFields(f).Debugf("CLA Group assigning foundationID to value of external ID: %s", project.ProjectExternalID)
			project.FoundationSFID = project.ProjectExternalID
		}
	}

	if project.FoundationSFID != "" {
		signed, checkErr := s.SignedAtFoundationLevel(ctx, project.FoundationSFID)
		if checkErr != nil {
			return nil, checkErr
		}
		project.FoundationLevelCLA = signed
	}

	return project, nil
}

// GetCLAGroupsByExternalSFID returns a list of projects based on the external SFID parameter
func (s service) GetCLAGroupsByExternalSFID(ctx context.Context, projectSFID string) (*models.Projects, error) {
	return s.GetCLAGroupsByExternalID(ctx, &project.GetProjectsByExternalIDParams{
		HTTPRequest: nil,
		NextKey:     nil,
		PageSize:    nil,
		ProjectSFID: projectSFID,
	})
}

// GetCLAGroupsByExternalID returns a list of projects based on the external ID parameters
func (s service) GetCLAGroupsByExternalID(ctx context.Context, params *project.GetProjectsByExternalIDParams) (*models.Projects, error) {
	f := logrus.Fields{
		"functionName":   "GetCLAGroupsByExternalID",
		utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
		"projectSFID":    params.ProjectSFID,
		"NextKey":        params.NextKey,
		"PageSize":       params.PageSize}
	log.Debugf("Project Service Handler - GetCLAGroupsByExternalID")
	projects, err := s.repo.GetCLAGroupsByExternalID(ctx, params, LoadRepoDetails)
	if err != nil {
		log.WithFields(f).Warnf("problem with query, error: %+v", err)
		return nil, err
	}
	numberOfProjects := len(projects.Projects)
	if numberOfProjects == 0 {
		return projects, nil
	}

	// Add repository information in the response model
	var wg sync.WaitGroup
	wg.Add(numberOfProjects)
	for i := range projects.Projects {
		go func(project *models.Project) {
			defer wg.Done()
			s.fillRepoInfo(ctx, project)
		}(&projects.Projects[i])
	}
	wg.Wait()

	return projects, nil
}

func (s service) fillRepoInfo(ctx context.Context, project *models.Project) {
	var wg sync.WaitGroup
	wg.Add(2)
	var ghrepos []*models.GithubRepositoriesGroupByOrgs
	var gerrits []*models.Gerrit
	go func() {
		defer wg.Done()
		var err error
		ghrepos, err = s.repositoriesRepo.GetCLAGroupRepositoriesGroupByOrgs(ctx, project.ProjectID, true)
		if err != nil {
			log.Error("unable to get github repositories for project.", err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		var gerritsList *models.GerritList
		gerritsList, err = s.gerritRepo.GetClaGroupGerrits(project.ProjectID, nil)
		if err != nil {
			log.Error("unable to get gerrit instances:w for project.", err)
			return
		}
		gerrits = gerritsList.List
	}()
	wg.Wait()
	project.GithubRepositories = ghrepos
	project.Gerrits = gerrits
}

// GetCLAGroupByName service method
func (s service) GetCLAGroupByName(ctx context.Context, projectName string) (*models.Project, error) {
	return s.repo.GetCLAGroupByName(ctx, projectName)
}

// DeleteCLAGroup service method
func (s service) DeleteCLAGroup(ctx context.Context, projectID string) error {
	return s.repo.DeleteCLAGroup(ctx, projectID)
}

// UpdateCLAGroup service method
func (s service) UpdateCLAGroup(ctx context.Context, projectModel *models.Project) (*models.Project, error) {
	// Updates to the CLA Group "projects" table will cause a DB trigger handler (separate lambda) to also update other
	// tables where we have the CLA Group name/description
	return s.repo.UpdateCLAGroup(ctx, projectModel)
}

// GetClaGroupsByFoundationSFID service method
func (s service) GetClaGroupsByFoundationSFID(ctx context.Context, foundationSFID string, loadRepoDetails bool) (*models.Projects, error) {
	return s.repo.GetClaGroupsByFoundationSFID(ctx, foundationSFID, loadRepoDetails)
}

// GetClaGroupByProjectSFID( service method
func (s service) GetClaGroupByProjectSFID(ctx context.Context, projectSFID string, loadRepoDetails bool) (*models.Project, error) {
	return s.repo.GetClaGroupByProjectSFID(ctx, projectSFID, loadRepoDetails)
}

// SignedAtFoundationLevel returns true if the specified foundation has a CLA Group at the foundation level, returns false otherwise.
func (s service) SignedAtFoundationLevel(ctx context.Context, foundationSFID string) (bool, error) {
	f := logrus.Fields{
		"functionName":   "SignedAtFoundationLevel",
		utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
		"foundationSFID": foundationSFID,
	}

	log.WithFields(f).Debug("querying foundation CLA Group entries...")
	entries, pcgErr := s.projectCGRepo.GetProjectsIdsForFoundation(foundationSFID)
	if pcgErr != nil {
		return false, pcgErr
	}
	log.WithFields(f).Debugf("loaded %d CLA Group entries", len(entries))

	// Check for number of claGroups for foundation
	foundationLevelCLAGroup := false
	for _, entry := range entries {
		if entry.ProjectSFID == entry.FoundationSFID {
			foundationLevelCLAGroup = true
			break
		}
	}

	return foundationLevelCLAGroup, nil
}
