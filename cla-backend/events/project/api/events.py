# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import logging
import os

from flask import Blueprint, request
from flask_restful import Api, Resource

from cla import config
from cla.models import DoesNotExist
from cla.models.dynamo_models import Event, EventModel
from cla.utils import get_event_instance

events_blueprint = Blueprint("events", __name__)
api = Api(events_blueprint)

logging.Formatter(fmt=config.LOG_FORMAT)


class EventsList(Resource):
    """
    Events endpoint that creates event
    """

    def post(self):
        post_data = request.get_json()
        response = {"status": "fail", "message": "Invalid payload"}
        if not post_data:
            return response, 400
        event_id = post_data.get("event_id")
        event_type = post_data.get("event_type")
        event_data = post_data.get("event_data")
        event_project_id = post_data.get("event_project_id",None)
        event_time = post_data.get("event_time",None)
        event_company_id = post_data.get("event_company_id",None)

        try:
            event = get_event_instance()
            event.set_event_id(event_id)
            event.set_event_type(event_type)
            event.set_event_data(event_data)
            event.set_event_company_id(event_company_id)
            event.set_event_project_id(event_project_id)
            event.set_event_time(event_time)
            event.save()
            response = {"status": "success", "message": "event was added"}

        except Exception as err:
            logging.error(err)
            return response, 400
        return response, 200

    def get(self):
        """
        Get all events
        """
        response = {
            'status': 'fail'
        }
        try:
            event = get_event_instance()
            response = {
                'status' : 'success',
                'data' : {
                    'events':[event.to_dict() for event in event.all()]
                }
            }
        except Exception as err:
            logging.error(err)
            return response, 404
        return response, 200



class Events(Resource):
    """
    Endpoint that gets event by event_id
    """

    def get(self, event_id):
        response = {"status": "fail", "message": "Event does not exist"}
        try:
            event_instance = get_event_instance()
            try:
                event_instance.load(event_id)
            except DoesNotExist:
                return response, 404
            response = {"status": "success", "data": f"{event_instance}"}
            return response, 200

        except DoesNotExist:
            return response, 404


api.add_resource(EventsList, "/events")
api.add_resource(Events, "/events/<event_id>")
 