"""
Tests having to do with the signature endpoints.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class SignatureTestCase(CLATestCase):
    def test_get_user_signatures(self):
        """Test for getting a list of a user's signatures."""
        user = self.create_user()
        project1 = self.create_project('test-project1')
        project2 = self.create_project('test-project2')
        self.create_document(project1['project_id'])
        self.create_document(project2['project_id'])
        self.create_signature(project1['project_id'], user['user_id'], 'user')
        self.create_signature(project2['project_id'], user['user_id'], 'user')
        self.create_signature(project2['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/signatures/user/' + user['user_id'])
        self.assertEqual(len(response.data), 3)
        response = hug.test.get(cla.routes, '/v1/signatures/user/' + user['user_id'] +
                                '/project/' + project1['project_id'])
        self.assertEqual(len(response.data), 1)
        response = hug.test.get(cla.routes, '/v1/signatures/user/' + user['user_id'] +
                                '/project/' + project2['project_id'])
        self.assertEqual(len(response.data), 2)

    def test_various_major_versions(self):
        """Test out-dated and invalidated signatures."""
        user_data = self.create_user()
        project = self.create_project('test-project')
        self.create_document(project['project_id'])
        self.create_signature(project['project_id'], user_data['user_id'], 'user')
        user = cla.utils.get_user_instance()
        user.load(user_data['user_id'])
        signed = cla.utils.user_signed_project_signature(user, project['project_id'])
        self.assertTrue(signed)
        self.create_document(project['project_id'])
        signed = cla.utils.user_signed_project_signature(user, project['project_id'])
        self.assertTrue(signed)
        self.create_document(project['project_id'], new_major_version=True)
        signed = cla.utils.user_signed_project_signature(user, project['project_id'])
        self.assertFalse(signed)

if __name__ == '__main__':
    unittest.main()
