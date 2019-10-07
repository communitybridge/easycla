# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import unittest

import cla


class TestUserModels(unittest.TestCase):

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
        # # TODO - should use mock data
        # user_instance = utils.get_user_instance()
        # u = user_instance.get_user_by_username('ddeal')
        # self.assertIsNotNone(u, 'User lookup by username is not None')
        #
        # # some invalid username
        # u = user_instance.get_user_by_username('foo')
        # self.assertIsNone(u, 'User lookup by username is None')

    def test_user_get_user_by_email(self) -> None:
        """
        Test that we can get a user by email
        """
        # # TODO - should use mock data
        # u = User().get_user_by_email('ddeal@linuxfoundation.org')
        # self.assertIsNotNone(u, 'User lookup by email is not None')
        #
        # # some invalid email
        # u = User().get_user_by_email('foo@bar.org')
        # self.assertIsNone(u, 'User lookup by email is None')

    def test_user_get_user_by_github_id(self) -> None:
        """
        Test that we can get a user by github id
        """
        # # TODO - should use mock data
        # u = User().get_user_by_github_id(519609)
        # self.assertIsNotNone(u, 'User lookup by github id is not None')
        #
        # # some invalid number
        # u = User().get_user_by_github_id(9999999)
        # self.assertIsNone(u, 'User lookup by github id is None')

    def test_user_get_user_by_github_username(self) -> None:
        """
        Test that we can get a user by github username
        """
        # TODO - should use mock data
        # u = User().get_user_by_github_username('dealako')
        # self.assertIsNotNone(u, 'User lookup by github username is not None')
        #
        # # some invalid username
        # u = User().get_user_by_github_username('foooooo')
        # self.assertIsNone(u, 'User lookup by github username is None')


if __name__ == '__main__':
    unittest.main()
