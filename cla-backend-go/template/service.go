package template

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/docraptor"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aymerick/raymond"
	"github.com/spf13/viper"
)

type Service interface {
	GetTemplates(ctx context.Context) ([]models.Template, error)
	CreateCLAGroupTemplate(ctx context.Context, claGroupID string, claGroupFields *models.CreateClaGroupTemplate) (string, string, error)
}

type service struct {
	templateRepo    Repository
	docraptorClient docraptor.DocraptorClient
	s3Client        *s3manager.Uploader
}

type pdfUrls struct {
	iclaPdfUrl string
	cclaPdfUrl string
}

func NewService(templateRepo Repository, docraptorClient docraptor.DocraptorClient, awsSession *session.Session) service {
	return service{
		templateRepo:    templateRepo,
		docraptorClient: docraptorClient,
		s3Client:        s3manager.NewUploader(awsSession),
	}
}

func (s service) GetTemplates(ctx context.Context) ([]models.Template, error) {
	templates, err := s.templateRepo.GetTemplates()
	if err != nil {
		return nil, err
	}

	// Remove HTML from template
	for i, template := range templates {
		template.IclaHTMLBody = ""
		template.CclaHTMLBody = ""
		templates[i] = template
	}

	return templates, nil
}

func (s service) CreateCLAGroupTemplate(ctx context.Context, claGroupID string, claGroupFields *models.CreateClaGroupTemplate) (pdfUrls, error) {
	// Verify claGroupID matches an existing CLA Group
	_, err := s.templateRepo.GetCLAGroup(claGroupID)
	if err != nil {
		return pdfUrls{}, err
	}

	// Verify the caller is authorized for the project that owns this CLA Group

	// Get Template
	template, err := s.templateRepo.GetTemplate(claGroupFields.TemplateID)
	if err != nil {
		return pdfUrls{}, err
	}

	// Apply template fields
	iclaTemplateHTML, cclaTemplateHTML, err := s.InjectProjectInformationIntoTemplate(template, claGroupFields.MetaFields)
	if err != nil {
		return pdfUrls{}, err
	}

	// Create PDF
	iclaPdf, err := s.docraptorClient.CreatePDF(iclaTemplateHTML)
	if err != nil {
		return pdfUrls{}, err
	}
	defer iclaPdf.Close()
	cclaPdf, err := s.docraptorClient.CreatePDF(cclaTemplateHTML)
	if err != nil {
		return pdfUrls{}, err
	}
	defer cclaPdf.Close()

	// Concatenate s3 bucket name with stage
	var buffer bytes.Buffer
	buffer.WriteString("cla-signature-files-")
	buffer.WriteString(viper.GetString("STAGE"))

	// Save PDF to S3
	bucket := buffer.String()
	fileNameTemplate := "contract-group/%s/template/%s"
	iclaFileName := fmt.Sprintf(fileNameTemplate, claGroupID, "icla.pdf")
	cclaFileName := fmt.Sprintf(fileNameTemplate, claGroupID, "ccla.pdf")

	iclaFileURL, err := s.SaveTemplateToS3(bucket, iclaFileName, iclaPdf)
	if err != nil {
		return pdfUrls{}, err
	}

	cclaFileURL, err := s.SaveTemplateToS3(bucket, cclaFileName, cclaPdf)
	if err != nil {
		return pdfUrls{}, err
	}

	// Concatenate tablename with stage
	buffer.Reset()
	buffer.WriteString("cla-")
	buffer.WriteString(viper.GetString("STAGE"))
	buffer.WriteString("-projects")

	// Save Template to Dynamodb
	tableName := buffer.String()
	err = s.templateRepo.UpdateDynamoContractGroupTemplates(ctx, claGroupID, tableName, template)
	if err != nil {
		return pdfUrls{}, err
	}
	template.IclaHTMLBody = iclaTemplateHTML
	template.CclaHTMLBody = cclaTemplateHTML

	pdfUrls := pdfUrls{
		iclaPdfUrl: iclaFileURL,
		cclaPdfUrl: cclaFileURL,
	}
	return pdfUrls, nil
}

func (s service) InjectProjectInformationIntoTemplate(template models.Template, metaFields []*models.MetaField) (string, string, error) {
	// TODO: Verify all template fields in template.MetaFields are present

	lookupMap := map[string]models.MetaField{}
	for _, field := range template.MetaFields {
		lookupMap[field.Name] = *field
	}

	metaFieldsMap := map[string]string{}
	for _, metaField := range metaFields {

		val, ok := lookupMap[metaField.Name]
		if !ok {
			continue
		}

		if val.Name == metaField.Name && val.TemplateVariable == metaField.TemplateVariable {
			metaFieldsMap[metaField.TemplateVariable] = metaField.Value
		}
	}
	if len(template.MetaFields) != len(metaFieldsMap) {
		return "", "", errors.New("Required fields for template were not found")
	}

	iclaTemplateHTML, err := raymond.Render(template.IclaHTMLBody, metaFieldsMap)
	if err != nil {
		return "", "", err
	}

	cclaTemplateHTML, err := raymond.Render(template.CclaHTMLBody, metaFieldsMap)
	if err != nil {
		return "", "", err
	}

	return iclaTemplateHTML, cclaTemplateHTML, nil
}

func (s service) SaveTemplateToS3(bucket, filepath string, template io.ReadCloser) (string, error) {
	defer template.Close()

	// Upload the file to S3.
	result, err := s.s3Client.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
		Body:   template,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3 Bucket, %v", err)
	}

	return result.Location, nil
}
