// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"fmt"
	"strings"
	"sync"
	"time"

	organizationService "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"

	acsService "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

	"github.com/aws/aws-lambda-go/events"
	claEvents "github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	"github.com/sirupsen/logrus"
)

// ProjectClaGroup is database model for projects_cla_group table
type ProjectClaGroup struct {
	ProjectSFID       string `json:"project_sfid"`
	ClaGroupID        string `json:"cla_group_id"`
	FoundationSFID    string `json:"foundation_sfid"`
	RepositoriesCount int64  `json:"repositories_count"`
}

// ProjectServiceEnableCLAServiceHandler handles enabling the CLA Service attribute from the project service
func (s *service) ProjectServiceEnableCLAServiceHandler(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "ProjectServiceEnableCLAServiceHandler",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("processing request")
	var newProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.NewImage, &newProject)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("project decoding add event")
		return err
	}

	f["projectSFID"] = newProject.ProjectSFID
	f["claGroupID"] = newProject.ClaGroupID
	f["foundationSFID"] = newProject.FoundationSFID

	psc := v2ProjectService.GetClient()
	log.WithFields(f).Debug("enabling CLA service...")
	projectDetails, prjerr := psc.GetProject(newProject.ProjectSFID)
	if prjerr != nil {
		log.WithError(err).Warn("enable to get project details")
	}
	projectName := newProject.ProjectSFID
	if projectDetails != nil {
		projectName = projectDetails.Name
	}
	start, _ := utils.CurrentTime()
	err = psc.EnableCLA(newProject.ProjectSFID)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("enabling CLA service failed")
		return err
	}
	finish, _ := utils.CurrentTime()
	log.WithFields(f).Debugf("enabling CLA service completed - took: %s", finish.Sub(start).String())

	// Log the event
	eventErr := s.eventsRepo.CreateEvent(&models.Event{
		ContainsPII:            false,
		EventData:              fmt.Sprintf("enabled CLA service for project: %s", projectName),
		EventSummary:           fmt.Sprintf("enabled CLA service for project: %s", projectName),
		EventFoundationSFID:    newProject.FoundationSFID,
		EventProjectExternalID: newProject.ProjectSFID,
		EventProjectID:         newProject.ClaGroupID,
		EventProjectSFID:       newProject.ProjectSFID,
		EventType:              claEvents.ProjectServiceCLAEnabled,
		LfUsername:             "easycla system",
		UserID:                 "easycla system",
		UserName:               "easycla system",
		// EventProjectName:       "",
		EventProjectSFName: projectName,
	})
	if eventErr != nil {
		log.WithFields(f).WithError(eventErr).Warn("problem logging event for enabling CLA service")
		// Ok - don't fail for now
	}

	return nil
}

// ProjectServiceDisableCLAServiceHandler handles disabling/removing the CLA Service attribute from the project service
func (s *service) ProjectServiceDisableCLAServiceHandler(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "ProjectServiceDisableCLAServiceHandler",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("processing request")
	var oldProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.OldImage, &oldProject)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem unmarshalling stream image")
		return err
	}

	// Add more fields for the logger
	f["ProjectSFID"] = oldProject.ProjectSFID
	f["ClaGroupID"] = oldProject.ClaGroupID
	f["FoundationSFID"] = oldProject.FoundationSFID

	psc := v2ProjectService.GetClient()
	// Gathering metrics - grab the time before the API call
	before, _ := utils.CurrentTime()
	log.WithFields(f).Debug("disabling CLA service")
	err = psc.DisableCLA(oldProject.ProjectSFID)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("disabling CLA service failed")
		return err
	}
	log.WithFields(f).Debugf("disabling CLA service took %s", time.Since(before).String())

	// Log the event
	eventErr := s.eventsRepo.CreateEvent(&models.Event{
		ContainsPII:            false,
		EventData:              fmt.Sprintf("disabled CLA service for project: %s", oldProject.ProjectSFID),
		EventSummary:           fmt.Sprintf("disabled CLA service for project: %s", oldProject.ProjectSFID),
		EventFoundationSFID:    oldProject.FoundationSFID,
		EventProjectExternalID: oldProject.ProjectSFID,
		EventProjectID:         oldProject.ClaGroupID,
		EventProjectSFID:       oldProject.ProjectSFID,
		EventType:              claEvents.ProjectServiceCLADisabled,
		LfUsername:             "easycla system",
		UserID:                 "easycla system",
		UserName:               "easycla system",
		// EventProjectName:       "",
		// EventProjectSFName:     "",
	})
	if eventErr != nil {
		log.WithFields(f).WithError(eventErr).Warn("problem logging event for disabling CLA service")
		// Ok - don't fail for now
	}

	return nil
}

