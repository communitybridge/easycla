// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"fmt"
	"time"

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
	log.Debug("ProjectAddedEvent called")
	var newProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.NewImage, &newProject)
	if err != nil {
		return err
	}
	psc := v2ProjectService.GetClient()
	log.WithField("project_sfid", newProject.ProjectSFID).Debug("enabling CLA service")
	err = psc.EnableCLA(newProject.ProjectSFID)
	if err != nil {
		log.WithField("project_sfid", newProject.ProjectSFID).Error("enabling CLA service failed")
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
	log.Debug("ProjectDeletedEvent called")
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
	log.WithFields(f).Debugf("disabling CLA service took %s", time.Since(before).String())

	if err != nil {
		log.WithFields(f).Warnf("disabling CLA service failed, error: %+v", err)
		return err
	}

	// Log the event
	eventErr := s.eventsRepo.CreateEvent(&models.Event{
		ContainsPII:            false,
		EventData:              fmt.Sprintf("disabled CLA service for Project: %s", oldProject.ProjectSFID),
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

	return nil
}
