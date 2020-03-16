# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from unittest.mock import Mock, patch

import pytest

import cla
from cla.auth import AuthUser
from cla.controllers import company as company_controller
from cla.models.dynamo_models import Company
from cla.models.event_types import EventType


@pytest.fixture()
def create_event_company():
    company_controller.create_event = Mock()


@pytest.fixture()
def auth_user():
    with patch.object(AuthUser, "__init__", lambda self: None):
        auth_user = AuthUser()
        yield auth_user

@patch('cla.controllers.company.Event.create_event')
def test_create_company_event(mock_event, auth_user, create_event_company, user, company):
    """ Test create company event """
    cla.controllers.user.get_or_create_user = Mock(return_value=user)
    company_controller.get_companies = Mock(return_value=[])
    Company.save = Mock()
    company_name="new_company"
    Company.get_company_name = Mock(return_value=company_name)
    Company.get_company_id = Mock(return_value='manager_id')
    company_controller.create_company(
        auth_user,
        company_name=company_name,
        company_manager_id="manager_id",
        company_manager_user_name="user name",
        company_manager_user_email="email",
        user_id=user.get_user_id(),
    )
    event_data = "Company-{} created".format(company_name)
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_type=EventType.CreateCompany,
        event_company_id=company.get_company_id(),
        event_user_id=user.get_user_id(),
        contains_pii=False,
    )

@patch('cla.controllers.company.Event.create_event')
def test_update_company_event(mock_event, create_event_company, company):
    """ Test update company """
    event_type = EventType.UpdateCompany
    Company.load = Mock()
    company_controller.company_acl_verify = Mock()
    Company.save = Mock()
    company_name = 'new name'
    company_controller.update_company(
        company.get_company_id(),
        company_name=company_name,
    )
    event_data = "company_name updated to {} \n".format(company_name)
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_type=event_type,
        event_company_id=company.get_company_id(),
        contains_pii=False,
    )

@patch('cla.controllers.company.Event.create_event')
def test_delete_company(mock_event, create_event_company, company):
    """ Test delete company event """
    event_type = EventType.DeleteCompany
    Company.load = Mock()
    company_controller.company_acl_verify = Mock()
    Company.delete = Mock()
    event_data = f'Company- {company.get_company_name()} deleted'
    company_controller.delete_company(
        company.get_company_id()
    )
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_type=event_type,
        event_company_id=company.get_company_id(),
        contains_pii=False,
    )

@patch('cla.controllers.company.Event.create_event')
def test_add_permission(mock_event, create_event_company, auth_user, company):
    """ Test add permission event """
    event_type = EventType.AddCompanyPermission
    Company.load = Mock()
    Company.add_company_acl = Mock()
    auth_user.username='ddeal'
    username = 'foo_username'
    company_controller.add_permission(
        auth_user,
        username,
        company.get_company_id(),
        ignore_auth_user=True
    )
    event_data = f'Permissions added to user {username} for Company {company.get_company_name()}'
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_type=event_type,
        event_company_id=company.get_company_id(),
        contains_pii = True,
    )


@patch('cla.controllers.company.Event.create_event')
def test_remove_permission(mock_event, create_event_company, auth_user, company):
    """Test remove permissions """
    event_type=EventType.RemoveCompanyPermission
    Company.load = Mock()
    Company.remove_company_acl = Mock()
    Company.save = Mock()
    auth_user.username='ddeal'
    company_id = company.get_company_id()
    username='remover'
    event_data = 'company ({}) removed for ({}) by {}'.format(company_id, username, auth_user.username)
    company_controller.remove_permission(
        auth_user,
        username,
        company_id
    )
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_company_id=company_id,
        event_type=event_type,
        contains_pii = True,
    )

