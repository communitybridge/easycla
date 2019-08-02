// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package template

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
)

var (
	// ErrTemplateNotFound error
	ErrTemplateNotFound = errors.New("template not found")
)

// Repository interface functions
type Repository interface {
	GetTemplates() ([]models.Template, error)
	GetTemplate(templateID string) (models.Template, error)
	GetCLAGroup(claGroupID string) (CLAGroup, error)
	UpdateDynamoContractGroupTemplates(ctx context.Context, ContractGroupID string, template models.Template, pdfUrls models.TemplatePdfs) error
}

type repository struct {
	stage          string // The AWS stage (dev, staging, prod)
	dynamoDBClient *dynamodb.DynamoDB
}

// CLAGroup structure
type CLAGroup struct {
}

// DynamoProjectCorporateDocuments model
type DynamoProjectCorporateDocuments struct {
	DynamoProjectDocument []DynamoProjectDocument `json:":project_corporate_documents"`
}

// DynamoProjectIndividualDocuments model
type DynamoProjectIndividualDocuments struct {
	DynamoProjectDocument []DynamoProjectDocument `json:":project_individual_documents"`
}

// DynamoProjectDocument model
type DynamoProjectDocument struct {
	DocumentName            string        `json:"document_name"`
	DocumentFileID          string        `json:"document_file_id"`
	DocumentContentType     string        `json:"document_content_type"`
	DocumentMajorVersion    int           `json:"document_major_version"`
	DocumentMinorVersion    int           `json:"document_minor_version"`
	DocumentCreationDate    string        `json:"document_creation_date"`
	DocumentPreamble        string        `json:"document_preamble"`
	DocumentLegalEntityName string        `json:"document_legal_entity_name"`
	DocumentAuthorName      string        `json:"document_author_name"`
	DocumentS3URL           string        `json:"document_s3_url"`
	DocumentTabs            []DocumentTab `json:"document_tabs"`
}

// DocumentTab structure
type DocumentTab struct {
	// Note: these are arranged to optimize structure memory alignment - see the lint: maligned
	DocumentTabType                     string `json:"document_tab_type"`
	DocumentTabID                       string `json:"document_tab_id"`
	DocumentTabName                     string `json:"document_tab_name"`
	DocumentTabAnchorString             string `json:"document_tab_anchor_string"`
	DocumentTabPage                     int64  `json:"document_tab_page"`
	DocumentTabWidth                    int64  `json:"document_tab_width"`
	DocumentTabAnchorXOffset            int64  `json:"document_tab_anchor_x_offset"`
	DocumentTabAnchorYOffset            int64  `json:"document_tab_anchor_y_offset"`
	DocumentTabPositionX                int64  `json:"document_tab_position_x"`
	DocumentTabPositionY                int64  `json:"document_tab_position_y"`
	DocumentTabHeight                   int64  `json:"document_tab_height"`
	DocumentTabIsLocked                 bool   `json:"document_tab_is_locked"`
	DocumentTabIsRequired               bool   `json:"document_tab_is_required"`
	DocumentTabAnchorIgnoreIfNotPresent bool   `json:"document_tab_anchor_ignore_if_not_present"`
}

