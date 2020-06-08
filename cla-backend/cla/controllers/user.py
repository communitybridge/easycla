# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Controller related to user operations.
"""

import uuid

import cla
from cla.models import DoesNotExist
from cla.models.dynamo_models import User, Company, Project, Event
from cla.models.event_types import EventType
from cla.utils import get_user_instance, get_email_service


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
        users = get_user_instance().get_user_by_email(str(user_email).lower())
        if users is None:
            return {'errors': {'user_email': 'User not found'}}
        # Use the first user for now - need to revisit - what if multiple are returned?
        user = users[0]
    elif user_github_id is not None:
        users = get_user_instance().get_user_by_github_id(user_github_id)
        if users is None:
            return {'errors': {'user_github_id': 'User not found'}}
        # Use the first user for now - need to revisit - what if multiple are returned?
        user = users[0]
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


def request_company_whitelist(user_id: str, company_id: str, user_name: str, user_email: str, project_id: str,
                              message: str = None, recipient_name: str = None, recipient_email: str = None):
    """
    Sends email to the specified company manager notifying them that a user has requested to be
    added to their approval list.

    :param user_id: The ID of the user requesting to be added to the company's approval list.
    :type user_id: string
    :param company_id: The ID of the company that the request is going to.
    :type company_id: string
    :param user_name: The name hat this user wants to be approved
    :type user_name: string
    :param user_email: The email address that this user wants to be approved. Must exist in the
        user's list of emails.
    :type user_email: string
    :param project_id: The ID of the project that the request is going to.
    :type project_id: string
    :param message: A custom message to add to the email sent out to the manager.
    :type message: string
    :param recipient_name: An optional recipient name for requesting the company approval list
    :type recipient_name: string
    :param recipient_email: An optional recipient email for requesting the company approval list
    :type recipient_email: string
    """
    if project_id is None:
        return {'errors': {'project_id': 'Project ID is missing from the request'}}
    if company_id is None:
        return {'errors': {'company_id': 'Company ID is missing from the request'}}
    if user_id is None:
        return {'errors': {'user_id': 'User ID is missing from the request'}}
    if user_name is None:
        return {'errors': {'user_name': 'User Name is missing from the request'}}
    if user_email is None:
        return {'errors': {'user_email': 'User Email is missing from the request'}}
    if recipient_name is None:
        return {'errors': {'recipient_name': 'Recipient Name is missing from the request'}}
    if recipient_email is None:
        return {'errors': {'recipient_email': 'Recipient Email is missing from the request'}}
    if message is None:
        return {'errors': {'message': 'Message is missing from the request'}}

    user = User()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}

    if user_email not in user.get_user_emails():
        return {'errors': {'user_email': 'User\'s email must match one of the user\'s existing emails in their profile'}}

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

    company_name = company.get_company_name()
    project_name = project.get_project_name()

    msg = ''
    if message is not None:
        msg += f'<p>{user_name} included the following message in the request:</p>'
        msg += f'<p>{message}</p>'

    subject = f'EasyCLA: Request to Authorize {user_name} for {project_name}'
    body = f'''
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello {recipient_name},</p>
<p>This is a notification email from EasyCLA regarding the project {project_name}.</p>
<p>{user_name} ({user_email}) has requested to be added to the Allow List as an authorized contributor from
{company_name} to the project {project_name}. You are receiving this message as a CLA Manager from {company} for
{project_name}.</p>
{msg}
<p>If you want to add them to the Allow List, please
<a href="https://{cla.conf['CORPORATE_BASE_URL']}#/company/{company_id}" target="_blank">log into the EasyCLA Corporate
Console</a>, where you can approve this user's request by selecting the 'Manage Approved List' and adding the
contributor's email, the contributor's entire email domain, their GitHub ID or the entire GitHub Organization for the
repository. This will permit them to begin contributing to {project_name} on behalf of {company}.</p>
<p>If you are not certain whether to add them to the Allow List, please reach out to them directly to discuss.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or 
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>
'''

    cla.log.debug(f'request_company_approval_list - sending email '
                  f'to recipient {recipient_name}/{recipient_email} '
                  f'for user {user_name}/{user_email} '
                  f'for project {project_name} '
                  f'assigned to company {company_name}')
    email_service = get_email_service()
    email_service.send(subject, body, recipient_email)

    # Create event
    event_data = (f'CLA: contributor {user_name} requests to be Approved for the '
                  f'project: {project_name} '
                  f'organization: {company_name} '
                  f'as {user_name} <{user_email}>')
    Event.create_event(
        event_user_id=user_id,
        event_project_id=project_id,
        event_company_id=company_id,
        event_type=EventType.RequestCompanyWL,
        event_data=event_data,
        contains_pii=True,
    )


def invite_cla_manager(contributor_id, contributor_name, contributor_email, cla_manager_name, cla_manager_email,
                       project_name, company_name):
    """
    Sends email to the specified CLA Manager to sign up through the Corporate
    console and adds the requested user to the Approved List request queue.

    :param contributor_id: The id of the user inviting the CLA Manager
    :param contributor_name: The name of the user inviting the CLA Manager
    :param contributor_email: The email address that this user wants to be added to the Approved List. Must exist in the user's list of emails.
    :param cla_manager_name: The name of the CLA manager
    :param cla_manager_email: The email address of the CLA manager
    :param project_name: The name of the project
    :param company_name: The name of the organization/company
    """
    user = User()
    try:
        user.load(contributor_id)
    except DoesNotExist as err:
        msg = f'unable to load user by id: {contributor_id} for inviting company admin - error: {err}'
        cla.log.warning(msg)
        return {'errors': {'user_id': contributor_id, 'message': msg, 'error': str(err)}}

    # We'll use the user's provided contributor name - if not provided use what we have in the DB
    if contributor_name is None:
        contributor_name = user.get_user_name()

    log_msg = (f'sent email to CLA Manager: {cla_manager_name} with email {cla_manager_email} '
               f'for project {project_name} and company {company_name} '
               f'to user {contributor_name} with email {contributor_email}')
    # Send email to the admin. set account_exists=False since the admin needs to sign up through the Corporate Console.
    cla.log.info(log_msg)
    send_email_to_cla_manager(contributor_name, contributor_email,
                              cla_manager_name, cla_manager_email,
                              project_name, company_name, False)

    Event.create_event(
        event_user_id=contributor_id,
        event_project_name=project_name,
        event_data=log_msg,
        event_type=EventType.InviteAdmin,
        contains_pii=True,
    )


def request_company_ccla(user_id, user_email, company_id, project_id):
    """
    Sends email to all company administrators in the company ACL to sign a CCLA for the given project.
    """
    user = User()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    user_name = user.get_user_name()

    company = Company()
    try:
        company.load(company_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    company_name = company.get_company_name()

    project = Project()
    try:
        project.load(project_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    project_name = project.get_project_name()

    # Send an email to sign the ccla for the project for every member in the company ACL
    # account_exists=True since company already exists.
    for admin in company.get_managers():
        send_email_to_cla_manager(user_name, user_email, admin.get_user_name(),
                                  admin.get_lf_email(), project_name, company_name, True)

    # Audit event
    event_data = f'Sent email to sign ccla for {project.get_project_name()}'
    Event.create_event(
        event_data=event_data,
        event_type=EventType.RequestCCLA,
        event_user_id=user_id,
        event_company_id=company_id,
        contains_pii=False,
    )


def send_email_to_cla_manager(contributor_name, contributor_email, cla_manager_name, cla_manager_email, project_name,
                              company_name, account_exists):
    """
    Helper function to send an email to a prospective CLA Manager.

    :param contributor_name: The name of the user sending the email.
    :param contributor_email: The email address that this user wants to be added to the approval list. Must exist in the
           user's list of emails.
    :param cla_manager_name: The name of the CLA manager
    :param cla_manager_email: The email address of the CLA manager
    :param project_name: The name of the project
    :param company_name: The name of the organization/company
    :param account_exists: boolean to check whether the email is being sent to a proposed admin(false), or an admin for
           an existing company(true).
     """

    # account_exists=True send email to the CLA Manager of the existing company
    # account_exists=False send email to a proposed CLA Manager who needs to register the company through
    # the Corporate Console.
    subject = f'EasyCLA: Request to start CLA signature process for {project_name}'
    body = f'''
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello {cla_manager_name},</p>
<p>This is a notification email from EasyCLA regarding the project {project_name}.</p>
<p>{project_name} uses EasyCLA to ensure that before a contribution is accepted, the contributor is covered under a
signed CLA.</p>
<p>{contributor_name} ({contributor_email}) has designated you as the proposed initial CLA Manager for contributions
from {company_name if company_name else 'your company'} to {project_name}. This would mean that, after the CLA is
signed, you would be able to maintain the list of employees allowed to contribute to {project_name} on behalf of your
company, as well as the list of your company’s CLA Managers for {project_name}.</p>
<p>If you can be the initial CLA Manager from your company for {project_name}, please log into the EasyCLA Corporate
Console at {cla.conf['CLA_LANDING_PAGE']} to begin the CLA signature process. You might not be authorized to sign the
CLA yourself on behalf of your company; if not, the signature process will prompt you to designate somebody else who is
authorized to sign the CLA.</p>
<p>If you need help or have questions about EasyCLA, you can
<a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or
<a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for
support</a>.</p>

<p>Thanks,</p>
<p>EasyCLA support team</p>
</body>
</html>
'''
    recipient = cla_manager_email
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
        last_signature['requires_resigning'] = last_signature['latest_document_major_version'] != last_signature[
            'signature_document_major_version']
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
        last_signature['requires_resigning'] = last_signature['latest_document_major_version'] != last_signature[
            'signature_document_major_version']
    return last_signature


# For GitHub user creating, see models.github_models.get_or_create_user(self, request)
def get_or_create_user(auth_user):
    user = User()

    # Returns None or List[User] objects - could be more than one
    users = user.get_user_by_username(str(auth_user.username))

    if users is None:
        user.set_user_id(str(uuid.uuid4()))
        user.set_user_name(auth_user.name)
        user.set_lf_email(auth_user.email.lower())
        user.set_lf_username(auth_user.username)
        user.set_lf_sub(auth_user.sub)

        user.save()

        event_data = f'CLA user added for {auth_user.username}'
        Event.create_event(
            event_data=event_data,
            event_type=EventType.CreateUser,
            contains_pii=True,
        )

        return user

    # Just return the first matching record
    return users[0]


def request_company_admin_access(user_id, company_id):
    """
    Send Email to company admins to inform that that a user is requesting to be a CLA Manager for their company.
    """

    user = User()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    user_name = user.get_user_name()

    company = Company()
    try:
        company.load(company_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    subject = 'EasyCLA: New CLA Manager Access Request'

    # Send emails to every CLA manager
    for admin in company.get_managers():
        # TODO: Review this - need to load project
        body = f'''
<html>
<head>
<style>
body {{font-family: Arial, Helvetica, sans-serif; font-size: 1.2em;}}
</style>
</head>
<body>
<p>Hello {admin.get_user_name()},
<p>This is a notification email from EasyCLA regarding the project {project}.
<p>You are currently listed as a CLA Manager from {company.get_company_name()} for the project {project}. This means that you are able to maintain the list of employees allowed to contribute to {project} on behalf of your company, as well as the list of your company’s CLA Managers for {project}.
<p>{user_name} ({user.get_user_email}) has requested to be added as another CLA Manager from {company.get_company_name()} for {project}. This would permit them to maintain the lists of approved contributors and CLA Managers as well.
<p>If you want to permit this, please log into the EasyCLA Corporate Console at https://{corporate_console_url}, where you can approve this user as an additional CLA Manager.
<p>If you need help or have questions about EasyCLA, you can <a href="https://docs.linuxfoundation.org/docs/communitybridge/communitybridge-easycla" target="_blank">read the documentation</a> or <a href="https://jira.linuxfoundation.org/servicedesk/customer/portal/4/create/143" target="_blank">reach out to us for support</a>.</p>
<p>Thanks,
<p>EasyCLA support team</p>
</body>
</html>
'''
        recipient = admin.get_lf_email()
        email_service = get_email_service()
        email_service.send(subject, body, recipient)
