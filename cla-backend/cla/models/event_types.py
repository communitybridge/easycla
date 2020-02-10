# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from enum import Enum

class EventType(Enum):
    """
    Enumeraters representing type of CLA events
    across projects, users, signatures, whitelists

    """
    CreateUser = "Create User"
    UpdateUser = "Update User"
    DeleteUser = "Delete User"
    CreateProject = "Create Project"
    UpdateProject = "Update Project"
    DeleteProject = "Delete Project"
    CreateCompany = "Create Company"
    DeleteCompany = "Delete Company"
    UpdateCompany = "Update Company"
    CreateProjectDocument= "Create Project Document"
    CreateProjectDocumentTemplate = "Create Project Document with Template"
    DeleteProjectDocument = "Delete Project Document"
    AddPermission = "Add Permission"
    RemovePermission = "Remove Pemrission"
    AddProjectManager = "Add Project Manager"
    RemoveProjectManager = "Remove Project Manager"
    RequestCompanyWL = "Request Company Whitelist"
    InviteAdmin = "Invite Admin"
    RequestCCLA = "Request Company CCLA"
    RequestCompanyAdmin = "Request Company Admin access"
    AddCompanyPermission = "Add Company Permissions"
    RemoveCompanyPermission = "Remove Company Permissions"
    CreateSignature = "Create Signature"
    DeleteSignature = "Delete Signature"
    UpdateSignature = "Update Signature"
    AddCLAManager = "Add CLA Manager"
    RemoveCLAManager = "Remove CLA Manager"
    NotifyWLChange = "Notify WL Change"

