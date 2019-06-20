# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

import cla
import requests
import os
import json
from http import HTTPStatus

import cla.auth
from cla.models.dynamo_models import UserPermissions

from oauthlib.oauth2 import LegacyApplicationClient
from requests_oauthlib import OAuth2Session
from jose import jwt

stage = os.environ.get('STAGE', '')

sf_instance_url = os.environ.get('SF_INSTANCE_URL', '')
sf_client_id = os.environ.get('SF_CLIENT_ID', '')
sf_client_secret = os.environ.get('SF_CLIENT_SECRET', '')
sf_username = os.environ.get('SF_USERNAME', '')
sf_password = os.environ.get('SF_PASSWORD', '')

cla_logo_url = os.environ.get('CLA_BUCKET_LOGO_URL', '')

def format_response(status_code, headers, body):
    """
    Helper function: Generic response formatter
    """
    response = {
        'statusCode': status_code,
        'headers': headers,
        'body': body
    }
    return response

def format_json_cors_response(status_code, body):
    """
    Helper function: Formats json responses with cors headers
    """
    body = json.dumps(body)
    cors_headers = {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Credentials': True,
        'Access-Control-Allow-Headers': 'Content-Type, Authorization'
    }
    response = format_response(status_code, cors_headers, body)
    return response

def get_sf_oauth_access():
    data = {
        'grant_type': 'password',
        'client_id': sf_client_id,
        'client_secret': sf_client_secret,
        'username': sf_username,
        'password': sf_password
    }

    token_url = 'https://{}/services/oauth2/token'.format(sf_instance_url)
    oauth_response = requests.post(token_url, data=data)
    oauth_response = oauth_response.json()

    return oauth_response

def get_projects(event, context):
    """
    Gets list of all projects from Salesforce
    """

    cla.log.info('event: {}'.format(event))

    try:
        auth_user = cla.auth.authenticate_user(event.get('headers'))
    except cla.auth.AuthError as e:
        cla.log.error('Authorization error: {}'.format(e))
        return format_json_cors_response(401, 'Error parsing Bearer token')
    except Exception as e:
        cla.log.error('Unknown authorization error: {}'.format(e))
        return format_json_cors_response(401, 'Error parsing Bearer token')

    # Get project access list for user
    user_permissions = UserPermissions()
    try:
        user_permissions.load(auth_user.username)
    except Exception as e:
        cla.log.error('Error invalid username: {}. error: {}'.format(auth_user.username, e))
        return format_json_cors_response(400, 'Error invalid username')

    user_permissions = user_permissions.to_dict()

    authorized_projects = user_permissions.get('projects')
    if authorized_projects is None:
        cla.log.error('Error user not authorized to access projects: {}'.format(user_permissions))
        return format_json_cors_response(403, 'Error user not authorized to access projects')

    project_list = ', '.join('\'' + project_id + '\'' for project_id in authorized_projects)

    oauth_response = get_sf_oauth_access()
    token = oauth_response['access_token']
    instance_url = oauth_response['instance_url']

    headers = {
        'Authorization': 'Bearer {}'.format(token),
        'Content-Type': 'application/json',
    }

    query_url = '{}/services/data/v20.0/query/'.format(instance_url)
    query = {'q': 'SELECT id, Name, Description__c from Project__c WHERE id IN ({})'.format(project_list)}
    r = requests.get(query_url, headers=headers, params=query)

    response = r.json()
    status_code = r.status_code
    if status_code != HTTPStatus.OK:
        cla.log.error('Error retrieving projects: %s', response[0].get('message'))
        return format_json_cors_response(status_code, 'Error retrieving projects')
    records = response.get('records')

    projects = []
    for project in records:
        logo_url = None
        project_id = project.get('Id')
        if project_id:
            logo_url = '{}/{}.png'.format(cla_logo_url, project_id)

        projects.append({
            'name': project.get('Name'),
            'id': project_id,
            'description': project.get('Description__c'),
            'logoUrl': logo_url
        })

    return format_json_cors_response(status_code, projects)

def get_project(event, context):
    """
    Given project id, gets project details from Salesforce
    """

    cla.log.info('event: {}'.format(event))

    project_id = event.get('queryStringParameters').get('id')
    if project_id is None:
        return format_json_cors_response(400, 'Missing project ID')

    try:
        auth_user = cla.auth.authenticate_user(event.get('headers'))
    except cla.auth.AuthError as e:
        cla.log.error('Authorization error: {}'.format(e))
        return format_json_cors_response(401, 'Error parsing Bearer token')
    except Exception as e:
        cla.log.error('Unknown authorization error: {}'.format(e))
        return format_json_cors_response(401, 'Error parsing Bearer token')

    # Get project access list for user
    user_permissions = UserPermissions()
    try:
        user_permissions.load(auth_user.username)
    except:
        cla.log.error('Error invalid username: {}'.format(auth_user.username))
        return format_json_cors_response(400, 'Error invalid username')

    user_permissions = user_permissions.to_dict()

    authorized_projects = user_permissions.get('projects')
    if authorized_projects is None:
        cla.log.error('Error user not authorized to access projects: {}'.format(user_permissions))
        return format_json_cors_response(403, 'Error user not authorized to access projects')
    
    if project_id not in authorized_projects:
        cla.log.error('Error user not authorized')
        return format_json_cors_response(403, 'Error user not authorized')

    oauth_response = get_sf_oauth_access()
    token = oauth_response['access_token']
    instance_url = oauth_response['instance_url']

    headers = {
        'Authorization': 'Bearer {}'.format(token)
    }

    url = '{}/services/data/v20.0/sobjects/Project__c/{}'.format(instance_url, project_id)
    cla.log.info('Calling salesforce api for project info..')
    r = requests.get(url, headers=headers)

    response = r.json()
    status_code = r.status_code
    if status_code != HTTPStatus.OK:
        cla.log.error('Error retrieving project: %s', response[0].get('message'))
        return format_json_cors_response(status_code, 'Error retrieving project')

    logo_url = None
    if response.get('id'):
        logo_url = '{}/{}.png'.format(cla_logo_url, response.get('id'))

    project = {
        'name': response.get('Name'),
        'id': response.get('Id'),
        'description': response.get('Description__c'),
        'logoUrl': logo_url
    }
    return format_json_cors_response(status_code, project)
