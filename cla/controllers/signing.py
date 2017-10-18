"""
Controller related to the signed callback.
"""

import falcon
import cla
from cla.utils import get_signing_service, get_signature_instance
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
