# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Controller related to the signed callback.
"""

import uuid
import hug
import falcon
import cla
from cla.utils import get_signing_service, get_signature_instance, get_project_instance, \
                      get_company_instance, get_user_instance, get_email_service
from cla.models import DoesNotExist

def request_individual_signature(project_id, user_id, return_url_type, return_url=None):
    """
    Handle POST request to send ICLA signature request to user.

    :param project_id: The project to sign for.
    :type project_id: string
    :param user_id: The ID of the user that will sign.
    :type user_id: string
    :param return_url_type: Refers to the return url provider type: Gerrit or Github
    :type return_url_type: string
    :param return_url: The URL to return the user to after signing is complete.
    :type return_url: string
    """
    signing_service = get_signing_service()
    if return_url_type == "Gerrit":
        return signing_service.request_individual_signature_gerrit(str(project_id), str(user_id), return_url)
    elif return_url_type == "Github":
        return signing_service.request_individual_signature(str(project_id), str(user_id), return_url)

def request_corporate_signature(auth_user, project_id, company_id, send_as_email=False, 
                                authority_name=None, authority_email=None, return_url_type=None, return_url=None):
    """
    Creates CCLA signature object that represents a company signing a CCLA.

    :param project_id: The ID of the project the company is signing a CCLA for.
    :type project_id: string
    :param company_id: The ID of the company that is signing the CCLA.
    :type company_id: string
    :param return_url: The URL to return the user to after signing is complete.
    :type return_url: string
    """
    return get_signing_service().request_corporate_signature(auth_user, str(project_id), str(company_id), send_as_email, authority_name, authority_email, return_url_type, return_url)

def request_employee_signature(project_id, company_id, user_id, return_url_type, return_url=None):
    """
    Creates placeholder signature object that represents a user signing a CCLA as an employee.

    :param project_id: The ID of the project the user is signing a CCLA for.
    :type project_id: string
    :param company_id: The ID of the company the employee belongs to.
    :type company_id: string
    :param user_id: The ID of the user.
    :type user_id: string
    :param return_url_type: Refers to the return url provider type: Gerrit or Github
    :type return_url_type: string
    :param return_url: The URL to return the user to after signing is complete.
    """

    signing_service = get_signing_service()
    if return_url_type == "Gerrit":
        return signing_service.request_employee_signature_gerrit(str(project_id), str(company_id), str(user_id), return_url)
    elif return_url_type == "Github":
        return signing_service.request_employee_signature(str(project_id), str(company_id), str(user_id), return_url)

def check_and_prepare_employee_signature(project_id, company_id, user_id):
    """
    Checks that 
    1. The given project, company, and user exists 
    2. The company signatory has signed the CCLA for their company. 
    3. The user is included as part of the whitelist of the CCLA that the company signed. 

    :param project_id: The ID of the project the user is signing a CCLA for.
    :type project_id: string
    :param company_id: The ID of the company the employee belongs to.
    :type company_id: string
    :param user_id: The ID of the user.
    :type user_id: string
    """
    return get_signing_service().check_and_prepare_employee_signature(str(project_id), str(company_id), str(user_id))

# Deprecated in favor of sending the email through DocuSign
def send_authority_email(company_name, project_name, authority_name, authority_email):
    """
    Sends email to the specified corporate authority to sign the CCLA Docusign file. 
    """

    subject = 'CLA: Invitation to Sign a Corporate Contributor License Agreement'
    body = '''Hello %s, 
    
Your organization: %s, 
    
has requested a Corporate Contributor License Agreement Form to be signed for the following project:

%s

Please read the agreement carefully and sign the attached file. 
    

- Linux Foundation CLA System
''' %(authority_name, company_name, project_name)
    recipient = authority_email
    email_service = get_email_service()
    email_service.send(subject, body, recipient)

def post_individual_signed(content, installation_id, github_repository_id, change_request_id):
    """
    Handle the posted callback from the signing service after ICLA signature.

    :param content: The POST body from the signing service callback.
    :type content: string
    :param repository_id: The ID of the repository that this signature was requested for.
    :type repository_id: string
    :param change_request_id: The ID of the change request or pull request that
        initiated this signature.
    :type change_request_id: string
    """
    get_signing_service().signed_individual_callback(content, installation_id, github_repository_id, change_request_id)

def post_individual_signed_gerrit(content, user_id):
    """
    Handle the posted callback from the signing service after ICLA signature for Gerrit.

    :param content: The POST body from the signing service callback.
    :type content: string
    :param user_id: The ID of the user that signed. 
    :type user_id: string
    """
    get_signing_service().signed_individual_callback_gerrit(content, user_id)

def post_corporate_signed(content, project_id, company_id):
    """
    Handle the posted callback from the signing service after CCLA signature.

    :param content: The POST body from the signing service callback.
    :type content: string
    :param project_id: The ID of the project that was signed.
    :type project_id: string
    :param company_id: The ID of the company that signed.
    :type company_id: string
    """
    get_signing_service().signed_corporate_callback(content, project_id, company_id)

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
        if event == 'ttl_expired' and not signature.get_signature_signed():
            # Need to re-generate a sign_url and try again.
            cla.log.info('DocuSign URL used was expired, re-generating sign_url')
            callback_url = signature.get_signature_callback_url()
            get_signing_service().populate_sign_url(signature, callback_url)
            signature.save()
            raise falcon.HTTPFound(signature.get_signature_sign_url())
    ret_url = signature.get_signature_return_url()
    if ret_url is not None:
        cla.log.info('Signature success - sending user to return_url: %s', ret_url)
        raise falcon.HTTPFound(ret_url)
    cla.log.info('No return_url set for signature - returning success message')
    return {'success': 'Thank you for signing'}
