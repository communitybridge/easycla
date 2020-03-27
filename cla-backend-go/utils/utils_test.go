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

func TestParseDateTimeDefault(t *testing.T) {
	validInput := []string{
		"2020-03-27T20:19:07+08:00",
		"2020-03-27T12:00:17+00:00",
		"2013-01-31T04:00:17-05:00",
	}

	for _, dateTimeStr := range validInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.Nil(t, err, fmt.Sprintf("Check for Valid Date Time: %s - %+v", dateTimeStr, dateTimeValue))
		assert.True(t, dateTimeValue.After(time.Unix(0, 0)))
	}

	inValidInput := []string{
		// Note: Formatter doesn't seem to catch invalid time zone values
		"20201-11-27T12:00:17+00:00", // Invalid Year
		"2020-13-27T12:00:17+00:00",  // Invalid Month
		"2013-01-32T04:00:17-05:00",  // Invalid Day
		"2013-01-02T24:01:17-05:00",  // Invalid Hour
		"2013-01-02T14:71:17-05:00",  // Invalid Minute
		"2013-01-02T14:21:87-05:00",  // Invalid Second
	}

	for _, dateTimeStr := range inValidInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.NotNil(t, err, fmt.Sprintf("Check for Invalid Date Time: %s - %v", dateTimeStr, dateTimeValue))
	}
}

func TestParseDateTimeZulu(t *testing.T) {
	validInput := []string{
		"2020-03-27T15:04:05Z+07:00",
		"2020-09-02T15:04:05Z-07:00",
		"2016-01-02T15:04:05Z-07:00",
	}

	for _, dateTimeStr := range validInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.Nil(t, err, fmt.Sprintf("Check for Valid Date Time: %s - %+v", dateTimeStr, dateTimeValue))
		assert.True(t, dateTimeValue.After(time.Unix(0, 0)))
	}

	inValidInput := []string{
		// Note: Formatter doesn't seem to catch invalid time zone values
		"20016-01-02T15:04:05Z-07:00", // Invalid Year
		"2006-21-02T15:04:05Z-07:00",  // Invalid Month
		"22006-01-92T15:04:05Z-07:00", // Invalid Day
		"22006-01-92T25:04:05Z-07:00", // Invalid Hour
		"22006-01-92T15:74:05Z-07:00", // Invalid Minute
		"22006-01-92T25:04:85Z-07:00", // Invalid Second
	}

	for _, dateTimeStr := range inValidInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.NotNil(t, err, fmt.Sprintf("Check for Invalid Date Time: %s - %+v", dateTimeStr, dateTimeValue))
	}
}

func TestParseDateTimeMS(t *testing.T) {
	validInput := []string{
		"2020-03-27T15:04:05.000000+0000",
		"2016-12-02T05:14:05.000000+0800",
		"2006-08-31T10:24:05.000000-1000",
	}

	for _, dateTimeStr := range validInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.Nil(t, err, fmt.Sprintf("Check for Valid Date Time: %s - %+v", dateTimeStr, dateTimeValue))
		assert.True(t, dateTimeValue.After(time.Unix(0, 0)))
	}

	inValidInput := []string{
		// Note: Formatter doesn't seem to catch invalid time zone values
		"20062-01-02T15:04:05.000000+0000", // Invalid Year
		"2006-21-02T15:04:05.000000+0000",  // Invalid Month
		"2006-01-92T15:04:05.000000+0000",  // Invalid Day
		"2020-03-27T35:04:05.000000+0000",  // Invalid Hour
		"2020-03-27T15:94:05.000000+0000",  // Invalid Minute
		"2020-03-27T15:04:85.000000+0000",  // Invalid Second
	}

	for _, dateTimeStr := range inValidInput {
		dateTimeValue, err := ParseDateTime(dateTimeStr)
		assert.NotNil(t, err, fmt.Sprintf("Check for Invalid Date Time: %s - %+v", dateTimeStr, dateTimeValue))
	}
}
