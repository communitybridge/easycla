package docraptor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errInvalidKey = errors.New("invalid key")
)

const (
	docraptorURL = "https://%s@docraptor.com/docs"
)

type DocraptorClient struct {
	apiKey   string
	url      string
	testMode bool
}

func NewDocraptorClient(key string, testMode bool) (DocraptorClient, error) {
	if key == "" {
		return DocraptorClient{}, errInvalidKey
	}

	url := fmt.Sprintf(docraptorURL, key)

	return DocraptorClient{
		apiKey:   key,
		url:      url,
		testMode: testMode,
	}, nil
}

func (dc DocraptorClient) CreatePDF(html string) (io.ReadCloser, error) {
	document := `{
  		"type": "pdf",
  		"document_content": "%s",
  		"test": %v
	}`
	document = fmt.Sprintf(document, html, dc.testMode)

	req, err := http.NewRequest(http.MethodPost, dc.url, bytes.NewBufferString(document))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
