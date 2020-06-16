// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package events

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// CSVEventsResponse creates a new response handler for CSV files
func CSVEventsResponse(filename string, events *models.EventList) middleware.Responder {
	return &CSVEventsResponderFunc{
		Filename: filename,
		Events:   events,
	}
}

// CSVEventsResponderFunc wraps a func as a Responder interface
type CSVEventsResponderFunc struct {
	Filename string
	Events   *models.EventList
}

// WriteResponse writes to the response
func (fn CSVEventsResponderFunc) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fn.Filename))
	if err := pr.Produce(rw, fn.convertToCSV(fn.Events, ",")); err != nil {
		msg := fmt.Sprintf("issue converting event models to CSV format - error: %+v.", err)
		log.Warn(msg)
		rw.Header().Set(runtime.HeaderContentType, runtime.JSONMime)
		_, writeErr := rw.Write([]byte("{error=\"" + msg + "\"}"))
		if writeErr != nil {
			log.Warnf("issue writing error response response - error: %+v", writeErr)
		}
	}
}

func (fn CSVEventsResponderFunc) convertToCSV(events *models.EventList, delimiter string) []string {
	var csvLines []string
	var b bytes.Buffer
	for _, event := range events.Events {
		b.WriteString(event.UserName)
		b.WriteString(delimiter)
		csvLines = append(csvLines, b.String())
		b.Reset()
	}

	return csvLines
}
