# Project external ID.
PROJECT_EXTERNAL_ID = 'salesforce-id-here'
# The GitHub user/org used for testing purposes.
GITHUB_ORGANIZATION_NAME = 'linuxfoundation'
GITHUB_INSTALLATION_ID = 49309 # NOT THE APP ID - find it in the webhook request JSON or URL when viewing installed apps.

import sys
sys.path.append('../')

import cla
from cla.utils import get_project_instance, get_github_organization_instance

# Organisation
project = get_project_instance().get_projects_by_external_id(PROJECT_EXTERNAL_ID)[0]
cla.log.info('Creating GitHub Organization: %s' %GITHUB_ORGANIZATION_NAME)
github_org = get_github_organization_instance()
github_org.set_organization_name(GITHUB_ORGANIZATION_NAME)
github_org.set_organization_project_id(project.get_project_id())
github_org.set_organization_installation_id(GITHUB_INSTALLATION_ID)
github_org.save()
