// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package sign

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	v1Models "github.com/linuxfoundation/easycla/cla-backend-go/gen/v1/models"
	mock_signatures "github.com/linuxfoundation/easycla/cla-backend-go/signatures/mocks"
	"github.com/stretchr/testify/assert"
)

// a completion callback is applied once and only for the signature's own envelope (lfx-self-serve-ops#347)
func TestSignedIndividualCallbackGithubIgnoresRepeatedOrForeignEnvelopes(t *testing.T) {
	payload := envelopeXML(recipientStatuses(DocusignCompleted), envelopeDocumentStatuses) // envelope e1, signature s1
	cases := []struct {
		name      string
		signature *v1Models.Signature
	}{
		{"already signed", &v1Models.Signature{SignatureID: "s1", ProjectID: "group-a", SignatureSigned: true, SignatureEnvelopeID: "e1"}},
		{"foreign envelope", &v1Models.Signature{SignatureID: "s1", ProjectID: "group-a", SignatureSigned: false, SignatureEnvelopeID: "e2"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			signatureService := mock_signatures.NewMockSignatureService(ctrl)
			signatureService.EXPECT().GetSignature(gomock.Any(), "s1").Return(tc.signature, nil)
			svc := &service{signatureService: signatureService}
			assert.NoError(t, svc.SignedIndividualCallbackGithub(context.Background(), payload, "7", "1", "42"), "ignored callbacks are not errors: nothing is written and no status is posted")
		})
	}
}

func TestRepositoryBindingErrorsAreRecognisable(t *testing.T) {
	err := repositoryBoundToCLAGroup(&v1Models.GithubRepository{Enabled: true, RepositoryClaGroupID: "group-b"}, &v1Models.GithubOrganization{OrganizationInstallationID: 7}, 42, 7, "group-a")
	assert.True(t, errors.Is(err, errRepositoryNotBound))
}
