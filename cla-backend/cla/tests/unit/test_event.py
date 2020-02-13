# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from cla.models.dynamo_models import Event, User, Project, Company
from cla.models import event_types
from unittest.mock import patch, Mock
import pytest
import datetime
import cla
import time


@pytest.fixture()
def mock_event():
    event = Event()
    event.model.save = Mock()
    yield event

def test_event_user_id(user_instance):
    """ Test event_user_id """
    Event.save = Mock()
    User.load = Mock()
    event_data = "test user id"
    response = Event.create_event(
        event_type=event_types.EventType.CreateProject,
        event_data=event_data,
        event_user_id=user_instance.get_user_id()
    )
    assert 'data' in response

def test_event_company_id(company):
    """ Test creation of event instance """
    #Case for creating Company
    Event.save = Mock()
    Company.load = Mock()
    event_data = 'test company created'
    response = Event.create_event(
        event_type=event_types.EventType.DeleteCompany,
        event_data=event_data,
        event_company_id=company.get_company_id()
    )
    assert 'data' in response

def test_event_project_id(project):
    """ Test event with event_project_id """
    Event.save = Mock()
    Project.load = Mock()
    event_data = 'project id loaded'
    response = Event.create_event(
        event_data=event_data,
        event_type=event_types.EventType.DeleteProject,
        event_project_id=project.get_project_id()
    )
    assert 'data' in response

def test_event_user_id_attribute(user_instance, mock_event):
    """ Test event_user_id attribute """
    mock_event.set_event_user_id(user_instance.get_user_id())
    mock_event.save()
    assert mock_event.get_event_user_id() == user_instance.get_user_id()

def test_event_company_name_lower_attribute(mock_event):
    """ Test company_name_lower attribute """
    mock_event.set_event_company_name("Company_lower")
    mock_event.save()
    assert mock_event.get_event_company_name_lower() == "company_lower"

def test_event_username_attribute(mock_event):
    """ Test event_username attribute """
    mock_event.set_event_user_name("foo_username")
    mock_event.save()
    assert mock_event.get_event_user_name() == "foo_username"

def test_event_user_name_lower_attribute(mock_event):
    """ Test event_user_name_lower attribute """
    mock_event.set_event_user_name("Username")
    mock_event.save()
    assert mock_event.get_event_user_name_lower() == "username"

def test_event_project_name_lower_attribute(mock_event):
    """ Test gettting project """
    mock_event.set_event_project_name("Project")
    mock_event.save()
    assert mock_event.get_event_project_name_lower() == "project"

def test_event_time(mock_event):
    """ Test event time  """
    mock_event.save()
    assert mock_event.get_event_time() <= datetime.datetime.now()

def test_event_time_epoch(mock_event):
    """ Test event time epoch """
    mock_event.save()
    assert mock_event.get_event_time_epoch() <= time.time()
