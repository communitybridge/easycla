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

    def test_employee_signature(self):
        """Test for creating an employee signature."""
        user = self.create_user(user_email='test@test.org')
        project = self.create_project('test-project')
        company = self.create_company(company_whitelist=['test@test.org'])
        self.create_document(project['project_id'], 'corporate')
        # Create Corporate signature.
        self.create_signature(project['project_id'], company['company_id'], 'company')
        # Create Employee signature.
        data = {'project_id': project['project_id'],
                'company_id': company['company_id'],
                'user_id': user['user_id']}
        hug.test.post(cla.routes, '/v1/request-employee-signature', data)

    def test_get_user_signatures(self):
        """Test getting ICLA and employee signatures."""
        user = self.create_user()
        project = self.create_project('test-project')
        company = self.create_company()
        self.create_document(project['project_id'], 'individual')
        self.create_document(project['project_id'], 'corporate')
        # Create Corporate signature.
        self.create_signature(project['project_id'], company['company_id'], 'company')
        # Create Employee signature.
        emp_sig = self.create_signature(project['project_id'], user['user_id'], 'user', signature_user_ccla_company_id=company['company_id'])
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'] + '/project/' + project['project_id'] + '/last-signature/' + company['company_id'])
        self.assertEqual(response.data['signature_user_ccla_company_id'], company['company_id'])
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'] + '/project/' + project['project_id'] + '/last-signature')
        self.assertTrue(response.data is None)
        # Create Individual signature.
        indiv_sig = self.create_signature(project['project_id'], user['user_id'], 'user')
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'] + '/project/' + project['project_id'] + '/last-signature')
        self.assertFalse(response.data['requires_resigning'])
        self.assertEqual(response.data['signature_id'], indiv_sig['signature_id'])
        # Make employee signature stale.
        self.create_document(project['project_id'], 'corporate', new_major_version=True)
        response = hug.test.get(cla.routes, '/v1/user/' + user['user_id'] + '/project/' + project['project_id'] + '/last-signature/' + company['company_id'])
        self.assertTrue(response.data['requires_resigning'])

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
