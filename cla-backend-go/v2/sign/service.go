// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

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
	companyRepo company.IRepository
	projectRepo ProjectRepo
}

// NewService returns an instance of v2 project service
func NewService(compRepo company.IRepository, projectRepo ProjectRepo) Service {
	return &service{
		companyRepo: compRepo,
		projectRepo: projectRepo,
	}
}

type RequestCorporateSignatureInput struct {
	ProjectID      string `json:"project_id,omitempty"`
	CompanyID      string `json:"company_id,omitempty"`
	SendAsEmail    bool   `json:"send_as_email,omitempty"`
	AuthorityName  string `json:"authority_name,omitempty"`
	AuthorityEmail string `json:"authority_email,omitempty"`
	ReturnURL      string `json:"return_url,omitempty"`
}

type RequestCorporateSignatureOutput struct {
	ProjectID   string `json:"project_id"`
	CompanyID   string `json:"company_id"`
	SignatureID string `json:"signature_id"`
	SignURL     string `json:"sign_url"`
}

func (in *RequestCorporateSignatureOutput) toModel() *models.CorporateSignatureOutput {
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
	out, err := requestCorporateSignature(authorizationHeader, "https://api.dev.lfcla.com", &RequestCorporateSignatureInput{
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

func requestCorporateSignature(authToken string, apiURL string, input *RequestCorporateSignatureInput) (*RequestCorporateSignatureOutput, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
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
	if strings.Contains(string(responseBody), "Company has already signed CCLA with this project") {
		return nil, errors.New("company has already signed CCLA with this project")
	}
	var out RequestCorporateSignatureOutput
	err = json.Unmarshal(responseBody, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
