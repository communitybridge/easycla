"""
Controller related to agreement operations.
"""

import uuid
import hug.types
import cla.hug_types
from cla.utils import get_agreement_instance, get_user_instance, get_organization_instance, \
                      get_project_instance, get_email_service
from cla.models import DoesNotExist

def get_agreements():
    """
    Returns a list of agreements in the CLA system.

    :return: List of agreements in dict format.
    :rtype: [dict]
    """
    agreements = [agreement.to_dict() for agreement in get_agreement_instance().all()]
    return agreements

def get_agreement(agreement_id):
    """
    Returns the CLA agreement requested by UUID.

    :param agreement_id: The agreement UUID.
    :type agreement_id: UUID
    :return: dict representation of the agreement object.
    :rtype: dict
    """
    agreement = get_agreement_instance()
    try:
        agreement.load(agreement_id=str(agreement_id))
    except DoesNotExist as err:
        return {'errors': {'agreement_id': str(err)}}
    return agreement.to_dict()

def create_agreement(agreement_project_id, # pylint: disable=too-many-arguments
                     agreement_reference_id,
                     agreement_reference_type,
                     agreement_type='cla',
                     agreement_approved=False,
                     agreement_signed=False,
                     agreement_return_url=None,
                     agreement_sign_url=None):
    """
    Creates an agreement and returns the newly created agreement in dict format.

    :param agreement_project_id: The project ID for this new agreement.
    :type agreement_project_id: string
    :param agreement_reference_id: The user or organization ID for this agreement.
    :type agreement_reference_id: string
    :param agreement_reference_type: The type of reference ('user' or 'organization')
    :type agreement_reference_type: string
    :param agreement_type: The agreement type ('cla' or 'dco')
    :type agreement_type: string
    :param agreement_signed: Whether or not the agreement has been signed.
    :type agreement_signed: boolean
    :param agreement_approved: Whether or not the agreement has been approved.
    :type agreement_approved: boolean
    :param agreement_return_url: The URL the user will be redirected to after signing.
    :type agreement_return_url: string
    :param agreement_sign_url: The URL the user must visit to sign the agreement.
    :type agreement_sign_url: string
    :return: A dict of a newly created agreement.
    :rtype: dict
    """
    agreement = get_agreement_instance()
    agreement.set_agreement_id(str(uuid.uuid4()))
    project = get_project_instance()
    try:
        project.load(project_id=str(agreement_project_id))
    except DoesNotExist as err:
        return {'errors': {'agreement_project_id': str(err)}}
    agreement.set_agreement_project_id(str(agreement_project_id))
    if agreement_reference_type == 'user':
        user = get_user_instance()
        try:
            user.load(agreement_reference_id)
        except DoesNotExist as err:
            return {'errors': {'agreement_reference_id': str(err)}}
        try:
            document = project.get_project_individual_document()
        except DoesNotExist as err:
            return {'errors': {'agreement_project_id': str(err)}}
    else:
        organization = get_organization_instance()
        try:
            organization.load(agreement_reference_id)
        except DoesNotExist as err:
            return {'errors': {'agreement_reference_id': str(err)}}
        try:
            document = project.get_project_corporate_document()
        except DoesNotExist as err:
            return {'errors': {'agreement_project_id': str(err)}}
    agreement.set_agreement_document_revision(document.get_document_revision())
    agreement.set_agreement_reference_id(str(agreement_reference_id))
    agreement.set_agreement_reference_type(agreement_reference_type)
    agreement.set_agreement_type(agreement_type)
    agreement.set_agreement_signed(agreement_signed)
    agreement.set_agreement_approved(agreement_approved)
    agreement.set_agreement_return_url(agreement_return_url)
    agreement.set_agreement_sign_url(agreement_sign_url)
    agreement.save()
    return agreement.to_dict()

