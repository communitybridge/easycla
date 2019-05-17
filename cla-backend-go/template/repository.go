package template

import (
	"context"
	"errors"
	"fmt"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

type Repository interface {
	GetTemplates() ([]models.Template, error)
	GetTemplate(templateID string) (models.Template, error)
	GetCLAGroup(claGroupID string) (CLAGroup, error)
	UpdateDynamoContractGroupTemplates(ctx context.Context, ContractGroupID string, tableName string, template models.Template) error
}

type repository struct {
	dynamoDBClient *dynamodb.DynamoDB
}

type CLAGroup struct {
}

type DynamoProjectCorporateDocuments struct {
	DynamoProjectDocument []DynamoProjectDocument `json:":project_corporate_documents"`
}

type DynamoProjectDocument struct {
	DocumentName         string        `json:"document_name"`
	DocumentFileID       string        `json:"document_file_id"`
	DocumentContentType  string        `json:"document_content_type"`
	DocumentMajorVersion int           `json:"document_major_version"`
	DocumentMinorVersion int           `json:"document_minor_version"`
	DocumentTabs         []DocumentTab `json:"document_tabs"`
}

type DocumentTab struct {
	DocumentTabType                     string `json:"document_tab_type"`
	DocumentTabID                       string `json:"document_tab_id"`
	DocumentTabName                     string `json:"document_tab_name"`
	DocumentTabPage                     int    `json:"document_tab_page"`
	DocumentTabWidth                    int    `json:"document_tab_width"`
	DocumentTabHeight                   int    `json:"document_tab_height"`
	DocumentTabIsLocked                 bool   `json:"document_tab_is_locked"`
	DocumentTabAnchorString             string `json:"document_tab_anchor_string"`
	DocumentTabAnchorIgnoreIfNotPresent bool   `json:"document_tab_anchor_ignore_if_not_present"`
	DocumentTabAnchorXOffset            int    `json:"document_tab_anchor_x_offset"`
	DocumentTabAnchorYOffset            int    `json:"document_tab_anchor_y_offset"`
}

func NewRepository(awsSession *session.Session) repository {
	return repository{
		dynamoDBClient: dynamodb.New(awsSession),
	}
}

func (r repository) GetTemplates() ([]models.Template, error) {
	templates := []models.Template{}
	for _, template := range templateMap {
		templates = append(templates, template)
	}

	return templates, nil
}

func (r repository) GetTemplate(templateID string) (models.Template, error) {
	template, ok := templateMap[templateID]
	if !ok {
		return models.Template{}, ErrTemplateNotFound
	}

	return template, nil
}

// This method belongs in the contractgroup package. We are leaving it here
// because it accesses DynamoDB, but the contractgroup repository is designed
// to connect to postgres
func (r repository) GetCLAGroup(claGroupID string) (CLAGroup, error) {
	result, err := r.dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("cla-dev-projects"),
		Key: map[string]*dynamodb.AttributeValue{
			"project_id": {
				S: aws.String(claGroupID),
			},
		},
	})
	if err != nil {
		return CLAGroup{}, err
	}

	fmt.Println(result)

	return CLAGroup{}, nil
}

func (repo repository) UpdateDynamoContractGroupTemplates(ctx context.Context, ContractGroupID string, tableName string, template models.Template) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create dynamodb Client
	svc := dynamodb.New(sess)

	// Map the fields to the dynamo model as the attribute names are different

	// Map Template Fields into DocumentTab
	cclaDocumentTabs := []DocumentTab{}

	for _, field := range template.CclaFields {
		dynamoTab := DocumentTab{
			DocumentTabType:                     field.FieldType,
			DocumentTabID:                       field.Name,
			DocumentTabPage:                     1,
			DocumentTabWidth:                    int(field.Width),
			DocumentTabHeight:                   int(field.Height),
			DocumentTabIsLocked:                 field.IsEditable,
			DocumentTabAnchorString:             field.AnchorString,
			DocumentTabAnchorIgnoreIfNotPresent: field.IsOptional,
			DocumentTabAnchorXOffset:            int(field.OffsetX),
			DocumentTabAnchorYOffset:            int(field.OffsetY),
		}
		cclaDocumentTabs = append(cclaDocumentTabs, dynamoTab)
	}

	// Map CCLA Template to Document
	dynamoCorporateProjectDocument := DynamoProjectDocument{
		DocumentName:         template.Name,
		DocumentFileID:       template.ID,
		DocumentContentType:  "storage+pdf",
		DocumentMajorVersion: 1,
		DocumentMinorVersion: 1,
		DocumentTabs:         cclaDocumentTabs,
	}

	dynamoCorporateProjectDocuments := []DynamoProjectDocument{}
	dynamoCorporateProjectDocuments = append(dynamoCorporateProjectDocuments, dynamoCorporateProjectDocument)

	dynamoCorporateProject := DynamoProjectCorporateDocuments{
		DynamoProjectDocument: dynamoCorporateProjectDocuments,
	}

	expr, err := dynamodbattribute.MarshalMap(dynamoCorporateProject)
	if err != nil {
		fmt.Println("Error marshalling Template:")
	}

	fmt.Println(expr)

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

	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	// // Map ICLA Template Fields into DocumentTab
	// iclaDocumentTabs := []DocumentTab{}

	// for _, field := range template.IclaFields {
	// 	dynamoTab := DocumentTab{
	// 		DocumentTabType:                     field.FieldType,
	// 		DocumentTabID:                       field.Name,
	// 		DocumentTabPage:                     1,
	// 		DocumentTabWidth:                    field.Width,
	// 		DocumentTabHeight:                   field.Height,
	// 		DocumentTabIsLocked:                 field.IsEditable,
	// 		DocumentTabAnchorString:             field.AnchorString,
	// 		DocumentTabAnchorIgnoreIfNotPresent: field.IsOptional,
	// 		DocumentTabAnchorXOffset:            field.OffsetX,
	// 		DocumentTabAnchorYOffset:            field.OffsetY,
	// 	}
	// 	iclaDocumentTabs = append(cclaDocumentTabs, dynamoTab)
	// }

	// // Map Template to Document
	// dynamoIndividualDocument := DynamoProjectDocument{
	// 	DocumentName:         template.Name,
	// 	DocumentFileID:       template.ID,
	// 	DocumentContentType:  "storage+pdf",
	// 	DocumentMajorVersion: 1,
	// 	DocumentMinorVersion: 1,
	// 	DocumentTabs:         iclaDocumentTabs,
	// }

	// expr, err = dynamodbattribute.MarshalMap(dynamoIndividualDocument)
	// if err != nil {
	// 	fmt.Println("Error marshalling Template:")
	// }

	// input = &dynamodb.UpdateItemInput{
	// 	ExpressionAttributeValues: expr,
	// 	TableName:                 aws.String(tableName),
	// 	Key:                       key,
	// 	ReturnValues:              aws.String("UPDATED_NEW"),
	// 	UpdateExpression:          aws.String("set project_individual_documents = :project_individual_documents"),
	// }

	// _, err = svc.UpdateItem(input)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	return err
}

