// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

	"github.com/aws/aws-lambda-go/events"
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

func (s *service) ProjectAddedEvent(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "ProjectAddedEvent",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}
	log.WithFields(f).Debug("ProjectAddedEvent called")
	var newProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.NewImage, &newProject)
	if err != nil {
		log.WithFields(f).Warnf("project decoding add event, error: %+v", err)
		return err
	}

	f["projectSFID"] = newProject.ProjectSFID
	f["claGroupID"] = newProject.ClaGroupID
	f["foundationSFID"] = newProject.FoundationSFID

	psc := v2ProjectService.GetClient()
	log.WithFields(f).Debug("enabling CLA service...")
	err = psc.EnableCLA(newProject.ProjectSFID)
	if err != nil {
		log.WithFields(f).Warnf("enabling CLA service failed, error: %+v", err)
		return err
	}

	return nil
}

// ProjectDeleteEvent handles the CLA Group (projects table) delete event
func (s *service) ProjectDeletedEvent(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "ProjectDeletedEvent",
		"eventID":      event.EventID,
		"eventName":    event.EventName,
		"eventSource":  event.EventSource,
	}
	log.WithFields(f).Debug("ProjectDeletedEvent called")
	var oldProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.OldImage, &oldProject)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling stream image, error: %+v", err)
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
		log.WithFields(f).Warnf("disabling CLA service failed, error: %+v", err)
		return err
	}
	log.WithFields(f).Debugf("disabling CLA service took %s", time.Since(before).String())

	// Log the event
	eventErr := s.eventsRepo.CreateEvent(&models.Event{
		ContainsPII:            false,
		EventData:              fmt.Sprintf("disabled CLA service for Project: %s", oldProject.ProjectSFID),
		EventSummary:           fmt.Sprintf("disabled CLA service for Project: %s", oldProject.ProjectSFID),
		EventFoundationSFID:    oldProject.FoundationSFID,
		EventProjectExternalID: oldProject.ProjectSFID,
		EventProjectID:         oldProject.ClaGroupID,
		EventProjectSFID:       oldProject.ProjectSFID,
		EventType:              "disable.cla",
		LfUsername:             "easycla system",
		UserID:                 "easycla system",
		UserName:               "easycla system",
		// EventProjectName:       "",
		// EventProjectSFName:     "",
	})
	if eventErr != nil {
		log.WithFields(f).Warnf("problem logging event for disabling CLA service, error: %+v", eventErr)
		// Ok - don't fail for now
	}

	// remove any CLA related permissions
	permErr := s.removeCLAPermissions(oldProject.ProjectSFID)
	if permErr != nil {
		log.WithFields(f).Warnf("problem removing CLA permissions for projectSFID, error: %+v", permErr)
		// Ok - don't fail for now
	}

	return nil
}

// ProjectDeleteEvent handles the CLA Group (projects table) delete event
func (s *service) removeCLAPermissions(projectSFID string) error {
	f := logrus.Fields{
		"functionName": "removeCLAPermissions",
		"projectSFID":  projectSFID,
	}
	log.WithFields(f).Debug("removing CLA permissions...")

	roleErr := s.removeCLAPermissionsByProjectRole(projectSFID, utils.CLAManagerRole)
	if roleErr != nil {
		return roleErr
	}
	roleErr = s.removeCLAPermissionsByProjectRole(projectSFID, utils.CLADesigneeRole)
	if roleErr != nil {
		return roleErr
	}
	roleErr = s.removeCLAPermissionsByProjectRole(projectSFID, utils.CLASignatoryRole)
	if roleErr != nil {
		return roleErr
	}

	return nil
}

// removeCLAPermissionsByProjectRole handles removal of the specified role for the given SF Project
func (s *service) removeCLAPermissionsByProjectRole(projectSFID, roleName string) error {
	f := logrus.Fields{
		"functionName": "removeCLAPermissionsByProjectRole",
		"projectSFID":  projectSFID,
		"roleName":     roleName,
	}
	log.WithFields(f).Debugf("removing CLA permissions for %s...", roleName)
	client := acs_service.GetClient()
	//roleID, roleLookupErr := client.GetRoleID(roleName)
	_, roleLookupErr := client.GetRoleID(roleName)
	if roleLookupErr != nil {
		log.WithFields(f).Warnf("problem looking up role ID for %s, error: %+v", roleName, roleLookupErr)
		return roleLookupErr
	}

	// TODO: figure out how to query ACS for the list of role assignments matching the provided projectSFID and roleID

	return nil
}

// removeCLAPermissionsByProjectOrganizationRole handles removal of the specified role for the given SF Project and SF Organization
func (s *service) removeCLAPermissionsByProjectOrganizationRole(projectSFID, organizationSFID, roleName string) error {
	f := logrus.Fields{
		"functionName":     "removeCLAPermissionsByProjectOrganizationRole",
		"projectSFID":      projectSFID,
		"organizationSFID": organizationSFID,
		"roleName":         roleName,
	}
	log.WithFields(f).Debug("removing CLA permissions...")
	client := acs_service.GetClient()

	log.WithFields(f).Debugf("locating users with assigned role of %s with matching scope...", roleName)
	assignedRoles, assignedRolesErr := client.GetAssignedRoles(roleName, projectSFID, organizationSFID)
	if assignedRolesErr != nil {
		log.WithFields(f).Warnf("problem looking up assigned roles of %s, error: %+v", roleName, assignedRolesErr)
		return assignedRolesErr
	}

	var eg errgroup.Group
	for _, assignedRoleData := range assignedRoles.Data {
		// Grab the user name for the printout/log
		username := "<not defined>"
		if assignedRoleData.User != nil && assignedRoleData.User.Username != "" {
			username = assignedRoleData.User.Username
		}

		for _, assignedRole := range assignedRoleData.Roles {
			// Only delete the roles which match the provided role name - the above query returns ALL roles
			// that match this object (e.g. project|organization) and
			// the scope of: {projectSFID}|{organizationSFID}
			// Still need to filter on the name, e.g. cla-manager, cla-manager-designee or cla-signatory
			if roleName == assignedRole.Name {
				eg.Go(func() error {
					log.WithFields(f).Debugf("deleting role using role ID: %s with name: %s for user: %s",
						assignedRole.ID, assignedRole.Name, username)
					deleteErr := client.DeleteRoleByID(assignedRole.ID)
					if deleteErr != nil {
						log.WithFields(f).Warnf("problem deleting  assigned role of %s, error: %+v", roleName, deleteErr)
						return deleteErr
					}

					log.WithFields(f).Debugf("deleted role using role ID: %s with name: %s for user: %s",
						assignedRole.ID, assignedRole.Name, username)
					return nil
				})
			}
		}
	}

	log.WithFields(f).Debug("waiting for role cleanup...")
	if loadErr := eg.Wait(); loadErr != nil {
		return loadErr
	}
	return nil
}
