# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import unittest

import cla


class TestGerritsModels(unittest.TestCase):

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

    def test_get_gerrit_by_project_id(self) -> None:
        """
        Test that we can get a gerrit by project id
        """
        pass
        # project_id = '6bba5291-5007-4aaf-abba-3b97875e2224'
        # response = gerrit.get_gerrit_by_project_id(project_id=project_id)
        # self.assertTrue(len(response) == 1, f'Found project id: {project_id}')


if __name__ == '__main__':
    unittest.main()
