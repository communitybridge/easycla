// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import "fmt"

// SFProjectNotFound is an error model for Salesforce Project not found errors
type SFProjectNotFound struct {
	ProjectSFID string
	Err         error
}

// Error is an error string function for Salesforce Project not found errors
func (e *SFProjectNotFound) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("salesforce project %s not found", e.ProjectSFID)
	}
	return fmt.Sprintf("salesforce project %s not found: %+v", e.ProjectSFID, e.Err)
}

// Unwrap method returns its contained error
func (e *SFProjectNotFound) Unwrap() error {
	return e.Err
}

// CLAGroupNotFound is an error model for CLA Group not found errors
type CLAGroupNotFound struct {
	CLAGroupID string
	Err        error
}

// Error is an error string function for CLA Group not found errors
func (e *CLAGroupNotFound) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("cla group %s not found", e.CLAGroupID)
	}
	return fmt.Sprintf("cla group %s not found: %+v", e.CLAGroupID, e.Err)
}

// Unwrap method returns its contained error
func (e *CLAGroupNotFound) Unwrap() error {
	return e.Err
}

// CLAGroupICLANotConfigured is an error model for CLA Group ICLA not configured
type CLAGroupICLANotConfigured struct {
	CLAGroupID   string
	CLAGroupName string
	Err          error
}

// Error is an error string function for CLA Group ICLA not configured
func (e *CLAGroupICLANotConfigured) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("cla group %s (%s) is not configured for ICLAs", e.CLAGroupName, e.CLAGroupID)
	}
	return fmt.Sprintf("cla group %s (%s) is not configured for ICLAs: %+v", e.CLAGroupName, e.CLAGroupID, e.Err)
}

// Unwrap method returns its contained error
func (e *CLAGroupICLANotConfigured) Unwrap() error {
	return e.Err
}

// CLAGroupCCLANotConfigured is an error model for CLA Group CCLA not configured
type CLAGroupCCLANotConfigured struct {
	CLAGroupID   string
	CLAGroupName string
	Err          error
}

// Error is an error string function for CLA Group CCLA not configured
func (e *CLAGroupCCLANotConfigured) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("cla group %s (%s) is not configured for CCLAs", e.CLAGroupName, e.CLAGroupID)
	}
	return fmt.Sprintf("cla group %s (%s) is not configured for CCLAs: %+v", e.CLAGroupName, e.CLAGroupID, e.Err)
}

// Unwrap method returns its contained error
func (e *CLAGroupCCLANotConfigured) Unwrap() error {
	return e.Err
}

// ProjectCLAGroupMappingNotFound is an error model for project CLA Group not found errors
type ProjectCLAGroupMappingNotFound struct {
	ProjectSFID string
	CLAGroupID  string
	Err         error
}

// Error is an error string function for project CLA Group not found errors
func (e *ProjectCLAGroupMappingNotFound) Error() string {
	if e.CLAGroupID == "" {
		if e.Err == nil {
			return fmt.Sprintf("project %s to cla group mapping not found", e.ProjectSFID)
		}
		return fmt.Sprintf("project %s cla group mapping not found: %+v", e.ProjectSFID, e.Err)
	}
	if e.ProjectSFID == "" {
		if e.Err == nil {
			return fmt.Sprintf("project to cla group %s mapping not found", e.CLAGroupID)
		}
		return fmt.Sprintf("project cla group %s mapping not found: %+v", e.CLAGroupID, e.Err)
	}

	return fmt.Sprintf("project %s cla group %s mapping not found: %+v", e.ProjectSFID, e.CLAGroupID, e.Err)
}

// Unwrap method returns its contained error
func (e *ProjectCLAGroupMappingNotFound) Unwrap() error {
	return e.Err
}
