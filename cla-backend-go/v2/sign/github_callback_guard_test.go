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
	"github.com/linuxfoundation/easycla/cla-backend-go/signatures"
	mock_signatures "github.com/linuxfoundation/easycla/cla-backend-go/signatures/mocks"
	"github.com/stretchr/testify/assert"
)

// a completion callback is applied only for the signature's own envelope, read from the raw row (lfx-self-serve-ops#347)
func TestSignedIndividualCallbackGithubIgnoresForeignOrUnnamedEnvelopes(t *testing.T) {
	payload := envelopeXML(recipientStatuses(DocusignCompleted), envelopeDocumentStatuses) // envelope e1, signature s1
	unnamed := []byte(strings.Replace(string(payload), "<EnvelopeID>e1</EnvelopeID>", "<EnvelopeID></EnvelopeID>", 1))
	cases := []struct {
		name    string
		payload []byte
		stored  string
	}{
		{"foreign envelope", payload, "e2"},
		{"signature without an envelope", payload, ""},
		{"callback without an envelope", unnamed, "e1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			signatureService := mock_signatures.NewMockSignatureService(ctrl)
			signatureService.EXPECT().GetSignature(gomock.Any(), "s1").Return(&v1Models.Signature{SignatureID: "s1", ProjectID: "group-a"}, nil)
			signatureService.EXPECT().GetItemSignature(gomock.Any(), "s1").Return(&signatures.ItemSignature{SignatureID: "s1", SignatureEnvelopeID: tc.stored}, nil)
			svc := &service{signatureService: signatureService}
			assert.NoError(t, svc.SignedIndividualCallbackGithub(context.Background(), tc.payload, "7", "1", "42"), "ignored callbacks are not errors: nothing is written and no status is posted")
		})
	}
}

// the API model never carries the envelope, so the guard must read the raw row: a matching row lets the completion proceed
func TestSignedIndividualCallbackGithubReadsTheEnvelopeFromTheRow(t *testing.T) {
	payload := envelopeXML(recipientStatuses(DocusignCompleted), envelopeDocumentStatuses)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	signatureService := mock_signatures.NewMockSignatureService(ctrl)
	signatureService.EXPECT().GetSignature(gomock.Any(), "s1").Return(&v1Models.Signature{SignatureID: "s1", ProjectID: "group-a"}, nil)
	signatureService.EXPECT().GetItemSignature(gomock.Any(), "s1").Return(&signatures.ItemSignature{SignatureID: "s1", SignatureEnvelopeID: "e1"}, nil)
	stop := errors.New("stop after the signature update")
	signatureService.EXPECT().UpdateSignature(gomock.Any(), "s1", gomock.Any()).Return(stop).Times(1)
	svc := &service{signatureService: signatureService}
	err := svc.SignedIndividualCallbackGithub(context.Background(), payload, "7", "1", "42")
	assert.True(t, errors.Is(err, stop), "the completion proceeded to the signature update")

	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	failing := mock_signatures.NewMockSignatureService(ctrl2)
	failing.EXPECT().GetSignature(gomock.Any(), "s1").Return(&v1Models.Signature{SignatureID: "s1", ProjectID: "group-a"}, nil)
	rowErr := errors.New("dynamodb unavailable")
	failing.EXPECT().GetItemSignature(gomock.Any(), "s1").Return(nil, rowErr)
	err = (&service{signatureService: failing}).SignedIndividualCallbackGithub(context.Background(), payload, "7", "1", "42")
	assert.True(t, errors.Is(err, rowErr), "a failing row read is returned, not treated as a mismatch")
}

// a repeated completion of the signature's own envelope is processed again so a DocuSign retry can finish the remaining steps
func TestIgnoredCompletionReprocessesOwnEnvelope(t *testing.T) {
	assert.Empty(t, ignoredCompletion("s1", "e1", "e1"))
	assert.NotEmpty(t, ignoredCompletion("s1", "e1", "E1"))
}

func TestRepositoryBindingErrorsAreRecognisable(t *testing.T) {
	err := repositoryBoundToCLAGroup(&v1Models.GithubRepository{Enabled: true, RepositoryClaGroupID: "group-b"}, &v1Models.GithubOrganization{OrganizationInstallationID: 7}, 42, 7, "group-a")
	assert.True(t, errors.Is(err, errRepositoryNotBound))
}