func (s *service) ProjectUnenrolledDisableRepositoryHandler(event events.DynamoDBEventRecord) error {
	ctx := utils.NewContext()
	f := logrus.Fields{
		"functionName":   "ProjectUnenrolledDisableRepositoryHandler",
		utils.XREQUESTID: ctx.Value(utils.XREQUESTID),
		"eventID":        event.EventID,
		"eventName":      event.EventName,
		"eventSource":    event.EventSource,
	}

	log.WithFields(f).Debug("processing request")
	var oldProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.OldImage, &oldProject)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem unmarshalling stream image")
		return err
	}

	// Add more fields for the logger
	f["ProjectSFID"] = oldProject.ProjectSFID
	f["ClaGroupID"] = oldProject.ClaGroupID
	f["FoundationSFID"] = oldProject.FoundationSFID

	// Disable GitHub repos associated with this project
	enabled := true // only care about enabled repos
	gitHubRepos, githubRepoErr := s.repositoryService.GetRepositoryByProjectSFID(ctx, oldProject.ProjectSFID, &enabled)
	if githubRepoErr != nil {
		log.WithFields(f).WithError(githubRepoErr).Warn("problem listing github repositories by project sfid")
		return githubRepoErr
	}
	if gitHubRepos != nil && len(gitHubRepos.List) > 0 {
		log.WithFields(f).Debugf("discovered %d github repositories for project with sfid: %s - disabling repositories...",
			len(gitHubRepos.List), oldProject.ProjectSFID)

		// For each GitHub repository...
		for _, gitHubRepo := range gitHubRepos.List {
			log.WithFields(f).Debugf("disabling github repository: %s with id: %s for project with sfid: %s",
				gitHubRepo.RepositoryName, gitHubRepo.RepositoryID, gitHubRepo.ProjectSFID)
			disableErr := s.repositoryService.DisableRepository(ctx, gitHubRepo.RepositoryID)
			if disableErr != nil {
				log.WithFields(f).WithError(disableErr).Warnf("problem disabling github repository: %s with id: %s", gitHubRepo.RepositoryName, gitHubRepo.RepositoryID)
				return disableErr
			}
		}
	} else {
		log.WithFields(f).Debugf("no github repositories for project with sfid: %s - nothing to disable",
			oldProject.ProjectSFID)
	}

	gerrits, gerritRepoErr := s.gerritService.GetGerritsByProjectSFID(ctx, oldProject.ProjectSFID)
	if gerritRepoErr != nil {
		log.WithFields(f).WithError(gerritRepoErr).Warn("problem listing gerrit repositories by project sfid")
		return gerritRepoErr
	}
	if gerrits != nil && len(gerrits.List) > 0 {
		log.WithFields(f).Debugf("discovered %d gerrit repositories for project with sfid: %s - deleting gerrit instances...",
			len(gerrits.List), oldProject.ProjectSFID)
		for _, gerritRepo := range gerrits.List {
			log.WithFields(f).Debugf("deleting gerrit instance: %s with id: %s for project with sfid: %s",
				gerritRepo.GerritName, gerritRepo.GerritID.String(), gerritRepo.ProjectSFID)
			gerritDeleteErr := s.gerritService.DeleteGerrit(ctx, gerritRepo.GerritID.String())
			if gerritDeleteErr != nil {
				log.WithFields(f).WithError(gerritDeleteErr).Warnf("problem deleting gerrit instance: %s with id: %s",
					gerritRepo.GerritName, gerritRepo.GerritID.String())
				return gerritDeleteErr
			}
		}
	} else {
		log.WithFields(f).Debugf("no gerrit instances for project with sfid: %s - nothing to delete",
			oldProject.ProjectSFID)
	}

	return nil
}

