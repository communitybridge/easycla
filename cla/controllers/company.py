"""
Controller related to company operations.
"""

import uuid
import hug.types
from cla.utils import get_company_instance
from cla.models import DoesNotExist


def get_companies():
    """
    Returns a list of companies in the CLA system.

    :return: List of companies in dict format.
    :rtype: [dict]
    """
    return [company.to_dict() for company in get_company_instance().all()]


def get_company(company_id):
    """
    Returns the CLA company requested by ID.

    :param company_id: The company's ID.
    :type company_id: ID
    :return: dict representation of the company object.
    :rtype: dict
    """
    company = get_company_instance()
    try:
        company.load(company_id=str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    return company.to_dict()


def create_company(company_name=None,
                   company_whitelist=None,
                   company_whitelist_patterns=None):
    """
    Creates an company and returns the newly created company in dict format.

    :param company_name: The company name.
    :type company_name: string
    :param company_whitelist: The list of whitelisted domain names for this company.
    :type company_whitelist: [string]
    :param company_whitelist_patterns: List of whitelisted email patterns.
    :type company_whitelist_patterns: [string]
    :return: dict representation of the company object.
    :rtype: dict
    """
    company = get_company_instance()
    company.set_company_id(str(uuid.uuid4()))
    company.set_company_name(company_name)
    # TODO: Need to validate these values.
    company.set_company_whitelist(company_whitelist)
    company.set_company_whitelist_patterns(company_whitelist_patterns)
    company.save()
    return company.to_dict()


def update_company(company_id, # pylint: disable=too-many-arguments
                   company_name=None,
                   company_whitelist=None,
                   company_whitelist_patterns=None):
    """
    Updates an company and returns the newly updated company in dict format.
    A value of None means the field should not be updated.

    :param company_id: ID of the company to update.
    :type company_id: ID
    :param company_name: New company name.
    :type company_name: string | None
    :param company_whitelist: New whitelist for this company.
    :type company_whitelist: [string] | None
    :param company_whitelist_patterns: New company whitelisted email patterns.
    :type company_whitelist_patterns: [string] | None
    :return: dict representation of the company object.
    :rtype: dict
    """
    company = get_company_instance()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    if company_name is not None:
        company.set_company_name(company_name)
    # TODO: Need to validate these values.
    if company_whitelist is not None:
        val = hug.types.multiple(company_whitelist)
        company.set_company_whitelist(val)
    # TODO: Need to validate these values.
    if company_whitelist_patterns is not None:
        val = hug.types.multiple(company_whitelist_patterns)
        company.set_company_whitelist_patterns(val)
    company.save()
    return company.to_dict()


def delete_company(company_id):
    """
    Deletes an company based on ID.

    :param company_id: The ID of the company.
    :type company_id: ID
    """
    company = get_company_instance()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    company.delete()
    return {'success': True}
