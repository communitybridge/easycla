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

AUTH0_DOMAIN = os.environ.get('AUTH0_DOMAIN')
AUTH0_AUDIENCE = os.environ.get('AUTH_AUDIENCE')
AUTH0_CLIENT_ID = os.environ.get('AUTH0_CLIENT_ID')
AUTH0_CLIENT_SECRET = os.environ.get('AUTH0_CLIENT_SECRET')
AUTH0_USERNAME_CLAIM = os.environ.get('AUTH0_USERNAME_CLAIM')


def acquire_access_token():
    """
    Utility function to create access token to be used in testing env
    """
    url = f'https://{AUTH0_DOMAIN}/oauth/token'
    headers = {'content-type ' : "application/json" }
    postdata = {
        'client_id' : AUTH0_CLIENT_ID,
        'client_secret' : AUTH0_CLIENT_SECRET,
        'audience' : AUTH0_AUDIENCE,
        'grant_type' : 'client_credentials'
    }
    response = requests.post(url, data=postdata, headers=headers)

    return response.json()


def test_create_company_duplicate():
    """
     Test creating duplicate company names
    """
    token = acquire_access_token()['access_token']
    url = '/v1/company'
    company_name = 'test_company_name'
    data = {
        'company_id' : uuid.uuid4() ,
        'company_name' : company_name ,
    }
    headers = {
        'Authorization' : f'Bearer {token}'
    }

    req = hug.test.post(routes, url, data=data, headers=headers)
    assert req.status == HTTP_200

    # add duplicate company
    data = {
        'company_id' : uuid.uuid4(),
        'company_name' : company_name
    }
    req = hug.test.post(routes, url, data=data, headers=headers)
    assert req.status == HTTP_409
