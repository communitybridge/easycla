// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// Service handles gerrit Repository service
type Service interface {
	DeleteClaGroupGerrits(claGroupID string) (int, error)
	DeleteGerrit(gerritID string) error
	GetGerrit(gerritID string) (*models.Gerrit, error)
	AddGerrit(claGroupID string, projectSFID string, input *models.AddGerritInput, claGroupModel *models.ClaGroup) (*models.Gerrit, error)
	GetClaGroupGerrits(claGroupID string, projectSFID *string) (*models.GerritList, error)
	GetGerritRepos(gerritName string) (*models.GerritRepoList, error)
}

type service struct {
	repo    Repository
	lfGroup *LFGroup
}

// NewService creates a new gerrit service
func NewService(repo Repository, lfg *LFGroup) Service {
	return service{
		repo:    repo,
		lfGroup: lfg,
	}
}

func (s service) DeleteClaGroupGerrits(claGroupID string) (int, error) {
	gerrits, err := s.repo.GetClaGroupGerrits(claGroupID, nil)
	if err != nil {
		return 0, err
	}
	if len(gerrits.List) > 0 {
		log.Debugf(fmt.Sprintf("Deleting gerrits for cla-group :%s ", claGroupID))
		for _, gerrit := range gerrits.List {
			err = s.repo.DeleteGerrit(gerrit.GerritID.String())
			if err != nil {
				return 0, err
			}
		}
	}
	return len(gerrits.List), nil
}

func (s service) DeleteGerrit(gerritID string) error {
	return s.repo.DeleteGerrit(gerritID)
}

func (s service) GetGerrit(gerritID string) (*models.Gerrit, error) {
	return s.repo.GetGerrit(gerritID)
}

func (s service) AddGerrit(claGroupID string, projectSFID string, params *models.AddGerritInput, claGroupModel *models.ClaGroup) (*models.Gerrit, error) {
	f := logrus.Fields{
		"claGroupID":  claGroupID,
		"projectSFID": projectSFID,
	}
	if params.GroupIDIcla == "" && params.GroupIDCcla == "" {
		return nil, errors.New("should specify at least a LDAP group for ICLA or CCLA")
	}

	log.WithFields(f).Debugf("cla groupID %s", claGroupID)
	log.WithFields(f).Debugf("project Model %+v", claGroupModel)

	if claGroupModel.ProjectCCLAEnabled && claGroupModel.ProjectICLAEnabled {
		if params.GroupIDCcla == "" {
			return nil, errors.New("please provide GroupIDCcla")
		}
		if params.GroupIDIcla == "" {
			return nil, errors.New("please provide GroupIDIcla")
		}
	} else if claGroupModel.ProjectCCLAEnabled {
		if params.GroupIDCcla == "" {
			return nil, errors.New("please provide GroupIDCcla")
		}
	} else if claGroupModel.ProjectICLAEnabled {
		if params.GroupIDIcla == "" {
			return nil, errors.New("please provide GroupIDIcla")
		}
	}

	if params.GroupIDIcla == params.GroupIDCcla {
		return nil, errors.New("LDAP group for ICLA and CCLA are same")
	}

	if params.GerritName == nil {
		return nil, errors.New("gerrit_name required")
	}

	gerritObject, err := s.repo.ExistsByName(*params.GerritName)
	if err != nil {
		message := fmt.Sprintf("unable to get gerrit by name : %s", *params.GerritName)
		log.WithFields(f).WithError(err).Warnf(message)
	}

	if len(gerritObject) > 0 {
		return nil, errors.New("gerrit_name already present in the system")
	}

	gerritCcla, err := s.repo.GetGerritsByID(params.GroupIDCcla, "CCLA")
	if err != nil {
		message := fmt.Sprintf("unable to get gerrit by ccla id : %s", params.GroupIDCcla)
		log.WithFields(f).WithError(err).Warnf(message)
	}

	if len(gerritCcla.List) > 0 {
		return nil, errors.New("gerrit_ccla id already present in the system")
	}

	gerritIcla, err := s.repo.GetGerritsByID(params.GroupIDIcla, "ICLA")
	if err != nil {
		message := fmt.Sprintf("unable to get gerrit by icla : %s", params.GroupIDIcla)
		log.WithFields(f).WithError(err).Warnf(message)
	}

	if len(gerritIcla.List) > 0 {
		return nil, errors.New("gerrit_icla id already present in the system")
	}

	if params.GerritURL == nil {
		return nil, errors.New("gerrit_url required")
	}

	var groupNameCcla, groupNameIcla string
	if params.GroupIDIcla != "" {
		group, err := s.lfGroup.GetGroup(params.GroupIDIcla)
		if err != nil {
			message := fmt.Sprintf("unable to get LDAP ICLA Group: %s", params.GroupIDIcla)
			log.WithFields(f).WithError(err).Warnf(message)
			return nil, errors.New(message)
		}
		groupNameIcla = group.Title
	}
	if params.GroupIDCcla != "" {
		group, err := s.lfGroup.GetGroup(params.GroupIDCcla)
		if err != nil {
			message := fmt.Sprintf("unable to get LDAP CCLA Group: %s", params.GroupIDCcla)
			log.WithFields(f).WithError(err).Warnf(message)
			return nil, errors.New(message)
		}
		groupNameCcla = group.Title
	}
	input := &models.Gerrit{
		GerritName:    utils.StringValue(params.GerritName),
		GerritURL:     strfmt.URI(*params.GerritURL),
		GroupIDCcla:   params.GroupIDCcla,
		GroupIDIcla:   params.GroupIDIcla,
		GroupNameCcla: groupNameCcla,
		GroupNameIcla: groupNameIcla,
		ProjectID:     claGroupID,
		ProjectSFID:   projectSFID,
		Version:       params.Version,
	}
	return s.repo.AddGerrit(input)
}

