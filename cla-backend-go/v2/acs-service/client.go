// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package acs_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	ProjectOrgScope     = "project|organization"
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
		Context: context.Background(),
	}
	if scope == ProjectOrgScope && projectID == nil {
		log.WithFields(f).Warnf("Project ID required for project|organization scope, error: %+v", ErrProjectIDMissing)
		return ErrProjectIDMissing
	}
	if scope == ProjectOrgScope {
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
	url := fmt.Sprintf("%s/acs/v1/api/roles?search=%s", ac.apiGwURL, roleName)
	tok, err := token.GetToken()
	if err != nil {
		log.WithFields(f).Warnf("problem obtaining token, error: %+v", err)
		return "", err
	}
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
