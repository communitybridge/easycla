// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/projects_cla_groups"

	"github.com/sirupsen/logrus"

	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

	organization_service "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"

	project_service "github.com/communitybridge/easycla/cla-backend-go/v2/project-service"
	user_service "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

// constants
const (
	DontLoadRepoDetails = false
)

// errors
var (
	ErrCCLANotEnabled        = errors.New("corporate license agreement is not enabled with this project")
	ErrTemplateNotConfigured = errors.New("cla template not configured for this project")
	ErrNotInOrg              error
)

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetCLAGroupByID(projectID string, loadRepoDetails bool) (*v1Models.Project, error)
}

// Service interface defines the sign service methods
type Service interface {
	RequestCorporateSignature(lfUsername string, authorizationHeader string, input *models.CorporateSignatureInput) (*models.CorporateSignatureOutput, error)
}

// service
type service struct {
	ClaV1ApiURL          string
	companyRepo          company.IRepository
	projectRepo          ProjectRepo
	projectClaGroupsRepo projects_cla_groups.Repository
	companyService       company.IService
}

// NewService returns an instance of v2 project service
func NewService(apiURL string, compRepo company.IRepository, projectRepo ProjectRepo, pcgRepo projects_cla_groups.Repository, compService company.IService) Service {
	return &service{
		ClaV1ApiURL:          apiURL,
		companyRepo:          compRepo,
		projectRepo:          projectRepo,
		projectClaGroupsRepo: pcgRepo,
		companyService:       compService,
	}
}

type requestCorporateSignatureInput struct {
	ProjectID      string `json:"project_id,omitempty"`
	CompanyID      string `json:"company_id,omitempty"`
	SendAsEmail    bool   `json:"send_as_email,omitempty"`
	AuthorityName  string `json:"authority_name,omitempty"`
	AuthorityEmail string `json:"authority_email,omitempty"`
	ReturnURL      string `json:"return_url,omitempty"`
}

type requestCorporateSignatureOutput struct {
	ProjectID   string `json:"project_id"`
	CompanyID   string `json:"company_id"`
	SignatureID string `json:"signature_id"`
	SignURL     string `json:"sign_url"`
}

func (in *requestCorporateSignatureOutput) toModel() *models.CorporateSignatureOutput {
	return &models.CorporateSignatureOutput{
		SignURL:     in.SignURL,
		SignatureID: in.SignatureID,
	}
}

func validateCorporateSignatureInput(input *models.CorporateSignatureInput) error {
	if input.SendAsEmail {
		if input.AuthorityName == "" {
			return errors.New("require authority_name")
		}
		if input.AuthorityEmail == "" {
			return errors.New("require authority_email")
		}
	} else {
		if input.ReturnURL.String() == "" {
			return errors.New("require return_url")
		}
	}
	if input.ProjectSfid == nil || *input.ProjectSfid == "" {
		return errors.New("require project_sfid")
	}
	if input.CompanySfid == nil || *input.CompanySfid == "" {
		return errors.New("require company_sfid")
	}
	return nil
}

func (s *service) RequestCorporateSignature(lfUsername string, authorizationHeader string, input *models.CorporateSignatureInput) (*models.CorporateSignatureOutput, error) {
	err := validateCorporateSignatureInput(input)
	if err != nil {
		return nil, err
	}
	comp, err := s.companyRepo.GetCompanyByExternalID(utils.StringValue(input.CompanySfid))
	if err != nil {
		return nil, err
	}
	psc := project_service.GetClient()
	project, err := psc.GetProject(utils.StringValue(input.ProjectSfid))
	if err != nil {
		return nil, err
	}
	var claGroupID string
	if project.Parent == "" {
		// this is root project
		cgmlist, perr := s.projectClaGroupsRepo.GetProjectsIdsForFoundation(utils.StringValue(input.ProjectSfid))
		if perr != nil {
			return nil, perr
		}
		if len(cgmlist) == 0 {
			// no cla group is link with root_project
			return nil, projects_cla_groups.ErrProjectNotAssociatedWithClaGroup
		}
		claGroups := utils.NewStringSet()
		for _, cg := range cgmlist {
			claGroups.Add(cg.ClaGroupID)
		}
		if claGroups.Length() > 1 {
			// multiple cla group are linked with root_project
			// so we can not determine which cla-group to use
			return nil, errors.New("invalid project_sfid. multiple cla-groups are associated with this project_sfid")
		}
		claGroupID = (claGroups.List())[0]
	} else {
		cgm, perr := s.projectClaGroupsRepo.GetClaGroupIDForProject(utils.StringValue(input.ProjectSfid))
		if perr != nil {
			return nil, perr
		}
		claGroupID = cgm.ClaGroupID
	}

	proj, err := s.projectRepo.GetCLAGroupByID(claGroupID, DontLoadRepoDetails)
	if err != nil {
		return nil, err
	}
	if !proj.ProjectCCLAEnabled {
		return nil, ErrCCLANotEnabled
	}
	if len(proj.ProjectCorporateDocuments) == 0 {
		return nil, ErrTemplateNotConfigured
	}
	if input.SendAsEmail {
		// this would be used only in case of cla-signatory
		err = prepareUserForSigning(input.AuthorityEmail.String(), utils.StringValue(input.CompanySfid), utils.StringValue(input.ProjectSfid))
		if err != nil {
			return nil, err
		}
	}
	out, err := requestCorporateSignature(authorizationHeader, s.ClaV1ApiURL, &requestCorporateSignatureInput{
		ProjectID:      proj.ProjectID,
		CompanyID:      comp.CompanyID,
		SendAsEmail:    input.SendAsEmail,
		AuthorityName:  input.AuthorityName,
		AuthorityEmail: input.AuthorityEmail.String(),
		ReturnURL:      input.ReturnURL.String(),
	})
	if err != nil {
		// remove role
		removeErr := removeSignatoryRole(input.AuthorityEmail.String(), utils.StringValue(input.CompanySfid), utils.StringValue(input.ProjectSfid))
		if removeErr != nil {
			return nil, removeErr
		}
		return nil, err
	}

	// Update the company ACL
	companyACLError := s.companyService.AddUserToCompanyAccessList(comp.CompanyID, lfUsername)
	if companyACLError != nil {
		log.Warnf("AddCLAManager- Unable to add user to company ACL, companyID: %s, user: %s, error: %+v", *input.CompanySfid, lfUsername, companyACLError)
		return nil, companyACLError
	}

	return out.toModel(), nil
}

