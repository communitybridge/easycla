// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	v1Models "github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	mock_signatures "github.com/linuxfoundation/easycla/cla-backend-go/signatures/mocks"
	"github.com/stretchr/testify/assert"
)

// a completion callback is applied only for the signature's own envelope, named exactly (lfx-self-serve-ops#347)
func TestSignedIndividualCallbackGithubIgnoresForeignOrUnnamedEnvelopes(t *testing.T) {
	payload := envelopeXML(recipientStatuses(DocusignCompleted), envelopeDocumentStatuses) // envelope e1, signature s1
	unnamed := []byte(strings.Replace(string(payload), "<EnvelopeID>e1</EnvelopeID>", "<EnvelopeID></EnvelopeID>", 1))
	cases := []struct {
		name      string
		payload   []byte
		signature *v1Models.Signature
	}{
		{"foreign envelope", payload, &v1Models.Signature{SignatureID: "s1", ProjectID: "group-a", SignatureEnvelopeID: "e2"}},
		{"signature without an envelope", payload, &v1Models.Signature{SignatureID: "s1", ProjectID: "group-a"}},
		{"callback without an envelope", unnamed, &v1Models.Signature{SignatureID: "s1", ProjectID: "group-a", SignatureEnvelopeID: "e1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			signatureService := mock_signatures.NewMockSignatureService(ctrl)
			signatureService.EXPECT().GetSignature(gomock.Any(), "s1").Return(tc.signature, nil)
			svc := &service{signatureService: signatureService}
			assert.NoError(t, svc.SignedIndividualCallbackGithub(context.Background(), tc.payload, "7", "1", "42"), "ignored callbacks are not errors: nothing is written and no status is posted")
		})
	}
}

// a repeated completion of the signature's own envelope is processed again so a DocuSign retry can finish the remaining steps
func TestIgnoredCompletionReprocessesOwnEnvelope(t *testing.T) {
	assert.Empty(t, ignoredCompletion(&v1Models.Signature{SignatureID: "s1", SignatureSigned: true, SignatureEnvelopeID: "e1"}, "e1"))
	assert.NotEmpty(t, ignoredCompletion(&v1Models.Signature{SignatureID: "s1", SignatureEnvelopeID: "e1"}, "E1"))
}

func TestRepositoryBindingErrorsAreRecognisable(t *testing.T) {
	err := repositoryBoundToCLAGroup(&v1Models.GithubRepository{Enabled: true, RepositoryClaGroupID: "group-b"}, &v1Models.GithubOrganization{OrganizationInstallationID: 7}, 42, 7, "group-a")
	assert.True(t, errors.Is(err, errRepositoryNotBound))
}
