// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"encoding/json"
	"reflect"

	"github.com/communitybridge/easycla/cla-backend-go/cmd/functional_tests/test_models"
	"github.com/verdverm/frisby"
)

var (
	token string
)

// TestBehaviour data model
type TestBehaviour struct {
	apiURL      string
	auth0Config test_models.Auth0Config
}

// NewTestBehaviour creates a new test behavior model
func NewTestBehaviour(apiURL string, auth0Config test_models.Auth0Config) *TestBehaviour {
	return &TestBehaviour{
		apiURL,
		auth0Config,
	}
}

// RunGetToken acquires the Auth0 token
func (t *TestBehaviour) RunGetToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0Config.Auth0UserName,
		"password":   t.auth0Config.Auth0Password,
		"client_id":  t.auth0Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create("CLA Group - Get Token").
		Post("https://linuxfoundation-dev.auth0.com/oauth/token").
		SetJson(authTokenReqPayload).
		Send().
		ExpectStatus(200).
		ExpectJsonType("access_token", reflect.String).
		ExpectJsonType("id_token", reflect.String).
		ExpectJsonType("scope", reflect.String).
		ExpectJsonType("expires_in", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			//log.Debugf("JSON: %+v", text)
			var auth0Response test_models.Auth0Response
			unmarshallErr := json.Unmarshal([]byte(text), &auth0Response)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &auth0Response == nil {
				F.AddError("Auth0Response is nil")
			}
			if auth0Response.IDToken == "" {
				F.AddError("Auth0Response id_token is empty")
			}
			token = auth0Response.IDToken
			//log.Debugf("ID Token is: %s", token)
		})
}

// RunCreateCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunCreateCLAManagerRequestNoAuth() {
	frisby.Create("CLA Group - Create CLA Manager Request - No Auth").
		Post(t.apiURL + "/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests").
		Send().
		ExpectStatus(401)
}

// RunGetCLAManagerRequestsNoAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestsNoAuth() {
	frisby.Create("CLA Group - Get CLA Manager Requests - No Auth").
		Get(t.apiURL + "/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests").
		Send().
		ExpectStatus(401)
}

// RunGetCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestNoAuth() {
	frisby.Create("CLA Group - Get CLA Manager Request - No Auth").
		Get(t.apiURL + "/v3/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests/4e0ed114-9c68-496f-800e-fa302f1c3edc").
		Send().
		ExpectStatus(401)
}

// RunApproveCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunApproveCLAManagerRequestNoAuth() {
	frisby.Create("CLA Group - Approve CLA Manager Request - No Auth").
		Post(t.apiURL + "/v3/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests/4e0ed114-9c68-496f-800e-fa302f1c3edc/approve").
		Send().
		ExpectStatus(401)
}

// RunDenyCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunDenyCLAManagerRequestNoAuth() {
	frisby.Create("CLA Group - Deny CLA Manager Request - No Auth").
		Post(t.apiURL + "/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests/4e0ed114-9c68-496f-800e-fa302f1c3edc/deny").
		Send().
		ExpectStatus(401)
}

// RunGetCLAManagerRequestsAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestsAuth() {
	frisby.Create("CLA Group - Get CLA Manager Requests - Auth").
		Get(t.apiURL+"/company/ee965ea2-ca83-4482-8a1b-94a468d9dcfa/project/d5412846-5dda-4c58-8f62-4c111a3cd0d3/cla-manager/requests").
		SetHeader("Authorization", "Bearer "+token).
		Send().
		ExpectStatus(200)
}

// RunAllTests runs all the CLA Group tests
func (t *TestBehaviour) RunAllTests() {
	t.RunGetToken()
	t.RunCreateCLAManagerRequestNoAuth()
	t.RunGetCLAManagerRequestsNoAuth()
	t.RunGetCLAManagerRequestNoAuth()
	t.RunApproveCLAManagerRequestNoAuth()
	t.RunDenyCLAManagerRequestNoAuth()
	t.RunGetCLAManagerRequestsAuth()
}
