"""
Convenience script to drop/re-create the database.

Based on the database storage configuration found in config.py and config_instance.py.
"""

# Temporary for testing purposes - create the Project and Document objects.
TEST_DOCUMENT_URL = 'https://github.com/cncf/cla/raw/master/individual-cla.pdf'
# The GitHub user/org used for testing purposes.
GITHUB_PROJECT_NAME = 'cla-test'

import sys
sys.path.append('../')

import cla
from cla.models.github_models import GitHub
from cla.models.gitlab_models import GitLab
from cla.utils import create_database, delete_database
delete_database()
create_database()

import uuid
import base64
import urllib.request
from cla.utils import get_project_instance, get_document_instance
## Project
cla.log.info('Creating new project: %s', 'Test Project')
github_project = get_project_instance()
github_project.set_project_id((uuid.uuid4()))
github_project.set_project_name('Test GitHub Project')
github_project.set_project_external_id(GITHUB_PROJECT_NAME)
github_project.set_project_ccla_requires_icla_signature(False)
github_project.save()
# Document
# Slower as the document is fetched every time a document signature is initiated.
#document = Document(str(uuid.uuid4()), 'Test Document', 'url+pdf', TEST_DOCUMENT_URL)
resource = urllib.request.urlopen(TEST_DOCUMENT_URL)
data = base64.b64encode(resource.read()) # Document expects base64 encoded data.
document = get_document_instance()
document.set_document_name('Test Document')
document.set_document_file_id(str(uuid.uuid4()))
document.set_document_content_type('storage+pdf')
document.set_document_content(data)
document.set_document_major_version(1)
document.set_document_minor_version(0)
cla.log.info('Adding CLA document to new project: %s', TEST_DOCUMENT_URL)
github_project.add_project_individual_document(document)
github_project.save()
cla.log.info('Done')
