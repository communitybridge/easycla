package dynamo_events

import (
	"github.com/aws/aws-lambda-go/events"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
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

func (s *service) ProjectDeletedEvent(event events.DynamoDBEventRecord) error {
	log.Debug("ProjectDeletedEvent called")
	var oldProject ProjectClaGroup
	err := unmarshalStreamImage(event.Change.OldImage, &oldProject)
	if err != nil {
		return err
	}
	psc := v2ProjectService.GetClient()
	log.WithField("project_sfid", oldProject.ProjectSFID).Debug("disabling CLA service")
	err = psc.DisableCLA(oldProject.ProjectSFID)
	if err != nil {
		log.WithField("project_sfid", oldProject.ProjectSFID).Error("disabling CLA service failed")
		return err
	}
	return nil
}
