// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package tests

import (
	"fmt"
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidCompanyName(t *testing.T) {
	// Valid company names
	validInput := []string{
		"This is a company name",
		"23 Air Jordan",
		"sdfsdfsdfááÀÁ(test)[test]0343+_7343-9(@)/the world",
		"ááÀÁ test",
		"世界",
	}

	for _, str := range validInput {
		utils.ValidCompanyName(str)
		assert.True(t, utils.ValidCompanyName(str), fmt.Sprintf("ValidCompanyName - %s", str))
	}

	var longString string
	for i := 1; i <= 256; i++ {
		longString += "a"
	}
	//log.Printf("longString = %d", len(longString))

	// Invalid company names
	inValidInput := []string{
		"+invalid name",
		"%invalid name",
		"!invalid name",
		".test",
		"?ááÀÁ test",
		"",  // min 2 chars
		"1", // min 2 chars
		longString,
	}

	for _, str := range inValidInput {
		assert.False(t, utils.ValidCompanyName(str), fmt.Sprintf("ValidCompanyName - %s", str))
	}
}
