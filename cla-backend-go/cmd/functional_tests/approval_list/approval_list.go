// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package approval_list

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"

	"github.com/communitybridge/easycla/cla-backend-go/cmd/functional_tests/test_models"
	"github.com/verdverm/frisby"
)

var (
	claManagerToken            string
	claProspectiveManagerToken string
	authHeaders                = map[string]string{
		"Authorization":   "Bearer " + claManagerToken,
		"Content-Type":    "application/json",
		"Accept-Encoding": "application/json",
		"X-ACL":           "ewogICAgImFsbG93ZWQiOnRydWUsCiAgICAiaXNBZG1pbiI6IGZhbHNlLAogICAgInJlc291cmNlIjoiY29tcGFueV9zaWduYXR1cmVzIiwKICAgICJzY29wZXMiOiBbCiAgICAgICAgeyAidHlwZSI6ICJwcm9qZWN0fG9yZ2FuaXphdGlvbiIsCiAgICAgICAgICAiaWQiOiAiYTA5NDEwMDAwMDVvdUpGQUFZfGV4dGVybmFsLWVlOTY1ZWEyLWNhODMtNDQ4Mi04YTFiLTk0YTQ2OGQ5ZGNmYSIsCiAgICAgICAgICAicm9sZSI6ICJjbGEtbWFuYWdlciIsCiAgICAgICAgICAiY29udGV4dCI6ICJzdGFmZiIKICAgICAgICB9LAogICAgICAgIHsgInR5cGUiOiAicHJvamVjdHxvcmdhbml6YXRpb24iLAogICAgICAgICAgImlkIjogImEwOTQxMDAwMDA1b3VKRkFBWXwwMDE0MTAwMDAwVGUxQ2FBQUoiLAogICAgICAgICAgInJvbGUiOiAiY2xhLW1hbmFnZXIiLAogICAgICAgICAgImNvbnRleHQiOiAic3RhZmYiCiAgICAgICAgfSwKICAgICAgICB7ICJ0eXBlIjogInByb2plY3R8b3JnYW5pemF0aW9uIiwKICAgICAgICAgICJpZCI6ICJhMDk0MTAwMDAwNW91SkZBQVl8MDAxNDEwMDAwMFRlMDJEQUFSIiwKICAgICAgICAgICJyb2xlIjogImNsYS1tYW5hZ2VyIiwKICAgICAgICAgICJjb250ZXh0IjogInN0YWZmIgogICAgICAgIH0sCiAgICAgICAgeyAidHlwZSI6ICJwcm9qZWN0fG9yZ2FuaXphdGlvbiIsCiAgICAgICAgICAiaWQiOiAiYTA5NDEwMDAwMDVvdUpGQUFZfDAwMTJNMDAwMDJEZG85T1FBUiIsCiAgICAgICAgICAicm9sZSI6ICJjbGEtbWFuYWdlciIsCiAgICAgICAgICAiY29udGV4dCI6ICJzdGFmZiIKICAgICAgICB9CiAgICBdCn0K",
		"X-Email":         "ddeal+cla+dev+functional+test+cla+manager+user@linuxfoundation.org",
		"X-USERNAME":      "cladevfunctionaltestclamanageruser",
	}
)

const (
	//claManagerCompanyID           = "ee965ea2-ca83-4482-8a1b-94a468d9dcfa"
	//claManagerProjectID           = "d5412846-5dda-4c58-8f62-4c111a3cd0d3"
	claProspectiveManagerLFID = "cladevfunctionaltestuser"
	//claProspectiveManagerUserName = "CLA Functional Test User Linux Foundation"
	claProspectiveManagerEmail = "ddeal+cla+dev+functional+test+user@linuxfoundation.org"
	claProjectSFID             = "a0941000005ouJFAAY"
	claManagerCompanySFID      = "external-ee965ea2-ca83-4482-8a1b-94a468d9dcfa"
	claGroupID                 = "d5412846-5dda-4c58-8f62-4c111a3cd0d3"
	testGitHubUsername         = "dealako"
	testDomainName             = "dealako.com"
	testGitHubOrg              = "deal-test-org"
)

