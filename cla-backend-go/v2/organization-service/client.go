// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package organization_service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	"github.com/sirupsen/logrus"

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

const (
	projectOrganization = "project|organization"
)

var (
	organizationServiceClient *Client
	v1EventService            events.Service
)

// InitClient initializes the user_service client
func InitClient(APIGwURL string, eventService events.Service) {
	APIGwURL = strings.ReplaceAll(APIGwURL, "https://", "")
	organizationServiceClient = &Client{
		cl: client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
			Host:     APIGwURL,
			BasePath: "organization-service/v1",
			Schemes:  []string{"https"},
		}),
	}
	v1EventService = eventService
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
		Context:      context.Background(),
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

// IsUserHaveRoleScope checks if user have required role and scope
func (osc *Client) IsUserHaveRoleScope(rolename string, userSFID string, organizationID string, projectSFID string) (bool, error) {
	objectID := fmt.Sprintf("%s|%s", projectSFID, organizationID)
	var offset int64
	var pageSize int64 = 1000
	tok, err := token.GetToken()
	if err != nil {
		return false, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	for {
		params := &organizations.ListOrgUsrServiceScopesParams{
			Offset:       aws.String(strconv.FormatInt(offset, 10)),
			PageSize:     aws.String(strconv.FormatInt(pageSize, 10)),
			SalesforceID: organizationID,
			Rolename:     []string{rolename},
			Context:      context.Background(),
		}
		result, err := osc.cl.Organizations.ListOrgUsrServiceScopes(params, clientAuth)
		if err != nil {
			return false, err
		}
		for _, userRole := range result.Payload.Userroles {
			// loop until we find user
			if userRole.Contact.ID != userSFID {
				continue
			}
			for _, rolescope := range userRole.RoleScopes {
				for _, scope := range rolescope.Scopes {
					if scope.ObjectTypeName == projectOrganization && scope.ObjectID == objectID {
						return true, nil
					}
				}
				return false, nil
			}
			return false, nil
		}
		if result.Payload.Metadata.TotalSize < offset+pageSize {
			break
		}
		offset = offset + pageSize
	}
	return false, nil
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
		Context:      context.Background(),
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

// DeleteRolePermissions removes the specified Org/Project user permissions for with the given role
func (osc *Client) DeleteRolePermissions(organizationID, projectID, role string, authUser *auth.User) error {
	f := logrus.Fields{
		"functionName":   "DeleteRolePermissions",
		"organizationID": organizationID,
		"projectID":      projectID,
		"role":           role,
	}

	// First, query the organization for the list of permissions (scopes)
	scopeResponse, err := osc.ListOrgUserScopes(organizationID, []string{role})
	if err != nil {
		log.WithFields(f).Warnf("problem listing org user scopes, error: %+v", err)
		return err
	}

	// For each result...
	for _, userRoleScopes := range scopeResponse.Userroles {
		userName := userRoleScopes.Contact.Username
		userEmail := userRoleScopes.Contact.EmailAddress

		for _, roleScopes := range userRoleScopes.RoleScopes {
			roleID := roleScopes.RoleID
			for _, scope := range roleScopes.Scopes {
				// Encoded as ProjectID|OrganizationID - split them out
				objectList := strings.Split(scope.ObjectID, "|")
				// check objectID having project|organization scope
				if len(objectList) == 2 {
					if scope.ObjectTypeName == projectOrganization && projectID == objectList[0] {
						delErr := osc.DeleteOrgUserRoleOrgScopeProjectOrg(organizationID, roleID, scope.ScopeID, &userName, &userEmail)
						if delErr != nil {
							f["userName"] = userName
							f["userEmail"] = userEmail
							log.WithFields(f).Warnf("problem removing user from role, error: %+v", err)
							return delErr
						}

						// Log Event...
						v1EventService.LogEvent(&events.LogEventArgs{
							EventType:         events.ClaManagerRoleDeleted,
							ProjectID:         projectID,
							ProjectModel:      nil,
							CompanyID:         organizationID,
							CompanyModel:      nil,
							LfUsername:        authUser.UserName,
							UserID:            authUser.UserName,
							UserModel:         nil,
							ExternalProjectID: projectID,
							EventData: &events.ClaManagerRoleCreatedData{
								Role:      role,                 // cla-manager
								Scope:     scope.ObjectTypeName, // project|organization
								UserName:  userName,             // bstonedev
								UserEmail: userEmail,            // bstone+dev@linuxfoundation.org
							},
						})
					}
				}
			}
		}
	}
	return nil
}

// DeleteOrgUserRoleOrgScopeProjectOrg removes role scope for user
func (osc *Client) DeleteOrgUserRoleOrgScopeProjectOrg(organizationID string, roleID string, scopeID string, userName *string, userEmail *string) error {

	params := &organizations.DeleteOrgUsrRoleScopesParams{
		SalesforceID: organizationID,
		RoleID:       roleID,
		ScopeID:      scopeID,
		XUSERNAME:    userName,
		XEMAIL:       userEmail,
		Context:      context.Background(),
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
func (osc *Client) GetScopeID(organizationID string, projectID string, roleName string, objectTypeName string, userLFID string) (string, error) {
	tok, err := token.GetToken()
	if err != nil {
		return "", err
	}
	params := &organizations.ListOrgUsrServiceScopesParams{
		SalesforceID: organizationID,
		Context:      context.Background(),
	}
	clientAuth := runtimeClient.BearerToken(tok)
	result, err := osc.cl.Organizations.ListOrgUsrServiceScopes(params, clientAuth)
	if err != nil {
		return "", err
	}
	data := result.Payload
	for _, userRole := range data.Userroles {
		// Check scopes for given user
		if userRole.Contact.Username == userLFID {
			for _, roleScopes := range userRole.RoleScopes {
				if roleScopes.RoleName == roleName {
					for _, scope := range roleScopes.Scopes {
						// Check object ID and and objectTypeName
						objectList := strings.Split(scope.ObjectID, "|")
						// check objectID having project|organization scope
						if len(objectList) == 2 {
							if scope.ObjectTypeName == objectTypeName && projectID == objectList[0] {
								return scope.ScopeID, nil
							}
						}
					}
				}
			}
		}
	}
	return "", nil
}

// SearchOrganization search organization by name. It will return
// array of organization matching with the orgName.
func (osc *Client) SearchOrganization(orgName string) ([]*models.Organization, error) {
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	var offset int64
	var pageSize int64 = 1000
	clientAuth := runtimeClient.BearerToken(tok)
	var orgs []*models.Organization
	for {
		params := &organizations.SearchOrgParams{
			Name:     aws.String(orgName),
			Offset:   aws.String(strconv.FormatInt(offset, 10)),
			PageSize: aws.String(strconv.FormatInt(pageSize, 10)),
			Context:  context.TODO(),
		}
		result, err := osc.cl.Organizations.SearchOrg(params, clientAuth)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, result.Payload.Data...)
		if result.Payload.Metadata.TotalSize > offset+pageSize {
			offset += pageSize
		} else {
			break
		}
	}
	return orgs, nil
}

// GetOrganization gets organization from organization id
func (osc *Client) GetOrganization(orgID string) (*models.Organization, error) {
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params := &organizations.GetOrgParams{
		SalesforceID: orgID,
		Context:      context.Background(),
	}
	result, err := osc.cl.Organizations.GetOrg(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload, nil
}

// ListOrgUserAdminScopes returns admin role scope of organization
func (osc *Client) ListOrgUserAdminScopes(orgID string) (*models.UserrolescopesList, error) {
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params := &organizations.ListOrgUsrAdminScopesParams{
		SalesforceID: orgID,
		Context:      context.Background(),
	}
	result, err := osc.cl.Organizations.ListOrgUsrAdminScopes(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload, nil
}

// ListOrgUserScopes returns role scope of organization
// rolename is optional filter
func (osc *Client) ListOrgUserScopes(orgID string, rolename []string) (*models.UserrolescopesList, error) {
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	params := &organizations.ListOrgUsrServiceScopesParams{
		SalesforceID: orgID,
		Context:      context.Background(),
	}
	if len(rolename) != 0 {
		params.Rolename = rolename
	}
	result, err := osc.cl.Organizations.ListOrgUsrServiceScopes(params, clientAuth)
	if err != nil {
		return nil, err
	}
	return result.Payload, nil
}

// CreateOrg creates company based on name and website with additional data for required fields
func (osc *Client) CreateOrg(companyName string, companyWebsite string) (*models.Organization, error) {
	tok, err := token.GetToken()
	if err != nil {
		return nil, err
	}
	// use linux foundation logo as default
	linuxFoundation, err := osc.SearchOrganization("Linux Foundation")
	if err != nil {
		return nil, err
	}
	clientAuth := runtimeClient.BearerToken(tok)
	description := "No Description"
	companyType := "No Type"
	companySource := "No Source"
	industry := "No Industry"
	logoURL := linuxFoundation[0].LogoURL

	org := models.CreateOrg{
		Description: &description,
		Name:        &companyName,
		Website:     &companyWebsite,
		Industry:    &industry,
		Source:      &companySource,
		Type:        &companyType,
		LogoURL:     &logoURL,
	}

	params := &organizations.CreateOrgParams{
		Org:     &org,
		Context: context.Background(),
	}

	result, err := osc.cl.Organizations.CreateOrg(params, clientAuth)

	if err != nil {
		log.Warnf("Failed to create salesforce Company :%s , err: %+v ", companyName, err)
		return nil, err
	}

	return result.Payload, err
}
