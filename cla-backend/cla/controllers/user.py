"""
Controller related to user operations.
"""

import uuid
import hug
from cla.utils import get_user_instance, get_company_instance, get_email_service
from cla.models import DoesNotExist
from cla.models.dynamo_models import User, Company, Project
import cla

def get_users():
    """
    Returns a list of users in the CLA system.

    :return: List of users in dict format.
    :rtype: [dict]
    """
    return [user.to_dict() for user in get_user_instance().all()]

def get_user(user_id=None, user_email=None, user_github_id=None):
    """
    Returns the CLA user requested by ID or email.

    :param user_id: The user's ID.
    :type user_id: string
    :param user_email: The user's email address.
    :type user_email: string
    :param user_github_id: The user's github ID.
    :type user_github_id: integer
    :return: dict representation of the user object.
    :rtype: dict
    """
    if user_id is not None:
        user = get_user_instance()
        try:
            user.load(user_id)
        except DoesNotExist as err:
            return {'errors': {'user_id': str(err)}}
    elif user_email is not None:
        user = get_user_instance().get_user_by_email(str(user_email).lower())
        if user is None:
            return {'errors': {'user_email': 'User not found'}}
    elif user_github_id is not None:
        user = get_user_instance().get_user_by_github_id(user_github_id)
        if user is None:
            return {'errors': {'user_github_id': 'User not found'}}
    return user.to_dict()

def get_user_signatures(user_id):
    """
    Given a user ID, returns the user's signatures.

    :param user_id: The user's ID.
    :type user_id: string
    :return: list of signature data for this user.
    :rtype: [dict]
    """
    user = get_user_instance()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    signatures = user.get_user_signatures()
    return [agr.to_dict() for agr in signatures]

def get_users_company(user_company_id):
    """
    Fetches all users that are associated with the company specified.

    :param user_company_id: The ID of the company in question.
    :type user_company_id: string
    :return: A list of user data in dict format.
    :rtype: [dict]
    """
    users = get_user_instance().get_users_by_company(user_company_id)
    return [user.to_dict() for user in users]

def request_company_whitelist(user_id, company_id, user_email, project_id, message=None):
    """
    Sends email to the specified company manager notifying them that a user has requested to be
    added to their whitelist.

    :param user_id: The ID of the user requesting to be added to the company's whitelist.
    :type user_id: string
    :param company_id: The ID of the company that the request is going to.
    :type company_id: string
    :param user_email: The email address that this user wants to be whitelisted. Must exist in the
        user's list of emails.
    :type user_email: string
    :param messsage: A custom message to add to the email sent out to the manager.
    :type message: string
    """
    user = User()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    emails = user.get_user_emails()
    if user_email not in emails:
        return {'errors': {'user_email': 'Must provide one of the user\'s existing emails'}}
    company = Company()
    try:
        company.load(company_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    project = Project()
    try: 
        project.load(project_id)
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}

    user_name = user.get_user_name()    
    company_name = company.get_company_name()
    project_name = project.get_project_name()

    subject = '''CLA: %s is requesting to be whitelisted for %s project ''' %(user_name, project_name)

    body = '''%s is requesting to be whitelisted as a contributor for your organization (%s):

    %s <%s>

The message that was attached to the request:

    %s

You can whitelist %s in the CLA Corporate console. If the email above is the personal email of one of your employees, please request that they add their organization email to their GitHub profile and try signing the CLA again. If you are unsure about this request, it may be prudent to get in touch with %s to clarify.
Please follow up with the user as necessary.

Click on the following link to navigate to the CLA Corporate Console.

 %s  

- Linux Foundation CLA System
''' %(user_name, company_name, user_name, user_email, message,
    user_name, user_name, 'https://{}'.format(cla.conf['CORPORATE_BASE_URL']))

    manager_id = company.get_company_manager_id()
    manager = get_user_instance()
    try:
        manager.load(manager_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': 'No manager exists for this company - can not send email'}}
    recipient = manager.get_user_email()
    email_service = get_email_service()
    email_service.send(subject, body, recipient)


def invite_company_admin(user_id, user_email, admin_name, admin_email):
    """
    Sends email to the specified company administrator to sign up through the CCLA console and add the requested user to the whitelist. 

    :param user_id: The ID of the user requesting to be added to the company's whitelist.
    :type user_id: string
    :param user_email: The email address that this user wants to be whitelisted. Must exist in the
        user's list of emails.
    :type user_email: string
    :param messsage: A message to be sent out to the administrator. 
    """
    user = get_user_instance()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}

    subject = 'CLA: Invitation to Sign Up for Corporate CLA'
    body = '''Hello %s, 
    
    The following user is requesting to be whitelisted as a contributor for your organization:

    %s <%s>

Please click the following link to sign up for Corporate CLA and add this user to your organization. 

%s

- Linux Foundation CLA System
''' %(admin_name, user.get_user_name(), user.get_user_email(), 'https://{}'.format(cla.conf['CORPORATE_BASE_URL']))
    recipient = admin_email
    email_service = get_email_service()
    email_service.send(subject, body, recipient)

