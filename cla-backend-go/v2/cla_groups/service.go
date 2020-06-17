// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_groups

import (
	"context"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/jinzhu/copier"

	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Project "github.com/communitybridge/easycla/cla-backend-go/project"
	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"
	v1Template "github.com/communitybridge/easycla/cla-backend-go/template"
	v2ProjectService "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	"github.com/sirupsen/logrus"
)

type service struct {
	v1ProjectService      v1Project.Service
	v1TemplateService     v1Template.Service
	projectsClaGroupsRepo projects_cla_groups.Repository
}

// Service interface
type Service interface {
	CreateCLAGroup(input *models.CreateClaGroupInput, projectManagerLFID string) (*models.ClaGroup, error)
	ValidateCLAGroup(input *models.ClaGroupValidationRequest) (bool, []string)
}

// NewService returns instance of CLA group service
func NewService(projectService v1Project.Service, templateService v1Template.Service, projectsClaGroupsRepo projects_cla_groups.Repository) Service {
	return &service{
		v1ProjectService:      projectService, // aka cla_group service of v1
		v1TemplateService:     templateService,
		projectsClaGroupsRepo: projectsClaGroupsRepo,
	}
}

// CreateCLAGroup is the service handler for creating a new CLA Group
func (s *service) CreateCLAGroup(input *models.CreateClaGroupInput, projectManagerLFID string) (*models.ClaGroup, error) {
	f := logrus.Fields{"function": "CreateCLAGroup"}
	// Validate the input
	log.WithFields(f).WithField("input", input).Debugf("validating create cla group input")
	err := s.validateClaGroupInput(input)
	if err != nil {
		log.WithFields(f).Warnf("validation of create cla group input failed")
		return nil, err
	}

	// Create cla group
	log.WithFields(f).WithField("input", input).Debugf("creating cla group")
	claGroup, err := s.v1ProjectService.CreateProject(&v1Models.Project{
		FoundationSFID:          input.FoundationSfid,
		ProjectDescription:      input.ClaGroupDescription,
		ProjectCCLAEnabled:      input.CclaEnabled,
		ProjectCCLARequiresICLA: input.CclaRequiresIcla,
		ProjectExternalID:       input.FoundationSfid,
		ProjectACL:              []string{projectManagerLFID},
		ProjectICLAEnabled:      input.IclaEnabled,
		ProjectName:             input.ClaGroupName,
		Version:                 "v2",
	})
	if err != nil {
		log.WithFields(f).Errorf("creating cla group failed. error = %s", err.Error())
		return nil, err
	}
	log.WithFields(f).WithField("cla_group", claGroup).Debugf("cla group created")
	f["cla_group_id"] = claGroup.ProjectID

	// Attach template with cla group
	var templateFields v1Models.CreateClaGroupTemplate
	err = copier.Copy(&templateFields, &input.TemplateFields)
	if err != nil {
		log.WithFields(f).Error("unable to create v1 create cla group template model", err)
		return nil, err
	}
	log.WithFields(f).Debug("attaching cla_group_template")
	if templateFields.TemplateID == "" {
		log.WithFields(f).Debug("using apache style template as template_id is not passed")
		templateFields.TemplateID = v1Template.ApacheStyleTemplateID
	}
	pdfUrls, err := s.v1TemplateService.CreateCLAGroupTemplate(context.Background(), claGroup.ProjectID, &templateFields)
	if err != nil {
		log.WithFields(f).Error("attaching cla_group_template failed", err)
		log.WithFields(f).Debug("deleting created cla group")
		deleteErr := s.v1ProjectService.DeleteProject(claGroup.ProjectID)
		if deleteErr != nil {
			log.WithFields(f).Error("deleting created cla group failed.", deleteErr)
		}
		return nil, err
	}
	log.WithFields(f).Debug("cla_group_template attached", pdfUrls)

	// Associate projects with cla group
	err = s.enrollProjects(claGroup.ProjectID, input.FoundationSfid, input.ProjectSfidList)
	if err != nil {
		log.WithFields(f).Debug("deleting created cla group")
		deleteErr := s.v1ProjectService.DeleteProject(claGroup.ProjectID)
		if deleteErr != nil {
			log.WithFields(f).Error("deleting created cla group failed.", deleteErr)
		}
		return nil, err
	}

	return &models.ClaGroup{
		CclaEnabled:         claGroup.ProjectCCLAEnabled,
		CclaPdfURL:          pdfUrls.CorporatePDFURL,
		CclaRequiresIcla:    claGroup.ProjectCCLARequiresICLA,
		ClaGroupDescription: claGroup.ProjectDescription,
		ClaGroupID:          claGroup.ProjectID,
		ClaGroupName:        claGroup.ProjectName,
		FoundationSfid:      claGroup.FoundationSFID,
		IclaEnabled:         claGroup.ProjectICLAEnabled,
		IclaPdfURL:          pdfUrls.IndividualPDFURL,
		ProjectSfidList:     input.ProjectSfidList,
	}, nil
}

