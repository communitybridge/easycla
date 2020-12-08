# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import json
import os
from urllib.parse import quote

import requests
from boto3 import client

from cla import log

STAGE = os.environ.get('STAGE', '')
REGION = 'us-east-1'


class UserService:
    """
    UserService Handles external salesforce Users
    """

    def __init__(self):
        self.platform_gateway_url = self.get_ssm_key(REGION, f'cla-auth0-platform-api-gw-{STAGE}')

    def get_ssm_key(self, region, key):
        """
        Fetches the specified SSM key value from the SSM key store
        :param region: aws region
        :type region: string
        :parm key: key
        :type key: string
        """
        ssm_client = client('ssm', region_name=region)

        log.debug(f'Fetching Key: {key}')
        response = ssm_client.get_parameter(Name=key, WithDecryption=True)
        log.debug(f'Fetched Key: {key}, value: {response["Parameter"]["Value"]}')
        return response['Parameter']['Value']

    def get_user_by_sf_id(self, sf_user_id: str):
        """
        Queries the platform user service for the specified user id. The
        result will return all the details for the user as a dictionary.
        """
        headers = {
            'Authorization': f'bearer {self.get_access_token()}',
            'accept': 'application/json'
        }

        try:
            url = f'{self.platform_gateway_url}/user-service/v1/users/{sf_user_id}'
            log.debug(f'Sending GET request to {url}')
            r = requests.get(url, headers=headers)
            r.raise_for_status()
            response_model = json.loads(r.text)
            return response_model
        except requests.exceptions.HTTPError as err:
            msg = f'Could not get user: {sf_user_id}, error: {err}'
            log.warning(msg)
            return None

    def _get_users_by_key_value(self, key: str, value: str ):
        """
        Queries the platform user service for the specified criteria.
        The result will return summary information for the users as a
        dictionary.
        """
        headers = {
            'Authorization': f'bearer {self.get_access_token()}',
            'accept': 'application/json'
        }

        users = []
        offset = 0
        pagesize = 1000

        while True:
            try:
                log.info(f'Search User using key: {key} with value: {value}')
                url = f'{self.platform_gateway_url}/user-service/v1/users/search?' \
                      f'{key}={quote(value)}&pageSize={pagesize}&offset={offset}'
                log.debug(f'Sending GET request to {url}')
                r = requests.get(url, headers=headers)
                r.raise_for_status()
                response_model = json.loads(r.text)
                total = response_model['Metadata']['TotalSize']
                if response_model['Data']:
                    users = users + response_model['Data']
                if total < (pagesize + offset):
                    break
                offset = offset + pagesize
            except requests.exceptions.HTTPError as err:
                msg = f'Could not get projects, error: {err}'
                log.warning(msg)
                return None

        log.debug('total users : {}'.format(len(users)))
        return users

    def get_users_by_username(self, user_name: str):
        return self._get_users_by_key_value("username", user_name)

    def get_users_by_firstname(self, first_name: str):
        return self._get_users_by_key_value("firstname", first_name)

    def get_users_by_lastname(self, last_name: str):
        return self._get_users_by_key_value("lastname", last_name)

    def get_users_by_email(self, email: str):
        return self._get_users_by_key_value("email", email)

    def get_access_token(self):
        auth0_url = self.get_ssm_key(REGION, f'cla-auth0-platform-url-{STAGE}')
        platform_client_id = self.get_ssm_key(REGION, f'cla-auth0-platform-client-id-{STAGE}')
        platform_client_secret = self.get_ssm_key(REGION, f'cla-auth0-platform-client-secret-{STAGE}')
        platform_audience = self.get_ssm_key(REGION, f'cla-auth0-platform-audience-{STAGE}')

        auth0_payload = {
            'grant_type': 'client_credentials',
            'client_id': platform_client_id,
            'client_secret': platform_client_secret,
            'audience': platform_audience
        }

        headers = {
            'content-type': 'application/x-www-form-urlencoded',
            'accept': 'application/json'
        }

        access_token = ''
        try:
            # logger.debug(f'Sending POST to {auth0_url} with payload: {auth0_payload}')
            log.debug(f'Sending POST to {auth0_url}')
            r = requests.post(auth0_url, data=auth0_payload, headers=headers)
            r.raise_for_status()
            json_data = json.loads(r.text)
            access_token = json_data["access_token"]
            return access_token
        except requests.exceptions.HTTPError as err:
            log.warning(f'Could not get auth token, error: {err}')
            return None
