package acs_service

import (
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/client"

	"errors"

	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/communitybridge/easycla/cla-backend-go/v2/acs-service/client/role"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// Client is client for acs_service
type Client struct {
	cl *client.CentralAuthorizationLayerForTheLFXPlatform
}

var (
	acsServiceClient *Client
)

// errors
var (
	ErrRoleNotFound = errors.New("role not found")
)

// InitClient initializes the acs_service client
func InitClient(APIGwURL string) {
	APIGwURL = strings.ReplaceAll(APIGwURL, "https://", "")
	acsServiceClient = &Client{
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     APIGwURL,
			BasePath: "acs-service/v1",
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
	tok, err := token.GetToken()
	if err != nil {
		return "", err
	}
	params := &role.GetRolesParams{
		Search: &roleName,
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := ac.cl.Role.GetRoles(params, clientAuth)
	if err != nil {
		return "", err
	}
	for _, r := range result.Payload {
		if utils.StringValue(r.RoleName) == roleName {
			return r.RoleID, nil
		}
	}
	return "", ErrRoleNotFound
}
