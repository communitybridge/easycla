package docraptor

import (
	"bytes"
	"encoding/json"
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
	document := map[string]interface{}{
		"document_type":    "pdf",
		"document_content": html,
		"name":             "docraptor-go.pdf",
		"test":             dc.testMode,
	}

	documentBytes, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(dc.url, "application/json", bytes.NewBuffer(documentBytes))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
