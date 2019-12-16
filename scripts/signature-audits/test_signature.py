# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import json
import os
import uuid

import boto3
import pytest
from moto import mock_dynamodb2

from audit import AuditSignature
from audit import ErrorType

@pytest.fixture(scope="function")
def aws_credentials():
    """Mocked aws credentials for moto"""
    os.environ["AWS_ACCESS_KEY_ID"] = "testing"
    os.environ["AWS_SECRET_ACCESS_KEY"] = "testing"
    os.environ["AWS_SECURITY_TOKEN"] = "testing"
    os.environ["AWS_SESSION_TOKEN"] = "testing"


@pytest.fixture(scope="function")
def dynamodb(aws_credentials):
    with mock_dynamodb2():
        session = boto3.Session()
        dynamodb = session.resource("dynamodb")
        yield dynamodb


@pytest.fixture(scope="function")
def signature_table(dynamodb):
    signature_table = dynamodb.create_table(
        TableName="cla-test-signatures",
        AttributeDefinitions=[
            {"AttributeName": "signature_id", "AttributeType": "S",},
        ],
        KeySchema=[{"AttributeName": "signature_id", "KeyType": "HASH"}],
        ProvisionedThroughput={"ReadCapacityUnits": 5, "WriteCapacityUnits": 5},
    )
    yield signature_table


@pytest.fixture(scope="function")
def company_table(dynamodb):
    company_table = dynamodb.create_table(
        TableName="cla-test-companies",
        AttributeDefinitions=[{"AttributeName": "company_id", "AttributeType": "S"}],
        KeySchema=[{"AttributeName": "company_id", "KeyType": "HASH"}],
        ProvisionedThroughput={"ReadCapacityUnits": 5, "WriteCapacityUnits": 5},
    )
    yield company_table


@pytest.fixture(scope="function")
def user_table(dynamodb):
    user_table = dynamodb.create_table(
        TableName="cla-test-users",
        AttributeDefinitions=[{"AttributeName": "user_id", "AttributeType": "S",}],
        KeySchema=[{"AttributeName": "user_id", "KeyType": "HASH"}],
        ProvisionedThroughput={"ReadCapacityUnits": 5, "WriteCapacityUnits": 5},
    )
    yield user_table


@pytest.fixture(scope="function")
def audit(dynamodb):
    audit = AuditSignature(dynamodb)
    audit.set_companies_table(dynamodb.Table("cla-test-companies"))
    audit.set_signatures_table(dynamodb.Table("cla-test-signatures"))
    audit.set_users_table(dynamodb.Table("cla-test-users"))
    yield audit


def get_signature_table(signature_table, signature_id):
    response = signature_table.get_item(Key={"signature_id": signature_id})
    return response["Item"]


def test_no_related_company_user_signature_reference_id(
    audit, dynamodb, signature_table, user_table
):
    """
    Test signature_reference_id referencing valid user record or company record
    """
    signature_id = str(uuid.uuid4())
    signature_reference_id = str(uuid.uuid4())
    table = dynamodb.Table("cla-test-signatures")
    table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_id": signature_reference_id,
        }
    )
    response = table.get_item(Key={"signature_id": signature_id})
    record = response["Item"]
    result = audit.validate_signature_reference_id(record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_reference_id",
        "data": signature_reference_id,
        "is_valid": False,
        "error_type": ErrorType.CCLA
    }

    assert result == expected_result


def test_valid_user_signature_reference_id(
    audit, dynamodb, signature_table, user_table
):
    """
    Test signature_reference_id with valid user id
    """
    signature_id = str(uuid.uuid4())
    user_id = str(uuid.uuid4())
    user_table = dynamodb.Table("cla-test-users")
    signature_table = dynamodb.Table("cla-test-signatures")
    user_table.put_item(Item={"user_id": user_id})
    signature_table.put_item(
        Item={"signature_id": signature_id, "signature_reference_id": user_id}
    )
    response = signature_table.get_item(Key={"signature_id": signature_id})

    record = response["Item"]
    result = audit.validate_signature_reference_id(record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_reference_id",
        "is_valid": True,
    }
    assert result == expected_result


def test_valid_company_signature_reference_id(
    audit, dynamodb, signature_table, company_table
):
    """
    Test signature_reference_id with valid company id
    """
    signature_id = str(uuid.uuid4())
    company_id = str(uuid.uuid4())

    company_table.put_item(Item={"company_id": company_id})
    signature_table.put_item(
        Item={"signature_id": signature_id, "signature_reference_id": company_id}
    )
    response = signature_table.get_item(Key={"signature_id": signature_id})
    record = response["Item"]
    result = audit.validate_signature_reference_id(record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_reference_id",
        "is_valid": True,
    }

    assert result == expected_result


def test_invalid_signature_reference_type(audit, signature_table):
    """
    Test that ensures records have a signature_reference_type
    of either user or company
    """
    invalid_value = "dummy_company"
    signature_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={"signature_id": signature_id, "signature_reference_type": invalid_value}
    )
    record = get_signature_table(signature_table, signature_id)
    result = audit.validate_signature_reference_type(record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_reference_type",
        "data": invalid_value,
        "is_valid": False,
        "error_type": ErrorType.CCLA
    }
    assert result == expected_result