var templateMap = map[string]models.Template{
	"fb4cc144-a76c-4c17-8a52-c648f158fded": models.Template{
		ID:          "fb4cc144-a76c-4c17-8a52-c648f158fded",
		Name:        "Apache Style",
		Description: "For use of projects under the Apache style of CLA.",
		MetaFields: []*models.MetaField{
			&models.MetaField{
				Name:             "Project Name",
				Description:      "Project's Full Name.",
				TemplateVariable: "PROJECT_NAME",
			},
			&models.MetaField{
				Name:             "Short Project Name",
				Description:      "The short version of the project’s name, used as a reference in the CLA.",
				TemplateVariable: "SHORT_PROJECT_NAME",
			},
			&models.MetaField{
				Name:             "Contact Email Address",
				Description:      "The E-Mail Address of the Person managing the CLA. ",
				TemplateVariable: "CONTACT_EMAIL",
			},
		},
		IclaFields: []*models.Field{
			&models.Field{
				Name:         "Full Name",
				AnchorString: "Full name:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        360,
				Height:       20,
				OffsetX:      72,
				OffsetY:      -8,
			},
			&models.Field{
				Name:         "Public Name",
				AnchorString: "Public name:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        345,
				Height:       20,
				OffsetX:      84,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Mailing Address1",
				AnchorString: "Mailing Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        325,
				Height:       20,
				OffsetX:      117,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Mailing Address2",
				AnchorString: "Mailing Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        420,
				Height:       20,
				OffsetX:      0,
				OffsetY:      29,
			},
			&models.Field{
				Name:         "Country",
				AnchorString: "Country:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        350,
				Height:       20,
				OffsetX:      60,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Telephone",
				AnchorString: "Telephone:",
				FieldType:    "text_unlocked",
				IsOptional:   true,
				IsEditable:   false,
				Width:        350,
				Height:       20,
				OffsetX:      70,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Email",
				AnchorString: "E-Mail:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        380,
				Height:       20,
				OffsetX:      50,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Please Sign",
				AnchorString: "Please Sign:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      140,
				OffsetY:      -5,
			},
			&models.Field{
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      60,
				OffsetY:      -7,
			},
		},
		CclaFields: []*models.Field{
			&models.Field{
				Name:         "Corporation Name",
				AnchorString: "Corporation Name:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        355,
				Height:       20,
				OffsetX:      140,
				OffsetY:      -5,
			},
			&models.Field{
				Name:         "Corporation Address1",
				AnchorString: "Corporation Address:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      140,
				OffsetY:      -8,
			},
			&models.Field{
				Name:         "Corporation Address2",
				AnchorString: "Corporation Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        400,
				Height:       20,
				OffsetX:      0,
				OffsetY:      29,
			},
			&models.Field{
				Name:         "Corporation Address3",
				AnchorString: "Corporation Address:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        400,
				Height:       20,
				OffsetX:      0,
				OffsetY:      64,
			},
			&models.Field{
				Name:         "Point of Contact",
				AnchorString: "Point of Contact:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      120,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Email",
				AnchorString: "E-Mail:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        340,
				Height:       20,
				OffsetX:      50,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Telephone",
				AnchorString: "Telephone:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        405,
				Height:       20,
				OffsetX:      70,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Please Sign",
				AnchorString: "Please sign:",
				FieldType:    "sign",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      140,
				OffsetY:      -6,
			},
			&models.Field{
				Name:         "Date",
				AnchorString: "Date:",
				FieldType:    "date",
				IsOptional:   false,
				IsEditable:   false,
				Width:        0,
				Height:       0,
				OffsetX:      80,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Title",
				AnchorString: "Title:",
				FieldType:    "text_unlocked",
				IsOptional:   false,
				IsEditable:   false,
				Width:        430,
				Height:       20,
				OffsetX:      40,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Corporation",
				AnchorString: "Corporation:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        385,
				Height:       20,
				OffsetX:      100,
				OffsetY:      -7,
			},
			&models.Field{
				Name:         "Schedule A",
				AnchorString: "Schedule A:",
				FieldType:    "text",
				IsOptional:   false,
				IsEditable:   false,
				Width:        550,
				Height:       600,
				OffsetX:      0,
				OffsetY:      150,
			},
		},
		HTMLBody: `<html>Hello World!</html>`,
	},
}
