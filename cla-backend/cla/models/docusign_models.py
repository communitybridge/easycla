# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Easily perform signing workflows using DocuSign signing service with pydocusign.

NOTE: This integration uses DocuSign's Legacy Authentication REST API Integration.
https://developers.docusign.com/esign-rest-api/guides/post-go-live

"""

import io
import os
import urllib.request
import uuid
import xml.etree.ElementTree as ET
from typing import Dict, Any, Optional
from urllib.parse import urlparse

import cla
import pydocusign  # type: ignore
from cla.controllers.lf_group import LFGroup
from cla.models import signing_service_interface, DoesNotExist
from cla.models.dynamo_models import Signature, User, \
    Project, Company, Gerrit, \
    Document
from cla.models.s3_storage import S3Storage
from pydocusign.exceptions import DocuSignException  # type: ignore

api_base_url = os.environ.get('CLA_API_BASE', '')
root_url = os.environ.get('DOCUSIGN_ROOT_URL', '')
username = os.environ.get('DOCUSIGN_USERNAME', '')
password = os.environ.get('DOCUSIGN_PASSWORD', '')
integrator_key = os.environ.get('DOCUSIGN_INTEGRATOR_KEY', '')

lf_group_client_url = os.environ.get('LF_GROUP_CLIENT_URL', '')
lf_group_client_id = os.environ.get('LF_GROUP_CLIENT_ID', '')
lf_group_client_secret = os.environ.get('LF_GROUP_CLIENT_SECRET', '')
lf_group_refresh_token = os.environ.get('LF_GROUP_REFRESH_TOKEN', '')
lf_group = LFGroup(lf_group_client_url, lf_group_client_id, lf_group_client_secret, lf_group_refresh_token)


class ProjectDoesNotExist(Exception):
    pass


class CompanyDoesNotExist(Exception):
    pass


class UserDoesNotExist(Exception):
    pass


class CCLANotFound(Exception):
    pass


class UserNotWhitelisted(Exception):
    pass


class SigningError(Exception):
    def __init__(self, response):
        self.response = response


class DocuSign(signing_service_interface.SigningService):
    """
    CLA signing service backed by DocuSign.
    """
    TAGS = {'envelope_id': '{http://www.docusign.net/API/3.0}EnvelopeID',
            'type': '{http://www.docusign.net/API/3.0}Type',
            'email': '{http://www.docusign.net/API/3.0}Email',
            'user_name': '{http://www.docusign.net/API/3.0}UserName',
            'routing_order': '{http://www.docusign.net/API/3.0}RoutingOrder',
            'sent': '{http://www.docusign.net/API/3.0}Sent',
            'decline_reason': '{http://www.docusign.net/API/3.0}DeclineReason',
            'status': '{http://www.docusign.net/API/3.0}Status',
            'recipient_ip_address': '{http://www.docusign.net/API/3.0}RecipientIPAddress',
            'client_user_id': '{http://www.docusign.net/API/3.0}ClientUserId',
            'custom_fields': '{http://www.docusign.net/API/3.0}CustomFields',
            'tab_statuses': '{http://www.docusign.net/API/3.0}TabStatuses',
            'account_status': '{http://www.docusign.net/API/3.0}AccountStatus',
            'recipient_id': '{http://www.docusign.net/API/3.0}RecipientId',
            'recipient_statuses': '{http://www.docusign.net/API/3.0}RecipientStatuses',
            'recipient_status': '{http://www.docusign.net/API/3.0}RecipientStatus'}

    def __init__(self):
        self.client = None
        self.s3storage = None

    def initialize(self, config):
        self.client = pydocusign.DocuSignClient(root_url=root_url,
                                                username=username,
                                                password=password,
                                                integrator_key=integrator_key)

        try:
            login_data = self.client.login_information()
            login_account = login_data['loginAccounts'][0]
            base_url = login_account['baseUrl']
            account_id = login_account['accountId']
            url = urlparse(base_url)
            parsed_root_url = '{}://{}/restapi/v2'.format(url.scheme, url.netloc)
        except Exception as e:
            cla.log.error('Error logging in to DocuSign: {}'.format(e))
            return {'errors': {'Error initializing DocuSign'}}

        self.client = pydocusign.DocuSignClient(root_url=parsed_root_url,
                                                account_url=base_url,
                                                account_id=account_id,
                                                username=username,
                                                password=password,
                                                integrator_key=integrator_key)
        self.s3storage = S3Storage()
        self.s3storage.initialize(None)

    def request_individual_signature(self, project_id, user_id, return_url=None):
        request_info = 'project: {project_id}, user: {user_id} with return_url: {return_url}'.format(
            project_id=project_id, user_id=user_id, return_url=return_url)
        cla.log.debug('Individual Signature - creating new signature for: {}'.format(request_info))

        # Ensure this is a valid user
        user_id = str(user_id)
        try:
            user = User()
            user.load(user_id)
            cla.log.debug('Individual Signature - loaded user name: {}, '
                          'user email: {}, gh user: {}, gh id: {}'.
                          format(user.get_user_name(), user.get_user_email(), user.get_github_username(),
                                 user.get_user_github_id()))
        except DoesNotExist as err:
            cla.log.warning('Individual Signature - user ID was NOT found for: {}'.format(request_info))
            return {'errors': {'user_id': str(err)}}

        # Ensure the project exists
        try:
            project = Project()
            project.load(project_id)
            cla.log.debug('Individual Signature - loaded project id: {}, name: {}, '.
                          format(project.get_project_id(), project.get_project_name()))
        except DoesNotExist as err:
            cla.log.warning('Individual Signature - project ID NOT found for: {}'.format(request_info))
            return {'errors': {'project_id': str(err)}}

        # Check for active signature object with this project. If the user has
        # signed the most recent major version, they do not need to sign again.
        cla.log.debug('Individual Signature - loading latest user signature for user: {}, project: {}'.
                      format(user, project))
        latest_signature = user.get_latest_signature(str(project_id))
        cla.log.debug('Individual Signature - loaded latest user signature for user: {}, project: {}'.
                      format(user, project))

        cla.log.debug('Individual Signature - loading latest individual document for project: {}'.
                      format(project))
        last_document = project.get_latest_individual_document()
        cla.log.debug('Individual Signature - loaded latest individual document for project: {}'.
                      format(project))

        if latest_signature is not None and \
                last_document.get_document_major_version() == latest_signature.get_signature_document_major_version():
            cla.log.debug('Individual Signature - user already has a signatures with this project: {}'.
                          format(latest_signature.get_signature_id()))
            return {'user_id': user_id,
                    'project_id': project_id,
                    'signature_id': latest_signature.get_signature_id(),
                    'sign_url': latest_signature.get_signature_sign_url()}
        else:
            cla.log.debug('Individual Signature - user does NOT have a signatures with this project: {}'.
                          format(project))

        # Generate signature callback url
        cla.log.debug('Individual Signature - get active signature metadata')
        signature_metadata = cla.utils.get_active_signature_metadata(user_id)
        cla.log.debug('Individual Signature - get active signature metadata: {}'.format(signature_metadata))

        cla.log.debug('Individual Signature - get individual signature callback url')
        callback_url = cla.utils.get_individual_signature_callback_url(user_id, signature_metadata)
        cla.log.debug('Individual Signature - get individual signature callback url: {}'.format(callback_url))

        # Get signature return URL
        if return_url is None:
            return_url = cla.utils.get_active_signature_return_url(user_id, signature_metadata)
            cla.log.debug('Individual Signature - setting signature return_url to {}'.format(return_url))

        if return_url is None:
            cla.log.warning('No active signature found for user - cannot generate '
                            'return_url without knowing where the user came from')
            return {'user_id': str(user_id),
                    'project_id': str(project_id),
                    'signature_id': None,
                    'sign_url': None,
                    'error': 'No active signature found for user - cannot generate return_url without knowing where the user came from'}

        # Get latest document
        try:
            cla.log.debug('Individual Signature - loading project latest individual document...')
            document = project.get_latest_individual_document()
            cla.log.debug('Individual Signature - loaded project latest individual document: {}'.format(document))
        except DoesNotExist as err:
            cla.log.warning('Individual Signature - project individual document does NOT exist for: {}'.
                            format(request_info))
            return {'errors': {'project_id': project_id, 'message': str(err)}}

        # If the CCLA/ICLA template is missing (not created in the project console), we won't have a document
        # return an error
        if not document:
            return {'errors': {'project_id': project_id, 'message': 'missing template document'}}

        cla.log.debug('Individual Signature - creating default individual values for user: {}'.format(user))
        default_cla_values = create_default_individual_values(user)
        cla.log.debug('Individual Signature - created default individual values: {}'.format(default_cla_values))

        # Create new Signature object
        cla.log.debug('Individual Signature - creating new signature document '
                      'project_id: {}, user_id: {}, return_url: {}, callback_url: {}'.
                      format(project_id, user_id, return_url, callback_url))
        signature = Signature(signature_id=str(uuid.uuid4()),
                              signature_project_id=project_id,
                              signature_document_major_version=document.get_document_major_version(),
                              signature_document_minor_version=document.get_document_minor_version(),
                              signature_reference_id=user_id,
                              signature_reference_type='user',
                              signature_type='cla',
                              signature_return_url_type='Github',
                              signature_signed=False,
                              signature_approved=True,
                              signature_return_url=return_url,
                              signature_callback_url=callback_url)

        # Set signature ACL
        cla.log.debug('Individual Signature - setting ACL using user GH id: {}'.format(user.get_user_github_id()))
        signature.set_signature_acl('github:{}'.format(user.get_user_github_id()))

        # Populate sign url
        self.populate_sign_url(signature, callback_url, default_values=default_cla_values)

        # Save signature
        signature.save()
        cla.log.debug('Individual Signature - Saved signature for: {}'.format(request_info))

        response = {'user_id': str(user_id),
                    'project_id': project_id,
                    'signature_id': signature.get_signature_id(),
                    'sign_url': signature.get_signature_sign_url()}

        cla.log.debug('Individual Signature - returning response: {}'.format(response))
        return response

    def request_individual_signature_gerrit(self, project_id, user_id, return_url=None):
        request_info = 'project: {project_id}, user: {user_id} with return_url: {return_url}'.format(
            project_id=project_id, user_id=user_id, return_url=return_url)
        cla.log.info('Creating new Gerrit signature for {}'.format(request_info))

        # Ensure this is a valid user
        user_id = str(user_id)
        try:
            user = User()
            user.load(user_id)
        except DoesNotExist as err:
            cla.log.warning('User ID does NOT found when requesting a signature for: {}'.format(request_info))
            return {'errors': {'user_id': str(err)}}

        # Ensure the project exists
        try:
            project = Project()
            project.load(project_id)
        except DoesNotExist as err:
            cla.log.warning('Project ID does NOT found when requesting a signature for: {}'.format(request_info))
            return {'errors': {'project_id': str(err)}}

        # Check for active signature object with this project. If the user has
        # signed the most recent major version, they do not need to sign again.
        latest_signature = user.get_latest_signature(str(project_id))
        last_document = project.get_latest_individual_document()
        if latest_signature is not None and \
                last_document.get_document_major_version() == latest_signature.get_signature_document_major_version():
            cla.log.info('User already has a signatures with this project: %s', \
                         latest_signature.get_signature_id())
            return {'user_id': user_id,
                    'project_id': project_id,
                    'signature_id': latest_signature.get_signature_id(),
                    'sign_url': latest_signature.get_signature_sign_url()}

        # the github flow has an option to have the return_url as a blank field,
        # and retrieves the return_url from the signature's metadata (github org id, PR id, etc.)
        # It will return the user to the pull request page.
        # For Gerrit users, we want the return_url to be the link to the Gerrit Instance's page.
        # Since Gerrit users will be able to make changes once they are part of the LDAP Group, 
        # They do not need to be directed to a specific code submission on Gerrit. 
        if return_url is None:
            try:
                gerrits = Gerrit().get_gerrit_by_project_id(project_id)
                if len(gerrits) >= 1:
                    # Github sends the user back to the pull request. Gerrit should send it back to the Gerrit instance url. 
                    return_url = gerrits[0].get_gerrit_url()
            except DoesNotExist as err:
                cla.log.error('Gerrit Instance not found by the given project ID: %s',
                              project_id)
                return {'errors': {'project_id': str(err)}}

        try:
            document = project.get_project_individual_document()
        except DoesNotExist as err:
            cla.log.warning('Document does NOT exist when searching for ICLA for: {}'.format(request_info))
            return {'errors': {'project_id': str(err)}}

        callback_url = self._generate_individual_signature_callback_url_gerrit(user_id)

        default_cla_values = create_default_individual_values(user)

        # Create new Signature object
        signature = Signature(signature_id=str(uuid.uuid4()),
                              signature_project_id=project_id,
                              signature_document_major_version=document.get_document_major_version(),
                              signature_document_minor_version=document.get_document_minor_version(),
                              signature_reference_id=user_id,
                              signature_reference_type='user',
                              signature_type='cla',
                              signature_return_url_type='Gerrit',
                              signature_signed=False,
                              signature_approved=True,
                              signature_return_url=return_url,
                              signature_callback_url=callback_url)

        # Set signature ACL
        signature.set_signature_acl(user.get_lf_username())
        cla.log.info('Set the signature ACL for: {}'.format(request_info))

        # Populate sign url
        self.populate_sign_url(signature, callback_url, default_values=default_cla_values)

        # Save signature
        signature.save()
        cla.log.info('Saved the signature for: {}'.format(request_info))

        return {'user_id': str(user_id),
                'project_id': project_id,
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    @staticmethod
    def check_and_prepare_employee_signature(project_id, company_id, user_id) -> dict:

        # Before an employee begins the signing process, ensure that 
        # 1. The given project, company, and user exists 
        # 2. The company signatory has signed the CCLA for their company. 
        # 3. The user is included as part of the whitelist of the CCLA that the company signed. 
        # Returns an error if any of the above is false. 

        request_info = f'project: {project_id}, company: {company_id}, user: {user_id}'
        cla.log.debug(f'Check and prepare employee signature for {request_info}')

        # Ensure the project exists
        project = Project()
        try:
            project.load(str(project_id))
        except DoesNotExist:
            cla.log.warning('Project does NOT exist for: {}'.format(request_info))
            return {'errors': {'project_id': f'Project ({project_id}) does not exist.'}}
        cla.log.debug(f'Project exists for: {request_info}')

        # Ensure the company exists
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist:
            cla.log.warning(f'Company does NOT exist for: {request_info}')
            return {'errors': {'company_id': f'Company ({company_id}) does not exist.'}}
        cla.log.debug(f'Company exists for: {request_info}')

        # Ensure the user exists
        user = User()
        try:
            user.load(str(user_id))
        except DoesNotExist:
            cla.log.warning(f'User does NOT exist for: {request_info}')
            return {'errors': {'user_id': f'User ({user_id}) does not exist.'}}
        cla.log.debug(f'User exists for: {request_info}')

        # Ensure the company actually has a CCLA with this project.
        # ccla_signatures = Signature().get_signatures_by_project(
        #    project_id,
        #    signature_reference_type='company',
        #    signature_reference_id=company.get_company_id()
        # )
        ccla_signatures = Signature().get_ccla_signatures_by_company_project(
            company_id=company.get_company_id(),
            project_id=project_id
        )
        if len(ccla_signatures) < 1:
            cla.log.warning(f'Company does not have CCLA for: {request_info}')
            return {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}

        cla.log.debug(f'Company has {len(ccla_signatures)} CCLAs for: {request_info}')

        # TODO - DAD: why only grab the first one???
        ccla_signature = ccla_signatures[0]

        # Ensure user is whitelisted for this company.
        if not user.is_whitelisted(ccla_signature):
            # TODO: DAD - update this warning message
            cla.log.warning('No user email whitelisted for this CCLA: {}'.format(request_info))
            return {'errors': {'ccla_whitelist': 'user not whitelisted for this ccla'}}

        cla.log.info(f'User is whitelisted for this CCLA: {request_info}')

        # Assume this company is the user's employer.
        # TODO: DAD - we should check to see if they already have a company id assigned
        user.set_user_company_id(str(company_id))

        # Take a moment to update the user record's github information
        github_username = user.get_user_github_username()
        github_id = user.get_user_github_id()

        if github_username is None and github_id is not None:
            github_username = cla.utils.lookup_user_github_username(github_id)
            if github_username is not None:
                cla.log.debug(f'Updating user record - adding github username: {github_username}')
                user.set_user_github_username(github_username)

        # Attempt to fetch the github id based on the github username
        if github_id is None and github_username is not None:
            github_username = github_username.strip()
            github_id = cla.utils.lookup_user_github_id(github_username)
            if github_id is not None:
                cla.log.debug(f'Updating user record - adding github id: {github_id}')
                user.set_user_github_id(github_id)

        user.save()
        cla.log.info(f'Assigned company ID to user. Employee is ready to sign the CCLA: {request_info}')

        return {'success': {'the employee is ready to sign the CCLA'}}

    def request_employee_signature(self, project_id, company_id, user_id, return_url=None):

        request_info = 'project: {project_id}, company: {company_id}, user: {user_id} with return_url: {return_url}'.format(
            project_id=project_id, company_id=company_id, user_id=user_id, return_url=return_url)
        cla.log.info('Processing request_employee_signature request with {}'.format(request_info))

        check_and_prepare_signature = self.check_and_prepare_employee_signature(project_id, company_id, user_id)
        # Check if there are any errors while preparing the signature. 
        if 'errors' in check_and_prepare_signature:
            cla.log.warning('Error in check_and_prepare_signature with: {} - signatures: {}'.format(request_info,
                                                                                                    check_and_prepare_signature))
            return check_and_prepare_signature

        # Ensure user hasn't already signed this signature.
        employee_signatures = Signature().get_signatures_by_project(
            project_id,
            signature_reference_type='user',
            signature_reference_id=user_id,
            signature_user_ccla_company_id=company_id
        )
        # Return existing signature if employee already signed it. 
        if len(employee_signatures) > 0:
            cla.log.info('Employee has signed for: {} - signatures: {}'.format(request_info, employee_signatures))
            return employee_signatures[0].to_dict()

        cla.log.info('Employee has NOT signed it for: {}'.format(request_info))

        # Requires us to know where the user came from.
        signature_metadata = cla.utils.get_active_signature_metadata(user_id)
        if return_url is None:
            cla.log.info('No return URL for: {}'.format(request_info))
            return_url = cla.utils.get_active_signature_return_url(user_id, signature_metadata)
            cla.log.info('Set return URL for: {} to: {}'.format(request_info, return_url))

        # project has already been checked from check_and_prepare_employee_signature. Load project with project ID.
        project = Project()
        project.load(project_id)
        cla.log.info('Loaded project details for: {}'.format(request_info))

        # Get project's latest corporate document to get major/minor version numbers.
        last_document = project.get_latest_corporate_document()
        cla.log.info('Loaded last project document details for: {}'.format(request_info))

        # return_url may still be empty at this point - the console will deal with it
        new_signature = Signature(signature_id=str(uuid.uuid4()),
                                  signature_project_id=project_id,
                                  signature_document_minor_version=last_document.get_document_minor_version(),
                                  signature_document_major_version=last_document.get_document_major_version(),
                                  signature_reference_id=user_id,
                                  signature_reference_type='user',
                                  signature_type='cla',
                                  signature_signed=True,
                                  signature_approved=True,
                                  signature_return_url=return_url,
                                  signature_user_ccla_company_id=company_id)
        cla.log.info('Created new signature document for: {} - signature: {}'.format(request_info, new_signature))

        # Set signature ACL (user already validated in 'check_and_prepare_employee_signature')
        user = User()
        user.load(str(user_id))
        new_signature.set_signature_acl('github:{}'.format(user.get_user_github_id()))

        # Save signature
        new_signature.save()
        cla.log.info('Set and saved signature for: {}'.format(request_info))

        # If the project does not require an ICLA to be signed, update the pull request and remove the active
        # signature metadata.
        if not project.get_project_ccla_requires_icla_signature():
            cla.log.info('Project does not requires ICLA signature from employee - updating PR')
            github_repository_id = signature_metadata['repository_id']
            change_request_id = signature_metadata['pull_request_id']

            # Get repository
            installation_id = cla.utils.get_installation_id_from_github_repository(github_repository_id)
            if installation_id is None:
                return {'errors': {'github_repository_id': 'The given github repository ID does not exist. '}}

            update_repository_provider(installation_id, github_repository_id, change_request_id)

            cla.utils.delete_active_signature_metadata(user_id)
        else:
            cla.log.info('Project requires ICLA signature from employee - PR has been left unchanged')

        cla.log.info('Returning new signature for: {} - signature: {}'.format(request_info, new_signature))
        return new_signature.to_dict()

    def request_employee_signature_gerrit(self, project_id, company_id, user_id, return_url=None):

        request_info = 'project: {project_id}, company: {company_id}, user: {user_id} with return_url: {return_url}'.format(
            project_id=project_id, company_id=company_id, user_id=user_id, return_url=return_url)
        cla.log.info('Processing request_employee_signature_gerrit request with {}'.format(request_info))

        check_and_prepare_signature = self.check_and_prepare_employee_signature(project_id, company_id, user_id)
        # Check if there are any errors while preparing the signature. 
        if 'errors' in check_and_prepare_signature:
            cla.log.warning('Error in request_employee_signature_gerrit with: {} - signatures: {}'.format(request_info,
                                                                                                          check_and_prepare_signature))
            return check_and_prepare_signature

        # Ensure user hasn't already signed this signature.
        employee_signatures = Signature().get_signatures_by_project(
            project_id,
            signature_reference_type='user',
            signature_reference_id=user_id,
            signature_user_ccla_company_id=company_id
        )
        # Return existing signature if employee already signed it. 
        if len(employee_signatures) > 0:
            cla.log.info('Employee has signed for: {} - signatures: {}'.format(request_info, employee_signatures))
            return employee_signatures[0].to_dict()

        cla.log.info('Employee has NOT signed it for: {}'.format(request_info))

        # Retrieve Gerrits by Project reference ID
        try:
            gerrits = Gerrit().get_gerrit_by_project_id(project_id)
        except DoesNotExist as err:
            cla.log.error('Cannot load Gerrit instance for: %s', request_info)
            return {'errors': {'missing_gerrit': str(err)}}

        # project has already been checked from check_and_prepare_employee_signature. Load project with project ID.
        project = Project()
        project.load(project_id)
        cla.log.info('Loaded project for: %s', request_info)

        # Get project's latest corporate document to get major/minor version numbers. 
        last_document = project.get_latest_corporate_document()

        new_signature = Signature(signature_id=str(uuid.uuid4()),
                                  signature_project_id=project_id,
                                  signature_document_minor_version=last_document.get_document_minor_version(),
                                  signature_document_major_version=last_document.get_document_major_version(),
                                  signature_reference_id=user_id,
                                  signature_reference_type='user',
                                  signature_type='cla',
                                  signature_signed=True,
                                  signature_approved=True,
                                  signature_return_url=return_url,
                                  signature_user_ccla_company_id=company_id)

        # Set signature ACL (user already validated in 'check_and_prepare_employee_signature')
        user = User()
        user.load(str(user_id))
        new_signature.set_signature_acl(user.get_lf_username())

        # Save signature before adding user to the LDAP Group. 
        new_signature.save()
        cla.log.info('Loaded user, set signature ACL, and saved for: %s', request_info)

        for gerrit in gerrits:
            # For every Gerrit Instance of this project, add the user to the LDAP Group.
            # this way we are able to keep track of signed signatures when user fails to be added to the LDAP GROUP.
            group_id = gerrit.get_group_id_ccla()
            # Add the user to the LDAP Group
            try:
                lf_group.add_user_to_group(group_id, user.get_lf_username())
            except Exception as e:
                cla.log.error('Failed in adding user to the LDAP group.{} - {}'.format(e, request_info))
                return

        return new_signature.to_dict()

    def _generate_individual_signature_callback_url_gerrit(self, user_id):
        """
        Helper function to get a user's active signature callback URL for Gerrit

        """
        return os.path.join(api_base_url, 'v2/signed/gerrit/individual', str(user_id))

    def _get_corporate_signature_callback_url(self, project_id, company_id):
        """
        Helper function to get the callback_url of a CCLA signature.

        :param project_id: The ID of the project this CCLA is for.
        :type project_id: string
        :param company_id: The ID of the company signing the CCLA.
        :type company_id: string
        :return: The callback URL hit by the signing provider once the signature is complete.
        :rtype: string
        """
        return os.path.join(api_base_url, 'v2/signed/corporate', str(project_id), str(company_id))

    def request_corporate_signature(self,
                                    auth_user,
                                    project_id,
                                    company_id,
                                    send_as_email=False,
                                    authority_name=None,
                                    authority_email=None,
                                    return_url_type=None,
                                    return_url=None):

        cla.log.debug('Request corporate signature - '
                      f'project id: {project_id}, '
                      f'company id: {company_id}, '
                      f'authority name: {authority_name}, '
                      f'authority email: {authority_email} '
                      f'send email: {send_as_email}')

        # DAD: Auth user is the currently logged in user
        # DAD: Authority Name and Authority Email are from the web form
        # (this is the signatory_user in this function)

        # Ensure the user exists in our database - load the record
        signatory_users = User().get_user_by_username(auth_user.username)
        if signatory_users is None:
            cla.log.warning('unable to lookup auth_user by username: {}'.format(auth_user.username))
            return {'errors': {'user_error': 'user does not exist'}}
        if len(signatory_users) > 1:
            cla.log.warning(f'More than one user record was returned ({len(signatory_users)}) from user '
                            f'username: {auth_user.username} query')
        # Just grab the first record
        signatory_user = signatory_users[0]

        # Ensure the project exists
        project = Project()
        try:
            project.load(str(project_id))
        except DoesNotExist as err:
            cla.log.warning('unable to load project by id: {}'.format(project_id))
            return {'errors': {'project_id': str(err)}}

        # Ensure the company exists
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist as err:
            cla.log.warning('unable to load company by id: {}'.format(company_id))
            return {'errors': {'company_id': str(err)}}

        # Ensure the managers list is not empty - returns a list of User objects
        # TODO: Should we flag the first one as the Initial CLA manager?
        managers = company.get_managers()
        if len(managers) == 0:
            cla.log.warning('No managers for company: {} - ACL is empty'.format(cla.utils.fmt_company(company)))
            return {'errors': {'company_acl': 'Company ACL is empty'}}

        cla.log.debug('Loaded {} managers for company: {}. Will use {} as the initial CLA manager'.
                      format(len(managers), company, cla.utils.fmt_user(managers[0])))

        found_authority = False
        if authority_name is not None:
            for manager in managers:
                if manager.get_user_name() is None:
                    cla.log.warning(f'Manager name {manager} is missing the user name value. '
                                    f'Unable to compare with provided authority_name={authority_name}. '
                                    'Skipping this manager.')
                    continue
                if manager.get_user_name().lower() == authority_name.lower():
                    found_authority = True
                    cla.log.debug(f'Authority name {authority_name} provided is in the '
                                  f'{company.get_company_name()} manager list.')
                    break

        # Get CLA Managers. In the future, we will support contributors
        scheduleA = generate_manager_and_contributor_list(
            [(manager.get_user_name(), manager.get_user_email()) for manager in managers]
        )

        # We may see the provided authority (from the web form) in the manager list...
        # If they are not, then we will send the request via e-mail below
        if not found_authority:
            cla.log.debug(f'Authority name: {authority_name} / {authority_email} '
                          f'NOT found in {company.get_company_name()} '
                          f'manager list: {cla.utils.fmt_users(managers)}')

        # TODO: DAD should we create the values based on:
        # web user? signatory_user / signatory_email
        # values from web form? authority_name / authority_email
        # Maybe we should decide this based on if the web user was already in
        # the manager list? the found_authority flag or do we key off the email flag?
        if authority_name:
            cla.log.debug('Using authority name from web form to create the default company values')
            cla_template_values = create_default_company_values(company,
                                                                authority_name,
                                                                authority_email,
                                                                managers[0].get_user_name(),
                                                                managers[0].get_user_email(),
                                                                scheduleA)
        else:
            cla.log.debug('Using current_user/signatory_user to create the default company values.')
            cla_template_values = create_default_company_values(company,
                                                                signatory_user.get_user_name(),
                                                                signatory_user.get_user_email(),
                                                                managers[0].get_user_name(),
                                                                managers[0].get_user_email(),
                                                                scheduleA)

        # Ensure the contract group has a CCLA
        last_document = project.get_latest_corporate_document()
        if last_document is None or \
                last_document.get_document_major_version() is None or \
                last_document.get_document_minor_version() is None:
            cla.log.info('Contract Group does not have a CCLA: {}'.format(project_id))
            return {'errors': {'project_id': 'Contract Group does not support CCLAs.'}}

        # Ensure the company doesn't already have a CCLA with this project. 
        # and the user is about to sign the ccla manually 
        cla.log.info('Checking if a signature exists for project: {}'.format(project_id))
        latest_signature = company.get_latest_signature(str(project_id))
        if latest_signature is not None and \
                last_document.get_document_major_version() == latest_signature.get_signature_document_major_version():
            cla.log.info('CCLA signature object ({}) already exists '
                         'for company {} on project {}'.
                         format(latest_signature.get_signature_id(), cla.utils.fmt_company(company),
                                cla.utils.fmt_project(project)))

            if latest_signature.get_signature_signed():
                cla.log.info('CCLA signature object ({}) is already signed'.format(latest_signature.get_signature_id()))
                return {'errors': {'signature_id': 'Company has already signed CCLA with this project'}}
            else:
                cla.log.info('CCLA signature object ({}) is NOT signed'.format(latest_signature.get_signature_id()))
                callback_url = self._get_corporate_signature_callback_url(str(project_id), str(company_id))
                cla.log.info('Signing callback URL is: {}'.format(callback_url))

                # Populate sign url
                # TOOD DAD: Need to provide: authority_name, authority_email ??
                if authority_name:
                    cla.log.debug('Using authority name from web form to populate the signing request.')
                    self.populate_sign_url(latest_signature, callback_url,
                                           authority_name,
                                           authority_email,
                                           send_as_email,
                                           managers[0].get_user_name(),
                                           managers[0].get_user_email(),
                                           cla_template_values)
                else:
                    cla.log.debug('Using current user/signatory user to populate the signing request.')
                    if signatory_user.get_user_name() is not None:
                        signatory_name = signatory_user.get_user_name()
                    else:
                        signatory_name = signatory_user.get_lf_username()
                    self.populate_sign_url(latest_signature, callback_url,
                                           signatory_name,
                                           signatory_user.get_user_email(),
                                           send_as_email,
                                           managers[0].get_user_name(),
                                           managers[0].get_user_email(),
                                           cla_template_values)

                return {'company_id': str(company_id),
                        'project_id': str(project_id),
                        'signature_id': latest_signature.get_signature_id(),
                        'sign_url': latest_signature.get_signature_sign_url()}

        # No signature exists, create the new Signature.
        cla.log.info('Creating new signature for company %s on project %s', company_id, project_id)
        signature = Signature(signature_id=str(uuid.uuid4()),
                              signature_project_id=project_id,
                              signature_document_minor_version=last_document.get_document_minor_version(),
                              signature_document_major_version=last_document.get_document_major_version(),
                              signature_reference_id=company_id,
                              signature_reference_type='company',
                              signature_type='ccla',
                              signature_signed=False,
                              signature_approved=True)

        callback_url = self._get_corporate_signature_callback_url(str(project_id), str(company_id))
        cla.log.info('Setting callback_url: %s', callback_url)
        signature.set_signature_callback_url(callback_url)

        if not send_as_email:  # get return url only for manual signing through console
            cla.log.info('Setting signature return_url to %s', return_url)
            signature.set_signature_return_url(return_url)

        # Set signature ACL
        signature.set_signature_acl(signatory_user.get_lf_username())

        # Populate sign url
        if authority_name:
            cla.log.debug('Using authority name from web form to populate the signing request.')
            self.populate_sign_url(signature, callback_url,
                                   authority_name,
                                   authority_email,
                                   send_as_email,
                                   managers[0].get_user_name(),  # authority_name
                                   managers[0].get_user_email(),  # authority_email
                                   cla_template_values)
        else:
            cla.log.debug('Using current user/signatory user to populate the signing request.')
            self.populate_sign_url(signature, callback_url,
                                   signatory_user.get_user_name(),
                                   signatory_user.get_user_email(),
                                   send_as_email,
                                   managers[0].get_user_name(),  # authority_name
                                   managers[0].get_user_email(),  # authority_email
                                   cla_template_values)

        # Save signature
        signature.save()

        return {'company_id': str(company_id),
                'project_id': str(project_id),
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    def populate_sign_url(self, signature, callback_url=None,
                          authority_or_signatory_name=None,
                          authority_or_signatory_email=None,
                          send_as_email=False,
                          cla_manager_name=None, cla_manager_email=None,
                          default_values: Optional[Dict[str, Any]] = None):  # pylint: disable=too-many-locals

        sig_type = signature.get_signature_reference_type()

        cla.log.debug(f'populate_sign_url - Populating sign_url for signature {signature.get_signature_id()} '
                      f'using callback: {callback_url} '
                      f'with authority_or_signatory_name {authority_or_signatory_name} '
                      f'with authority_or_signatory_email {authority_or_signatory_email} '
                      f'with cla manager name: {cla_manager_name} '
                      f'with cla manager email: {cla_manager_email} '
                      f'send as email: {send_as_email} '
                      f' reference type: {sig_type}')

        company = Company()
        user = User()

        # We use user name/email non-email docusign user ICLA
        user_signature_name = 'Unknown'
        user_signature_email = 'Unknown'

        # Assume the company manager is signing the CCLA
        cla.log.debug('populate_sign_url - processing {} signature...'.format(sig_type))

        if sig_type == 'company':
            try:
                # Grab the company id from the signature
                cla.log.debug(f'populate_sign_url - Loading company id: {signature.get_signature_reference_id()}')
                company.load(signature.get_signature_reference_id())
                cla.log.debug(f'populate_sign_url - Loaded company: {company}')

                # Grab the company manager id (id of the user)
                cla.log.debug(f'populate_sign_url - Loading user manager id: {company.get_company_manager_id()}')
                user.load(company.get_company_manager_id())
                cla.log.debug(f'populate_sign_url - Loaded user: {user}')

                user_signature_name = user.get_user_name()
                user_signature_email = user.get_lf_email()
            except DoesNotExist:
                cla.log.warning('populate_sign_url - No CLA manager associated with this company - can not sign CCLA')
                return
            except Exception as e:
                cla.log.warning('populate_sign_url - No CLA manager lookup error: '.format(e))
                return
        elif sig_type == 'user':
            if not send_as_email:
                try:
                    cla.log.debug('populate_sign_url - Loading user by id: {}'.
                                  format(signature.get_signature_reference_id()))
                    user.load(signature.get_signature_reference_id())
                    cla.log.debug('populate_sign_url - Loaded user by id: {} - name: {}, email: {}'.
                                  format(user.get_user_id(), user.get_user_name(), user.get_user_email()))
                    if not user.get_user_name() is None:
                        user_signature_name = user.get_user_name()
                    if not user.get_user_email() is None:
                        user_signature_email = user.get_user_email()
                except DoesNotExist:
                    cla.log.warning('populate_sign_url - no user associated with this signature '
                                    'id: {} - can not sign CCLA'.
                                    format(signature.get_signature_reference_id()))
                    return
                except Exception as e:
                    cla.log.warning('populate_sign_url - no user associated with this signature id: {}, error: {}'.
                                    format(signature.get_signature_reference_id(), e))
                    return
        else:
            cla.log.warning('populate_sign_url - unsupported signature type: {}'.format(sig_type))
            return

        # Fetch the document to sign.
        cla.log.debug('populate_sign_url - Loading project by id: {}'.format(signature.get_signature_project_id()))
        project = Project()
        project.load(signature.get_signature_project_id())
        cla.log.debug('populate_sign_url - Loaded project by id: {} - project: {}'.
                      format(signature.get_signature_project_id(), project))

        # Load the appropriate document
        if sig_type == 'company':
            document = project.get_project_corporate_document()
            if document is None:
                cla.log.error('populate_sign_url - could not get sign url for project: {}. '
                              'Project has no corporate CLA document set'.format(project))
                return
        else:  # sig_type == 'user'
            document = project.get_project_individual_document()
            if document is None:
                cla.log.error('populate_sign_url - Could not get sign url for project {}. '
                              'Project has no individual CLA document set', project)
                return

        # Void the existing envelope to prevent multiple envelopes pending for a signer. 
        envelope_id = signature.get_signature_envelope_id()
        if envelope_id is not None:
            try:
                message = 'You are getting this message because your DocuSign Session ' \
                          'for project {} expired. A new session will be in place for ' \
                          'your signing process.'.format(project.get_project_name())
                cla.log.debug(message)
                self.client.void_envelope(envelope_id, message)
            except Exception as e:
                cla.log.warning('populate_sign_url - DocuSign error while voiding the envelope - '
                                'regardless, continuing on..., error: {}'.format(e))

        # Not sure what should be put in as documentId.
        document_id = uuid.uuid4().int & (1 << 16) - 1  # Random 16bit integer -.pylint: disable=no-member
        tabs = get_docusign_tabs_from_document(document, document_id, default_values=default_values)

        if send_as_email:
            # Sending email to authority
            signatory_email = authority_or_signatory_email
            signatory_name = authority_or_signatory_name

            # Not assigning a clientUserId sends an email.
            project_name = project.get_project_name()
            company_name = company.get_company_name()

            cla.log.debug('populate_sign_url - Sending document as email with name: {}, email: {} '
                          'project name: {}, company: {}'.
                          format(signatory_name, signatory_email, project_name, company_name))

            email_subject = 'CLA Sign Request for {}'.format(project_name)
            email_body = '''{cla_manager_name} has designated you as being an authorized signatory for {company_name}. In order for employees of your company to contribute to the open source project {project_name}, they must do so under a Contributor License Agreement signed by someone with authority to sign on behalf of your company.

After you sign, {cla_manager_name} (as the initial CLA Manager for your company) will be able to maintain the list of specific employees authorized to contribute to the project under this signed CLA.

If you have questions, or if you are not an authorized signatory of this company, please contact the requester at {cla_manager_email}.

            '''.format(cla_manager_name=cla_manager_name,
                       company_name=company_name,
                       project_name=project_name,
                       cla_manager_email=cla_manager_email)

            cla.log.debug('populate_sign_url - Generating a docusign signer object form email with name: {}, email: {}'.
                          format(signatory_name, signatory_email))
            signer = pydocusign.Signer(email=signatory_email,
                                       name=signatory_name,
                                       recipientId=1,
                                       tabs=tabs,
                                       emailSubject=email_subject,
                                       emailBody=email_body,
                                       supportedLanguage='en',
                                       )
        else:
            # This will be the Initial CLA Manager
            signatory_name = user_signature_name
            signatory_email = user_signature_email

            # Assigning a clientUserId does not send an email.
            # It assumes that the user handles the communication with the client. 
            # In this case, the user opened the docusign document to manually sign it. 
            # Thus the email does not need to be sent.
            cla.log.debug('populate_sign_url - Generating a docusign signer object with name: {}, email: {}'.
                          format(signatory_name, signatory_email))
            signer = pydocusign.Signer(email=signatory_email, name=signatory_name,
                                       recipientId=1, clientUserId=signature.get_signature_id(),
                                       tabs=tabs,
                                       emailSubject='CLA Sign Request',
                                       emailBody='CLA Sign Request for {}'.format(user.get_user_email()),
                                       supportedLanguage='en',
                                       )

        content_type = document.get_document_content_type()
        if document.get_document_s3_url() is not None:
            pdf = self.get_document_resource(document.get_document_s3_url())
        elif content_type.startswith('url+'):
            pdf_url = document.get_document_content()
            pdf = self.get_document_resource(pdf_url)
        else:
            content = document.get_document_content()
            pdf = io.BytesIO(content)

        doc_name = document.get_document_name()
        cla.log.debug('populate_sign_url - Docusign document name: {}, id: {}, content type: {}'.
                      format(doc_name, document_id, content_type))
        document = pydocusign.Document(name=doc_name, documentId=document_id, data=pdf)

        if callback_url is not None:
            # Webhook properties for callbacks after the user signs the document.
            # Ensure that a webhook is returned on the status "Completed" where 
            # all signers on a document finish signing the document. 
            recipient_events = [{"recipientEventStatusCode": "Completed"}]
            event_notification = pydocusign.EventNotification(url=callback_url,
                                                              loggingEnabled=True,
                                                              recipientEvents=recipient_events)
            envelope = pydocusign.Envelope(documents=[document],
                                           emailSubject='CLA Sign Request',
                                           emailBlurb='CLA Sign Request',
                                           eventNotification=event_notification,
                                           status=pydocusign.Envelope.STATUS_SENT,
                                           recipients=[signer])
        else:
            envelope = pydocusign.Envelope(documents=[document],
                                           emailSubject='CLA Sign Request',
                                           emailBlurb='CLA Sign Request',
                                           status=pydocusign.Envelope.STATUS_SENT,
                                           recipients=[signer])

        envelope = self.prepare_sign_request(envelope)

        if not send_as_email:
            recipient = envelope.recipients[0]

            # The URL the user will be redirected to after signing.
            # This route will be in charge of extracting the signature's return_url and redirecting.
            return_url = os.path.join(api_base_url, 'v2/return-url', str(recipient.clientUserId))

            cla.log.debug("populate_sign_url - generating signature sign_url, using return-url as: {}".
                          format(return_url))
            sign_url = self.get_sign_url(envelope, recipient, return_url)
            cla.log.debug('populate_sign_url - setting signature sign_url as: {}'.format(sign_url))
            signature.set_signature_sign_url(sign_url)

        # Save Envelope ID in signature.
        cla.log.debug('populate_sign_url - saving signature to database...')
        signature.set_signature_envelope_id(envelope.envelopeId)
        signature.save()

    def signed_individual_callback(self, content, installation_id, github_repository_id, change_request_id):
        """
        Will be called on ICLA signature callback, but also when a document has been
        opened by a user - no action required then.
        """
        cla.log.debug('Docusign ICLA signed callback POST data: {}'.format(content))
        tree = ET.fromstring(content)
        # Get envelope ID.
        envelope_id = tree.find('.//' + self.TAGS['envelope_id']).text
        # Assume only one signature per signature.
        signature_id = tree.find('.//' + self.TAGS['client_user_id']).text
        signature = cla.utils.get_signature_instance()
        try:
            signature.load(signature_id)
        except DoesNotExist:
            cla.log.error('DocuSign ICLA callback returned signed info on invalid signature: %s',
                          content)
            return
        # Iterate through recipients and update the signature signature status if changed.
        elem = tree.find('.//' + self.TAGS['recipient_statuses'] +
                         '/' + self.TAGS['recipient_status'])
        status = elem.find(self.TAGS['status']).text
        if status == 'Completed' and not signature.get_signature_signed():
            cla.log.info('ICLA signature signed (%s) - Notifying repository service provider',
                         signature_id)
            signature.set_signature_signed(True)
            # Save signature
            signature.save()
            # Send user their signed document.
            user = User()
            user.load(signature.get_signature_reference_id())
            # Remove the active signature metadata.
            cla.utils.delete_active_signature_metadata(user.get_user_id())
            # Get signed document
            document_data = self.get_signed_document(envelope_id, user)
            # Send email with signed document.
            self.send_signed_document(document_data, user)

            # Verify user id exist for saving on storage
            user_id = user.get_user_id()
            if user_id is None:
                raise SigningError('Missing user_id on ICLA for saving signed file on s3 storage.')

            # Store document on S3
            project_id = signature.get_signature_project_id()
            self.send_to_s3(document_data, project_id, signature_id, 'icla', user_id)

            # Update the repository provider with this change.
            update_repository_provider(installation_id, github_repository_id, change_request_id)

    def signed_individual_callback_gerrit(self, content, user_id):
        cla.log.debug('Docusign Gerrit ICLA signed callback POST data: %s', content)
        tree = ET.fromstring(content)
        # Get envelope ID.
        envelope_id = tree.find('.//' + self.TAGS['envelope_id']).text
        # Assume only one signature per signature.
        signature_id = tree.find('.//' + self.TAGS['client_user_id']).text
        signature = cla.utils.get_signature_instance()
        try:
            signature.load(signature_id)
        except DoesNotExist:
            cla.log.error('DocuSign Gerrit ICLA callback returned signed info on invalid signature: %s',
                          content)
            return
        # Iterate through recipients and update the signature signature status if changed.
        elem = tree.find('.//' + self.TAGS['recipient_statuses'] +
                         '/' + self.TAGS['recipient_status'])
        status = elem.find(self.TAGS['status']).text
        if status == 'Completed' and not signature.get_signature_signed():
            cla.log.info('ICLA signature signed (%s) - Notifying repository service provider',
                         signature_id)
            # Get User
            user = cla.utils.get_user_instance()
            user.load(user_id)

            # Save signature before adding user to LDAP Groups.
            signature.set_signature_signed(True)
            # Save signature
            signature.save()

            gerrits = Gerrit().get_gerrit_by_project_id(signature.get_signature_project_id())
            for gerrit in gerrits:
                # Get Gerrit Group ID
                group_id = gerrit.get_group_id_icla()

                # Check if Group id is none
                if group_id is not None:
                    lf_username = user.get_lf_username()
                    # Add the user to the LDAP Group
                    try:
                        lf_group.add_user_to_group(group_id, lf_username)
                    except Exception as e:
                        cla.log.error('Failed in adding user to the LDAP group: %s', e)
                        return

            # Get signed document
            document_data = self.get_signed_document(envelope_id, user)
            # Send email with signed document.
            self.send_signed_document(document_data, user)

            # Verify user id exist for saving on storage
            if user_id is None:
                raise SigningError('Missing user_id on ICLA for saving signed file on s3 storage.')

            # Store document on S3
            project_id = signature.get_signature_project_id()
            self.send_to_s3(document_data, project_id, signature_id, 'icla', user_id)

    def signed_corporate_callback(self, content, project_id, company_id):
        """
        Will be called on CCLA signature callback, but also when a document has been
        opened by a user - no action required then.
        """
        cla.log.debug('Docusign CCLA signed callback POST data: %s', content)
        tree = ET.fromstring(content)
        # Get envelope ID.
        envelope_id = tree.find('.//' + self.TAGS['envelope_id']).text

        # Get Company with company ID. 
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist as err:
            return {'errors': {'Docusign callback failed: Invalid company_id {}'.format(company_id): str(err)}}

        # Assume only one signature per signature.
        client_user_id = tree.find('.//' + self.TAGS['client_user_id'])
        if client_user_id is not None:
            signature_id = client_user_id.text
            signature = cla.utils.get_signature_instance()
            try:
                signature.load(signature_id)
            except DoesNotExist:
                cla.log.error('DocuSign callback returned signed info on invalid signature: %s',
                              content)
                return
        else:
            # If client_user_id is None, the callback came from the email that finished signing. 
            # Retrieve the latest signature with projectId and CompanyId.
            signature = company.get_latest_signature(str(project_id))
            signature_id = signature.get_signature_id()

        # Get User
        user = cla.utils.get_user_instance()
        if signature.get_signature_reference_type() == 'user':
            user.load(signature.get_signature_reference_id())
        elif signature.get_signature_reference_type() == 'company':
            # Get company manager if reference id is of a company's ID. 
            user.load(company.get_company_manager_id())

        # Iterate through recipients and update the signature signature status if changed.
        elem = tree.find('.//' + self.TAGS['recipient_statuses'] +
                         '/' + self.TAGS['recipient_status'])
        status = elem.find(self.TAGS['status']).text
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist:
            cla.log.error('Received CCLA signed callback from signing service provider for an unknown company: %s',
                          company_id)
            return
        if status == 'Completed' and not signature.get_signature_signed():
            cla.log.info('CCLA signature signed (%s)', signature_id)
            signature.set_signature_signed(True)
            # Save signature
            signature.save()

            # Check if the callback is for a Gerrit Instance
            try:
                gerrits = Gerrit().get_gerrit_by_project_id(signature.get_signature_project_id())
            except DoesNotExist:
                gerrits = []

            # Get LF user name. 
            lf_username = user.get_lf_username()
            for gerrit in gerrits:
                # Get Gerrit Group ID
                group_id = gerrit.get_group_id_ccla()

                # Check if Group id is none
                if group_id is not None:
                    # Add the user to the LDAP Group (corporate authority)
                    try:
                        lf_group.add_user_to_group(group_id, lf_username)
                    except Exception as e:
                        cla.log.error('Failed in adding user to the LDAP group: %s', e)
                        return

            # Send manager their signed document.
            manager = User()
            manager.load(company.get_company_manager_id())
            # Get signed document
            document_data = self.get_signed_document(envelope_id, user)
            # Send email with signed document.
            self.send_signed_document(document_data, user)

            # verify company_id is not none
            if company_id is None:
                raise SigningError('Missing company_id on CCLA for saving signed file on s3 storage.')

            # Store document on S3
            self.send_to_s3(document_data, project_id, signature_id, 'ccla', company_id)

    def get_signed_document(self, envelope_id, user):
        """Helper method to get the signed document from DocuSign."""

        cla.log.debug('Fetching signed CLA document for envelope: %s', envelope_id)
        envelope = pydocusign.Envelope()
        envelope.envelopeId = envelope_id
        try:
            documents = envelope.get_document_list(self.client)
        except Exception as err:
            cla.log.error('Unknown error when trying to load signed document: %s', str(err))
            return
        if documents is None or len(documents) < 1:
            cla.log.error('Could not find signed document envelope %s and user %s',
                          envelope_id, user.get_user_email())
            return
        document = documents[0]
        if 'documentId' not in document:
            cla.log.error('Not document ID found in document response: %s', str(document))
            return
        try:
            # TODO: Also send the signature certificate? envelope.get_certificate()
            document_file = envelope.get_document(document['documentId'], self.client)
        except Exception as err:
            cla.log.error('Unknown error when trying to fetch signed document content ' + \
                          'for document ID %s: %s', document['documentId'], str(err))
            return
        return document_file.read()

    def send_signed_document(self, document_data, user, icla=True):
        """Helper method to send the user their signed document."""

        subject = 'CLA Signed Document'
        body = 'Thank you for signing the CLA! Your signed document is attached to this email.'

        recipient = user.get_user_email()
        filename = recipient + '-cla.pdf'
        attachment = {'type': 'content',
                      'content': document_data,
                      'content-type': 'application/pdf',
                      'filename': filename}
        # Third, send the email.
        cla.log.info('Sending signed CLA document to %s', recipient)
        cla.utils.get_email_service().send(subject, body, recipient, attachment)

    def send_to_s3(self, document_data, project_id, signature_id, cla_type, identifier):
        # cla_type could be: icla or ccla (String)
        # identifier could be: user_id or company_id
        filename = str.join('/',
                            ('contract-group', str(project_id), cla_type, str(identifier), str(signature_id) + '.pdf'))
        self.s3storage.store(filename, document_data)

    def get_document_resource(self, url):  # pylint: disable=no-self-use
        """
        Mockable method to fetch the PDF for signing.

        :param url: The URL of the PDF file to sign.
        :type url: string
        :return: A resource that can be read()'d.
        :rtype: Resource
        """
        return urllib.request.urlopen(url)

    def prepare_sign_request(self, envelope):
        """
        Mockable method for sending a signature request to DocuSign.

        :param envelope: The envelope to send to DocuSign.
        :type envelope: pydocusign.Envelope
        :return: The new envelope to work with after the request has been sent.
        :rtype: pydocusign.Envelope
        """
        try:
            self.client.create_envelope_from_documents(envelope)
            envelope.get_recipients()
            return envelope
        except DocuSignException as err:
            cla.log.error('Error while fetching DocuSign envelope recipients: %s', str(err))

    def get_sign_url(self, envelope, recipient, return_url):  # pylint:disable=no-self-use
        """
        Mockable method for getting a signing url.

        :param envelope: The envelope in question.
        :type envelope: pydocusign.Envelope
        :param recipient: The recipient inside this envelope.
        :type recipient: pydocusign.Recipient
        :param return_url: The URL to return the user after successful signing.
        :type return_url: string
        :return: A URL for the recipient to hit for signing.
        :rtype: string
        """
        return envelope.post_recipient_view(recipient, returnUrl=return_url)


class MockDocuSign(DocuSign):
    """
    Mock object to test DocuSign service implementation.
    """

    def get_document_resource(self, url):
        """
        Need to implement fake resource here.
        """
        return open(cla.utils.get_cla_path() + '/tests/resources/test.pdf', 'rb')

    def prepare_sign_request(self, envelope):
        """
        Don't actually send the request when running tests.
        """
        recipients = []
        for recipient in envelope.recipients:
            recip = lambda: None
            recip.clientUserId = recipient.clientUserId
            recipients.append(recip)
        envelope = MockRecipient()
        envelope.recipients = recipients
        return envelope

    def get_sign_url(self, envelope, recipient, return_url):
        """
        Don't communicate with DocuSign when running tests.
        """
        return 'http://signing-service.com/send-user-here'

    def send_signed_document(self, envelope_id, user):
        """Mock method to send a signed DocuSign document to the user's email."""
        pass


