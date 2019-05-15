package docraptor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aymerick/raymond"
)

type DocraptorClient struct {
	APIKey    string
	HTMLInput string
}

func NewDocraptorClient(key, html string) (Client, error) {
	if key == "" {
		return DocraptorClient{}, errors.New("invalid key")
	}

	if html == "" {
		return DocraptorClient{}, errors.New("invalid html")
	}

	return DocRaptorClient{
		APIKey:   key,
		HTMLPage: html,
	}, nil
}

func (dc DocraptorClient) SendHTMLToDocRaptor(APIKey, HTML string) io.ReadCloser {
	URL = "https://%s@docraptor.com/docs"
	URL = fmt.Sprint(URL, APIKey)
	document := `{
  		"type": "pdf",
  		"document_content": "%s",
  		"test":true
	}`
	document = fmt.Sprintf(document, HTML)

	request, err := http.NewRequest(http.MethodPost, URL, bytes.NewBufferString(document))
	if err != nil {
		fmt.Printf("failed to create request to submit data to API: %s", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to submit data to DocRaptorAPI: %s", err)
	}

	fmt.Printf("API Response Status Code: %s\n", response.Status)

	return response.Body
}

func (dc DocraptorClient) InjectProjectInformationIntoTemplate(projectName, shortProjectName, documentType, majorVersion, minorVersion, contactEmail string) string {
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
