// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cla_manager

import "github.com/aws/aws-sdk-go/service/dynamodb/expression"

// buildRequestProjection returns the database field projection for the table
func buildRequestProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("request_id"),
		expression.Name("company_id"),
		//expression.Name("company_external_id"),
		expression.Name("company_name"),
		expression.Name("project_id"),
		expression.Name("project_external_id"),
		expression.Name("project_name"),
		expression.Name("user_id"),
		//expression.Name("user_external_id"),
		expression.Name("user_name"),
		expression.Name("user_email"),
		expression.Name("status"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
	)
}
