// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package approval_list

import (
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/stretchr/testify/assert"
)

func TestRequestApprovedEmailToRecipientContent(t *testing.T) {
	subject, body, recipients := requestApprovedEmailToRecipientContent(
		&models.Company{
			CompanyName: "gardenerLtd"},
		&models.ClaGroup{Version: "v2"},
		"john",
		"john@john.com")

	assert.Equal(t, "EasyCLA: Company Manager Access Approved for gardenerLtd", subject)
	assert.Equal(t, []string{"john@john.com"}, recipients)
	assert.Contains(t, body, "Hello john,")
	assert.Contains(t, body, "This is a notification email from EasyCLA regarding the company gardenerLtd")
	assert.Contains(t, body, "You have now been approved as a Company Manager for gardenerLtd")
}