class MockRecipient(object):
    def __init__(self):
        self.recipients = None
        self.envelopeId = None


def update_repository_provider(installation_id, github_repository_id, change_request_id):
    """Helper method to notify the repository provider of successful signature."""
    repo_service = cla.utils.get_repository_service('github')
    repo_service.update_change_request(installation_id, github_repository_id, change_request_id)


def get_org_from_return_url(repo_provider_type, return_url, orgs):
    """
    Helper method to find specific org from list of orgs under same contract group
    This is a hack solution since it totally depends on return_url and repo service provider
    However, based on the current implementation, it's a simple way to invovled minimal refactor
    BTW, I don't believe the last team can do a successful demo without doing any tweaks like this

    :param repo_provider_type: The repo service provider.
    :type repo_provider_type: string
    :param return_url: The URL will be redirected after signature done.
    :type return_url: string
    :return: List of Organizations of any repo service provider.
    :rtype: [any_repo_service_provider.Organization]
    """
    if repo_provider_type is 'github':
        split_url = return_url.split('/')  # parse repo name from URL
        target_org_name = split_url[3]
        for org in orgs:
            if org.get_organization_name() == target_org_name:
                return org
        raise Exception('Not found org: {} under current CLA project'.format(target_org_name))
    else:
        raise Exception('Repo service: {} not supported'.format(repo_provider_type))


