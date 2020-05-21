package organization_service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/communitybridge/easycla/cla-backend-go/token"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/client/organizations"
	"github.com/communitybridge/easycla/cla-backend-go/v2/organization-service/models"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// Client is client for organization_service
type Client struct {
	cl *client.OrganziationService
}

var (
	organizationServiceClient *Client
	//ErrScopeNotFound returns error when getting scopeID
	ErrScopeNotFound = errors.New("scope not found")
)

// InitClient initializes the user_service client
func InitClient(APIGwURL string) {
	APIGwURL = strings.ReplaceAll(APIGwURL, "https://", "")
	organizationServiceClient = &Client{
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     APIGwURL,
			BasePath: "user-service/v1",
			Schemes:  []string{"https"},
		}),
	}
}

// GetClient return user_service client
func GetClient() *Client {
	return organizationServiceClient
}

// CreateOrgUserRoleOrgScope attached role scope for particular org and user
func (osc *Client) CreateOrgUserRoleOrgScope(emailID string, organizationID string, roleID string) error {
	params := &organizations.CreateOrgUsrRoleScopesParams{
		CreateRoleScopes: &models.CreateRolescopes{
			EmailAddress: &emailID,
			ObjectID:     &organizationID,
			ObjectType:   aws.String("organization"),
			RoleID:       &roleID,
		},
		SalesforceID: organizationID,
	}
	tok, err := token.GetToken()
	if err != nil {
		return err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	log.Debugf("CreateOrgUserRoleOrgScope: called with args emailID: %s, organizationID: %s, roleID: %s\n", emailID, organizationID, roleID)
	result, err := osc.cl.Organizations.CreateOrgUsrRoleScopes(params, clientAuth)
	if err != nil {
		log.Error("CreateOrgUserRoleOrgScope failed", err)
		return err
	}
	log.Debugf("CreateOrgUserRoleOrgScope: result: %#v\n", result)
	return nil
}

// CreateOrgUserRoleOrgScopeProjectOrg assigns role scope to user
func (osc *Client) CreateOrgUserRoleOrgScopeProjectOrg(emailID string, projectID string, organizationID string, roleID string) error {

	params := &organizations.CreateOrgUsrRoleScopesParams{
		CreateRoleScopes: &models.CreateRolescopes{
			EmailAddress: &emailID,
			ObjectID:     aws.String(fmt.Sprintf("%s|%s", projectID, organizationID)),
			ObjectType:   aws.String("project|organization"),
			RoleID:       &roleID,
		},
		SalesforceID: organizationID,
	}
	tok, err := token.GetToken()
	if err != nil {
		return err
	}

	clientAuth := runtimeClient.BearerToken(tok)
	log.Debugf("CreateOrgUserRoleScope: called with args emailID: %s, projectID: %s, organizationID: %s, roleID: %s", emailID, projectID, organizationID, roleID)
	result, err := osc.cl.Organizations.CreateOrgUsrRoleScopes(params, clientAuth)
	if err != nil {
		log.Error("CreateOrgUserRoleScope failed", err)
		return err
	}
	log.Debugf("CreateOrgUserRoleOrgScope: result: %#v\n", result)
	return nil
}

// DeleteOrgUserRoleOrgScopeProjectOrg removes role scope for user
func (osc *Client) DeleteOrgUserRoleOrgScopeProjectOrg(organizationID string, roleID string, scopeID string, userName string, userEmail string) error {

	params := &organizations.DeleteOrgUsrRoleScopesParams{
		SalesforceID: organizationID,
		RoleID:       roleID,
		ScopeID:      scopeID,
		XUSERNAME:    &userName,
		XEMAIL:       &userEmail,
	}
	tok, err := token.GetToken()
	if err != nil {
		return err
	}

	clientAuth := runtimeClient.BearerToken(tok)
	log.Debugf("DeleteOrgUserRoleOrgScopeProjectOrg called with organizationID: %s, roleID: %s, scopeID: %s", organizationID, roleID, scopeID)
	result, deleteErr := osc.cl.Organizations.DeleteOrgUsrRoleScopes(params, clientAuth)
	if deleteErr != nil {
		log.Error("DeleteOrgUserRoleOrgScopeProjectOrg failed", deleteErr)
		return deleteErr
	}
	log.Debugf("CreateOrgUserRoleOrgScope: result: %#v\n", result)
	return nil
}

// GetScopeID will return scopeID for a give role
func (osc *Client) GetScopeID(organizationID string, roleName string, objectTypeName string) (string, error) {
	tok, err := token.GetToken()
	if err != nil {
		return "", err
	}
	params := &organizations.ListOrgUsrServiceScopesParams{
		SalesforceID: organizationID,
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := osc.cl.Organizations.ListOrgUsrServiceScopes(params, clientAuth)
	if err != nil {
		return "", err
	}
	data := result.Payload
	for _, userRole := range data.Userroles {
		for _, roleScopes := range userRole.RoleScopes {
			if roleScopes.RoleName == roleName {
				for _, scope := range roleScopes.Scopes {
					if scope.ObjectTypeName == objectTypeName {
						return scope.ScopeID, nil
					}
				}
			}
		}
	}
	return "", ErrScopeNotFound
}
