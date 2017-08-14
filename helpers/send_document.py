"""
Convenience script to send a DocuSign document to the user's email address.

The user_id specified must exist in the database.
The docusign_document_id specified must exist in DocuSign.
"""

import sys
sys.path.append('../')

import cla
from cla.utils import get_signing_service, get_user_instance

user_id = '023643e7-1964-4fea-a877-fce7e4bc8eb5'
docusign_document_id = '7f06adaf-bfd9-4206-88bd-4561ccab2822'

user = get_user_instance()
user.load(user_id)
cla.log.info('Sending DocuSign document (%s) to user\'s email: %s', docusign_document_id, user.get_user_email())
get_signing_service().send_signed_document(docusign_document_id, user)
cla.log.info('Done')
