// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a missing key is expected until the allow-list is provisioned, any other failure is not, so
// the two must stay distinguishable in the logs even when the SDK error arrives wrapped
func TestIsParameterNotFound(t *testing.T) {
	assert.True(t, isParameterNotFound(awserr.New(ssm.ErrCodeParameterNotFound, "not found", nil)))
	assert.True(t, isParameterNotFound(fmt.Errorf("wrapped: %w", awserr.New(ssm.ErrCodeParameterNotFound, "not found", nil))))
	assert.False(t, isParameterNotFound(awserr.New("AccessDeniedException", "denied", nil)))
	assert.False(t, isParameterNotFound(awserr.New("ThrottlingException", "slow down", nil)))
	assert.False(t, isParameterNotFound(errors.New("connection reset")))
}

func TestParseTrustedClientIDs(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		clientIDs []string
	}{
		{"unset", "", nil},
		{"blank", "   ", nil},
		{"separators only", " , ,, ", nil},
		{"single", "ss-client", []string{"ss-client"}},
		{"padded", "  ss-client  ", []string{"ss-client"}},
		{"multiple", "ss-client, other-client ,third-client", []string{"ss-client", "other-client", "third-client"}},
		{"trailing separator", "ss-client,", []string{"ss-client"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.clientIDs, parseTrustedClientIDs(test.value))
		})
	}
}

func TestLoadConfigRejectsUnreadableSSS(t *testing.T) {
	oldConfig := easyCLAConfig
	t.Cleanup(func() { easyCLAConfig = oldConfig })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name  string
			Names []string
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/x-amz-json-1.1")
		if strings.HasPrefix(input.Name, "cla-sss-") || len(input.Names) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"__type":"ThrottlingException","message":"Rate exceeded"}`))
			assert.NoError(t, err)
			return
		}
		value := "test-value"
		if strings.HasPrefix(input.Name, "cla-gh-app-id-") {
			value = "123"
		} else if strings.Contains(input.Name, "-enabled-") || strings.Contains(input.Name, "-enable-services-for-parent-") {
			value = "false"
		}
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"Parameter": map[string]string{"Name": input.Name, "Value": value, "Type": "String"},
		})
		assert.NoError(t, err)
	}))
	defer server.Close()
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String(server.URL),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		MaxRetries:  aws.Int(0),
	})
	require.NoError(t, err)
	_, err = LoadConfig("", awsSession, "test")
	require.Error(t, err, "an unreadable SSS policy must not become optional or unconfigured")
	assert.Contains(t, err.Error(), "ThrottlingException")
}

func TestLoadOptionalSSSConfig(t *testing.T) {
	tests := []struct {
		name       string
		values     map[string]string
		failKey    string
		errorCode  string
		incomplete bool
		want       SSS
	}{
		{name: "not provisioned", want: SSS{Enabled: true}},
		{
			name: "configured and trimmed",
			values: map[string]string{
				"base-url": " https://sss.invalid ", "auth0-audience": " sss-audience ",
				"required": " true ", "enabled": "true",
			},
			want: SSS{BaseURL: "https://sss.invalid", Audience: "sss-audience", Required: true, Enabled: true},
		},
		{
			name:   "kill switch",
			values: map[string]string{"required": "true", "enabled": "false"},
			want:   SSS{Required: true},
		},
		{
			name:   "missing flags keep defaults",
			values: map[string]string{"base-url": "https://sss.invalid", "auth0-audience": "sss-audience"},
			want:   SSS{BaseURL: "https://sss.invalid", Audience: "sss-audience", Enabled: true},
		},
		{
			name:   "malformed flags keep defaults",
			values: map[string]string{"required": "invalid", "enabled": "invalid"},
			want:   SSS{Enabled: true},
		},
		{name: "incomplete success is not absence", failKey: "base-url", incomplete: true},
		{name: "access denied is not absence", failKey: "required", errorCode: "AccessDeniedException"},
	}
	for _, key := range []string{"base-url", "auth0-audience", "required", "enabled"} {
		tests = append(tests, struct {
			name       string
			values     map[string]string
			failKey    string
			errorCode  string
			incomplete bool
			want       SSS
		}{name: "throttled " + key, failKey: key, errorCode: "ThrottlingException"})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assert.Equal(t, "AmazonSSM.GetParameter", r.Header.Get("X-Amz-Target"))
				var input ssm.GetParameterInput
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					t.Error(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				key := strings.TrimSuffix(strings.TrimPrefix(aws.StringValue(input.Name), "cla-sss-"), "-test")
				w.Header().Set("Content-Type", "application/x-amz-json-1.1")
				if key == tt.failKey && tt.incomplete {
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{}))
					return
				}
				value, found := tt.values[key]
				code := ""
				if key == tt.failKey {
					code = tt.errorCode
				} else if !found {
					code = ssm.ErrCodeParameterNotFound
				}
				if code != "" {
					w.WriteHeader(http.StatusBadRequest)
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]string{"__type": code, "message": "test SSM error"}))
					return
				}
				assert.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
					"Parameter": map[string]string{"Name": aws.StringValue(input.Name), "Value": value},
				}))
			}))
			defer server.Close()
			awsSession, err := session.NewSession(&aws.Config{
				Region:      aws.String("us-east-1"),
				Endpoint:    aws.String(server.URL),
				Credentials: credentials.NewStaticCredentials("test", "test", ""),
				MaxRetries:  aws.Int(0),
			})
			require.NoError(t, err)
			original := SSS{BaseURL: "unchanged", Required: true}
			cfg := Config{SSS: original}
			err = loadOptionalSSSConfig(ssm.New(awsSession), "test", &cfg, nil)
			if tt.failKey != "" {
				require.Error(t, err)
				assert.Equal(t, original, cfg.SSS, "do not publish partial or defaulted policy after a read failure")
				assert.Contains(t, err.Error(), tt.failKey)
				if tt.errorCode != "" {
					assert.Contains(t, err.Error(), tt.errorCode)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, cfg.SSS)
				assert.Equal(t, int32(4), calls.Load())
			}
		})
	}
}
