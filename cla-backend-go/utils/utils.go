// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// FmtDuration is a helper function to format a duration in the HH:MM:SS.sss format
func FmtDuration(d time.Duration) string {
	//days := int64(d.Hours() / 24)
	hours := int64(math.Mod(d.Hours(), 24))
	minutes := int64(math.Mod(d.Minutes(), 60))
	seconds := int64(math.Mod(d.Seconds(), 60))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, d.Milliseconds())
}

// TimeToString returns time in the RFC3339 format
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// CurrentTime returns the current UTC time and current Time in the RFC3339 format
func CurrentTime() (time.Time, string) {
	t := time.Now().UTC()
	return t, TimeToString(t)
}

// ParseDateTime attempts to convert the string to one of our supported date time formats
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	dateTimeStrTrimmed := strings.TrimSpace(dateTimeStr)

	supportedFormats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05Z-07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000000-0700",
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
	}

	for _, format := range supportedFormats {
		dateTimeValue, err := time.Parse(format, dateTimeStrTrimmed)
		if err == nil {
			return dateTimeValue, err
		}
	}

	return time.Unix(0, 0), fmt.Errorf("unsupported date/time format: %s", dateTimeStrTrimmed)
}

// AddStringAttribute adds string attribute to dynamodb input map
func AddStringAttribute(item map[string]*dynamodb.AttributeValue, key string, value string) {
	if value != "" {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}
}

// AddNumberAttribute adds number attribute to dynamodb input map
func AddNumberAttribute(item map[string]*dynamodb.AttributeValue, key string, value int64) {
	numString := strconv.FormatInt(value, 10)
	item[key] = &dynamodb.AttributeValue{N: aws.String(numString)}
}

// StringInSlice returns true if the specified string value exists in the slice, otherwise returns false
func StringInSlice(a string, list []string) bool {
	if list == nil {
		return false
	}

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// HostInSlice returns true if the specified host value exists in the slice, otherwise returns false
func HostInSlice(a string, list []string) bool {
	if list == nil {
		return false
	}

	for _, b := range list {
		b = strings.Split(b, ":")[0]
		if b == a {
			return true
		}
	}
	return false
}
