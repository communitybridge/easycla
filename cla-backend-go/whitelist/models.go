// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package whitelist

// CLARequestModel data model
type CLARequestModel struct {
	RequestID          string   `dynamodbav:"request_id"`
	RequestStatus      string   `dynamodbav:"request_status"`
	CompanyID          string   `dynamodbav:"company_id"`
	CompanyName        string   `dynamodbav:"company_name"`
	ProjectID          string   `dynamodbav:"project_id"`
	ProjectName        string   `dynamodbav:"project_name"`
	UserID             string   `dynamodbav:"user_id"`
	UserEmails         []string `dynamodbav:"user_emails"`
	UserName           string   `dynamodbav:"user_name"`
	UserGithubID       string   `dynamodbav:"user_github_id"`
	UserGithubUsername string   `dynamodbav:"user_github_username"`
	DateCreated        string   `dynamodbav:"date_created"`
	DateModified       string   `dynamodbav:"date_modified"`
	Version            string   `dynamodbav:"version"`
}
