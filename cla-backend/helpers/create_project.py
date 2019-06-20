# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

# Project external ID.
PROJECT_EXTERNAL_ID1 = 'a090t0000008DEiAAM'
PROJECT_EXTERNAL_ID2 = 'a090t0000008E7iAAE'

import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_project_instance

## Project
cla.log.info('Creating first project for %s', PROJECT_EXTERNAL_ID1)
github_project = get_project_instance()
github_project.set_project_id(str(uuid.uuid4()))
github_project.set_project_external_id(PROJECT_EXTERNAL_ID1)
github_project.set_project_name('Test Project One')
github_project.set_project_ccla_requires_icla_signature(False)
github_project.save()
cla.log.info('Creating second project for %s', PROJECT_EXTERNAL_ID2)
github_project = get_project_instance()
github_project.set_project_id(str(uuid.uuid4()))
github_project.set_project_external_id(PROJECT_EXTERNAL_ID2)
github_project.set_project_name('Test Project Two')
github_project.set_project_ccla_requires_icla_signature(True)
github_project.save()
