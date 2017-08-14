"""
Tests having to do with the agreements.
"""

import unittest
import uuid
import hug

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_email_service

class AgreementTestCase(CLATestCase):
    """Agreement test cases."""
    def test_get_agreements(self):
        """Tests for getting all agreements."""
        response = hug.test.get(cla.routes, '/v1/agreement')
        self.assertEqual(response.data, [])
        project = self.create_project()
        user = self.create_user()
        organization = self.create_organization()
        self.create_document(project['project_id'])
        agreement1 = self.create_agreement(project['project_id'],
                                           user['user_id'],
                                           'user')
        agreement2 = self.create_agreement(project['project_id'],
                                           organization['organization_id'],
                                           'organization')
        self.assertEqual(agreement2,
                          {'errors': {'agreement_project_id': \
                                      'No corporate document exists for this project'}})
        self.create_document(project['project_id'], document_type='corporate')
        agreement3 = self.create_agreement(project['project_id'],
                                           organization['organization_id'],
                                           'organization')
        response = hug.test.get(cla.routes, '/v1/agreement')
        self.assertEqual(len(response.data), 2)

    def test_get_agreement(self):
        """Tests for getting individual agreements."""
        response = hug.test.get(cla.routes, '/v1/agreement/1')
        self.assertEqual(response.data, {'errors': {'agreement_id': 'Invalid UUID provided'}})
        response = hug.test.get(cla.routes, '/v1/agreement/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'agreement_id': 'Agreement not found'}})
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        agreement = self.create_agreement(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/agreement/' + agreement['agreement_id'])
        self.assertEqual(response.data['agreement_id'], agreement['agreement_id'])

    def test_post_agreement(self):
        """Tests for creating agreements."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        response = self.create_agreement(project['project_id'],
                                         user['user_id'],
                                         'invalid')
        self.assertEqual(response, {'errors': {'agreement_reference_type': \
                                    'Invalid value passed. The accepted values are: (organization|user)'}})
        response = self.create_agreement(project['project_id'],
                                         user['user_id'],
                                         'user', agreement_type='invalid')
        self.assertEqual(response, {'errors': {'agreement_type': \
                                    'Invalid value passed. The accepted values are: (cla|dco)'}})
        response = self.create_agreement(project['project_id'], user['user_id'], 'user', agreement_signed='1')
        self.assertTrue(response['agreement_signed'])
        response = self.create_agreement(project['project_id'], user['user_id'], 'user', agreement_approved='true')
        self.assertTrue(response['agreement_approved'])

    def test_put_agreement(self):
        """Tests for updating agreements."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        agreement = self.create_agreement(project['project_id'], user['user_id'], 'user')
        response = hug.test.put(cla.routes, '/v1/agreement', {'meh': 'dco'})
        self.assertEqual(response.data, \
            {'errors': {'agreement_id': "Required parameter 'agreement_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': str(uuid.uuid4()), 'meh': 'dco'})
        self.assertEqual(response.data, {'errors': {'agreement_id': 'Agreement not found'}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_type': 'dco'})
        self.assertEqual(response.data, \
            {'errors': {'agreement_id': "Required parameter 'agreement_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': str(uuid.uuid4()), 'agreement_type': 'dco'})
        self.assertEqual(response.data, {'errors': {'agreement_id': 'Agreement not found'}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': agreement['agreement_id'], 'agreement_type': 'test'})
        self.assertEqual(response.data, {'errors': {'agreement_type': "Invalid value passed. The accepted values are: (cla|dco)"}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': agreement['agreement_id'], 'agreement_type': 'dco'})
        self.assertEqual(response.data['agreement_type'], 'dco')
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': agreement['agreement_id'], 'agreement_signed': 'test'})
        self.assertEqual(response.data, {'errors': {'agreement_signed': "Invalid value passed in for true/false field"}})
        response = hug.test.put(cla.routes, '/v1/agreement', {'agreement_id': agreement['agreement_id'], 'agreement_approved': 'false'})
        self.assertEqual(response.data['agreement_approved'], False)

    def test_delete_agreement(self):
        """Tests for deleting agreements."""
        project = self.create_project()
        user = self.create_user()
        self.create_document(project['project_id'])
        agreement = self.create_agreement(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/agreement')
        self.assertTrue(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/agreement/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'agreement_id': 'Agreement not found'}})
        response = hug.test.delete(cla.routes, '/v1/agreement/' + agreement['agreement_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/agreement')
        self.assertEqual(len(response.data), 0)

if __name__ == '__main__':
    unittest.main()