// NewRepository creates a new instance of the repository service
func NewRepository(awsSession *session.Session, stage string) repository {
	return repository{
		stage:          stage,
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

// GetTemplates returns a list containing all the template models
func (r repository) GetTemplates() ([]models.Template, error) {
	templates := []models.Template{}
	for _, template := range templateMap {
		// DEBUG
		// Only show the LF style template in dev for now
		if template.Name == "LF Style Template" && strings.ToLower(r.stage) != "dev" {
			log.Debugf("Skipping '%s' template since we are in stage: %s - only shown in 'dev'", template.Name, r.stage)
		} else {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

// GetTemplate returns the template based on the template ID
func (r repository) GetTemplate(templateID string) (models.Template, error) {
	template, ok := templateMap[templateID]
	if !ok {
		return models.Template{}, ErrTemplateNotFound
	}

	return template, nil
}

// GetCLAGroup This method belongs in the contractgroup package. We are leaving it here
// because it accesses DynamoDB, but the contractgroup repository is designed
// to connect to postgres
func (r repository) GetCLAGroup(claGroupID string) (CLAGroup, error) {
	tableName := fmt.Sprintf("cla-%s-projects", r.stage)

	_, err := r.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(claGroupID),
			},
		},
	})
	if err != nil {
		return CLAGroup{}, err
	}

	return CLAGroup{}, nil
}

// UpdateDynamoContractGroupTemplates updates the templates in the data store
func (r repository) UpdateDynamoContractGroupTemplates(ctx context.Context, ContractGroupID string, template models.Template, pdfUrls models.TemplatePdfs) error {
	tableName := fmt.Sprintf("cla-%s-projects", r.stage)

	// Map the fields to the dynamo model as the attribute names are different
	// Map Template Fields into DocumentTab
	cclaDocumentTabs := []DocumentTab{}

	for _, field := range template.CclaFields {
		dynamoTab := DocumentTab{
			DocumentTabType:                     field.FieldType,
			DocumentTabID:                       field.ID,
			DocumentTabPage:                     1,
			DocumentTabName:                     field.Name,
			DocumentTabWidth:                    field.Width,
			DocumentTabHeight:                   field.Height,
			DocumentTabIsLocked:                 field.IsEditable,
			DocumentTabAnchorString:             field.AnchorString,
			DocumentTabIsRequired:               field.IsOptional,
			DocumentTabAnchorIgnoreIfNotPresent: field.IsOptional,
			DocumentTabAnchorXOffset:            field.OffsetX,
			DocumentTabAnchorYOffset:            field.OffsetY,
			DocumentTabPositionX:                0,
			DocumentTabPositionY:                0,
		}
		cclaDocumentTabs = append(cclaDocumentTabs, dynamoTab)
	}

	currentTime := time.Now().Format(time.RFC3339)

	// Map CCLA Template to Document
	dynamoCorporateProjectDocument := DynamoProjectDocument{
		DocumentName:            template.Name,
		DocumentFileID:          template.ID,
		DocumentContentType:     "storage+pdf",
		DocumentMajorVersion:    2,
		DocumentMinorVersion:    0,
		DocumentCreationDate:    currentTime,
		DocumentPreamble:        template.Name,
		DocumentLegalEntityName: template.Name,
		DocumentAuthorName:      template.Name,
		DocumentS3URL:           pdfUrls.CorporatePDFURL,
		DocumentTabs:            cclaDocumentTabs,
	}

	// project_corporate_documents is a List type, and thus the item needs to be in a slice
	dynamoCorporateProjectDocuments := []DynamoProjectDocument{}
	dynamoCorporateProjectDocuments = append(dynamoCorporateProjectDocuments, dynamoCorporateProjectDocument)

	dynamoCorporateProject := DynamoProjectCorporateDocuments{
		DynamoProjectDocument: dynamoCorporateProjectDocuments,
	}

	// Marshal object into dynamodb attribute
	expr, err := dynamodbattribute.MarshalMap(dynamoCorporateProject)
	if err != nil {
		return err
	}

	// Find Contract Group to update the Templates on
	key := map[string]*dynamodb.AttributeValue{
		"project_id": {
			S: aws.String(ContractGroupID),
		},
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(tableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set project_corporate_documents =  list_append(project_corporate_documents, :project_corporate_documents)"),
	}

	_, err = r.dynamoDBClient.UpdateItem(input)
	if err != nil {
		return err
	}

	// Map ICLA Template Fields into DocumentTab
	iclaDocumentTabs := []DocumentTab{}

	for _, field := range template.IclaFields {
		dynamoTab := DocumentTab{
			DocumentTabType:                     field.FieldType,
			DocumentTabID:                       field.ID,
			DocumentTabPage:                     1,
			DocumentTabName:                     field.Name,
			DocumentTabWidth:                    field.Width,
			DocumentTabHeight:                   field.Height,
			DocumentTabIsLocked:                 field.IsEditable,
			DocumentTabIsRequired:               field.IsOptional,
			DocumentTabAnchorString:             field.AnchorString,
			DocumentTabAnchorIgnoreIfNotPresent: field.IsOptional,
			DocumentTabAnchorXOffset:            field.OffsetX,
			DocumentTabAnchorYOffset:            field.OffsetY,
			DocumentTabPositionX:                0,
			DocumentTabPositionY:                0,
		}
		iclaDocumentTabs = append(iclaDocumentTabs, dynamoTab)
	}

	currentTime = time.Now().Format(time.RFC3339)
	// Map Template to Document
	dynamoIndividualDocument := DynamoProjectDocument{
		DocumentName:            template.Name,
		DocumentFileID:          template.ID,
		DocumentContentType:     "storage+pdf",
		DocumentMajorVersion:    2,
		DocumentMinorVersion:    0,
		DocumentCreationDate:    currentTime,
		DocumentPreamble:        template.Name,
		DocumentLegalEntityName: template.Name,
		DocumentAuthorName:      template.Name,
		DocumentS3URL:           pdfUrls.IndividualPDFURL,
		DocumentTabs:            iclaDocumentTabs,
	}

	dynamoProjectIndividualDocuments := []DynamoProjectDocument{}
	dynamoProjectIndividualDocuments = append(dynamoProjectIndividualDocuments, dynamoIndividualDocument)

	dynamoIndividualProject := DynamoProjectIndividualDocuments{
		DynamoProjectDocument: dynamoProjectIndividualDocuments,
	}

	expr, err = dynamodbattribute.MarshalMap(dynamoIndividualProject)
	if err != nil {
		return err
	}

	input = &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(tableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set project_individual_documents =  list_append(project_individual_documents, :project_individual_documents)"),
	}

	_, err = r.dynamoDBClient.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

// templateMap contains a list of our template models
var templateMap = map[string]models.Template{
	"fb4cc144-a76c-4c17-8a52-c648f158fded": {
		ID:          "fb4cc144-a76c-4c17-8a52-c648f158fded",
		Name:        "Apache Style",
		Description: "For use of projects under the Apache style of CLA.",
		MetaFields: []*models.MetaField{
			{
				Name:             "Project Name",
				Description:      "Project's Full Name.",
				TemplateVariable: "PROJECT_NAME",
			},
			{
				Name:             "Project Entity Name",
				Description:      "The Full Entity Name of the Project.",
				TemplateVariable: "PROJECT_ENTITY_NAME",
			},
			{
				Name:             "Contact Email Address",
				Description:      "The E-Mail Address of the Person managing the CLA. ",
				TemplateVariable: "CONTACT_EMAIL",
			},
		},
		IclaFields: []*models.Field{
			{
				ID:           "full_name",
				Name:         "Full Name",
				AnchorString: "Full name:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      55,
				OffsetY:      -8,
			},
			{
				ID:           "public_name",
				Name:         "Public Name",
				AnchorString: "Public name:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      68,
				OffsetY:      -8,
			},
			{
				ID:           "mailing_address1",
				Name:         "Mailing Address",
				AnchorString: "Mailing Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        300,
				Height:       20,
				OffsetX:      95,
				OffsetY:      -7,
			},
			{
				ID:           "mailing_address2",
				Name:         "Mailing Address",
				AnchorString: "Mailing Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      0,
				OffsetY:      24,
			},
			{
				ID:           "mailing_address3",
				Name:         "Mailing Address",
				AnchorString: "Mailing Address:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      0,
				OffsetY:      47,
			},
			{
				ID:           "country",
				Name:         "Country",
				AnchorString: "Country:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        300,
				Height:       20,
				OffsetX:      50,
				OffsetY:      -7,
			},
			{
				ID:           "telephone",
				AnchorString: "Telephone:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        350,
				Height:       20,
				OffsetX:      70,
				OffsetY:      -7,
			},
			{
				ID:           "email",
				Name:         "Email",
				AnchorString: "E-Mail:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        320,
				Height:       20,
				OffsetX:      40,
				OffsetY:      -8,
			},
			{
				ID:           "sign",
				Name:         "Please Sign",
				AnchorString: "Please Sign:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      70,
				OffsetY:      -5,
			},
			{
				ID:           "date",
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      30,
				OffsetY:      -7,
			},
		},
		CclaFields: []*models.Field{
			{
				ID:           "corporation_name",
				Name:         "Corporation Name",
				AnchorString: "Corporation Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        355,
				Height:       20,
				OffsetX:      100,
				OffsetY:      -5,
			},
			{
				ID:           "corporation_address1",
				Name:         "Corporation Address1",
				AnchorString: "Corporation Address:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        230,
				Height:       20,
				OffsetX:      100,
				OffsetY:      -6,
			},
			{
				ID:           "corporation_address2",
				Name:         "Corporation Address2",
				AnchorString: "Corporation Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      0,
				OffsetY:      23,
			},
			{
				ID:           "corporation_address3",
				Name:         "Corporation Address3",
				AnchorString: "Corporation Address:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      0,
				OffsetY:      49,
			},
			{
				ID:           "point_of_contact",
				Name:         "Point of Contact",
				AnchorString: "Point of Contact:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        260,
				Height:       20,
				OffsetX:      95,
				OffsetY:      -8,
			},
			{
				ID:           "email",
				Name:         "Email",
				AnchorString: "Contact E-Mail:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        260,
				Height:       20,
				OffsetX:      95,
				OffsetY:      -7,
			},
			{
				ID:           "sign",
				Name:         "Please Sign",
				AnchorString: "Please sign:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      100,
				OffsetY:      -6,
			},
			{
				ID:           "date",
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      40,
				OffsetY:      -7,
			},
			{
				ID:           "title",
				Name:         "Title",
				AnchorString: "Title:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        320,
				Height:       20,
				OffsetX:      30,
				OffsetY:      -7,
			},
			{
				ID:           "cla_manager_name",
				Name:         "Initial CLA Manager Name",
				AnchorString: "CLA Manager Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        385,
				Height:       20,
				OffsetX:      130,
				OffsetY:      -7,
			},
			{
				ID:           "cla_manager_email",
				Name:         "Initial CLA Manager Email",
				AnchorString: "CLA Manager E-Mail:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        385,
				Height:       20,
				OffsetX:      130,
				OffsetY:      -7,
			},
		},
		IclaHTMLBody: `
		<html><body>
		<p>Project Name: {{ PROJECT_NAME }} </p>

		<p>Project Entity:	{{ PROJECT_ENTITY_NAME }} </p>

		<p>If emailing signed PDF, send to: {{ CONTACT_EMAIL }}</p> 

		<h3 style="text-align: center">Individual Contributor License Agreement (“Agreement”) v2.0</h3>

		<p>Thank you for your interest in the project specified above (the “Project”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Project must have a Contributor License Agreement (CLA) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of the Project and its users; it does not change your rights to use your own Contributions for any other purpose. </p>

		<p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Project or its third-party service providers, or email a PDF of the signed agreement to the email address specified above. Please read this document carefully before signing and keep a copy for your records.</p>

		<p>You accept and agree to the following terms and conditions for Your present and future Contributions submitted to the Project. In return, the Project shall not use Your Contributions in a way that is contrary to the public benefit or inconsistent with its charter at the time of the Contribution. Except for the license granted herein to the Project and recipients of software distributed by the Project, You reserve all right, title, and interest in and to Your Contributions.</p>

		<p>1. Definitions.</p>

		<p>“You” (or “Your”) shall mean the copyright owner or legal entity authorized by the copyright owner that is making this Agreement with the Project. For legal entities, the entity making a Contribution and all other entities that control, are controlled by, or are under common control with that entity are considered to be a single Contributor. For the purposes of this definition, “control” means (i) the power, direct or indirect, to cause the direction or management of such entity, whether by contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the outstanding shares, or (iii) beneficial ownership of such entity.</p>

		<p>“Contribution” shall mean the code, documentation or other original works of authorship, including any modifications or additions to an existing work, that is intentionally submitted by You to the Project for inclusion in, or documentation of, any of the products owned or managed by the Project (the “Work”). For the purposes of this definition, “submitted” means any form of electronic, verbal, or written communication sent to the Project or its representatives, including but not limited to communication on electronic mailing lists, source code control systems, and issue tracking systems that are managed by, or on behalf of, the Project for the purpose of discussing and improving the Work, but excluding communication that is conspicuously marked or otherwise designated in writing by You as “Not a Contribution.”</p>

		<p>2. Grant of Copyright License. Subject to the terms and conditions of this Agreement, You hereby grant to the Project and to recipients of software distributed by the Project a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Your Contributions and such derivative works.</p>

		<p>3. Grant of Patent License. Subject to the terms and conditions of this Agreement, You hereby grant to the Project and to recipients of software distributed by the Project a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable (except as stated in this section) patent license to make, have made, use, offer to sell, sell, import, and otherwise transfer the Work, where such license applies only to those patent claims licensable by You that are necessarily infringed by Your Contribution(s) alone or by combination of Your Contribution(s)  with the Work to which such Contribution(s) were submitted. If any entity institutes patent litigation against You or any other entity (including a cross-claim or counterclaim in a lawsuit) alleging that your Contribution, or the Work to which you have contributed, constitutes direct or contributory patent infringement, then any patent licenses granted to that entity under this Agreement for that Contribution or Work shall terminate as of the date such litigation is filed.</p>

		<p>4. You represent that you are legally entitled to grant the above license. If your employer(s) has rights to intellectual property that you create that includes your Contributions, you represent that you have received permission to make Contributions on behalf of that employer, that your employer has waived such rights for your Contributions to the Project, or that your employer has executed a separate Corporate CLA with the Project.</p>

		<p>5. You represent that each of Your Contributions is Your original creation (see section 7 for submissions on behalf of others). You represent that Your Contribution submissions include complete details of any third-party license or other restriction (including, but not limited to, related patents and trademarks) of which you are personally aware and which are associated with any part of Your Contributions.</p>

		<p>6. You are not expected to provide support for Your Contributions, except to the extent You desire to provide support. You may provide support for free, for a fee, or not at all. Unless required by applicable law or agreed to in writing, You provide Your Contributions on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.</p>

		<p>7. Should You wish to submit work that is not Your original creation, You may submit it to the Project separately from any Contribution, identifying the complete details of its source and of any license or other restriction (including, but not limited to, related patents, trademarks, and license agreements) of which you are personally aware, and conspicuously marking the work as “Submitted on behalf of a third-party: [named here]”.</p>

		<p>8. You agree to notify the Project of any facts or circumstances of which you become aware that would make these representations inaccurate in any respect.</p>

		<p style="page-break-after: always; text-align: center">[Please complete and sign on the next page.]</p>

		<p>Please sign: __________________________________ Date: _______________ </p>
		<p>Full name: ______________________________________________________</p>
		<p>Public name: ____________________________________________________</p>
		<p>Mailing Address:	_____________________________________</p>
		<p>____________________________________________________</p>
		<p>____________________________________________________</p>		
		<p>Country: ________________________________________</p>
		<p>E-Mail: _________________________________________</p>
		</body></html>
		`,
		CclaHTMLBody: `
		<html><body>
		<p>Project Name: {{ PROJECT_NAME }} </p>

		<p>Project Entity:	{{ PROJECT_ENTITY_NAME }} </p>

		<p>If emailing signed PDF, send to: {{ CONTACT_EMAIL }}</p> 

		<h3 style="text-align: center"> Software Grant and Corporate Contributor License Agreement (“Agreement”) v2.0 </h3>

		<p>Thank you for your interest in the project specified above (the “Project”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Project must have a Contributor License Agreement (CLA) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of the Project and its users; it does not change your rights to use your own Contributions for any other purpose. </p>

		<p>This version of the Agreement allows an entity (the “Corporation”) to submit Contributions to the Project, to authorize Contributions submitted by its designated employees to the Project, and to grant copyright and patent licenses thereto. </p> 

		<p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Project or its third-party service providers, or email a PDF of the signed agreement to the email address specified above. Please read this document carefully before signing and keep a copy for your records. </p>

		<p>You accept and agree to the following terms and conditions for Your present and future Contributions submitted to the Project. In return, the Project shall not use Your Contributions in a way that is contrary to the public benefit or inconsistent with its charter at the time of the Contribution. Except for the license granted herein to the Project and recipients of software distributed by the Project, You reserve all right, title, and interest in and to Your Contributions. </p>

		<p>1. Definitions. </p>

		<p>“You” (or “Your”) shall mean the copyright owner or legal entity authorized by the copyright owner that is making this Agreement with the Project. For legal entities, the entity making a Contribution and all other entities that control, are controlled by, or are under common control with that entity are considered to be a single Contributor. For the purposes of this definition, “control” means (i) the power, direct or indirect, to cause the direction or management of such entity, whether by contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the outstanding shares, or (iii) beneficial ownership of such entity.</p>

		<p>“Contribution” shall mean the code, documentation or other original works of authorship, including any modifications or additions to an existing work, that is intentionally submitted by You to the Project for inclusion in, or documentation of, any of the products owned or managed by the Project (the “Work”). For the purposes of this definition, “submitted” means any form of electronic, verbal, or written communication sent to the Project or its representatives, including but not limited to communication on electronic mailing lists, source code control systems, and issue tracking systems that are managed by, or on behalf of, the Project for the purpose of discussing and improving the Work, but excluding communication that is conspicuously marked or otherwise designated in writing by You as “Not a Contribution.” </p>

		<p>2. Grant of Copyright License. Subject to the terms and conditions of this Agreement, You hereby grant to the Project and to recipients of software distributed by the Project a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Your Contributions and such derivative works.</p>

		<p>3. Grant of Patent License. Subject to the terms and conditions of this Agreement, You hereby grant to the Project and to recipients of software distributed by the Project a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable (except as stated in this section) patent license to make, have made, use, offer to sell, sell, import, and otherwise transfer the Work, where such license applies only to those patent claims licensable by You that are necessarily infringed by Your Contribution(s) alone or by combination of Your Contribution(s) with the Work to which such Contribution(s) were submitted. If any entity institutes patent litigation against You or any other entity (including a cross-claim or counterclaim in a lawsuit) alleging that your Contribution, or the Work to which you have contributed, constitutes direct or contributory patent infringement, then any patent licenses granted to that entity under this Agreement for that Contribution or Work shall terminate as of the date such litigation is filed. </p>

		<p>4. You represent that You are legally entitled to grant the above license. You represent further that the employee of the Corporation designated as the Initial CLA Manager below (and each who is designated in a subsequent written modification to the list of CLA Managers) (each, a “CLA Manager”) is authorized to maintain (1) the list of employees of the Corporation who are authorized to submit Contributions on behalf of the Corporation, and (2) the list of CLA Managers; in each case, using the designated system for managing such lists (the “CLA Tool”).</p>

		<p>5. You represent that each of Your Contributions is Your original creation (see section 7 for submissions on behalf of others).</p>

		<p>6. You are not expected to provide support for Your Contributions, except to the extent You desire to provide support. You may provide support for free, for a fee, or not at all. Unless required by applicable law or agreed to in writing, You provide Your Contributions on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.</p>

		<p>7. Should You wish to submit work that is not Your original creation, You may submit it to the Project separately from any Contribution, identifying the complete details of its source and of any license or other restriction (including, but not limited to, related patents, trademarks, and license agreements) of which you are personally aware, and conspicuously marking the work as “Submitted on behalf of a third-party: [named here]”.</p>

		<p>8. It is your responsibility to use the CLA Tool when any change is required to the list of designated employees authorized to submit Contributions on behalf of the Corporation, or to the list of the CLA Managers. It is your responsibility to notify the Project when any change is required to the Corporation’s Point of Contact with the Project.</p>

		<p style="page-break-after: always; text-align: center">[Please complete and sign on the next page.]</p>

		<p>Please sign: __________________________________ Date: _______________ </p>
		<p>Title: _______________________________________________</p>
		<p>Corporation name:_____________________________________</p>
		<p>Corporation address:__________________________________</p>
		<p>___________________________________________________</p>
		<p>___________________________________________________</p>						
		<p>Point of Contact:_________________________________________</p>								
		<p>Point of Contact E-Mail:__________________________________</p>
		<p>Initial CLA Manager Name:_________________________________</p>
		<p>Initial CLA Manager E-Mail:_______________________________</p>
		</body></html>`,
	},
	"eb017ab5-0fa0-4230-9060-3239720c6494": {
		ID:          "eb017ab5-0fa0-4230-9060-3239720c6494",
		Name:        "LF Style Template",
		Description: "For use of projects under the LF Style of CLA.",
		MetaFields: []*models.MetaField{
			{
				Name:             "Project Name",
				Description:      "Project's Full Name.",
				TemplateVariable: "PROJECT_NAME",
			},
			{
				Name:             "Project Entity Name",
				Description:      "The Full Entity Name of the Project.",
				TemplateVariable: "PROJECT_ENTITY_NAME",
			},
			{
				Name:             "Contact Email Address",
				Description:      "The E-Mail Address of the Person managing the CLA. ",
				TemplateVariable: "CONTACT_EMAIL",
			},
		},
		IclaFields: []*models.Field{
			{
				ID:           "sign",
				Name:         "Signature",
				AnchorString: "Signature:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      70,
				OffsetY:      -5,
			},
			{
				ID:           "full_name",
				Name:         "Full Name",
				AnchorString: "Name:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      55,
				OffsetY:      -8,
			},
			{
				ID:           "date",
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      30,
				OffsetY:      -7,
			},
		},
		CclaFields: []*models.Field{
			{
				ID:           "cla_manager_name",
				Name:         "Initial CLA Manager Name",
				AnchorString: "Initial CLA Manager Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        385,
				Height:       20,
				OffsetX:      190,
				OffsetY:      -7,
			},
			{
				ID:           "cla_manager_email",
				Name:         "Initial CLA Manager Email",
				AnchorString: "Initial CLA Manager E-Mail:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        385,
				Height:       20,
				OffsetX:      190,
				OffsetY:      -7,
			},
			{
				ID:           "corporation_name",
				Name:         "Corporation Name",
				AnchorString: "Company Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   true,
				Width:        355,
				Height:       20,
				OffsetX:      120,
				OffsetY:      -5,
			},
			{
				ID:           "sign",
				Name:         "Please Sign",
				AnchorString: "Signature:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      120,
				OffsetY:      -6,
			},
			{
				ID:           "point_of_contact",
				Name:         "Point of Contact",
				AnchorString: "CLA Manager Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   true,
				Width:        320,
				Height:       20,
				OffsetX:      130,
				OffsetY:      -8,
			},
			{
				ID:           "title",
				Name:         "Title",
				AnchorString: "CLA Manager Title:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   true,
				Width:        320,
				Height:       20,
				OffsetX:      130,
				OffsetY:      -7,
			},
			{
				ID:           "date",
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      40,
				OffsetY:      -7,
			},
		},
		IclaHTMLBody: `
		<html>
		<body>
		<p>Individual Contributor License Agreement ("Agreement")</p>

		<p>
			Project Name: {{ PROJECT_NAME }}</br>
		    Project Entity:	{{ PROJECT_ENTITY_NAME }}</br>
		    If emailing signed PDF, send to: {{ CONTACT_EMAIL }}
		</p> 

		<p>Thank you for your interest in the {{ PROJECT_ENTITY_NAME }} (hereinafter the "Project"). In order to clarify the intellectual property licenses granted with Contributions from any corporate entity, the Project is required to have an Individual Contributor License Agreement (ICLA) on file that has been signed by each contributing individual. (For legal entities, please use the Corporate Contributor License Agreement (CCLA).)</p>
		<p>Each contributing individual ("You") must accept and agree that, for any Contribution (as defined below), You are bound by the licenses granted and representations made herein.</p> 
		<p>"Contribution" means any code, documentation or other original work of authorship that is submitted to the Project for inclusion in the Project by You or by another person authorized to make the submission on Your behalf.</p>
		<p>You accept and agree that all of Your present and future Contributions to the Project shall be:</p>
		<p>Submitted under a Developer's Certificate of Origin v. 1.1 (DCO); and Licensed under the BSD-3-Clause License.</p>

		<p>Signature: ________________________________________________</p>
		<p>Name: _____________________________________________________</p>
		<p>Date: _____________________________________________________</p>
		</body></html>
		`,
		CclaHTMLBody: `
		<html><body>
		<p>Corporate Contributor License Agreement ("Agreement")</p>

		<p>
			Project Name: {{ PROJECT_NAME }}</br>
		    Project Entity:	{{ PROJECT_ENTITY_NAME }}</br>
		    If emailing signed PDF, send to: {{ CONTACT_EMAIL }}
		</p> 

		<p>Thank you for your interest in the {{ PROJECT_ENTITY_NAME }} (hereinafter "Project"). In order to clarify the intellectual property licenses granted with Contributions from any corporate entity to the Project, the Project is required to have a Corporate Contributor License Agreement (CCLA) on file that has been signed by each contributing corporation.</p>
		<p>Each contributing corporation ("You") must accept and agree that, for any Contribution (as defined below), You and all other individuals and entities that control You, are controlled by You, or are under common control with You, are bound by the licenses granted and representations made as though all such individuals and entities are a single contributor. For the purposes of this definition, "control" means (i) the power, direct or indirect, to cause the direction or management of such entity, whether by contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the outstanding shares, or (iii) beneficial ownership of such entity.</p> 
		<p>"Contribution" means any code, documentation or other original work of authorship that is submitted to the Project for inclusion in the Project by Your employee or by any individual or entity that controls You, or is under common control with You or is controlled by You and authorized to make the submission on Your behalf. You accept and agree that all of Your present and future Contributions to the Project shall be:</p>
		<p>Submitted under a Developer's Certificate of Origin v. 1.1 (DCO); and Licensed under the BSD-3-Clause License.</p>
		<p>You agree that You shall be bound by the terms of the BSD-3-Clause License for all contributions made by You and Your employees. Your designated employees are those listed by Your CLA Manager(s) on the system of record for the Project. You agree to identify Your initial CLA Manager below and thereafter maintain current CLA Manager records in the Project’s system of record.</p>

		<p>Initial CLA Manager Name:   __________________________________ </p>
		<p>Initial CLA Manager E-mail: __________________________________ </p>
		<p>Corporate Signature</p>
		<p>Company Name: _____________________________________________</p>
		<p>Signature: ________________________________________________</p>
		<p>CLA Manager Name: _____________________________________________________</p>						
		<p>CLA Manager Title: ____________________________________________________</p>								
		<p>Date: _____________________________________________________</p>
		</body></html>`,
	},
}
