# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from abc import ABC, abstractmethod
from enum import Enum
from mimetypes import MimeTypes
from urllib.error import URLError
from urllib.request import Request, urlopen


class ErrorType(Enum):
    """
    Enumeration for error types for Signature records
    """

    ICLA = "ICLA"
    CCLA_COMPANY = "CCLA(company)"
    CCLA_USER = "CCLA(user)"
    CCLA = "CCLA"
    NULL = "NOT NULL"
    INVALID_SIG = "Invalid Signature"
    MISSING_LINK = "Missing Link"
    INVALID_LINK = "Invalid Link"


class Audit(ABC):
    def __init__(self, dynamodb, batch=None):
        self.dynamodb = dynamodb
        self.batch = batch

    @abstractmethod
    def process_batch(self):
        pass


class ProjectAudit(Audit):
    def __init__(self, dynamodb, batch=None):
        self.dynamodb = dynamodb
        self.batch = batch
        self._projects_table = None
        self._field_exists = {"individual_s3_url": False, "corporate_s3_url": False}

    def get_projects_table(self):
        """
        Gets dynamodb projects table
        """
        return self._projects_table

    def set_projects_table(self, projects_table):
        """
        Sets projects table
        """
        self._projects_table = projects_table

    def process_batch(self):
        """
        Function that processes batch list of project table records for invalid
        s3_url links for individuals and companies
        """
        audited_records = []
        invalid_records = []
        for record in self.batch:
            if record:
                audited_records.append(self.validate_individual_s3_url(record))
                audited_records.append(self.validate_project_corporate_document(record))
                audited_records.append(self.validate_project_individual_document(record))
                audited_records.append(self.validate_s3_url(record))

        for rec in audited_records:
            if rec:
                if not rec["is_valid"]:
                    invalid_records.append(rec)

        return invalid_records

    def validate_project_corporate_document(self, record):
        """
        Function that ensures project_corporate_document is not null
        """
        try:
            is_valid = False
            project_id = record["project_id"]
            project_corporate_documents = record["project_corporate_documents"]
            if project_corporate_documents:
                is_valid = True

        except KeyError:
            pass
        finally:
            result = {"is_valid": is_valid, "column": "project_corporate_documents"}
            if not is_valid:
                result["project_id"] = project_id
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_s3_url(self, record):
        """
        Function that validates project_document_s3_url
        """
        try:
            is_valid = False
            missing = False
            project_id = record["project_id"]
            try:
                project_corporate_documents = record["project_corporate_documents"]
                if project_corporate_documents:
                    s3_url = project_corporate_documents[0]["M"]["document_s3_url"]["S"]
                    self._field_exists["corporate_s3_url"] = True
                    mime = MimeTypes()
                    mime_type, _ = mime.guess_type(s3_url)
                    response = urlopen(s3_url)
                    if response.getcode() == 200 and mime_type == "application/pdf":
                        is_valid = True

            except KeyError:
                missing = True
        except URLError:
            pass
        finally:
            result = {
                "is_valid": is_valid,
                "project_id": project_id,
                "column": "project_corporate_documents",
            }
            if not is_valid:
                result["error_type"] = ErrorType.MISSING_LINK if missing else ErrorType.INVALID_LINK
                result["data"] = s3_url if self._field_exists["corporate_s3_url"] else None

        return result

    def validate_project_individual_document(self, record):
        """
        Function that checks if project_individual_document exists
        """
        try:
            is_valid = False
            project_id = record["project_id"]
            project_individual_documents = record["project_individual_documents"]
            if project_individual_documents:
                is_valid = True
        except KeyError:
            pass
        finally:
            result = {"is_valid": is_valid, "column": "project_individual_documents"}
            if not is_valid:
                result["project_id"] = project_id
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_individual_s3_url(self, record):
        """
        Function that checks if a individual_document_s3_url is valid
        """
        try:
            is_valid = False
            missing = False
            project_id = record["project_id"]
            try:
                project_individual_documents = record["project_individual_documents"]
                if project_individual_documents:
                    s3_url = project_individual_documents[0]["M"]["document_s3_url"]["S"]
                    self._field_exists["individual_s3_url"] = True
                    mime = MimeTypes()
                    mime_type, _ = mime.guess_type(s3_url)
                    response = urlopen(s3_url)
                    if response.getcode() == 200 and mime_type == "application/pdf":
                        is_valid = True
            except KeyError:
                missing = True
        except (Exception,URLError):
            pass
        finally:
            result = {
                "is_valid": is_valid,
                "project_id": project_id,
                "column": "project_individual_documents",
            }
            if not is_valid:
                result["error_type"] = ErrorType.MISSING_LINK if missing else ErrorType.INVALID_LINK
                result["data"] = s3_url if self._field_exists["individual_s3_url"] else None

        return result
