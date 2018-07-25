"""
Tests having to do with the signatures.
"""

import unittest
import uuid
import hug

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_email_service

class SignatureTestCase(CLATestCase):
    """Signature test cases."""
    def test_get_signatures(self):
        """Tests for getting all signatures."""
        response = hug.test.get(cla.routes, '/v1/signature')
        self.assertEqual(response.data, [])
        project = self.create_project()
        user = self.create_user()
        company = self.create_company()
        self.create_document(project['project_id'])
        signature1 = self.create_signature(project['project_id'],
                                           user['user_id'],
                                           'user')
        signature2 = self.create_signature(project['project_id'],
                                           company['company_id'],
                                           'company')
        self.assertEqual(signature2,
                          {'errors': {'signature_project_id': \
                                      'No corporate document exists for this project'}})
        self.create_document(project['project_id'], document_type='corporate')
        signature3 = self.create_signature(project['project_id'],
                                           company['company_id'],
                                           'company')
        response = hug.test.get(cla.routes, '/v1/signature')
        self.assertEqual(len(response.data), 2)

    def test_get_signature(self):
        """Tests for getting individual signatures."""
        response = hug.test.get(cla.routes, '/v1/signature/1')
        self.assertEqual(response.data, {'errors': {'signature_id': 'Invalid UUID provided'}})
        response = hug.test.get(cla.routes, '/v1/signature/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'signature_id': 'Signature not found'}})
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        signature = self.create_signature(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/signature/' + signature['signature_id'])
        self.assertEqual(response.data['signature_id'], signature['signature_id'])

    def test_post_signature(self):
        """Tests for creating signatures."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        response = self.create_signature(project['project_id'],
                                         user['user_id'],
                                         'invalid')
        self.assertEqual(response, {'errors': {'signature_reference_type': \
                                    'Invalid value passed. The accepted values are: (company|user)'}})
        response = self.create_signature(project['project_id'],
                                         user['user_id'],
                                         'user', signature_type='invalid')
        self.assertEqual(response, {'errors': {'signature_type': \
                                    'Invalid value passed. The accepted values are: (cla|dco)'}})
        response = self.create_signature(project['project_id'], user['user_id'], 'user', signature_signed='1')
        self.assertTrue(response['signature_signed'])
        response = self.create_signature(project['project_id'], user['user_id'], 'user', signature_approved='true')
        self.assertTrue(response['signature_approved'])

    def test_put_signature(self):
        """Tests for updating signatures."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        signature = self.create_signature(project['project_id'], user['user_id'], 'user')
        response = hug.test.put(cla.routes, '/v1/signature', {'meh': 'dco'})
        self.assertEqual(response.data, \
            {'errors': {'signature_id': "Required parameter 'signature_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': str(uuid.uuid4()), 'meh': 'dco'})
        self.assertEqual(response.data, {'errors': {'signature_id': 'Signature not found'}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_type': 'dco'})
        self.assertEqual(response.data, \
            {'errors': {'signature_id': "Required parameter 'signature_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': str(uuid.uuid4()), 'signature_type': 'dco'})
        self.assertEqual(response.data, {'errors': {'signature_id': 'Signature not found'}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': signature['signature_id'], 'signature_type': 'test'})
        self.assertEqual(response.data, {'errors': {'signature_type': "Invalid value passed. The accepted values are: (cla|dco)"}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': signature['signature_id'], 'signature_type': 'dco'})
        self.assertEqual(response.data['signature_type'], 'dco')
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': signature['signature_id'], 'signature_signed': 'test'})
        self.assertEqual(response.data, {'errors': {'signature_signed': "Invalid value passed in for true/false field"}})
        response = hug.test.put(cla.routes, '/v1/signature', {'signature_id': signature['signature_id'], 'signature_approved': 'false'})
        self.assertEqual(response.data['signature_approved'], False)

    def test_delete_signature(self):
        """Tests for deleting signatures."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        signature = self.create_signature(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/signature')
        self.assertTrue(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/signature/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'signature_id': 'Signature not found'}})
        response = hug.test.delete(cla.routes, '/v1/signature/' + signature['signature_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/signature')
        self.assertEqual(len(response.data), 0)

if __name__ == '__main__':
    unittest.main()
