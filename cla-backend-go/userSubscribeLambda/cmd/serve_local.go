// +build !aws_lambda

// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/gommon/log"
)

type fn func(ctx context.Context, sqsEvent events.SQSEvent) error

var handler fn

func postCollabSQSEvent(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.Method != http.MethodPost {
		log.Warn("404 not found.")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	dataByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn("Failed to read body")
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
	}

	var sEvent events.SQSEvent
	sEvent.Records = append(sEvent.Records, events.SQSMessage{
		Body:        string(dataByte),
		MessageId:   "LocalID",
		EventSource: "localhost",
	})

	if handler(r.Context(), sEvent) != nil {
		log.Fatal(err)
	}
}

// Start lambda handler function
func Start(hf fn) error {
	handler = hf
	http.HandleFunc("/", postCollabSQSEvent)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
	return nil
}
