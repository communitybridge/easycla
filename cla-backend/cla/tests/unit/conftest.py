# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import pytest

from unittest.mock import patch
from cla.tests.unit.data import COMPANY_TABLE_DATA, USER_TABLE_DATA, SIGNATURE_TABLE_DATA
from cla.utils import get_user_instance, get_signature_instance, get_company_instance

PATCH_METHOD = "pynamodb.connection.Connection._make_api_call"


@pytest.fixture()
def signature_instance():
    """
    Mock signature instance
    """
    with patch(PATCH_METHOD) as req:
        req.return_value = SIGNATURE_TABLE_DATA
        instance = get_signature_instance()
        instance.set_signature_id("sig_id")
        instance.set_signature_project_id("proj_id")
        instance.set_signature_reference_id("ref_id")
        instance.set_signature_type("type")
        instance.set_signature_project_external_id("proj_id")
        instance.set_signature_company_signatory_id("comp_sig_id")
        instance.set_signature_company_signatory_name("name")
        instance.set_signature_company_signatory_email("email")
        instance.set_signature_company_initial_manager_id("manager_id")
        instance.set_signature_company_initial_manager_name("manager_name")
        instance.set_signature_company_initial_manager_email("manager_email")
        instance.set_signature_company_secondary_manager_list({"foo": "bar"})
        instance.set_signature_document_major_version(1)
        instance.set_signature_document_minor_version(2)
        instance.save()
        yield instance


@pytest.fixture()
def user_instance():
    """
    Mock user instance
    """
    with patch(PATCH_METHOD) as req:
        req.return_value = USER_TABLE_DATA
        instance = get_user_instance()
        instance.set_user_id("foo")
        instance.set_user_external_id("bar")
        instance.save()
        yield instance


@pytest.fixture()
def company_instance():
    """
    Mock Company instance
    """
    with patch(PATCH_METHOD) as req:
        req.return_value = COMPANY_TABLE_DATA
        instance = get_company_instance()
        instance.set_company_id("uuid")
        instance.set_company_name("co")
        instance.set_company_external_id("external id")
        instance.save()
        yield instance
