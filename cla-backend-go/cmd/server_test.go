// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	mySlice := []string{"aaaa", "bbbb", "cccc"}
	assert.True(t, stringInSlice("aaaa", mySlice))
	assert.True(t, stringInSlice("bbbb", mySlice))
	assert.True(t, stringInSlice("cccc", mySlice))
	assert.False(t, stringInSlice("aaaa1", mySlice))
	assert.False(t, stringInSlice("aaaa", nil))
	assert.False(t, stringInSlice("aaaa", []string{}))
}
