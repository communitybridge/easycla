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

func TestStringInSliceHostname(t *testing.T) {
	mySlice := []string{"project.dev.lfcla.com", "corporate.dev.lfcla.com", "contributor.dev.lfcla.com", "api.dev.lfcla.com", "dev.lfcla.com"}
	assert.True(t, stringInSlice("project.dev.lfcla.com", mySlice))
	assert.False(t, stringInSlice("*.dev.lfcla.com", mySlice))
	assert.True(t, stringInSlice("corporate.dev.lfcla.com", mySlice))
	assert.True(t, stringInSlice("contributor.dev.lfcla.com", mySlice))
	assert.True(t, stringInSlice("api.dev.lfcla.com", mySlice))
	assert.False(t, stringInSlice("https://api.dev.lfcla.com", mySlice))
	assert.True(t, stringInSlice("dev.lfcla.com", mySlice))
	assert.False(t, stringInSlice("https://dev.lfcla.com", mySlice))
}