func (s service) GetClaGroupGerrits(claGroupID string, projectSFID *string) (*models.GerritList, error) {
	f := logrus.Fields{"functionName": "GetClaGroupGerrits", "claGroupID": claGroupID, "projectSFID": *projectSFID}
	responseModel, err := s.repo.GetClaGroupGerrits(claGroupID, projectSFID)
	if err != nil {
		log.WithFields(f).Warnf("problem getting CLA Group gerrits, error: %+v", err)
		return nil, err
	}

	// Add the repo list to the response model
	for _, gerrit := range responseModel.List {
		log.WithFields(f).Debugf("Processing gerrit URL: %s", gerrit.GerritURL)

		var gerritHost = gerrit.GerritURL.String()
		gerritHost, err = extractGerritHost(gerritHost, f)
		if err != nil {
			return nil, err
		}

		log.WithFields(f).Debugf("fetching gerrit repos from host: %s", gerritHost)
		gerritRepoList, getRepoErr := s.GetGerritRepos(gerritHost)
		if getRepoErr != nil {
			log.WithFields(f).Warnf("problem fetching gerrit repos from host: %s, error: %+v", gerritHost, err)
			log.Error("skipping", getRepoErr)
			continue
			//return nil, getRepoErr
		}

		// Set the connected flag - for now, we just set this value to true
		for _, repo := range gerritRepoList.Repos {
			repo.Connected = true
		}

		gerrit.GerritRepoList = gerritRepoList
	}

	return responseModel, err
}

func extractGerritHost(gerritHost string, f logrus.Fields) (string, error) {
	if strings.HasPrefix(gerritHost, "http") {
		log.WithFields(f).Debugf("extracting gerrit host from URL: %s", gerritHost)
		u, urlErr := url.Parse(gerritHost)
		if urlErr != nil {
			log.WithFields(f).Warnf("problem converting gerrit URL: %s, error: %+v", gerritHost, urlErr)
			return "", urlErr
		}
		gerritHost = u.Host
		log.WithFields(f).Debugf("extracted gerrit host is: %s", gerritHost)
	}
	return gerritHost, nil
}

func (s service) GetGerritRepos(gerritHost string) (*models.GerritRepoList, error) {
	f := logrus.Fields{
		"functionName": "GetGerritRepos",
		"gerritName":   gerritHost,
	}

	gerritRepos, err := listGerritRepos(gerritHost)
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit host, error: %+v", err)
		return nil, err
	}

	gerritConfig, err := getGerritConfig(gerritHost)
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit config, error: %+v", err)
		return nil, err
	}

	return convertModel(gerritRepos, gerritConfig), nil
}