def get_docusign_tabs_from_document(document: Document,
                                    document_id: int,
                                    default_values: Optional[Dict[str, Any]] = None):
    """
    Helper function to extract the DocuSign tabs out of a document object.

    :param document: The document to extract the tabs from.
    :type document: cla.models.model_interfaces.Document
    :param document_id: The ID of the document to use for grouping of the tabs.
    :type document_id: int
    :return: List of formatted tabs for consumption by pydocusign.
    :rtype: [pydocusign.Tab]
    """
    tabs = []
    for tab in document.get_document_tabs():
        args = {
            'documentId': document_id,
            'pageNumber': tab.get_document_tab_page(),
            'xPosition': tab.get_document_tab_position_x(),
            'yPosition': tab.get_document_tab_position_y(),
            'width': tab.get_document_tab_width(),
            'height': tab.get_document_tab_height(),
            'customTabId': tab.get_document_tab_id(),
            'tabLabel': tab.get_document_tab_id(),
            'name': tab.get_document_tab_name()
        }

        if tab.get_document_tab_anchor_string() is not None:
            # Set only when anchor string exists 
            args['anchorString'] = tab.get_document_tab_anchor_string()
            args['anchorIgnoreIfNotPresent'] = tab.get_document_tab_anchor_ignore_if_not_present()
            args['anchorXOffset'] = tab.get_document_tab_anchor_x_offset()
            args['anchorYOffset'] = tab.get_document_tab_anchor_y_offset()
            # Remove x,y coordinates since offsets will define them
        # del args['xPosition']
        # del args['yPosition']

        if default_values is not None and \
                default_values.get(tab.get_document_tab_id()) is not None:
            args['value'] = default_values[tab.get_document_tab_id()]

        tab_type = tab.get_document_tab_type()
        if tab_type == 'text':
            tab_class = pydocusign.TextTab
        elif tab_type == 'text_unlocked':
            tab_class = TextUnlockedTab
            args['locked'] = False
        elif tab_type == 'text_optional':
            tab_class = TextOptionalTab
            args['required'] = False
        elif tab_type == 'number':
            tab_class = pydocusign.NumberTab
        elif tab_type == 'sign':
            tab_class = pydocusign.SignHereTab
        elif tab_type == 'date':
            tab_class = pydocusign.DateSignedTab
        else:
            cla.log.warning('Invalid tab type specified (%s) in document file ID %s',
                            tab_type, document.get_document_file_id())
            continue

        tab_obj = tab_class(**args)
        tabs.append(tab_obj)

    return tabs


