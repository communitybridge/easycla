from http import HTTPStatus
import requests
import os
import json
from requests_oauthlib import OAuth2Session

sf_instance_url = os.environ.get('SF_INSTANCE_URL', '')

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
        'Access-Control-Allow-Credentials': True
    }
    response = format_response(status_code, cors_headers, body)
    return response


def get_projects():
    query = {'q': 'SELECT+name+from+Project__c'}
    url = 'https://{}/services/data/v20.0/query/'.format(sf_instance_url)
    ##Oauth stuff
    token = '' #need to get using fetch token thing
    oauth2 = OAuth2Session(client_id, token=token)
    #client_id='3MVG9Xjf0O2Peyd6s.Y3iz9Ev_DVUX_IElCRz4Yhi1Itv5WPAdfCx_KTVcjYxBTEg_PRfdkoApZ4rvFqV8KbE',
    oauth2.fetch_token(code=None, token_url='https://cs93.salesforce.com/services/oauth2/token',
    client_secret='7907730997014943382', username='ernest@twobulls.com', 
    password='GB{9TB9B}8CDPv2oz&naca[JmRTH04K6rxIrM9WfqrOs9mO0Z', grant_type='password')
    request = oauth2.get(url)
    
    ##fix
    response = requests.get(url, params=query)
    status_code = response.status_code
    if status_code != HTTPStatus.OK:
        return format_json_cors_response(status_code, 'Error retrieving projects')
    records = response.json()
    projects = [
        {'name': project.get('Name'), 
        'id': project.get('attributes').get('url').split('/')[-1]}
        for project in records]
    return format_json_cors_response(status_code, projects)

