import json
import os
import uuid
from datetime import datetime

import boto3
import pytest
from moto import mock_dynamodb2

from audit import CompanyAudit, ErrorType


def get_company_table(company_table, company_id):
    response = company_table.get_item(Key={"company_id": company_id})
    return response["Item"]


def test_missing_company_manager_id(audit_companies, dynamodb, company_table):
    """
    Test missing company_manager_id field
    """
    company_id = str(uuid.uuid4())
    table = dynamodb.Table("cla-test-companies")
    table.put_item(
        Item={"company_id": company_id,}
    )
    record = get_company_table(company_table, company_id)
    result = audit_companies.validate_company_manager_id(record)
    expected_result = {
        "company_id": company_id,
        "column": "company_manager_id",
        "is_valid": False,
        "error_type": ErrorType.INVALID,
        "data": None
    }

    assert result == expected_result


def test_invalid_user_manager_company_id(
    audit_companies, dynamodb, company_table, user_table
):
    """
    Test invalid company_manager_id that references invalid user_id
    """
    company_id = str(uuid.uuid4())
    user_id = str(uuid.uuid4())
    company_table.put_item(Item={"company_id": company_id, "company_manager_id": user_id})
    record = get_company_table(company_table,company_id)
    result = audit_companies.validate_company_manager_id(record)
    expected_result = {
        "company_id": company_id,
        "error_type": ErrorType.INVALID,
        "is_valid": False,
        "column": "company_manager_id",
        "data": user_id
    }
    assert result == expected_result
    user_id = str(uuid.uuid4())
    user_table.put_item(Item={"user_id": user_id})
    company_table.put_item(Item={
        "company_id":company_id,
        "company_manager_id": user_id
    })
    record = get_company_table(company_table,company_id)
    result = audit_companies.validate_company_manager_id(record)
    expected_result = {
        "company_id": company_id,
        "is_valid": True,
        "column": "company_manager_id",
    }

    assert result == expected_result


def test_company_name_missing_blank(
    audit_companies, dynamodb, company_table, user_table
):
    """
    Test company name missing or blank
    """
    company_id = str(uuid.uuid4())
    table = dynamodb.Table("cla-test-companies")
    table.put_item(
        Item={"company_id": company_id,}
    )
    response = table.get_item(Key={"company_id": company_id})
    record = response["Item"]
    result = audit_companies.validate_company_name(record)
    expected_result = {
        "company_id": company_id,
        "column": "company_name",
        "error_type": ErrorType.NULL_BLANK,
        "data": None,
        "is_valid": False,
    }
    assert result == expected_result


def test_missing_date_modified(audit_companies, dynamodb, company_table):
    company_id = str(uuid.uuid4())
    company_table.put_item(Item={"company_id": company_id})
    record = get_company_table(company_table,company_id)
    result = audit_companies.validate_date_modified(record)
    expected_result = {
        "company_id": company_id,
        "is_valid": False,
        "error_type": ErrorType.NULL,
        "data": None,
        "column": "date_modified"
    }

    assert result == expected_result

    company_table.put_item(Item = {"company_id": company_id, "date_modified": str(datetime.now())})
    record = get_company_table(company_table,company_id)
    result = audit_companies.validate_date_modified(record)
    expected_result = {
        "company_id": company_id,
        "column": "date_modified",
        "is_valid": True,
    }

    assert result == expected_result


def test_missing_date_created(audit_companies, dynamodb, company_table):
    """
    Test missing  date_created column
    """

    company_id = str(uuid.uuid4())
    # table = dynamodb.Table("cla-test_companies")
    company_table.put_item(Item={"company_id": company_id})
    record = get_company_table(company_table, company_id)
    result = audit_companies.validate_date_created(record)
    expected_result = {
        "company_id": company_id,
        "is_valid": False,
        "error_type": ErrorType.NULL,
        "column": "date_created",
        "data": None,
    }

    assert result == expected_result

    company_table.put_item(
        Item={"company_id": company_id, "date_created": str(datetime.now())}
    )

    record = get_company_table(company_table, company_id)

    result = audit_companies.validate_date_created(record)
    expected_result = {
        "company_id": company_id,
        "is_valid": True,
        "column": "date_created",
    }

    assert result == expected_result


