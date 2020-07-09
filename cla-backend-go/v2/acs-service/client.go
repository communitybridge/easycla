package acs_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	ErrRoleNotFound = errors.New("role not found")
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
	roleName string, scope string, organizationID string, inviteType string) error {
	tok, err := token.GetToken()
	if err != nil {
		return err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params := &invite.CreateUserInviteParams{
		SendInvite: &models.CreateInvite{
			Email:    email,
			Scope:    scope,
			ScopeID:  organizationID,
			RoleName: roleName,
			Type:     inviteType,
		},
		Context: context.Background(),
	}
	result, inviteErr := ac.cl.Invite.CreateUserInvite(params, clientAuth)
	log.Debugf("CreateUserinvite called with args email: %s, scope: %s, roleName: %s, type: %s, scopeID: %s",
		*email, scope, roleName, inviteType, organizationID)
	if inviteErr != nil {
		log.Error("CreateUserInvite failed", err)
		return err
	}
	log.Debugf("CreatedUserInvite :%+v", result.Payload)
	return nil
}

// GetRoleID will return roleID for the provided role name
func (ac *Client) GetRoleID(roleName string) (string, error) {
	url := fmt.Sprintf("%s/acs/v1/api/roles?search=%s", ac.apiGwURL, roleName)
	tok, err := token.GetToken()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-KEY", ac.apiKey)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var roles []struct {
		RoleName string `json:"role_name"`
		RoleID   string `json:"role_id"`
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &roles)
	if err != nil {
		return "", err
	}
	for _, role := range roles {
		if role.RoleName == roleName {
			return role.RoleID, nil
		}
	}
	return "", ErrRoleNotFound
}
