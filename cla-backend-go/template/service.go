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
	"github.com/spf13/viper"
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
	HTML := s.InjectProjectInformationIntoTemplate(template, claGroupFields.MetaFields)
	PDF := s.docraptorClient.CreatePDF(HTML)
	err := s.SaveFileToS3Bucket(PDF)
	if err != nil {
		fmt.Println("Error saving to S3 bucket : ", err)
		return err
	}
	// Save Template to Dynamodb once method is finalized
	return nil
}

func (s service) InjectProjectInformationIntoTemplate(template string, fields []*models.MetaField) string {
	templateBefore := template
	fieldsMap := map[string]string{}
	for _, field := range fields {
		fieldsMap[field.TemplateVariable] = field.Value
	}

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

func (s service) SaveFileToS3Bucket(file io.ReadCloser) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("AWS_REGION")),
	},
	))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("AWS_BUCKET_NAME")),
		Key:    aws.String("savedFile"),
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
