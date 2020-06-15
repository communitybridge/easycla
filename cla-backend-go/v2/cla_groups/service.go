package cla_groups

import (
	"context"
	"fmt"
	"strconv"

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

// constants
const (
	DontLoadDetails = false
	LoadDetails     = true
)

type service struct {
	v1ProjectService      v1Project.Service
	v1TemplateService     v1Template.Service
	projectsClaGroupsRepo projects_cla_groups.Repository
}

// Service interface
type Service interface {
	CreateCLAGroup(input *models.CreateClaGroupInput, projectManagerLFID string) (*models.ClaGroup, error)
	EnrollProjectsInClaGroup(claGroupID string, foundationSFID string, projectSFIDList []string) error
	DeleteCLAGroup(claGroupID string) error
	ListClaGroupsUnderFoundation(foundationSFID string) (*models.ClaGroupList, error)
}

// NewService returns instance of CLA group service
func NewService(projectService v1Project.Service, templateService v1Template.Service, projectsClaGroupsRepo projects_cla_groups.Repository) Service {
	return &service{
		v1ProjectService:      projectService, // aka cla_group service of v1
		v1TemplateService:     templateService,
		projectsClaGroupsRepo: projectsClaGroupsRepo,
	}
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

func (s *service) EnrollProjectsInClaGroup(claGroupID string, foundationSFID string, projectSFIDList []string) error {
	f := logrus.Fields{"cla_group_id": claGroupID, "foundation_sfid": foundationSFID, "project_sfid_list": projectSFIDList}
	log.WithFields(f).Debug("validating enroll project input")
	err := s.validateEnrollProjectsInput(foundationSFID, projectSFIDList)
	if err != nil {
		log.WithFields(f).Errorf("validating enroll project input failed. error = %s", err)
		return err
	}
	log.WithFields(f).Debug("validating enroll project input passed")
	log.WithFields(f).Debug("enrolling projects in cla_group")
	err = s.enrollProjects(claGroupID, foundationSFID, projectSFIDList)
	if err != nil {
		log.WithFields(f).Errorf("enrolling projects in cla_group failed. error = %s", err)
		return err
	}
	log.WithFields(f).Debug("projects enrolled successfully in cla_group")
	return nil
}

func (s *service) DeleteCLAGroup(claGroupID string) error {
	f := logrus.Fields{"cla_group_id": claGroupID}
	log.WithFields(f).Debug("deleting cla_group")
	log.WithFields(f).Debug("deleting cla_group project association")
	err := s.projectsClaGroupsRepo.RemoveProjectAssociatedWithClaGroup(claGroupID, []string{}, true)
	if err != nil {
		return nil
	}
	log.WithFields(f).Debug("deleting cla_group from dynamodb")
	err = s.v1ProjectService.DeleteProject(claGroupID)
	if err != nil {
		log.WithFields(f).Errorf("deleting cla_group from dynamodb failed. error = %s", err.Error())
		return err
	}
	return nil
}

func getS3Url(claGroupID string, docs []v1Models.ProjectDocument) string {
	if len(docs) == 0 {
		return ""
	}
	var version int64
	var url string
	for _, doc := range docs {
		maj, err := strconv.Atoi(doc.DocumentMajorVersion)
		if err != nil {
			log.WithField("cla_group_id", claGroupID).Error("invalid major number in cla_group")
			continue
		}
		min, err := strconv.Atoi(doc.DocumentMinorVersion)
		if err != nil {
			log.WithField("cla_group_id", claGroupID).Error("invalid minor number in cla_group")
			continue
		}
		docVersion := int64(maj)<<32 | int64(min)
		if docVersion > version {
			url = doc.DocumentS3URL
		}
	}
	return url
}

func (s *service) ListClaGroupsUnderFoundation(foundationSFID string) (*models.ClaGroupList, error) {
	out := &models.ClaGroupList{List: make([]*models.ClaGroup, 0)}
	v1ClaGroups, err := s.v1ProjectService.GetClaGroupsByFoundationSFID(foundationSFID, DontLoadDetails)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*models.ClaGroup)
	for _, v1ClaGroup := range v1ClaGroups.Projects {
		cg := &models.ClaGroup{
			CclaEnabled:         v1ClaGroup.ProjectCCLAEnabled,
			CclaRequiresIcla:    v1ClaGroup.ProjectCCLARequiresICLA,
			ClaGroupDescription: v1ClaGroup.ProjectDescription,
			ClaGroupID:          v1ClaGroup.ProjectID,
			ClaGroupName:        v1ClaGroup.ProjectName,
			FoundationSfid:      v1ClaGroup.FoundationSFID,
			IclaEnabled:         v1ClaGroup.ProjectICLAEnabled,
			CclaPdfURL:          getS3Url(v1ClaGroup.ProjectID, v1ClaGroup.ProjectCorporateDocuments),
			IclaPdfURL:          getS3Url(v1ClaGroup.ProjectID, v1ClaGroup.ProjectIndividualDocuments),
			ProjectSfidList:     make([]string, 0),
		}
		m[cg.ClaGroupID] = cg
	}
	// Fill projectSFID list in cla group
	cgprojects, err := s.projectsClaGroupsRepo.GetProjectsIdsForFoundation(foundationSFID)
	if err != nil {
		return nil, err
	}
	for _, cgproject := range cgprojects {
		cg, ok := m[cgproject.ClaGroupID]
		if !ok {
			log.Warnf("stale data present in cla-group-projects table. cla_group_id : %s", cgproject.ClaGroupID)
			continue
		}
		cg.ProjectSfidList = append(cg.ProjectSfidList, cgproject.ProjectSFID)
	}
	// now build output array
	for _, cg := range m {
		out.List = append(out.List, cg)
	}
	return out, nil
}
