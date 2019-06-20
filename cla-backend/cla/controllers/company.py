"""
Controller related to company operations.
"""

import uuid
import hug.types
from cla.models import DoesNotExist
import cla
import cla.controllers.user
from cla.auth import AuthUser, admin_list
from cla.models.dynamo_models import Company, User
from falcon import HTTP_409, HTTP_200, HTTPForbidden

def get_companies():
    """
    Returns a list of companies in the CLA system.

    :return: List of companies in dict format.
    :rtype: [dict]
    """
    return [company.to_dict() for company in Company().all()]

def get_companies_by_user(username):
    """
    Returns a list of companies for a user in the CLA system.

    :return: List of companies in dict format.
    :rtype: [dict]
    """
    all_companies = [company.to_dict() for company in Company().all() if username in company.get_company_acl()]

    return all_companies

def company_acl_verify(username, company):
    if username in company.get_company_acl():
        return True

    raise HTTPForbidden('Unauthorized',
        'Provided Token credentials does not have sufficient permissions to access resource')

def get_company(company_id):
    """
    Returns the CLA company requested by ID.

    :param company_id: The company's ID.
    :type company_id: ID
    :return: dict representation of the company object.
    :rtype: dict
    """
    company = Company()
    try:
        company.load(company_id=str(company_id))
        
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    return company.to_dict()

def create_company(auth_user,
                   company_name=None,
                   company_manager_id=None,
                   company_manager_user_name=None,
                   company_manager_user_email=None,
                   user_id=None):
    """
    Creates an company and returns the newly created company in dict format.

    :param company_name: The company name.
    :type company_name: string
    :param company_manager_id: The ID of the company manager user.
    :type company_manager_id: string
    :param company_manager_user_name: The user name of the company manager user.
    :type company_manager_user_name: string
    :param company_manager_user_email: The user email of the company manager user.
    :type company_manager_user_email: string
    :return: dict representation of the company object.
    :rtype: dict
    """

    manager = cla.controllers.user.get_or_create_user(auth_user)

    for company in get_companies():
        if company.get("company_name") == company_name:
            cla.log.error({"error": "Company already exists"})
            return {"status_code": HTTP_409,
                    "data": {"error":"Company already exists.",
                            "company_id": company.get("company_id")}
                    }

    company = Company()
    company.set_company_id(str(uuid.uuid4()))
    company.set_company_name(company_name)
    company.set_company_manager_id(manager.get_user_id())
    company.set_company_acl(manager.get_lf_username())

    company.save()

    return {"status_code": HTTP_200,
            "data": company.to_dict()
            }

def update_company(company_id, # pylint: disable=too-many-arguments
                   company_name=None,
                   company_manager_id=None,
                   username=None):
    """
    Updates an company and returns the newly updated company in dict format.
    A value of None means the field should not be updated.

    :param company_id: ID of the company to update.
    :type company_id: ID
    :param company_name: New company name.
    :type company_name: string | None
    :param company_manager_id: The ID of the company manager user.
    :type company_manager_id: string
    :return: dict representation of the company object.
    :rtype: dict
    """
    company = Company()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    company_acl_verify(username, company)

    if company_name is not None:
        company.set_company_name(company_name)
    if company_manager_id is not None:
        val = hug.types.uuid(company_manager_id)
        company.set_company_manager_id(str(val))
    company.save()
    return company.to_dict()

def update_company_whitelist_csv(content, company_id, username=None):
    """
    Adds the CSV of email addresse to this company's whitelist.

    :param content: The content posted to this endpoint (CSV data).
    :type content: string
    :param company_id: The ID of the company to add to the whitelist.
    :type company_id: UUID
    """
    company = Company()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    company_acl_verify(username, company)

    # Ready email addresses.
    emails = content.split('\n')
    emails = [email for email in emails if '@' in email]
    current_whitelist = company.get_company_whitelist()
    new_whitelist = list(set(current_whitelist + emails))
    company.set_company_whitelist(new_whitelist)
    company.save()
    return company.to_dict()

def delete_company(company_id, username=None):
    """
    Deletes an company based on ID.

    :param company_id: The ID of the company.
    :type company_id: ID
    """
    company = Company()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    company_acl_verify(username, company)

    company.delete()
    return {'success': True}


def get_manager_companies(manager_id):
    companies = Company().get_companies_by_manager(manager_id)
    return companies

def add_permission(auth_user: AuthUser, username: str, company_id: str, ignore_auth_user=False):
    if not ignore_auth_user and auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('company ({}) added for user ({}) by {}'.format(company_id, username, auth_user.username))

    company = Company()
    try:
        company.load(company_id)
    except Exception as err:
        print('Unable to update company permission: {}'.format(err))
        return {'error': str(err)}

    company.add_company_acl(username)
    company.save()

def remove_permission(auth_user: AuthUser, username: str, company_id: str):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('company ({}) removed for ({}) by {}'.format(company_id, username, auth_user.username))

    company = Company()
    try:
        company.load(company_id)
    except Exception as err:
        print('Unable to update company permission: {}'.format(err))
        return {'error': str(err)}

    company.remove_company_acl(username)
    company.save()
