// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package organization_service

import "errors"

// These are Organization Service specific Error objects that match their code/impl
// so, instead of pulling in all their code from a private repo, we only need their
// error codes to deal with HTTP response errors
var (
	// ErrNotFound is a not found error
	ErrNotFound = errors.New("not found")
	// ErrInvalid is an invalid request error
	ErrInvalid = errors.New("invalid request")
	// ErrUnauthorized is an unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrConflict is an conflict/duplicate error
	ErrConflict = errors.New("conflict")
	// ErrForbidden is an forbidden error
	ErrForbidden = errors.New("forbidden")
	// ErrPendingOrg comes when status of an org is in pending
	ErrPendingOrg = errors.New("org_pending")
)