// convertModel is a helper function to create a GerritRepoList response model
func convertModel(responseModel map[string]GerritRepoInfo, serverInfo *ServerInfo) *models.GerritRepoList {
	var gerritRepos []*models.GerritRepo
	for name, repo := range responseModel {

		var weblinks []*models.GerritRepoWebLinksItems0
		for _, weblink := range repo.WebLinks {
			weblinks = append(weblinks, &models.GerritRepoWebLinksItems0{
				Name: weblink.Name,
				URL:  strfmt.URI(weblink.URL),
			})
		}

		claEnabled := false
		if serverInfo != nil && serverInfo.Auth.UseContributorAgreements {
			claEnabled = true
		}

		gerritRepos = append(gerritRepos, &models.GerritRepo{
			ID:                    repo.ID,
			Name:                  name,
			Description:           repo.Description,
			State:                 repo.State,
			ClaEnabled:            claEnabled,
			ContributorAgreements: buildContributorAgreementDetails(serverInfo),
			WebLinks:              weblinks,
		})
	}

	return &models.GerritRepoList{
		Repos: gerritRepos,
	}
}

// buildContributorAgreementDetails helper function to extract and conver the gerrit server info contributor agreement information into a response data model
func buildContributorAgreementDetails(serverInfo *ServerInfo) []*models.GerritRepoContributorAgreementsItems0 {
	var response []*models.GerritRepoContributorAgreementsItems0

	for _, agreement := range serverInfo.Auth.ContributorAgreements {
		response = append(response, &models.GerritRepoContributorAgreementsItems0{
			Name:        agreement.Name,
			Description: agreement.Description,
			URL:         strfmt.URI(agreement.URL),
		})
	}

	return response
}

// listGerritRepos returns a list of gerrit repositories for the given gerrit host
func listGerritRepos(gerritHost string) (map[string]GerritRepoInfo, error) {
	f := logrus.Fields{
		"functionName": "listGerritRepos",
		"gerritHost":   gerritHost,
	}
	client := resty.New()

	gerritAPIPath, gerritAPIPathErr := getGerritAPIPath(gerritHost)
	if gerritAPIPathErr != nil {
		return nil, gerritAPIPathErr
	}

	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("https://%s/%s/projects/?d&pp=0", gerritHost, gerritAPIPath))
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit host: %s, error: %+v", gerritHost, err)
		return nil, err
	}

	if resp.IsError() {
		msg := fmt.Sprintf("non-success response from list gerrit host repos for gerrit %s, error code: %s", gerritHost, resp.Status())
		log.WithFields(f).Warn(msg)
		return nil, errors.New(msg)
	}

	var result map[string]GerritRepoInfo
	// Need to strip off the leading "magic prefix line" from the response payload, which is: )]}'
	// See: https://gerrit.linuxfoundation.org/infra/Documentation/rest-api.html#output
	err = json.Unmarshal(resp.Body()[4:], &result)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response for gerrit host: %s, error: %+v", gerritHost, err)
		return nil, err
	}

	return result, nil
}

// getGerritConfig returns the gerrit configuration for the specified host
func getGerritConfig(gerritHost string) (*ServerInfo, error) {
	f := logrus.Fields{
		"functionName": "getGerritConfig",
		"gerritHost":   gerritHost,
	}
	client := resty.New()

	gerritAPIPath, gerritAPIPathErr := getGerritAPIPath(gerritHost)
	if gerritAPIPathErr != nil {
		return nil, gerritAPIPathErr
	}

	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("https://%s/%s/config/server/info", gerritHost, gerritAPIPath))
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit config, error: %+v", err)
		return nil, err
	}

	if resp.IsError() {
		msg := fmt.Sprintf("non-success response from list gerrit host config query, error code: %s", resp.Status())
		log.WithFields(f).Warn(msg)
		return nil, errors.New(msg)
	}

	var result ServerInfo
	// Need to strip off the leading "magic prefix line" from the response payload, which is: )]}'
	// See: https://gerrit.linuxfoundation.org/infra/Documentation/rest-api.html#output
	err = json.Unmarshal(resp.Body()[4:], &result)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response for gerrit host: %s, error: %+v", gerritHost, err)
		return nil, err
	}

	return &result, nil
}

// getGerritAPIPath returns the path to the API based on the gerrit host
func getGerritAPIPath(gerritHost string) (string, error) {
	f := logrus.Fields{"functionName": "getGerritAPIPath", "gerritHost": gerritHost}
	switch gerritHost {
	case "gerrit.linuxfoundation.org":
		return "infra", nil
	case "gerrit.onap.org":
		return "r", nil
	case "gerrit.o-ran-sc.org":
		return "r", nil
	case "gerrit.tungsten.io":
		return "r", nil
	case "gerrit.opnfv.org":
		return "gerrit", nil
	default:
		msg := fmt.Sprintf("unsupport gerrit host: %s", gerritHost)
		log.WithFields(f).Warnf(msg)
		return "", errors.New(msg)
	}
}
