"""
Tests having to do with the companies.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import email_whitelisted, get_company_instance

class CompanyTestCase(CLATestCase):
    def test_whitelist(self):
        """Test for company whitelists."""
        whitelist = ['test@test.com']
        patterns = ['*@ibm.com', 'info@*.ibm.co.uk', 'some@other.email']
        company_data = self.create_company(company_whitelist=whitelist,
                                           company_whitelist_patterns=patterns)
        company = get_company_instance()
        company.load(company_data['company_id'])
        self.assertTrue(email_whitelisted('test@test.com', company))
        self.assertFalse(email_whitelisted('test2@test.com', company))
        self.assertTrue(email_whitelisted('test@ibm.com', company))
        self.assertFalse(email_whitelisted('test@ibm.ca', company))
        self.assertTrue(email_whitelisted('info@test.ibm.co.uk', company))
        self.assertTrue(email_whitelisted('some@other.email', company))
        self.assertFalse(email_whitelisted('Some@other.email', company))
        self.assertFalse(email_whitelisted('some@other.email.com', company))

if __name__ == '__main__':
    unittest.main()