// AddCLAPermissions handles adding CLA permissions
func (s *service) AddCLAPermissions(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "AddCLAPermissions",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("processing event")
	var newProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.NewImage, &newProject)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem unmarshalling stream image")
		return err
	}

	// Add more fields for the logger
	f["ProjectSFID"] = newProject.ProjectSFID
	f["ClaGroupID"] = newProject.ClaGroupID
	f["FoundationSFID"] = newProject.FoundationSFID

	// Add any relevant CLA Manager permissions for this CLA Group/Project SFID
	permErr := s.addPermissions(newProject.ClaGroupID, newProject.FoundationSFID, newProject.ProjectSFID, utils.CLAManagerRole)
	if permErr != nil {
		log.WithFields(f).WithError(permErr).Warn("problem adding CLA Manager permissions for projectSFID")
		// Ok - don't fail for now
	}

	// Add any relevant CLA Manager Designee permissions for this CLA Group/Project SFID
	permErr = s.addPermissions(newProject.ClaGroupID, newProject.FoundationSFID, newProject.ProjectSFID, utils.CLADesigneeRole)
	if permErr != nil {
		log.WithFields(f).WithError(permErr).Warn("problem adding CLA Manager Designee permissions for projectSFID")
		// Ok - don't fail for now
	}

	return nil
}

// RemoveCLAPermissions handles removing existing CLA permissions
func (s *service) RemoveCLAPermissions(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "RemoveCLAPermissions",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}

	log.WithFields(f).Debug("processing event")
	var oldProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.OldImage, &oldProject)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem unmarshalling stream image")
		return err
	}

	// Add more fields for the logger
	f["ProjectSFID"] = oldProject.ProjectSFID
	f["ClaGroupID"] = oldProject.ClaGroupID
	f["FoundationSFID"] = oldProject.FoundationSFID

	// Remove any CLA related permissions
	permErr := s.removeCLAPermissions(oldProject.ProjectSFID)
	if permErr != nil {
		log.WithFields(f).WithError(permErr).Warn("problem removing CLA permissions for projectSFID")
		// Ok - don't fail for now
	}

	return nil
}

