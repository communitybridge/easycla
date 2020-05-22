package acs_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/token"

	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/client"

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
