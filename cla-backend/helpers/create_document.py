# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

PROJECT_EXTERNAL_ID1 = 'a090t0000008DEiAAM'
PROJECT_EXTERNAL_ID2 = 'a090t0000008E7iAAE'
GITHUB_INSTALLATION_ID = 72228 # NOT THE APP ID - find it in the webhook request JSON or URL when viewing installed apps.

import sys
sys.path.append('../')

import cla
import uuid
import base64
import urllib.request
from cla.utils import get_document_instance, get_github_organization_instance, get_project_instance, get_pdf_service

from cla.resources.contract_templates import CNCFTemplate
template = CNCFTemplate(document_type='Corporate',
                 major_version=1,
                 minor_version=0)
individual_template = CNCFTemplate(document_type='Individual',
                 major_version=1,
                 minor_version=0)
content = template.get_html_contract("", "")
pdf_generator = get_pdf_service()
pdf_content = pdf_generator.generate(content)

# Organisation
github_org = get_github_organization_instance().get_organization_by_installation_id(GITHUB_INSTALLATION_ID)
# Project
github_project1 = get_project_instance().get_projects_by_external_id(PROJECT_EXTERNAL_ID1)[0]
github_project2 = get_project_instance().get_projects_by_external_id(PROJECT_EXTERNAL_ID2)[0]
# Document
# ICLA Project1
individual_document = get_document_instance()
individual_document.set_document_name('Test ICLA Document')
individual_document.set_document_file_id(str(uuid.uuid4()))
individual_document.set_document_content_type('storage+pdf')
individual_document.set_document_content(pdf_content, b64_encoded=False)
individual_document.set_document_major_version(1)
individual_document.set_document_minor_version(0)
individual_document.set_raw_document_tabs(template.get_tabs())
github_project1.add_project_individual_document(individual_document)
github_project2.add_project_individual_document(individual_document)
github_project1.save()
github_project2.save()

# CCLA
corporate_document = get_document_instance()
corporate_document.set_document_name('Test CCLA Document')
corporate_document.set_document_file_id(str(uuid.uuid4()))
corporate_document.set_document_content_type('storage+pdf')
corporate_document.set_document_content(pdf_content, b64_encoded=False)
corporate_document.set_document_major_version(1)
corporate_document.set_document_minor_version(0)
corporate_document.set_raw_document_tabs(template.get_tabs())
github_project1.add_project_corporate_document(corporate_document)
github_project2.add_project_corporate_document(corporate_document)
github_project1.save()
github_project2.save()
