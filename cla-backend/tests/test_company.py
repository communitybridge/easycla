"""
Tests having to do with the CLA companies.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import user_whitelisted, get_company_instance, get_user_instance

class CompanyTestCase(CLATestCase):
    """Company test cases."""
    def test_get_companies(self):
        """Tests for getting all companies."""
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(response.data, [])
        self.create_company()
        self.create_company()
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(len(response.data), 2)

    def test_get_companies(self):
        """Tests for getting individual companies."""
        response = hug.test.get(cla.routes, '/v1/company/1')
        self.assertEqual(response.data, {'errors': {'company_id': 'Company not found'}})
        company = self.create_company()
        response = hug.test.get(cla.routes, '/v1/company/' + company['company_id'])
        self.assertEqual(response.data['company_id'], company['company_id'])

    def test_post_companies(self):
        """Tests for creating companies."""
        response = self.create_company()
        self.assertTrue('errors' not in response)
        data = {'company_name': 'Org Name',
                'company_whitelist': [],
                'company_whitelist_patterns': []}
        missing_name = data.copy()
        del missing_name['company_name']
        response = hug.test.post(cla.routes, '/v1/company', missing_name)
        self.assertEqual(response.data, {'errors': {'company_name': "Required parameter 'company_name' not supplied"}})
        missing_whitelist = data.copy()
        del missing_whitelist['company_whitelist']
        response = hug.test.post(cla.routes, '/v1/company', missing_whitelist)
        self.assertEqual(response.data, {'errors': {'company_whitelist': "Required parameter 'company_whitelist' not supplied"}})

    def test_put_companies(self):
        """Tests for updating companies."""
        company = self.create_company()
        response = hug.test.put(cla.routes, '/v1/company', {'meh': 'test'})
        self.assertEqual(response.data, {'errors': {'company_id': "Required parameter 'company_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/company', {'company_id': str(uuid.uuid4()), 'meh': 'test'})
        self.assertEqual(response.data, {'errors': {'company_id': 'Company not found'}})
        response = hug.test.put(cla.routes, '/v1/company', {'company_name': 'New Org Name'})
        self.assertEqual(response.data, {'errors': {'company_id': "Required parameter 'company_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/company', {'company_id': company['company_id'],
                                                             'company_name': 'New Org Name'})
        self.assertEqual(response.data['company_name'], 'New Org Name')
        response = hug.test.put(cla.routes, '/v1/company', {'company_id': company['company_id'],
                                                             'company_whitelist': ['user@very-safe.org', 'another-user@super-safe.org']})
        self.assertTrue('user@very-safe.org' in response.data['company_whitelist'])
        self.assertTrue('another-user@super-safe.org' in response.data['company_whitelist'])
        response = hug.test.put(cla.routes, '/v1/company', {'company_id': company['company_id'],
                                                             'company_whitelist_patterns': ['@ibm.com', '@info.ibm.*']})
        self.assertTrue('@ibm.com' in response.data['company_whitelist_patterns'])
        self.assertTrue('@info.ibm.*' in response.data['company_whitelist_patterns'])

    def test_delete_companies(self):
        """Tests for deleting companies."""
        company = self.create_company()
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertTrue(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/company/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'company_id': 'Company not found'}})
        response = hug.test.delete(cla.routes, '/v1/company/' + company['company_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(len(response.data), 0)

    def test_whitelist(self):
        """Test for company whitelists."""
        user = get_user_instance()
        whitelist = ['test@test.com']
        patterns = ['*@ibm.com', 'info@*.ibm.co.uk', 'some@other.email']
        company_data = self.create_company(company_whitelist=whitelist,
                                           company_whitelist_patterns=patterns)
        company = get_company_instance()
        company.load(company_data['company_id'])
        user.model.user_emails = ['test@test.com']
        self.assertTrue(user_whitelisted(user, company))
        user.model.user_emails = ['test2@test.com']
        self.assertFalse(user_whitelisted(user, company))
        user.model.user_emails = ['test@ibm.com']
        self.assertTrue(user_whitelisted(user, company))
        user.model.user_emails = ['test@ibm.ca']
        self.assertFalse(user_whitelisted(user, company))
        user.model.user_emails = ['info@test.ibm.co.uk']
        self.assertTrue(user_whitelisted(user, company))
        user.model.user_emails = ['some@other.email']
        self.assertTrue(user_whitelisted(user, company))
        user.model.user_emails = ['Some@other.email']
        self.assertFalse(user_whitelisted(user, company))
        user.model.user_emails = ['some@other.email.com']
        self.assertFalse(user_whitelisted(user, company))

if __name__ == '__main__':
    unittest.main()
