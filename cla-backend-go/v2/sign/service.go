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

	"github.com/sirupsen/logrus"

	acs_service "github.com/communitybridge/easycla/cla-backend-go/v2/acs-service"

	organization_service "github.com/communitybridge/easycla/cla-backend-go/v2/organization-service"

	user_service "github.com/communitybridge/easycla/cla-backend-go/v2/user-service"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/company"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
)

// ProjectRepo contains project repo methods
type ProjectRepo interface {
	GetProjectByID(projectID string) (*v1Models.Project, error)
}

// Service interface defines the sign service methods
type Service interface {
	RequestCorporateSignature(authorizationHeader string, input *models.CorporateSignatureInput) (*models.CorporateSignatureOutput, error)
}

// service
type service struct {
	ClaV1ApiURL string
	companyRepo company.IRepository
	projectRepo ProjectRepo
}

// NewService returns an instance of v2 project service
func NewService(apiURL string, compRepo company.IRepository, projectRepo ProjectRepo) Service {
	return &service{
		ClaV1ApiURL: apiURL,
		companyRepo: compRepo,
		projectRepo: projectRepo,
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
		if input.AuthorityName == "" || input.AuthorityEmail == "" {
			return errors.New("require authority_name and authority_email")
		}
	}
	return nil
}

func (s *service) RequestCorporateSignature(authorizationHeader string, input *models.CorporateSignatureInput) (*models.CorporateSignatureOutput, error) {
	err := validateCorporateSignatureInput(input)
	if err != nil {
		return nil, err
	}
	comp, err := s.companyRepo.GetCompanyByExternalID(utils.StringValue(input.CompanySfid))
	if err != nil {
		return nil, err
	}
	proj, err := s.projectRepo.GetProjectByID(input.ClaGroupID.String())
	if err != nil {
		return nil, err
	}
	if proj.ProjectExternalID != utils.StringValue(input.ProjectSfid) {
		return nil, errors.New("project_sfid does not match with cla_groups project_sfid")
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
		return nil, err
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
	}
	var out requestCorporateSignatureOutput
	err = json.Unmarshal(responseBody, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func prepareUserForSigning(userEmail string, companySFID, projectSFID string) error {
	role := "cla-signatory"
	f := logrus.Fields{"user_email": userEmail, "company_sfid": companySFID, "project_sfid": projectSFID}
	log.WithFields(f).Debug("prepareUserForSigning called")
	usc := user_service.GetClient()
	// search user
	log.WithFields(f).Debug("searching user by email")
	user, err := usc.SearchUserByEmail(userEmail)
	if err != nil {
		log.WithFields(f).Errorf("search user by email failed: %v", err)
		return err
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
	log.WithFields(f).Debugf("checking if user have role of %s", role)
	haveRole, err := osc.IsUserHaveRoleScope(roleID, user.ID, companySFID, projectSFID)
	if err != nil {
		log.WithFields(f).Errorf("checking user have role of %s. failed: %v", role, err)
		return err
	}
	log.WithFields(f).Debugf("user have role %s: status %v", role, haveRole)
	// make user cla-signatory
	if !haveRole {
		log.WithFields(f).Debugf("assigning user role of %s", role)
		err = osc.CreateOrgUserRoleOrgScopeProjectOrg(userEmail, projectSFID, companySFID, roleID)
		if err != nil {
			log.WithFields(f).Errorf("assigning user role of %s failed: %v", role, err)
			return err
		}
	}
	return nil
}
