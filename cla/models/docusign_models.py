"""
Easily perform signing workflows using DocuSign signing service with pydocusign.
"""

import io
import uuid
import urllib.request
import xml.etree.ElementTree as ET
import pydocusign
from pydocusign.exceptions import DocuSignException
import cla
from cla.models import signing_service_interface, DoesNotExist

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

    def initialize(self, config):
        root_url = config['DOCUSIGN_ROOT_URL']
        username = config['DOCUSIGN_USERNAME']
        password = config['DOCUSIGN_PASSWORD']
        integrator_key = config['DOCUSIGN_INTEGRATOR_KEY']
        self.client = pydocusign.DocuSignClient(root_url=root_url,
                                                username=username,
                                                password=password,
                                                integrator_key=integrator_key)

    def request_signature(self, project_id, user_id, return_url=None):
        cla.log.info('Creating new signature for user %s on project %s', user_id, project_id)
        # Ensure this is a valid user.
        user_id = str(user_id)
        try:
            user = cla.utils.get_user_instance()
            user.load(user_id)
        except DoesNotExist as err:
            cla.log.warning('User ID not found when trying to request a signature: %s',
                            user_id)
            return {'errors': {'user_id': str(err)}}
        # Check for active signature object with this project.
        latest_signature = cla.utils.get_user_latest_signature(user, project_id)
        if latest_signature is not None:
            cla.log.info('User already has a signatures with this project: %s', \
                         latest_signature.get_signature_id())
            # TODO: Check versioning to determine whether to create new one or not.
            return {'user_id': user_id,
                    'project_id': project_id,
                    'signature_id': latest_signature.get_signature_id(),
                    'sign_url': latest_signature.get_signature_sign_url()}
        # Create new signature.
        signature = cla.utils.get_signature_instance()
        signature.set_signature_id(str(uuid.uuid4()))
        try:
            project = cla.utils.get_project_instance()
            project.load(project_id)
        except DoesNotExist as err:
            cla.log.error('Project ID not found when trying to request a signature: %s',
                        project_id)
            return {'errors': {'project_id': str(err)}}
        signature.set_signature_project_id(project_id)
        signature_metadata = cla.utils.get_active_signature_metadata(user_id)
        callback_url = cla.utils.get_active_signature_callback_url(user_id, signature_metadata)
        cla.log.info('Setting callback_url: %s', callback_url)
        signature.set_signature_callback_url(callback_url)
        # Requires us to know where the user came from.
        if return_url is None:
            return_url = cla.utils.get_active_signature_return_url(user_id, signature_metadata)
        if return_url is None:
            return {'user_id': str(user_id),
                    'project_id': str(project_id),
                    'signature_id': None,
                    'sign_url': None,
                    'error': 'No active signature found for user - cannot generate return_url \
                              without knowing where the user came from'}

        # Assume ICLA only for now.
        try:
            document = project.get_project_individual_document()
        except DoesNotExist as err:
            return {'errors': {'project_id': str(err)}}
        signature.set_signature_document_major_version(document.get_document_major_version())
        signature.set_signature_document_minor_version(document.get_document_minor_version())
        signature.set_signature_signed(False)
        signature.set_signature_approved(True)
        signature.set_signature_type('icla')
        signature.set_signature_reference_id(user_id)
        signature.set_signature_reference_type('user')
        cla.log.info('Setting signature return_url to %s', return_url)
        signature.set_signature_return_url(return_url)
        self.populate_sign_url(signature, callback_url)
        signature.save()
        return {'user_id': str(user_id),
                'project_id': project_id,
                'signature_id': signature.get_signature_id(),
                'sign_url': signature.get_signature_sign_url()}

    def populate_sign_url(self, signature, callback_url=None): # pylint: disable=too-many-locals
        cla.log.debug('Populating sign_url for signature %s', signature.get_signature_id())
        user = cla.utils.get_user_instance()
        user.load(signature.get_signature_reference_id())
        name = user.get_user_name()
        if name is None:
            name = 'Unknown'
        # Not sure what should be put in as documentId.
        document_id = uuid.uuid4().int & (1<<16)-1 # Random 16bit integer -.pylint: disable=no-member
        tab = pydocusign.SignHereTab(documentId=document_id,
                                     pageNumber=1,
                                     xPosition=280,
                                     yPosition=700)
        signer = pydocusign.Signer(email=user.get_user_email(),
                                   name=name,
                                   recipientId=1,
                                   clientUserId=signature.get_signature_id(),
                                   tabs=[tab], # Can be placed in DocuSign UI
                                   emailSubject='CLA Sign Request',
                                   emailBody='CLA Sign Request for %s'
                                   %user.get_user_email(),
                                   supportedLanguage='en')
        # Fetch the document to sign.
        # TODO: Need to support corporate CLAs?
        project = cla.utils.get_project_instance()
        project.load(signature.get_signature_project_id())
        document = project.get_project_individual_document()
        if document is None:
            cla.log.error('Could not get sign url for project %s: Project has no individual \
                          CLA document set', project.get_project_id())
            return
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
            event_notification = pydocusign.EventNotification(url=callback_url)
            envelope = pydocusign.Envelope(documents=[document],
                                           emailSubject='CLA Sign Request',
                                           emailBlurb='CLA Sign Request',
                                           eventNotification=event_notification,
                                           status=pydocusign.Envelope.STATUS_SENT, # Send now.
                                           recipients=[signer])
        else:
            envelope = pydocusign.Envelope(documents=[document],
                                           emailSubject='CLA Sign Request',
                                           emailBlurb='CLA Sign Request',
                                           status=pydocusign.Envelope.STATUS_SENT, # Send now.
                                           recipients=[signer])
        envelope = self.prepare_sign_request(envelope)
        cla.log.info('New envelope created in DocuSign: %s' %envelope.envelopeId)
        recipient = envelope.recipients[0]
        # The URL the user will be redirected to after signing.
        # This route will be in charge of extracting the signature's return_url and redirecting.
        return_url = cla.conf['BASE_URL'] + '/v1/return-url/' + str(recipient.clientUserId)
        sign_url = self.get_sign_url(envelope, recipient, return_url)
        cla.log.info('Setting signature sign_url to %s', sign_url)
        signature.set_signature_sign_url(sign_url)

    def signed_callback(self, content, installation_id, github_repository_id, change_request_id):
        """
        Will be called on signature callback, but also when a document has been
        opened by a user - no action required then.
        """
        cla.log.debug('Docusign signed callback POST data: %s', content)
        tree = ET.fromstring(content)
        # Get envelope ID.
        envelope_id = tree.find('.//' + self.TAGS['envelope_id']).text
        # Assume only one signature per signature.
        signature_id = tree.find('.//' + self.TAGS['client_user_id']).text
        signature = cla.utils.get_signature_instance()
        try:
            signature.load(signature_id)
        except DoesNotExist:
            cla.log.error('DocuSign callback returned signed info on invalid signature: %s',
                          content)
            return
        # Iterate through recipients and update the signature signature status if changed.
        elem = tree.find('.//' + self.TAGS['recipient_statuses'] +
                         '/' + self.TAGS['recipient_status'])
        status = elem.find(self.TAGS['status']).text
        if status == 'Completed' and not signature.get_signature_signed():
            cla.log.info('CLA signature signed (%s) - Notifying repository service provider',
                         signature_id)
            signature.set_signature_signed(True)
            signature.save()
            # Send user their signed document.
            # TODO: This currently only supports ICLAs.
            if signature.get_signature_reference_type() != 'user':
                cla.log.error('Trying to handle CCLA as a ICLA - not implemented yet')
                raise NotImplementedError()
            user = cla.utils.get_user_instance()
            user.load(signature.get_signature_reference_id())
            # Remove the active signature metadata.
            cla.utils.delete_active_signature_metadata(user.get_user_id())
            # Send email with signed document.
            self.send_signed_document(envelope_id, user)
            # Update the repository provider with this change.
            update_repository_provider(installation_id, github_repository_id, change_request_id)

    def send_signed_document(self, envelope_id, user):
        """Helper method to send the user their signed document."""
        # First, get the signed document from DocuSign.
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
            document_content = envelope.get_document(document['documentId'], self.client)
        except Exception as err:
            cla.log.error('Unknown error when trying to fetch signed document content ' + \
                          'for document ID %s: %s', document['documentId'], str(err))
            return
        # Second, prepare the email to the user.
        subject = 'CLA Signed Document'
        body = 'Thank you for signing the CLA! Your signed document is attached to this email.'
        recipient = user.get_user_email()
        filename = recipient + '-cla.pdf'
        attachment = {'type': 'content',
                      'content': document_content.read(),
                      'content-type': 'application/pdf',
                      'filename': filename}
        # Third, send the email.
        cla.log.info('Sending signed CLA document to %s', recipient)
        cla.utils.get_email_service().send(subject, body, recipient, attachment)

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
        return open('resources/test.pdf', 'rb')

    def prepare_sign_request(self, envelope):
        """
        Don't actually send the request when running tests.
        """
        recipients = []
        for recipient in envelope.recipients:
            recip = lambda: None
            recip.clientUserId = recipient.clientUserId
            recipients.append(recip)
        envelope = lambda: None
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

def update_repository_provider(installation_id, github_repository_id, change_request_id):
    """Helper method to notify the repository provider of successful signature."""
    repo_service = cla.utils.get_repository_service('github')
    repo_service.update_change_request(installation_id, github_repository_id, change_request_id)
