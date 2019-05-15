package template

import (
	"context"
	"fmt"
	"io"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/docraptor"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aymerick/raymond"
)

type Service interface {
	GetTemplates(ctx context.Context) ([]models.Template, error)
}

type service struct {
	templateRepo    Repository
	docraptorClient docraptor.DocraptorClient
}

func NewService(templateRepo Repository, docraptorClient docraptor.DocraptorClient) service {
	return service{
		templateRepo:    templateRepo,
		docraptorClient: docraptorClient,
	}
}

func (s service) CreateCLAGroupTemplate(ctx context.Context, claGroupID string, claGroupFields *models.CreateClaTemplateGroup) error {
	// template, err := s.templateRepo.CreateCLAGroupTemplate(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// use some methods to retrieve template from repository using templateID
	// template := s.repo.GetTemplate(claGroupFields.TemplateID)
	template := "testvalue"
	bucketName := "testvalue"
	fileName := "testvalue"
	region := "testvalue"
	HTML := s.InjectProjectInformationIntoTemplate(template, claGroupFields.MetaFields)
	PDF := s.docraptorClient.CreatePDF(HTML)
	err := s.SaveFileToS3Bucket(PDF, bucketName, fileName, region)
	if err != nil {
		fmt.Println("Error saving to S3 bucket : ", err)
		return err
	}
	// Save Template to Dynamodb once method is finalized
	return nil
}

func (s service) InjectProjectInformationIntoTemplate(template string, fields []*models.MetaField) string {
	// DocRaptor API likes HTML in single line
	//templateBefore := `<html><body><p style=\"text-align: center\">{{projectName}}<br />{{documentType}} Contributor License Agreement (\"Agreement\")v{{majorVersion}}.{{minorVersion}}</p><p>Thank you for your interest in {{projectName}} project (“{{shortProjectName}}”) of The Linux Foundation (the “Foundation”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Foundation must have a Contributor License Agreement (“CLA”) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of {{shortProjectName}}, the Foundation and its users; it does not change your rights to use your own Contributions for any other purpose.</p><p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Foundation or its third-party service providers, or email a PDF of the signed agreement to {{contactEmail}}. Please read this document carefully before signing and keep a copy for your records.</p></body></html>`
	// add logic to parse fields to inject metafields in proper places
	// loop through fields to populate map
	templateBefore := template
	fieldsMap := map[string]string{}

	templateAfter, err := raymond.Render(templateBefore, fieldsMap)
	if err != nil {
		fmt.Println("Failed to enter fields into HTML", err)
	}

	return templateAfter
}

func (s service) SaveTemplateToDynamoDB(template models.Template, templateName, tableName, contractGroupID, region string) error {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	item, err := dynamodbattribute.MarshalMap(template)
	if err != nil {
		fmt.Println("Error marshaling values into item: ", err)
		return err
	}

	// Create item in table
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Error putting item in database: ", err)
		return err
	}

	fmt.Println("Successfully put item in database.")
	return nil

}

func (s service) SaveFileToS3Bucket(file io.ReadCloser, bucketName, fileName, region string) error {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	},
	))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3 Bucket, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	defer file.Close()

	return nil
}

func (s service) GetTemplates(ctx context.Context) ([]models.Template, error) {
	templates, err := s.templateRepo.GetTemplates(ctx)
	if err != nil {
		return nil, err
	}

	return templates, nil
}
