// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package metrics

// ItemSignature represent item of signature table
type ItemSignature struct {
	SignatureID            string   `json:"signature_id"`
	SignatureReferenceID   string   `json:"signature_reference_id"`
	SignatureReferenceName string   `json:"signature_reference_name"`
	SignatureACL           []string `json:"signature_acl"`
	SignatureUserCompanyID string   `json:"signature_user_ccla_company_id"`
	SignatureType          string   `json:"signature_type"`
	SignatureReferenceType string   `json:"signature_reference_type"`
	SignatureProjectID     string   `json:"signature_project_id"`
}

// ItemRepository represent item of repositories table
type ItemRepository struct {
	RepositoryProjectID string `json:"repository_project_id"`
}

// ItemCompany represent item of companies table
type ItemCompany struct {
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
}

// ItemGerritInstance represent item of gerrit instance table
type ItemGerritInstance struct {
	ProjectID string `json:"project_id"`
}

// ItemProject represent item of projects table
type ItemProject struct {
	ProjectID         string `json:"project_id"`
	ProjectExternalID string `json:"project_external_id"`
	ProjectName       string `json:"project_name"`
}

// ItemUser represent item of users table
type ItemUser struct {
	LfUsername string `json:"lf_username"`
}

// Metrics contain all metrics related to easycla
type Metrics struct {
	TotalCountMetrics       *TotalCountMetrics       `json:"total_metrics"`
	CompanyMetrics          *CompanyMetrics          `json:"company_metrics"`
	ProjectMetrics          *ProjectMetrics          `json:"project_metrics"`
	CompanyProjectMetrics   *CompanyProjectMetrics   `json:"company_project_metrics"`
	ClaManagersDistribution *ClaManagersDistribution `json:"cla_managers_distribution"`
	CalculatedAt            string                   `json:"calculated_at"`
}

// TotalCountMetrics contains all metrics related to total count
type TotalCountMetrics struct {
	CorporateContributorsCount        int64  `json:"corporate_contributors_count"`
	IndividualContributorsCount       int64  `json:"individual_contributors_count"`
	ClaManagersCount                  int64  `json:"cla_managers_count"`
	ContributorsCount                 int64  `json:"contributors_count"`
	ProjectsCount                     int64  `json:"projects_count"`
	GithubRepositoriesCount           int64  `json:"github_repositories_count"`
	GerritRepositoriesCount           int64  `json:"gerrit_repositories_count"`
	RepositoriesCount                 int64  `json:"repositories_count"`
	CompaniesCount                    int64  `json:"companies_count"`
	CompaniesProjectContributionCount int64  `json:"companies_project_contribution_count"`
	LfMembersCLACount                 int64  `json:"lf_members_cla_count"`
	NonLfMembersCLACount              int64  `json:"non_lf_members_cla_count"`
	CreatedAt                         string `json:"created_at"`

	corporateContributors        map[string]interface{}
	individualContributors       map[string]interface{}
	claManagers                  map[string]interface{}
	contributors                 map[string]interface{}
	companiesProjectContribution map[string]interface{}
}

// CompanyMetric contains all metrics related with particular company
type CompanyMetric struct {
	ID                         string `json:"id"`
	CompanyName                string `json:"company_name"`
	ProjectCount               int64  `json:"project_count"`
	CorporateContributorsCount int64  `json:"corporate_contributors_count"`
	ClaManagersCount           int64  `json:"cla_managers_count"`
	CreatedAt                  string `json:"created_at"`
	corporateContributors      map[string]interface{}
	claManagers                map[string]interface{}
}

// ProjectMetric contains all metrics related with particular project
type ProjectMetric struct {
	ID                          string `json:"id"`
	SalesforceID                string `json:"salesforce_id,omitempty"`
	CompaniesCount              int64  `json:"companies_count"`
	ClaManagersCount            int64  `json:"cla_managers_count"`
	CorporateContributorsCount  int64  `json:"corporate_contributors_count"`
	IndividualContributorsCount int64  `json:"individual_contributors_count"`
	TotalContributorsCount      int64  `json:"total_contributors_count"`
	RepositoriesCount           int64  `json:"repositories_count"`
	CreatedAt                   string `json:"created_at"`
	ExternalProjectID           string `json:"external_project_id"`
	ProjectName                 string `json:"project_name"`
	companies                   map[string]interface{}
	claManagers                 map[string]interface{}
	corporateContributors       map[string]interface{}
	individualContributors      map[string]interface{}
}

// ClaManagersDistribution tells distribution of number of cla mangers associated with company
type ClaManagersDistribution struct {
	OneClaManager        int64  `json:"one_cla_manager"`
	TwoClaManager        int64  `json:"two_cla_manager"`
	ThreeClaManager      int64  `json:"three_cla_manager"`
	FourOrMoreClaManager int64  `json:"four_or_more_cla_manager"`
	CreatedAt            string `json:"created_at"`
}
