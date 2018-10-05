import cla
import requests
import os
import json
from http import HTTPStatus
from oauthlib.oauth2 import LegacyApplicationClient
from requests_oauthlib import OAuth2Session

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
    token_url = 'https://{}/services/oauth2/token'.format(sf_instance_url)

    print('token_url: {}\nclient_secret: {}\nclient_id: {}\nusername: {}\npassword: {}'.format(token_url, sf_client_secret, sf_client_id, sf_username, sf_password))

    oauth2 = OAuth2Session(client=LegacyApplicationClient(client_id=sf_client_id))
    token = oauth2.fetch_token(token_url=token_url, client_secret=sf_client_secret,
    client_id=sf_client_id, username=sf_username, password=sf_password)

    query = {'q': 'SELECT name from Project__c'}
    headers = {'Content-Type': 'application/json'}
    url = 'https://{}/services/data/v20.0/query/'.format(sf_instance_url)
    cla.log.info('Calling salesforce api for project list...')
    r = oauth2.request('GET', url, params=query, headers=headers)

    response = r.json()
    status_code = r.status_code
    if status_code != HTTPStatus.OK:
        cla.log.error('Error retrieving projects: %s', response[0].get('message'))
        return format_json_cors_response(status_code, 'Error retrieving projects')
    records = response.get('records')
    projects = [
        {'name': project.get('Name'),
        'id': project.get('attributes').get('url').split('/')[-1],
        'logoRef': project.get('Image_File_for_PDF__c')}
        for project in records]
    return format_json_cors_response(status_code, projects)


def get_project(event, context):
    """
    Given project id, gets project details from Salesforce
    """
    project_id = event.get('queryStringParameters').get('id')

    token_url = 'https://{}/services/oauth2/token'.format(sf_instance_url)
    oauth2 = OAuth2Session(client=LegacyApplicationClient(client_id=sf_client_id))
    token = oauth2.fetch_token(token_url=token_url, client_secret=sf_client_secret,
    client_id=sf_client_id, username=sf_username, password=sf_password)

    url = 'https://{}/services/data/v20.0/sobjects/Project__c/{}'.format(sf_instance_url, project_id)
    cla.log.info('Calling salesforce api for project info..')
    r = oauth2.request('GET', url)

    response = r.json()
    status_code = r.status_code
    if status_code != HTTPStatus.OK:
        cla.log.error('Error retrieving project: %s', response[0].get('message'))
        return format_json_cors_response(status_code, 'Error retrieving project')
    project = {
        'name': response.get('Name'),
        'id': response.get('Id'),
        'logoRef': response.get('Image_File_for_PDF__c'),
        'description': response.get('Description__c')
    }
    return format_json_cors_response(status_code, project)
