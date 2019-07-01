# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Convenience script to send a DocuSign document to the user's email address.

The user_id specified must exist in the database.
The docusign_document_id specified must exist in DocuSign.
"""

import sys
sys.path.append('../')

import cla
from cla.utils import get_signing_service, get_user_instance

docusign_document_id = 'dcd9a52b-bed8-4c6f-9a71-ce00252a4e5d'

user = get_user_instance().get_user_by_github_id(123)
cla.log.info('Sending DocuSign document (%s) to user\'s email: %s', docusign_document_id, user.get_user_email())
get_signing_service().send_signed_document(docusign_document_id, user)
cla.log.info('Done')
