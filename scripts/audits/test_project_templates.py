# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import uuid


import pytest

from audit import ErrorType


def get_table(table_name, id):
    """
    Utility function that gets given table by its Key
    """
    response = table_name.get_item(Key={"project_id": id})
    record = response["Item"]
    return record


def test_missing_individual_documents(audit_projects, project_table):
    """
    Function to test missing project_individual_documents
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(Item={"project_id": project_id, "project_individual_documents": []})
    record = get_table(project_table, project_id)
    result = audit_projects.validate_project_individual_document(record)
    expected_result = {
        "project_id": project_id,
        "column": "project_individual_documents",
        "is_valid": False,
        "data": None,
        'error_type': ErrorType.NULL
    }
    assert result == expected_result


def test_missing_individual_document_s3_url(audit_projects, project_table):
    """
    Function to test project_individual_documents_s3_url for missing links
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(
        Item={
            "project_id": project_id,
            "project_individual_documents": [{"M": {"document_author_name": {"S": "Apache_Style"}}}],
        }
    )
    record = get_table(project_table, project_id)
    result = audit_projects.validate_individual_s3_url(record)
    expected_result = {
        "project_id": project_id,
        "error_type": ErrorType.MISSING_LINK,
        "data": None,
        "column": "project_individual_documents",
        "is_valid": False,
    }
    assert result == expected_result


def test_invalid_individual_document_s3_url(audit_projects, project_table):
    """
    Function to test if document_s3_url is valid
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(
        Item={
            "project_id": project_id,
            "project_individual_documents": [{"M": {"document_s3_url": {"S": "https://google.com"}}}],
        }
    )
    record = get_table(project_table, project_id)
    s3_url = record["project_individual_documents"][0]["M"]["document_s3_url"]["S"]
    result = audit_projects.validate_individual_s3_url(record)
    expected_result = {
        "is_valid": False,
        "project_id": project_id,
        "error_type": ErrorType.INVALID_LINK,
        "data": s3_url,
        "column": "project_individual_documents",
    }
    assert result == expected_result


def test_missing_project_corporate_documents(audit_projects, project_table):
    """
    Function to test missing project_corporate_document
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(Item={"project_id": project_id, "project_corporate_documents": []})
    record = get_table(project_table, project_id)
    result = audit_projects.validate_project_corporate_document(record)
    expected_result = {
        "project_id": project_id,
        "column": "project_corporate_documents",
        "is_valid": False,
        "data": None,
        "error_type": ErrorType.NULL,
    }
    assert result == expected_result


def test_missing_document_s3_url(audit_projects, project_table):
    """
    Function to test project_corporate_documents
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(
        Item={
            "project_id": project_id,
            "project_corporate_documents": [{"M": {"document_author_name": {"S": "Apache_Style"}}}],
        }
    )
    record = get_table(project_table, project_id)
    result = audit_projects.validate_s3_url(record)
    expected_result = {
        "project_id": project_id,
        "error_type": ErrorType.MISSING_LINK,
        "data": None,
        "column": "project_corporate_documents",
        "is_valid": False,
    }
    assert result == expected_result


def test_invalid_document_s3_url(audit_projects, project_table):
    """
    Function to test if document_s3_url is valid
    """
    project_id = str(uuid.uuid4())
    project_table.put_item(
        Item={
            "project_id": project_id,
            "project_corporate_documents": [{"M": {"document_s3_url": {"S": "https://google.com"}}}],
        }
    )
    record = get_table(project_table, project_id)
    s3_url = record["project_corporate_documents"][0]["M"]["document_s3_url"]["S"]
    result = audit_projects.validate_s3_url(record)
    expected_result = {
        "is_valid": False,
        "project_id": project_id,
        "error_type": ErrorType.INVALID_LINK,
        "data": s3_url,
        "column": "project_corporate_documents",
    }
    assert result == expected_result
