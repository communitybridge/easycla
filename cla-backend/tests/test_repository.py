"""
Tests having to do with the repository.
"""

import unittest
import uuid
import hug
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class RepositoryTestCase(CLATestCase):
    def test_get_repositories(self):
        """Test for getting a list of all repositories."""
        response = hug.test.get(cla.routes, '/v1/repository')
        self.assertEqual(response.data, [])
        project = self.create_project()
        project_id = project['project_id']
        self.create_repository(project_id)
        self.create_repository(project_id)
        self.create_repository(project_id)
        response = hug.test.get(cla.routes, '/v1/repository')
        self.assertEqual(len(response.data), 3)

    def test_get_repository(self):
        """Test for getting repository information."""
        response = hug.test.get(cla.routes, '/v1/repository/1')
        self.assertEqual(response.data, {'errors': {'repository_id': 'Repository not found'}})
        project = self.create_project()
        repo = self.create_repository(project['project_id'])
        response = hug.test.get(cla.routes, '/v1/repository/' + repo['repository_id'])
        self.assertEqual(response.data['repository_id'], repo['repository_id'])
        self.assertEqual(response.data['repository_name'], repo['repository_name'])

    def test_post_repository(self):
        """Test for creating a new repository."""
        response = hug.test.post(cla.routes, '/v1/repository',
                                 {'repository_name': 'Test Repo',
                                  'repository_type': 'mock_github',
                                  'repository_url': 'http://repo-url.com'})
        self.assertEqual(response.data, {'errors': {'repository_project_id': "Required parameter 'repository_project_id' not supplied"}})
        project = self.create_project()
        repo = self.create_repository(project['project_id'], repository_type='invalid')
        accepted_values = '|'.join(cla.utils.get_supported_repository_providers().keys())
        self.assertEqual(repo, {'errors': {'repository_type': 'Invalid value passed. The accepted values are: (%s)' %accepted_values}})
        repo = self.create_repository(project['project_id'], repository_url='invalid')
        self.assertEqual(repo, {'errors': {'repository_url': 'Invalid URL specified'}})

    def test_put_repository(self):
        """Test for updating a repository."""
        response = hug.test.put(cla.routes, '/v1/repository', {'repository_name': 'New Repo Name'})
        self.assertEqual(response.data, {'errors': {'repository_id': "Required parameter 'repository_id' not supplied"}})
        response = hug.test.put(cla.routes, '/v1/repository', {'repository_id': str(uuid.uuid4()), 'repository_name': 'New Repo Name'})
        self.assertEqual(response.data, {'errors': {'repository_id': 'Repository not found'}})
        project = self.create_project()
        repository = self.create_repository(project['project_id'])
        response = hug.test.put(cla.routes, '/v1/repository', {'repository_id': repository['repository_id'], 'repository_name': 'New Repo Name'})
        self.assertEqual(response.data['repository_name'], 'New Repo Name')
        response = hug.test.put(cla.routes, '/v1/repository', {'repository_id': repository['repository_id'], 'repository_url': 'invalid'})
        self.assertEqual(response.data, {'errors': {'repository_url': 'Invalid URL specified'}})
        accepted_values = '|'.join(cla.utils.get_supported_repository_providers().keys())
        response = hug.test.put(cla.routes, '/v1/repository', {'repository_id': repository['repository_id'], 'repository_type': 'invalid'})
        self.assertEqual(response.data, {'errors': {'repository_type': 'Invalid value passed. The accepted values are: (%s)' %accepted_values}})
        new_data = {'repository_id': repository['repository_id'],
                    'repository_name': 'New Repo Name Again',
                    'repository_url': 'http://new-repo-url.com',
                    'repository_type': 'github'}
        response = hug.test.put(cla.routes, '/v1/repository', new_data)
        self.assertEqual(response.data['repository_name'], 'New Repo Name Again')
        self.assertEqual(response.data['repository_url'], 'http://new-repo-url.com')
        self.assertEqual(response.data['repository_type'], 'github')

    def test_delete_repository(self):
        """Test deleting repository."""
        response = hug.test.delete(cla.routes, '/v1/repository/' + str(uuid.uuid4()))
        self.assertEqual(response.data, {'errors': {'repository_id': 'Repository not found'}})
        project = self.create_project()
        repository = self.create_repository(project['project_id'])
        response = hug.test.get(cla.routes, '/v1/repository')
        self.assertEqual(len(response.data), 1)
        response = hug.test.delete(cla.routes, '/v1/repository/' + repository['repository_id'])
        self.assertEqual(response.data, {'success': True})
        response = hug.test.get(cla.routes, '/v1/repository')
        self.assertEqual(len(response.data), 0)

if __name__ == '__main__':
    unittest.main()
