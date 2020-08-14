// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package approval_list

import "github.com/aws/aws-sdk-go/service/dynamodb/expression"

// buildProjects builds the response model projection for a given query
func buildProjection() expression.ProjectionBuilder {
	// These are the columns we want returned
	return expression.NamesList(
		expression.Name("request_id"),
		expression.Name("request_status"),
		expression.Name("company_id"),
		expression.Name("company_name"),
		expression.Name("project_id"),
		expression.Name("project_name"),
		expression.Name("user_id"),
		expression.Name("user_emails"),
		expression.Name("user_name"),
		expression.Name("user_github_id"),
		expression.Name("user_github_username"),
		expression.Name("date_created"),
		expression.Name("date_modified"),
		expression.Name("version"),
	)
}
