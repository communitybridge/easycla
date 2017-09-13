"""
Convenience script to drop/re-create the database.

Based on the database storage configuration found in config.py and config_instance.py.
"""

# Temporary for testing purposes - create the Project and Document objects.
TEST_DOCUMENT_URL = 'https://github.com/cncf/cla/raw/master/individual-cla.pdf'
# Project external ID.
PROJECT_EXTERNAL_ID = 'salesforce-id-here'
# The GitHub user/org used for testing purposes.
GITHUB_ORGANIZATION_NAME = 'linuxfoundation'
GITHUB_INSTALLATION_ID = 52353 # NOT THE APP ID - find it in the webhook request JSON

import sys
sys.path.append('../')

import cla
from cla.models.github_models import GitHub
from cla.utils import create_database, delete_database
delete_database()
create_database()

import uuid
import base64
import urllib.request
from cla.utils import get_project_instance, get_document_instance, get_github_organization_instance
## Project
cla.log.info('Creating new project for %s', PROJECT_EXTERNAL_ID)
github_project = get_project_instance()
github_project.set_project_id((uuid.uuid4()))
github_project.set_project_external_id(PROJECT_EXTERNAL_ID)
github_project.set_project_name('Test Project')
github_project.set_project_ccla_requires_icla_signature(False)
github_project.save()

# Organisation
github_org = get_github_organization_instance()
github_org.set_organization_name(GITHUB_ORGANIZATION_NAME)
github_org.set_organization_project_id(github_project.get_project_id())
github_org.set_organization_installation_id(GITHUB_INSTALLATION_ID)
github_org.save()

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
