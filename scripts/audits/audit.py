# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Easyclav1 - Audits the Signatures tables
"""

import logging
from abc import ABC, abstractmethod
from enum import Enum
from typing import Dict

from botocore.exceptions import ClientError


class ErrorType(Enum):
    """
    Enumeration for error types for Signature records
    """

    ICLA = "ICLA"
    CCLA_COMPANY = "CCLA(company)"
    CCLA_USER = "CCLA(user)"
    CCLA = "CCLA"
    NULL = "NOT NULL"
    INVALID = "Invalid Column"
    NULL_BLANK = "NULL OR BLANK"


class Audit(ABC):
    """
    Abstract Audit class
    """

    def __init__(self, dynamo_db, batch=None):
        self.dynamo_db = dynamo_db
        self.batch = batch

    def process_batch(self):
        pass


class CompanyAudit(Audit):
    """
    Handles field validations for the company table
    """

    def __init__(self, dynamo_db, batch=None):
        self._dynamo_db = dynamo_db
        self._batch = batch
        self._companies_table = None
        self._users_table = None
        self._field_exists = {"company_name": False, "company_manager_id": False}

    def process_batch(self) -> Dict:
        """
        Function that processes batch list of records
        :return : A list of invalid audited records
        """
        audited_records = []

        for record in self._batch:
            if record:
                audited_records.append(self.validate_company_manager_id(record))
                audited_records.append(self.validate_company_name(record))
                audited_records.append(self.validate_date_created(record))
                audited_records.append(self.validate_date_modified(record))

        # Filter invalid audited records
        invalid_records = []
        for rec in audited_records:
            if rec:
                if not rec["is_valid"]:
                    invalid_records.append(rec)
        return invalid_records

    def get_companies_table(self):
        """
        Gets companies table
        """
        return self._companies_table

    def set_companies_table(self, companies_table):
        """
        Sets companies table
        """
        self._companies_table = companies_table

    def get_users_table(self):
        """
        Gets users table
        """
        return self._users_table

    def set_users_table(self, users_table):
        """
        Sets users table
        """
        self._users_table = users_table

    def validate_date_created(self, record):
        """
        Function that ensures the date_created column in companies table is not null
        """
        try:
            is_valid = False
            company_id = record["company_id"]
            record["date_created"]
            is_valid = True
        except KeyError:
            pass
        finally:
            result = {
                "company_id": company_id,
                "is_valid": is_valid,
                "column": "date_created",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_date_modified(self, record):
        """
        Function that ensures the date_modfied column is not null
        """
        try:
            is_valid = False
            company_id = record["company_id"]
            record["date_modified"]
            is_valid = True
        except KeyError:
            pass
        finally:
            result = {
                "company_id": company_id,
                "is_valid": is_valid,
                "column": "date_modified",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_company_name(self, record):
        """
        Function that ensures company_name is neither blank nor null
        """
        try:
            is_valid = False
            company_id = record["company_id"]
            company_name = record["company_name"]
            if company_name is not None:
                is_valid = True
        except KeyError:
            pass
        finally:
            result = {
                "company_id": company_id,
                "is_valid": is_valid,
                "column": "company_name",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL_BLANK
                result["data"] = None
        return result

    def validate_company_manager_id(self, record):
        """
        Function that ensures referenced user_id(company_manager_id) is valid
        """
        try:
            is_valid = False
            company_id = record["company_id"]
            company_manager_id = record["company_manager_id"]
            self._field_exists["company_manager_id"] = True
            response = self._users_table.get_item(Key={"user_id": company_manager_id})
            company_manager_id = response["Item"]
            is_valid = True
        except (ClientError, KeyError):
            pass
        finally:
            result = {
                "company_id": company_id,
                "column": "company_manager_id",
                "is_valid": is_valid,
            }
            if not is_valid:
                result["error_type"] = ErrorType.INVALID
                result["data"] = company_manager_id if self._field_exists["company_manager_id"] else None
        return result
