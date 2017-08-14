"""
Tests having to do with the CLA organizations.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class OrganizationTestCase(CLATestCase):
    """Organization test cases."""
    def test_get_organizations(self):
        """Tests for getting all organizations."""
        response = hug.test.get(cla.routes, '/v1/organization')
        self.assertEqual(response.data, [])
        self.create_organization()
        self.create_organization()
        response = hug.test.get(cla.routes, '/v1/organization')
        self.assertEqual(len(response.data), 2)

    def test_get_organizations(self):
        """Tests for getting individual organizations."""
        response = hug.test.get(cla.routes, '/v1/organization/1')
        self.assertEqual(response.data, {'errors': {'organization_id': 'Organization not found'}})
        organization = self.create_organization()
        response = hug.test.get(cla.routes, '/v1/organization/' + organization['organization_id'])
        self.assertEqual(response.data['organization_id'], organization['organization_id'])

    def test_post_organizations(self):
        """Tests for creating organizations."""
        response = self.create_organization()
        self.assertTrue('errors' not in response)
        data = {'organization_name': 'Org Name',
                'organization_whitelist': [],
                'organization_exclude_patterns': []}
        missing_name = data.copy()
        del missing_name['organization_name']
        response = hug.test.post(cla.routes, '/v1/organization', missing_name)
        self.assertEqual(response.data, {'errors': {'organization_name': "Required parameter 'organization_name' not supplied"}})
        missing_whitelist = data.copy()
        del missing_whitelist['organization_whitelist']
        response = hug.test.post(cla.routes, '/v1/organization', missing_whitelist)
        self.assertEqual(response.data, {'errors': {'organization_whitelist': "Required parameter 'organization_whitelist' not supplied"}})
        missing_exclude_patterns = data.copy()
        del missing_exclude_patterns['organization_exclude_patterns']
        response = hug.test.post(cla.routes, '/v1/organization', missing_exclude_patterns)
        self.assertEqual(response.data, {'errors': {'organization_exclude_patterns': "Required parameter 'organization_exclude_patterns' not supplied"}})

    def test_put_organizations(self):
        """Tests for updating organizations."""
        organization = self.create_organization()
        response = hug.test.put(cla.routes, '/v1/organization', {'meh': 'test'})
        self.assertEqual(response.data, {'errors': {'organization_id': "Required parameter 'organization_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/organization', {'organization_id': str(uuid.uuid4()), 'meh': 'test'})
        self.assertEqual(response.data, {'errors': {'organization_id': 'Organization not found'}})
        response = hug.test.put(cla.routes, '/v1/organization', {'organization_name': 'New Org Name'})
        self.assertEqual(response.data, {'errors': {'organization_id': "Required parameter 'organization_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/organization', {'organization_id': organization['organization_id'],
                                                             'organization_name': 'New Org Name'})
        self.assertEqual(response.data['organization_name'], 'New Org Name')
        response = hug.test.put(cla.routes, '/v1/organization', {'organization_id': organization['organization_id'],
                                                             'organization_whitelist': ['@very-safe.org', '@super-safe.org']})
        self.assertTrue('@very-safe.org' in response.data['organization_whitelist'])
        self.assertTrue('@super-safe.org' in response.data['organization_whitelist'])
        response = hug.test.put(cla.routes, '/v1/organization', {'organization_id': organization['organization_id'],
                                                             'organization_exclude_patterns': ['^admin@*', '^info@*']})
        self.assertTrue('^admin@*' in response.data['organization_exclude_patterns'])
        self.assertTrue('^info@*' in response.data['organization_exclude_patterns'])

    def test_delete_organizations(self):
        """Tests for deleting organizations."""
        organization = self.create_organization()
        response = hug.test.get(cla.routes, '/v1/organization')
        self.assertTrue(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/organization/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'organization_id': 'Organization not found'}})
        response = hug.test.delete(cla.routes, '/v1/organization/' + organization['organization_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/organization')
        self.assertEqual(len(response.data), 0)

if __name__ == '__main__':
    unittest.main()
