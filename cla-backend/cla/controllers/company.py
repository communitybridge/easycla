"""
Controller related to company operations.
"""

import uuid
import hug.types
from cla.utils import get_company_instance
from cla.utils import get_user_instance
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
                   company_whitelist_patterns=None,
                   company_manager_id=None,
                   company_manager_user_name=None,
                   company_manager_user_email=None):
    """
    Creates an company and returns the newly created company in dict format.

    :param company_name: The company name.
    :type company_name: string
    :param company_whitelist: The list of whitelisted domain names for this company.
    :type company_whitelist: [string]
    :param company_whitelist_patterns: List of whitelisted email patterns.
    :type company_whitelist_patterns: [string]
    :param company_manager_id: The ID of the company manager user.
    :type company_manager_id: string
    :param company_manager_user_name: The user name of the company manager user.
    :type company_manager_user_name: string
    :param company_manager_user_email: The user email of the company manager user.
    :type company_manager_user_email: string
    :return: dict representation of the company object.
    :rtype: dict
    """
    if company_manager_user_email not None:
        company_manager_id = get_or_create_user_as_manager(company_manager_user_email, company_manager_user_name)
    company = get_company_instance()
    company.set_company_id(str(uuid.uuid4()))
    company.set_company_name(company_name)
    # TODO: Need to validate these values.
    company.set_company_whitelist(company_whitelist)
    company.set_company_whitelist_patterns(company_whitelist_patterns)
    company.set_company_manager_id(str(company_manager_id))
    company.save()
    return company.to_dict()

def update_company(company_id, # pylint: disable=too-many-arguments
                   company_name=None,
                   company_whitelist=None,
                   company_whitelist_patterns=None,
                   company_manager_id=None):
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
    :param company_manager_id: The ID of the company manager user.
    :type company_manager_id: string
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
    if company_manager_id is not None:
        val = hug.types.uuid(company_manager_id)
        company.set_company_manager_id(str(val))
    company.save()
    return company.to_dict()

def update_company_whitelist_csv(content, company_id):
    """
    Adds the CSV of email addresse to this company's whitelist.

    :param content: The content posted to this endpoint (CSV data).
    :type content: string
    :param company_id: The ID of the company to add to the whitelist.
    :type company_id: UUID
    """
    company = get_company_instance()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    # Ready email addresses.
    emails = content.split('\n')
    emails = [email for email in emails if '@' in email]
    current_whitelist = company.get_company_whitelist()
    new_whitelist = list(set(current_whitelist + emails))
    company.set_company_whitelist(new_whitelist)
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


def get_manager_companies(manager_id):
    companies = get_company_instance().get_companies_by_manager(manager_id)
    return companies

def get_or_create_user_as_manager(email, user_name):
    """
    Helper method to either get or create a user as company manager based on the email address.
    And return user id as manager id for a company.

    :param email: The email address of a user who is creating a company.
    :type email: string
    :param user_name: The user name of a user who is creating a company.
    :type user_name: string
    :return: user_id
    :type: string
    """
    user = get_user_instance()
    found = user.get_user_by_email(email)
    if found is not None:
        # Found user by email, set the GitHub ID
        cla.log.info('Loaded user by email: %s - (%s)', found.get_user_name(), found.get_user_emails())
        return str(found.get_user_id())
    # User not found, create.
    cla.log.debug('Could not find user by email: %s', email)
    user_uuid = str(uuid.uuid4())
    user.set_user_id(user_uuid)
    user.set_user_emails(emails)
    user.set_user_name(userName)
    user.save()
    return user_uuid
