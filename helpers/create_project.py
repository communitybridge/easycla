# Project external ID.
PROJECT_EXTERNAL_ID = 'salesforce-id-here'

import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_project_instance

## Project
cla.log.info('Creating new project for %s', PROJECT_EXTERNAL_ID)
github_project = get_project_instance()
github_project.set_project_id(str(uuid.uuid4()))
github_project.set_project_external_id(PROJECT_EXTERNAL_ID)
github_project.set_project_name('Test Project')
github_project.set_project_ccla_requires_icla_signature(False)
github_project.save()
