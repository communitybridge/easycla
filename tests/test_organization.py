"""
Tests having to do with the CLA companys.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class CompanyTestCase(CLATestCase):
    """Company test cases."""
    def test_get_companys(self):
        """Tests for getting all companys."""
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(response.data, [])
        self.create_company()
        self.create_company()
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(len(response.data), 2)

    def test_get_companys(self):
        """Tests for getting individual companys."""
        response = hug.test.get(cla.routes, '/v1/company/1')
        self.assertEqual(response.data, {'errors': {'company_id': 'Company not found'}})
        company = self.create_company()
        response = hug.test.get(cla.routes, '/v1/company/' + company['company_id'])
        self.assertEqual(response.data['company_id'], company['company_id'])

    def test_post_companys(self):
        """Tests for creating companys."""
        response = self.create_company()
        self.assertTrue('errors' not in response)
        data = {'company_name': 'Org Name',
                'company_whitelist': [],
                'company_exclude_patterns': []}
        missing_name = data.copy()
        del missing_name['company_name']
        response = hug.test.post(cla.routes, '/v1/company', missing_name)
        self.assertEqual(response.data, {'errors': {'company_name': "Required parameter 'company_name' not supplied"}})
        missing_whitelist = data.copy()
        del missing_whitelist['company_whitelist']
        response = hug.test.post(cla.routes, '/v1/company', missing_whitelist)
        self.assertEqual(response.data, {'errors': {'company_whitelist': "Required parameter 'company_whitelist' not supplied"}})
        missing_exclude_patterns = data.copy()
        del missing_exclude_patterns['company_exclude_patterns']
        response = hug.test.post(cla.routes, '/v1/company', missing_exclude_patterns)
        self.assertEqual(response.data, {'errors': {'company_exclude_patterns': "Required parameter 'company_exclude_patterns' not supplied"}})

    def test_put_companys(self):
        """Tests for updating companys."""
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
                                                             'company_whitelist': ['@very-safe.org', '@super-safe.org']})
        self.assertTrue('@very-safe.org' in response.data['company_whitelist'])
        self.assertTrue('@super-safe.org' in response.data['company_whitelist'])
        response = hug.test.put(cla.routes, '/v1/company', {'company_id': company['company_id'],
                                                             'company_exclude_patterns': ['^admin@*', '^info@*']})
        self.assertTrue('^admin@*' in response.data['company_exclude_patterns'])
        self.assertTrue('^info@*' in response.data['company_exclude_patterns'])

    def test_delete_companys(self):
        """Tests for deleting companys."""
        company = self.create_company()
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertTrue(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/company/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'company_id': 'Company not found'}})
        response = hug.test.delete(cla.routes, '/v1/company/' + company['company_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/company')
        self.assertEqual(len(response.data), 0)

if __name__ == '__main__':
    unittest.main()
