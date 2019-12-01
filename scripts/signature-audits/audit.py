# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Easyclav1 - Audits the Signatures tables
"""

import logging
from typing import Dict
from enum import Enum

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
    INVALID_SIG = "Invalid Signature"


class AuditSignature:
    """
    Handles field validations for signatures
    """

    def __init__(self, dynamo_db, batch=None):
        self._dynamo_db = dynamo_db
        self._signatures_table = None
        self._users_table = None
        self._companies_table = None
        self._batch = batch
        self._field_exists = {
            "signature_document_major_version": False,
            "signature_user_ccla_company_id": False,
            "signature_reference_type": False,
            "signature_type": False,
            "signature_reference_id": False,
        }

    def get_signatures_table(self):
        """
        Get signatures tabble(dynamodb)
        """
        return self._signatures_table

    def get_users_table(self):
        """
        Get users table(dynamodb)
        """
        return self._users_table

    def get_companies_table(self):
        """
        Get companies table(dynamodb)
        """
        return self._companies_table

    def set_signatures_table(self, signatures_table):
        """
        Set signatures table
        """
        self._signatures_table = signatures_table

    def set_users_table(self, users_table):
        """
        Set users table
        """
        self._users_table = users_table

    def set_companies_table(self, companies_table):
        """
        Set companies table
        """
        self._companies_table = companies_table

    def process_batch(self) -> Dict:
        """
        Function that process batch list of records
        :return : A list of invalid audited_records dicts
        """
        audited_records = []

        for record in self._batch:
            if record:
                audited_records.append(self.validate_ccla_signature(record))
                audited_records.append(self.validate_employee_signature(record))
                audited_records.append(self.validate_icla_signature(record))
                audited_records.append(self.validate_signature_reference_type(record))
                audited_records.append(self.validate_signature_signed(record))
                audited_records.append(self.validate_signature_type(record))
                audited_records.append(self.validate_signature_version(record))
                audited_records.append(self.validate_signature_approved(record))
                audited_records.append(
                    self.validate_signature_user_ccla_company(record)
                )
                audited_records.append(self.validate_signature_document_version(record))
                audited_records.append(self.validate_signature_reference_id(record))

        # filter invalid audited_records
        invalid_records = []
        for rec in audited_records:
            if rec:
                if not rec["is_valid"]:
                    invalid_records.append(rec)

        return invalid_records

    def validate_signature_document_version(self, record: Dict) -> Dict:
        """
        Function that ensures signature_document_version field is not null
        """
        try:
            is_valid = False
            signature_id = record["signature_id"]
            version = record["signature_document_major_version"]
            self._field_exists["signature_document_major_version"] = True
            is_valid = True
        except KeyError:
            pass
        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_document_major_version",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_signature_version(self, record: Dict) -> Dict:
        """
        Function that validates signature_version field ensuring
        field is not null
        :param record: signature record to be validated
        """
        try:
            is_valid = False
            signature_id = record["signature_id"]
            record["version"]
            is_valid = True

        except KeyError:
            pass
        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "version",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
            return result

    def validate_signature_signed(self, record: Dict) -> Dict:
        """
        Function that ensures validate_signature_signed is not null
        """
        try:
            is_valid = False
            signature_id = record["signature_id"]
            record["signature_signed"]
            is_valid = True

        except KeyError:
            pass
        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_signed",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_signature_approved(self, record: Dict) -> Dict:
        try:
            is_valid = False
            signature_id = record["signature_id"]
            record["signature_approved"]
            is_valid = True

        except KeyError:
            pass
        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_approved",
            }
            if not is_valid:
                result["error_type"] = ErrorType.NULL
                result["data"] = None
        return result

    def validate_signature_user_ccla_company(self, record: Dict) -> Dict:
        """
        Function that ensures signature_user_ccla_company relates to a valid
        company
        """
        try:
            is_valid = False
            signature_id = record["signature_id"]
            signature_user_ccla_company = record["signature_user_ccla_company_id"]
            try:
                company = self._companies_table.get_item(
                    Key={"company_id": signature_user_ccla_company}
                )
            except ClientError as err:
                logging.error(err.response["Error"]["Message"])
            else:
                if company and "Item" in company:
                    is_valid = True

        except (KeyError) as err:
            is_valid = True
            logging.info(err)
        finally:
            result = {
                "is_valid": is_valid,
                "signature_id": signature_id,
                "column": "signature_user_ccla_company_id",
            }
            if not result["is_valid"]:
                result["error_type"] = ErrorType.CCLA_COMPANY
                result["data"] = signature_user_ccla_company
        return result

    def validate_employee_signature(self, record: Dict):
        try:
            signature_id = record["signature_id"]
            is_valid = False
            try:
                signature_reference_type = record["signature_reference_type"]
                signature_type = record["signature_type"]
            except KeyError:
                return

            if signature_reference_type == "user" and signature_type == "cla":
                signature_user_ccla_company_id = record[
                    "signature_user_ccla_company_id"
                ]
                self._field_exists["signature_user_ccla_company_id"] = True
                # resolve valid company
                try:
                    company = self._companies_table.get_item(
                        Key={"company_id": signature_user_ccla_company_id}
                    )
                except ClientError as err:
                    logging.error(err.response["Error"]["Message"])
                else:
                    if company and "Item" in company:
                        is_valid = True

        except (Exception, ClientError) as err:
            logging.error(err)
        finally:
            result = {
                "signature_id": signature_id,
                "column": "signature_user_ccla_company_id",
                "is_valid": is_valid,
            }
            if not is_valid:
                result["error_type"] = ErrorType.ICLA
                result["data"] = (
                    signature_user_ccla_company_id
                    if self._field_exists["signature_user_ccla_company_id"]
                    else None
                )

        return result

    def validate_ccla_signature(self, record: Dict):
        try:
            is_valid = True
            signature_id = record["signature_id"]
            signature_reference_type = record["signature_reference_type"]
            self._field_exists["signature_reference_type"] = True
            signature_type = record["signature_type"]
            self._field_exists["signature_type"] = True
            if signature_reference_type == "company" and signature_type == "ccla":
                try:
                    company_id = record["signature_user_ccla_company_id"]
                    self._field_exists["signature_user_ccla_company_id"] = True
                    is_valid = False
                except KeyError:
                    # If Key error then record is valid
                    pass

        except KeyError:
            pass

        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_user_ccla_company_id",
            }
            if not is_valid:
                result["error_type"] = ErrorType.CCLA_COMPANY
                result["data"] = (
                    company_id
                    if self._field_exists["signature_user_ccla_company_id"]
                    else None
                )
        return result

    def validate_icla_signature(self, record: Dict):
        try:
            is_valid = True
            signature_id = record["signature_id"]
            signature_reference_type = record["signature_reference_type"]
            signature_type = record["signature_type"]
            if signature_reference_type == "user" and signature_type == "ccla":
                try:
                    record["signature_user_ccla_company_id"]
                    is_valid = False
                except KeyError:
                    # If Key error then record is valid
                    pass

        except KeyError as err:
            pass

        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_user_ccla_company_id",
            }
            if not is_valid:
                result["error_type"] = ErrorType.CCLA_USER
                result["data"] = None
        return result

    def validate_signature_type(self, record: Dict) -> Dict:
        try:
            is_valid = False
            valid_signature_type = ["cla", "ccla"]
            signature_id = record["signature_id"]
            signature_type = record["signature_type"]
            self._field_exists["signature_type"] = True
            is_valid = signature_type in valid_signature_type

        except KeyError:
            pass

        finally:
            result = {
                "signature_id": signature_id,
                "column": "signature_type",
                "data": signature_type
                if self._field_exists["signature_type"]
                else None,
                "is_valid": is_valid,
            }
            if not is_valid:
                result["error_type"] = ErrorType.INVALID_SIG
        return result

    def validate_signature_reference_type(self, record: Dict) -> Dict:
        """
        Function that checks the signature_reference type column and ensures
        values in (company,user)
        """
        try:
            is_valid = False
            valid_signature_reference_type = ["company", "user"]
            signature_id = record["signature_id"]
            signature_reference_type = record["signature_reference_type"]
            self._field_exists["signature_reference_type"] = True
            is_valid = signature_reference_type in valid_signature_reference_type

        except (KeyError) as err:
            pass

        finally:
            result = {
                "signature_id": signature_id,
                "column": "signature_reference_type",
                "is_valid": is_valid,
            }
            if not is_valid:
                result["data"] = (
                    signature_reference_type
                    if self._field_exists["signature_reference_type"]
                    else None
                )
                result["error_type"] = ErrorType.CCLA
        return result

    def validate_signature_reference_id(self, record: Dict) -> Dict:
        """
        Function that validates signature_reference_id column
        :param record: A signatures table record parsed
        """

        try:
            is_valid = False
            signature_id = record["signature_id"]
            signature_reference_id = record["signature_reference_id"]
            self._field_exists["signature_reference_id"] = True
            try:
                company = self._companies_table.get_item(
                    Key={"company_id": signature_reference_id}
                )
            except ClientError as err:
                logging.error(err.response["Error"]["Message"])
            else:
                if company and "Item" in company:
                    is_valid = True
            try:
                user = self._users_table.get_item(
                    Key={"user_id": signature_reference_id}
                )
            except ClientError as err:
                logging.error(err.response["Error"]["Message"])
            else:
                if user and "Item" in user:
                    is_valid = True

        except KeyError as err:
            logging.error(err)

        finally:
            result = {
                "signature_id": signature_id,
                "is_valid": is_valid,
                "column": "signature_reference_id",
            }
            if not is_valid:
                result["error_type"] = ErrorType.CCLA
                result["data"] = (
                    signature_reference_id
                    if self._field_exists["signature_reference_id"]
                    else None
                )
        return result
