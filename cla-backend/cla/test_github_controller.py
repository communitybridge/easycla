# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import unittest

from cla.controllers.github import get_org_name_from_event


class TestGitHubController(unittest.TestCase):

    @classmethod
    def setUpClass(cls) -> None:
        pass

    @classmethod
    def tearDownClass(cls) -> None:
        pass

    def setUp(self) -> None:
        pass

    def tearDown(self) -> None:
        pass

    def test_get_org_name_from_event(self) -> None:
        # Webhook event payload
        # see: https://developer.github.com/v3/activity/events/types/#webhook-payload-example-12
        # body['installation']['account']['login']
        body = {
            'action': 'created',
            'installation': {
                'id': 2,
                'account': {
                    'login': 'Linux Foundation',
                    'id': 1,
                    'node_id': 'MDQ6VXNlcjE=',
                    'avatar_url': 'https://github.com/images/error/octocat_happy.gif',
                    'gravatar_id': '',
                    'url': 'https://api.github.com/users/octocat',
                    'html_url': 'https://github.com/octocat',
                    'followers_url': 'https://api.github.com/users/octocat/followers',
                    'following_url': 'https://api.github.com/users/octocat/following{/other_user}',
                    'gists_url': 'https://api.github.com/users/octocat/gists{/gist_id}',
                    'starred_url': 'https://api.github.com/users/octocat/starred{/owner}{/repo}',
                    'subscriptions_url': 'https://api.github.com/users/octocat/subscriptions',
                    'organizations_url': 'https://api.github.com/users/octocat/orgs',
                    'repos_url': 'https://api.github.com/users/octocat/repos',
                    'events_url': 'https://api.github.com/users/octocat/events{/privacy}',
                    'received_events_url': 'https://api.github.com/users/octocat/received_events',
                    'type': 'User',
                    'site_admin': False
                },
                'repository_selection': 'selected',
                'access_tokens_url': 'https://api.github.com/installations/2/access_tokens',
                'repositories_url': 'https://api.github.com/installation/repositories',
                'html_url': 'https://github.com/settings/installations/2',
                'app_id': 5725,
                'target_id': 3880403,
                'target_type': 'User',
                'permissions': {
                    'metadata': 'read',
                    'contents': 'read',
                    'issues': 'write'
                },
                'events': [
                    'push',
                    'pull_request'
                ],
                'created_at': 1525109898,
                'updated_at': 1525109899,
                'single_file_name': 'config.yml'
            }
        }
        self.assertEqual('Linux Foundation', get_org_name_from_event(body), 'GitHub Org Matches')

    def test_get_org_name_from_event_empty(self) -> None:
        self.assertIsNone(get_org_name_from_event({}), 'GitHub Org Does Not Match')


if __name__ == '__main__':
    unittest.main()