// TestBehaviour data model
type TestBehaviour struct {
	apiURL           string
	serviceEndpoint  string
	auth0User1Config test_models.Auth0Config
	auth0User2Config test_models.Auth0Config
}

// NewTestBehaviour creates a new test behavior model
func NewTestBehaviour(apiURL string, auth0User1Config, auth0User2Config test_models.Auth0Config) *TestBehaviour {
	return &TestBehaviour{
		apiURL + "/v4",
		fmt.Sprintf("%s/signatures/project/%s/company/%s/clagroup/%s/approval-list", apiURL+"/v4", claProjectSFID, claManagerCompanySFID, claGroupID),
		auth0User1Config,
		auth0User2Config,
	}
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

// RunUpdateApprovalListNoAuth test
func (t *TestBehaviour) RunUpdateApprovalListNoAuth() {
	frisby.Create("CLA Approval List - Update Approval List - No Auth").
		Put(t.serviceEndpoint).
		Send().
		ExpectStatus(401)
}

// RunUpdateApprovalListUnauthorized test
func (t *TestBehaviour) RunUpdateApprovalListUnauthorized() {
	// ProspectiveManagerToken is not approved yet - shouldn't be able to approve requests
	frisby.Create(fmt.Sprintf("CLA Approval List - Update Approval List - Unauthorized - %s", claProspectiveManagerLFID)).
		Put(t.serviceEndpoint).
		SetHeaders(map[string]string{
			"Authorization":   "Bearer " + claProspectiveManagerToken,
			"Content-Type":    "application/json",
			"Accept-Encoding": "application/json",
		}).
		SetJson(map[string][]string{
			"AddEmailApprovalList": {claProspectiveManagerEmail},
		}).
		Send().
		ExpectStatus(401)
}

// RunUpdateApprovalListAddEmail test
func (t *TestBehaviour) RunUpdateApprovalListAddEmail() {
	frisby.Create("CLA Approval List - Update Approval List - Add Email").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListRemoveEmail() {
	frisby.Create("CLA Approval List - Update Approval List - Remove Email").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListAddDomain() {
	frisby.Create("CLA Approval List - Update Approval List - Add Domain").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListRemoveDomain() {
	frisby.Create("CLA Approval List - Update Approval List - Remove Domain").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListAddGitHubUsername() {
	frisby.Create("CLA Approval List - Update Approval List - Add GitHub Username").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListRemoveGitHubUsername() {
	frisby.Create("CLA Approval List - Update Approval List - Remove GitHub Username").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListAddGitHubOrg() {
	frisby.Create("CLA Approval List - Update Approval List - Add GitHub Org").
		Put(t.serviceEndpoint).
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
func (t *TestBehaviour) RunUpdateApprovalListRemoveGitHubOrg() {
	frisby.Create("CLA Approval List - Update Approval List - Remove GitHub Org").
		Put(t.serviceEndpoint).
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
	t.RunGetCLAManagerToken()
	t.RunGetCLAProspectiveManagerToken()

	// No Credentials/Auth Tests - these should return 401
	t.RunUpdateApprovalListNoAuth()

	// Shouldn't be allowed to Update the Approval List unless you are already a CLA Manager
	t.RunUpdateApprovalListUnauthorized()

	t.RunUpdateApprovalListAddEmail()
	t.RunUpdateApprovalListRemoveEmail()
	t.RunUpdateApprovalListAddDomain()
	t.RunUpdateApprovalListRemoveDomain()
	t.RunUpdateApprovalListAddGitHubUsername()
	t.RunUpdateApprovalListRemoveGitHubUsername()
	t.RunUpdateApprovalListAddGitHubOrg()
	t.RunUpdateApprovalListRemoveGitHubOrg()
}
