// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

// DBProjectModel data model
type DBProjectModel struct {
	DateCreated                      string   `dynamodbav:"date_created"`
	DateModified                     string   `dynamodbav:"date_modified"`
	ProjectExternalID                string   `dynamodbav:"project_external_id"`
	ProjectID                        string   `dynamodbav:"project_id"`
	ProjectName                      string   `dynamodbav:"project_name"`
	Version                          string   `dynamodbav:"version"`
	ProjectCclaEnabled               bool     `dynamodbav:"project_ccla_enabled"`
	ProjectCclaRequiresIclaSignature bool     `dynamodbav:"project_ccla_requires_icla_signature"`
	ProjectIclaEnabled               bool     `dynamodbav:"project_icla_enabled"`
	ProjectACL                       []string `dynamodbav:"project_acl"`
}
