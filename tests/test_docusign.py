"""
Tests having to do with the DocuSign signing service.
"""

import unittest
import hug

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_user_instance, get_signing_service

class DocuSignTestCase(CLATestCase):
    """DocuSign test cases."""
    def test_request_individual_cncf_signature(self):
        """
        Tests for requesting a individual signature.
        Communicates with the signing service provider.
        """

        return # Only enabled temporarily when testing out new templates.

        test_signing_service = cla.conf['SIGNING_SERVICE']
        cla.conf['SIGNING_SERVICE'] = 'DocuSign'
        project = self.create_project()

        data = {'document_name': 'test.pdf',
                'document_preamble': 'Preamble here',
                'document_legal_entity_name': 'Legal Entity Name',
                'template_name': 'CNCFTemplate'}
        hug.test.post(cla.routes, '/v1/project/' + project['project_id'] + '/document/template/individual', data)
        user = self.create_user()
        data = {'project_id': project['project_id'],
                'user_id': user['user_id'],
                'return_url': 'http://return-url.com/done'}
        result = hug.test.post(cla.routes, '/v1/request-individual-signature', data)
        print(result.data)
        cla.conf['SIGNING_SERVICE'] = test_signing_service

    def test_request_corporate_cncf_signature(self):
        """
        Tests for requesting a corporate signature.
        Communicates with the signing service provider.
        """

        #return # Only enabled temporarily when testing out new templates.

        test_signing_service = cla.conf['SIGNING_SERVICE']
        cla.conf['SIGNING_SERVICE'] = 'DocuSign'
        project = self.create_project()
        manager = self.create_user()
        company = self.create_company(manager['user_id'])
        data = {'document_name': 'test.pdf',
                'document_preamble': 'Preamble here',
                'document_legal_entity_name': 'Legal Entity Name',
                'template_name': 'CNCFTemplate'}
        hug.test.post(cla.routes, '/v1/project/' + project['project_id'] + '/document/template/corporate', data)
        data = {'project_id': project['project_id'],
                'company_id': company['company_id'],
                'return_url': 'http://return-url.com/done'}
        result = hug.test.post(cla.routes, '/v1/request-corporate-signature', data)
        print(result.data)
        cla.conf['SIGNING_SERVICE'] = test_signing_service

    def test_request_signature(self):
        """Tests for requesting signature. Uses mock objects."""
        project = self.create_project()
        document = self.create_document(project['project_id'])
        user = self.create_user()
        data = {'project_id': project['project_id'],
                'user_id': user['user_id'],
                'return_url': 'http://return-url.com/done'}
        result = hug.test.post(cla.routes, '/v1/request-individual-signature', data)
        signature_id = result.data['signature_id']
        expected = {'user_id': user['user_id'],
                    'project_id': project['project_id'],
                    'signature_id': signature_id,
                    'sign_url': 'http://signing-service.com/send-user-here'}
        self.assertEqual(result.data, expected)

    def test_signed_callback(self):
        """Tests for the DocusSign signed callback."""
        # TODO: Implement this test.
        return
        content = open('resources/docusign_callback_payload.xml').read()
        get_signing_service().signed_callback(content, None, None)

    def test_send_signed_document(self):
        """Tests for the DocusSign send signed document method."""
        user = get_user_instance()
        user.set_user_email('test@test.com')
        get_signing_service().send_signed_document('fake-envelope-id', user)

if __name__ == '__main__':
    unittest.main()
