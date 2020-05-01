// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager_requests

// CLAManagerRequests data model
type CLAManagerRequests struct {
	Requests []CLAManagerRequest
}

// CLAManagerRequest data model
type CLAManagerRequest struct {
	RequestID         string `json:"request_id"`
	CompanyID         string `json:"company_id"`
	CompanyExternalID string `json:"company_external_id"`
	CompanyName       string `json:"company_name"`
	ProjectID         string `json:"project_id"`
	ProjectExternalID string `json:"project_external_id"`
	ProjectName       string `json:"project_name"`
	UserID            string `json:"user_id"`
	UserExternalID    string `json:"user_external_id"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	Status            string `json:"status"`
	Created           string `json:"date_created"`
	Updated           string `json:"date_modified"`
}
