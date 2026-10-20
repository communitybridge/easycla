// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package githublegacy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// the endpoint is unauthenticated: the SSRF guard applies unconditionally and nothing is fetched for a rejected URL
// (lfx-self-serve-ops#788)
func TestValidateOrganizationRejectsNonGitHubEndpointsBeforeFetching(t *testing.T) {
	var fetched int32
	rt := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		atomic.AddInt32(&fetched, 1)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
	})
	svc := New(&http.Client{Transport: rt})
	for _, endpoint := range []string{
		"http://127.0.0.1:9001/",
		"http://169.254.169.254/latest/meta-data/",
		"https://[::1]/",
		"https://example.com/org",
		"https://github.com.evil.example/",
		"ftp://github.com/org",
		"file:///etc/passwd",
	} {
		t.Run(endpoint, func(t *testing.T) {
			body, status, err := svc.ValidateOrganization(context.Background(), endpoint)
			if err == nil || status != http.StatusBadRequest || body != nil {
				t.Fatalf("expected a 400 rejection, got body=%v status=%d err=%v", body, status, err)
			}
		})
	}
	if atomic.LoadInt32(&fetched) != 0 {
		t.Fatalf("rejected endpoints must not be fetched, got %d requests", fetched)
	}
}

func TestValidateOrganizationFetchesGitHubEndpoints(t *testing.T) {
	var seen string
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen = r.URL.String()
		rec := httptest.NewRecorder()
		rec.WriteString(`<html><body itemtype="http://schema.org/Organization"></body></html>`)
		return rec.Result(), nil
	})
	svc := New(&http.Client{Transport: rt})
	body, status, err := svc.ValidateOrganization(context.Background(), "https://github.com/linuxfoundation")
	if err != nil || status != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("expected an ok validation, got body=%v status=%d err=%v", body, status, err)
	}
	if seen != "https://github.com/linuxfoundation" {
		t.Fatalf("unexpected fetched URL: %s", seen)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
