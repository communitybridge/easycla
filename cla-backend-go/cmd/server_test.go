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

func TestHostInSlice(t *testing.T) {
	mySlice := []string{"project.dev.lfcla.com", "corporate.dev.lfcla.com", "contributor.dev.lfcla.com", "api.dev.lfcla.com", "dev.lfcla.com", "localhost", "localhost:8100", "localhost:8101"}
	assert.True(t, hostInSlice("project.dev.lfcla.com", mySlice))
	assert.False(t, hostInSlice("*.dev.lfcla.com", mySlice))
	assert.True(t, hostInSlice("corporate.dev.lfcla.com", mySlice))
	assert.True(t, hostInSlice("contributor.dev.lfcla.com", mySlice))
	assert.True(t, hostInSlice("api.dev.lfcla.com", mySlice))
	assert.False(t, hostInSlice("https://api.dev.lfcla.com", mySlice))
	assert.True(t, hostInSlice("dev.lfcla.com", mySlice))
	assert.False(t, hostInSlice("https://dev.lfcla.com", mySlice))
	assert.True(t, hostInSlice("localhost", []string{"localhost", "foo", "bar"}))
	assert.True(t, hostInSlice("localhost", []string{"foo", "localhost:8100", "localhost:8101"}))
}
