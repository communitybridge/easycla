// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

// TheLinuxFoundation is the name of the super parent for many Salesforce Foundations/Project Groups
const TheLinuxFoundation = "The Linux Foundation"

// XREQUESTID is the client request id - used to trace a client request through the system/logs
const XREQUESTID = "x-request-id"

// CLADesigneeRole CLA manager designee role identifier
const CLADesigneeRole = "cla-manager-designee"

// CLAManagerRole CLA manager role identifier
const CLAManagerRole = "cla-manager"

// CLASignatoryRole CLA signatory role identifier
const CLASignatoryRole = "cla-signatory"

// Lead representing type of user
const Lead = "lead"

// ProjectOrgScope is the ACS project + organiztion scope
const ProjectOrgScope = "project|organization"

// ClaTypeICLA represents individual contributor CLA records
const ClaTypeICLA = "icla"

// ClaTypeECLA represents employee contributor CLA records (acknowledgements)
const ClaTypeECLA = "ecla"

// ClaTypeCCLA represents corporate CLA records (includes approval lists)
const ClaTypeCCLA = "ccla"

// SignatureTypeCLA is the cla signature type in the DB
const SignatureTypeCLA = "cla"

// SignatureTypeCCLA is the ccla signature type in the DB
const SignatureTypeCCLA = "ccla"

// SignatureReferenceTypeUser is the signature reference type for user signatures - individual and employee
const SignatureReferenceTypeUser = "user"

// SignatureReferenceTypeCompany is the signature reference type for corporate signatures - signed by CLA Signatories, managed by CLA Managers
const SignatureReferenceTypeCompany = "company"

// ProjectTypeProjectGroup is the string that represents the Project Group type in a Project Service record
const ProjectTypeProjectGroup = "Project Group"