// ValidateCLAGroup is the service handler for validating a CLA Group
func (s *service) ValidateCLAGroup(input *models.ClaGroupValidationRequest) (bool, []string) {

	var valid = true
	var validationErrors []string

	if input.ClaGroupName != nil {
		claGroupModel, err := s.v1ProjectService.GetProjectByName(*input.ClaGroupName)
		if err != nil {
			valid = false
			validationErrors = append(validationErrors, fmt.Sprintf("unable to query project service - error: %+v", err))
		}
		if claGroupModel != nil {
			valid = false
			validationErrors = append(validationErrors, fmt.Sprintf("CLA Group with name %s already exist", *input.ClaGroupName))
		}

		if len(*input.ClaGroupName) < 3 {
			valid = false
			validationErrors = append(validationErrors, "CLA Group name should be at least 3 characters")
		}
		if len(*input.ClaGroupName) > 256 {
			valid = false
			validationErrors = append(validationErrors, "description maximum length of the CLA Group name is 256 characters")
		}
	}

	if input.ClaGroupDescription != nil {
		if len(*input.ClaGroupDescription) < 3 {
			valid = false
			validationErrors = append(validationErrors, "description should be at least 3 characters")
		}
		if len(*input.ClaGroupDescription) > 256 {
			valid = false
			validationErrors = append(validationErrors, "description maximum length of the description is 256 characters")
		}
	}

	return valid, validationErrors
}

func (s *service) validateClaGroupInput(input *models.CreateClaGroupInput) error {
	if !input.IclaEnabled && !input.CclaEnabled {
		return fmt.Errorf("bad request: can not create cla group with both icla and ccla disabled")
	}
	if input.CclaRequiresIcla {
		if !(input.IclaEnabled && input.CclaEnabled) {
			return fmt.Errorf("bad request: ccla_requires_icla can not be enabled if one of icla/ccla is disabled")
		}
	}
	claGroupModel, err := s.v1ProjectService.GetProjectByName(input.ClaGroupName)
	if err != nil {
		return err
	}
	if claGroupModel != nil {
		return fmt.Errorf("bad request: cla_group with name %s already exist", input.ClaGroupName)
	}
	err = s.validateEnrollProjectsInput(input.FoundationSfid, input.ProjectSfidList)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) validateEnrollProjectsInput(foundationSFID string, projectSFIDList []string) error {
	psc := v2ProjectService.GetClient()

	// fetch foundation and its sub projects
	foundationProjectDetails, err := psc.GetProject(foundationSFID)
	if err != nil {
		return err
	}

	// check if it is foundation
	if foundationProjectDetails.ProjectType != "Foundation" {
		return fmt.Errorf("bad request: invalid foundation_sfid: %s", foundationSFID)
	}

	// check if all enrolled projects are part of foundation
	foundationProjectList := utils.NewStringSet()
	for _, pr := range foundationProjectDetails.Projects {
		foundationProjectList.Add(pr.ID)
	}
	for _, projectSFID := range projectSFIDList {
		if !foundationProjectList.Include(projectSFID) {
			return fmt.Errorf("bad request: invalid project_sfid: %s. This project is not under foundation", projectSFID)
		}
	}

	// check if projects are not already enabled
	enabledProjects, err := s.projectsClaGroupsRepo.GetProjectsIdsForFoundation(foundationSFID)
	if err != nil {
		return err
	}
	enabledProjectList := utils.NewStringSet()
	for _, pr := range enabledProjects {
		enabledProjectList.Add(pr.ProjectSFID)
	}
	for _, projectSFID := range projectSFIDList {
		if enabledProjectList.Include(projectSFID) {
			return fmt.Errorf("bad request: invalid project_sfid passed : %s. This project is already part of another cla_group", projectSFID)
		}
	}
	return nil
}

func (s *service) enrollProjects(claGroupID string, foundationSFID string, projectSFIDList []string) error {
	f := logrus.Fields{"function": "enrollProjects"}
	for _, projectSFID := range projectSFIDList {
		log.WithFields(f).Debugf("associating cla_group with project : %s", projectSFID)
		err := s.projectsClaGroupsRepo.AssociateClaGroupWithProject(claGroupID, projectSFID, foundationSFID)
		if err != nil {
			log.WithFields(f).Errorf("associating cla_group with project : %s failed", projectSFID)
			log.WithFields(f).Debug("deleting stale entries from cla_group project association")
			deleteErr := s.projectsClaGroupsRepo.RemoveProjectAssociatedWithClaGroup(claGroupID, projectSFIDList, false)
			if deleteErr != nil {
				log.WithFields(f).Error("deleting stale entries from cla_group project association failed", deleteErr)
			}
			return err
		}
	}
	return nil
}
