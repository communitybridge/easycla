// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package approval_list

import (
	"encoding/json"
	"fmt"
	"reflect"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"

	"github.com/communitybridge/easycla/cla-backend-go/cmd/functional_tests/test_models"
	"github.com/verdverm/frisby"
)

var (
	claProspectiveManagerToken string
	claManagerToken            string
	claManagerIntelToken       string
	claManagerATTToken         string
)

const (
	//claManagerCompanyID           = "ee965ea2-ca83-4482-8a1b-94a468d9dcfa"
	//claManagerProjectID           = "d5412846-5dda-4c58-8f62-4c111a3cd0d3"
	claProspectiveManagerLFID = "cladevfunctionaltestuser"
	//claProspectiveManagerUserName = "CLA Functional Test User Linux Foundation"
	claProspectiveManagerEmail = "ddeal+cla+dev+functional+test+user@linuxfoundation.org"

	projectSFID = "a092M00001If9v8QAB"
	intelSFID   = "00117000015vpjXAAQ"
	attSFID     = "0014100000Te1CaAAJ"
	claGroupID  = "0e011a1a-a67d-498a-a698-df247481dbb6"

	// test data
	testGitHubUsername = "dealako"
	testDomainName     = "dealako.com"
	testGitHubOrg      = "deal-test-org"
	// Project: Deal Project (fake): d5412846-5dda-4c58-8f62-4c111a3cd0d3 : a0941000005ouJFAAY
	// Company: Deal Gateway : 4c96ad67-f43f-4eee-a462-f79816c2c3f2 : 0014100001b1vOqAAI
)

// TestBehaviour data model
type TestBehaviour struct {
	apiURL           string
	auth0User1Config test_models.Auth0Config
	auth0User2Config test_models.Auth0Config
	auth0User3Config test_models.Auth0Config
	auth0User4Config test_models.Auth0Config
}

// NewTestBehaviour creates a new test behavior model
func NewTestBehaviour(apiURL string, auth0User1Config, auth0User2Config, auth0User3Config, auth0User4Config test_models.Auth0Config) *TestBehaviour {
	return &TestBehaviour{
		apiURL,
		auth0User1Config,
		auth0User2Config,
		auth0User3Config,
		auth0User4Config,
	}
}

func (t *TestBehaviour) getProspectiveCLAManagerHeaders() map[string]string {
	return map[string]string{
		"Authorization":   "Bearer " + claProspectiveManagerToken,
		"Content-Type":    "application/json",
		"Accept-Encoding": "application/json",
		"X-Email":         t.auth0User2Config.Auth0Email,
		"X-USERNAME":      t.auth0User2Config.Auth0UserName,
	}
}

func (t *TestBehaviour) getCLAManagerHeaders() map[string]string {
	return map[string]string{
		"Authorization":   "Bearer " + claManagerToken,
		"Content-Type":    "application/json",
		"Accept-Encoding": "application/json",
		"X-Email":         t.auth0User1Config.Auth0Email,
		"X-USERNAME":      t.auth0User1Config.Auth0UserName,
	}
}

func (t *TestBehaviour) getURLApprovalListColorIOIntel() string {
	return fmt.Sprintf("%s/v4/signatures/project/%s/company/%s/clagroup/%s/approval-list",
		t.apiURL, projectSFID, intelSFID, claGroupID)
}

func (t *TestBehaviour) getURLApprovalListColorIOATT() string {
	return fmt.Sprintf("%s/v4/signatures/project/%s/company/%s/clagroup/%s/approval-list",
		t.apiURL, projectSFID, attSFID, claGroupID)
}

func (t *TestBehaviour) getCLAManagerIntelHeaders() map[string]string {
	return map[string]string{
		"Authorization":   "Bearer " + claManagerIntelToken,
		"Content-Type":    "application/json",
		"Accept-Encoding": "application/json",
		"X-Email":         t.auth0User3Config.Auth0Email,
		"X-USERNAME":      t.auth0User3Config.Auth0UserName,
	}
}

func (t *TestBehaviour) getCLAManagerATTHeaders() map[string]string {
	return map[string]string{
		"Authorization":   "Bearer " + claManagerATTToken,
		"Content-Type":    "application/json",
		"Accept-Encoding": "application/json",
		"X-Email":         t.auth0User4Config.Auth0Email,
		"X-USERNAME":      t.auth0User4Config.Auth0UserName,
	}
}

// RunGetCLAProspectiveManagerToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAProspectiveManagerToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0User1Config.Auth0UserName,
		"password":   t.auth0User1Config.Auth0Password,
		"client_id":  t.auth0User1Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}

	frisby.Create(fmt.Sprintf("CLA Manager - Get Token - Prospective CLA Manager - %s", t.auth0User1Config.Auth0UserName)).
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

// RunGetCLAManagerToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAManagerToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0User2Config.Auth0UserName,
		"password":   t.auth0User2Config.Auth0Password,
		"client_id":  t.auth0User2Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create(fmt.Sprintf("CLA Manager - Get Token - CLA Manager - %s", t.auth0User2Config.Auth0UserName)).
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

// RunGetCLAManagerIntelToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAManagerIntelToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0User3Config.Auth0UserName,
		"password":   t.auth0User3Config.Auth0Password,
		"client_id":  t.auth0User3Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create(fmt.Sprintf("CLA Manager - Get Token - CLA Manager Intel - %s", t.auth0User3Config.Auth0UserName)).
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
			claManagerIntelToken = auth0Response.IDToken
			//log.Debugf("CLA Manager Intel Token is: %s", claManagerIntelToken)
		})
}

// RunGetCLAManagerATTToken acquires the Auth0 token
func (t *TestBehaviour) RunGetCLAManagerATTToken() {
	authTokenReqPayload := map[string]string{
		"grant_type": "http://auth0.com/oauth/grant-type/password-realm",
		"realm":      "Username-Password-Authentication",
		"username":   t.auth0User4Config.Auth0UserName,
		"password":   t.auth0User4Config.Auth0Password,
		"client_id":  t.auth0User4Config.Auth0ClientID,
		"audience":   "https://api-gw.dev.platform.linuxfoundation.org/",
		"scope":      "access:api openid profile email",
	}
	frisby.Create(fmt.Sprintf("CLA Manager - Get Token - CLA Manager AT&T - %s", t.auth0User4Config.Auth0UserName)).
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
			claManagerATTToken = auth0Response.IDToken
			//log.Debugf("CLA Manager ATT Token is: %s", claManagerATTToken)
		})
}

// RunUpdateApprovalListNoAuth test
func (t *TestBehaviour) RunUpdateApprovalListNoAuth(endpoint string) {
	frisby.Create("CLA Approval List - Update Approval List - No Auth").
		Put(endpoint).
		Send().
		ExpectStatus(401)
}

// RunUpdateApprovalListUnauthorized test
func (t *TestBehaviour) RunUpdateApprovalListUnauthorized(endpoint string, authHeaders map[string]string) {
	// ProspectiveManagerToken is not approved yet - shouldn't be able to approve requests
	frisby.Create(fmt.Sprintf("CLA Approval List - Update Approval List - Unauthorized - %s", claProspectiveManagerLFID)).
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"AddEmailApprovalList": {claProspectiveManagerEmail},
		}).
		Send().
		ExpectStatus(403)
}

// RunUpdateApprovalListAddEmail test
func (t *TestBehaviour) RunUpdateApprovalListAddEmail(endpoint string, authHeaders map[string]string) {
	log.Debugf("CLA Approval List - Update Approval List - Add Email - URL: %s", endpoint)
	frisby.Create(fmt.Sprintf("CLA Approval List - Update Approval List - Add Email - %s", authHeaders["X-USERNAME"])).
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"AddEmailApprovalList": {claProspectiveManagerEmail},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if !listContains(sig.EmailApprovalList, claProspectiveManagerEmail) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - Email should contain: %s",
					claProspectiveManagerEmail))
			}
		})
}

// RunUpdateApprovalListRemoveEmail test
func (t *TestBehaviour) RunUpdateApprovalListRemoveEmail(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Remove Email").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"RemoveEmailApprovalList": {claProspectiveManagerEmail},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if listContains(sig.EmailApprovalList, claProspectiveManagerEmail) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - Email should not contain: %s",
					claProspectiveManagerEmail))
			}
		})
}

// RunUpdateApprovalListAddDomain test
func (t *TestBehaviour) RunUpdateApprovalListAddDomain(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Add Domain").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"AddDomainApprovalList": {testDomainName},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if !listContains(sig.DomainApprovalList, testDomainName) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - Domain not contain: %s",
					testDomainName))
			}
		})
}

// RunUpdateApprovalListRemoveDomain test
func (t *TestBehaviour) RunUpdateApprovalListRemoveDomain(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Remove Domain").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"RemoveDomainApprovalList": {testDomainName},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if listContains(sig.DomainApprovalList, testDomainName) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - Domain should not contain: %s",
					testDomainName))
			}
		})
}

// RunUpdateApprovalListAddGitHubUsername test
func (t *TestBehaviour) RunUpdateApprovalListAddGitHubUsername(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Add GitHub Username").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"AddGithubUsernameApprovalList": {testGitHubUsername},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if !listContains(sig.GithubUsernameApprovalList, testGitHubUsername) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - GH Username should contain: %s",
					testGitHubUsername))
			}
		})
}

