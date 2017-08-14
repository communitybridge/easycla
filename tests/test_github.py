"""
Tests having to do with the GitHub repository service provider.
"""

import unittest
import falcon

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_user_instance, get_agreement_instance, get_repository_instance
from cla.models.github_models import MockGitHub, MockGitHubPullRequest, \
                                     get_pull_request_commit_authors

class GitHubTestCase(CLATestCase):
    """GitHub test cases."""
    def test_received_activity(self):
        """Tests for requesting signature."""
        data = {}
        github = MockGitHub()
        # Test non-pull request.
        result = github.received_activity(data)
        self.assertEqual(result, {'message': 'Not a pull request - no action performed'})
        # Test closed PR.
        data['pull_request'] = True
        data['action'] = 'closed'
        result = github.received_activity(data)
        self.assertEqual(result, None)

    def test_sign_request(self):
        """Tests for a request to sign a document."""
        github = MockGitHub(oauth2_token=True)
        github.initialize({'GITHUB_USERNAME': 'username', 'GITHUB_TOKEN': 'token'})
        request = None
        project = self.create_project()
        self.create_document(project['project_id'])
        repository = self.create_repository(project['project_id'])
        change_request_id = 1
        with self.assertRaises(falcon.redirects.HTTPFound) as context:
            github.sign_request(repository['repository_id'], change_request_id, request)
        self.assertTrue('http://signing-service.com/send-user-here' in str(context.exception))

    def test_oauth2_redirect(self):
        """Tests for the OAuth2 redirect from our repository provider."""
        github = MockGitHub(oauth2_token=True)
        github.initialize({'GITHUB_USERNAME': 'username', 'GITHUB_TOKEN': 'token'})
        state = 'wrong-state'
        code = 'code-here'
        project = self.create_project()
        self.create_document(project['project_id'])
        repository = self.create_repository(project['project_id'])
        change_request_id = 1
        request = None
        with self.assertRaises(falcon.HTTPBadRequest) as context:
            github.oauth2_redirect(state, code, repository['repository_id'],
                                   change_request_id, request)
        self.assertTrue('Invalid OAuth2 state' in str(context.exception))
        state = 'random-state'
        with self.assertRaises(falcon.HTTPFound) as context:
            github.oauth2_redirect(state, code, repository['repository_id'],
                                   change_request_id, request)
        self.assertTrue('http://signing-service.com/send-user-here' in str(context.exception))

    def test_get_pull_request_commit_authors(self): # pylint: disable=invalid-name
        """Tests the get_pull_request_commit_authors() function."""
        pull_request = MockGitHubPullRequest(1)
        authors = get_pull_request_commit_authors(pull_request)
        self.assertEqual(len(authors), 1)
        commit = authors[0][0]
        author = authors[0][1]
        self.assertEqual(commit.sha, 'sha-test-commit')
        self.assertEqual(author.email, 'user@github.com')

    def test_update_change_request(self):
        """Tests for the update_change_request method."""
        github = MockGitHub(oauth2_token=True)
        github.initialize({'GITHUB_USERNAME': 'username', 'GITHUB_TOKEN': 'token'})
        user_data = self.create_user(user_email='user@github.com', user_github_id=1)
        project = self.create_project()
        repo_data = self.create_repository(project['project_id'])
        repository = get_repository_instance()
        repository.load(repo_data['repository_id'])
        # Update with a signed and approved agreement.
        agreement_data = self.create_agreement(agreement_project_id=project['project_id'],
                                               agreement_reference_id=user_data['user_id'],
                                               agreement_reference_type='user')
        change_request_id = 1
        github.update_change_request(repository, change_request_id)
        # Update with a signed agreement, not approved.
        agreement_data = self.create_agreement(agreement_project_id=project['project_id'],
                                               agreement_reference_id=user_data['user_id'],
                                               agreement_reference_type='user',
                                               agreement_approved=False)
        change_request_id = 1
        github.update_change_request(repository, change_request_id)

    def test_get_or_create_user(self):
        """Tests for the get_or_create_user() method."""
        github = MockGitHub(oauth2_token=True)
        github.initialize({'GITHUB_USERNAME': 'username', 'GITHUB_TOKEN': 'token'})
        # All other tests assume missing user. Check for user found here.
        user_data = self.create_user(user_email='test@user.com', user_github_id=1)
        github.get_or_create_user(None)
        # Ensure the GitHub ID was updated to the 123 that the mock data uses.
        user = get_user_instance()
        user.load(user_data['user_id'])
        self.assertEqual(user.get_user_github_id(), 123)

if __name__ == '__main__':
    unittest.main()
