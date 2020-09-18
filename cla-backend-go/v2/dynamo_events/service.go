// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package dynamo_events

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/approval_list"
	"github.com/communitybridge/easycla/cla-backend-go/cla_manager"

	claevent "github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/communitybridge/easycla/cla-backend-go/company"

	"github.com/communitybridge/easycla/cla-backend-go/signatures"

	"github.com/sirupsen/logrus"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// constants
const (
	Insert = "INSERT"
	Modify = "MODIFY"
	Remove = "REMOVE"
)

// EventHandlerFunc is type for dynamoDB event handler function
type EventHandlerFunc func(event events.DynamoDBEventRecord) error

type service struct {
	// key : tablename:action
	functions                map[string][]EventHandlerFunc
	signatureRepo            signatures.SignatureRepository
	companyRepo              company.IRepository
	projectsClaGroupRepo     projects_cla_groups.Repository
	eventsRepo               claevent.Repository
	projectRepo              project.ProjectRepository
	projectService           project.Service
	claManagerRequestsRepo   cla_manager.IRepository
	approvalListRequestsRepo approval_list.IRepository
}

// Service implements DynamoDB stream event handler service
type Service interface {
	ProcessEvents(event events.DynamoDBEvent)
}

// NewService creates DynamoDB stream event handler service
func NewService(stage string,
	signatureRepo signatures.SignatureRepository,
	companyRepo company.IRepository,
	pcgRepo projects_cla_groups.Repository,
	eventsRepo claevent.Repository,
	projectRepo project.ProjectRepository,
	projService project.Service,
	claManagerRequestsRepo cla_manager.IRepository,
	approvalListRequestsRepo approval_list.IRepository) Service {
	SignaturesTable := fmt.Sprintf("cla-%s-signatures", stage)
	eventsTable := fmt.Sprintf("cla-%s-events", stage)
	projectsCLAGroupsTable := fmt.Sprintf("cla-%s-projects-cla-groups", stage)
	repositoryTableName := fmt.Sprintf("cla-%s-repositories", stage)
	claGroupsTable := fmt.Sprintf("cla-%s-projects", stage)

	s := &service{
		functions:                make(map[string][]EventHandlerFunc),
		signatureRepo:            signatureRepo,
		companyRepo:              companyRepo,
		projectsClaGroupRepo:     pcgRepo,
		eventsRepo:               eventsRepo,
		projectRepo:              projectRepo,
		projectService:           projService,
		claManagerRequestsRepo:   claManagerRequestsRepo,
		approvalListRequestsRepo: approvalListRequestsRepo,
	}

	s.registerCallback(SignaturesTable, Modify, s.SignatureSignedEvent)
	s.registerCallback(SignaturesTable, Modify, s.SignatureAddSigTypeSignedApprovedID)
	s.registerCallback(SignaturesTable, Insert, s.SignatureAddSigTypeSignedApprovedID)
	s.registerCallback(SignaturesTable, Insert, s.SignatureAddUsersDetails)

	s.registerCallback(eventsTable, Insert, s.EventAddedEvent)

	s.registerCallback(projectsCLAGroupsTable, Insert, s.ProjectAddedEvent)
	s.registerCallback(projectsCLAGroupsTable, Remove, s.ProjectDeletedEvent)

	s.registerCallback(repositoryTableName, Insert, s.GithubRepoAddedEvent)
	s.registerCallback(repositoryTableName, Remove, s.GithubRepoDeletedEvent)

	s.registerCallback(claGroupsTable, Modify, s.ProcessCLAGroupUpdateEvents)

	return s
}

func (s *service) registerCallback(tableName, eventName string, callbackFunction EventHandlerFunc) {
	key := fmt.Sprintf("%s:%s", tableName, eventName)
	funcArr := s.functions[key]
	funcArr = append(funcArr, callbackFunction)
	s.functions[key] = funcArr
}

func (s *service) ProcessEvents(events events.DynamoDBEvent) {
	for _, event := range events.Records {
		tableName := strings.Split(event.EventSourceArn, "/")[1]
		fields := logrus.Fields{
			"table_name":  tableName,
			"eventID":     event.EventID,
			"eventName":   event.EventName,
			"eventSource": event.EventSource,
			// Dumping the event is super verbose
			// "event":      event,
		}
		// Generates a ton of output
		// b, _ := json.Marshal(events) // nolint
		//fields["events_data"] = string(b)
		log.WithFields(fields).Debug("processing event record")
		key := fmt.Sprintf("%s:%s", tableName, event.EventName)
		for _, f := range s.functions[key] {
			fields["key"] = key
			fields["functionType"] = fmt.Sprintf("%T", f)
			log.WithFields(fields).Debug("invoking handler")
			err := f(event)
			if err != nil {
				log.WithFields(fields).WithField("event", event).Error("unable to process event", err)
			}
		}
	}
}

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to struct
func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)
	for k, v := range attribute {
		var dbAttr dynamodb.AttributeValue
		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}
		err := json.Unmarshal(bytes, &dbAttr)
		if err != nil {
			return err
		}
		dbAttrMap[k] = &dbAttr
	}
	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)
}
