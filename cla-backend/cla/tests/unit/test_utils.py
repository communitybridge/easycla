# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import unittest

import cla
from cla import utils
from cla.models.dynamo_models import User


class TestUtils(unittest.TestCase):
    tests_enabled = False

    @classmethod
    def setUpClass(cls) -> None:
        pass

    @classmethod
    def tearDownClass(cls) -> None:
        pass

    def setUp(self) -> None:
        # Only show critical logging stuff
        cla.log.level = logging.CRITICAL

    def tearDown(self) -> None:
        pass

    def test_user_get_user_by_username(self) -> None:
        """
        Test that we can get a user by username
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            user_instance = utils.get_user_instance()
            users = user_instance.get_user_by_username('ddeal')
            self.assertIsNotNone(users, 'User lookup by username is not None')
            self.assertEqual(len(users), 1, 'User lookup by username found 1')

            # some invalid username
            users = user_instance.get_user_by_username('foo')
            self.assertIsNone(users, 'User lookup by username is None')

    def test_user_get_user_by_email(self) -> None:
        """
        Test that we can get a user by email
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            users = User().get_user_by_email('ddeal@linuxfoundation.org')
            self.assertIsNotNone(users, 'User lookup by email is not None')
            self.assertEqual(len(users), 1, 'User lookup by email found 1')

            # some invalid email
            users = User().get_user_by_email('foo@bar.org')
            self.assertIsNone(users, 'User lookup by email is None')

    def test_user_get_user_by_github_id(self) -> None:
        """
        Test that we can get a user by github id
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            users = User().get_user_by_github_id(519609)
            self.assertIsNotNone(users, 'User lookup by github id is not None')
            self.assertEqual(len(users), 2, 'User lookup by github id found 2')

            # some invalid number
            users = User().get_user_by_github_id(9999999)
            self.assertIsNone(users, 'User lookup by github id is None')

    def test_user_get_user_by_github_username(self) -> None:
        """
        Test that we can get a user by github username
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            users = User().get_user_by_github_username('dealako')
            self.assertIsNotNone(users, 'User lookup by github username is not None')
            self.assertEqual(len(users), 1, 'User lookup by github username found 1')

            # some invalid username
            users = User().get_user_by_github_username('foooooo')
            self.assertIsNone(users, 'User lookup by github username is None')

    def test_lookup_user_github_username(self) -> None:
        """
        Test that we can lookup a given gihub user by id
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            self.assertEqual('dealako', cla.utils.lookup_user_github_username(519609), 'Found github username')
            # some invalid username
            self.assertIsNone(cla.utils.lookup_user_github_username(5196090000), 'None response from invalid github id')

    def test_lookup_user_github_id(self) -> None:
        """
        Test that we can lookup a given gihub id by the username
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            self.assertEqual(519609, cla.utils.lookup_user_github_id('dealako'), 'Found github id')
            # some invalid username
            self.assertIsNone(cla.utils.lookup_user_github_id('dealakooooooooo'),
                              'None response from invalid github username')

    def test_lookup_github_organizations(self) -> None:
        """
        Test that we can lookup a user's github organizations
        """
        # TODO - should use mock data - disable tests for now :-(
        if self.tests_enabled:
            organizations = cla.utils.lookup_github_organizations('dealako')
            self.assertEqual(3, len(organizations), 'Find github organizations')
            # some invalid username
            organizations = cla.utils.lookup_github_organizations('dealakooooooooo')
            self.assertTrue('error' in organizations, 'Find 0 github organizations')


if __name__ == '__main__':
    unittest.main()
