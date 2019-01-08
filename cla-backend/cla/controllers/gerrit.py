"""
Controller related to repository operations.
"""

import uuid
import cla.hug_types
import os
from cla.utils import get_gerrit_instance
from cla.models import DoesNotExist
from cla.controllers.lf_group import LFGroup

def get_gerrits():
    """
    Returns a list of gerrit instances in the CLA system.

    :return: List of gerrit instances in dict format.
    :rtype: [dict]
    """
    return [gerrit_instance.to_dict() for gerrit_instance in get_gerrit_instance().all()]


def get_gerrit(gerrit_id):
    """
    Returns the CLA Gerrit Instance requested by ID.

    :param gerrit_id: The repository ID.
    :type gerrit_id: ID
    :return: dict representation of the Gerrit object.
    :rtype: dict
    """
    gerrit = get_gerrit_instance()
    try:
        gerrit.load(str(gerrit_id))
    except DoesNotExist as err:
        return {'errors': {'gerrit_id': str(err)}}
    return gerrit.to_dict()


def get_gerrits_by_project_id(project_id):
    gerrit = get_gerrit_instance()
    try:
        gerrits = gerrit.get_gerrits_by_project_id(project_id)
    except DoesNotExist as err:
        return {'errors': {'gerrit_id': str(err)}}
    return [gerrit.to_dict() for gerrit in gerrits]


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
    
    gerrit = get_gerrit_instance()
    gerrit.set_gerrit_id(str(uuid.uuid4()))
    gerrit.set_project_id(str(project_id))
    gerrit.set_gerrit_url(gerrit_url)
    gerrit.set_gerrit_name(gerrit_name)
    
    #check if LDAP group exists
    lf_group_client_url = os.environ.get('LF_GROUP_CLIENT_URL', '')
    lf_group_client_id = os.environ.get('LF_GROUP_CLIENT_ID', '')
    lf_group_client_secret = os.environ.get('LF_GROUP_CLIENT_SECRET', '')
    lf_group_refresh_token = os.environ.get('LF_GROUP_REFRESH_TOKEN', '')
    lf_group = LFGroup(lf_group_client_url, lf_group_client_id, lf_group_client_secret, lf_group_refresh_token)
     
    # returns 'error' if the LDAP group does not exist
    ldap_group_icla = lf_group.get_group(group_id_icla)
    if ldap_group_icla.get('error') is not None:
        return {'error_icla': 'The specified LDAP group for ICLA does not exist. '}
    gerrit.set_group_name_icla(ldap_group_icla.get('title'))
    gerrit.set_group_id_icla(str(group_id_icla))

    ldap_group_ccla = lf_group.get_group(group_id_ccla)
    if ldap_group_ccla.get('error') is not None:
        return {'error_ccla': 'The specified LDAP group for CCLA does not exist. '}
    gerrit.set_group_name_ccla(ldap_group_ccla.get('title'))
    gerrit.set_group_id_ccla(str(group_id_ccla))

    gerrit.save()
    return gerrit.to_dict()


def update_gerrit(gerrit_id, # pylint: disable=too-many-arguments
                    project_id=None,
                    gerrit_name=None,
                    gerrit_url=None,
                    group_id_icla=None,
                    group_id_ccla=None):
    """
    Updates a repository and returns the newly updated gerrt instance in dict format.
    Values of None means the field will not be updated.

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
    gerrit = get_gerrit_instance()
    try:
        gerrit.load(str(gerrit_id))
    except DoesNotExist as err:
        return {'errors': {'gerrit_id': str(err)}}
    # TODO: Ensure project_id exists.
    if project_id is not None:
        gerrit.set_project_id(str(project_id))
    if gerrit_name is not None:
        gerrit.set_gerrit_name(gerrit_name)
    if gerrit_url is not None:
        try:
            val = cla.hug_types.url(gerrit_url)
            gerrit.set_gerrit_url(val)
        except ValueError as err:
            return {'errors': {'gerrit_url': 'Invalid URL specified'}}
    if group_id_icla is not None:
        gerrit.set_group_id_icla(group_id_icla)
    if group_id_ccla is not None:
        gerrit.set_group_id_ccla(group_id_ccla)
    gerrit.save()
    return gerrit.to_dict()


def delete_gerrit(gerrit_id):
    """
    Deletes a gerrit instance

    :param gerrit_id: The ID of the gerrit instance.
    """
    gerrit = get_gerrit_instance()
    try:
        gerrit.load(str(gerrit_id))
    except DoesNotExist as err:
        return {'errors': {'gerrit_id': str(err)}}
    gerrit.delete()
    return {'success': True}


def get_agreement_html(gerrit_id, contract_type):
    return """
        <html>
            <a href="https://contributor.dev.lfcla.com/gerrit/{gerrit_id}/{contract_type}">Click on the link to Sign the CLA Agreement. </a>
        <html>""".format(gerrit_id=gerrit_id, contract_type=contract_type)
