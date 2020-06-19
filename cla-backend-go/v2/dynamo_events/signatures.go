package dynamo_events

import (
	"github.com/aws/aws-lambda-go/events"
)

type Signature struct {
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

// should be called when we modify signature
func (s *service) SignatureSignedEvent(event events.DynamoDBEventRecord) error {
	var newSignature, oldSignature Signature
	err := unmarshalStreamImage(event.Change.OldImage, &oldSignature)
	if err != nil {
		return err
	}
	err = unmarshalStreamImage(event.Change.NewImage, &newSignature)
	if err != nil {
		return err
	}
	// check if signature signed event is received
	if oldSignature.SignatureSigned == false && newSignature.SignatureSigned == true {

	}
	return nil
}
