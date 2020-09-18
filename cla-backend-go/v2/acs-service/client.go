// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package acs_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/sirupsen/logrus"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/token"

	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/client"
	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/client/invite"
	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/models"
	runtimeClient "github.com/go-openapi/runtime/client"

	"errors"

	"github.com/go-openapi/strfmt"
)

// Client is client for acs_service
type Client struct {
	apiKey   string
	apiGwURL string
	cl       *client.CentralAuthorizationLayerForTheLFXPlatform
}

var (
	acsServiceClient *Client
)

// errors
var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrProjectIDMissing = errors.New("project ID missing")
)

// InitClient initializes the acs_service client
func InitClient(APIGwURL string, apiKey string) {
	url := strings.ReplaceAll(APIGwURL, "https://", "")
	acsServiceClient = &Client{
		apiKey:   apiKey,
		apiGwURL: APIGwURL,
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     url,
			BasePath: "acs/v1/api",
			Schemes:  []string{"https"},
		}),
	}
}

// GetClient return user_service client
func GetClient() *Client {
	return acsServiceClient
}

// SendUserInvite invites users to the LFX platform
func (ac *Client) SendUserInvite(email *string,
	roleName string, scope string, projectID *string, organizationID string, inviteType string, subject *string, emailContent *string, automate bool) error {
	f := logrus.Fields{
		"functionName":   "SendUserInvite",
		"roleName":       roleName,
		"scope":          scope,
		"organizationID": organizationID,
		"inviteType":     inviteType,
	}

	if email != nil {
		f["email"] = *email
	}
	if projectID != nil {
		f["projectID"] = *projectID
	}
	if subject != nil {
		f["subject"] = *subject
	}

	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params := &invite.CreateUserInviteParams{
		SendInvite: &models.CreateInvite{
			Automate: automate,
			Email:    email,
			Scope:    scope,
			RoleName: roleName,
			Type:     inviteType,
		},
		Context: utils.NewContext(),
	}
	if scope == utils.ProjectOrgScope && projectID == nil {
		log.WithFields(f).Warnf("Project ID required for project|organization scope, error: %+v", ErrProjectIDMissing)
		return ErrProjectIDMissing
	}
	if scope == utils.ProjectOrgScope {
		// Set project|organization scope
		params.SendInvite.ScopeID = fmt.Sprintf("%s|%s", *projectID, organizationID)
	} else {
		params.SendInvite.ScopeID = organizationID
	}
	if subject != nil {
		f["subject"] = *subject
		params.SendInvite.Subject = *subject
	}
	// Pass emailContent if passed in the args
	if emailContent != nil {
		params.SendInvite.Body = *emailContent
	}
	result, inviteErr := ac.cl.Invite.CreateUserInvite(params, clientAuth)
	log.Debugf("CreateUserinvite called with args email: %s, scope: %s, roleName: %s, type: %s, scopeID: %s",
		*email, scope, roleName, inviteType, organizationID)
	if inviteErr != nil {
		log.WithFields(f).Errorf("CreateUserInvite failed for payload : %+v : %v", params, inviteErr)
		return nil
	}

	log.WithFields(f).Debugf("CreatedUserInvite :%+v", result.Payload)
	return nil
}

// GetRoleID will return roleID for the provided role name
func (ac *Client) GetRoleID(roleName string) (string, error) {
	f := logrus.Fields{
		"functionName": "GetRoleID",
		"roleName":     roleName,
	}

	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return "", err
	}

	url := fmt.Sprintf("%s/acs/v1/api/roles?search=%s", ac.apiGwURL, roleName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithFields(f).Warnf("problem making a new GET request for url: %s, error: %+v", url, err)
		return "", err
	}
	req.Header.Set("X-API-KEY", ac.apiKey)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(f).Warnf("problem invoking http GET request to url: %s, error: %+v", url, err)
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.WithFields(f).Warnf("error closing resource: %+v", closeErr)
		}
	}()
	var roles []struct {
		RoleName string `json:"role_name"`
		RoleID   string `json:"role_id"`
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(f).Warnf("problem reading response body, error: %+v", err)
		return "", err
	}
	err = json.Unmarshal(b, &roles)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response body, error: %+v", err)
		return "", err
	}
	for _, role := range roles {
		if role.RoleName == roleName {
			return role.RoleID, nil
		}
	}

	return "", ErrRoleNotFound
}

