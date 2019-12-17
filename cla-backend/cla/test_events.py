# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import uuid
import json

import boto3
import hug
import pytest
from falcon import HTTP_200, HTTP_400, HTTP_404

import cla
from cla import routes
from cla.utils import (
    get_event_instance,
    get_project_instance,
    get_company_instance,
    get_user_instance,
)
from datetime import datetime


def add_event(event_id):
    """
    Utility function that adds an event
    """
    event = get_event_instance()
    event.set_event_id(event_id)
    event.save()
    return event


def test_get_events():
    """
    Test getting events
    """
    response = hug.test.get(routes, "/v1/events")
    assert response.status == HTTP_200
    assert response.data["events"]


def test_get_event_by_event_id():
    """
    Test get event given an event_id
    """
    # add event
    event_id = str(uuid.uuid4())
    add_event(event_id)
    # hit endpoint
    url = "/v1/events/{}".format(event_id)
    response = hug.test.get(routes, url)
    assert response.status == HTTP_200
    assert response.data["event_id"] == event_id

    url = "/v1/events/{}".format("123123")
    response = hug.test.get(routes, url)
    assert response.status == HTTP_404


def test_create_event():
    """
    Test endpoint that creates event instance
    """
    url = "/v1/events"
    event_project_id = str(uuid.uuid4())
    event_company_id = str(uuid.uuid4())
    user_id = str(uuid.uuid4())

    payload = {
        "event_type": "cla",
        "event_company_id": event_company_id,
        "event_project_id": event_project_id,
        "user_id": user_id,
        "event_data": json.dumps({"foo": "bar"}),
    }

    response = hug.test.post(routes, url, payload)
    assert response.status == HTTP_400

    user = get_user_instance()
    user.set_user_id(user_id)
    user.save()
    company = get_company_instance()
    company.set_company_id(event_company_id)
    company.set_company_name("foobar")
    company.save()
    project = get_project_instance()
    project.set_project_name("project")
    project.set_project_external_id(str(uuid.uuid4()))
    project.set_project_id(event_project_id)
    project.save()
    response = hug.test.post(routes, url, payload)
    assert response.status == HTTP_200
    assert response.data["data"]["event_type"] == "cla"
    assert response.data["data"]["event_company_id"] == event_company_id
    assert response.data["data"]["event_project_id"] == event_project_id
    assert response.data["data"]["user_id"] == user_id
    assert response.data["data"]["event_data"] == json.dumps({"foo": "bar"})


def test_search_event():
    """
    Test searching event
    """
    event_id = str(uuid.uuid4())
    event_type = "test_cla"
    event = get_event_instance()
    event.set_event_id(event_id)
    event.set_event_type(event_type)
    event_company_id = str(uuid.uuid4())
    event.set_event_company_id(event_company_id)
    event.save()
    url = "/v1/events"

    response = hug.test.get(routes, url, params={"event_type": "test_cla"})
    assert response.status == HTTP_200
    assert response.data["events"][0]["event_type"] == "test_cla"

    params = {"event_type": "random"}
    response = hug.test.get(routes, url, params=params)
    assert response.status == HTTP_404

    params = {"event_type": "test_cla", "event_company_id": event_company_id}
    response = hug.test.get(routes, url, params=params)
    assert response.status == HTTP_200
    assert response.data["events"][0]["event_type"] == "test_cla"
    assert response.data["events"][0]["event_company_id"] == event_company_id
