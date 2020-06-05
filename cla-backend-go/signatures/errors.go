// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package signatures

// NewBadRequestError returns an error that formats as the given text.
func NewBadRequestError(text string) error {
	return &BadRequestError{text}
}

// BadRequestError is a trivial implementation of error.
type BadRequestError struct {
	s string
}

// Error is the to string method for an error
func (e BadRequestError) Error() string {
	return e.s
}

// NewForbiddenError returns an error that formats as the given text.
func NewForbiddenError(text string) error {
	return &ForbiddenError{text}
}

// ForbiddenError is a trivial implementation of error.
type ForbiddenError struct {
	s string
}

// Error is the to string method for an error
func (e ForbiddenError) Error() string {
	return e.s
}
