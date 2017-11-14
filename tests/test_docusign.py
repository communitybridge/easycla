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
    def test_request_signature(self):
        """Tests for requesting signature."""
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
