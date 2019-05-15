package docraptor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DocraptorClient struct {
	APIKey string
	URL    string
}

func NewDocraptorClient(key string) (DocraptorClient, error) {
	if key == "" {
		return DocraptorClient{}, errors.New("invalid key")
	}
	URL := fmt.Sprintf("https://%s@docraptor.com/docs", key)

	return DocraptorClient{
		APIKey: key,
		URL:    URL,
	}, nil
}

func (dc DocraptorClient) CreatePDF(HTML string) io.ReadCloser {
	document := `{
  		"type": "pdf",
  		"document_content": "%s",
  		"test":true
	}`
	document = fmt.Sprintf(document, HTML)

	req, err := http.NewRequest(http.MethodPost, dc.URL, bytes.NewBufferString(document))
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
