"""
Tests having to do with the user endpoints.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class UserTestCase(CLATestCase):
    def test_get_users(self):
        """Test for getting a list of all users."""
        response = hug.test.get(cla.routes, '/v1/user')
        self.assertEqual(response.data, [])
        self.create_user()
        self.create_user()
        self.create_user()
        response = hug.test.get(cla.routes, '/v1/user')
        self.assertEqual(len(response.data), 3)

    def test_get_user(self):
        response = hug.test.get(cla.routes, '/v1/user/fake')
        self.assertEqual(response.data, {'errors': {'user_id': 'Invalid UUID provided'}})
        response = hug.test.get(cla.routes, '/v1/user/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'user_id': 'User not found'}})
        user = self.create_user()
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'])
        data = response.data
        del data['date_created']
        del data['date_modified']
        self.assertEqual(response.data, data)

    def test_post_user(self):
        self.create_user(user_email='user1@email.com', user_github_id=11111)

    def test_put_user(self):
        response = hug.test.put(cla.routes, '/v1/user', {'user_github_id': 99999})
        self.assertEqual(response.data,
                         {'errors':
                             {'user_id': "Required parameter 'user_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/user', {'user_id': 'fake-id', 'user_github_id': 99999})
        self.assertEqual(response.data, {'errors': {'user_id': 'Invalid UUID provided'}})
        user = self.create_user()
        self.assertEqual(user['user_github_id'], '12345')
        response = hug.test.put(cla.routes, '/v1/user', {'user_id': user['user_id'], 'user_github_id': 99999})
        self.assertEqual(response.data['user_github_id'], '99999')

    def test_delete_user(self):
        user = self.create_user()
        response = hug.test.delete(cla.routes, '/v1/user/' + user['user_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'])
        self.assertEqual(response.data, {'errors': {'user_id': 'User not found'}})

    def test_get_user_signatures(self):
        user = self.create_user()
        project = self.create_project('test-project')
        self.create_document(project['project_id'])
        self.create_signature(project['project_id'], user['user_id'], 'user')
        self.create_signature(project['project_id'], user['user_id'], 'user')
        self.create_signature(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/user/%s/signatures' %user['user_id'])
        self.assertTrue('errors' not in response.data)
        self.assertTrue(len(response.data) == 3)

    def test_get_user_by_email(self):
        user = self.create_user()
        response = hug.test.get(cla.routes, '/v1/user/email/' + user['user_emails'][0])
        self.assertEqual(response.data['user_id'], user['user_id'])

    def test_get_user_by_github_id(self):
        user = self.create_user()
        response = hug.test.get(cla.routes, '/v1/user/github/' + user['user_github_id'])
        self.assertEqual(response.data['user_id'], user['user_id'])

    #def test_get_users_by_company(self):
        #company = self.create_company()
        #self.create_user(user_company_id=company['company_id'])
        #self.create_user()
        #self.create_user(user_company_id=company['company_id'])
        #response = hug.test.get(cla.routes, '/v1/users/company/' + \
                                            #company['company_id'])
        #self.assertTrue(len(response.data) == 2)

if __name__ == '__main__':
    unittest.main()
