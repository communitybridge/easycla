// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package tests

import (
	"testing"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestTrimSpaceFromItems(t *testing.T) {
	testInputs := [][]string{
		nil,
		{},
		{"  "},
		{" test1"},
		{"test1 "},
		{" test1      "},
		{" test 1      "},
		{" test    1      "},
		{" test1", "test 2 ", " test  3 "},
	}
	expectedResults := [][]string{
		nil,
		{},
		{""},
		{"test1"},
		{"test1"},
		{"test1"},
		{"test1"},
		{"test1"},
		{"test1", "test 2", "test  3"},
	}

	for i := range testInputs {
		assert.ObjectsAreEqualValues(expectedResults[i], utils.TrimSpaceFromItems(testInputs[i]))
	}
}
