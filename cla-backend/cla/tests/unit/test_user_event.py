# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from unittest.mock import patch, Mock

import pytest

from cla.models.dynamo_models import User, Project, Company
from cla.models.event_types import EventType
from cla.controllers.event import create_event
from cla.controllers import user as user_controller
from cla.auth import AuthUser


@pytest.fixture
def create_event_user():
    user_controller.create_event = Mock()


def test_request_company_whitelist(create_event_user, project, company, user):
    """ Test user requesting to be whitelisted event """
    event_type = EventType.RequestCompanyWL
    User.load = Mock()
    User.get_user_emails = Mock(return_value=[user.get_user_email()])
    User.get_user_email = Mock(return_value=user.get_user_email())
    Company.load = Mock()
    Project.load = Mock()
    User.get_user_name = Mock(return_value=user.get_user_name())
    Company.get_company_name = Mock(return_value=company.get_company_name())
    Project.get_project_name = Mock(return_value=project.get_project_name())
    user_controller.get_email_service = Mock()
    user_controller.send = Mock()
    user_controller.request_company_whitelist(
        user.get_user_id(),
        company.get_company_id(),
        user.get_user_email(),
        project.get_project_id(),
        message="Please add",
    )
    event_data = (
        f"CLA: {user.get_user_name()} requests to be whitelisted for the organization {company.get_company_name}"
        f"{user.get_user_name()} <{user.get_user_email()}>"
    )

    user_controller.create_event.assert_called_once_with(
        user_id=user.get_user_id(),
        event_project_id=project.get_project_id(),
        event_company_id=company.get_company_id(),
        event_type=event_type,
        event_data=event_data,
    )


def test_invite_company_admin(create_event_user, user):
    """ Test send email to CLA manager event """
    User.load = Mock()
    user_controller.send_email_to_admin = Mock()
    user_id = user.get_user_id()
    user_email = user.get_user_email()
    admin_name = "admin"
    admin_email = "foo@admin.com"
    project_name = "foo_project"
    event_data = f"{user_id} with {user_email} sends to {admin_name}/{admin_email} for project: {project_name}"
    user_controller.invite_company_admin(
        user_id, user_email, admin_name, admin_email, project_name
    )
    user_controller.create_event.assert_called_once_with(
        user_id=user_id,
        event_project_name=project_name,
        event_data=event_data,
        event_type=EventType.InviteAdmin,
    )


def test_request_company_ccla(create_event_user, user, project, company):
    """ Test request company ccla event """
    User.load = Mock()
    User.get_user_name = Mock(return_value=user.get_user_name())
    email = user.get_user_email()
    Company.load = Mock()
    Project.load = Mock()
    Project.get_project_name = Mock(return_value=project.get_project_name())
    manager = User(lf_username="harold", user_email="foo@gmail.com")
    Company.get_managers = Mock(return_value=[manager,])
    event_data = f"Sent email to sign ccla for {project.get_project_name() }"
    user_controller.request_company_ccla(
        user.get_user_id(), email, company.get_company_id(), project.get_project_id()
    )
    user_controller.create_event.assert_called_once_with(
        event_data=event_data,
        event_type=EventType.RequestCCLA,
        user_id=user.get_user_id(),
        event_company_id=company.get_company_id(),
    )




