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
		Get(fmt.Sprintf("https://%s/infra/projects/?d&pp=0", gerritName))
	if err != nil {
		log.WithFields(f).Warnf("problem querying gerrit %s, error: %+v", gerritName, err)
		return nil, err
	}

	if resp.IsError() {
		msg := fmt.Sprintf("non-success response from list gerrit repos for gerrit %s, error code: %s", gerritName, resp.Status())
		log.WithFields(f).Warn(msg)
		return nil, errors.New(msg)
	}

	// Clean up the output - has this weird prefix at the beginning of the JSON
	//log.Debugf("Body: %s", resp.Body())
	// responseClean := strings.Trim(string(resp.Body()), ")]}'")
	//log.Debugf("Body clean: %s", resp.Body()[4:])

	//log.WithFields(f).Debug("Decoding response...")
	//var responseModel GerritRepoListResponse
	//err = json.Unmarshal(resp.Body()[4:], &responseModel)
	var result map[string]GerritRepoInfo
	err = json.Unmarshal(resp.Body()[4:], &result)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response for gerrit: %s, error: %+v", gerritName, err)
		return nil, err
	}

	//log.WithFields(f).Debugf("response model: %+v", responseModel)
	//log.WithFields(f).Debugf("response model: %+v", result)

	return convertModel(result), nil
}

func convertModel(responseModel map[string]GerritRepoInfo) *models.GerritRepoList {
	var gerritRepos []*models.GerritRepo
	for name, repo := range responseModel {
		gerritRepos = append(gerritRepos, &models.GerritRepo{
			ID:          repo.ID,
			Name:        name,
			Description: repo.Description,
			State:       repo.State,
		})
	}

	return &models.GerritRepoList{
		List: gerritRepos,
	}
}
