// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmtDuration(t *testing.T) {
	now := time.Now()
	duration, err := time.ParseDuration("2h45m35s")
	assert.Nil(t, err, fmt.Sprintf("Time parse error: %+v", err))
	future := now.Add(duration)
	strDuration := FmtDuration(future.Sub(now))
	assert.True(t, strings.HasPrefix(strDuration, "02:45:35"))
}