def test_valid_signature_reference_type_company(audit, signature_table):
    valid_value = "company"
    signature_id_2 = str(uuid.uuid4())
    signature_table.put_item(
        Item={"signature_id": signature_id_2, "signature_reference_type": valid_value}
    )
    valid_record = get_signature_table(signature_table, signature_id_2)
    result = audit.validate_signature_reference_type(valid_record)
    expected_result = {
        "signature_id": signature_id_2,
        "column": "signature_reference_type",
        "is_valid": True,
    }

    assert result == expected_result


def test_valid_signature_reference_type_user(audit, signature_table):
    valid_value = "user"
    signature_id_2 = str(uuid.uuid4())
    signature_table.put_item(
        Item={"signature_id": signature_id_2, "signature_reference_type": valid_value}
    )
    valid_record = get_signature_table(signature_table, signature_id_2)
    result = audit.validate_signature_reference_type(valid_record)
    expected_result = {
        "signature_id": signature_id_2,
        "column": "signature_reference_type",
        "is_valid": True,
    }

    assert result == expected_result


def test_signature_employee_field(audit, signature_table):
    signature_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_type": "user",
            "signature_type": "cla",
        }
    )
    invalid_record = get_signature_table(signature_table, signature_id)
    result = audit.validate_employee_signature(invalid_record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_user_ccla_company_id",
        "is_valid": False,
        "error_type": ErrorType.INVALID,
        "data":None
    }

    assert result == expected_result


def test_valid_signature_employee_field(audit, signature_table, company_table):
    signature_id = str(uuid.uuid4())
    company_id = str(uuid.uuid4())
    company_table.put_item(Item={"company_id": company_id})
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_type": "user",
            "signature_type": "cla",
            "signature_user_ccla_company_id": company_id,
        }
    )
    invalid_record = get_signature_table(signature_table, signature_id)
    result = audit.validate_employee_signature(invalid_record)
    expected_result = {
        "signature_id": signature_id,
        "column": "signature_user_ccla_company_id",
        "is_valid": True,
    }

    assert result == expected_result


def test_icla_signature_field(audit, signature_table):
    signature_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_type": "user",
            "signature_type": "ccla",
            "signature_user_ccla_company_id": str(uuid.uuid4()),
        }
    )
    invalid_record = get_signature_table(signature_table, signature_id)
    expected_result = {
        "signature_id": signature_id,
        "is_valid": False,
        "column": "signature_user_ccla_company_id",
        "error_type": ErrorType.CCLA_USER,
        "data": None
    }
    result = audit.validate_icla_signature(invalid_record)
    assert result == expected_result


def test_ccla_signature_field(audit, signature_table):
    signature_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_type": "company",
            "signature_type": "ccla",
        }
    )
    expected_result = {
        "signature_id": signature_id,
        "is_valid": True,
        "column": "signature_user_ccla_company_id",
    }
    record = get_signature_table(signature_table, signature_id)
    result = audit.validate_ccla_signature(record)
    assert result == expected_result

    other_company_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_reference_type": "company",
            "signature_type": "ccla",
            "signature_user_ccla_company_id": other_company_id,
        }
    )
    record = get_signature_table(signature_table, signature_id)
    result = audit.validate_ccla_signature(record)
    expected = {
        "signature_id": signature_id,
        "is_valid": False,
        "column": "signature_user_ccla_company_id",
        "data": other_company_id,
        "error_type": ErrorType.CCLA_COMPANY,
    }

    assert result == expected


def test_valid_signature_document_version(audit, signature_table):
    signature_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={"signature_id": signature_id,}
    )
    invalid_record = get_signature_table(signature_table, signature_id)
    expected_result = {
        "signature_id": signature_id,
        "is_valid": False,
        "error_type": ErrorType.NULL,
        "column": "signature_document_major_version",
        "data":None,
    }

    result = audit.validate_signature_document_version(invalid_record)

    assert result == expected_result
    signature_id_2 = str(uuid.uuid4())
    signature_table.put_item(
        Item={"signature_id": signature_id_2, "signature_document_major_version": "v1"}
    )
    valid_record = get_signature_table(signature_table, signature_id_2)
    expected_result = {
        "signature_id": signature_id_2,
        "is_valid": True,
        "column": "signature_document_major_version",
    }
    result = audit.validate_signature_document_version(valid_record)

    assert expected_result == result


def test_signature_user_ccla_company_id_field(audit, signature_table, company_table):
    signature_id = str(uuid.uuid4())
    company_id = str(uuid.uuid4())
    signature_user_ccla_company_id = str(uuid.uuid4())
    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_user_ccla_company_id": signature_user_ccla_company_id,
        }
    )
    expected_result = {
        "signature_id": signature_id,
        "is_valid": False,
        "error_type": ErrorType.CCLA_COMPANY,
        "column": "signature_user_ccla_company_id",
        "data": signature_user_ccla_company_id
    }

    record = get_signature_table(signature_table, signature_id)
    result = audit.validate_signature_user_ccla_company(record)

    assert result == expected_result

    signature_table.put_item(
        Item={
            "signature_id": signature_id,
            "signature_user_ccla_company_id": company_id,
        }
    )
    company_table.put_item(
        Item = {
            'company_id':company_id
        }
    )

    expected_result = {
        "signature_id": signature_id,
        "is_valid": True,
        "column": "signature_user_ccla_company_id",
    }
    valid_record = get_signature_table(signature_table, signature_id)
    result_2 = audit.validate_signature_user_ccla_company(valid_record)

    assert result_2 == expected_result
