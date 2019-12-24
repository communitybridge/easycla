# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import json
import os
import requests
import uuid

import hug
import pytest
from falcon import HTTP_200, HTTP_409

import cla
from cla import routes

ID_TOKEN = os.environ.get('ID_TOKEN')
API_URL = os.environ.get('API_URL')

def test_create_company_duplicate():
    """
     Test creating duplicate company names
    """
    import pdb;pdb.set_trace()
    url = f'{API_URL}/v1/company'
    company_name = 'test_company_name'
    data = {
        'company_id' : uuid.uuid4() ,
        'company_name' : company_name ,
    }
    headers = {
        'Authorization' : f'Bearer {ID_TOKEN}'
    }
    response = requests.post(url, data=data, headers=headers)
    assert response.status == HTTP_200

    # add duplicate company
    data = {
        'company_id' : uuid.uuid4(),
        'company_name' : company_name
    }
    req = hug.test.post(routes, url, data=data, headers=headers)
    assert req.status == HTTP_409