def update_agreement(agreement_id, # pylint: disable=too-many-arguments,too-many-return-statements,too-many-branches
                     agreement_project_id=None,
                     agreement_reference_id=None,
                     agreement_reference_type=None,
                     agreement_type=None,
                     agreement_approved=None,
                     agreement_signed=None,
                     agreement_return_url=None,
                     agreement_sign_url=None):
    """
    Updates an agreement and returns the newly updated agreement in dict format.
    A value of None means the field should not be updated.

    :param agreement_id: ID of the agreement.
    :type agreement_id: ID | None
    :param agreement_project_id: Project ID for this agreement.
    :type agreement_project_id: string | None
    :param agreement_reference_id: Reference ID for this agreement.
    :type agreement_reference_id: string | None
    :param agreement_reference_type: Reference type for this agreement.
    :type agreement_reference_type: ['user' | 'organization'] | None
    :param agreement_type: New agreement type ('cla' or 'dco').
    :type agreement_type: string | None
    :param agreement_signed: Whether this agreement is signed or not.
    :type agreement_signed: boolean | None
    :param agreement_approved: Whether this agreement is approved or not.
    :type agreement_approved: boolean | None
    :param agreement_return_url: The URL the user will be sent to after signing.
    :type agreement_return_url: string | None
    :param agreement_sign_url: The URL the user must visit to sign the agreement.
    :type agreement_sign_url: string | None
    :return: dict representation of the agreement object.
    :rtype: dict
    """
    agreement = get_agreement_instance()
    try: # Try to load the agreement to update.
        agreement.load(str(agreement_id))
    except DoesNotExist as err:
        return {'errors': {'agreement_id': str(err)}}
    if agreement_project_id is not None:
        try:
            agreement.set_agreement_project_id(agreement_project_id)
        except DoesNotExist as err:
            return {'errors': {'agreement_project_id': str(err)}}
    # TODO: Ensure agreement_reference_id exists.
    if agreement_reference_id is not None:
        agreement.set_agreement_reference_id(agreement_reference_id)
    if agreement_reference_type is not None:
        agreement.set_agreement_reference_type(agreement_reference_type)
    if agreement_type is not None:
        if agreement_type in ['cla', 'dco']:
            agreement.set_agreement_type(agreement_type)
        else:
            return {'errors': {'agreement_type': \
                               'Invalid value passed. The accepted values are: (cla|dco)'}}
    if agreement_signed is not None:
        try:
            val = hug.types.smart_boolean(agreement_signed)
            agreement.set_agreement_signed(val)
        except KeyError as err:
            return {'errors': {'agreement_signed':
                               'Invalid value passed in for true/false field'}}
    if agreement_approved is not None:
        try:
            val = hug.types.smart_boolean(agreement_approved)
            update_agreement_approved(agreement, val)
        except KeyError as err:
            return {'errors': {'agreement_approved':
                               'Invalid value passed in for true/false field'}}
    if agreement_return_url is not None:
        try:
            val = cla.hug_types.url(agreement_return_url)
            agreement.set_agreement_return_url(val)
        except KeyError as err:
            return {'errors': {'agreement_return_url':
                               'Invalid value passed in for URL field'}}
    if agreement_sign_url is not None:
        try:
            val = cla.hug_types.url(agreement_sign_url)
            agreement.set_agreement_sign_url(val)
        except KeyError as err:
            return {'errors': {'agreement_sign_url':
                               'Invalid value passed in for URL field'}}
    agreement.save()
    return agreement.to_dict()

def update_agreement_approved(agreement, value):
    """Helper function to update the agreement approval status and send emails if necessary."""
    previous = agreement.get_agreement_approved()
    agreement.set_agreement_approved(value)
    email_approval = cla.conf['EMAIL_ON_AGREEMENT_APPROVED']
    if email_approval and not previous and value: # Just got approved.
        subject, body, recipients = get_agreement_approved_email_content(agreement)
        get_email_service().send(subject, body, recipients)

def get_agreement_approved_email_content(agreement): # pylint: disable=invalid-name
    """Helper function to get agreement approval email subject, body, and recipients."""
    if agreement.get_agreement_reference_type() != 'user':
        cla.log.info('Not sending agreement approved emails for CCLAs')
        return
    subject = 'CLA Agreement Approved'
    user = get_user_instance()
    user.load(agreement.get_agreement_reference_id())
    project = get_project_instance()
    project.load(agreement.get_agreement_project_id())
    recipients = [user.get_user_id()]
    body = 'Hello %s. Your Contributor License Agreement for %s has been approved!' \
           %(user.get_user_name(), project.get_project_name())
    return subject, body, recipients

def delete_agreement(agreement_id):
    """
    Deletes an agreement based on UUID.

    :param agreement_id: The UUID of the agreement.
    :type agreement_id: UUID
    """
    agreement = get_agreement_instance()
    try: # Try to load the agreement to delete.
        agreement.load(str(agreement_id))
    except DoesNotExist as err:
        # Should we bother sending back an error?
        return {'errors': {'agreement_id': str(err)}}
    agreement.delete()
    return {'success': True}

def get_user_agreements(user_id):
    """
    Get all agreements for user.

    :param user_id: The ID of the user in question.
    :type user_id: string
    """
    agreements = get_agreement_instance().get_agreements_by_reference(user_id, 'user')
    return [agreement.to_dict() for agreement in agreements]

def get_organization_agreements(organization_id):
    """
    Get all agreements for organization.

    :param organization_id: The ID of the organization in question.
    :type organization_id: string
    """
    agreements = get_agreement_instance().get_agreements_by_reference(organization_id,
                                                                      'organization')
    return [agreement.to_dict() for agreement in agreements]

def get_project_agreements(project_id):
    """
    Get all agreements for project.

    :param project_id: The ID of the project in question.
    :type project_id: string
    """
    agreements = get_agreement_instance().get_agreements_by_project(project_id)
    return [agreement.to_dict() for agreement in agreements]
