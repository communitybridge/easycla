package contractgroup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aymerick/raymond"
)

var (
	userID = "***REMOVED***"
)

type Service interface {
	CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (models.ContractGroup, error)
	GetContractGroups(ctx context.Context, projectID string) ([]models.ContractGroup, error)

	CreateContractTemplate(ctx context.Context, contractTemplate models.ContractTemplate, contractID string) (models.ContractTemplate, error)

	CreateGitHubOrganization(ctx context.Context, contractID string, githubOrg models.Github) (models.Github, error)

	CreateGerritInstance(ctx context.Context, projectSfdcID, contractID, userID string, gerritInstance models.Gerrit) (models.Gerrit, error)
	DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error

	GetContractGroupSignatures(ctx context.Context, projectSfdcID string, contractID string) (models.ContractGroupSignatures, error)
}

type service struct {
	contractGroupRepo Repository
}

func (s Service) SaveTemplateToDynamoDB(template Template, templateName, tableName, contractGroupID, region string) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	item := dynamodbattribute.MarshalMap(template)

	// Create item in table Movies
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

func (s Service) SaveFileToS3Bucket(file io.ReadCloser, bucketName, fileName, region string) error {
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
func (s Service) SendHTMLToDocRaptor(URL, HTML string) io.Reader {
	document := `{
  		"type": "pdf",
  		"document_content": "%s",
  		"test":true
	}`
	document = fmt.Sprintf(document, HTML)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBufferString(document))
	if err != nil {
		fmt.Printf("failed to create request to submit data to API: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to submit data to DocRaptorAPI: %s", err)
	}

	fmt.Printf("API Response Status Code: %s\n", resp.Status)

	return resp.Body
}

func (s service) InjectProjectInformationIntoTemplate(projectName, shortProjectName, documentType, majorVersion, minorVersion, contactEmail string) string {
	// DocRaptor API likes HTML in single line
	templateBefore := `<html><body><p style=\"text-align: center\">{{projectName}}<br />{{documentType}} Contributor License Agreement (\"Agreement\")v{{majorVersion}}.{{minorVersion}}</p><p>Thank you for your interest in {{projectName}} project (“{{shortProjectName}}”) of The Linux Foundation (the “Foundation”). In order to clarify the intellectual property license granted with Contributions from any person or entity, the Foundation must have a Contributor License Agreement (“CLA”) on file that has been signed by each Contributor, indicating agreement to the license terms below. This license is for your protection as a Contributor as well as the protection of {{shortProjectName}}, the Foundation and its users; it does not change your rights to use your own Contributions for any other purpose.</p><p>If you have not already done so, please complete and sign this Agreement using the electronic signature portal made available to you by the Foundation or its third-party service providers, or email a PDF of the signed agreement to {{contactEmail}}. Please read this document carefully before signing and keep a copy for your records.</p></body></html>`
	fieldsMap := map[string]string{
		"projectName":      projectName,
		"shortProjectName": shortProjectName,
		"documentType":     documentType,
		"majorVersion":     majorVersion,
		"minorVersion":     minorVersion,
		"contactEmail":     contactEmail,
	}

	templateAfter, err := raymond.Render(templateBefore, fieldsMap)
	if err != nil {
		fmt.Println("Failed to enter fields into HTML", err)
	}

	return templateAfter
}

func NewService(contractGroupRepo Repository) service {
	return service{
		contractGroupRepo: contractGroupRepo,
	}
}

func (s service) CreateContractGroup(ctx context.Context, projectSfdcID string, contractGroup models.ContractGroup) (models.ContractGroup, error) {
	contractGroupID, err := s.contractGroupRepo.CreateContractGroup(ctx, projectSfdcID, contractGroup)
	if err != nil {
		return models.ContractGroup{}, err
	}

	contractGroup.ContractGroupID = contractGroupID
	contractGroup.ProjectSfdcID = projectSfdcID

	return contractGroup, nil
}

func (s service) GetContractGroups(ctx context.Context, projectSfdcID string) ([]models.ContractGroup, error) {
	contractGroups, err := s.contractGroupRepo.GetContractGroups(ctx, projectSfdcID)
	if err != nil {
		return nil, err
	}

	for i, contractGroup := range contractGroups {
		contractGroups[i].CclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "CCLA")
		if err != nil {
			return nil, err
		}

		contractGroups[i].IclaTemplate, err = s.contractGroupRepo.GetLatestContractTemplate(ctx, contractGroup.ContractGroupID, "ICLA")
		if err != nil {
			return nil, err
		}

		contractGroups[i].GithubOrganizations, err = s.contractGroupRepo.GetGithubOrganizatons(ctx, contractGroup.ContractGroupID)
		if err != nil {
			return nil, err
		}

		contractGroups[i].GerritInstances, err = s.contractGroupRepo.GetGerritInstances(ctx, contractGroup.ContractGroupID)
		if err != nil {
			return nil, err
		}
	}

	return contractGroups, nil
}

func (s service) CreateContractTemplate(ctx context.Context, contractTemplate models.ContractTemplate, contractID string) (models.ContractTemplate, error) {
	contractTemplateID, err := s.contractGroupRepo.CreateContractTemplate(ctx, contractID, contractTemplate)
	if err != nil {
		return models.ContractTemplate{}, err
	}

	contractTemplate.ContractTemplateID = contractTemplateID

	return contractTemplate, nil
}

func (s service) CreateGitHubOrganization(ctx context.Context, contractID string, githubOrg models.Github) (models.Github, error) {
	githubOrgID, err := s.contractGroupRepo.CreateGitHubOrganization(ctx, contractID, userID, githubOrg)
	if err != nil {
		return models.Github{}, err
	}

	githubOrg.GithubOrganizationID = githubOrgID

	return githubOrg, nil
}

func (s service) CreateGerritInstance(ctx context.Context, projectSFDCID, contractID, userID string, gerritInstance models.Gerrit) (models.Gerrit, error) {
	gerritInstanceID, err := s.contractGroupRepo.CreateGerritInstance(ctx, projectSFDCID, contractID, userID, gerritInstance)
	if err != nil {
		return models.Gerrit{}, err
	}

	gerritInstance.GerritInstanceID = gerritInstanceID

	return gerritInstance, nil

}

func (s service) DeleteGerritInstance(ctx context.Context, projectSfdcID string, contractID string, gerritInstanceID string) error {
	err := s.contractGroupRepo.DeleteGerritInstance(ctx, projectSfdcID, contractID, gerritInstanceID)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetContractGroupSignatures(ctx context.Context, projectSFDCID string, contractID string) (models.ContractGroupSignatures, error) {
	contractGoupSignatures := models.ContractGroupSignatures{ContractGroupID: contractID}
	var err error

	contractGoupSignatures.CclaSignatures, err = s.contractGroupRepo.GetContractGroupCCLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		return models.ContractGroupSignatures{}, err
	}

	contractGoupSignatures.IclaSignatures, err = s.contractGroupRepo.GetContractGroupICLASignatures(ctx, projectSFDCID, contractID)
	if err != nil {
		return models.ContractGroupSignatures{}, err
	}

	return contractGoupSignatures, nil
}
