// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package logging

import log "github.com/sirupsen/logrus"

// UTCFormatter structure for logging
type UTCFormatter struct {
	log.Formatter
}

// Format handler for UTC time - usage: log.SetFormatter(UTCFormatter{&log.JSONFormatter{}})
func (u UTCFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
