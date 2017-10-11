"""
Convenience script to create a user and signature in the CLA system.
"""
import uuid
import sys
sys.path.append('../')

import cla
from cla.utils import get_signature_instance, get_user_instance, get_project_instance

PROJECT_EXTERNAL_ID = 'salesforce-id-here'
USER_GITHUB_ID = 123

user = get_user_instance().get_user_by_github_id(USER_GITHUB_ID)
project = get_project_instance().get_project_by_external_id(PROJECT_EXTERNAL_ID)

# Test Agreement.
cla.log.info('Creating CLA signature for user %s and project %s' %(user.get_user_name(), project.get_project_external_id()))
signature = get_signature_instance()
signature.set_signature_id(str(uuid.uuid4()))
signature.set_signature_project_id(project.get_project_id())
signature.set_signature_signed(True)
signature.set_signature_approved(True)
signature.set_signature_type('cla')
signature.set_signature_reference_id(user.get_user_id())
signature.set_signature_reference_type('user')
signature.set_signature_document_major_version(1)
signature.set_signature_document_minor_version(0)
signature.save()