func (s *service) addPermissions(claGroupID, foundationSFID, projectSFID, role string) error {
	ctx := utils.NewContext()
	f := logrus.Fields{
		"functionName": "addPermissions",
		"claGroupID":   claGroupID,
		"projectSFID":  projectSFID,
		"role":         role,
	}

	//handle userscopes per project(users with Designee role)
	userScopes := make([]acsService.UserScope, 0)

	log.WithFields(f).Debugf("adding %s permissions...", role)
	// Check if signed at Foundation
	signedAtFoundationLevel, signedErr := s.projectService.SignedAtFoundationLevel(ctx, foundationSFID)
	if signedErr != nil {
		log.WithFields(f).Warnf("Problem getting level of CLA Group Signature for CLAGroup: %s ", claGroupID)
		return signedErr
	}
	orgClient := organizationService.GetClient()
	acsClient := acsService.GetClient()

	log.WithFields(f).Debugf("locating role ID for role: %s", role)
	roleID, roleErr := acsClient.GetRoleID(role)
	if roleErr != nil {
		log.WithFields(f).Warnf("problem looking up details for role: %s, error: %+v", utils.CLADesigneeRole, roleErr)
		return roleErr
	}

	if signedAtFoundationLevel {
		// Determine if any users have the CLA Manager Designee Role at the Foundation Level
		log.WithFields(f).Debugf("Getting users with role: %s for foundationSFID: %s  ", role, foundationSFID)
		foundationUserScopes, err := acsClient.GetProjectRoleUsersScopes(foundationSFID, role)
		if err != nil {
			log.WithFields(f).Warnf("problem getting userscopes for foundationSFID: %s and role: %s ", foundationSFID, role)
			return err
		}
		//Tabulating userscopes for new ProjectSFID assignment
		userScopes = append(userScopes, foundationUserScopes...)
		log.WithFields(f).Debugf("Found userscopes: %+v for foundationSFID: %s ", userScopes, foundationSFID)

	} else {
		// Signed at Project level Use case
		pcgs, err := s.projectsClaGroupRepo.GetProjectsIdsForClaGroup(claGroupID)
		if err != nil {
			log.WithFields(f).WithError(err).Warnf("problem getting project cla Groups for claGroupID: %s", claGroupID)
			return err
		}
		for _, pcg := range pcgs {
			//ignore newly added project
			if pcg.ProjectSFID == projectSFID {
				continue
			}
			projectUserScopes, err := acsClient.GetProjectRoleUsersScopes(pcg.ProjectSFID, role)
			if err != nil {
				log.WithFields(f).Warnf("problem getting userscopes for projectSFID: %s ", pcg.ProjectSFID)
				return err
			}
			//Tabulating userscopes for new ProjectSFID assignment
			userScopes = append(userScopes, projectUserScopes...)
			log.WithFields(f).Debugf("Found userscopes : %+v for project: %s ", userScopes, pcg.ProjectSFID)
		}

	}

	if len(userScopes) > 0 {
		log.WithFields(f).Debugf("Identified users: %+v to be updated with role: %s for project : %s ", userScopes, role, projectSFID)
		// If so, for each user, add the CLA Manager Designee role for this projectSFID
		var wg sync.WaitGroup
		wg.Add(len(userScopes))

		for _, userScope := range userScopes {
			go func(userScope acsService.UserScope) {
				defer wg.Done()

				orgID := strings.Split(userScope.ObjectID, "|")[1]
				email := userScope.Email

				roleErr := orgClient.CreateOrgUserRoleOrgScopeProjectOrg(ctx, email, projectSFID, orgID, roleID)
				if roleErr != nil {
					log.WithFields(f).WithError(roleErr).Warnf("%s, role assignment for user %s failed for this project: %s, company: %s ",
						utils.CLADesigneeRole, email, projectSFID, orgID)
					return
				}
			}(userScope)
		}

		// Wait for the goroutines to finish
		log.WithFields(f).Debugf("waiting for role:%s  assignment to complete for project: %s ", role, projectSFID)
		wg.Wait()
	}

	return nil
}

// removeCLAPermissions handles removing CLA Group (projects table) permissions for the specified project
func (s *service) removeCLAPermissions(projectSFID string) error {
	f := logrus.Fields{
		"functionName": "removeCLAPermissions",
		"projectSFID":  projectSFID,
	}
	log.WithFields(f).Debug("removing CLA permissions...")

	client := acsService.GetClient()
	err := client.RemoveCLAUserRolesByProject(projectSFID, []string{utils.CLAManagerRole, utils.CLADesigneeRole, utils.CLASignatoryRole})
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem removing CLA user roles by projectSFID")
	}

	return err
}

// removeCLAPermissionsByProjectOrganizationRole handles removal of the specified role for the given SF Project and SF Organization
func (s *service) removeCLAPermissionsByProjectOrganizationRole(projectSFID, organizationSFID string, roleNames []string) error {
	f := logrus.Fields{
		"functionName":     "removeCLAPermissionsByProjectOrganizationRole",
		"projectSFID":      projectSFID,
		"organizationSFID": organizationSFID,
		"roleNames":        strings.Join(roleNames, ","),
	}

	log.WithFields(f).Debug("removing CLA permissions...")
	client := acsService.GetClient()
	err := client.RemoveCLAUserRolesByProjectOrganization(projectSFID, organizationSFID, roleNames)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem removing CLA user roles by projectSFID and organizationSFID")
	}

	return err
}
