// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

// ItemSignature database model
type ItemSignature struct {
	SignatureID                   string   `json:"signature_id"`
	DateCreated                   string   `json:"date_created"`
	DateModified                  string   `json:"date_modified"`
	SignatureApproved             bool     `json:"signature_approved"`
	SignatureSigned               bool     `json:"signature_signed"`
	SignatureDocumentMajorVersion string   `json:"signature_document_major_version"`
	SignatureDocumentMinorVersion string   `json:"signature_document_minor_version"`
	SignatureReferenceID          string   `json:"signature_reference_id"`
	SignatureReferenceName        string   `json:"signature_reference_name"`
	SignatureReferenceNameLower   string   `json:"signature_reference_name_lower"`
	SignatureProjectID            string   `json:"signature_project_id"`
	SignatureReferenceType        string   `json:"signature_reference_type"`
	SignatureType                 string   `json:"signature_type"`
	SignatureUserCompanyID        string   `json:"signature_user_ccla_company_id"`
	EmailWhitelist                []string `json:"email_whitelist"`
	DomainWhitelist               []string `json:"domain_whitelist"`
	GitHubWhitelist               []string `json:"github_whitelist"`
	GitHubOrgWhitelist            []string `json:"github_org_whitelist"`
	SignatureACL                  []string `json:"signature_acl"`
}

// DBManagersModel is a database model for only the ACL/Manager column
type DBManagersModel struct {
	SignatureID  string   `json:"signature_id"`
	SignatureACL []string `json:"signature_acl"`
}
