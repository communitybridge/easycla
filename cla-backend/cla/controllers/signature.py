# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Controller related to signature operations.
"""

import uuid
import hug.types
import cla.hug_types
from cla.utils import get_email_service
from cla.models import DoesNotExist
from cla.models.dynamo_models import User, Project, Signature, Company
from cla.controllers import company

def get_signatures():
    """
    Returns a list of signatures in the CLA system.

    :return: List of signatures in dict format.
    :rtype: [dict]
    """
    signatures = [signature.to_dict() for signature in Signature().all()]
    return signatures

def get_signature(signature_id):
    """
    Returns the CLA signature requested by UUID.

    :param signature_id: The signature UUID.
    :type signature_id: UUID
    :return: dict representation of the signature object.
    :rtype: dict
    """
    signature = Signature()
    try:
        signature.load(signature_id=str(signature_id))
    except DoesNotExist as err:
        return {'errors': {'signature_id': str(err)}}
    return signature.to_dict()

def create_signature(signature_project_id, # pylint: disable=too-many-arguments
                     signature_reference_id,
                     signature_reference_type,
                     signature_type='cla',
                     signature_approved=False,
                     signature_signed=False,
                     signature_return_url=None,
                     signature_sign_url=None,
                     signature_user_ccla_company_id=None,
                     signature_acl=None):
    """
    Creates an signature and returns the newly created signature in dict format.

    :param signature_project_id: The project ID for this new signature.
    :type signature_project_id: string
    :param signature_reference_id: The user or company ID for this signature.
    :type signature_reference_id: string
    :param signature_reference_type: The type of reference ('user' or 'company')
    :type signature_reference_type: string
    :param signature_type: The signature type ('cla' or 'dco')
    :type signature_type: string
    :param signature_signed: Whether or not the signature has been signed.
    :type signature_signed: boolean
    :param signature_approved: Whether or not the signature has been approved.
    :type signature_approved: boolean
    :param signature_return_url: The URL the user will be redirected to after signing.
    :type signature_return_url: string
    :param signature_sign_url: The URL the user must visit to sign the signature.
    :type signature_sign_url: string
    :param signature_user_ccla_company_id: The company ID if creating an employee signature.
    :type signature_user_ccla_company_id: string
    :return: A dict of a newly created signature.
    :rtype: dict
    """
    signature = Signature()
    signature.set_signature_id(str(uuid.uuid4()))
    project = Project()
    try:
        project.load(project_id=str(signature_project_id))
    except DoesNotExist as err:
        return {'errors': {'signature_project_id': str(err)}}
    signature.set_signature_project_id(str(signature_project_id))
    if signature_reference_type == 'user':
        user = User()
        try:
            user.load(signature_reference_id)
        except DoesNotExist as err:
            return {'errors': {'signature_reference_id': str(err)}}
        try:
            document = project.get_project_individual_document()
        except DoesNotExist as err:
            return {'errors': {'signature_project_id': str(err)}}
    else:
        company = Company()
        try:
            company.load(signature_reference_id)
        except DoesNotExist as err:
            return {'errors': {'signature_reference_id': str(err)}}
        try:
            document = project.get_project_corporate_document()
        except DoesNotExist as err:
            return {'errors': {'signature_project_id': str(err)}}

    # Set username to this signature ACL
    if signature_acl is not None:
        signature.set_signature_acl(signature_acl)

    signature.set_signature_document_minor_version(document.get_document_minor_version())
    signature.set_signature_document_major_version(document.get_document_major_version())
    signature.set_signature_reference_id(str(signature_reference_id))
    signature.set_signature_reference_type(signature_reference_type)
    signature.set_signature_type(signature_type)
    signature.set_signature_signed(signature_signed)
    signature.set_signature_approved(signature_approved)
    signature.set_signature_return_url(signature_return_url)
    signature.set_signature_sign_url(signature_sign_url)
    if signature_user_ccla_company_id is not None:
        signature.set_signature_user_ccla_company_id(str(signature_user_ccla_company_id))
    signature.save()
    return signature.to_dict()

def update_signature(signature_id, # pylint: disable=too-many-arguments,too-many-return-statements,too-many-branches
                     signature_project_id=None,
                     signature_reference_id=None,
                     signature_reference_type=None,
                     signature_type=None,
                     signature_approved=None,
                     signature_signed=None,
                     signature_return_url=None,
                     signature_sign_url=None,
                     domain_whitelist=None,
                     email_whitelist=None,
                     github_whitelist=None):
    """
    Updates an signature and returns the newly updated signature in dict format.
    A value of None means the field should not be updated.

    :param signature_id: ID of the signature.
    :type signature_id: ID | None
    :param signature_project_id: Project ID for this signature.
    :type signature_project_id: string | None
    :param signature_reference_id: Reference ID for this signature.
    :type signature_reference_id: string | None
    :param signature_reference_type: Reference type for this signature.
    :type signature_reference_type: ['user' | 'company'] | None
    :param signature_type: New signature type ('cla' or 'dco').
    :type signature_type: string | None
    :param signature_signed: Whether this signature is signed or not.
    :type signature_signed: boolean | None
    :param signature_approved: Whether this signature is approved or not.
    :type signature_approved: boolean | None
    :param signature_return_url: The URL the user will be sent to after signing.
    :type signature_return_url: string | None
    :param signature_sign_url: The URL the user must visit to sign the signature.
    :type signature_sign_url: string | None
    :return: dict representation of the signature object.
    :rtype: dict
    """
    signature = Signature()
    try: # Try to load the signature to update.
        signature.load(str(signature_id))
    except DoesNotExist as err:
        return {'errors': {'signature_id': str(err)}}
    if signature_project_id is not None:
        try:
            signature.set_signature_project_id(str(signature_project_id))
        except DoesNotExist as err:
            return {'errors': {'signature_project_id': str(err)}}
    # TODO: Ensure signature_reference_id exists.
    if signature_reference_id is not None:
        signature.set_signature_reference_id(signature_reference_id)
    if signature_reference_type is not None:
        signature.set_signature_reference_type(signature_reference_type)
    if signature_type is not None:
        if signature_type in ['cla', 'dco']:
            signature.set_signature_type(signature_type)
        else:
            return {'errors': {'signature_type': \
                               'Invalid value passed. The accepted values are: (cla|dco)'}}
    if signature_signed is not None:
        try:
            val = hug.types.smart_boolean(signature_signed)
            signature.set_signature_signed(val)
        except KeyError as err:
            return {'errors': {'signature_signed':
                               'Invalid value passed in for true/false field'}}
    if signature_approved is not None:
        try:
            val = hug.types.smart_boolean(signature_approved)
            update_signature_approved(signature, val)
        except KeyError as err:
            return {'errors': {'signature_approved':
                               'Invalid value passed in for true/false field'}}
    if signature_return_url is not None:
        try:
            val = cla.hug_types.url(signature_return_url)
            signature.set_signature_return_url(val)
        except KeyError as err:
            return {'errors': {'signature_return_url':
                               'Invalid value passed in for URL field'}}
    if signature_sign_url is not None:
        try:
            val = cla.hug_types.url(signature_sign_url)
            signature.set_signature_sign_url(val)
        except KeyError as err:
            return {'errors': {'signature_sign_url':
                               'Invalid value passed in for URL field'}}

    if domain_whitelist is not None:
        try:
            domain_whitelist = hug.types.multiple(domain_whitelist)
            signature.set_domain_whitelist(domain_whitelist)
        except KeyError as err:
            return {'errors': {
                'domain_whitelist': 'Invalid value passed in for the domain whitelist'
            }}

    if email_whitelist is not None:
        try:
            email_whitelist = hug.types.multiple(email_whitelist)
            signature.set_email_whitelist(email_whitelist)
        except KeyError as err:
            return {'errors': {
                'email_whitelist': 'Invalid value passed in for the email whitelist'
            }}

    if github_whitelist is not None:
        try:
            github_whitelist = hug.types.multiple(github_whitelist)
            signature.set_github_whitelist(github_whitelist)
        except KeyError as err:
            return {'errors': {
                'github_whitelist': 'Invalid value passed in for the github whitelist'
            }}

    signature.save()
    return signature.to_dict()

def update_signature_approved(signature, value):
    """Helper function to update the signature approval status and send emails if necessary."""
    previous = signature.get_signature_approved()
    signature.set_signature_approved(value)
    email_approval = cla.conf['EMAIL_ON_SIGNATURE_APPROVED']
    if email_approval and not previous and value: # Just got approved.
        subject, body, recipients = get_signature_approved_email_content(signature)
        get_email_service().send(subject, body, recipients)

def get_signature_approved_email_content(signature): # pylint: disable=invalid-name
    """Helper function to get signature approval email subject, body, and recipients."""
    if signature.get_signature_reference_type() != 'user':
        cla.log.info('Not sending signature approved emails for CCLAs')
        return
    subject = 'CLA Signature Approved'
    user = User()
    user.load(signature.get_signature_reference_id())
    project = Project()
    project.load(signature.get_signature_project_id())
    recipients = [user.get_user_id()]
    body = 'Hello %s. Your Contributor License Agreement for %s has been approved!' \
           %(user.get_user_name(), project.get_project_name())
    return subject, body, recipients

def delete_signature(signature_id):
    """
    Deletes an signature based on UUID.

    :param signature_id: The UUID of the signature.
    :type signature_id: UUID
    """
    signature = Signature()
    try: # Try to load the signature to delete.
        signature.load(str(signature_id))
    except DoesNotExist as err:
        # Should we bother sending back an error?
        return {'errors': {'signature_id': str(err)}}
    signature.delete()
    return {'success': True}

def get_user_signatures(user_id):
    """
    Get all signatures for user.

    :param user_id: The ID of the user in question.
    :type user_id: string
    """
    signatures = Signature().get_signatures_by_reference(str(user_id), 'user')
    return [signature.to_dict() for signature in signatures]

def get_user_project_signatures(user_id, project_id, signature_type=None):
    """
    Get all signatures for user filtered by a project.

    :param user_id: The ID of the user in question.
    :type user_id: string
    :param project_id: The ID of the project to filter by.
    :type project_id: string
    :param signature_type: The signature type to filter by.
    :type signature_type: string (one of 'individual', 'employee')
    :return: The list of signatures requested.
    :rtype: [cla.models.model_interfaces.Signature]
    """
    sig = Signature()
    signatures = sig.get_signatures_by_project(str(project_id),
                                               signature_reference_type='user',
                                               signature_reference_id=str(user_id))
    ret = []
    for signature in signatures:
        if signature_type is not None:
            if signature_type == 'individual' and \
               signature.get_signature_user_ccla_employee_id() is not None:
                continue
            elif signature_type == 'employee' and \
                 signature.get_signature_user_ccla_employee_id() is None:
                continue
        ret.append(signature.to_dict())
    return ret

def get_company_signatures(company_id):
    """
    Get all signatures for company.

    :param company_id: The ID of the company in question.
    :type company_id: string
    """
    signatures = Signature().get_signatures_by_reference(company_id,
                                                                      'company')

    return [signature.to_dict() for signature in signatures]

def get_company_signatures_by_acl(username, company_id):
    """
    Get all signatures for company filtered by it's ACL.
    A company's signature will be returned only if the provided
    username appears in the signature's ACL.

    :param username: The username of the authenticated user
    :type username: string
    :param company_id: The ID of the company in question.
    :type company_id: string
    """
    # Get signatures by company reference
    all_signatures = Signature().get_signatures_by_reference(company_id, 'company')

    # Filter signatures this manager is authorized to see
    signatures = []
    for signature in all_signatures:
        if username in signature.get_signature_acl():
            signatures.append(signature)

    return [signature.to_dict() for signature in signatures]

def get_project_signatures(project_id):
    """
    Get all signatures for project.

    :param project_id: The ID of the project in question.
    :type project_id: string
    """
    signatures = Signature().get_signatures_by_project(str(project_id))
    return [signature.to_dict() for signature in signatures]


def get_project_company_signatures(company_id, project_id):
    """
    Get all company signatures for project specified and a company specified

    :param company_id: The ID of the company in question
    :param project_id: The ID of the project in question
    :type company_id: string
    :type project_id: string
    """
    signatures = Signature().get_signatures_by_company_project(str(company_id),
                                                                            str(project_id))
    return signatures

def get_project_employee_signatures(company_id, project_id):
    """
    Get all employee signatures for project specified and a company specified

    :param company_id: The ID of the company in question
    :param project_id: The ID of the project in question
    :type company_id: string
    :type project_id: string
    """
    signatures = Signature().get_employee_signatures_by_company_project(str(company_id),
                                                                            str(project_id))
    return signatures

def get_cla_managers(username, signature_id):
    """
    Returns CLA managers from the CCLA signature ID.

    :param username: The LF username
    :type username: string
    :param signature_id: The Signature ID of the CCLA signed. 
    :type signature_id: string
    :return: dict representation of the project managers.
    :rtype: dict
    """
    signature = Signature()
    try:
        signature.load(str(signature_id))
    except DoesNotExist as err:
        return {'errors': {'signature_id': str(err)}}

    # Get Signature ACL
    signature_acl = signature.get_signature_acl()

    if username not in signature_acl:
        return {'errors': {'user_id': 'You are not authorized to see the managers.'}}

    return get_managers_dict(signature_acl)

def add_cla_manager(auth_user, signature_id, lfid):
    """
    Adds the LFID to the signature ACL and returns a new list of CLA Managers. 

    :param username: username of the user
    :type username: string
    :param signature_id: The ID of the project
    :type signature_id: UUID
    :param lfid: the lfid (manager username) to be added to the project acl
    :type lfid: string
    """
    # Find project
    signature = Signature()
    try:
        signature.load(str(signature_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}

    # Get Signature ACL
    signature_acl = signature.get_signature_acl()

    if auth_user.username not in signature_acl:
        return {'errors': {'user_id': 'You are not authorized to see the managers.'}}

    company.add_permission(auth_user, lfid, signature.get_signature_reference_id(), ignore_auth_user=True)

    # Add lfid to acl
    signature.add_signature_acl(lfid)
    signature.save()

    return get_managers_dict(signature_acl)

def remove_cla_manager(username, signature_id, lfid):
    """
    Removes the LFID from the project ACL

    :param username: username of the user
    :type username: string
    :param project_id: The ID of the project
    :type project_id: UUID
    :param lfid: the lfid (manager username) to be removed to the project acl
    :type lfid: string
    """
    # Find project
    signature = Signature()
    try:
        signature.load(str(signature_id))
    except DoesNotExist as err:
        return {'errors': {'signature_id': str(err)}}

    # Validate user is the manager of the project
    signature_acl = signature.get_signature_acl()
    if username not in signature_acl:
        return {'errors': {'user': "You are not authorized to manage this CCLA."}}

    # Avoid to have an empty acl
    if len(signature_acl) == 1 and username == lfid:
        return {'errors': {'user': "You cannot remove this manager because a CCLA must have at least one CLA manager."}}
    
    # Remove LFID from the acl
    signature.remove_signature_acl(lfid)
    signature.save()

    # Return modified managers
    return get_managers_dict(signature_acl)


def get_managers_dict(signature_acl):
    # Helper function to get a list of all cla managers from a CCLA Signature ACL
    # Generate managers dict
    managers_dict = []
    for lfid in signature_acl:
        user = User()
        user = user.get_user_by_username(str(lfid))
        if user is not None:
            # Manager found, fill with it's information
            managers_dict.append({
                'name': user.get_user_name(),
                'email': user.get_user_email(),
                'lfid': user.get_lf_username()
            })
        else:
            # Manager not in database yet, only set the lfid
            managers_dict.append({
                'lfid': str(lfid)
            })

    return managers_dict
