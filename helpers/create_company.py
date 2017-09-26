"""
Convenience script to create companies.

Based on the database storage configuration found in config.py and config_instance.py.
"""

MANAGER_USER_ID = 'c75637ae-8aa8-4ad1-926c-15c7bb20915b'

import sys
sys.path.append('../')

import cla
import uuid
import base64
import urllib.request
from cla.utils import get_company_instance

cla.log.info('Creating new company with manager ID: %s', MANAGER_USER_ID)
company = get_company_instance()
company.set_company_id(str(uuid.uuid4()))
company.set_company_external_id('test-company')
company.set_company_manager_id(MANAGER_USER_ID)
company.set_company_name('Test Company')
company.set_company_whitelist([])
company.set_company_whitelist_patterns([])
#company.set_company_employees()
company.save()
cla.log.info('Done')
