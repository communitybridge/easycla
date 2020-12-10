# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import datetime
import json
import os
from urllib.parse import quote

import requests

import cla
from cla import log

STAGE = os.environ.get('STAGE', '')
REGION = 'us-east-1'


class UserServiceInstance:
    """
    UserService Handles external salesforce Users
    """

    access_token = None
    access_token_expires = datetime.datetime.now() + datetime.timedelta(minutes=30)

    def __init__(self):
        self.platform_gateway_url = cla.config.PLATFORM_GATEWAY_URL

    def get_user_by_sf_id(self, sf_user_id: str):
        """
        Queries the platform user service for the specified user id. The
        result will return all the details for the user as a dictionary.
        """
        fn = 'user_service.get_user_by_sf_id'

        headers = {
            'Authorization': f'bearer {self.get_access_token()}',
            'accept': 'application/json'
        }

        try:
            url = f'{self.platform_gateway_url}/user-service/v1/users/{sf_user_id}'
            log.debug(f'{fn} - sending GET request to {url}')
            r = requests.get(url, headers=headers)
            r.raise_for_status()
            response_model = json.loads(r.text)
            return response_model
        except requests.exceptions.HTTPError as err:
            msg = f'{fn} - Could not get user: {sf_user_id}, error: {err}'
            log.warning(msg)
            return None

    def _get_users_by_key_value(self, key: str, value: str):
        """
        Queries the platform user service for the specified criteria.
        The result will return summary information for the users as a
        dictionary.
        """
        fn = 'user_service._get_users_by_key_value'

        headers = {
            'Authorization': f'bearer {self.get_access_token()}',
            'accept': 'application/json'
        }

        users = []
        offset = 0
        pagesize = 1000

        while True:
            try:
                log.info(f'{fn} - Search User using key: {key} with value: {value}')
                url = f'{self.platform_gateway_url}/user-service/v1/users/search?' \
                      f'{key}={quote(value)}&pageSize={pagesize}&offset={offset}'
                log.debug(f'{fn} - sending GET request to {url}')
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
                log.warning(f'{fn} - Could not get projects, error: {err}')
                return None

        log.debug(f'{fn} - total users : {len(users)}')
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
        fn = 'user_service.get_access_token'
        # Use previously cached value, if not expired
        if self.access_token and datetime.datetime.now() < self.access_token_expires:
            cla.log.debug(f'{fn} - using cached access token: {self.access_token[0:10]}...')
            return self.access_token

        auth0_url = cla.config.AUTH0_PLATFORM_URL
        platform_client_id = cla.config.AUTH0_PLATFORM_CLIENT_ID
        platform_client_secret = cla.config.AUTH0_PLATFORM_CLIENT_SECRET
        platform_audience = cla.config.AUTH0_PLATFORM_AUDIENCE

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

        try:
            # logger.debug(f'Sending POST to {auth0_url} with payload: {auth0_payload}')
            log.debug(f'{fn} - sending POST to {auth0_url}')
            r = requests.post(auth0_url, data=auth0_payload, headers=headers)
            r.raise_for_status()
            json_data = json.loads(r.text)
            self.access_token = json_data["access_token"]
            self.access_token_expires = datetime.datetime.now() + datetime.timedelta(minutes=30)
            log.debug(f'{fn} - successfully obtained access_token: {self.access_token[0:10]}...')
            return self.access_token
        except requests.exceptions.HTTPError as err:
            log.warning(f'{fn} - could not get auth token, error: {err}')
            return None


UserService = UserServiceInstance()
