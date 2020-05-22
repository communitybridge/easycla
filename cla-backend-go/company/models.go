// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package company

// Company data model
type Company struct {
	CompanyID         string   `dynamodbav:"company_id" json:"company_id"`
	CompanyName       string   `dynamodbav:"company_name" json:"company_name"`
	CompanyACL        []string `dynamodbav:"company_acl" json:"company_acl"`
	CompanyExternalID string   `dynamodbav:"company_external_id" json:"company_external_id"`
	Created           string   `dynamodbav:"date_created" json:"date_created"`
	Updated           string   `dynamodbav:"date_modified" json:"date_modified"`
	CompanyManagerID  string   `dynamodbav:"company_manager_id" json:"company_manager_id"`
	Version           string   `dynamodbav:"version" json:"version"`
}

// Invite data model
type Invite struct {
	CompanyInviteID    string `dynamodbav:"company_invite_id"`
	RequestedCompanyID string `dynamodbav:"requested_company_id"`
	UserID             string `dynamodbav:"user_id"`
	Status             string `dynamodbav:"status"`
	Created            string `dynamodbav:"date_created"`
	Updated            string `dynamodbav:"date_modified"`
}

// InviteModel data model
type InviteModel struct {
	CompanyInviteID    string `json:"company_invite_id"`
	RequestedCompanyID string `json:"requested_company_id"`
	CompanyName        string `json:"company_name"`
	UserID             string `json:"user_id"`
	UserName           string `json:"user_name"`
	UserEmail          string `json:"user_email"`
	Status             string `json:"status"`
	Created            string `json:"date_created"`
	Updated            string `json:"date_modified"`
}
