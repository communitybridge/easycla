"""
Tests having to do with the CLA documents.
"""

import os
import shutil
import hug
import unittest

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.utils import get_document_instance, get_project_instance

class LocalStorageTestCase(CLATestCase):
    """Local storage test cases."""
    def test_local_storage(self):
        """Tests for local storage of documents."""
        local_storage_folder = cla.conf['LOCAL_STORAGE_FOLDER'] + '/'
        # Clear the previous files in case of failed tests.
        try: shutil.rmtree(local_storage_folder)
        except: pass
        os.makedirs(local_storage_folder)
        project = self.create_project()
        with open(cla.utils.get_cla_path() + '/tests/resources/test.pdf', 'rb') as fhandle:
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
        doc.set_document_major_version(document['document_major_version'])
        doc.set_document_minor_version(document['document_minor_version'])
        doc.set_document_file_id(document['document_file_id'])
        proj = get_project_instance()
        proj.load(project['project_id'])
        self.assertTrue(len(proj.get_project_individual_documents()) == 1)
        proj.remove_project_individual_document(doc)
        self.assertTrue(len(proj.get_project_individual_documents()) == 0)
        self.assertFalse(os.path.exists(local_storage_folder + filename))
        self.assertTrue(os.listdir(local_storage_folder) == [])
        os.rmdir(local_storage_folder)

if __name__ == '__main__':
    unittest.main()
