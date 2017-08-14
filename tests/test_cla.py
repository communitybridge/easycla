"""
Base test case for CLA system.
"""
# Ensure the CLA app is in our path for tests.
import os
import sys
TESTDIR = os.path.abspath(os.path.dirname(__file__))
sys.path.insert(0, TESTDIR + '/..')

# Ensure no debug messages from the boto package (and others) don't appear.
import logging
logging.disable(logging.ERROR)

import unittest
import uuid
import hug
import cla

# For now, force DynamoDB for testing (should probably be SQlite :memory: in the future).
cla.conf['DATABASE'] = 'DynamoDB'
cla.conf['DATABASE_HOST'] = 'http://localhost:8000'
# Use a mock signing service for testing.
cla.conf['SIGNING_SERVICE'] = 'MockDocuSign'
cla.conf['BASE_URL'] = 'http://cla-system.com'
# Use a mock email service for testing.
cla.conf['EMAIL_SERVICE'] = 'MockSMTP'
# Use exclusively local storage for tests.
cla.conf['STORAGE_SERVICE'] = 'LocalStorage'

from cla import routes
from cla.utils import create_database, delete_database

class CLATestCase(unittest.TestCase):
    def setUp(self):
        """Setup method executes before every test."""
        create_database()

    def create_user(self,
                    user_email='user@email.com',
                    user_name='User Name',
                    user_organization_id=None,
                    user_github_id=12345):
        """Helper method to create a user."""
        data = {'user_email': user_email,
                'user_name': user_name,
                'user_github_id': user_github_id}
        if user_organization_id is not None:
            data['user_organization_id'] = user_organization_id
        response = hug.test.post(routes, '/v1/user', data)
        return response.data

    def create_agreement(self,
                         agreement_project_id,
                         agreement_reference_id,
                         agreement_reference_type,
                         agreement_type='cla',
                         agreement_signed=True,
                         agreement_approved=True,
                         agreement_return_url='http://test-github.com/user/repo/1',
                         agreement_sign_url='http://link-to-agreement.com/sign-here'):
        """Helper method to create an agreements."""
        data = {'agreement_project_id': agreement_project_id,
                'agreement_reference_id': agreement_reference_id,
                'agreement_reference_type': agreement_reference_type,
                'agreement_type': agreement_type,
                'agreement_signed': agreement_signed,
                'agreement_approved': agreement_approved,
                'agreement_return_url': agreement_return_url,
                'agreement_sign_url': agreement_sign_url}
        response = hug.test.post(routes, '/v1/agreement', data)
        return response.data

    def create_repository(self, # pylint: disable=too-many-arguments
                          repository_project_id,
                          repository_name='Repo Name',
                          repository_type='mock_github',
                          repository_url='https://some-github-url.com/repo-name',
                          repository_external_id=1):
        """Helper method to create a repository."""
        data = {'repository_project_id': repository_project_id,
                'repository_name': repository_name,
                'repository_type': repository_type,
                'repository_url': repository_url,
                'repository_external_id': repository_external_id}
        response = hug.test.post(routes, '/v1/repository', data)
        return response.data

    def create_organization(self,
                            organization_name='Org Name',
                            organization_whitelist=['safe.org'],
                            organization_exclude_patterns=['^info@.*']):
        """Helper method to create organizations."""
        data = {'organization_name': organization_name,
                'organization_whitelist': organization_whitelist,
                'organization_exclude_patterns': organization_exclude_patterns}
        response = hug.test.post(routes, '/v1/organization', data)
        return response.data

    def create_document(self, project_id, document_type='individual',
                        document_name='doc_name.pdf',
                        document_content_type='url+pdf',
                        document_content='http://url.com/document.pdf'):
        """Helper method to create a document."""
        data = {'document_name': document_name,
                'document_content_type': document_content_type,
                'document_content': document_content}
        response = hug.test.post(routes,
                                 '/v1/project/' + project_id +
                                 '/document/' + document_type, data)
        return response.data

    def create_project(self, project_id=None, project_name='Project Name'):
        """Helper method to create a project."""
        if project_id is None:
            project_id = str(uuid.uuid4())
        data = {'project_id': project_id, 'project_name': project_name}
        response = hug.test.post(routes, '/v1/project', data)
        return response.data

    def tearDown(self):
        """Tear down method executes after each test."""
        delete_database()
