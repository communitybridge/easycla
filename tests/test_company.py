"""
Tests having to do with the companies.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import user_whitelisted, get_company_instance, get_user_instance

class CompanyTestCase(CLATestCase):
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
