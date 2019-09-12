// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"
	"math"
	"time"
)

// FmtDuration is a helper function to format a duration in the HH:MM:SS.sss format
func FmtDuration(d time.Duration) string {
	//days := int64(d.Hours() / 24)
	hours := int64(math.Mod(d.Hours(), 24))
	minutes := int64(math.Mod(d.Minutes(), 60))
	seconds := int64(math.Mod(d.Seconds(), 60))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, d.Milliseconds())
}
