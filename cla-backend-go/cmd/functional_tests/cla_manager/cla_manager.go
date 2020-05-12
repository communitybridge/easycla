// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/cmd/functional_tests/test_models"
	"github.com/verdverm/frisby"
)

var (
	claManagerToken            string
	claProspectiveManagerToken string
	claManagerCreateRequestID  string
)

const (
	claManagerCompanyID = "ee965ea2-ca83-4482-8a1b-94a468d9dcfa"
	claManagerProjectID = "d5412846-5dda-4c58-8f62-4c111a3cd0d3"
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

// RunGetCLAManagerToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAManagerToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0Config.Auth0UserName, // TODO - use diff creds for this person
		"password":   t.auth0Config.Auth0Password,
		"client_id":  t.auth0Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create("CLA Manager - Get Token").
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
			claManagerToken = auth0Response.IDToken
			//log.Debugf("ID Token is: %s", token)
		})
}

// RunGetCLAProspectiveManagerToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAProspectiveManagerToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0Config.Auth0UserName, // TODO - use diff creds for this person
		"password":   t.auth0Config.Auth0Password,
		"client_id":  t.auth0Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create("CLA Manager - Get Token").
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
			claProspectiveManagerToken = auth0Response.IDToken
			//log.Debugf("ID Token is: %s", token)
		})
}

// RunCreateCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunCreateCLAManagerRequestNoAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests",
		t.apiURL, claManagerCompanyID, claManagerProjectID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Create CLA Manager Request - No Auth").
		Post(url).
		Send().
		ExpectStatus(401)
}

// RunGetCLAManagerRequestsNoAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestsNoAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests",
		t.apiURL, claManagerCompanyID, claManagerProjectID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Get CLA Manager Requests - No Auth").
		Get(url).
		Send().
		ExpectStatus(401)
}

// RunGetCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestNoAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s",
		t.apiURL, claManagerCompanyID, claManagerProjectID, "test-request-id")
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Get CLA Manager Request - No Auth").
		Get(url).
		Send().
		ExpectStatus(401)
}

// RunApproveCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunApproveCLAManagerRequestNoAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s/approve",
		t.apiURL, claManagerCompanyID, claManagerProjectID, "test-request-id")
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Approve CLA Manager Request - No Auth").
		Put(url).
		Send().
		ExpectStatus(401)
}

// RunDenyCLAManagerRequestNoAuth test
func (t *TestBehaviour) RunDenyCLAManagerRequestNoAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s/deny",
		t.apiURL, claManagerCompanyID, claManagerProjectID, "test-request-id")
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Deny CLA Manager Request - No Auth").
		Put(url).
		Send().
		ExpectStatus(401)
}

// RunCreateCLAManagerRequestAuth test
func (t *TestBehaviour) RunCreateCLAManagerRequestAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests",
		t.apiURL, claManagerCompanyID, claManagerProjectID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Create CLA Manager Request - Auth").
		Post(url).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claProspectiveManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		SetJson(map[string]string{
			"userName":  "Deal Test User",
			"userEmail": "ddeal+test@linuxfoundation.org",
			"userLFID":  "ddealtest",
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyID", reflect.String).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("projectName", reflect.String).
		ExpectJsonType("projectExternalID", reflect.String).
		ExpectJsonType("userID", reflect.String).
		ExpectJsonType("userName", reflect.String).
		ExpectJsonType("userEmail", reflect.String).
		ExpectJsonType("created", reflect.String).
		ExpectJsonType("updated", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			//log.Debugf("Create CLA Manager Response JSON: %+v", text)
			var claManagerResp models.ClaManagerRequest
			unmarshallErr := json.Unmarshal([]byte(text), &claManagerResp)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &claManagerResp == nil {
				F.AddError("CLA Manager Response is nil")
			}
			claManagerCreateRequestID = claManagerResp.RequestID
			log.Debugf("Saved CLA Manager request ID: %s", claManagerCreateRequestID)
		})
}

// RunGetCLAManagerRequestsAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestsAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests",
		t.apiURL, claManagerCompanyID, claManagerProjectID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Get CLA Manager Requests - Auth").
		Get(url).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		Send().
		ExpectStatus(200)
}

// RunGetCLAManagerRequestAuth test
func (t *TestBehaviour) RunGetCLAManagerRequestAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s",
		t.apiURL, claManagerCompanyID, claManagerProjectID, claManagerCreateRequestID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Get CLA Manager Request - Auth").
		Get(url).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyID", reflect.String).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("projectName", reflect.String).
		ExpectJsonType("projectExternalID", reflect.String).
		ExpectJsonType("userID", reflect.String).
		ExpectJsonType("userName", reflect.String).
		ExpectJsonType("userEmail", reflect.String).
		ExpectJsonType("created", reflect.String).
		ExpectJsonType("updated", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			//log.Debugf("Create CLA Manager Response JSON: %+v", text)
			var claManagerResp models.ClaManagerRequest
			unmarshallErr := json.Unmarshal([]byte(text), &claManagerResp)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &claManagerResp == nil {
				F.AddError("CLA Manager Response is nil")
			}
			// Company ID's Match
			F.Expect(func(F *frisby.Frisby) (bool, string) {
				return claManagerCompanyID == claManagerResp.CompanyID, fmt.Sprintf("Company IDs Match")
			})

			// Project ID's Match
			F.Expect(func(F *frisby.Frisby) (bool, string) {
				return claManagerProjectID == claManagerResp.ProjectID, fmt.Sprintf("Project IDs Match")
			})

			// Request ID's Match
			F.Expect(func(F *frisby.Frisby) (bool, string) {
				return claManagerCreateRequestID == claManagerResp.RequestID, fmt.Sprintf("Request IDs Match")
			})
		})
}

// RunDeleteCLAManagerRequestUnauthorized test
func (t *TestBehaviour) RunDeleteCLAManagerRequestUnauthorized() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s",
		t.apiURL, claManagerCompanyID, claManagerProjectID, claManagerCreateRequestID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Delete CLA Manager Requests - Auth").
		Delete(url).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		Send().
		ExpectStatus(401)
}

// RunDeleteCLAManagerRequestAuth test
func (t *TestBehaviour) RunDeleteCLAManagerRequestAuth() {
	url := fmt.Sprintf("%s/company/%s/project/%s/cla-manager/requests/%s",
		t.apiURL, claManagerCompanyID, claManagerProjectID, claManagerCreateRequestID)
	//log.Debugf("URL: %s", url)
	frisby.Create("CLA Manager - Delete CLA Manager Requests - Auth").
		Delete(url).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		Send().
		ExpectStatus(200)
}

// RunAllTests runs all the CLA Manager tests
func (t *TestBehaviour) RunAllTests() {
	t.RunGetCLAManagerToken()
	t.RunGetCLAProspectiveManagerToken()
	t.RunCreateCLAManagerRequestNoAuth()
	t.RunGetCLAManagerRequestsNoAuth()
	t.RunGetCLAManagerRequestNoAuth()
	t.RunApproveCLAManagerRequestNoAuth()
	t.RunDenyCLAManagerRequestNoAuth()
	t.RunCreateCLAManagerRequestAuth()
	t.RunGetCLAManagerRequestsAuth()
	t.RunGetCLAManagerRequestAuth()
	t.RunDeleteCLAManagerRequestUnauthorized()
	t.RunDeleteCLAManagerRequestAuth()
}
