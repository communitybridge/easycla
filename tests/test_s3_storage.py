"""
Tests having to do with the CLA documents.
"""

import hug
import unittest

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_document_instance, get_project_instance
# Ready the storage interface.
from cla.models.s3_storage import MockS3Storage as Storage
cla.conf['STORAGE_SERVICE'] = 'MockS3Storage'
# OR, if you want to test directly on S3:
#from cla.models.s3_storage import S3Storage as Storage
#cla.conf['STORAGE_SERVICE'] = 'S3Storage'

class S3StorageTestCase(CLATestCase):
    """S3 storage test cases."""
    def test_s3_storage(self):
        """Tests for S3 storage of documents."""
        project = self.create_project()
        with open('resources/test.pdf', 'rb') as fhandle:
            content = fhandle.read()
            ret = self.create_document(project['project_id'],
                                       document_content_type='storage+pdf',
                                       document_content=content)
        document = ret['project_individual_documents'][0]
        filename = document['document_file_id']
        path = '/v1/project/' + project['project_id'] + '/document/individual'
        response = hug.test.get(cla.routes, path)
        self.assertEqual(response.data, document)
        doc = get_document_instance()
        doc.set_document_content_type('storage+pdf')
        doc.set_document_major_version(document['document_major_version'])
        doc.set_document_minor_version(document['document_minor_version'])
        doc.set_document_file_id(document['document_file_id'])
        content = doc.get_document_content()
        with open('resources/test.pdf', 'rb') as fhandle:
            self.assertEqual(content, fhandle.read())
        proj = get_project_instance()
        proj.load(project['project_id'])
        self.assertTrue(len(proj.get_project_individual_documents()) == 1)
        # This should also delete the file from the S3 bucket.
        proj.remove_project_individual_document(doc)
        self.assertTrue(len(proj.get_project_individual_documents()) == 0)
        storage = Storage()
        storage.initialize(cla.conf)
        storage.retrieve(filename)

if __name__ == '__main__':
    unittest.main()
