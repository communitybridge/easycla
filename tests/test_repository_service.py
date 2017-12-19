"""
Tests having to do with the repository service webhooks.
"""

import unittest
import hug
import json

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class RepositoryServiceTestCase(CLATestCase):
    """Repository service test cases."""
    def test_github_webhook(self):
        """Tests for the GitHub webhook."""
        # Opened PR.
        fhandle = open(cla.utils.get_cla_path() + '/tests/resources/github_open_pr.json')
        github_payload = json.load(fhandle)
        fhandle.close()
        response = hug.test.post(cla.routes, '/v1/repository-provider/mock_github/activity', github_payload)
        self.assertEqual(response.status, '200 OK')
        # TODO: Check if things were handled properly in DB.
        # Re-opened PR.
        fhandle = open(cla.utils.get_cla_path() + '/tests/resources/github_reopen_pr.json')
        github_payload = json.load(fhandle)
        fhandle.close()
        response = hug.test.post(cla.routes, '/v1/repository-provider/mock_github/activity', github_payload)
        self.assertEqual(response.status, '200 OK')
        # TODO: Check if things were handled properly in DB.
        project = self.create_project()
        repository = self.create_repository(project['project_id'])
        change_request_id = 1
        url = '/v1/repository-provider/mock_github/sign/999/' + repository['repository_id'] + '/' + str(change_request_id)
        response = hug.test.get(cla.routes, url)
        self.assertEqual(response.status, '302 Found')
        self.assertEqual(response.headers_dict['location'], 'http://authorization.url')
        # TODO: Check if things were handled properly in DB.

if __name__ == '__main__':
    unittest.main()