def get_active_signature(user_id):
    """
    Returns information on the user's active signature - if there is one.

    :param user_id: The ID of the user.
    :type user_id: string
    :return: A dictionary of all the active signature's metadata, along with the return_url.
    :rtype: dict | None
    """
    metadata = cla.utils.get_active_signature_metadata(user_id)
    if metadata is None:
        return None
    return_url = cla.utils.get_active_signature_return_url(user_id, metadata)
    metadata['return_url'] = return_url
    return metadata

def get_user_project_last_signature(user_id, project_id):
    """
    Returns the user's last signature object for a project.

    :param user_id: The ID of the user.
    :type user_id: string
    :param project_id: The project in question.
    :type project_id: string
    :return: The signature object that was last signed by the user for this project.
    :rtype: cla.models.model_interfaces.Signature
    """
    user = get_user_instance()
    try:
        user.load(str(user_id))
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    last_signature = user.get_latest_signature(str(project_id))
    if last_signature is not None:
        last_signature = last_signature.to_dict()
        latest_doc = cla.utils.get_project_latest_individual_document(str(project_id))
        last_signature['latest_document_major_version'] = str(latest_doc.get_document_major_version())
        last_signature['latest_document_minor_version'] = str(latest_doc.get_document_minor_version())
        last_signature['requires_resigning'] = last_signature['latest_document_major_version'] != last_signature['signature_document_major_version']
    return last_signature

def get_user_project_company_last_signature(user_id, project_id, company_id):
    """
    Returns the user's last signature object for a project.

    :param user_id: The ID of the user.
    :type user_id: string
    :param project_id: The project in question.
    :type project_id: string
    :param company_id: The ID of the company that this employee belongs to.
    :type company_id: string
    :return: The signature object that was last signed by the user for this project.
    :rtype: cla.models.model_interfaces.Signature
    """
    user = get_user_instance()
    try:
        user.load(str(user_id))
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    last_signature = user.get_latest_signature(str(project_id), company_id=str(company_id))
    if last_signature is not None:
        last_signature = last_signature.to_dict()
        latest_doc = cla.utils.get_project_latest_corporate_document(str(project_id))
        last_signature['latest_document_major_version'] = str(latest_doc.get_document_major_version())
        last_signature['latest_document_minor_version'] = str(latest_doc.get_document_minor_version())
        last_signature['requires_resigning'] = last_signature['latest_document_major_version'] != last_signature['signature_document_major_version']
    return last_signature

# For GitHub user creating, see models.github_models.get_or_create_user(self, request)
def get_or_create_user(auth_user):
    user = User()

    existing_user = user.get_user_by_username(str(auth_user.username))
    
    if existing_user is None:
        user.set_user_id(str(uuid.uuid4()))
        user.set_user_name(auth_user.name)
        user.set_lf_email(auth_user.email.lower())
        user.set_lf_username(auth_user.username)
        user.set_lf_sub(auth_user.sub)

        user.save()

        return user

    return existing_user
