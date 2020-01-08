# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import uuid
import json
from datetime import datetime

import hug.types
from falcon import HTTP_200, HTTP_400, HTTP_404, HTTPError

import cla
from cla.auth import AuthUser, admin_list
from cla.models import DoesNotExist
from cla.utils import get_event_instance


def events(request, response=None):
    """
    Returns a list of events in the CLA system.
    if parameters are passed returns filtered lists

    :return: List of events in dict format
    """

    event = get_event_instance()
    events = [event.to_dict() for event in event.all()]
    if request.params:
        results = event.search_events(**request.params)
        if results:
            events = [ev.to_dict() for ev in results]
        else:
            # return empty list if search fails
            response.status = HTTP_404
            return {"events": []}

    return {"events": events}


def get_event(event_id=None, response=None):
    """
    Returns an event given an event_id

    :param event_id: The event's ID
    :type event_id: string
    :rtype: dict
    """
    if event_id is not None:
        event = get_event_instance()
        try:
            event.load(event_id)
        except DoesNotExist as err:
            response.status = HTTP_404
            return {"errors": {"event_id": str(err)}}
        return event.to_dict()
    else:
        response.status = HTTP_404
        return {"errors": "Id is not passed"}


def create_event(
    response=None, event_type=None, event_project_id=None, event_company_id=None, event_data=None, user_id=None,
):
    """
    Creates an event returns the newly created event in dict format.

    :param event_type: The type of event
    :type event_type: string
    :param user_id: The user that is assocaited with the event
    :type user_id: string
    :param event_project_id: The project associated with event
    :type event_project_id: string

    """
    try:
        event = get_event_instance()
        if event_project_id:
            try:
                project = cla.utils.get_project_instance()
                project.load(str(event_project_id))
                event.set_event_project_id(event_project_id)
            except DoesNotExist as err:
                response.status = HTTP_400
                return {"errors": {"event_project_id": str(err)}}
        if event_company_id:
            try:
                company = cla.utils.get_company_instance()
                company.load(str(event_company_id))
                event.set_event_company_id(event_company_id)
            except DoesNotExist as err:
                response.status = HTTP_400
                return {"errors": {"event_company_id": str(err)}}
        if user_id:
            try:
                user = cla.utils.get_user_instance()
                user.load(str(user_id))
                event.set_user_id(user_id)
            except DoesNotExist as err:
                response.status = HTTP_400
                return {"errors": {"user_id": str(err)}}

        event.set_event_id(str(uuid.uuid4()))
        if event_type:
            event.set_event_type(event_type)
        event.set_event_data(event_data)
        event.set_event_time(str(datetime.now()))
        event.save()
        return {"status_code": HTTP_200, "data": event.to_dict()}

    except Exception as err:
        return {"errors": {"event_id": str(err)}}
