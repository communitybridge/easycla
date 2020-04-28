// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package template

import (
	"context"
	"errors"
	"fmt"
	"io"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/docraptor"
	"github.com/communitybridge/easycla/cla-backend-go/gen/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aymerick/raymond"
)

// Service interface
type Service interface {
	GetTemplates(ctx context.Context) ([]models.Template, error)
	CreateCLAGroupTemplate(ctx context.Context, claGroupID string, claGroupFields *models.CreateClaGroupTemplate) (models.TemplatePdfs, error)
}

type service struct {
	stage           string // The AWS stage (dev, staging, prod)
	templateRepo    Repository
	docraptorClient docraptor.Client
	s3Client        *s3manager.Uploader
}

// NewService API call
func NewService(stage string, templateRepo Repository, docraptorClient docraptor.Client, awsSession *session.Session) service {
	return service{
		stage:           stage,
		templateRepo:    templateRepo,
		docraptorClient: docraptorClient,
		s3Client:        s3manager.NewUploader(awsSession),
	}
}

// GetTemplates API call
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

// CreateCLAGroupTemplate
func (s service) CreateCLAGroupTemplate(ctx context.Context, claGroupID string, claGroupFields *models.CreateClaGroupTemplate) (models.TemplatePdfs, error) {
	// Verify claGroupID matches an existing CLA Group
	claGroup, err := s.templateRepo.GetCLAGroup(claGroupID)
	if err != nil {
		log.Warnf("Unable to fetch CLA group by id: %s, error: %v - returning empty template PDFs", claGroupID, err)
		return models.TemplatePdfs{}, err
	}

	// Verify the caller is authorized for the project that owns this CLA Group

	// Get Template
	template, err := s.templateRepo.GetTemplate(claGroupFields.TemplateID)
	if err != nil {
		log.Warnf("Unable to fetch template fields: %s, error: %v - returning empty template PDFs",
			claGroupFields.TemplateID, err)
		return models.TemplatePdfs{}, err
	}

	// Apply template fields
	iclaTemplateHTML, cclaTemplateHTML, err := s.InjectProjectInformationIntoTemplate(template, claGroupFields.MetaFields)
	if err != nil {
		log.Warnf("Unable to inject metadata details into template, error: %v - returning empty template PDFs", err)
		return models.TemplatePdfs{}, err
	}

	bucket := fmt.Sprintf("cla-signature-files-%s", s.stage)
	fileNameTemplate := "contract-group/%s/template/%s"

	// Create PDF
	var pdfUrls models.TemplatePdfs
	var iclaFileURL string
	var cclaFileURL string

	if claGroup.ProjectICLAEnabled {
		iclaPdf, iclaErr := s.docraptorClient.CreatePDF(iclaTemplateHTML)
		if iclaErr != nil {
			log.Warnf("Problem generating ICLA template via docraptor client, error: %v - returning empty template PDFs", err)
			return models.TemplatePdfs{}, err
		}
		defer func() {
			closeErr := iclaPdf.Close()
			if closeErr != nil {
				log.Warnf("error closing ICLA PDF, error: %v", closeErr)
			}
		}()
		iclaFileName := fmt.Sprintf(fileNameTemplate, claGroupID, "icla.pdf")
		iclaFileURL, err = s.SaveTemplateToS3(bucket, iclaFileName, iclaPdf)
		if err != nil {
			log.Warnf("Problem uploading ICLA PDF: %s to s3, error: %v - returning empty template PDFs", iclaFileName, err)
			return models.TemplatePdfs{}, err
		}
		template.IclaHTMLBody = iclaTemplateHTML
	}

	if claGroup.ProjectCCLAEnabled {
		cclaPdf, cclaErr := s.docraptorClient.CreatePDF(cclaTemplateHTML)
		if cclaErr != nil {
			log.Warnf("Problem generating CCLA template via docraptor client, error: %v - returning empty template PDFs", err)
			return models.TemplatePdfs{}, err
		}
		defer func() {
			closeErr := cclaPdf.Close()
			if closeErr != nil {
				log.Warnf("error closing CCLA PDF, error: %v", closeErr)
			}
		}()
		cclaFileName := fmt.Sprintf(fileNameTemplate, claGroupID, "ccla.pdf")
		cclaFileURL, err = s.SaveTemplateToS3(bucket, cclaFileName, cclaPdf)
		if err != nil {
			log.Warnf("Problem uploading CCLA PDF: %s to s3, error: %v - returning empty template PDFs", cclaFileName, err)
			return models.TemplatePdfs{}, err
		}
		template.CclaHTMLBody = cclaTemplateHTML
	}

	if claGroup.ProjectICLAEnabled && claGroup.ProjectCCLAEnabled {
		pdfUrls = models.TemplatePdfs{
			IndividualPDFURL: iclaFileURL,
			CorporatePDFURL:  cclaFileURL,
		}
	} else if claGroup.ProjectCCLAEnabled {
		pdfUrls = models.TemplatePdfs{
			CorporatePDFURL: cclaFileURL,
		}
	} else if claGroup.ProjectICLAEnabled {
		pdfUrls = models.TemplatePdfs{
			IndividualPDFURL: iclaFileURL,
		}
	}

	// Save Template to DynamoDB
	err = s.templateRepo.UpdateDynamoContractGroupTemplates(ctx, claGroupID, template, pdfUrls, claGroup.ProjectCCLAEnabled, claGroup.ProjectICLAEnabled)
	if err != nil {
		log.Warnf("Problem updating the database with ICLA/CCLA new PDF details, error: %v - returning empty template PDFs", err)
		return models.TemplatePdfs{}, err
	}

	return pdfUrls, nil
}

// InjectProjectInformationIntoTemplate
func (s service) InjectProjectInformationIntoTemplate(template models.Template, metaFields []*models.MetaField) (string, string, error) {
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
		return "", "", errors.New("required fields for template were not found")
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

// SaveTemplateToS3
func (s service) SaveTemplateToS3(bucket, filepath string, template io.ReadCloser) (string, error) {
	defer template.Close()

	// Upload the file to S3.
	result, err := s.s3Client.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filepath),
		Body:        template,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3 Bucket: %s / %s, %v", bucket, filepath, err)
	}

	return result.Location, nil
}
