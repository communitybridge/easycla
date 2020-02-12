# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import unittest
from unittest.mock import Mock, patch, MagicMock

from github import Github

import cla
from cla.models.github_models import get_pull_request_commit_authors, handle_commit_from_user
from cla.models.dynamo_models import Signature, Project


class TestGitHubModels(unittest.TestCase):

    @classmethod
    def setUpClass(cls) -> None:
        cls.mock_user_patcher = patch('cla.models.github_models.cla.utils.get_user_instance')
        cls.mock_signature_patcher = patch('cla.models.github_models.cla.utils.get_signature_instance')
        cls.mock_utils_patcher = patch('cla.models.github_models.cla.utils')
        cls.mock_utils_get = cls.mock_utils_patcher.start()
        cls.mock_user_get = cls.mock_user_patcher.start()
        cls.mock_signature_get = cls.mock_signature_patcher.start()

    @classmethod
    def tearDownClass(cls) -> None:
        cls.mock_user_patcher.stop()
        cls.mock_signature_patcher.stop()
        cls.mock_utils_patcher.stop()

    def setUp(self) -> None:
        # Only show critical logging stuff
        cla.log.level = logging.CRITICAL
        self.assertTrue(cla.conf['GITHUB_OAUTH_TOKEN'] != '',
                        'Missing GITHUB_OAUTH_TOKEN environment variable - required to run unit tests')
        # cla.log.debug('Using GITHUB_OAUTH_TOKEN: {}...'.format(cla.conf['GITHUB_OAUTH_TOKEN'][:5]))

    def tearDown(self) -> None:
        pass

    def test_commit_authors_with_named_user(self) -> None:
        """
        Test that we can load commit authors from a pull request that does have the traditional
        github.NamedUser.NamedUser object filled out
        """
        g = Github(cla.conf['GITHUB_OAUTH_TOKEN'])
        repo = g.get_repo(27729926)  # grpc/grpc-java
        pr = repo.get_pull(6142)  # example: https://github.com/grpc/grpc-java/pull/6142
        cla.log.info("Retrieved GitHub PR: {}".format(pr))
        commits = pr.get_comments()
        cla.log.info("Retrieved GitHub PR: {}, commits: {}".format(pr, commits))

        # Returns a list tuples, which look like (commit_sha_string, (author_id, author_username, author_email),
        # which, as you can see, the second element of the tuple is another tuple containing the author information
        commit_authors = get_pull_request_commit_authors(pr)
        # cla.log.info("Result: {}".format(commit_authors))
        # cla.log.info([author_info[1] for commit, author_info in commit_authors])
        self.assertTrue(4779759 in [author_info[0] for commit, author_info in commit_authors])

    def test_commit_authors_no_named_user(self) -> None:
        """
        Test that we can load commit authors from a pull request that does NOT have the traditional
        github.NamedUser.NamedUser object filled out
        """
        # We need to mock this service so that we can test our business logic - disabling this test for now
        # as they closed the PR
        g = Github(cla.conf['GITHUB_OAUTH_TOKEN'])
        repo = g.get_repo(27729926)  # grpc/grpc-java
        pr = repo.get_pull(6152)  # example: https://github.com/grpc/grpc-java/pull/6152
        cla.log.info("Retrieved GitHub PR: {}".format(pr))
        commits = pr.get_comments()
        cla.log.info("Retrieved GitHub PR: {}, commits: {}".format(pr, commits))

        # Returns a list tuples, which look like (commit_sha_string, (author_id, author_username, author_email),
        # which, as you can see, the second element of the tuple is another tuple containing the author information
        # commit_authors = get_pull_request_commit_authors(pr)
        # cla.log.info("Result: {}".format(commit_authors))
        # cla.log.info([author_info[1] for commit, author_info in commit_authors])
        # self.assertTrue('snalkar' in [author_info[1] for commit, author_info in commit_authors])

    def test_handle_commit_author_whitelisted(self) -> None:
        """
        Test case where commit authors have no signatures but have been whitelisted and should
        return missing list containing a whitelisted flag
        """
        # Mock user not existing and happens to be whitelisted
        self.mock_user_get.return_value.get_user_by_github_id.return_value = None
        self.mock_user_get.return_value.get_user_by_email.return_value = None
        self.mock_signature_get.return_value.get_signatures_by_project.return_value = [Signature()]
        self.mock_utils_get.return_value.is_whitelisted.return_value = True
        missing = []
        signed = []
        project = Project()
        project.set_project_id('fake_project_id')
        handle_commit_from_user(project, 'fake_sha', (123,'foo','foo@gmail.com'), signed, missing)
        self.assertListEqual(missing,[('fake_sha', [123, 'foo', 'foo@gmail.com', True])])
        self.assertEqual(signed, [])


if __name__ == '__main__':
    unittest.main()
