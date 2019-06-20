# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Convenience script to create a new user signature request (simulate a user clicking on the sign icon in GitHub).
"""
import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_user_instance, get_project_instance, set_active_signature_metadata

PROJECT_EXTERNAL_ID1 = 'a090t0000008DEiAAM'

# Create new user so as to not conflict with the create_user.py script.
user = get_user_instance()
user.set_user_id(str(uuid.uuid4()))
user.set_user_name('Signing User')
user.set_user_email('signing@domain.org')
user.set_user_github_id(234)
user.save()

user_id = user.get_user_id()
project = get_project_instance().get_projects_by_external_id(PROJECT_EXTERNAL_ID1)[0]
project_id = project.get_project_id()

repository_id = '96820382'
pull_request_id = '16'
url = 'https://github.com/linuxfoundation/CLA-Test/pull/16'

cla.log.info('Creating new active signature for project %s, user %s, repository %s, PR %s - %s',
             project_id, user_id, repository_id, pull_request_id, url)

# Store data on signature.
set_active_signature_metadata(user_id, project_id, repository_id, pull_request_id)
