import cla
import requests
import os
import json
from http import HTTPStatus

import cla.utils

from oauthlib.oauth2 import LegacyApplicationClient
from requests_oauthlib import OAuth2Session
from jose import jwt

stage = os.environ.get('STAGE', '')

sf_instance_url = os.environ.get('SF_INSTANCE_URL', '')
sf_client_id = os.environ.get('SF_CLIENT_ID', '')
sf_client_secret = os.environ.get('SF_CLIENT_SECRET', '')
sf_username = os.environ.get('SF_USERNAME', '')
sf_password = os.environ.get('SF_PASSWORD', '')

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

def get_projects(event, context):
    """
    Gets list of all projects from Salesforce
    """

    cla.log.info('event: {}'.format(event))

    # Get userID from token
    if stage == 'dev':
        headers = event.headers
        bearer_token = headers.get('AUTHORIZATION')
    else:
        headers = event.get('headers')
        bearer_token = headers.get('Authorization')

    if headers is None:
        cla.log.error('Error reading headers')
        return format_json_cors_response(400, 'Error reading headers')

    if bearer_token is None:
        cla.log.error('Error reading authorization header')
        return format_json_cors_response(400, 'Error reading authorization header')
    
    bearer_token = bearer_token.replace('Bearer ', '')
    try:
        token_params = jwt.get_unverified_claims(bearer_token)
    except Exception as e:
        cla.log.error('Error parsing Bearer token: {}'.format(e))
        return format_json_cors_response(400, 'Error parsing Bearer token')

    user_id = token_params.get('sub')
    if user_id is None:
        cla.log.error('Error parsing user ID. event')
        return format_json_cors_response(400, 'Error parsing user ID')

    # Get project access list for user
    user_permissions = cla.utils.get_user_permissions_instance()
    try:
        user_permissions.load(user_id)
    except Exception as e:
        cla.log.error('Error invalid user ID: {}. error: {}'.format(user_id, e))
        return format_json_cors_response(400, 'Error invalid user ID')

    user_permissions = user_permissions.to_dict()

    authorized_projects = user_permissions.get('projects')
    if authorized_projects is None:
        cla.log.error('Error user not authorized to access projects: {}'.format(user_permissions))
        return format_json_cors_response(403, 'Error user not authorized to access projects')

    token_url = 'https://{}/services/oauth2/token'.format(sf_instance_url)

    oauth2 = OAuth2Session(client=LegacyApplicationClient(client_id=sf_client_id))
    token = oauth2.fetch_token(token_url=token_url, client_secret=sf_client_secret,
    client_id=sf_client_id, username=sf_username, password=sf_password)

    project_list = ', '.join('\'' + project_id + '\'' for project_id in authorized_projects)
    query = {'q': 'SELECT id, Name, Description__c from Project__c WHERE id IN ({})'.format(project_list)}
    headers = {'Content-Type': 'application/json'}
    url = '{}/services/data/v20.0/query/'.format(token['instance_url'])
    cla.log.info('Calling salesforce api for project list...')
    r = oauth2.request('GET', url, params=query, headers=headers)

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
            logo_url = 'https://s3.amazonaws.com/cla-logo-{}/{}.png'.format(stage, project_id)

        projects.append({
            'name': project.get('Name'),
            'id': project_id,
            'description': project.get('Description'),
            'logoUrl': logo_url
        })

    return format_json_cors_response(status_code, projects)

def get_project(event, context):
    """
    Given project id, gets project details from Salesforce
    """

    cla.log.info('event: {}'.format(event))

    # Get userID from token
    if stage == 'dev':
        headers = event.headers
        bearer_token = headers.get('AUTHORIZATION')
        project_id = event.params
        cla.log.info(project_id)
    else:
        headers = event.get('headers')
        bearer_token = headers.get('Authorization')
        project_id = event.get('queryStringParameters').get('id')

    if headers is None:
        cla.log.error('Error reading headers. event')
        return format_json_cors_response(400, 'Error reading headers')

    if bearer_token is None:
        cla.log.error('Error reading authorization header')
        return format_json_cors_response(400, 'Error reading authorization header')

    bearer_token = bearer_token.replace('Bearer ', '')
    try:
        token_params = jwt.get_unverified_claims(bearer_token)
    except:
        cla.log.error('Error parsing Bearer token')
        return format_json_cors_response(400, 'Error parsing Bearer token')

    user_id = token_params.get('sub')
    if user_id is None:
        cla.log.error('Error parsing user ID')
        return format_json_cors_response(400, 'Error parsing user ID')

    # Get project access list for user
    user_permissions = cla.utils.get_user_permissions_instance()
    try:
        user_permissions.load(user_id)
    except:
        cla.log.error('Error invalid user ID: {}'.format(user_id))
        return format_json_cors_response(400, 'Error invalid user ID')

    user_permissions = user_permissions.to_dict()

    authorized_projects = user_permissions.get('projects')
    if authorized_projects is None:
        cla.log.error('Error user not authorized to access projects: {}'.format(user_permissions))
        return format_json_cors_response(403, 'Error user not authorized to access projects')

    
    if project_id not in authorized_projects:
        cla.log.error('Error user not authorized')
        return format_json_cors_response(403, 'Error user not authorized')

    token_url = 'https://{}/services/oauth2/token'.format(sf_instance_url)
    oauth2 = OAuth2Session(client=LegacyApplicationClient(client_id=sf_client_id))
    token = oauth2.fetch_token(token_url=token_url, client_secret=sf_client_secret,
    client_id=sf_client_id, username=sf_username, password=sf_password)

    url = '{}/services/data/v20.0/sobjects/Project__c/{}'.format(token['instance_url'], project_id)
    cla.log.info('Calling salesforce api for project info..')
    r = oauth2.request('GET', url)

    response = r.json()
    status_code = r.status_code
    if status_code != HTTPStatus.OK:
        cla.log.error('Error retrieving project: %s', response[0].get('message'))
        return format_json_cors_response(status_code, 'Error retrieving project')

    logo_url = None
    if response.get('id'):
        logo_url = 'https://s3.amazonaws.com/cla-logo-{}/{}.png'.format(stage, response.get('id'))

    project = {
        'name': response.get('Name'),
        'id': response.get('Id'),
        'description': response.get('Description__c'),
        'logoUrl': logo_url
    }
    return format_json_cors_response(status_code, project)
