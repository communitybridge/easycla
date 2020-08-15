// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"github.com/aws/aws-lambda-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/sirupsen/logrus"
)

func (s *service) ProcessCLAGroupUpdateEvents(event events.DynamoDBEventRecord) error {
	f := logrus.Fields{
		"functionName": "ProjectUpdatedEvent",
	}

	f["eventID"] = event.EventID
	f["eventName"] = event.EventName
	f["eventSource"] = event.EventSource
	log.WithFields(f).Debug("processing event")

	var updatedProject v1Models.Project
	err := unmarshalStreamImage(event.Change.NewImage, &updatedProject)
	if err != nil {
		log.WithFields(f).Warnf("unable to unmarshal project model, error: %+v", err)
		return err
	}
	log.WithFields(f).Debugf("decoded project record from stream: %+v", updatedProject)

	// Update any DB records that have CLA Approval Requests from Contributors - need to update Name, etc. if that has changed
	log.WithFields(f).Debugf("updating any CLA approval requests from contributors for this CLA Group")
	approvalListRequestErr := s.approvalListRequestsRepo.UpdateRequestsByCLAGroup(updatedProject)
	if approvalListRequestErr != nil {
		log.WithFields(f).Warnf("unable to update contributor approval list requests with updated CLA Group information, error: %+v", approvalListRequestErr)
	}

	log.WithFields(f).Debugf("updating any CLA manager requests for this CLA Group")
	managerRequestErr := s.claManagerRequestsRepo.UpdateRequestsByCLAGroup(updatedProject)
	if managerRequestErr != nil {
		log.WithFields(f).Warnf("unable to update cla manager request with updated CLA Group information, error: %+v", approvalListRequestErr)
	}

	// TODO - update other tables:
	//  cla-%s-metrics,
	//  cla-%s-projects-cla-groups,
	//  cla-%s-gerrit-instances,
	// possibly add/update cla_group_name/project_name to other tables:
	//  cla-%-repositories
	//  cla-%-signatures

	return nil
}
