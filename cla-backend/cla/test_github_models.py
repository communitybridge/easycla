# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import unittest

import cla
import github
from cla import config
from cla.models.github_models import get_pull_request_commit_authors
from github import Github


class TestGitHubModels(unittest.TestCase):

    def setUp(self):
        self.assertTrue(cla.conf['GITHUB_OAUTH_TOKEN'] != '',
                        'Missing GITHUB_OAUTH_TOKEN environment variable - required to run unit tests')
        cla.log.debug('Using GITHUB_OAUTH_TOKEN: {}...'.format(cla.conf['GITHUB_OAUTH_TOKEN'][:5]))

    def test_commit_authors_with_named_user(self):
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

    def test_commit_authors_no_named_user(self):
        """
        Test that we can load commit authors from a pull request that does NOT have the traditional
        github.NamedUser.NamedUser object filled out
        """
        g = Github(cla.conf['GITHUB_OAUTH_TOKEN'])
        repo = g.get_repo(27729926)  # grpc/grpc-java
        pr = repo.get_pull(6152)  # example: https://github.com/grpc/grpc-java/pull/6152
        cla.log.info("Retrieved GitHub PR: {}".format(pr))
        commits = pr.get_comments()
        cla.log.info("Retrieved GitHub PR: {}, commits: {}".format(pr, commits))

        # Returns a list tuples, which look like (commit_sha_string, (author_id, author_username, author_email),
        # which, as you can see, the second element of the tuple is another tuple containing the author information
        commit_authors = get_pull_request_commit_authors(pr)
        # cla.log.info("Result: {}".format(commit_authors))
        # cla.log.info([author_info[1] for commit, author_info in commit_authors])
        self.assertTrue('snalkar' in [author_info[1] for commit, author_info in commit_authors])


if __name__ == '__main__':
    unittest.main()
