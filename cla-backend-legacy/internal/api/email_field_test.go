// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package api

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestEmailFieldEscapesAndBounds(t *testing.T) {
	cases := map[string]struct {
		in   string
		max  int
		want string
	}{
		"anchor":    {`<a href="https://evil.invalid">click</a>`, 60, "&lt;a href=&#34;https://evil.invalid&#34;&gt;click&lt;/a&gt;"},
		"img":       {" <img src=x onerror=alert(1)> ", 60, "&lt;img src=x onerror=alert(1)&gt;"},
		"truncated": {strings.Repeat("a", 61), 60, strings.Repeat("a", 60)},
		"multibyte": {strings.Repeat("é", 31), 61, strings.Repeat("é", 30)},
		"plain":     {"Jane Doe", 60, "Jane Doe"},
	}
	for name, tc := range cases {
		if got := emailField(tc.in, tc.max); got != tc.want {
			t.Errorf("%s: emailField(%q, %d) = %q, want %q", name, tc.in, tc.max, got, tc.want)
		}
	}
}

func TestPublicSignatureItemDropsTheDocuSignEnvelope(t *testing.T) {
	item := map[string]types.AttributeValue{
		"signature_id":          &types.AttributeValueMemberS{Value: "sig-1"},
		"signature_signed":      &types.AttributeValueMemberBOOL{Value: true},
		"user_docusign_raw_xml": &types.AttributeValueMemberS{Value: "<DocuSignEnvelopeInformation/>"},
	}
	out := publicSignatureItem(item)
	if out["signature_id"] != "sig-1" || out["signature_signed"] != true {
		t.Fatalf("unexpected serialization: %#v", out)
	}
	if _, present := out["user_docusign_raw_xml"]; present {
		t.Fatal("user_docusign_raw_xml must never be returned to a client")
	}
}
