import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_project_instance, get_user_instance, get_document_instance, get_github_organization_instance, get_project_instance, get_pdf_service, get_company_instance, get_signature_instance
from cla.resources.contract_templates import CNCFTemplate


# Project external ID.
PROJECT_EXTERNAL_ID = 'a0941000002wByZAAU'

## Project
cla.log.info('Creating first project for %s', PROJECT_EXTERNAL_ID)
project = get_project_instance()
project.set_project_id(str(uuid.uuid4()))
project.set_project_external_id(PROJECT_EXTERNAL_ID)
project.set_project_name('Test Project One')
project.set_project_ccla_requires_icla_signature(False)
project.save()

## Create CCLA Document
corporate_template = CNCFTemplate(document_type='Corporate',
                 major_version=1,
                 minor_version=0)
content = corporate_template.get_html_contract("", "")
pdf_generator = get_pdf_service()
pdf_content = pdf_generator.generate(content)


# CCLA
corporate_document = get_document_instance()
corporate_document.set_document_name('Test CCLA Document')
corporate_document.set_document_file_id(str(uuid.uuid4()))
corporate_document.set_document_content_type('storage+pdf')
corporate_document.set_document_content(pdf_content, b64_encoded=False)
corporate_document.set_document_major_version(1)
corporate_document.set_document_minor_version(0)
corporate_document.set_raw_document_tabs(corporate_template.get_tabs())
project.add_project_corporate_document(corporate_document)
project.save()

## Create Github Org

## User (For Company Management)
cla.log.info('Creating company manager user')
manager = get_user_instance()
manager.set_user_id(str(uuid.uuid4()))
manager.set_user_name('First User')
manager.set_user_email('firstuser@domain.org')
manager.set_user_email('***REMOVED***@linuxfoundation.org')
manager.set_user_github_id(123)
manager.save()

## Company
cla.log.info('Creating new company with manager ID: %s', manager.get_user_id())
company = get_company_instance()
company.set_company_id(str(uuid.uuid4()))
company.set_company_external_id('company-external-id')
company.set_company_manager_id(manager.get_user_id())
company.set_company_name('Test Company')
company.set_company_whitelist([])
company.set_company_whitelist_patterns(['*@listed.org'])
company.save()

## Company Signature
company_signature_id = str(uuid.uuid4())
cla.log.info('Creating CCLA signature for company %s and project %s: %s' \
             %(company.get_company_external_id(), project.get_project_external_id(), company_signature_id))
signature = get_signature_instance()
signature.set_signature_id(company_signature_id)
signature.set_signature_project_id(project.get_project_id())
signature.set_signature_signed(True)
signature.set_signature_approved(True)
signature.set_signature_type('cla')
signature.set_signature_reference_id(company.get_company_id())
signature.set_signature_reference_type('company')
signature.set_signature_document_major_version(1)
signature.set_signature_document_minor_version(0)
signature.save()

## User (For Individual Contributor)

## Signature (Individual Contributor)
