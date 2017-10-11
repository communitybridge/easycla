PROJECT_EXTERNAL_ID = 'salesforce-id-here'
TEST_DOCUMENT_URL = 'https://github.com/cncf/cla/raw/master/individual-cla.pdf'
GITHUB_INSTALLATION_ID = 49309 # NOT THE APP ID - find it in the webhook request JSON or URL when viewing installed apps.

import sys
sys.path.append('../')

import cla
import uuid
import base64
import urllib.request
from cla.utils import get_document_instance, get_github_organization_instance, get_project_instance

# Organisation
github_org = get_github_organization_instance().get_organization_by_installation_id(GITHUB_INSTALLATION_ID)
# Project
github_project = get_project_instance().get_project_by_external_id(PROJECT_EXTERNAL_ID)
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
cla.log.info('Adding CLA document to project: %s', TEST_DOCUMENT_URL)
github_project.add_project_individual_document(document)
github_project.save()
