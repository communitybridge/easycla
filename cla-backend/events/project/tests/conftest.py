import os
import pytest
import boto3
from botocore.errorfactory import ClientError
import logging

from project import create_app
from cla import config
from cla.models.dynamo_models import EventModel

AWS_REGION = os.environ.get("AWS_REGION")


@pytest.fixture(scope="module")
def test_app():
    """
    Flask app for testing events_log api functionality
    """
    app = create_app()
    app.config.from_object("project.config.TestingConfig")
    with app.app_context():
        yield app


@pytest.fixture(scope="module")
def events_table():
    """
    Fixture that creates events table for test purporses
    """
    try:
        table = EventModel.create_table(
            read_capacity_units=config.DYNAMO_READ_UNITS, write_capacity_units=config.DYNAMO_WRITE_UNITS
        )
        yield table
    except ClientError as err:
        logging.error(err)
