// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"fmt"
	"strings"
	"time"

	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

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
	log.WithFields(f).Debug("ProjectServiceEnableCLAServiceHandler called")
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

	// Log the event
	eventErr := s.eventsRepo.CreateEvent(&models.Event{
		ContainsPII:            false,
		EventData:              fmt.Sprintf("enabled CLA service for project: %s", newProject.ProjectSFID),
		EventSummary:           fmt.Sprintf("enabled CLA service for project: %s", newProject.ProjectSFID),
		EventFoundationSFID:    newProject.FoundationSFID,
		EventProjectExternalID: newProject.ProjectSFID,
		EventProjectID:         newProject.ClaGroupID,
		EventProjectSFID:       newProject.ProjectSFID,
		EventType:              claEvents.ProjectServiceCLAEnabled,
		LfUsername:             "easycla system",
		UserID:                 "easycla system",
		UserName:               "easycla system",
		// EventProjectName:       "",
		// EventProjectSFName:     "",
	})
	if eventErr != nil {
		log.WithFields(f).Warnf("problem logging event for enabling CLA service, error: %+v", eventErr)
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
	log.WithFields(f).Debug("ProjectServiceDisableCLAServiceHandler called")
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
		log.WithFields(f).Warnf("problem logging event for disabling CLA service, error: %+v", eventErr)
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
	log.WithFields(f).Debug("RemoveCLAPermissions called")
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

	// Remove any CLA related permissions
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

	client := acs_service.GetClient()
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
	client := acs_service.GetClient()
	err := client.RemoveCLAUserRolesByProjectOrganization(projectSFID, organizationSFID, roleNames)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("problem removing CLA user roles by projectSFID and organizationSFID")
	}

	return err
}
