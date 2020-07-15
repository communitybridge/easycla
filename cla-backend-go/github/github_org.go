package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/google/go-github/github"
)

// errors
var (
	ErrGithubOrganizationNotFound = errors.New("github organization name not found")
)

// GetOrganization gets github organization
func GetOrganization(organizationName string) (*github.Organization, error) {
	client := newGithubOauthClient()
	org, resp, err := client.Organizations.Get(context.TODO(), organizationName)
	if err != nil {
		logging.Warnf("GetOrganization %s failed. error = %s", organizationName, err.Error())
		if resp.StatusCode == 404 {
			return nil, ErrGithubOrganizationNotFound
		}
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("GetOrganization %s failed with no success response code %d. error = %s", organizationName, resp.StatusCode, err.Error())
		logging.Warnf(msg)
		return nil, errors.New(msg)
	}
	return org, nil
}