# Returns a dictionary of document id to value
def create_default_company_values(company: Company,
                                  signatory_name: str,
                                  signatory_email: str,
                                  manager_name: str,
                                  manager_email: str,
                                  schedule_a: str) -> Dict[str, Any]:
    values = {}

    if company is not None and company.get_company_name() is not None:
        values['corporation_name'] = company.get_company_name()
        values['corporation'] = company.get_company_name()

    if signatory_name is not None:
        values['signatory_name'] = signatory_name

    if signatory_email is not None:
        values['signatory_email'] = signatory_email

    if manager_name is not None:
        values['point_of_contact'] = manager_name
        values['cla_manager_name'] = manager_name

    if manager_email is not None:
        values['email'] = manager_email
        values['cla_manager_email'] = manager_email

    if schedule_a is not None:
        values['scheduleA'] = schedule_a

    return values


def create_default_individual_values(user: User) -> Dict[str, Any]:
    values = {}

    if user is None:
        return values

    if user.get_user_name() is not None:
        values['full_name'] = user.get_user_name()
        values['public_name'] = user.get_user_name()

    if user.get_user_email() is not None:
        values['email'] = user.get_user_email()

    return values


class TextOptionalTab(pydocusign.Tab):
    """Tab to show a free-form text field on the document.
    """
    attributes = pydocusign.Tab._common_attributes + pydocusign.Tab._formatting_attributes + [
        'name',
        'value',
        'height',
        'width',
        'locked',
        'required'
    ]
    tabs_name = 'textTabs'


class TextUnlockedTab(pydocusign.Tab):
    """Tab to show a free-form text field on the document.
    """
    attributes = pydocusign.Tab._common_attributes + pydocusign.Tab._formatting_attributes + [
        'name',
        'value',
        'height',
        'width',
        'locked'
    ]
    tabs_name = 'textTabs'


# managers and contributors are tuples of (name, email)
def generate_manager_and_contributor_list(managers, contributors=None):
    lines = []

    for manager in managers:
        lines.append('CLA Manager: {}, {}'.format(manager[0], manager[1]))

    if contributors is not None:
        for contributor in contributors:
            lines.append('{}, {}'.format(contributor[0], contributor[1]))

    lines = '\n'.join([str(line) for line in lines])

    return lines
