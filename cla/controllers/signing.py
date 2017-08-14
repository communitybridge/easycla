"""
Controller related to the signed callback.
"""

import falcon
import cla
from cla.utils import get_signing_service, get_agreement_instance
from cla.models import DoesNotExist

def request_signature(project_id, user_id, ret_url, callback_url=None):
    """
    Handle POST request to send signature request to user.

    :param project_id: The project to sign for.
    :type project_id: string
    :param user_id: The ID of the user that will sign.
    :type user_id: string
    :param ret_url: The URL the user will be sent to after successful signature.
    :type ret_url: string
    :param callback_url: Optional parameter to specify the URL the signing service should hit
        after confirmation of successful signature.
    :type callback_url: string
    """
    return get_signing_service().request_signature(project_id,
                                                   user_id,
                                                   ret_url,
                                                   callback_url)

def post_signed(content, repository_id, change_request_id):
    """
    Handle the posted callback from the signing service.

    :param content: The POST body from the signing service callback.
    :type content: string
    :param repository_id: The ID of the repository that this signature was requested for.
    :type repository_id: string
    :param change_request_id: The ID of the change request or pull request that
        initiated this signature.
    :type change_request_id: string
    """
    get_signing_service().signed_callback(content, repository_id, change_request_id)

def return_url(agreement_id, event=None): # pylint: disable=unused-argument
    """
    Handle the GET request from the user once they have successfully signed.

    :param agreement_id: The ID of the agreement they have just signed.
    :type agreement_id: string
    :param event: The event GET flag sent back from the signing service provider.
    :type event: string | None
    """
    try: # Load the agreement based on ID.
        agreement = get_agreement_instance()
        agreement.load(str(agreement_id))
    except DoesNotExist as err:
        cla.log.error('Invalid agreement_id provided when trying to send user back to their ' + \
                      'return_url after signing: %s', agreement_id)
        return {'errors': {'agreement_id': str(err)}}
    # Ensure everything went well on the signing service provider's side.
    if event is not None:
        # Expired signing URL - the user was redirected back immediately but still needs to sign.
        if event == 'ttl_expired' and agreement.get_agreement_signed() is False:
            # Need to re-generate a sign_url and try again.
            cla.log.info('DocuSign URL used was expired, re-generating sign_url')
            callback_url = agreement.get_agreement_callback_url()
            get_signing_service().populate_sign_url(agreement, callback_url)
            agreement.save()
            raise falcon.HTTPFound(agreement.get_agreement_sign_url())
    ret_url = agreement.get_agreement_return_url()
    raise falcon.HTTPFound(ret_url)