func requestCorporateSignature(authToken string, apiURL string, input *requestCorporateSignatureInput) (*requestCorporateSignatureOutput, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	log.Debugf("requesting corporate signatures: %#v\n", string(requestBody))
	req, err := http.NewRequest("POST", apiURL+"/v1/request-corporate-signature", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("corporate signature response: %#v\n", string(responseBody))
	log.Debugf("corporate signature response headers :%#v\n", resp.Header)
	if strings.Contains(string(responseBody), "Company has already signed CCLA with this project") {
		return nil, errors.New("company has already signed CCLA with this project")
	} else if strings.Contains(string(responseBody), "Contract Group does not support CCLAs.") {
		return nil, errors.New("contract Group does not support CCLAs")
	} else if strings.Contains(string(responseBody), "user_error': 'user does not exist") {
		return nil, errors.New("user_error': 'user does not exist")
	}
	var out requestCorporateSignatureOutput
	err = json.Unmarshal(responseBody, &out)
	if err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, errors.New(string(responseBody))
		}
		return nil, err
	}
	return &out, nil
}

func removeSignatoryRole(userEmail string, companySFID string, projectSFID string) error {
	f := logrus.Fields{"user_email": userEmail, "company_sfid": companySFID, "project_sfid": projectSFID}
	log.WithFields(f).Debug("removing role for user")

	usc := user_service.GetClient()
	// search user
	log.WithFields(f).Debug("searching user by email")
	user, err := usc.SearchUserByEmail(userEmail)
	if err != nil {
		log.WithFields(f).Debug("Failed to get user")
		return err
	}

	log.WithFields(f).Debug("Getting role id")
	acsClient := acs_service.GetClient()
	roleID, roleErr := acsClient.GetRoleID("cla-signatory")
	if roleErr != nil {
		log.WithFields(f).Debug("Failed to get role id for cla-signatory")
		return roleErr
	}
	// Get scope id
	log.WithFields(f).Debug("getting scope id")
	orgClient := organization_service.GetClient()
	scopeID, scopeErr := orgClient.GetScopeID(companySFID, projectSFID, "cla-signatory", "project|organization", user.Username)

	if scopeErr != nil {
		log.WithFields(f).Debug("Failed to get scope id for cla-signatory role")
		return scopeErr
	}

	//Unassign role
	log.WithFields(f).Debug("Unassigning role")
	deleteErr := orgClient.DeleteOrgUserRoleOrgScopeProjectOrg(companySFID, roleID, scopeID, &user.Username, &userEmail)

	if deleteErr != nil {
		log.WithFields(f).Debug("Failed to remove cla-signatory role")
		return deleteErr
	}

	return nil

}

func prepareUserForSigning(userEmail string, companySFID, projectSFID string) error {
	var ErrNotInOrg error
	role := "cla-signatory"
	f := logrus.Fields{"user_email": userEmail, "company_sfid": companySFID, "project_sfid": projectSFID}
	log.WithFields(f).Debug("prepareUserForSigning called")
	usc := user_service.GetClient()
	// search user
	log.WithFields(f).Debug("searching user by email")
	user, err := usc.SearchUserByEmail(userEmail)

	if err != nil {
		log.Debugf("User with email : %s does not have an LF login", userEmail)
		return nil
	}
	log.WithFields(f).Debugf("user type is %s", user.Type)
	if user.Type == "lead" {
		// convert user to contact
		log.WithFields(f).Debug("converting lead to contact")
		err = usc.ConvertToContact(user.ID)
		if err != nil {
			log.WithFields(f).Errorf("converting lead to contact failed: %v", err)
			return err
		}
	}
	ac := acs_service.GetClient()
	log.WithFields(f).Debugf("getting role_id for %s", role)
	roleID, err := ac.GetRoleID(role)
	if err != nil {
		fmt.Println("error", err)
		log.WithFields(f).Errorf("getting role_id for %s failed: %v", role, err.Error())
		return err
	}
	log.Debugf("role %s, role_id %s", role, roleID)
	// assign user role of cla signatory for this project
	osc := organization_service.GetClient()
	if err != nil {
		return err
	}
	// make user cla-signatory
	log.WithFields(f).Debugf("assigning user role of %s", role)
	err = osc.CreateOrgUserRoleOrgScopeProjectOrg(userEmail, projectSFID, companySFID, roleID)
	if err != nil {
		if strings.Contains(err.Error(), "associated with some organization") {
			ErrNotInOrg = fmt.Errorf("user: %s already associated with some organization", user.Username)
			return ErrNotInOrg
		}
		log.WithFields(f).Errorf("assigning user role of %s failed: %v", role, err)
		return err
	}
	return nil
}
