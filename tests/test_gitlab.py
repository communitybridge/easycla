"""
Tests having to do with the GitLab repository service provider.
"""

import unittest
import falcon

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_repository_instance
from cla.models.gitlab_models import MockGitLab

class GitLabTestCase(CLATestCase):
    """GitLab test cases."""
    def test_received_activity(self):
        """Tests for requesting signature."""
        return # Disable tests for GitLab until we support it again.
        data = {}
        gitlab = MockGitLab()
        # Test non-merge request.
        result = gitlab.received_activity(data)
        self.assertEqual(result, {'message': 'Not a merge request - no action performed'})
        # Test closed PR.
        data['object_kind'] = 'merge_request'
        data['object_attributes'] = {'action': 'close'}
        result = gitlab.received_activity(data)
        self.assertEqual(result, None)

    def test_sign_request(self):
        """Tests for a request to sign a document."""
        return # Disable tests for GitLab until we support it again.
        gitlab = MockGitLab(oauth2_token=True)
        gitlab.initialize({'GITLAB_DOMAIN': 'https://domain.name', 'GITLAB_TOKEN': 'token'})
        request = None
        project = self.create_project()
        self.create_document(project['project_id'])
        repository = self.create_repository(project['project_id'])
        change_request_id = 1
        with self.assertRaises(falcon.redirects.HTTPFound) as context:
            gitlab.sign_request(repository['repository_id'], change_request_id, request)
        self.assertTrue('http://signing-service.com/send-user-here' in str(context.exception))

    def test_oauth2_redirect(self):
        """Tests for the OAuth2 redirect from our repository provider."""
        return # Disable tests for GitLab until we support it again.
        gitlab = MockGitLab(oauth2_token=True)
        gitlab.initialize({'GITLAB_DOMAIN': 'https://domain.name', 'GITLAB_TOKEN': 'token'})
        state = 'wrong-state'
        code = 'code-here'
        project = self.create_project()
        self.create_document(project['project_id'])
        repository = self.create_repository(project['project_id'])
        change_request_id = 1
        request = None
        with self.assertRaises(falcon.HTTPBadRequest) as context:
            gitlab.oauth2_redirect(state, code, repository['repository_id'],
                                   change_request_id, request)
        self.assertTrue('Invalid OAuth2 state' in str(context.exception))
        state = 'random-state'
        with self.assertRaises(falcon.HTTPFound) as context:
            gitlab.oauth2_redirect(state, code, repository['repository_id'],
                                   change_request_id, request)
        self.assertTrue('http://signing-service.com/send-user-here' in str(context.exception))

    def test_update_change_request(self):
        """Tests for the update_change_request method."""
        return # Disable tests for GitLab until we support it again.
        gitlab = MockGitLab(oauth2_token=True)
        gitlab.initialize({'GITLAB_DOMAIN': 'https://domain.name', 'GITLAB_TOKEN': 'token'})
        user_data = self.create_user(user_email='user@gitlab.com')
        project = self.create_project()
        repo_data = self.create_repository(project['project_id'])
        repository = get_repository_instance()
        repository.load(repo_data['repository_id'])
        # Update with a signed and approved signature.
        signature_data = self.create_signature(signature_project_id=project['project_id'],
                                               signature_reference_id=user_data['user_id'],
                                               signature_reference_type='user')
        change_request_id = 1
        gitlab.update_change_request(repository, change_request_id)
        # Update with a signed signature, not approved.
        signature_data = self.create_signature(signature_project_id=project['project_id'],
                                               signature_reference_id=user_data['user_id'],
                                               signature_reference_type='user',
                                               signature_approved=False)
        change_request_id = 1
        gitlab.update_change_request(repository, change_request_id)

    def test_get_or_create_user(self):
        """Tests for the get_or_create_user() method."""
        return # Disable tests for GitLab until we support it again.
        gitlab = MockGitLab(oauth2_token=True)
        gitlab.initialize({'GITLAB_DOMAIN': 'https://domain.name', 'GITLAB_TOKEN': 'token'})
        # All other tests assume missing user. Check for user found here.
        self.create_user(user_email='test@user.com')
        user = gitlab.get_or_create_user(None)
        self.assertEqual(user.get_user_email(), 'test@user.com')

if __name__ == '__main__':
    unittest.main()
