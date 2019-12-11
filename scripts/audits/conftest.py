# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT


import os

import boto3
import pytest
from moto import mock_dynamodb2

from audit import CompanyAudit, ProjectAudit


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
def project_table(dynamodb):
    project_table = dynamodb.create_table(
        TableName="cla-test-projects",
        AttributeDefinitions=[{"AttributeName": "project_id", "AttributeType": "S"}],
        KeySchema=[{"AttributeName": "project_id", "KeyType": "HASH"}],
        ProvisionedThroughput={"ReadCapacityUnits": 5, "WriteCapacityUnits": 5},
    )
    yield project_table


@pytest.fixture(scope="function")
def audit_companies(dynamodb):
    audit = CompanyAudit(dynamodb)
    audit.set_companies_table(dynamodb.Table("cla-test-companies"))
    audit.set_users_table(dynamodb.Table("cla-test-users"))
    yield audit

@pytest.fixture(scope="function")
def audit_projects(dynamodb):
    audit = ProjectAudit(dynamodb)
    audit.set_projects_table(dynamodb.Table("cla-test-projects"))
    yield audit
