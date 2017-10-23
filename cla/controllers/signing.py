"""
Controller related to the signed callback.
"""

import uuid
import falcon
import cla
from cla.utils import get_signing_service, get_signature_instance, get_project_instance, \
                      get_company_instance, get_user_instance
from cla.models import DoesNotExist

def request_signature(project_id, user_id, return_url=None):
    """
    Handle POST request to send signature request to user.

    :param project_id: The project to sign for.
    :type project_id: string
    :param user_id: The ID of the user that will sign.
    :type user_id: string
    :param return_url: The URL to return the user to after signing is complete.
    :type return_url: string
    """
    return get_signing_service().request_signature(str(project_id), str(user_id), return_url)

def employee_signature(project_id, company_id, user_id):
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    company = get_company_instance()
    try:
        company.load(str(company_id))
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}
    user = get_user_instance()
    try:
        user.load(str(user_id))
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    # Ensure the company actually has a CCLA with this project.
    existing_signatures = get_signature_instance().get_signatures_by_project(
        project_id,
        signature_type='ccla',
        signature_reference_type='company',
        signature_reference_id=company.get_company_id()
    )
    if len(existing_signatures) < 1:
        return {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
    # Ensure user hasn't already signed this signature.
    existing_signatures = get_signature_instance().get_signatures_by_project(
        project_id,
        signature_type='ccla',
        signature_reference_type='user',
        signature_reference_id=user_id,
        signature_user_ccla_company_id=company_id
    )
    if len(existing_signatures) > 0:
        return {'errors': {'signature_id': 'User has already signed CCLA with this company'}}
    # Create the new Signature.
    new_signature = get_signature_instance()
    new_signature.set_signature_id(str(uuid.uuid4()))
    new_signature.set_signature_project_id(str(project_id))
    new_signature.set_signature_document_major_version(0)
    new_signature.set_signature_document_minor_version(0)
    new_signature.set_signature_signed(True)
    new_signature.set_signature_approved(True)
    new_signature.set_signature_type('ccla')
    new_signature.set_signature_reference_type('user')
    new_signature.set_signature_reference_id(user_id)
    new_signature.set_signature_user_ccla_company_id(company_id)
    new_signature.save()
    return new_signature.to_dict()

def post_signed(content, installation_id, github_repository_id, change_request_id):
    """
    Handle the posted callback from the signing service.

    :TODO: Update comments.

    :param content: The POST body from the signing service callback.
    :type content: string
    :param repository_id: The ID of the repository that this signature was requested for.
    :type repository_id: string
    :param change_request_id: The ID of the change request or pull request that
        initiated this signature.
    :type change_request_id: string
    """
    get_signing_service().signed_callback(content, installation_id, github_repository_id, change_request_id)

def return_url(signature_id, event=None): # pylint: disable=unused-argument
    """
    Handle the GET request from the user once they have successfully signed.

    :param signature_id: The ID of the signature they have just signed.
    :type signature_id: string
    :param event: The event GET flag sent back from the signing service provider.
    :type event: string | None
    """
    try: # Load the signature based on ID.
        signature = get_signature_instance()
        signature.load(str(signature_id))
    except DoesNotExist as err:
        cla.log.error('Invalid signature_id provided when trying to send user back to their ' + \
                      'return_url after signing: %s', signature_id)
        return {'errors': {'signature_id': str(err)}}
    # Ensure everything went well on the signing service provider's side.
    if event is not None:
        # Expired signing URL - the user was redirected back immediately but still needs to sign.
        if event == 'ttl_expired' and signature.get_signature_signed() is False:
            # Need to re-generate a sign_url and try again.
            cla.log.info('DocuSign URL used was expired, re-generating sign_url')
            callback_url = signature.get_signature_callback_url()
            get_signing_service().populate_sign_url(signature, callback_url)
            signature.save()
            raise falcon.HTTPFound(signature.get_signature_sign_url())
    ret_url = signature.get_signature_return_url()
    raise falcon.HTTPFound(ret_url)
