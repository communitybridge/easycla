"""
Easily perform signing workflows using DocuSign signing service with pydocusign.

NOTE: This integration uses DocuSign's Legacy Authentication REST API Integration.
https://developers.docusign.com/esign-rest-api/guides/post-go-live

"""

import io
import os
import uuid
import json
import urllib.request
from urllib.parse import urlparse
import xml.etree.ElementTree as ET
from typing import Dict, Any, Optional, Tuple, List

import pydocusign # type: ignore
from pydocusign.exceptions import DocuSignException # type: ignore

import cla
from cla.controllers.lf_group import LFGroup
from cla.models import signing_service_interface, DoesNotExist
from cla.models.s3_storage import S3Storage
from cla.models.dynamo_models import Signature, GitHubOrg, User, \
                                        Project, Company, Gerrit, \
                                        Document, Repository

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
        cla.log.info('Creating new signature for user %s on project %s', user_id, project_id)

        # Ensure this is a valid user
        user_id = str(user_id)
        try:
            user = User()
            user.load(user_id)
        except DoesNotExist as err:
            cla.log.warning('User ID not found when trying to request a signature: %s',
                            user_id)
            return {'errors': {'user_id': str(err)}}

        # Ensure the project exists
        try:
            project = Project()
            project.load(project_id)
        except DoesNotExist as err:
            cla.log.error('Project ID not found when trying to request a signature: %s',
                        project_id)
            return {'errors': {'project_id': str(err)}}

        # Check for active signature object with this project. If the user has
        # signed the most recent major version, they do not need to sign again.
        latest_signature = user.get_latest_signature(user, str(project_id))
        last_document = project.get_latest_individual_document()
        if latest_signature is not None and \
           last_document.get_document_major_version() == latest_signature.get_signature_document_major_version():
            cla.log.info('User already has a signatures with this project: %s', \
                         latest_signature.get_signature_id())
            return {'user_id': user_id,
                    'project_id': project_id,
                    'signature_id': latest_signature.get_signature_id(),
                    'sign_url': latest_signature.get_signature_sign_url()}

        # Generate signature callback url
        signature_metadata = cla.utils.get_active_signature_metadata(user_id)
        callback_url = cla.utils.get_individual_signature_callback_url(user_id, signature_metadata)
        cla.log.info('Setting callback_url: %s', callback_url)

        # Get signature return URL
        if return_url is None:
            return_url = cla.utils.get_active_signature_return_url(user_id, signature_metadata)
            cla.log.info('Setting signature return_url to %s', return_url)
        if return_url is None:
            return {'user_id': str(user_id),
                    'project_id': str(project_id),
                    'signature_id': None,
                    'sign_url': None,
                    'error': 'No active signature found for user - cannot generate return_url without knowing where the user came from'}

        # Get latest document
        try:
            document = project.get_latest_individual_document()
        except DoesNotExist as err:
            return {'errors': {'project_id': str(err)}}

        default_cla_values = create_default_individual_values(user)

        # Create new Signature object
        signature = Signature(signature_id=str(uuid.uuid4()),
                                signature_project_id=project_id,
                                signature_document_major_version=document.get_document_major_version(),
                                signature_document_minor_version=document.get_document_minor_version(),
                                signature_reference_id=user_id,
                                signature_reference_type='user',
                                signature_type='cla',
                                signature_return_url_type = 'Github',
                                signature_signed=False,
                                signature_approved=True,
                                signature_return_url=return_url,
                                signature_callback_url=callback_url)

        # Populate sign url
        self.populate_sign_url(signature, callback_url, default_values=default_cla_values)

        # Save signature
        signature.save()

        return {'user_id': str(user_id),
                'project_id': project_id,
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    def request_individual_signature_gerrit(self, project_id, user_id, return_url=None):
        cla.log.info('Creating new Gerrit signature for user %s on project %s', user_id, project_id)

        # Ensure this is a valid user
        user_id = str(user_id)
        try:
            user = User()
            user.load(user_id)
        except DoesNotExist as err:
            cla.log.warning('User ID not found when trying to request a signature: %s',
                            user_id)
            return {'errors': {'user_id': str(err)}}

        # Ensure the project exists
        try:
            project = Project()
            project.load(project_id)
        except DoesNotExist as err:
            cla.log.error('Project ID not found when trying to request a signature: %s',
                        project_id)
            return {'errors': {'project_id': str(err)}}


        # Check for active signature object with this project. If the user has
        # signed the most recent major version, they do not need to sign again.
        latest_signature = user.get_latest_signature(user, str(project_id))
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
                                signature_return_url_type = 'Gerrit',
                                signature_signed=False,
                                signature_approved=True,
                                signature_return_url=return_url,
                                signature_callback_url=callback_url)

        self.populate_sign_url(signature, callback_url, default_values=default_cla_values)

        # Save signature
        signature.save()

        return {'user_id': str(user_id),
                'project_id': project_id,
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    # Returns the current project, user and employee Signature, if such
    # a signature exists.
    def check_and_prepare_employee_signature(self,
                                             project_id: str,
                                             company_id: str,
                                             user_id: str
                                             ) -> Tuple[
                                                 Project,
                                                 User,
                                                 Optional[Signature]
                                             ] :
         # Ensure the project exists
        project = Project()
        try:
            project.load(str(project_id))
        except DoesNotExist:
            raise ProjectDoesNotExist

        # Ensure the company exists
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist:
            raise CompanyDoesNotExist

        # Ensure the user exists
        user = User()
        try:
            user.load(str(user_id))
        except DoesNotExist:
            raise UserDoesNotExist

        # Ensure the company actually has a CCLA with this project.
        ccla_signatures = Signature().get_signatures_by_project(
            project_id,
            signature_reference_type='company',
            signature_reference_id=company.get_company_id()
        )
        if len(ccla_signatures) < 1:
            raise CCLANotFound

        ccla_signature = ccla_signatures[0]

        # Ensure user hasn't already signed this signature.
        employee_signatures = Signature().get_signatures_by_project(
            project_id,
            signature_reference_type='user',
            signature_reference_id=user_id,
            signature_user_ccla_company_id=company_id
        )
        if len(employee_signatures) > 0:
            cla.log.info('Employee signature already exists for this project')
            return project, user, employee_signatures[0]

        # Ensure user is whitelisted for this company.
        if not user.is_whitelisted(ccla_signature):
            raise UserNotWhitelisted

        # Assume this company is the user's employer.
        user.set_user_company_id(str(company_id))
        user.save()

        return project, user, None

    def request_employee_signature(self, project_id, company_id, user_id, return_url=None):
        try:
            project, \
            user, \
            employee_signature = self.check_and_prepare_employee_signature(str(project_id),
                                                                           str(company_id),
                                                                           str(user_id))
        except ProjectDoesNotExist:
            return {'errors': {'project_id': 'Project ({}) does not exist.'.format(project_id)}}
        except CompanyDoesNotExist:
            return {'errors': {'company_id': 'Company ({}) does not exist.'.format(company_id)}}
        except UserDoesNotExist:
            return {'errors': {'user_id': 'User ({}) does not exist.'.format(user_id)}}
        except CCLANotFound:
            return {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
        except UserNotWhitelisted:
            return {'errors': {'ccla_whitelist': 'No user email whitelisted for this ccla'}}

        if employee_signature is not None:
            return employee_signature.to_dict()

        # Requires us to know where the user came from.
        signature_metadata = cla.utils.get_active_signature_metadata(user_id)
        if return_url is None:
            return_url = cla.utils.get_active_signature_return_url(user_id, signature_metadata)

        # return_url may still be empty at this point - the console will deal with it
        new_signature = Signature(signature_id=str(uuid.uuid4()),
                                signature_project_id=project_id,
                                signature_document_minor_version=0,
                                signature_document_major_version=0,
                                signature_reference_id=user_id,
                                signature_reference_type='user',
                                signature_type='cla',
                                signature_signed=True,
                                signature_approved=True,
                                signature_return_url=return_url,
                                signature_user_ccla_company_id=company_id)
        new_signature.save()

        # If the project does not require an ICLA to be signed, update the pull request and remove the active
        # signature metadata.
        if not project.get_project_ccla_requires_icla_signature():
            cla.log.info('Project does not requires ICLA signature from employee - updating PR')
            github_repository_id = signature_metadata['repository_id']
            change_request_id = signature_metadata['pull_request_id']
            
            # Get repository
            repository = Repository() 
            try: 
                repository.get_repository_by_external_id(github_repository_id, 'github')
            except DoesNotExist:
                return {'errors': {'repository_id': 'The given repository_id does not exist. '}}
                
            # Get organization that the repository is configured for
            organization = cla.utils.get_github_organization_instance()
            organization.get_organization(repository.get_repository_organization_name())
            
            # Get installation ID of the organization.
            installation_id = organization.get_organization_installation_id()
            update_repository_provider(installation_id, github_repository_id, change_request_id)

            cla.utils.delete_active_signature_metadata(user.get_user_id())
        else:
            cla.log.info('Project requires ICLA signature from employee - PR has been left unchanged')

        return new_signature.to_dict()

    def request_employee_signature_gerrit(self, project_id, company_id, user_id, return_url=None):
        try:
            _, \
            user, \
            employee_signature = self.check_and_prepare_employee_signature(str(project_id),
                                                                           str(company_id),
                                                                           str(user_id))
        except ProjectDoesNotExist:
            return {'errors': {'project_id': 'Project ({}) does not exist.'.format(project_id)}}
        except CompanyDoesNotExist:
            return {'errors': {'company_id': 'Company ({}) does not exist.'.format(company_id)}}
        except UserDoesNotExist:
            return {'errors': {'user_id': 'User ({}) does not exist.'.format(user_id)}}
        except CCLANotFound:
            return {'errors': {'missing_ccla': 'Company does not have CCLA with this project'}}
        except UserNotWhitelisted:
            return {'errors': {'ccla_whitelist': 'No user email whitelisted for this ccla'}}

        if employee_signature is not None:
            return employee_signature.to_dict()

        # Retrieve Gerrits by Project reference ID
        try:
            gerrits = Gerrit().get_gerrit_by_project_id(project_id)
        except DoesNotExist as err:
            cla.log.error('Cannot load Gerrit instance for the given project: %s',project_id)
            return {'errors': {'missing_gerrit': str(err)}}

        new_signature = Signature(signature_id=str(uuid.uuid4()),
                                signature_project_id=project_id,
                                signature_document_minor_version=0,
                                signature_document_major_version=0,
                                signature_reference_id=user_id,
                                signature_reference_type='user',
                                signature_type='cla',
                                signature_signed=True,
                                signature_approved=True,
                                signature_return_url=return_url,
                                signature_user_ccla_company_id=company_id)

        # Save signature before adding user to the LDAP Group. 
        new_signature.save()

        # Get lf_username from User 
        lf_username = user.get_lf_username()

        for gerrit in gerrits:
            # For every Gerrit Instance of this project, add the user to the LDAP Group.
            # this way we are able to keep track of signed signatures when user fails to be added to the LDAP GROUP.
            group_id = gerrit.get_group_id_ccla()
            # Add the user to the LDAP Group
            try:
                lf_group.add_user_to_group(group_id, lf_username)
            except Exception as e:
                cla.log.error('Failed in adding user to the LDAP group.%s', e)
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
        cla.log.info('Validating company %s on project %s', company_id, project_id)

        # Ensure the project exists
        project = Project()
        try:
            project.load(str(project_id))
        except DoesNotExist as err:
            return {'errors': {'project_id': str(err)}}
        
        # Ensure the company exists
        company = Company()
        try:
            company.load(str(company_id))
        except DoesNotExist as err:
            return {'errors': {'company_id': str(err)}}

        # Ensure the managers list is not empty
        company_model = Company()
        managers = company_model.get_managers_by_company_acl(company.get_company_acl())
        if len(managers) == 0:
            return {'errors': {'company_acl': 'Company ACL is empty'}}

        # Find the manager user object
        for manager in managers:
            if manager.get_lf_username() == auth_user.username:
                break
        else:
            # Return an error if the manager is not found in the managers list
            return {'errors': {'manager': 'CLA Manager not found'}}

        # Get CLA Managers. In the future, we will support contribuitors
        scheduleA = generate_manager_and_contributor_list(
            [(manager.get_user_name(), manager.get_user_email()) for manager in managers]
        )

        default_cla_values = create_default_company_values(company, manager.get_user_name(), manager.get_user_email(), scheduleA)

        # Ensure the contract group has a CCLA
        last_document = project.get_latest_corporate_document()
        if last_document is None or \
            last_document.get_document_major_version() is None or \
            last_document.get_document_minor_version() is None:
            cla.log.info('Contract Group does not have a CCLA: {}'.format(project_id))
            return {'errors': {'project_id': 'Contract Group does not support CCLAs.'}}

        # Ensure the company doesn't already have a CCLA with this project. 
        # and the user is about to sign the ccla manually 
        cla.log.info('Checking if a signature exists')
        latest_signature = company.get_latest_signature(str(project_id))
        if latest_signature is not None and \
        last_document.get_document_major_version() == latest_signature.get_signature_document_major_version():
            cla.log.info('CCLA signature object already exists for company %s on project %s', company_id, project_id)
            if latest_signature.get_signature_signed():
                cla.log.info('CCLA signature object already signed')
                return {'errors': {'signature_id': 'Company has already signed CCLA with this project'}}
            else:
                if not send_as_email:
                    #signature object exists but still has not been manually signed.
                    cla.log.info('CCLA signature object still missing signature')
                else:
                    #signature object exists and the user wants to send it to a corp authority.
                    callback_url = self._get_corporate_signature_callback_url(str(project_id), str(company_id))
                    self.populate_sign_url(latest_signature, callback_url, send_as_email, authority_name, authority_email, default_values=default_cla_values)
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

        if(not send_as_email): #get return url only for manual signing through console
            cla.log.info('Setting signature return_url to %s', return_url)
            signature.set_signature_return_url(return_url)

        self.populate_sign_url(signature, callback_url, send_as_email, authority_name, authority_email, default_values=default_cla_values)
        signature.save()

        return {'company_id': str(company_id),
                'project_id': str(project_id),
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    def populate_sign_url(self, signature, callback_url=None, send_as_email=False,
    authority_name=None, authority_email=None, default_values: Optional[Dict[str, Any]] = None): # pylint: disable=too-many-locals
        cla.log.debug('Populating sign_url for signature %s', signature.get_signature_id())
        sig_type = signature.get_signature_reference_type()
        user = User()
        
        # Assume the company manager is signing the CCLA
        if sig_type == 'company': 
            company = Company()
            company.load(signature.get_signature_reference_id())
            try:
                user.load(company.get_company_manager_id())
                name = user.get_user_name()
            except DoesNotExist:
                cla.log.error('No CLA manager associated with this company - can not sign CCLA')
                return
        else:
            if not send_as_email: 
                # sig_type == 'user'
                user.load(signature.get_signature_reference_id())            
                name = user.get_user_name()
                if name is None:
                    name = 'Unknown'
        
        
        # Fetch the document to sign.
        project = Project()
        project.load(signature.get_signature_project_id())
        if sig_type == 'company':
            document = project.get_project_corporate_document()
            if document is None:
                cla.log.error('Could not get sign url for project %s: Project has no corporate \
                               CLA document set', project.get_project_id())
                return
        else: # sig_type == 'user'
            document = project.get_project_individual_document()
            if document is None:
                cla.log.error('Could not get sign url for project %s: Project has no individual \
                               CLA document set', project.get_project_id())
                return

        # Not sure what should be put in as documentId.
        document_id = uuid.uuid4().int & (1<<16)-1 # Random 16bit integer -.pylint: disable=no-member
        tabs = get_docusign_tabs_from_document(document, document_id, default_values=default_values)

        if send_as_email:
            # Sending email to authority
            email = authority_email
            name = authority_name
        else:
            # User email
            email = user.get_user_email()

        if send_as_email: 
            # Not assigning a clientUserId sends an email. 
            signer = pydocusign.Signer(email=email,
                                    name=name,
                                    recipientId=1,
                                    tabs=tabs, 
                                    emailSubject='CLA Sign Request',
                                    emailBody='CLA Sign Request for %s'
                                    %authority_email,
                                    supportedLanguage='en',
                                    )
        else:
            # Assigning a clientUserId does not send an email.
            # It assumes that the user handles the communication with the client. 
            # In this case, the user opened the docusign document to manually sign it. 
            # Thus the email does not need to be sent. 
            signer = pydocusign.Signer(email=email,
                                    name=name,
                                    recipientId=1,
                                    clientUserId=signature.get_signature_id(),
                                    tabs=tabs, 
                                    emailSubject='CLA Sign Request',
                                    emailBody='CLA Sign Request for %s'
                                    %user.get_user_email(),
                                    supportedLanguage='en',
                                    )
        
        content_type = document.get_document_content_type()
        if content_type.startswith('url+'):
            pdf_url = document.get_document_content()
            pdf = self.get_document_resource(pdf_url)
        else:
            content = document.get_document_content()
            pdf = io.BytesIO(content)
        doc_name = document.get_document_name()
        document = pydocusign.Document(name=doc_name,
                                       documentId=document_id,
                                       data=pdf)                         

        if callback_url is not None:
            # Webhook properties for callbacks after the user signs the document.
            # Ensure that a webhook is returned on the status "Completed" where 
            # all signers on a document finish signing the document. 
            recipient_events = [{"recipientEventStatusCode": "Completed"}]
            event_notification= pydocusign.EventNotification(url=callback_url,
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

        recipient = envelope.recipients[0]

        if(not send_as_email):
            # The URL the user will be redirected to after signing.
            # This route will be in charge of extracting the signature's return_url and redirecting.
            return_url = os.path.join(api_base_url, 'v2/return-url', str(recipient.clientUserId))
            
            cla.log.info("return-url " + return_url)
            sign_url = self.get_sign_url(envelope, recipient, return_url)
            cla.log.info('Setting signature sign_url to %s', sign_url)
            signature.set_signature_sign_url(sign_url)

    def signed_individual_callback(self, content, installation_id, github_repository_id, change_request_id):
        """
        Will be called on ICLA signature callback, but also when a document has been
        opened by a user - no action required then.
        """
        cla.log.debug('Docusign ICLA signed callback POST data: %s', content)
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
            signature.save()
            # Send user their signed document.
            user = User()
            user.load(signature.get_signature_reference_id())
            # Remove the active signature metadata.
            cla.utils.delete_active_signature_metadata(user.get_user_id())
            # Get signed document
            document_data = self.get_signed_document(self, envelope_id, user)
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
            cla.log.error('Received CCLA signed callback from signing service provider for an unknown company: %s', company_id)
            return
        if status == 'Completed' and not signature.get_signature_signed():
            cla.log.info('CCLA signature signed (%s)', signature_id)
            signature.set_signature_signed(True)
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
        filename = str.join('/', ('', 'contract-group', project_id, cla_type, identifier, signature_id + '.pdf'))
        self.s3storage.store(filename, document_data)
        
    def get_document_resource(self, url): # pylint: disable=no-self-use
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

    def get_sign_url(self, envelope, recipient, return_url): # pylint:disable=no-self-use
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
        split_url = return_url.split('/') #parse repo name from URL
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
                                    scheduleA: str) -> Dict[str, Any]:
    values = {}

    if company is not None and \
        company.get_company_name() is not None:
        values['corporation_name'] = company.get_company_name()
        values['corporation'] = company.get_company_name()

    if signatory_name is not None:
        values['point_of_contact'] = signatory_name

    if signatory_email is not None:
        values['email'] = signatory_email

    if scheduleA is not None:
        values['scheduleA'] = scheduleA

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
def generate_manager_and_contributor_list( managers, contributors=None):
    lines = []

    for manager in managers:
        lines.append('CLA Manager: {}, {}'.format(manager[0], manager[1]))

    if contributors is not None:
        for contributor in contributors:
            lines.append('{}, {}'.format(contributor[0], contributor[1]))

    lines = '\n'.join([str(line) for line in lines])

    return lines