// GetObjectTypeIDByName will return object type ID for the provided role name
func (ac *Client) GetObjectTypeIDByName(objectType string) (int, error) {
	f := logrus.Fields{
		"functionName": "GetObjectTypeID",
		"objectType":   objectType,
	}

	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return 0, err
	}

	url := fmt.Sprintf("%s/acs/v1/api/object-types", ac.apiGwURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithFields(f).Warnf("problem making a new GET request for url: %s, error: %+v", url, err)
		return 0, err
	}
	req.Header.Set("X-API-KEY", ac.apiKey)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(f).Warnf("problem invoking http GET request to url: %s, error: %+v", url, err)
		return 0, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.WithFields(f).Warnf("error closing resource: %+v", closeErr)
		}
	}()
	var objectTypes []struct {
		TypeID    int    `json:"type_id"`
		Name      string `json:"name"`
		CreatedAt int    `json:"created_at"`
		UpdatedAt int    `json:"updated_at"`
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(f).Warnf("problem reading response body, error: %+v", err)
		return 0, err
	}
	err = json.Unmarshal(b, &objectTypes)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response body, error: %+v", err)
		return 0, err
	}

	for _, role := range objectTypes {
		if role.Name == objectType {
			return role.TypeID, nil
		}
	}

	return 0, ErrRoleNotFound
}

// GetAssignedRoles will return assigned roles based on the roleName, project and organization SFID
func (ac *Client) GetAssignedRoles(roleName, projectSFID, organizationSFID string) (*models.ObjectRoleScope, error) {
	f := logrus.Fields{
		"functionName":     "GetAssignedRole",
		"roleName":         roleName,
		"projectSFID":      projectSFID,
		"organizationSFID": organizationSFID,
	}

	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return nil, err
	}

	// Lookup the Project|Organization type
	objectTypeID, err := ac.GetObjectTypeIDByName(utils.ProjectOrgScope)
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/acs/v1/api/object-types/%d/roles?ojectid=%s|%s", ac.apiGwURL, objectTypeID, projectSFID, organizationSFID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithFields(f).Warnf("problem making a new GET request for url: %s, error: %+v", url, err)
		return nil, err
	}
	req.Header.Set("X-API-KEY", ac.apiKey)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(f).Warnf("problem invoking http GET request to url: %s, error: %+v", url, err)
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.WithFields(f).Warnf("error closing resource: %+v", closeErr)
		}
	}()

	var response *models.ObjectRoleScope
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(f).Warnf("problem reading response body, error: %+v", err)
		return nil, err
	}
	err = json.Unmarshal(b, response)
	if err != nil {
		log.WithFields(f).Warnf("problem unmarshalling response body, error: %+v", err)
		return nil, err
	}

	return response, nil
}

// DeleteRoleByID will delete the specified role by ID
func (ac *Client) DeleteRoleByID(roleID string) error {
	f := logrus.Fields{
		"functionName": "DeleteRoleByID",
		"roleID":       roleID,
	}

	if roleID == "" {
		log.WithFields(f).Warn("unable to delete role by ID - role ID is empty")
		return errors.New("empty role ID")
	}

	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return err
	}

	url := fmt.Sprintf("%s/acs/v1/api/roles/%s", ac.apiGwURL, roleID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.WithFields(f).Warnf("problem making a new DELETE request for url: %s, error: %+v", url, err)
		return err
	}
	req.Header.Set("X-API-KEY", ac.apiKey)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(f).Warnf("problem invoking http DELETE request to url: %s, error: %+v", url, err)
		return err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.WithFields(f).Warnf("error closing resource: %+v", closeErr)
		}
	}()

	if resp.StatusCode != 204 {
		log.WithFields(f).Warnf("non-success status code returned from delete operation: %d", resp.StatusCode)
	}

	return nil
}
