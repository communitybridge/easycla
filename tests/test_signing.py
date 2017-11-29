"""
Tests having to do with the signature endpoints.
"""

import unittest
import uuid
import hug

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_signature_instance

class SigningTestCase(CLATestCase):
    """Signing test cases."""
    def test_post_signed(self):
        """Tests for signing service callbacks."""
        # Will need to test all of the supported repository types.
        project = self.create_project()
        repo = self.create_repository(project['project_id'])
        user = self.create_user()
        self.create_document(project['project_id'])
        signature = self.create_signature(project['project_id'],
                                          user['user_id'],
                                          'user',
                                          signature_signed=False)
        change_id = 1 # Repository provider specific ID.
        # First one has status 'Sent', second one has status 'Completed'.
        fhandle = open('resources/docusign_callback_payload.xml')
        docusign_payload = fhandle.read()
        fhandle.close()
        data = docusign_payload %signature['signature_id']
        installation_id = 123
        signed_route = '/v1/signed/individual/%s/%s/%s' %(installation_id,
                                                          repo['repository_external_id'],
                                                          change_id)
        # Currently can't test this as we're too GitHub specific.
        # See the update_repository_provider function in docusign_models.py
        #response = hug.test.post(cla.routes, signed_route, data)
        #agr = get_signature_instance()
        #agr.load(signature['signature_id'])
        #self.assertTrue(agr.get_signature_signed())

    def test_return_url(self):
        """Tests for the user return URL after signing."""
        # Will need to test other repo types as well.
        project = self.create_project()
        repo = self.create_repository(project['project_id'])
        user = self.create_user()
        self.create_document(project['project_id'])
        signature = self.create_signature(project['project_id'], user['user_id'], 'user',
                                          signature_return_url='http://github.com/user/repo/1')
        url = '/v1/return-url/' + signature['signature_id']
        response = hug.test.get(cla.routes, url)
        self.assertEqual(response.status, '302 Found')
        self.assertEqual(response.headers_dict['location'], 'http://github.com/user/repo/1')

    def test_request_signature(self):
        """Tests for the request signature endpoint."""
        response = hug.test.post(cla.routes, '/v1/request-individual-signature', {'bad-data'})
        self.assertEqual(response.data, {
            'errors': {
                'project_id': "Required parameter 'project_id' not supplied",
                'user_id': "Required parameter 'user_id' not supplied"}})
        project = self.create_project()
        document = self.create_document(project['project_id'])
        user = self.create_user()
        repository = self.create_repository(project['project_id'])
        response = hug.test.post(cla.routes,
                                 '/v1/request-individual-signature',
                                 {'project_id': project['project_id'],
                                  'user_id': user['user_id'],
                                  'return_url': 'http://return.url/here'})
        signature_id = response.data['signature_id']
        self.assertEqual(response.data,  {'signature_id': signature_id,
                                          'project_id': project['project_id'],
                                          'sign_url': 'http://signing-service.com/send-user-here',
                                          'user_id': user['user_id']})

if __name__ == '__main__':
    unittest.main()
