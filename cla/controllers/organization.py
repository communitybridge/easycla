"""
Controller related to organization operations.
"""

import uuid
import hug.types
from cla.utils import get_organization_instance
from cla.models import DoesNotExist

def get_organizations():
    """
    Returns a list of organizations in the CLA system.

    :return: List of organizations in dict format.
    :rtype: [dict]
    """
    return [organization.to_dict() for organization in get_organization_instance().all()]

def get_organization(organization_id):
    """
    Returns the CLA organization requested by ID.

    :param organization_id: The organization's ID.
    :type organization_id: ID
    :return: dict representation of the organization object.
    :rtype: dict
    """
    organization = get_organization_instance()
    try:
        organization.load(organization_id=str(organization_id))
    except DoesNotExist as err:
        return {'errors': {'organization_id': str(err)}}
    return organization.to_dict()

def create_organization(organization_name=None,
                        organization_whitelist=None,
                        organization_exclude_patterns=None):
    """
    Creates an organization and returns the newly created organization in dict format.

    :param organization_name: The organization name.
    :type organization_name: string
    :param organization_whitelist: The list of whitelisted domain names for this organization.
    :type organization_whitelist: [string]
    :param organization_exclude_patterns: List of exclude patterns for email addresses.
    :type organization_exclude_patterns: [string]
    :return: dict representation of the organization object.
    :rtype: dict
    """
    organization = get_organization_instance()
    organization.set_organization_id(str(uuid.uuid4()))
    organization.set_organization_name(organization_name)
    organization.set_organization_whitelist(organization_whitelist)
    organization.set_organization_exclude_patterns(organization_exclude_patterns)
    organization.save()
    return organization.to_dict()

def update_organization(organization_id, # pylint: disable=too-many-arguments
                        organization_name=None,
                        organization_whitelist=None,
                        organization_exclude_patterns=None):
    """
    Updates an organization and returns the newly updated organization in dict format.
    A value of None means the field should not be updated.

    :param organization_id: ID of the organization to update.
    :type organization_id: ID
    :param organization_name: New organization name.
    :type organization_name: string | None
    :param organization_whitelist: New whitelist for this organization.
    :type organization_whitelist: [string] | None
    :param organization_exclude_patterns: New exclude patterns list for this organization.
    :type organization_exclude_patterns: [string] | None
    :return: dict representation of the organization object.
    :rtype: dict
    """
    organization = get_organization_instance()
    try:
        organization.load(str(organization_id))
    except DoesNotExist as err:
        return {'errors': {'organization_id': str(err)}}
    if organization_name is not None:
        organization.set_organization_name(organization_name)
    if organization_whitelist is not None:
        val = hug.types.multiple(organization_whitelist)
        organization.set_organization_whitelist(val)
    if organization_exclude_patterns is not None:
        val = hug.types.multiple(organization_exclude_patterns)
        organization.set_organization_exclude_patterns(val)
    organization.save()
    return organization.to_dict()

def delete_organization(organization_id):
    """
    Deletes an organization based on ID.

    :param organization_id: The ID of the organization.
    :type organization_id: ID
    """
    organization = get_organization_instance()
    try:
        organization.load(str(organization_id))
    except DoesNotExist as err:
        return {'errors': {'organization_id': str(err)}}
    organization.delete()
    return {'success': True}
