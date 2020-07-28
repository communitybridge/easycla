// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package gerrits

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	resty "github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Service contains functions of GithubOrganizations service
type Service interface {
	GetGerritRepos(gerritName string) (*models.GerritRepoList, error)
}

type service struct {
}

// NewService creates a new githubOrganizations service
func NewService() Service {
	return service{}
}

func (s service) GetGerritRepos(gerritName string) (*models.GerritRepoList, error) {
	f := logrus.Fields{
		"function":   "GetGerritRepos",
		"gerritName": gerritName,
	}

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("https://%s/%s/projects/?d&pp=0", gerritName, getGerritAPIPath(gerritName)))
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit %s, error: %+v", gerritName, err)
		return nil, err
	}

	if resp.IsError() {
		msg := fmt.Sprintf("non-success response from list gerrit repos for gerrit %s, error code: %s", gerritName, resp.Status())
		log.WithFields(f).Warn(msg)
		return nil, errors.New(msg)
	}

	var result map[string]GerritRepoInfo
	// Need to strip off the leading "magic prefix line" from the response payload, which is: )]}'
	// See: https://gerrit.linuxfoundation.org/infra/Documentation/rest-api.html#output
	err = json.Unmarshal(resp.Body()[4:], &result)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response for gerrit: %s, error: %+v", gerritName, err)
		return nil, err
	}

	return convertModel(result), nil
}

// convertModel is a helper function to create a GerritRepoList response model
func convertModel(responseModel map[string]GerritRepoInfo) *models.GerritRepoList {
	var gerritRepos []*models.GerritRepo
	for name, repo := range responseModel {

		var weblinks []*models.GerritRepoWebLinksItems0
		for _, weblink := range repo.WebLinks {
			weblinks = append(weblinks, &models.GerritRepoWebLinksItems0{
				Name: weblink.Name,
				URL:  weblink.URL,
			})
		}

		gerritRepos = append(gerritRepos, &models.GerritRepo{
			ID:          repo.ID,
			Name:        name,
			Description: repo.Description,
			State:       repo.State,
			WebLinks:    weblinks,
		})
	}

	return &models.GerritRepoList{
		List: gerritRepos,
	}
}

// getGerritAPIPath returns the path to the API based on the gerrit host
func getGerritAPIPath(gerritHost string) string {
	switch gerritHost {
	case "gerrit.linuxfoundation.org":
		return "infra"
	case "gerrit.onap.org":
		return "r"
	case "gerrit.o-ran-sc.org":
		return "r"
	case "gerrit.tungsten.io":
		return "r"
	case "gerrit.opnfv.org":
		return "gerrit"
	default:
		return "r"
	}
}
