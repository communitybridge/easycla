# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Controller related to repository operations.
"""

import uuid
import cla.hug_types
import os
import cla 
from cla.models.dynamo_models import Gerrit
from cla.models import DoesNotExist
from cla.controllers.lf_group import LFGroup

lf_group_client_url = os.environ.get('LF_GROUP_CLIENT_URL', '')
lf_group_client_id = os.environ.get('LF_GROUP_CLIENT_ID', '')
lf_group_client_secret = os.environ.get('LF_GROUP_CLIENT_SECRET', '')
lf_group_refresh_token = os.environ.get('LF_GROUP_REFRESH_TOKEN', '')
lf_group = LFGroup(lf_group_client_url, lf_group_client_id, lf_group_client_secret, lf_group_refresh_token)

def get_gerrit_by_project_id(project_id):
    gerrit = Gerrit()
    try:
        gerrits = gerrit.get_gerrit_by_project_id(project_id)
    except DoesNotExist:
        return []
    except Exception as e:
        return {'errors': {'a gerrit instance does not exist with the given project ID. ': str(e)}}

    if gerrits is None:
        return []

    return [gerrit.to_dict() for gerrit in gerrits]

def get_gerrit(gerrit_id):
    gerrit = Gerrit()
    try:
        gerrit.load(str(gerrit_id))
    except DoesNotExist as err:
        return {'errors': {'a gerrit instance does not exist with the given Gerrit ID. ': str(err)}}

    return gerrit.to_dict()


def create_gerrit(project_id, 
                    gerrit_name, 
                    gerrit_url, 
                    group_id_icla,
                    group_id_ccla):
    """
    Creates a gerrit instance and returns the newly created gerrit object dict format.

    :param gerrit_project_id: The project ID of the gerrit instance
    :type gerrit_project_id: string
    :param gerrit_name: The new gerrit instance name
    :type gerrit_name: string
    :param gerrit_url: The new Gerrit URL.
    :type gerrit_url: string
    :param group_id_icla: The id of the LDAP group for ICLA. 
    :type group_id_icla: string
    :param group_id_ccla: The id of the LDAP group for CCLA. 
    :type group_id_ccla: string
    """

    gerrit = Gerrit()

    # Check if at least ICLA or CCLA is specified 
    if group_id_icla is None and group_id_ccla is None:
        return {'error': 'Should specify at least a LDAP group for ICLA or CCLA.'}

    # Check if ICLA exists
    if group_id_icla is not None:
        ldap_group_icla = lf_group.get_group(group_id_icla)
        if ldap_group_icla.get('error') is not None:
            return {'error_icla': 'The specified LDAP group for ICLA does not exist. '}

        gerrit.set_group_name_icla(ldap_group_icla.get('title'))
        gerrit.set_group_id_icla(str(group_id_icla))

    # Check if CCLA exists
    if group_id_ccla is not None:
        ldap_group_ccla = lf_group.get_group(group_id_ccla)
        if ldap_group_ccla.get('error') is not None:
            return {'error_ccla': 'The specified LDAP group for CCLA does not exist. '}

        gerrit.set_group_name_ccla(ldap_group_ccla.get('title'))
        gerrit.set_group_id_ccla(str(group_id_ccla))

    # Save Gerrit Instance
    gerrit.set_gerrit_id(str(uuid.uuid4()))
    gerrit.set_project_id(str(project_id))
    gerrit.set_gerrit_url(gerrit_url)
    gerrit.set_gerrit_name(gerrit_name)
    gerrit.save()

    return gerrit.to_dict()

def delete_gerrit(gerrit_id):
    """
    Deletes a gerrit instance

    :param gerrit_id: The ID of the gerrit instance.
    """
    gerrit = Gerrit()
    try:
        gerrit.load(str(gerrit_id))
    except DoesNotExist as err:
        return {'errors': {'gerrit_id': str(err)}}
    gerrit.delete()
    return {'success': True}


def get_agreement_html(project_id, contract_type):
    contributor_base_url = cla.conf['CONTRIBUTOR_BASE_URL']
    return """
        <html>
            <a href="https://{contributor_base_url}/#/cla/gerrit/project/{project_id}/{contract_type}">Thank you. Unfortunately, your account is not authorized under a signed CLA. Please click here to proceed. </a>
        <html>""".format(
            contributor_base_url = contributor_base_url,
            project_id = project_id,
            contract_type = contract_type
        )