// RunUpdateApprovalListRemoveGitHubUsername test
func (t *TestBehaviour) RunUpdateApprovalListRemoveGitHubUsername(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Remove GitHub Username").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"RemoveGithubUsernameApprovalList": {testGitHubUsername},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if listContains(sig.GithubUsernameApprovalList, testGitHubUsername) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - GH Username should not contain: %s",
					testGitHubUsername))
			}
		})
}

// RunUpdateApprovalListAddGitHubOrg test
func (t *TestBehaviour) RunUpdateApprovalListAddGitHubOrg(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Add GitHub Org").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"AddGithubOrgApprovalList": {testGitHubOrg},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if !listContains(sig.GithubOrgApprovalList, testGitHubOrg) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - GH Org should contain: %s",
					testGitHubOrg))
			}
		})
}

// RunUpdateApprovalListRemoveGitHubOrg test
func (t *TestBehaviour) RunUpdateApprovalListRemoveGitHubOrg(endpoint string, authHeaders map[string]string) {
	frisby.Create("CLA Approval List - Update Approval List - Remove GitHub Org").
		Put(endpoint).
		SetHeaders(authHeaders).
		SetJson(map[string][]string{
			"RemoveGithubOrgApprovalList": {testGitHubOrg},
		}).
		Send().
		ExpectStatus(200).
		ExpectJsonType("companyName", reflect.String).
		ExpectJsonType("projectID", reflect.String).
		ExpectJsonType("signatureID", reflect.String).
		ExpectJsonType("signatureMajorVersion", reflect.String).
		ExpectJsonType("signatureMinorVersion", reflect.String).
		ExpectJsonType("signatureReferenceID", reflect.String).
		ExpectJsonType("signatureReferenceName", reflect.String).
		ExpectJsonType("signatureReferenceNameLower", reflect.String).
		ExpectJsonType("signatureReferenceType", reflect.String).
		ExpectJsonType("signatureType", reflect.String).
		ExpectJsonType("signatureApproved", reflect.Bool).
		ExpectJsonType("signatureSigned", reflect.Bool).
		ExpectJsonType("signatureCreated", reflect.String).
		ExpectJsonType("signatureModified", reflect.String).
		ExpectJsonType("version", reflect.String).
		AfterText(func(F *frisby.Frisby, text string, err error) {
			var sig models.Signature
			unmarshallErr := json.Unmarshal([]byte(text), &sig)
			if unmarshallErr != nil {
				F.AddError(unmarshallErr.Error())
			}
			if &sig == nil {
				F.AddError("CLA Approval List - Update Approval List - Response is nil")
			}
			if sig.ProjectID != claGroupID {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - CLA Group ID's do not match: %s vs %s",
					sig.ProjectID, claGroupID))
			}
			if listContains(sig.GithubOrgApprovalList, testGitHubOrg) {
				F.AddError(fmt.Sprintf("CLA Approval List - Update Approval List - GH Org should not contain: %s",
					testGitHubOrg))
			}
		})
}

func listContains(list []string, str string) bool {
	retVal := false
	for _, s := range list {
		if s == str {
			retVal = true
			break
		}
	}

	return retVal
}

// RunAllTests runs all the CLA Manager tests
func (t *TestBehaviour) RunAllTests() {
	// Need our authentication tokens for each persona/user
	t.RunGetCLAProspectiveManagerToken()
	t.RunGetCLAManagerToken()
	t.RunGetCLAManagerIntelToken()
	t.RunGetCLAManagerATTToken()

	// No Credentials/Auth Tests - these should return 401
	t.RunUpdateApprovalListNoAuth(t.getURLApprovalListColorIOIntel())

	// Shouldn't be allowed to Update the Approval List unless you are already a CLA Manager
	t.RunUpdateApprovalListUnauthorized(t.getURLApprovalListColorIOIntel(), t.getProspectiveCLAManagerHeaders())

	//t.RunUpdateApprovalListAddEmail(t.getCLAManagerHeaders())
	t.RunUpdateApprovalListAddEmail(t.getURLApprovalListColorIOIntel(), t.getCLAManagerIntelHeaders())
	t.RunUpdateApprovalListAddEmail(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
	if false {
		t.RunUpdateApprovalListRemoveEmail(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListAddDomain(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListRemoveDomain(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListAddGitHubUsername(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListRemoveGitHubUsername(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListAddGitHubOrg(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
		t.RunUpdateApprovalListRemoveGitHubOrg(t.getURLApprovalListColorIOATT(), t.getCLAManagerATTHeaders())
	}
	t.getCLAManagerHeaders()
}
