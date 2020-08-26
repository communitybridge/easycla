// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// WatermarkPdf adds watermark text to the given pdf blob and returns back the new one
func WatermarkPdf(pdf []byte, text string) ([]byte, error) {
	readSeek := bytes.NewReader(pdf)
	var b bytes.Buffer
	outWriter := bufio.NewWriter(&b)

	// this means it's a watermark
	onTop := false
	wm, err := pdfcpu.ParseTextWatermarkDetails(text, "points:48, scale:1, op:.4", onTop)
	if err != nil {
		return nil, err
	}

	err = api.AddWatermarks(readSeek, outWriter, nil, wm, nil)
	if err != nil {
		return nil, fmt.Errorf("applying watermark failed : %w", err)
	}

	return b.Bytes(), nil
}
