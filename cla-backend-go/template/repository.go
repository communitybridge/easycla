package template

import (
	"errors"
	"fmt"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

type Repository interface {
	GetTemplates() ([]models.Template, error)
	GetTemplate(templateID string) (models.Template, error)
	GetCLAGroup(claGroupID string) (CLAGroup, error)
}

type repository struct {
	dynamoDBClient *dynamodb.DynamoDB
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

type CLAGroup struct {
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

var templateMap = map[string]models.Template{
	"fb4cc144-a76c-4c17-8a52-c648f158fded": models.Template{
		TemplateID:  "fb4cc144-a76c-4c17-8a52-c648f158fded",
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
		HTMLBody: `<html><body><p style=\"text-align: center\">{{PROJECT_NAME}}<br /> Contributor License Agreement ("Agreement")v</p><p>Thank you for your interest in {{PROJECT_NAME}} project (“{{SHORT_PROJECT_NAME}}”) of The Linux Foundation (the “Foundation”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Foundation must have a Contributor License Agreement (“CLA”) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of {{SHORT_PROJECT_NAME}}, the Foundation and its users; it does not change your rights to use your own Contributions for any other purpose.</p><p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Foundation or its third-party service providers, or email a PDF of the signed agreement to {{CONTACT_EMAIL}}. Please read this document carefully before signing and keep a copy for your records.</p></body></html>`,
	},
}
